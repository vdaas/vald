//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package service manages the main logic of benchmark job.
package service

import (
	"context"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/job"
	benchjob "github.com/vdaas/vald/internal/k8s/vald/benchmark/job"
	benchscenario "github.com/vdaas/vald/internal/k8s/vald/benchmark/scenario"
	"github.com/vdaas/vald/internal/log"
)

type Scenario interface {
	PreStart(context.Context) error
	Start(context.Context) (<-chan error, error)
}

type scenario struct {
	jobs                    atomic.Value
	jobName                 string
	jobNamespace            string
	jobTemplate             string   // row manifest template data of rebalance job.
	jobObject               *job.Job // object generated from template.
	currentDeviationJobName atomic.Value

	rcd  time.Duration // reconcile check duration
	eg   errgroup.Group
	ctrl k8s.Controller
}

func New(opts ...Option) (Scenario, error) {
	sc := new(scenario)

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(sc); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	err := sc.initCtrl()
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// initCtrl creates the controller for reconcile k8s objects.
func (sc *scenario) initCtrl() (err error) {
	// watcher of vald benchmark scenario resource
	bs, err := benchscenario.New()
	if err != nil {
		return
	}
	// watcher of vald benchmark job resource
	bj, err := benchjob.New(
		benchjob.WithControllerName("benchmark job resource"),
		benchjob.WithNamespaces(sc.jobNamespace),
		benchjob.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		benchjob.WithOnReconcileFunc(sc.benchJobReconcile),
	)
	if err != nil {
		return
	}
	// watcher of configmap resource

	// watcher of job resource
	job, err := job.New(
		job.WithControllerName("benchmark job"),
		job.WithNamespaces(sc.jobNamespace),
		job.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		job.WithOnReconcileFunc(sc.jobReconcile),
	)
	if err != nil {
		return
	}

	sc.ctrl, err = k8s.New(
		k8s.WithControllerName("vald benchmark operator"),
		k8s.WithResourceController(bs),
		k8s.WithResourceController(bj),
		k8s.WithResourceController(job),
	)
	if err != nil {
		return
	}
	return
}

func (sc *scenario) jobReconcile(ctx context.Context, jobList map[string][]job.Job) {}

func (sc *scenario) benchJobReconcile(ctx context.Context, jobList map[string]benchjob.BenchmarkJobSpec) {
}

func (sc *scenario) PreStart(ctx context.Context) error {
	log.Infof("[benchmark scenario] start vald benchmark scenario")
	return nil
}

func (sc *scenario) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 2)
	sc.eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			}
		}
	})

	return ech, nil
}
