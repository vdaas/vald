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

// Package service manages the main logic of benchmark job.
package service

import (
	"context"
	"maps"
	"reflect"
	"strconv"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/reconciler"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	benchjob "github.com/vdaas/vald/internal/k8s/vald/benchmark/job"
	benchscenario "github.com/vdaas/vald/internal/k8s/vald/benchmark/scenario"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Operator interface {
	PreStart(context.Context) error
	Start(context.Context) (<-chan error, error)
	GetScenarioStatus() map[v1.ValdBenchmarkScenarioStatus]int64
	GetBenchmarkJobStatus() map[v1.BenchmarkJobStatus]int64
}

type scenario struct {
	Crd            *v1.ValdBenchmarkScenario
	BenchJobStatus map[string]v1.BenchmarkJobStatus
	// createdJobs lists, in spec order, the benchmark jobs already created for
	// Crd's generation while jobsIncomplete is set, so that the retry resumes
	// after them instead of creating them again.
	createdJobs    []string
	jobsIncomplete bool
	// jobsPruned records that a tracked job status was dropped because its
	// benchmark job was deleted, so an empty BenchJobStatus means "every job
	// is gone" rather than "no job observed yet" (e.g. after a restart).
	jobsPruned bool
}

const (
	Scenario           = "scenario"
	ScenarioKind       = "ValdBenchmarkScenario"
	Name               = "name"
	ServerName         = "vald-benchmark-job"
	BenchmarkName      = "benchmark-name"
	BeforeJobName      = "before-job-name"
	BeforeJobNamespace = "before-job-namespace"
)

type operator struct {
	eg                 errgroup.Group
	ctrl               k8s.Controller
	scenarios          atomic.Pointer[map[string]*scenario]
	benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
	jobs               atomic.Pointer[map[string]string]
	listRetry          map[string]v1.ValdBenchmarkScenario
	jobNamespace       string
	jobImageRepository string
	jobImageTag        string
	jobImagePullPolicy string
	configMapName      string
	rcd                time.Duration
	// mu serializes the writers of scenarios, benchjobs and jobs: each controller's
	// reconcile callback and the periodic status refresh in Start run on different
	// goroutines and update the maps by load-copy-store. Readers stay lock-free.
	// mu also guards listRetry: the scenarios of the last reconcile whose
	// existing jobs could not be listed; they stay uncached and refreshStatuses
	// retries them so that the retry does not depend on another scenario event.
	mu sync.Mutex
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
			log.Errorf("failed to reconcile benchmark scenario resource: %v", err)
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
			log.Errorf("failed to reconcile benchmark job resource: %v", err)
		}),
		benchjob.WithNamespaces(o.jobNamespace),
		benchjob.WithOnReconcileFunc(o.benchJobReconcile),
	)
	if err != nil {
		return err
	}

	// watcher of job resource
	job := reconciler.NewListReconciler(
		"benchmark job",
		new(k8s.Job),
		func(ctx context.Context, list *k8s.JobList) {
			jobs := jobsByAppName(list)
			o.jobReconcile(ctx, jobs)
			releaseJobsByAppName(jobs)
		},
		reconciler.WithOnError(func(err error) {
			log.Errorf("failed to reconcile job resource: %v", err)
		}),
		reconciler.WithNamespace(o.jobNamespace),
	)

	// create reconcile controller which watches valdbenchmarkscenario resource, valdbenchmarkjob resource, and job resource.
	o.ctrl, err = k8s.New(
		k8s.WithControllerName("vald benchmark scenario operator"),
		k8s.WithResourceController(benchScenario),
		k8s.WithResourceController(benchJob),
		k8s.WithResourceController(job),
	)
	return err
}

// loadAtomic returns the map stored in p, or M's zero value (nil for the map
// types used here) when nothing has been stored yet. The atomic fields are
// value-typed and usable from the operator's zero value, so the former
// lazy-init nil checks — a data race when reached from multiple goroutines —
// are unnecessary.
func loadAtomic[M any](p *atomic.Pointer[M]) M {
	if v := p.Load(); v != nil {
		return *v
	}
	var zero M
	return zero
}

func (o *operator) getAtomicScenario() map[string]*scenario {
	return loadAtomic(&o.scenarios)
}

func (o *operator) getAtomicBenchJob() map[string]*v1.ValdBenchmarkJob {
	return loadAtomic(&o.benchjobs)
}

func (o *operator) getAtomicJob() map[string]string {
	return loadAtomic(&o.jobs)
}

// k8sClient returns the controller-runtime client owned by the controller's
// manager. It is looked up on demand instead of cached in a field because
// o.ctrl is populated by initCtrl (or injected by tests) after construction.
func (o *operator) k8sClient() k8s.Client {
	return o.ctrl.GetManager().GetClient()
}

// singleOwnerRef returns the sole OwnerReference pointing at the resource
// identified by the given coordinates, shared by the benchmark-job and
// k8s-job creation paths that previously spelled the literal out inline.
func singleOwnerRef(apiVersion, kind, name string, uid k8s.UID) []k8s.OwnerReference {
	return []k8s.OwnerReference{{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
		UID:        uid,
	}}
}

