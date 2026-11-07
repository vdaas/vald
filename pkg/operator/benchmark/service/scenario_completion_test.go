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
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
)

const (
	completionScenario  = "scenario"
	completionNamespace = "default"
	completionUID       = k8s.UID("scenario-uid")

	// BenchmarkJobSpec.JobType values of insert and search jobs.
	jobTypeInsert = "insert"
	jobTypeSearch = "search"
)

func completionScenarioCRD(status v1.ValdBenchmarkScenarioStatus) v1.ValdBenchmarkScenario {
	return v1.ValdBenchmarkScenario{
		Name:       completionScenario,
		Namespace:  completionNamespace,
		UID:        completionUID,
		Generation: 1,
		Spec: v1.ValdBenchmarkScenarioSpec{
			Jobs: []*v1.BenchmarkJobSpec{{JobType: jobTypeInsert}, {JobType: jobTypeSearch}},
		},
		Status: status,
	}
}

func completedBenchJob(name string) v1.ValdBenchmarkJob {
	return v1.ValdBenchmarkJob{
		Name:      name,
		Namespace: completionNamespace,
		OwnerReferences: []k8s.OwnerReference{{
			Kind: ScenarioKind,
			Name: completionScenario,
			UID:  completionUID,
		}},
		Status: v1.BenchmarkJobCompleted,
	}
}

// TestOperator_refreshStatuses_CompletesScenarioOfDeletedJobs pins that a
// Healthy scenario whose completed benchmark jobs are deleted (and so pruned
// from its tracked statuses) is still completed by the status refresh.
func TestOperator_refreshStatuses_CompletesScenarioOfDeletedJobs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		listed map[string]v1.ValdBenchmarkJob
		name   string
	}{
		{
			name:   "one of two jobs deleted",
			listed: map[string]v1.ValdBenchmarkJob{"a": completedBenchJob("a")},
		},
		{
			name:   "every job deleted",
			listed: map[string]v1.ValdBenchmarkJob{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var statuses []v1.ValdBenchmarkScenarioStatus
			o := &operator{ctrl: newClientMockCtrl(nil, nil, func(_ context.Context, obj k8s.Object) error {
				if sc, ok := obj.(*v1.ValdBenchmarkScenario); ok {
					statuses = append(statuses, sc.Status)
				}
				return nil
			})}
			crd := completionScenarioCRD(v1.BenchmarkScenarioHealthy)
			scenarios := map[string]*scenario{completionScenario: {
				Crd: &crd,
				BenchJobStatus: map[string]v1.BenchmarkJobStatus{
					"a": v1.BenchmarkJobCompleted,
					"b": v1.BenchmarkJobCompleted,
				},
			}}
			o.scenarios.Store(&scenarios)

			o.benchJobReconcile(t.Context(), tc.listed)
			o.refreshStatuses(t.Context())

			want := []v1.ValdBenchmarkScenarioStatus{v1.BenchmarkScenarioCompleted}
			if !reflect.DeepEqual(statuses, want) {
				t.Errorf("scenario status writes = %v, want %v", statuses, want)
			}
		})
	}
}

