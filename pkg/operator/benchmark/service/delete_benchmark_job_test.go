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

// The Scenario label value is name+generation with no separator, so scenario
// "bench" at generation 12 and scenario "bench1" at generation 2 share the
// label "bench12". Deleting the old generation of one scenario must not delete
// the other scenario's jobs.
func Test_operator_deleteBenchmarkJob_onlyOwnedJobs(t *testing.T) {
	t.Parallel()
	const (
		insertJob = "bench-insert"
		searchJob = "bench-search"
	)
	job := func(name, owner string) v1.ValdBenchmarkJob {
		var bj v1.ValdBenchmarkJob
		bj.Name = name
		bj.SetOwnerReferences([]k8s.OwnerReference{{Kind: ScenarioKind, Name: owner}})
		return bj
	}
	tests := []struct {
		name      string
		deleteErr map[string]error
		want      []string
		wantErr   bool
	}{
		{
			name: "deletes only jobs owned by the scenario",
			want: []string{insertJob, searchJob},
		},
		{
			name:      "ignores NotFound and reports other delete errors",
			deleteErr: map[string]error{insertJob: k8smock.NewNotFoundError(insertJob), searchJob: errors.New("boom")},
			want:      []string{insertJob, searchJob},
			wantErr:   true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var deleted []string
			var gotLabel, gotNS string
			cli := &k8smock.ClientMock{
				ListFunc: func(_ context.Context, list k8s.ObjectList, opts ...k8s.ListOption) error {
					lo := new(k8s.ListOptions)
					for _, o := range opts {
						o.ApplyToList(lo)
					}
					gotNS = lo.Namespace
					if lo.LabelSelector != nil {
						gotLabel = lo.LabelSelector.String()
					}
					l, ok := list.(*v1.ValdBenchmarkJobList)
					if !ok {
						t.Fatalf("unexpected list type %T", list)
					}
					l.Items = []v1.ValdBenchmarkJob{
						job(insertJob, "bench"),
						job("bench1-insert", "bench1"),
						job(searchJob, "bench"),
					}
					return nil
				},
				DeleteFunc: func(_ context.Context, obj k8s.Object, _ ...k8s.DeleteOption) error {
					deleted = append(deleted, obj.GetName())
					return tc.deleteErr[obj.GetName()]
				},
			}
			o := &operator{
				jobNamespace: "ns",
				ctrl: &mockCtrl{GetManagerFunc: func() k8s.Manager {
					return &k8smock.ManagerMock{GetClientFunc: func() k8s.Client { return cli }}
				}},
			}
			err := o.deleteBenchmarkJob(t.Context(), "bench", "", 12)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if gotLabel != Scenario+"=bench12" || gotNS != "ns" {
				t.Errorf("list selector = %q ns = %q", gotLabel, gotNS)
			}
			if len(deleted) != len(tc.want) {
				t.Fatalf("deleted = %v, want %v", deleted, tc.want)
			}
			for i := range tc.want {
				if deleted[i] != tc.want[i] {
					t.Errorf("deleted = %v, want %v", deleted, tc.want)
				}
			}
		})
	}
}

// A scenario deleted and re-created with the same name shares the label with
// the old instance; deleting by the old instance's UID must keep the jobs of
// the new one.
func Test_operator_deleteBenchmarkJob_onlyOwnedUID(t *testing.T) {
	t.Parallel()
	job := func(name string, uid k8s.UID) v1.ValdBenchmarkJob {
		var bj v1.ValdBenchmarkJob
		bj.Name = name
		bj.SetOwnerReferences([]k8s.OwnerReference{{Kind: ScenarioKind, Name: "bench", UID: uid}})
		return bj
	}
	tests := []struct {
		name string
		uid  k8s.UID
		want []string
	}{
		{name: "old UID deletes only the old instance's jobs", uid: "uid-old", want: []string{"old-insert"}},
		{name: "new UID deletes only the new instance's jobs", uid: "uid-new", want: []string{"new-insert"}},
		{name: "unknown UID matches by name", uid: "", want: []string{"old-insert", "new-insert"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var deleted []string
			cli := &k8smock.ClientMock{
				ListFunc: func(_ context.Context, list k8s.ObjectList, _ ...k8s.ListOption) error {
					l, ok := list.(*v1.ValdBenchmarkJobList)
					if !ok {
						t.Fatalf("unexpected list type %T", list)
					}
					l.Items = []v1.ValdBenchmarkJob{job("old-insert", "uid-old"), job("new-insert", "uid-new")}
					return nil
				},
				DeleteFunc: func(_ context.Context, obj k8s.Object, _ ...k8s.DeleteOption) error {
					deleted = append(deleted, obj.GetName())
					return nil
				},
			}
			o := &operator{
				jobNamespace: "ns",
				ctrl: &mockCtrl{GetManagerFunc: func() k8s.Manager {
					return &k8smock.ManagerMock{GetClientFunc: func() k8s.Client { return cli }}
				}},
			}
			if err := o.deleteBenchmarkJob(t.Context(), "bench", tc.uid, 1); err != nil {
				t.Fatalf("err = %v", err)
			}
			if len(deleted) != len(tc.want) {
				t.Fatalf("deleted = %v, want %v", deleted, tc.want)
			}
			for i := range tc.want {
				if deleted[i] != tc.want[i] {
					t.Errorf("deleted = %v, want %v", deleted, tc.want)
				}
			}
		})
	}
}
