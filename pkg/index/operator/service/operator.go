// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/k8s/vald"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"

	//FIXME:
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	apiName = "vald/index/operator"
)

// Operator represents an interface for indexing.
type Operator interface {
	Start(ctx context.Context) (<-chan error, error)
}

type operator struct {
	ctrl                       k8s.Controller
	eg                         errgroup.Group
	namespace                  string
	client                     client.Client
	rotatorName                string
	targetReadReplicaIDEnvName string
	readReplicaEnabled         bool
	readReplicaLabelKey        string
}

// New returns Indexer object if no error occurs.
func New(namespace, agentName, rotatorName, targetReadReplicaIDEnvName string, opts ...Option) (o Operator, err error) {
	operator := new(operator)
	operator.namespace = namespace
	operator.targetReadReplicaIDEnvName = targetReadReplicaIDEnvName
	operator.rotatorName = rotatorName
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

	isAgent := func(pod *corev1.Pod) bool {
		return pod.Labels["app"] == agentName
	}

	podController := pod.New(
		pod.WithControllerName("pod reconciler for index operator"),
		pod.WithOnErrorFunc(func(err error) {
			log.Error("failed to reconcile:", err)
		}),
		pod.WithNamespace(operator.namespace),
		pod.WithOnReconcileFunc(operator.podOnReconcile),
		pod.WithLabels(map[string]string{
			"app": agentName,
		}),
		// To only reconcile for agent pods
		pod.WithForOpts(
			builder.WithPredicates(predicate.Funcs{
				CreateFunc: func(e event.CreateEvent) bool {
					pod, ok := e.Object.(*corev1.Pod)
					if !ok {
						return false
					}
					return isAgent(pod)
				},
				DeleteFunc: func(e event.DeleteEvent) bool {
					pod, ok := e.Object.(*corev1.Pod)
					if !ok {
						return false
					}
					return isAgent(pod)
				},
				UpdateFunc: func(e event.UpdateEvent) bool {
					pod, ok := e.ObjectNew.(*corev1.Pod)
					if !ok {
						return false
					}
					return isAgent(pod)
				},
				GenericFunc: func(e event.GenericEvent) bool {
					pod, ok := e.Object.(*corev1.Pod)
					if !ok {
						return false
					}
					return isAgent(pod)
				},
			}),
		),
	)

	operator.ctrl, err = k8s.New(
		k8s.WithResourceController(podController),
		k8s.WithLeaderElection(true, "vald-index-operator", operator.namespace),
	)
	if err != nil {
		return nil, err
	}

	client, err := client.New()
	if err != nil {
		return nil, err
	}
	operator.client = client

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

func (o *operator) podOnReconcile(ctx context.Context, podList map[string][]pod.Pod) {
	for k, v := range podList {
		for _, pod := range v {
			log.Debug("key", k, "name:", pod.Name, "annotations:", pod.Annotations)

			// rotate read replica if needed
			if o.readReplicaEnabled {
				if err := o.rotateIfNeeded(ctx, pod); err != nil {
					log.Error(err)
				}
			}
		}
	}
}

// rotateIfNeeded starts rotation job when the condition meets.
// This function is work in progress.
func (o *operator) rotateIfNeeded(ctx context.Context, pod pod.Pod) error {
	t, ok := pod.Annotations[vald.LastTimeSaveIndexTimestampAnnotationsKey]
	if !ok {
		log.Info("the agent pod has not saved index yet. skipping...")
		return nil
	}
	lastSavedTime, err := time.Parse(vald.TimeFormat, t)
	if err != nil {
		return fmt.Errorf("parsing last time saved time: %w", err)
	}

	podIdx, ok := pod.Labels[client.PodIndexLabel]
	if !ok {
		log.Info("no index label found. the agent is not StatefulSet? skipping...")
		return nil
	}

	var depList client.DeploymentList
	selector, err := o.client.LabelSelector(o.readReplicaLabelKey, client.SelectionOpEquals, []string{podIdx})
	if err != nil {
		return fmt.Errorf("creating label selector: %w", err)
	}
	listOpts := client.ListOptions{
		Namespace:     o.namespace,
		LabelSelector: selector,
	}
	if err := o.client.List(ctx, &depList, &listOpts); err != nil {
		return err
	}
	if len(depList.Items) == 0 {
		return errors.New("no readreplica deployment found")
	}
	dep := depList.Items[0]

	annotations := dep.GetAnnotations()
	t, ok = annotations[vald.LastTimeSnapshotTimestampAnnotationsKey]
	if ok {
		lastSnapshotTime, err := time.Parse(vald.TimeFormat, t)
		if err != nil {
			return fmt.Errorf("parsing last snapshot time: %w", err)
		}

		if lastSnapshotTime.After(lastSavedTime) {
			log.Info("snapshot taken after the last save. skipping...")
			return nil
		}
	}

	log.Infof("rotation required for agent(id: %s)", podIdx)
	if err := o.createRotationJob(ctx, podIdx); err != nil {
		return fmt.Errorf("creating rotation job: %w", err)
	}
	return nil
}

func (o *operator) createRotationJob(ctx context.Context, podIdx string) error {
	var cronJob batchv1.CronJob
	if err := o.client.Get(ctx, o.rotatorName, o.namespace, &cronJob); err != nil {
		return err
	}

	// get all the rotation jobs and make sure the job is not running
	var jobList batchv1.JobList
	selector, err := o.client.LabelSelector("app", client.SelectionOpEquals, []string{o.rotatorName})
	if err != nil {
		return fmt.Errorf("creating label selector: %w", err)
	}
	if err := o.client.List(ctx, &jobList, &client.ListOptions{
		Namespace:     o.namespace,
		LabelSelector: selector,
	}); err != nil {
		return fmt.Errorf("listing jobs: %w", err)
	}
	for _, job := range jobList.Items {
		// no need to check finished jobs
		if job.Status.Active == 0 {
			continue
		}

		envs := job.Spec.Template.Spec.Containers[0].Env
		// since latest append wins, checking backbards
		for i := len(envs) - 1; i >= 0; i-- {
			env := envs[i]
			if env.Name == o.targetReadReplicaIDEnvName {
				if env.Value == podIdx {
					log.Infof("rotation job for the agent(id: %s) is already running. skipping...", podIdx)
					return nil
				} else {
					break
				}
			}
		}
	}

	// now we actually needs to create the rotator job
	log.Info("no job is running to rotate the agent(id:%s). creating a new job...", podIdx)
	spec := *cronJob.Spec.JobTemplate.Spec.DeepCopy()
	spec.Template.Spec.Containers[0].Env = append(spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
		Name:  o.targetReadReplicaIDEnvName,
		Value: podIdx,
	})

	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cronJob.Name + "-",
			Namespace:    o.namespace,
		},
		Spec: spec,
	}

	if err := o.client.Create(ctx, &job); err != nil {
		return err
	}

	return nil
}
