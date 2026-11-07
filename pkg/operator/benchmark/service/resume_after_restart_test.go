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
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	k8smock "github.com/vdaas/vald/internal/test/mock/k8s"
)

// fakeBenchJobServer is a minimal API server for ValdBenchmarkJobs shared by
// the operator instances of a test, so a new instance (an operator restart)
// sees the jobs created by the previous one.
type fakeBenchJobServer struct {
	failCreate map[int]bool // 1-based Create call numbers that fail
	jobs       []v1.ValdBenchmarkJob
	statuses   []v1.ValdBenchmarkScenarioStatus
	creates    int
	failList   bool
}

func (s *fakeBenchJobServer) ctrl(t *testing.T) *mockCtrl {
	t.Helper()
	create := func(_ context.Context, obj k8s.Object) error {
		s.creates++
		if s.failCreate[s.creates] {
			return errors.New("create failed")
		}
		bj, ok := obj.(*v1.ValdBenchmarkJob)
		if !ok {
			return errors.Errorf("unexpected object %T", obj)
		}
		s.jobs = append(s.jobs, *bj.DeepCopy())
		return nil
	}
	statusUpdate := func(_ context.Context, obj k8s.Object) error {
		if sc, ok := obj.(*v1.ValdBenchmarkScenario); ok {
			s.statuses = append(s.statuses, sc.Status)
		}
		return nil
	}
	cli, ok := k8smock.NewHookManagerMock(create, nil, statusUpdate).GetClient().(*k8smock.ClientMock)
	if !ok {
		t.Fatal("NewHookManagerMock client is not a *k8smock.ClientMock")
	}
	cli.ListFunc = func(_ context.Context, list k8s.ObjectList, opts ...k8s.ListOption) error {
		if s.failList {
			return errors.New("list failed")
		}
		l, ok := list.(*v1.ValdBenchmarkJobList)
		if !ok {
			return errors.Errorf("unexpected list type %T", list)
		}
		var (
			ns  string
			sel map[string]string
		)
		for _, o := range opts {
			switch v := o.(type) {
			case k8s.InNamespace:
				ns = string(v)
			case k8s.MatchingLabels:
				sel = v
			default:
				return errors.Errorf("unexpected list option %T", o)
			}
		}
		l.Items = nil
	jobs:
		for _, bj := range s.jobs {
			if ns != "" && bj.GetNamespace() != ns {
				continue
			}
			for k, v := range sel {
				if bj.GetLabels()[k] != v {
					continue jobs
				}
			}
			l.Items = append(l.Items, *bj.DeepCopy())
		}
		return nil
	}
	return &mockCtrl{GetManagerFunc: func() k8s.Manager {
		return &k8smock.ManagerMock{GetClientFunc: func() k8s.Client { return cli }}
	}}
}

