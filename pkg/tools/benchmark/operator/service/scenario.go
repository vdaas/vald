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
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/job"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	benchjob "github.com/vdaas/vald/internal/k8s/vald/benchmark/job"
	benchscenario "github.com/vdaas/vald/internal/k8s/vald/benchmark/scenario"
	"github.com/vdaas/vald/internal/log"
	corev1 "k8s.io/api/core/v1"
	k8smeta "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	scenarios atomic.Value
	benchjobs atomic.Value

	rcd  time.Duration // reconcile check duration
	eg   errgroup.Group
	ctrl k8s.Controller
}

// New creates the new scenario struct to handle vald benchmark job scenario.
// When the input options are invalid, the error will be returned.
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
	bs, err := benchscenario.New(
		benchscenario.WithControllerName("benchmark scenario resource"),
		benchscenario.WithNamespaces(sc.jobNamespace),
		benchscenario.WithOnErrorFunc(func(err error) {
			log.Errorf("failed to reconcile:", err)
		}),
		benchscenario.WithOnReconcileFunc(sc.benchScenarioReconcile),
	)
	if err != nil {
		return
	}

	// watcher of vald benchmark job resource
	bj, err := benchjob.New(
		benchjob.WithControllerName("benchmark job resource"),
		benchjob.WithOnErrorFunc(func(err error) {
			log.Errorf("failed to reconcile:", err)
		}),
		benchjob.WithNamespaces(sc.jobNamespace),
		benchjob.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		benchjob.WithOnReconcileFunc(sc.benchJobReconcile),
	)
	if err != nil {
		return
	}

	// watcher of job resource
	job, err := job.New(
		job.WithControllerName("benchmark job"),
		job.WithNamespaces(sc.jobNamespace),
		job.WithOnErrorFunc(func(err error) {
			log.Errorf("failed to reconcile:", err)
		}),
		job.WithOnReconcileFunc(sc.jobReconcile),
	)
	if err != nil {
		return
	}

	// create reconcile controller which watches valdbenchmarkscenario resource, valdbenchmarkjob resource, and job resource.
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

// jobReconcile gets k8s job list and watches theirs STATUS.
// Then, it processes according STATUS.
func (sc *scenario) jobReconcile(ctx context.Context, jobList map[string][]job.Job) {
	log.Warn(jobList)
	// TODO: impl logic
	return
}

// benchmarkJobReconcile gets the vald benchmark job resource list and create Job for running benchmark job.
func (sc *scenario) benchJobReconcile(ctx context.Context, jobList map[string]v1.ValdBenchmarkJob) {
	log.Debug("[reconcile benchmark job resource]: %v", jobList)
	if len(jobList) == 0 {
		if ok := sc.benchjobs.Load(); ok == nil {
			sc.benchjobs.Store(make([]*v1.ValdBenchmarkJob, 0))
		} else {
			sc.benchjobs.Swap(make([]*v1.ValdBenchmarkJob, 0))
		}
		log.Infof("[reconcile benchmark job resource] job resource not found")
		return
	}
	var cbjl []*v1.ValdBenchmarkJob
	if ok := sc.benchjobs.Load(); ok == nil {
		cbjl = make([]*v1.ValdBenchmarkJob, 0)
	} else {
		cbjl = ok.([]*v1.ValdBenchmarkJob)
	}
	for _, job := range jobList {
		err := sc.createJob(ctx, job)
		if err != nil {
			log.Errorf("[reconcile benchmark job] failed to create job: %s", err.Error())
		}
		cbjl = append(cbjl, &job)
	}
	sc.benchjobs.Swap(cbjl)
}

// benchScenarioReconcile gets the vald benchmark scenario list and create vald benchmark job resource according to it.
func (sc *scenario) benchScenarioReconcile(ctx context.Context, scenarioList map[string]v1.ValdBenchmarkScenario) {
	log.Debug("[reconcile scenario]: %#v", scenarioList)
	if len(scenarioList) == 0 {
		if ok := sc.scenarios.Load(); ok == nil {
			sc.scenarios.Store(make([]*v1.ValdBenchmarkScenario, 0))
		} else {
			sc.scenarios.Swap(make([]*v1.ValdBenchmarkScenario, 0))
		}
		log.Infof("[reconcile scenario] scenario not found")
		return
	}
	var cbsl []*v1.ValdBenchmarkScenario
	if ok := sc.scenarios.Load(); ok == nil {
		cbsl = make([]*v1.ValdBenchmarkScenario, len(scenarioList))
	} else {
		cbsl = ok.([]*v1.ValdBenchmarkScenario)
	}
	for _, scenario := range scenarioList {
		err := sc.createBenchmarkJob(ctx, scenario)
		if err != nil {
			log.Errorf("[reconcile scenario] failed to create job: %s", err.Error())
		}
		cbsl = append(cbsl, &scenario)
	}
	sc.scenarios.Swap(cbsl)
}

