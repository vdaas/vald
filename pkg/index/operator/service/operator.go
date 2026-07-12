// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/reconciler"
	"github.com/vdaas/vald/internal/k8s/vald"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

const (
	apiName = "vald/index/operator"
)

type jobReconcileResult int

const (
	createRequired jobReconcileResult = iota
	createSkipped
	requeueRequired
)

// Operator represents an interface for indexing.
type Operator interface {
	Start(ctx context.Context) (<-chan error, error)
}

type operator struct {
	ctrl                              k8s.Controller
	eg                                errgroup.Group
	client                            k8s.ClientWithWatch
	deployments                       client.Client[*k8s.Deployment, *k8s.DeploymentList]
	jobs                              client.Client[*k8s.Job, *k8s.JobList]
	rotatorJob                        *k8s.Job
	namespace                         string
	rotatorName                       string
	targetReadReplicaIDAnnotationsKey string
	readReplicaLabelKey               string
	rotationJobConcurrency            uint
	readReplicaEnabled                bool
}

// New returns Indexer object if no error occurs.
func New(
	namespace, agentName, rotatorName, targetReadReplicaIDKey string,
	rotatorJob *k8s.Job,
	opts ...Option,
) (o Operator, err error) {
	operator := new(operator)
	operator.namespace = namespace
	operator.targetReadReplicaIDAnnotationsKey = targetReadReplicaIDKey
	operator.rotatorName = rotatorName
	operator.rotatorJob = rotatorJob
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(operator); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(err)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}

	isAgent := func(pod *k8s.Pod) bool {
		return pod.Labels["app"] == agentName
	}

	podController := reconciler.NewObjectReconciler[*k8s.Pod](
		"pod reconciler for index operator", new(k8s.Pod),
		reconciler.WithOnObjectReconcile(operator.podOnReconcile),
		reconciler.WithObjectOnError[*k8s.Pod](func(err error) {
			log.Error("failed to reconcile:", err)
		}),
		// To only reconcile for agent pods
		reconciler.WithObjectForOptions[*k8s.Pod](
			client.PodPredicates(isAgent),
		),
	)

	operator.ctrl, err = k8s.New(
		k8s.WithResourceController(podController),
		k8s.WithLeaderElection(true, "vald-index-operator", operator.namespace),
	)
	if err != nil {
		return nil, err
	}

	if operator.client == nil {
		c, err := client.New[*k8s.Pod, *k8s.PodList](new(k8s.Pod), new(k8s.PodList))
		if err != nil {
			return nil, err
		}
		operator.client = c.Raw()
	}
	operator.deployments = client.NewWithClient[*k8s.Deployment, *k8s.DeploymentList](operator.client, new(k8s.Deployment), new(k8s.DeploymentList))
	operator.jobs = client.NewWithClient[*k8s.Job, *k8s.JobList](operator.client, new(k8s.Job), new(k8s.JobList))

	return operator, nil
}

// Start starts indexing process.
func (o *operator) Start(ctx context.Context) (<-chan error, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/service/operator.Start")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	dech, err := o.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 2)
	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-dech:
				if err != nil {
					ech <- err
				}
			}
		}
	}))

	return ech, nil
}

func (o *operator) podOnReconcile(ctx context.Context, pod *k8s.Pod) (k8s.Result, error) {
	if o.readReplicaEnabled {
		rq, err := o.reconcileRotatorJob(ctx, pod)
		if err != nil {
			log.Errorf("reconciling rotator job: %s", err)
			return k8s.Result{}, fmt.Errorf("reconciling rotator job: %w", err)
		}
		// let controller-runtime backoff exponentially by not setting the backoff duration
		return k8s.Result{
			Requeue: rq,
		}, nil
	}

	return k8s.Result{}, nil
}

// reconcileRotatorJob starts rotation job when the condition meets.
// This function is work in progress.
func (o *operator) reconcileRotatorJob(
	ctx context.Context, pod *k8s.Pod,
) (requeue bool, err error) {
	podIdx, ok := pod.Labels[k8s.PodIndexLabel]
	if !ok {
		log.Info("no index label found. the agent is not StatefulSet? skipping...")
		return false, nil
	}

	// retrieve the readreplica deployment annotations for podIdx
	selector, err := client.NewLabelSelector(o.readReplicaLabelKey, k8s.SelectionOpEquals, []string{podIdx})
	if err != nil {
		return false, fmt.Errorf("creating label selector: %w", err)
	}
	readReplicaDeployments, err := o.deployments.List(ctx, &k8s.ListOptions{
		Namespace:     o.namespace,
		LabelSelector: selector,
	})
	if err != nil {
		return false, err
	}
	if len(readReplicaDeployments.Items) == 0 {
		return false, errors.New("no readreplica deployment found")
	}
	dep := readReplicaDeployments.Items[0]

	need, err := needsRotation(pod.Annotations, dep.Annotations)
	if err != nil {
		return false, fmt.Errorf("checking if rotation is required: %w", err)
	}
	if !need {
		return false, nil
	}

	log.Infof("rotation required for agent(id: %s)", podIdx)
	requeue, err = o.createRotationJobOrRequeue(ctx, podIdx)
	if err != nil {
		return false, fmt.Errorf("creating rotation job: %w", err)
	}
	return requeue, nil
}