// TestOperator_benchScenarioReconcile_ResumesAfterRestart pins that a
// scenario whose benchmark job creation was cut short is resumed by a
// restarted operator (empty cache): only the spec jobs that do not exist yet
// are created, chained after the existing ones, and the scenario becomes
// Healthy. Before the fix a NotReady scenario was never resumed after a
// restart, and an Available one re-created the existing jobs.
func TestOperator_benchScenarioReconcile_ResumesAfterRestart(t *testing.T) {
	t.Parallel()

	const (
		name = "scenario"
		ns   = "default"
	)
	jobTypes := []string{jobTypeInsert, jobTypeSearch, "upsert"}
	newScenario := func(status v1.ValdBenchmarkScenarioStatus) v1.ValdBenchmarkScenario {
		jobs := make([]*v1.BenchmarkJobSpec, 0, len(jobTypes))
		for _, jt := range jobTypes {
			jobs = append(jobs, &v1.BenchmarkJobSpec{JobType: jt})
		}
		sc := v1.ValdBenchmarkScenario{
			Name: name, Namespace: ns, Generation: 1,
			Spec:   v1.ValdBenchmarkScenarioSpec{Jobs: jobs},
			Status: status,

			Kind: ScenarioKind,
		}
		return sc
	}
	tests := []struct {
		name string
		// restartStatus is the scenario status seen by the restarted operator.
		restartStatus v1.ValdBenchmarkScenarioStatus
		// failOn is the Create call of the first operator that fails (0: none).
		failOn int
		// orphanAfter adds, before the restart, a second job of jobTypes[i]
		// after the same predecessor as the first one (a failed Create that
		// still reached the API server); -1: none.
		orphanAfter int
	}{
		{
			name:          "NotReady scenario with the first job created resumes",
			failOn:        2,
			restartStatus: v1.BenchmarkScenarioNotReady,
			orphanAfter:   -1,
		},
		{
			name:          "NotReady scenario with two jobs created resumes",
			failOn:        3,
			restartStatus: v1.BenchmarkScenarioNotReady,
			orphanAfter:   -1,
		},
		{
			name:          "NotReady scenario with no job created creates all",
			failOn:        1,
			restartStatus: v1.BenchmarkScenarioNotReady,
			orphanAfter:   -1,
		},
		{
			name:          "Available scenario whose jobs were all created before its status was written creates none",
			restartStatus: v1.BenchmarkScenarioAvailable,
			orphanAfter:   -1,
		},
		{
			name:          "a duplicate job after the same predecessor does not break the chain",
			failOn:        3,
			restartStatus: v1.BenchmarkScenarioNotReady,
			orphanAfter:   1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := &fakeBenchJobServer{failCreate: map[int]bool{tc.failOn: true}}
			first := &operator{jobNamespace: ns, ctrl: srv.ctrl(t)}
			first.benchScenarioReconcile(t.Context(), map[string]v1.ValdBenchmarkScenario{
				name: newScenario(v1.BenchmarkScenarioAvailable),
			})
			before := len(srv.jobs)
			wantBefore := len(jobTypes)
			if tc.failOn > 0 {
				wantBefore = tc.failOn - 1
			}
			if before != wantBefore {
				t.Fatalf("first operator created %d jobs, want %d", before, wantBefore)
			}
			// the earliest job after a predecessor wins: the orphan sorts after
			// the original because its name has a larger suffix.
			if tc.orphanAfter >= 0 {
				orphan := *srv.jobs[tc.orphanAfter].DeepCopy()
				orphan.Name += "9"
				srv.jobs = append(srv.jobs, orphan)
			}

			// restart: a new operator with an empty cache, same API server.
			srv.statuses = nil
			restarted := &operator{jobNamespace: ns, ctrl: srv.ctrl(t)}
			restarted.benchScenarioReconcile(t.Context(), map[string]v1.ValdBenchmarkScenario{
				name: newScenario(tc.restartStatus),
			})

			var chain []v1.ValdBenchmarkJob
			for _, bj := range srv.jobs {
				if tc.orphanAfter < 0 || bj.Name != srv.jobs[tc.orphanAfter].Name+"9" {
					chain = append(chain, bj)
				}
			}
			if len(chain) != len(jobTypes) {
				names := make([]string, 0, len(chain))
				for _, bj := range chain {
					names = append(names, bj.Name)
				}
				t.Fatalf("jobs after the restart = %v, want one per spec job %v", names, jobTypes)
			}
			for i, bj := range chain {
				if bj.Spec.JobType != jobTypes[i] {
					t.Errorf("job %d has type %q, want %q", i, bj.Spec.JobType, jobTypes[i])
				}
				wantPrev := ""
				if i > 0 {
					wantPrev = chain[i-1].Name
				}
				if got := bj.GetAnnotations()[BeforeJobName]; got != wantPrev {
					t.Errorf("job %d waits for %q, want %q", i, got, wantPrev)
				}
			}
			if len(srv.statuses) == 0 || srv.statuses[len(srv.statuses)-1] != v1.BenchmarkScenarioHealthy {
				t.Errorf("scenario status writes after the restart = %v, want it to end with %q", srv.statuses, v1.BenchmarkScenarioHealthy)
			}
			sc := restarted.getAtomicScenario()[name]
			if sc == nil {
				t.Fatal("restarted operator does not cache the scenario")
			}
			if sc.jobsIncomplete {
				t.Error("scenario is still marked incomplete after the resumed creation")
			}
			for _, bj := range chain {
				if _, ok := sc.BenchJobStatus[bj.Name]; !ok {
					t.Errorf("scenario does not track job %s: %v", bj.Name, sc.BenchJobStatus)
				}
			}

			// a further reconcile creates nothing.
			n := len(srv.jobs)
			restarted.benchScenarioReconcile(t.Context(), map[string]v1.ValdBenchmarkScenario{
				name: newScenario(v1.BenchmarkScenarioHealthy),
			})
			if len(srv.jobs) != n {
				t.Errorf("reconcile after the resumed creation created %d more jobs", len(srv.jobs)-n)
			}
		})
	}
}

