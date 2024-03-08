//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/job"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	benchjob "github.com/vdaas/vald/internal/k8s/vald/benchmark/job"
	benchscenario "github.com/vdaas/vald/internal/k8s/vald/benchmark/scenario"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Operator interface {
	PreStart(context.Context) error
	Start(context.Context) (<-chan error, error)
	LenBenchSC() map[v1.ValdBenchmarkScenarioStatus]int64
	LenBenchBJ() map[v1.BenchmarkJobStatus]int64
}

type scenario struct {
	Crd            *v1.ValdBenchmarkScenario
	BenchJobStatus map[string]v1.BenchmarkJobStatus
}

const (
	Scenario           = "scenario"
	ScenarioKind       = "ValdBenchmarkScenario"
	BenchmarkName      = "benchmark-name"
	BeforeJobName      = "before-job-name"
	BeforeJobNamespace = "before-job-namespace"
)

type operator struct {
	jobNamespace       string
	jobImage           string
	jobImagePullPolicy string
	scenarios          *atomic.Pointer[map[string]*scenario]
	benchjobs          *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
	jobs               *atomic.Pointer[map[string]string]
	rcd                time.Duration // reconcile check duration
	eg                 errgroup.Group
	ctrl               k8s.Controller
}

// New creates the new scenario struct to handle vald benchmark job scenario.
// When the input options are invalid, the error will be returned.
func New(opts ...Option) (Operator, error) {
	operator := new(operator)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(operator); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	err := operator.initCtrl()
	if err != nil {
		return nil, err
	}
	return operator, nil
}

// initCtrl creates the controller for reconcile k8s objects.
func (o *operator) initCtrl() error {
	// watcher of vald benchmark scenario resource
	benchScenario, err := benchscenario.New(
		benchscenario.WithControllerName("benchmark scenario resource"),
		benchscenario.WithNamespaces(o.jobNamespace),
		benchscenario.WithOnErrorFunc(func(err error) {
			log.Errorf("failed to reconcile benchmark scenario resource:", err)
		}),
		benchscenario.WithOnReconcileFunc(o.benchScenarioReconcile),
	)
	if err != nil {
		return err
	}

	// watcher of vald benchmark job resource
	benchJob, err := benchjob.New(
		benchjob.WithControllerName("benchmark job resource"),
		benchjob.WithOnErrorFunc(func(err error) {
			log.Errorf("failed to reconcile benchmark job resource:", err)
		}),
		benchjob.WithNamespaces(o.jobNamespace),
		benchjob.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		benchjob.WithOnReconcileFunc(o.benchJobReconcile),
	)
	if err != nil {
		return err
	}

	// watcher of job resource
	job, err := job.New(
		job.WithControllerName("benchmark job"),
		job.WithNamespaces(o.jobNamespace),
		job.WithOnErrorFunc(func(err error) {
			log.Errorf("failed to reconcile job resource:", err)
		}),
		job.WithOnReconcileFunc(o.jobReconcile),
	)
	if err != nil {
		return err
	}

	// create reconcile controller which watches valdbenchmarkscenario resource, valdbenchmarkjob resource, and job resource.
	o.ctrl, err = k8s.New(
		k8s.WithControllerName("vald benchmark scenario operator"),
		k8s.WithResourceController(benchScenario),
		k8s.WithResourceController(benchJob),
		k8s.WithResourceController(job),
	)
	return err
}

func (o *operator) getAtomicScenario() map[string]*scenario {
	if o.scenarios == nil {
		o.scenarios = &atomic.Pointer[map[string]*scenario]{}
		return nil
	}
	if v := o.scenarios.Load(); v != nil {
		return *(v)
	}
	return nil
}

func (o *operator) getAtomicBenchJob() map[string]*v1.ValdBenchmarkJob {
	if o.benchjobs == nil {
		o.benchjobs = &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
		return nil
	}
	if v := o.benchjobs.Load(); v != nil {
		return *(v)
	}
	return nil
}

func (o *operator) getAtomicJob() map[string]string {
	if o.jobs == nil {
		o.jobs = &atomic.Pointer[map[string]string]{}
		return nil
	}
	if v := o.jobs.Load(); v != nil {
		return *(v)
	}
	return nil
}

