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
	"github.com/vdaas/vald/internal/k8s/job"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/k8s/vald"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

const (
	apiName = "vald/index/operator"
)

// Operator represents an interface for indexing.
type Operator interface {
	Start(ctx context.Context) (<-chan error, error)
}

type operator struct {
	ctrl                k8s.Controller
	eg                  errgroup.Group
	namespace           string
	client              client.Client
	readReplicaEnabled  bool
	readReplicaLabelKey string
}

// New returns Indexer object if no error occurs.
func New(agentName string, opts ...Option) (o Operator, err error) {
	operator := new(operator)
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
	)

	jobController, err := job.New(
		job.WithControllerName("job reconciler for index operator"),
		job.WithNamespaces(operator.namespace),
		job.WithOnErrorFunc(func(err error) {
			log.Errorf("failed to reconcile job resource:", err)
		}),
		job.WithOnReconcileFunc(operator.jobOnReconcile),
	)
	if err != nil {
		return nil, err
	}

	operator.ctrl, err = k8s.New(
		k8s.WithResourceController(podController),
		k8s.WithResourceController(jobController),
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

// TODO: implement agent pod reconcile logic to detect conditions to start indexing and saving.
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

// TODO: implement job reconcile logic to detect save job completion and to start rotation.
func (o *operator) jobOnReconcile(ctx context.Context, jobList map[string][]job.Job) {
	for k, v := range jobList {
		for _, job := range v {
			log.Debug("key", k, "name:", job.Name, "status:", job.Status)
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

	log.Infof("rotation required for agent id: %s. creating rotator job...", podIdx)
	// TODO: check if the rotator job already exists or queued
	//       then create rotation job
	return nil
}