// jobReconcile gets k8s job list and watches theirs STATUS.
// Then, it processes according STATUS.
// skipcq: GO-R1005
func (o *operator) jobReconcile(ctx context.Context, jobList map[string][]k8s.Job) {
	o.mu.Lock()
	defer o.mu.Unlock()
	log.Debug("[reconcile job] start")
	cjobs := o.getAtomicJob()
	if cjobs == nil {
		cjobs = map[string]string{}
	} else {
		// copy-on-write: mutate a clone so a concurrent reader holding the
		// previous atomic snapshot never observes an in-progress mutation.
		cjobs = maps.Clone(cjobs)
	}
	if len(jobList) == 0 {
		log.Info("[reconcile job] no job is founded")
		o.jobs.Store(&(map[string]string{}))
		log.Debug("[reconcile job] finish")
		return
	}
	// benchmarkJobStatus is used for update benchmark job resource status
	benchmarkJobStatus := make(map[string]v1.BenchmarkJobStatus)
	// jobNames is used for check whether cjobs has deleted job.
	// If cjobs has the deleted job, it will be remove the end of jobReconcile function.
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
				benchmarkJobStatus[job.GetName()] = v1.BenchmarkJobHealthy
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
			log.Error(err.Error())
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
func (o *operator) benchJobReconcile(
	ctx context.Context, benchJobList map[string]v1.ValdBenchmarkJob,
) {
	o.mu.Lock()
	defer o.mu.Unlock()
	log.Debugf("[reconcile benchmark job resource] job list: %#v", benchJobList)
	cbjl := o.getAtomicBenchJob()
	if cbjl == nil {
		cbjl = make(map[string]*v1.ValdBenchmarkJob, 0)
	} else {
		cbjl = maps.Clone(cbjl)
	}
	if len(benchJobList) == 0 {
		// no early return: the scenarios still have to drop the statuses of
		// the deleted jobs below.
		log.Info("[reconcile benchmark job resource] job resource not found")
	}
	// jobStatus is used for update benchmarkJob CR status if updating is needed.
	jobStatus := make(map[string]v1.BenchmarkJobStatus)
	// scenarios is the published scenario snapshot; scenariosCopy is a lazily
	// created copy-on-write clone of it. Recording a job's status into its
	// owning scenario must never mutate the live map — nor the nested
	// BenchJobStatus map — that Start's ticker and checkAtomics range over
	// concurrently, so we clone the outer map and each touched entry (with its
	// nested map) once, mutate the clones, and Store the result once below.
	scenarios := o.getAtomicScenario()
	var scenariosCopy map[string]*scenario
	// ownScenario returns the copy of the published scenario name owned by
	// this pass, creating it on first use.
	ownScenario := func(name string) *scenario {
		if scenariosCopy == nil {
			scenariosCopy = maps.Clone(scenarios)
		}
		sc := scenariosCopy[name]
		// Deep-copy this entry (and its nested map) the first time it is
		// touched this pass: while sc still aliases the published entry
		// we replace it with an isolated copy we own; later jobs sharing
		// the same owner then mutate that copy in place.
		if sc == scenarios[name] {
			nsc := *sc
			nsc.BenchJobStatus = maps.Clone(sc.BenchJobStatus)
			if nsc.BenchJobStatus == nil {
				nsc.BenchJobStatus = map[string]v1.BenchmarkJobStatus{}
			}
			sc = &nsc
			scenariosCopy[name] = sc
		}
		return sc
	}
	for k := range benchJobList {
		job := benchJobList[k]
		// update scenario status
		if refs := job.GetOwnerReferences(); len(refs) > 0 && scenarios != nil {
			ownerName := refs[0].Name
			if _, ok := scenarios[ownerName]; ok {
				ownScenario(ownerName).BenchJobStatus[job.Name] = job.Status
			}
		}
		// update benchmark job
		if oldJob := cbjl[k]; oldJob != nil {
			if oldJob.GetGeneration() != job.GetGeneration() {
				if job.Status != "" && oldJob.Status != v1.BenchmarkJobCompleted {
					// delete old version job
					err := o.deleteJob(ctx, oldJob.GetName())
					if err != nil {
						log.Warnf("[reconcile benchmark job resource] failed to delete old version job: job name=%s, version=%d\t%s", oldJob.GetName(), oldJob.GetGeneration(), err.Error())
					}
					// create new version job
					if err = o.createJob(ctx, job); err != nil {
						// keep the old version in cbjl so the next reconcile still sees a
						// generation mismatch and retries delete+create. AlreadyExists is a
						// failure here too: the job name does not change across generations,
						// so it means the old version job is not gone yet (failed or
						// in-flight deletion) and the new spec has not been applied.
						log.Errorf("[reconcile benchmark job resource] failed to create new version job: %s", err.Error())
						continue
					}
					cbjl[k] = &job
				}
			} else if oldJob.Status == "" {
				jobStatus[oldJob.GetName()] = v1.BenchmarkJobAvailable
			}
		} else {
			if job.Status == "" || job.Status == v1.BenchmarkJobAvailable {
				log.Info("[reconcile benchmark job resource] create job: ", k)
				if err := o.createJob(ctx, job); err != nil {
					if !k8s.IsAlreadyExists(err) {
						log.Errorf("[reconcile benchmark job resource] failed to create job: %s", err.Error())
						// leave the job out of cbjl so the next reconcile sees it as
						// not created yet and retries, instead of caching it as healthy.
						continue
					}
					// the job already exists (e.g. created before a restart or before its
					// status was persisted); treat it as created so it is tracked and never
					// stalls on AlreadyExists.
					log.Infof("[reconcile benchmark job resource] job already exists: %s", k)
				}
				jobStatus[job.Name] = v1.BenchmarkJobHealthy
			}
			cbjl[k] = &job
		}
	}
	// prune job statuses of the scenarios whose benchmark job is gone from
	// benchJobList, so a deleted job does not keep the scenario from completing
	// or the atomics from matching. NotReady entries are kept: they are jobs
	// createBenchmarkJob just created that may not be in the informer cache
	// yet. Entries added above are all listed, so ranging over the published
	// statuses covers every candidate.
	for name, sc := range scenarios {
		for jn, st := range sc.BenchJobStatus {
			if _, ok := benchJobList[jn]; !ok && st != v1.BenchmarkJobNotReady {
				own := ownScenario(name)
				delete(own.BenchJobStatus, jn)
				own.jobsPruned = true
			}
		}
	}
	// delete benchmark job which is not be in `benchJobList` from cbjl.
	for k := range cbjl {
		if _, ok := benchJobList[k]; !ok {
			delete(cbjl, k)
		}
	}
	o.benchjobs.Store(&cbjl)
	if scenariosCopy != nil {
		o.scenarios.Store(&scenariosCopy)
	}
	if len(jobStatus) != 0 {
		_, err := o.updateBenchmarkJobStatus(ctx, jobStatus)
		if err != nil {
			log.Errorf("[reconcile benchmark job resource] failed update job status: %s", err)
		}
	}
	log.Debug("[reconcile benchmark job resource] finish")
}

// benchScenarioReconcile gets the vald benchmark scenario list and create vald benchmark job resource according to it.
func (o *operator) benchScenarioReconcile(
	ctx context.Context, scenarioList map[string]v1.ValdBenchmarkScenario,
) {
	o.mu.Lock()
	defer o.mu.Unlock()
	log.Debugf("[reconcile benchmark scenario resource] scenario list: %#v", scenarioList)
	cbsl := o.getAtomicScenario()
	if cbsl == nil {
		cbsl = map[string]*scenario{}
	} else {
		cbsl = maps.Clone(cbsl)
	}
	if len(scenarioList) == 0 {
		log.Info("[reconcile benchmark scenario resource]: scenario not found")
		o.listRetry = nil
		o.scenarios.Store(&(map[string]*scenario{}))
		log.Debug("[reconcile benchmark scenario resource] finish")
		return
	}
	scenarioStatus := make(map[string]v1.ValdBenchmarkScenarioStatus)
	listRetry := make(map[string]v1.ValdBenchmarkScenario)
	for name := range scenarioList {
		sc := scenarioList[name]
		if oldScenario := cbsl[name]; oldScenario == nil {
			// init atomic values for current scenario; the benchmark job
			// reconcile may have run first (e.g. after a restart) and skipped
			// the job statuses of this then-uncached scenario, so take them
			// from the cached benchmark jobs.
			cbsl[name] = &scenario{
				Crd:            &sc,
				BenchJobStatus: o.ownedBenchJobStatus(sc),
			}
			scenarioStatus[sc.GetName()] = sc.Status
			// apply new crd which is not set yet, or resume a NotReady one
			// whose creation state was lost with an operator restart.
			if sc.Status == "" || sc.Status == v1.BenchmarkScenarioAvailable ||
				sc.Status == v1.BenchmarkScenarioNotReady {
				// jobs of this generation may already exist (created before
				// a restart): resume after them instead of duplicating them.
				nsc, st, err := o.resumeScenarioJobs(ctx, sc)
				if err != nil {
					log.Errorf("[reconcile benchmark scenario resource] failed to list existing benchmark jobs: scenario name=%s\t%s", name, err.Error())
					// leave the scenario uncached so the next reconcile (or
					// refreshStatuses) retries instead of creating possibly
					// existing jobs.
					delete(cbsl, name)
					listRetry[name] = sc
					continue
				}
				cbsl[name] = nsc
				// use for updating benchmarkScenario CR status
				scenarioStatus[sc.GetName()] = st
			}
		} else {
			// apply updated crd which is already applied.
			if oldScenario.Crd.GetGeneration() < sc.GetGeneration() {
				// delete old job resource. If it is succeeded, job pod will be deleted automatically because of OwnerReference.
				err := o.deleteBenchmarkJob(ctx, oldScenario.Crd.GetName(), oldScenario.Crd.GetUID(), oldScenario.Crd.Generation)
				if err != nil {
					log.Warnf(
						"[reconcile benchmark scenario resource] failed to delete old version benchmark jobs: scenario name=%s, version=%d\t%s",
						oldScenario.Crd.GetName(),
						oldScenario.Crd.Generation,
						err.Error(),
					)
				}
				// create new benchmark job resources of new version.
				var ok bool
				cbsl[name], ok = o.createScenarioJobs(ctx, sc, nil)
				if !ok {
					scenarioStatus[name] = v1.BenchmarkScenarioNotReady
				}
			} else if oldScenario.jobsIncomplete {
				// a previous reconcile failed to create all jobs of this
				// generation: resume after the ones already created.
				var ok bool
				cbsl[name], ok = o.createScenarioJobs(ctx, sc, oldScenario)
				scenarioStatus[name] = createdScenarioStatus(ok)
			} else if oldScenario.Crd.Status != sc.Status {
				// only update status; copy-on-write the scenario and its Crd so a
				// reader holding the previous atomic snapshot never observes the
				// in-place field mutation.
				newSc := *oldScenario
				newCrd := *oldScenario.Crd
				newCrd.Status = sc.Status
				newSc.Crd = &newCrd
				cbsl[name] = &newSc
			}
		}
	}
	// delete stored crd which is not be in `scenarioList` from cbsl.
	for k := range cbsl {
		if _, ok := scenarioList[k]; !ok {
			delete(cbsl, k)
		}
	}
	o.listRetry = listRetry
	o.scenarios.Store(&cbsl)
	// Update scenario status
	_, err := o.updateBenchmarkScenarioStatus(ctx, scenarioStatus)
	if err != nil {
		log.Errorf("[reconcile benchmark scenario resource] failed to update benchmark scenario resource status: %s", err.Error())
	}
	log.Debug("[reconcile benchmark scenario resource] finish")
}

// ownedBenchJobStatus returns the statuses of the cached benchmark jobs owned
// by sc, or nil when there is none.
func (o *operator) ownedBenchJobStatus(
	sc v1.ValdBenchmarkScenario,
) map[string]v1.BenchmarkJobStatus {
	var s map[string]v1.BenchmarkJobStatus
	for name, bj := range o.getAtomicBenchJob() {
		if !ownedByScenario(bj.GetOwnerReferences(), sc.GetName(), sc.GetUID()) {
			continue
		}
		if s == nil {
			s = make(map[string]v1.BenchmarkJobStatus)
		}
		s[name] = bj.Status
	}
	return s
}

// deleteBenchmarkJob deletes benchmark job resource according to given scenario name and generation.
// The Scenario label is name+generation without a separator, so different
// scenarios can share a label value ("bench"+"12" == "bench1"+"2"); only jobs
// owned by the named scenario (with uid, when it is known) are deleted.
func (o *operator) deleteBenchmarkJob(
	ctx context.Context, name string, uid k8s.UID, generation int64,
) error {
	c := o.k8sClient()
	list := new(v1.ValdBenchmarkJobList)
	if err := c.List(ctx, list,
		k8s.MatchingLabels{Scenario: name + strconv.Itoa(int(generation))},
		k8s.InNamespace(o.jobNamespace),
	); err != nil {
		return err
	}
	var errs error
	for i := range list.Items {
		bj := &list.Items[i]
		if !ownedByScenario(bj.GetOwnerReferences(), name, uid) {
			continue
		}
		if err := c.Delete(ctx, bj); err != nil && !k8s.IsNotFound(err) {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

// listCreatedScenarioJobs rebuilds the job creation state of sc's generation
// from the ValdBenchmarkJobs that already exist, for a scenario that is not
// cached (e.g. after an operator restart). createBenchmarkJob creates the spec
// jobs in order and points each one at the previous one through the
// BeforeJobName annotation, so following that chain from the first job (empty
// BeforeJobName) while the job types match the spec yields the leading spec
// jobs already created. The result is the prev argument of
// createScenarioJobs; with no existing job it creates every spec job.
func (o *operator) listCreatedScenarioJobs(
	ctx context.Context, sc v1.ValdBenchmarkScenario,
) (*scenario, error) {
	list := new(v1.ValdBenchmarkJobList)
	if err := o.k8sClient().List(ctx, list,
		k8s.MatchingLabels{Scenario: sc.GetName() + strconv.Itoa(int(sc.Generation))},
		k8s.InNamespace(o.jobNamespace),
	); err != nil {
		return nil, err
	}
	next := make(map[string][]*v1.ValdBenchmarkJob, len(list.Items))
	for i := range list.Items {
		bj := &list.Items[i]
		// a scenario deleted and re-created with the same name shares the
		// label with the old instance, whose jobs may still be listed until
		// the GC and the informer catch up: only adopt jobs of this instance.
		if !ownedByScenario(bj.GetOwnerReferences(), sc.GetName(), sc.GetUID()) {
			continue
		}
		before := bj.GetAnnotations()[BeforeJobName]
		next[before] = append(next[before], bj)
	}
	var created []string
	status := make(map[string]v1.BenchmarkJobStatus, len(sc.Spec.Jobs))
	var before string
	for _, spec := range sc.Spec.Jobs {
		// a failed Create that still reached the API server can leave more
		// than one job after the same predecessor; follow the earliest one
		// (names end with the creation UnixNano) to stay deterministic.
		var bj *v1.ValdBenchmarkJob
		for _, c := range next[before] {
			if c.Spec.JobType == spec.JobType && (bj == nil || c.GetName() < bj.GetName()) {
				bj = c
			}
		}
		if bj == nil {
			break
		}
		created = append(created, bj.GetName())
		st := bj.Status
		if st == "" {
			st = v1.BenchmarkJobNotReady
		}
		status[bj.GetName()] = st
		before = bj.GetName()
	}
	return &scenario{
		Crd:            &sc,
		BenchJobStatus: status,
		createdJobs:    created,
	}, nil
}

// ownedByScenario reports whether owners include the scenario name; a
// non-empty uid must match too (createBenchmarkJob writes the scenario UID).
func ownedByScenario(owners []k8s.OwnerReference, name string, uid k8s.UID) bool {
	for _, ow := range owners {
		if ow.Kind == ScenarioKind && ow.Name == name && (uid == "" || ow.UID == uid) {
			return true
		}
	}
	return false
}

// deleteJob deletes job resource according to given benchmark job name and generation.
func (o *operator) deleteJob(ctx context.Context, name string) error {
	cj := new(k8s.Job)
	err := o.k8sClient().Get(ctx, k8s.ObjectKey{
		Namespace: o.jobNamespace,
		Name:      name,
	}, cj)
	if err != nil {
		return err
	}
	opts := new(k8s.DeleteOptions)
	deleteProgation := k8s.DeletePropagationBackground
	opts.PropagationPolicy = &deleteProgation
	return o.k8sClient().Delete(ctx, cj, opts)
}

// resumeScenarioJobs creates the jobs of an uncached scenario after the ones
// of its generation that already exist, and returns the scenario to cache and
// its status to write. A list error is returned before any job is created.
func (o *operator) resumeScenarioJobs(
	ctx context.Context, sc v1.ValdBenchmarkScenario,
) (*scenario, v1.ValdBenchmarkScenarioStatus, error) {
	prev, err := o.listCreatedScenarioJobs(ctx, sc)
	if err != nil {
		return nil, "", err
	}
	nsc, ok := o.createScenarioJobs(ctx, sc, prev)
	return nsc, createdScenarioStatus(ok), nil
}

// retryScenarioJobLists retries the scenarios whose existing jobs could not be
// listed by the last reconcile. The caller holds mu.
func (o *operator) retryScenarioJobLists(ctx context.Context) {
	if len(o.listRetry) == 0 {
		return
	}
	var next map[string]*scenario
	status := make(map[string]v1.ValdBenchmarkScenarioStatus)
	cbsl := o.getAtomicScenario()
	for name, sc := range o.listRetry {
		if _, ok := cbsl[name]; ok {
			delete(o.listRetry, name)
			continue
		}
		nsc, st, err := o.resumeScenarioJobs(ctx, sc)
		if err != nil {
			log.Errorf("[reconcile benchmark scenario resource] failed to list existing benchmark jobs: scenario name=%s\t%s", name, err.Error())
			continue
		}
		delete(o.listRetry, name)
		if next == nil {
			next = maps.Clone(cbsl)
			if next == nil {
				next = make(map[string]*scenario)
			}
		}
		next[name] = nsc
		status[name] = st
	}
	if next == nil {
		return
	}
	o.scenarios.Store(&next)
	if _, err := o.updateBenchmarkScenarioStatus(ctx, status); err != nil {
		log.Errorf("[reconcile benchmark scenario resource] failed to update benchmark scenario resource status: %s", err.Error())
	}
}

// retryIncompleteScenarioJobs resumes the job creation of the cached
// scenarios whose previous attempt failed. Their status is already NotReady,
// so a failed retry writes nothing and no scenario event would trigger the
// next one. The caller holds mu.
func (o *operator) retryIncompleteScenarioJobs(ctx context.Context) {
	cbsl := o.getAtomicScenario()
	var next map[string]*scenario
	status := make(map[string]v1.ValdBenchmarkScenarioStatus)
	for name, sc := range cbsl {
		if !sc.jobsIncomplete {
			continue
		}
		nsc, ok := o.createScenarioJobs(ctx, *sc.Crd, sc)
		if next == nil {
			next = maps.Clone(cbsl)
		}
		next[name] = nsc
		status[name] = createdScenarioStatus(ok)
	}
	if next == nil {
		return
	}
	o.scenarios.Store(&next)
	if _, err := o.updateBenchmarkScenarioStatus(ctx, status); err != nil {
		log.Errorf("[reconcile benchmark scenario resource] failed to update benchmark scenario resource status: %s", err.Error())
	}
}

// createScenarioJobs creates the benchmark jobs of sc and returns the scenario
// to cache plus whether every job of its spec is created now. prev is the
// cached scenario of the same generation whose job creation is incomplete, or
// the state rebuilt by listCreatedScenarioJobs for an uncached scenario (nil
// for a new generation); its created jobs are not created
// again and keep their tracked status. On failure the returned scenario
// records the jobs created so far so that the next reconcile resumes.
func (o *operator) createScenarioJobs(
	ctx context.Context, sc v1.ValdBenchmarkScenario, prev *scenario,
) (*scenario, bool) {
	var (
		created []string
		s       map[string]v1.BenchmarkJobStatus
		pruned  bool
	)
	if prev != nil {
		created = prev.createdJobs
		pruned = prev.jobsPruned
		s = maps.Clone(prev.BenchJobStatus)
		// keep the status written after the failed attempt (the listed one
		// may predate it) so that the caller's status update is not skipped.
		sc.Status = prev.Crd.Status
	}
	if s == nil {
		s = map[string]v1.BenchmarkJobStatus{}
	}
	jobNames, err := o.createBenchmarkJob(ctx, sc, created)
	// benchmark job resource status to store benchmarkScenario Atomic
	for _, v := range jobNames[len(created):] {
		s[v] = v1.BenchmarkJobNotReady
	}
	nsc := &scenario{
		Crd:            &sc,
		BenchJobStatus: s,
		jobsPruned:     pruned,
	}
	if err != nil {
		log.Errorf("[reconcile benchmark scenario resource] failed to create benchmark job resource: scenario name=%s, created %d of %d jobs\t%s",
			sc.GetName(), len(jobNames), len(sc.Spec.Jobs), err.Error())
		nsc.createdJobs = jobNames
		nsc.jobsIncomplete = true
		return nsc, false
	}
	return nsc, true
}

// createdScenarioStatus is the scenario status written after a job creation
// attempt: Healthy once all its jobs are created, NotReady otherwise so that
// the scenario is neither reported as running nor completed.
func createdScenarioStatus(ok bool) v1.ValdBenchmarkScenarioStatus {
	if ok {
		return v1.BenchmarkScenarioHealthy
	}
	return v1.BenchmarkScenarioNotReady
}

// createBenchmarkJob creates the ValdBenchmarkJob crd for running job.
// created lists the names of the leading spec jobs already created (in spec
// order); creation resumes after them. The returned names always start with
// created and include every job created before an error.
func (o *operator) createBenchmarkJob(
	ctx context.Context, scenario v1.ValdBenchmarkScenario, created []string,
) ([]string, error) {
	ownerRef := singleOwnerRef(scenario.APIVersion, scenario.Kind, scenario.Name, scenario.UID)
	jobs := scenario.Spec.Jobs
	if len(created) >= len(jobs) {
		return created, nil
	}
	jobNames := make([]string, 0, len(jobs))
	jobNames = append(jobNames, created...)
	var beforeJobName string
	if len(created) > 0 {
		beforeJobName = created[len(created)-1]
	}
	for _, job := range jobs[len(created):] {
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
		if err := o.k8sClient().Create(ctx, bj); err != nil {
			return jobNames, errors.ErrFailedToCreateBenchmarkJob(err, bj.GetName())
		}
		jobNames = append(jobNames, bj.Name)
		beforeJobName = bj.Name
	}
	return jobNames, nil
}

// createJob creates benchmark job from benchmark job resource.
func (o *operator) createJob(ctx context.Context, bjr v1.ValdBenchmarkJob) error {
	label := map[string]string{
		Name:          ServerName,
		BenchmarkName: bjr.GetName() + strconv.Itoa(int(bjr.GetGeneration())),
	}
	job, err := benchjob.NewBenchmarkJob(
		benchjob.WithContainerName(bjr.GetName()),
		benchjob.WithContainerImage(strings.Join([]string{o.jobImageRepository, o.jobImageTag}, ":")),
		benchjob.WithImagePullPolicy(benchjob.ImagePullPolicy(o.jobImagePullPolicy)),
		benchjob.WithOperatorConfigMap(o.configMapName),
	)
	if err != nil {
		return err
	}
	tpl, err := job.CreateJobTpl(
		benchjob.WithName(bjr.GetName()),
		benchjob.WithNamespace(bjr.Namespace),
		benchjob.WithLabel(label),
		benchjob.WithCompletions(int32(bjr.Spec.Repetition)), //nolint:gosec // CRD-validated small config value, no overflow
		benchjob.WithParallelism(int32(bjr.Spec.Replica)),    //nolint:gosec // CRD-validated small config value, no overflow
		benchjob.WithOwnerRef(singleOwnerRef(bjr.APIVersion, bjr.Kind, bjr.Name, bjr.UID)),
		benchjob.WithTTLSecondsAfterFinished(int32(bjr.Spec.TTLSecondsAfterFinished)), //nolint:gosec // CRD-validated small config value, no overflow
	)
	if err != nil {
		return err
	}
	// create job
	if err = o.k8sClient().Create(ctx, &tpl); err != nil {
		return errors.ErrFailedToCreateJob(err, tpl.GetName())
	}
	return nil
}

// updateCachedStatuses walks desired and, for every name where resolve
// returns a changed object (resolve builds a copy carrying the new status
// and returns it to persist, plus a commit func that writes the copy into
// the cache, or ok=false to skip), issues a status update and, only when it
// succeeds, commits the copy and collects the name. Failed updates are
// logged and skipped without touching the cache, so the next refresh still
// sees the old status and retries. Resolution stays with the caller because
// the cached shapes differ (scenario wrapper exposing Crd vs raw job object).
func updateCachedStatuses[S comparable](
	ctx context.Context,
	c k8s.Client,
	desired map[string]S,
	resolve func(name string, status S) (obj k8s.Object, commit func(), ok bool),
) (updated []string) {
	for name, status := range desired {
		obj, commit, ok := resolve(name, status)
		if !ok {
			continue
		}
		if err := c.Status().Update(ctx, obj); err != nil {
			log.Error(err.Error())
			continue
		}
		commit()
		updated = append(updated, name)
	}
	return updated
}

// updateBenchmarkScenarioStatus updates status of ValdBenchmarkScenarioResource.
//
//nolint:unparam // the updated-name list mirrors updateCachedStatuses for symmetry with updateBenchmarkJobStatus
func (o *operator) updateBenchmarkScenarioStatus(
	ctx context.Context, ss map[string]v1.ValdBenchmarkScenarioStatus,
) ([]string, error) {
	cbsl := o.getAtomicScenario()
	if cbsl == nil {
		return nil, nil
	}
	// copy-on-write: mutating sc.Crd.Status in place would race a concurrent
	// reader (GetScenarioStatus / Start's ticker) holding the published
	// snapshot, so clone the outer map plus each mutated entry and its Crd,
	// then Store the clone once so the cache still reflects the new status.
	var cbslCopy map[string]*scenario
	updated := updateCachedStatuses(ctx, o.k8sClient(), ss,
		func(name string, status v1.ValdBenchmarkScenarioStatus) (k8s.Object, func(), bool) {
			sc, ok := cbsl[name]
			if !ok || sc.Crd.Status == status {
				return nil, nil, false
			}
			newSc := *sc
			newCrd := *sc.Crd
			newCrd.Status = status
			newSc.Crd = &newCrd
			return newSc.Crd, func() {
				if cbslCopy == nil {
					cbslCopy = maps.Clone(cbsl)
				}
				cbslCopy[name] = &newSc
			}, true
		})
	if cbslCopy != nil {
		o.scenarios.Store(&cbslCopy)
	}
	return updated, nil
}

// updateBenchmarkJobStatus updates status of ValdBenchmarkJobResource.
//
//nolint:unparam // the updated-name list mirrors updateCachedStatuses for symmetry with updateBenchmarkScenarioStatus
func (o *operator) updateBenchmarkJobStatus(
	ctx context.Context, js map[string]v1.BenchmarkJobStatus,
) ([]string, error) {
	cbjl := o.getAtomicBenchJob()
	if cbjl == nil {
		return nil, nil
	}
	// copy-on-write: mutating bjob.Status in place would race a concurrent
	// reader (GetBenchmarkJobStatus / checkAtomics) holding the published
	// snapshot, so clone the outer map plus each mutated job, then Store the
	// clone once so the cache still reflects the new status.
	var cbjlCopy map[string]*v1.ValdBenchmarkJob
	updated := updateCachedStatuses(ctx, o.k8sClient(), js,
		func(name string, status v1.BenchmarkJobStatus) (k8s.Object, func(), bool) {
			bjob, ok := cbjl[name]
			if !ok || bjob.Status == status {
				return nil, nil, false
			}
			njob := *bjob
			njob.Status = status
			return &njob, func() {
				if cbjlCopy == nil {
					cbjlCopy = maps.Clone(cbjl)
				}
				cbjlCopy[name] = &njob
			}, true
		})
	if cbjlCopy != nil {
		o.benchjobs.Store(&cbjlCopy)
	}
	return updated, nil
}

func (o *operator) checkJobsStatus(ctx context.Context, jobs map[string]string) error {
	cbjl := o.getAtomicBenchJob()
	if jobs == nil || cbjl == nil {
		log.Infof("[check job status] no job launched")
		return nil
	}
	job := new(k8s.Job)
	c := o.k8sClient()
	jobStatus := map[string]v1.BenchmarkJobStatus{}
	var errs []error
	for name, ns := range jobs {
		err := c.Get(ctx, k8s.ObjectKey{
			Namespace: ns,
			Name:      name,
		}, job)
		if err != nil {
			// keep checking the remaining jobs; a single failing Get must not
			// block their status updates.
			errs = append(errs, err)
			continue
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
	if err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// checkAtomics checks each atomic keeps consistency.
// skipcq: GO-R1005
func (o *operator) checkAtomics() error {
	cjl := o.getAtomicJob()
	cbjl := o.getAtomicBenchJob()
	cbsl := o.getAtomicScenario()
	bjCompletedCnt := 0
	bjAvailableCnt := 0
	mismatch := func() error {
		log.Errorf("mismatch atomics: job=%v, benchjob=%v, scenario=%v", cjl, cbjl, cbsl)
		return errors.ErrMismatchBenchmarkAtomics(cjl, cbjl, cbsl)
	}

	if len(cbjl) == 0 && len(cbsl) > 0 && len(cjl) > 0 {
		return mismatch()
	}

	for _, bj := range cbjl {
		// check bench and job
		if bj.Status == v1.BenchmarkJobCompleted {
			bjCompletedCnt++
		} else {
			bjAvailableCnt++
			if ns, ok := cjl[bj.GetName()]; !ok || ns != bj.GetNamespace() {
				return mismatch()
			}
		}
		// check scenario resource and bench resource
		if owners := bj.GetOwnerReferences(); len(owners) > 0 {
			var scenarioName string
			for _, o := range owners {
				if o.Kind == ScenarioKind {
					scenarioName = o.Name
				}
			}
			if sc := cbsl[scenarioName]; sc != nil {
				if sc.BenchJobStatus[bj.Name] != bj.Status {
					return mismatch()
				}
			} else {
				return mismatch()
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
			return mismatch()
		}
	}
	return nil
}

func (o *operator) GetScenarioStatus() map[v1.ValdBenchmarkScenarioStatus]int64 {
	m := map[v1.ValdBenchmarkScenarioStatus]int64{
		v1.BenchmarkScenarioAvailable: 0,
		v1.BenchmarkScenarioHealthy:   0,
		v1.BenchmarkScenarioNotReady:  0,
		v1.BenchmarkScenarioCompleted: 0,
	}
	if sc := o.getAtomicScenario(); sc != nil {
		for _, s := range sc {
			m[s.Crd.Status] += 1
		}
	}
	return m
}

func (o *operator) GetBenchmarkJobStatus() map[v1.BenchmarkJobStatus]int64 {
	m := map[v1.BenchmarkJobStatus]int64{
		v1.BenchmarkJobAvailable: 0,
		v1.BenchmarkJobHealthy:   0,
		v1.BenchmarkJobNotReady:  0,
		v1.BenchmarkJobCompleted: 0,
	}
	if bjs := o.getAtomicBenchJob(); bjs != nil {
		for _, bj := range bjs {
			m[bj.Status] += 1
		}
	}
	return m
}

func (*operator) PreStart(context.Context) error {
	log.Infof("[benchmark scenario operator] start vald benchmark scenario operator")
	return nil
}

// refreshStatuses completes scenarios whose jobs have all completed and
// syncs the job statuses from the cluster. It holds mu because it writes the
// same maps as the reconcile callbacks.
func (o *operator) refreshStatuses(ctx context.Context) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.retryScenarioJobLists(ctx)
	o.retryIncompleteScenarioJobs(ctx)
	// determine whether benchmark scenario status should be updated.
	if cbsl := o.getAtomicScenario(); cbsl != nil {
		scenarioStatus := make(map[string]v1.ValdBenchmarkScenarioStatus)
		for name, scenario := range cbsl {
			// a scenario is complete once all its spec jobs are created and
			// every tracked job is completed. Deleted jobs are pruned from
			// BenchJobStatus, so it may track fewer jobs than the spec; but
			// no tracked job at all, unless the jobs were pruned or the spec
			// has none, means they are not observed yet (e.g. after a
			// restart), not that they finished.
			if scenario.Crd.Status != v1.BenchmarkScenarioCompleted &&
				!scenario.jobsIncomplete &&
				(len(scenario.BenchJobStatus) > 0 || scenario.jobsPruned ||
					len(scenario.Crd.Spec.Jobs) == 0) {
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
		if _, uerr := o.updateBenchmarkScenarioStatus(ctx, scenarioStatus); uerr != nil {
			log.Errorf("failed to update benchmark scenario to %s\terror: %s", v1.BenchmarkJobCompleted, uerr.Error())
		}
	}
	// get job and check status
	if jobs := o.getAtomicJob(); jobs != nil {
		if err := o.checkJobsStatus(ctx, jobs); err != nil {
			log.Error(err.Error())
		}
	}
}

// startErrChanBufferSize buffers the controller error stream plus this
// goroutine's own terminal error so neither send can block shutdown.
const startErrChanBufferSize = 2

// skipcq: GO-R1005
func (o *operator) Start(ctx context.Context) (<-chan error, error) {
	scch, err := o.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, startErrChanBufferSize)
	o.eg.Go(func() error {
		defer close(ech)
		rcticker := time.NewTicker(o.rcd)
		defer rcticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-rcticker.C:
				// check mismatch atomic; the reconcilers publish the three
				// atomics one at a time under mu, so check them under mu too.
				o.mu.Lock()
				err = o.checkAtomics()
				o.mu.Unlock()
				if err != nil {
					select {
					case ech <- err:
					case <-ctx.Done():
						return nil
					}
				}
				o.refreshStatuses(ctx)
			case serr, ok := <-scch:
				if !ok {
					// the controller closed its error stream (its manager
					// returned); a closed channel is always ready, so stop
					// selecting it or this loop busy-spins until ctx ends.
					scch = nil
					continue
				}
				if serr != nil {
					select {
					case ech <- serr:
					case <-ctx.Done():
						return nil
					}
				}
			}
		}
	})
	return ech, nil
}