// jobReconcile gets k8s job list and watches theirs STATUS.
// Then, it processes according STATUS.
// skipcq: GO-R1005
func (o *operator) jobReconcile(ctx context.Context, jobList map[string][]job.Job) {
	log.Debug("[reconcile job] start")
	cjobs := o.getAtomicJob()
	if cjobs == nil {
		cjobs = map[string]string{}
	}
	if len(jobList) == 0 {
		log.Info("[reconcile job] no job is founded")
		o.jobs.Store(&(map[string]string{}))
		log.Debug("[reconcile job] finish")
		return
	}
	// benchmarkJobStatus is used for update benchmark job resource status
	benchmarkJobStatus := make(map[string]v1.BenchmarkJobStatus)
	// jobNames is used for check whether cjobs has delted job.
	// If cjobs has the delted job, it will be remove the end of jobReconcile function.
	jobNames := map[string]struct{}{}
	for _, jobs := range jobList {
		cnt := len(jobs)
		var name string
		for idx := range jobs {
			job := jobs[idx]
			if job.GetNamespace() != o.jobNamespace {
				continue
			}
			jobNames[job.GetName()] = struct{}{}
			if _, ok := cjobs[job.Name]; !ok && job.Status.CompletionTime == nil {
				cjobs[job.GetName()] = job.Namespace
				benchmarkJobStatus[job.GetName()] = v1.BenchmarkJobAvailable
				continue
			}
			name = job.GetName()
			if job.Status.Active == 0 && job.Status.Succeeded != 0 {
				cnt--
			}
		}
		if cnt == 0 && name != "" {
			benchmarkJobStatus[name] = v1.BenchmarkJobCompleted
		}
	}
	if len(benchmarkJobStatus) != 0 {
		_, err := o.updateBenchmarkJobStatus(ctx, benchmarkJobStatus)
		if err != nil {
			log.Error(err.Error)
		}
	}
	// delete job which is not be in `jobList` from cj.
	for k := range cjobs {
		if _, ok := jobNames[k]; !ok {
			delete(cjobs, k)
		}
	}
	o.jobs.Store(&cjobs)
	log.Debug("[reconcile job] finish")
}

// benchJobReconcile gets the vald benchmark job resource list and create Job for running benchmark job.
// skipcq: GO-R1005
func (o *operator) benchJobReconcile(ctx context.Context, benchJobList map[string]v1.ValdBenchmarkJob) {
	log.Debugf("[reconcile benchmark job resource] job list: %#v", benchJobList)
	cbjl := o.getAtomicBenchJob()
	if cbjl == nil {
		cbjl = make(map[string]*v1.ValdBenchmarkJob, 0)
	}
	if len(benchJobList) == 0 {
		log.Info("[reconcile benchmark job resource] job resource not found")
		o.benchjobs.Store(&(map[string]*v1.ValdBenchmarkJob{}))
		log.Debug("[reconcile benchmark job resource] finish")
		return
	}
	// jobStatus is used for update benchmarkJob CR status if updating is needed.
	jobStatus := make(map[string]v1.BenchmarkJobStatus)
	for k := range benchJobList {
		// update scenario status
		job := benchJobList[k]
		hasOwner := false
		if len(job.GetOwnerReferences()) > 0 {
			hasOwner = true
		}
		if scenarios := o.getAtomicScenario(); scenarios != nil && hasOwner {
			on := job.GetOwnerReferences()[0].Name
			if _, ok := scenarios[on]; ok {
				if scenarios[on].BenchJobStatus == nil {
					scenarios[on].BenchJobStatus = map[string]v1.BenchmarkJobStatus{}
				}
				scenarios[on].BenchJobStatus[job.Name] = job.Status
			}
			o.scenarios.Store(&scenarios)
		}
		if oldJob := cbjl[k]; oldJob != nil {
			if oldJob.GetGeneration() != job.GetGeneration() {
				if job.Status != "" && oldJob.Status != v1.BenchmarkJobCompleted {
					// delete old version job
					err := o.deleteJob(ctx, oldJob.GetName())
					if err != nil {
						log.Warnf("[reconcile benchmark job resource] failed to delete old version job: job name=%s, version=%d\t%s", oldJob.GetName(), oldJob.GetGeneration(), err.Error())
					}
					// create new version job
					err = o.createJob(ctx, job)
					if err != nil {
						log.Errorf("[reconcile benchmark job resource] failed to create new version job: %s", err.Error())
					}
					cbjl[k] = &job
				}
			} else if oldJob.Status == "" {
				jobStatus[oldJob.GetName()] = v1.BenchmarkJobAvailable
			}
		} else if len(job.Status) == 0 || job.Status == v1.BenchmarkJobNotReady {
			log.Info("[reconcile benchmark job resource] create job: ", k)
			err := o.createJob(ctx, job)
			if err != nil {
				log.Errorf("[reconcile benchmark job resource] failed to create job: %s", err.Error())
			}
			jobStatus[job.Name] = v1.BenchmarkJobAvailable
			cbjl[k] = &job
		}
	}
	// delete benchmark job which is not be in `benchJobList` from cbjl.
	for k := range cbjl {
		if _, ok := benchJobList[k]; !ok {
			delete(cbjl, k)
		}
	}
	o.benchjobs.Store(&cbjl)
	if len(jobStatus) != 0 {
		_, err := o.updateBenchmarkJobStatus(ctx, jobStatus)
		if err != nil {
			log.Errorf("[reconcile benchmark job resource] failed update job status: %s", err)
		}
	}
	log.Debug("[reconcile benchmark job resource] finish")
}

