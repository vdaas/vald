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
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/job"
	"github.com/vdaas/vald/internal/k8s/pod"
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
	ctrl      k8s.Controller
	eg        errgroup.Group
	namespace string
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

	var k8sOpts []k8s.Option
	podLabelSelector := map[string]string{
		"app": agentName,
	}
	podOpts := k8s.WithResourceController(pod.New(
		pod.WithControllerName("pod reconciler for index operator"),
		pod.WithOnErrorFunc(func(err error) {
			log.Error("failed to reconcile:", err)
		}),
		pod.WithNamespace(operator.namespace),
		// TODO:
		pod.WithOnReconcileFunc(operator.podOnReconcile),
		pod.WithLabels(podLabelSelector),
	))
	k8sOpts = append(k8sOpts, podOpts)

	jobController, err := job.New(
		job.WithControllerName("job reconciler for index operator"),
		job.WithNamespaces(operator.namespace),
		job.WithOnErrorFunc(func(err error) {
			log.Errorf("failed to reconcile job resource:", err)
		}),
		// TODO:
		job.WithOnReconcileFunc(func(ctx context.Context, jobList map[string][]job.Job) {
			log.Debugf("reconciled job list: %v", jobList)
		}),
	)
	if err != nil {
		return nil, err
	}
	jobOpts := k8s.WithResourceController(jobController)
	k8sOpts = append(k8sOpts, jobOpts)

	operator.ctrl, err = k8s.New(k8sOpts...)
	if err != nil {
		return nil, err
	}
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
		log.Debug("key", k)
		for _, pod := range v {
			log.Debug("name:", pod.Name, "annotations:", pod.Annotations)
		}
	}
}