// TestOperator_benchScenarioReconcile_RestartListFailureCreatesNothing pins
// that when the existing jobs of an uncached scenario cannot be listed no job
// is created (it could duplicate one) and the scenario stays uncached, so the
// next reconcile retries and then resumes.
func TestOperator_benchScenarioReconcile_RestartListFailureCreatesNothing(t *testing.T) {
	t.Parallel()

	const (
		name = "scenario"
		ns   = "default"
	)
	sc := v1.ValdBenchmarkScenario{
		Name: name, Namespace: ns, Generation: 1,
		Spec: v1.ValdBenchmarkScenarioSpec{Jobs: []*v1.BenchmarkJobSpec{
			{JobType: jobTypeInsert}, {JobType: jobTypeSearch},
		}},
		Status: v1.BenchmarkScenarioNotReady,

		Kind: ScenarioKind,
	}
	list := map[string]v1.ValdBenchmarkScenario{name: sc}

	srv := &fakeBenchJobServer{failList: true}
	o := &operator{jobNamespace: ns, ctrl: srv.ctrl(t)}
	o.benchScenarioReconcile(t.Context(), list)
	if srv.creates != 0 {
		t.Fatalf("created %d jobs although the existing ones could not be listed", srv.creates)
	}
	if _, ok := o.getAtomicScenario()[name]; ok {
		t.Fatal("scenario is cached although its job state is unknown")
	}

	srv.failList = false
	o.benchScenarioReconcile(t.Context(), list)
	if len(srv.jobs) != len(sc.Spec.Jobs) {
		t.Fatalf("created %d jobs on the retry, want %d", len(srv.jobs), len(sc.Spec.Jobs))
	}
}

// TestOperator_refreshStatuses_RetriesFailedJobList pins that a scenario whose
// existing jobs could not be listed is retried by the periodic refresh, so the
// retry does not depend on another scenario event.
func TestOperator_refreshStatuses_RetriesFailedJobList(t *testing.T) {
	t.Parallel()

	const (
		name = "scenario"
		ns   = "default"
	)
	sc := v1.ValdBenchmarkScenario{
		Name: name, Namespace: ns, Generation: 1,
		Spec: v1.ValdBenchmarkScenarioSpec{Jobs: []*v1.BenchmarkJobSpec{
			{JobType: jobTypeInsert}, {JobType: jobTypeSearch},
		}},

		Kind: ScenarioKind,
	}

	srv := &fakeBenchJobServer{failList: true}
	o := &operator{jobNamespace: ns, ctrl: srv.ctrl(t)}
	o.benchScenarioReconcile(t.Context(), map[string]v1.ValdBenchmarkScenario{name: sc})
	if srv.creates != 0 {
		t.Fatalf("created %d jobs although the existing ones could not be listed", srv.creates)
	}

	// the list still fails: nothing is created and the retry is kept.
	o.refreshStatuses(t.Context())
	if srv.creates != 0 {
		t.Fatalf("refresh created %d jobs although the list failed", srv.creates)
	}

	srv.failList = false
	o.refreshStatuses(t.Context())
	if len(srv.jobs) != len(sc.Spec.Jobs) {
		t.Fatalf("refresh created %d jobs, want %d", len(srv.jobs), len(sc.Spec.Jobs))
	}
	got := o.getAtomicScenario()[name]
	if got == nil || got.Crd.Status != v1.BenchmarkScenarioHealthy {
		t.Fatalf("scenario after the retry = %+v, want it cached as %q", got, v1.BenchmarkScenarioHealthy)
	}
	if len(srv.statuses) == 0 || srv.statuses[len(srv.statuses)-1] != v1.BenchmarkScenarioHealthy {
		t.Errorf("scenario status writes = %v, want it to end with %q", srv.statuses, v1.BenchmarkScenarioHealthy)
	}

	// the retry is done once.
	o.refreshStatuses(t.Context())
	if len(srv.jobs) != len(sc.Spec.Jobs) {
		t.Errorf("a further refresh created %d more jobs", len(srv.jobs)-len(sc.Spec.Jobs))
	}
}