// benchScenarioReconcile gets the vald benchmark scenario list and create vald benchmark job resource according to it.
func (o *operator) benchScenarioReconcile(ctx context.Context, scenarioList map[string]v1.ValdBenchmarkScenario) {
	log.Debugf("[reconcile benchmark scenario resource] scenario list: %#v", scenarioList)
	cbsl := o.getAtomicScenario()
	if cbsl == nil {
		cbsl = map[string]*scenario{}
	}
	if len(scenarioList) == 0 {
		log.Info("[reconcile benchmark scenario resource]: scenario not found")
		o.scenarios.Store(&(map[string]*scenario{}))
		log.Debug("[reconcile benchmark scenario resource] finish")
		return
	}
	scenarioStatus := make(map[string]v1.ValdBenchmarkScenarioStatus)
	for name := range scenarioList {
		sc := scenarioList[name]
		if oldScenario := cbsl[name]; oldScenario == nil {
			// apply new crd which is not set yet.
			jobNames, err := o.createBenchmarkJob(ctx, sc)
			if err != nil {
				log.Errorf("[reconcile benchmark scenario resource] failed to create benchmark job resource: %s", err.Error())
			}
			cbsl[name] = &scenario{
				Crd: &sc,
				BenchJobStatus: func() map[string]v1.BenchmarkJobStatus {
					s := map[string]v1.BenchmarkJobStatus{}
					for _, v := range jobNames {
						s[v] = v1.BenchmarkJobNotReady
					}
					return s
				}(),
			}
			scenarioStatus[sc.GetName()] = v1.BenchmarkScenarioHealthy
		} else {
			// apply updated crd which is already applied.
			if oldScenario.Crd.GetGeneration() < sc.GetGeneration() {
				// delete old job resource. If it is succeeded, job pod will be deleted automatically because of OwnerReference.
				err := o.deleteBenchmarkJob(ctx, oldScenario.Crd.GetName(), oldScenario.Crd.Generation)
				if err != nil {
					log.Warnf("[reconcile benchmark scenario resource] failed to delete old version benchmark jobs: scenario name=%s, version=%d\t%s", oldScenario.Crd.GetName(), oldScenario.Crd.Generation, err.Error())
				}
				// create new benchmark job resources of new version.
				jobNames, err := o.createBenchmarkJob(ctx, sc)
				if err != nil {
					log.Errorf("[reconcile benchmark scenario resource] failed to create new version benchmark job resource: %s", err.Error())
				}
				cbsl[name] = &scenario{
					Crd: &sc,
					BenchJobStatus: func() map[string]v1.BenchmarkJobStatus {
						s := map[string]v1.BenchmarkJobStatus{}
						for _, v := range jobNames {
							s[v] = v1.BenchmarkJobNotReady
						}
						return s
					}(),
				}
			} else if oldScenario.Crd.Status != sc.Status {
				// only update status
				cbsl[name].Crd.Status = sc.Status
			}
		}
	}
	// delete stored crd which is not be in `scenarioList` from cbsl.
	for k := range cbsl {
		if _, ok := scenarioList[k]; !ok {
			delete(cbsl, k)
		}
	}
	o.scenarios.Store(&cbsl)
	// Update scenario status
	_, err := o.updateBenchmarkScenarioStatus(ctx, scenarioStatus)
	if err != nil {
		log.Errorf("[reconcile benchmark scenario resource] failed to update benchmark scenario resource status: %s", err.Error())
	}
	log.Debug("[reconcile benchmark scenario resource] finish")
}