func needsRotation(agentAnnotations, readReplicaAnnotations map[string]string) (bool, error) {
	t, ok := agentAnnotations[vald.LastTimeSaveIndexTimestampAnnotationsKey]
	if !ok {
		log.Info("the agent pod has not saved index yet. skipping...")
		return false, nil
	}
	lastSavedTime, err := time.Parse(vald.TimeFormat, t)
	if err != nil {
		return false, fmt.Errorf("parsing last time saved time: %w", err)
	}

	t, ok = readReplicaAnnotations[vald.LastTimeSnapshotTimestampAnnotationsKey]
	if ok {
		lastSnapshotTime, err := time.Parse(vald.TimeFormat, t)
		if err != nil {
			return false, fmt.Errorf("parsing last snapshot time: %w", err)
		}

		if lastSnapshotTime.After(lastSavedTime) {
			log.Info("snapshot taken after the last save. skipping...")
			return false, nil
		}
	}

	return true, nil
}

func (o *operator) createRotationJobOrRequeue(
	ctx context.Context, podIdx string,
) (rq bool, err error) {
	// get all the rotation jobs and make sure the job is not running
	res, err := o.ensureJobConcurrency(ctx, podIdx)
	if err != nil {
		return false, fmt.Errorf("checking if the same job exists: %w", err)
	}
	switch res {
	case createSkipped:
		log.Infof("rotation job for the agent(id: %s) is already running. skipping...", podIdx)
		return false, nil
	case requeueRequired:
		log.Infof("rotation job concurrency limit(%d) reached. rotation job for the agent(id: %s) will be requeued", o.rotationJobConcurrency, podIdx)
		return true, nil
	case createRequired:
		// continue to create a new job
		break
	}

	// now we actually need to create the rotator job
	log.Infof("no job is running to rotate the agent(id:%s). creating a new job...", podIdx)
	job := o.rotatorJob.DeepCopy()
	if job.Spec.Template.Annotations == nil {
		job.Spec.Template.Annotations = make(map[string]string)
	}
	job.Spec.Template.Annotations[o.targetReadReplicaIDAnnotationsKey] = podIdx
	job.ObjectMeta = k8s.ObjectMeta{
		GenerateName: fmt.Sprintf("%s-", o.rotatorName),
		Namespace:    o.namespace,
	}

	if err := o.jobs.Create(ctx, job); err != nil {
		return false, fmt.Errorf("creating job resource with k8s API: %w", err)
	}

	return false, nil
}

// ensureJobConcurrency controls the job concurrency. It cannot handle concurrent calls but it is fine because
// the MaxConcurrentReconciles defaults to 1 and we do not change it.
func (o *operator) ensureJobConcurrency(
	ctx context.Context, podIdx string,
) (jobReconcileResult, error) {
	// get all the rotation jobs and make sure the job is not running
	selector, err := client.NewLabelSelector("app", k8s.SelectionOpEquals, []string{o.rotatorName})
	if err != nil {
		return createSkipped, fmt.Errorf("creating label selector: %w", err)
	}
	jobList, err := o.jobs.List(ctx, &k8s.ListOptions{
		Namespace:     o.namespace,
		LabelSelector: selector,
	})
	if err != nil {
		return createSkipped, fmt.Errorf("listing jobs: %w", err)
	}

	// no need to check finished jobs
	jobList.Items = slices.DeleteFunc(jobList.Items, func(job k8s.Job) bool {
		return job.Status.Active == 0
	})

	if len(jobList.Items) >= int(o.rotationJobConcurrency) {
		return requeueRequired, nil
	}

	for _, job := range jobList.Items {
		annotations := job.Spec.Template.Annotations
		if annotations == nil {
			continue
		}
		id, ok := annotations[o.targetReadReplicaIDAnnotationsKey]
		if !ok {
			continue
		}
		if id == podIdx {
			// the same job is already running. no need to requeue
			return createSkipped, nil
		}
	}

	return createRequired, nil
}