// TestOperator_benchScenarioReconcile_RecreatedScenarioIgnoresOldJobs pins
// that a scenario deleted and re-created with the same name does not adopt the
// jobs of the old instance, which share its Scenario label until the GC
// removes them: the new instance creates its own jobs.
func TestOperator_benchScenarioReconcile_RecreatedScenarioIgnoresOldJobs(t *testing.T) {
	t.Parallel()

	const (
		name = "scenario"
		ns   = "default"
	)
	newScenario := func(uid k8s.UID) v1.ValdBenchmarkScenario {
		sc := v1.ValdBenchmarkScenario{
			Name: name, Namespace: ns, Generation: 1, UID: uid,
			Spec: v1.ValdBenchmarkScenarioSpec{Jobs: []*v1.BenchmarkJobSpec{
				{JobType: jobTypeInsert}, {JobType: jobTypeSearch},
			}},

			Kind: ScenarioKind,
		}
		return sc
	}

	srv := &fakeBenchJobServer{}
	old := &operator{jobNamespace: ns, ctrl: srv.ctrl(t)}
	old.benchScenarioReconcile(t.Context(), map[string]v1.ValdBenchmarkScenario{name: newScenario("uid-old")})
	if len(srv.jobs) != 2 {
		t.Fatalf("old instance created %d jobs, want 2", len(srv.jobs))
	}
	for i := range srv.jobs {
		srv.jobs[i].Status = v1.BenchmarkJobCompleted
	}

	// the re-created scenario is uncached: the old one left the scenario list.
	o := &operator{jobNamespace: ns, ctrl: srv.ctrl(t)}
	o.benchScenarioReconcile(t.Context(), map[string]v1.ValdBenchmarkScenario{name: newScenario("uid-new")})
	if len(srv.jobs) != 4 {
		t.Fatalf("jobs after the re-creation = %d, want the 2 old ones plus 2 new ones", len(srv.jobs))
	}
	sc := o.getAtomicScenario()[name]
	if sc == nil {
		t.Fatal("re-created scenario is not cached")
	}
	for _, bj := range srv.jobs[2:] {
		owners := bj.GetOwnerReferences()
		if len(owners) != 1 || owners[0].UID != "uid-new" {
			t.Errorf("job %s owners = %v, want the new instance", bj.Name, owners)
		}
		if st, ok := sc.BenchJobStatus[bj.Name]; !ok || st != v1.BenchmarkJobNotReady {
			t.Errorf("job %s tracked as %q (%v), want %q", bj.Name, st, ok, v1.BenchmarkJobNotReady)
		}
	}
	for _, bj := range srv.jobs[:2] {
		if _, ok := sc.BenchJobStatus[bj.Name]; ok {
			t.Errorf("re-created scenario tracks job %s of the old instance", bj.Name)
		}
	}
}