// deleteBenchmarkJob deletes benchmark job resource according to given scenario name and generation.
func (o *operator) deleteBenchmarkJob(ctx context.Context, name string, generation int64) error {
	opts := new(client.DeleteAllOfOptions)
	client.MatchingLabels(map[string]string{
		Scenario: name + strconv.Itoa(int(generation)),
	}).ApplyToDeleteAllOf(opts)
	client.InNamespace(o.jobNamespace).ApplyToDeleteAllOf(opts)
	return o.ctrl.GetManager().GetClient().DeleteAllOf(ctx, &v1.ValdBenchmarkJob{}, opts)
}

// deleteJob deletes job resource according to given benchmark job name and generation.
func (o *operator) deleteJob(ctx context.Context, name string) error {
	cj := new(job.Job)
	err := o.ctrl.GetManager().GetClient().Get(ctx, client.ObjectKey{
		Namespace: o.jobNamespace,
		Name:      name,
	}, cj)
	if err != nil {
		return err
	}
	opts := new(client.DeleteOptions)
	deleteProgation := client.DeletePropagationBackground
	opts.PropagationPolicy = &deleteProgation
	return o.ctrl.GetManager().GetClient().Delete(ctx, cj, opts)
}

// createBenchmarkJob creates the ValdBenchmarkJob crd for running job.
func (o *operator) createBenchmarkJob(ctx context.Context, scenario v1.ValdBenchmarkScenario) ([]string, error) {
	ownerRef := []k8s.OwnerReference{
		{
			APIVersion: scenario.APIVersion,
			Kind:       scenario.Kind,
			Name:       scenario.Name,
			UID:        scenario.UID,
		},
	}
	jobNames := make([]string, 0)
	var beforeJobName string
	for _, job := range scenario.Spec.Jobs {
		bj := new(v1.ValdBenchmarkJob)
		// set metadata.name, metadata.namespace, OwnerReference
		bj.Name = scenario.GetName() + "-" + job.JobType + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
		bj.Namespace = scenario.GetNamespace()
		bj.SetOwnerReferences(ownerRef)
		// set label
		labels := map[string]string{
			Scenario: scenario.GetName() + strconv.Itoa(int(scenario.Generation)),
		}
		bj.SetLabels(labels)
		// set annotations for wating before job
		annotations := map[string]string{
			BeforeJobName:      beforeJobName,
			BeforeJobNamespace: o.jobNamespace,
		}
		bj.SetAnnotations(annotations)
		// set specs
		bj.Spec = *job
		if bj.Spec.Target == nil {
			bj.Spec.Target = scenario.Spec.Target
		}
		if bj.Spec.Dataset == nil {
			bj.Spec.Dataset = scenario.Spec.Dataset
		}
		// set status
		bj.Status = v1.BenchmarkJobNotReady
		// create benchmark job resource
		c := o.ctrl.GetManager().GetClient()
		if err := c.Create(ctx, bj); err != nil {
			return nil, errors.ErrFailedToCreateBenchmarkJob(err, bj.GetName())
		}
		jobNames = append(jobNames, bj.Name)
		beforeJobName = bj.Name
	}
	return jobNames, nil
}

// createJob creates benchmark job from benchmark job resource.
func (o *operator) createJob(ctx context.Context, bjr v1.ValdBenchmarkJob) error {
	label := map[string]string{
		BenchmarkName: bjr.GetName() + strconv.Itoa(int(bjr.GetGeneration())),
	}
	job, err := benchjob.NewBenchmarkJob(
		benchjob.WithContainerName(bjr.GetName()),
		benchjob.WithContainerImage(o.jobImage),
		benchjob.WithImagePullPolicy(benchjob.ImagePullPolicy(o.jobImagePullPolicy)),
	)
	if err != nil {
		return err
	}
	tpl, err := job.CreateJobTpl(
		benchjob.WithName(bjr.GetName()),
		benchjob.WithNamespace(bjr.Namespace),
		benchjob.WithLabel(label),
		benchjob.WithCompletions(int32(bjr.Spec.Repetition)),
		benchjob.WithParallelism(int32(bjr.Spec.Replica)),
		benchjob.WithOwnerRef([]k8s.OwnerReference{
			{
				APIVersion: bjr.APIVersion,
				Kind:       bjr.Kind,
				Name:       bjr.Name,
				UID:        bjr.UID,
			},
		}),
		benchjob.WithTTLSecondsAfterFinished(int32(bjr.Spec.TTLSecondsAfterFinished)),
	)
	if err != nil {
		return err
	}
	// create job
	c := o.ctrl.GetManager().GetClient()
	if err = c.Create(ctx, &tpl); err != nil {
		return errors.ErrFailedToCreateJob(err, tpl.GetName())
	}
	return nil
}