// TestOperator_refreshStatuses_CompletesRestartedScenarioInAnyReconcileOrder
// pins that a Healthy scenario whose jobs have all completed is completed
// after an operator restart even when the benchmark job reconcile runs before
// the scenario is cached, and that it is not completed before its jobs are
// observed.
func TestOperator_refreshStatuses_CompletesRestartedScenarioInAnyReconcileOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		jobsFirst bool
	}{
		{name: "benchmark jobs reconciled first", jobsFirst: true},
		{name: "scenario reconciled first"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var statuses []v1.ValdBenchmarkScenarioStatus
			o := &operator{
				jobNamespace: completionNamespace,
				ctrl: newClientMockCtrl(
					func(_ context.Context, obj k8s.Object) error {
						return errors.Errorf("unexpected create of %s", obj.GetName())
					},
					nil,
					func(_ context.Context, obj k8s.Object) error {
						if sc, ok := obj.(*v1.ValdBenchmarkScenario); ok {
							statuses = append(statuses, sc.Status)
						}
						return nil
					},
				),
			}
			jobs := map[string]v1.ValdBenchmarkJob{
				"a": completedBenchJob("a"),
				"b": completedBenchJob("b"),
			}
			scenarios := map[string]v1.ValdBenchmarkScenario{
				completionScenario: completionScenarioCRD(v1.BenchmarkScenarioHealthy),
			}
			if tc.jobsFirst {
				o.benchJobReconcile(t.Context(), jobs)
				o.benchScenarioReconcile(t.Context(), scenarios)
			} else {
				o.benchScenarioReconcile(t.Context(), scenarios)
				// the jobs are not observed yet: not completed.
				o.refreshStatuses(t.Context())
				if len(statuses) != 0 {
					t.Fatalf("scenario status writes before its jobs are observed = %v, want none", statuses)
				}
				o.benchJobReconcile(t.Context(), jobs)
			}
			if err := o.checkAtomics(); err != nil {
				t.Errorf("checkAtomics() = %v, want nil", err)
			}
			o.refreshStatuses(t.Context())

			want := []v1.ValdBenchmarkScenarioStatus{v1.BenchmarkScenarioCompleted}
			if !reflect.DeepEqual(statuses, want) {
				t.Errorf("scenario status writes = %v, want %v", statuses, want)
			}
		})
	}
}

// TestOperator_refreshStatuses_RetriesIncompleteScenarioJobs pins that the
// status refresh keeps retrying the job creation of a scenario that is
// already NotReady (so a failed retry writes no status and fires no scenario
// event) until every job is created, and then marks it Healthy.
func TestOperator_refreshStatuses_RetriesIncompleteScenarioJobs(t *testing.T) {
	t.Parallel()

	var (
		creates  int
		created  []string
		statuses []v1.ValdBenchmarkScenarioStatus
	)
	o := &operator{
		jobNamespace: completionNamespace,
		ctrl: newClientMockCtrl(
			func(_ context.Context, obj k8s.Object) error {
				creates++
				// the first refresh fails; the second one succeeds.
				if creates == 1 {
					return errors.New("create failed")
				}
				created = append(created, obj.GetName())
				return nil
			},
			nil,
			func(_ context.Context, obj k8s.Object) error {
				if sc, ok := obj.(*v1.ValdBenchmarkScenario); ok {
					statuses = append(statuses, sc.Status)
				}
				return nil
			},
		),
	}
	crd := completionScenarioCRD(v1.BenchmarkScenarioNotReady)
	scenarios := map[string]*scenario{completionScenario: {
		Crd:            &crd,
		BenchJobStatus: map[string]v1.BenchmarkJobStatus{},
		jobsIncomplete: true,
	}}
	o.scenarios.Store(&scenarios)

	o.refreshStatuses(t.Context())
	if creates != 1 {
		t.Fatalf("first refresh made %d create calls, want 1", creates)
	}
	if len(statuses) != 0 {
		t.Fatalf("scenario status writes after a failed retry = %v, want none", statuses)
	}
	if sc := o.getAtomicScenario()[completionScenario]; !sc.jobsIncomplete {
		t.Fatalf("scenario not marked incomplete after a failed retry")
	}

	o.refreshStatuses(t.Context())
	if len(created) != len(crd.Spec.Jobs) {
		t.Fatalf("created jobs = %v, want %d", created, len(crd.Spec.Jobs))
	}
	want := []v1.ValdBenchmarkScenarioStatus{v1.BenchmarkScenarioHealthy}
	if !reflect.DeepEqual(statuses, want) {
		t.Errorf("scenario status writes = %v, want %v", statuses, want)
	}
	sc := o.getAtomicScenario()[completionScenario]
	if sc.jobsIncomplete || len(sc.BenchJobStatus) != len(crd.Spec.Jobs) {
		t.Errorf("scenario incomplete=%v tracks %v, want complete with %d jobs", sc.jobsIncomplete, sc.BenchJobStatus, len(crd.Spec.Jobs))
	}

	// a further refresh creates nothing.
	o.refreshStatuses(t.Context())
	if len(created) != len(crd.Spec.Jobs) {
		t.Errorf("refresh after a complete creation created %v", created[len(crd.Spec.Jobs):])
	}
}