// createBenchmarkJob creates the ValdBenchmarkJob crd for running job.
func (sc *scenario) createBenchmarkJob(ctx context.Context, scenario v1.ValdBenchmarkScenario) error {
	ownerRef := []k8smeta.OwnerReference{
		{
			APIVersion: scenario.APIVersion,
			Kind:       scenario.Kind,
			Name:       scenario.Name,
			UID:        scenario.UID,
		},
	}
	for _, job := range scenario.Spec.Jobs {
		bj := new(v1.ValdBenchmarkJob)
		// set metadata.name, metadata.namespace
		bj.Name = scenario.GetName() + "-" + job.JobType + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
		bj.Namespace = scenario.GetNamespace()
		bj.SetOwnerReferences(ownerRef)

		// set specs
		bj.Spec = *job
		if bj.Spec.Target == nil {
			bj.Spec.Target = scenario.Spec.Target
		}
		if bj.Spec.Dataset == nil {
			bj.Spec.Dataset = scenario.Spec.Dataset
		}
		// create benchmark job resource
		c := sc.ctrl.GetManager().GetClient()
		if err := c.Create(ctx, bj); err != nil {
			// TODO: create new custom error
			return err
		}
	}
	return nil
}

// createJobTemplate creates the job template for crating k8s job resource.
// ns and name are required to set job environment value.
func createJobTemplate(ns, name string) job.Job {
	j := new(job.Job)
	backoffLimit := int32(0)
	j.Spec.BackoffLimit = &backoffLimit
	j.Spec.Template.Spec.ServiceAccountName = "vald-benchmark-operator"
	j.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyNever
	j.Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:            "vald-benchmark-job",
			Image:           "local-registry:5000/vdaas/vald-benchmark-job:latest",
			ImagePullPolicy: corev1.PullAlways,
			LivenessProbe: &corev1.Probe{
				InitialDelaySeconds: int32(60),
				PeriodSeconds:       int32(10),
				TimeoutSeconds:      int32(300),
				ProbeHandler: corev1.ProbeHandler{
					Exec: &corev1.ExecAction{
						Command: []string{
							"/go/bin/job",
							"-v",
						},
					},
				},
			},
			StartupProbe: &corev1.Probe{
				FailureThreshold: int32(30),
				PeriodSeconds:    int32(10),
				TimeoutSeconds:   int32(300),
				ProbeHandler: corev1.ProbeHandler{
					Exec: &corev1.ExecAction{
						Command: []string{
							"/go/bin/job",
							"-v",
						},
					},
				},
			},
			Ports: []corev1.ContainerPort{
				{
					Name:          "liveness",
					Protocol:      corev1.ProtocolTCP,
					ContainerPort: int32(3000),
				},
				{
					Name:          "readiness",
					Protocol:      corev1.ProtocolTCP,
					ContainerPort: int32(3001),
				},
			},
			Env: []corev1.EnvVar{
				{
					Name:  "POD_NAMESPACE",
					Value: ns,
				},
				{
					Name:  "POD_NAME",
					Value: name,
				},
			},
		},
	}
	return *j
}

// createJob creates benchmark job from benchmark job resource.
func (sc *scenario) createJob(ctx context.Context, bjr v1.ValdBenchmarkJob) (err error) {
	bj := createJobTemplate(bjr.Namespace, bjr.Name)
	bj.Name = bjr.Name
	bj.Namespace = bjr.Namespace
	bj.SetOwnerReferences(bjr.GetOwnerReferences())
	// create job
	c := sc.ctrl.GetManager().GetClient()
	if err = c.Create(ctx, &bj); err != nil {
		// TODO: create new custom error
		return err
	}
	if ok := sc.jobs.Load(); ok == nil {
		sc.jobs.Store([]string{bj.Name})
		return
	} else {
		jobs := sc.jobs.Load().([]string)
		jobs = append(jobs, bj.Name)
		sc.jobs.Swap(jobs)
	}
	return
}

func (sc *scenario) PreStart(ctx context.Context) error {
	log.Infof("[benchmark scenario] start vald benchmark scenario")
	return nil
}

func (sc *scenario) Start(ctx context.Context) (<-chan error, error) {
	scch, err := sc.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 2)
	sc.eg.Go(func() error {
		defer close(ech)
		dt := time.NewTicker(sc.rcd)
		defer dt.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-dt.C:
				// TODO: Get Resource
				_, ok := sc.scenarios.Load().([]*v1.ValdBenchmarkScenario)
				if !ok {
					log.Info("benchmark scenario resource is empty")
					continue
				}
			case err = <-scch:
				if err != nil {
					ech <- err
				}
			}
		}
	})

	return ech, nil
}