// updateBenchmarkScenarioStatus updates status of ValdBenchmarkScenarioResource.
func (o *operator) updateBenchmarkScenarioStatus(ctx context.Context, ss map[string]v1.ValdBenchmarkScenarioStatus) ([]string, error) {
	var sns []string
	if cbsl := o.getAtomicScenario(); cbsl != nil {
		for name, status := range ss {
			if scenario, ok := cbsl[name]; ok {
				if scenario.Crd.Status == status {
					continue
				}
				scenario.Crd.Status = status
				cli := o.ctrl.GetManager().GetClient()
				err := cli.Status().Update(ctx, scenario.Crd)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				sns = append(sns, name)
			}
		}
	}
	return sns, nil
}

// updateBenchmarkJobStatus updates status of ValdBenchmarkJobResource.
func (o *operator) updateBenchmarkJobStatus(ctx context.Context, js map[string]v1.BenchmarkJobStatus) ([]string, error) {
	var jns []string
	if cbjl := o.getAtomicBenchJob(); cbjl != nil {
		for name, status := range js {
			if bjob, ok := cbjl[name]; ok {
				if bjob.Status == status {
					continue
				}
				bjob.Status = status
				cli := o.ctrl.GetManager().GetClient()
				err := cli.Status().Update(ctx, bjob)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				jns = append(jns, name)
			}
		}
	}
	return jns, nil
}

func (o *operator) checkJobsStatus(ctx context.Context, jobs map[string]string) error {
	cbjl := o.getAtomicBenchJob()
	if jobs == nil || cbjl == nil {
		log.Infof("[check job status] no job launched")
		return nil
	}
	job := new(job.Job)
	c := o.ctrl.GetManager().GetClient()
	jobStatus := map[string]v1.BenchmarkJobStatus{}
	for name, ns := range jobs {
		err := c.Get(ctx, client.ObjectKey{
			Namespace: ns,
			Name:      name,
		}, job)
		if err != nil {
			return err
		}
		if job.Status.Active != 0 || job.Status.Failed != 0 {
			continue
		}
		if job.Status.Succeeded != 0 {
			if job, ok := cbjl[name]; ok {
				if job.Status != v1.BenchmarkJobCompleted {
					jobStatus[name] = v1.BenchmarkJobCompleted
				}
			}
		}

	}
	_, err := o.updateBenchmarkJobStatus(ctx, jobStatus)
	return err
}

// checkAtomics checks each atomic keeps consistency.
// skipcq: GO-R1005
func (o *operator) checkAtomics() error {
	cjl := o.getAtomicJob()
	cbjl := o.getAtomicBenchJob()
	cbsl := o.getAtomicScenario()
	bjCompletedCnt := 0
	bjAvailableCnt := 0

	if len(cbjl) == 0 && len(cbsl) > 0 && len(cjl) > 0 {
		log.Errorf("mismatch atomics: job=%v, benchjob=%v, scenario=%v", cjl, cbjl, cbsl)
		return errors.ErrMismatchBenchmarkAtomics(cjl, cbjl, cbsl)
	}

	for _, bj := range cbjl {
		// check bench and job
		if bj.Status == v1.BenchmarkJobCompleted {
			bjCompletedCnt++
		} else {
			bjAvailableCnt++
			if ns, ok := cjl[bj.GetName()]; !ok || ns != bj.GetNamespace() {
				log.Errorf("mismatch atomics: job=%v, benchjob=%v, scenario=%v", cjl, cbjl, cbsl)
				return errors.ErrMismatchBenchmarkAtomics(cjl, cbjl, cbsl)
			}
		}
		// check scenario and bench
		if owners := bj.GetOwnerReferences(); len(owners) > 0 {
			var scenarioName string
			for _, o := range owners {
				if o.Kind == ScenarioKind {
					scenarioName = o.Name
				}
			}
			if sc := cbsl[scenarioName]; sc != nil {
				if sc.BenchJobStatus[bj.Name] != bj.Status {
					log.Errorf("mismatch atomics: job=%v, benchjob=%v, scenario=%v", cjl, cbjl, cbsl)
					return errors.ErrMismatchBenchmarkAtomics(cjl, cbjl, cbsl)
				}
			} else {
				log.Errorf("mismatch atomics: job=%v, benchjob=%v, scenario=%v", cjl, cbjl, cbsl)
				return errors.ErrMismatchBenchmarkAtomics(cjl, cbjl, cbsl)
			}
		}
	}
	// check benchmarkjob status list and scenario benchmark job status list
	if len(cbsl) > 0 {
		for _, sc := range cbsl {
			for _, status := range sc.BenchJobStatus {
				if status == v1.BenchmarkJobCompleted {
					bjCompletedCnt--
				} else {
					bjAvailableCnt--
				}
			}
		}
		if bjAvailableCnt != 0 || bjCompletedCnt != 0 {
			log.Errorf("mismatch atomics: job=%v, benchjob=%v, scenario=%v", cjl, cbjl, cbsl)
			return errors.ErrMismatchBenchmarkAtomics(cjl, cbjl, cbsl)
		}
	}
	return nil
}

func (o *operator) LenBenchSC() map[v1.ValdBenchmarkScenarioStatus]int64 {
	m := map[v1.ValdBenchmarkScenarioStatus]int64{
		v1.BenchmarkScenarioAvailable: 0,
		v1.BenchmarkScenarioHealthy:   0,
		v1.BenchmarkScenarioNotReady:  0,
		v1.BenchmarkScenarioCompleted: 0,
	}
	if sc := o.getAtomicScenario(); sc != nil {
		for _, s := range sc {
			if _, ok := m[s.Crd.Status]; ok {
				m[s.Crd.Status] += 1
			} else {
				m[s.Crd.Status] = 1
			}
		}
	}
	return m
}

func (o *operator) LenBenchBJ() map[v1.BenchmarkJobStatus]int64 {
	m := map[v1.BenchmarkJobStatus]int64{
		v1.BenchmarkJobAvailable: 0,
		v1.BenchmarkJobHealthy:   0,
		v1.BenchmarkJobNotReady:  0,
		v1.BenchmarkJobCompleted: 0,
	}
	if bjs := o.getAtomicBenchJob(); bjs != nil {
		for _, bj := range bjs {
			if _, ok := m[bj.Status]; ok {
				m[bj.Status] += 1
			} else {
				m[bj.Status] = 1
			}
		}
	}
	return m
}

func (*operator) PreStart(context.Context) error {
	log.Infof("[benchmark scenario operator] start vald benchmark scenario operator")
	return nil
}

// skipcq: GO-R1005
func (o *operator) Start(ctx context.Context) (<-chan error, error) {
	scch, err := o.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 2)
	o.eg.Go(func() error {
		defer close(ech)
		rcticker := time.NewTicker(o.rcd)
		defer rcticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-rcticker.C:
				// check mismatch atomic
				err = o.checkAtomics()
				if err != nil {
					ech <- err
				}
				// determine whether benchmark scenario status should be updated.
				if cbsl := o.getAtomicScenario(); cbsl != nil {
					scenarioStatus := make(map[string]v1.ValdBenchmarkScenarioStatus)
					for name, scenario := range cbsl {
						if scenario.Crd.Status != v1.BenchmarkScenarioCompleted {
							cnt := len(scenario.BenchJobStatus)
							for _, bjob := range scenario.BenchJobStatus {
								if bjob == v1.BenchmarkJobCompleted {
									cnt--
								}
							}
							if cnt == 0 {
								scenarioStatus[name] = v1.BenchmarkScenarioCompleted
							}
						}
					}
					if _, err := o.updateBenchmarkScenarioStatus(ctx, scenarioStatus); err != nil {
						log.Errorf("failed to update benchmark scenario to %s\terror: %s", v1.BenchmarkJobCompleted, err.Error())
					}

				}
				// get job and check status
				if jobs := o.getAtomicJob(); jobs != nil {
					err = o.checkJobsStatus(ctx, jobs)
					if err != nil {
						log.Error(err.Error())
					}
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
