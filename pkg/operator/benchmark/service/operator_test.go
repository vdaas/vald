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
	"time"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/internal/test/mock"
	k8smock "github.com/vdaas/vald/internal/test/mock/k8s"
)

// mockCtrl is used for mock the request to the Kubernetes API.
type mockCtrl struct {
	StartFunc      func(ctx context.Context) (<-chan error, error)
	GetManagerFunc func() k8s.Manager
}

func (m *mockCtrl) Start(ctx context.Context) (<-chan error, error) {
	return m.StartFunc(ctx)
}

func (m *mockCtrl) GetManager() k8s.Manager {
	return m.GetManagerFunc()
}

// storeInto copies the map carried by src (if any) into the operator's
// value-typed atomic field; test cases keep building *atomic.Pointer carriers
// while the operator itself owns value-typed fields.
func storeInto[M any](dst, src *atomic.Pointer[M]) {
	if src == nil {
		return
	}
	if v := src.Load(); v != nil {
		dst.Store(v)
	}
}

//nolint:goconst
func Test_operator_getAtomicScenario(t *testing.T) {
	t.Parallel()
	type fields struct {
		scenarios *atomic.Pointer[map[string]*scenario]
	}
	type want struct {
		want map[string]*scenario
	}
	type test struct {
		fields     fields
		want       want
		checkFunc  func(want, map[string]*scenario) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		name       string
	}
	defaultCheckFunc := func(w want, got map[string]*scenario) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name:   "get nil when atomic has no resource",
			fields: fields{},
			want: want{
				want: nil,
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(t *testing.T) {
				t.Helper()
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
			},
		},
		{
			name: "get scenarios when scenario list is stored",
			fields: fields{
				scenarios: func() *atomic.Pointer[map[string]*scenario] {
					ap := atomic.Pointer[map[string]*scenario]{}
					ap.Store(&map[string]*scenario{
						"scenario": {
							Crd: &v1.ValdBenchmarkScenario{
								Spec: v1.ValdBenchmarkScenarioSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									Jobs: []*v1.BenchmarkJobSpec{
										{
											JobType: "insert",
											InsertConfig: &config.InsertConfig{
												SkipStrictExistCheck: false,
												Timestamp:            "",
											},
										},
										{
											JobType: "search",
											SearchConfig: &config.SearchConfig{
												Epsilon:              0.1,
												Radius:               -1,
												Num:                  10,
												MinNum:               10,
												Timeout:              "10s",
												EnableLinearSearch:   false,
												AggregationAlgorithm: "",
											},
										},
									},
								},
								Status: v1.BenchmarkScenarioHealthy,
							},
							BenchJobStatus: map[string]v1.BenchmarkJobStatus{
								"scenario-insert": v1.BenchmarkJobAvailable,
								"scenario-search": v1.BenchmarkJobAvailable,
							},
						},
					})
					return &ap
				}(),
			},
			want: want{
				want: map[string]*scenario{
					"scenario": {
						Crd: &v1.ValdBenchmarkScenario{
							Spec: v1.ValdBenchmarkScenarioSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								Jobs: []*v1.BenchmarkJobSpec{
									{
										JobType: "insert",
										InsertConfig: &config.InsertConfig{
											SkipStrictExistCheck: false,
											Timestamp:            "",
										},
									},
									{
										JobType: "search",
										SearchConfig: &config.SearchConfig{
											Epsilon:              0.1,
											Radius:               -1,
											Num:                  10,
											MinNum:               10,
											Timeout:              "10s",
											EnableLinearSearch:   false,
											AggregationAlgorithm: "",
										},
									},
								},
							},
							Status: v1.BenchmarkScenarioHealthy,
						},
						BenchJobStatus: map[string]v1.BenchmarkJobStatus{
							"scenario-insert": v1.BenchmarkJobAvailable,
							"scenario-search": v1.BenchmarkJobAvailable,
						},
					},
				},
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(t *testing.T) {
				t.Helper()
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := new(operator)
			storeInto(&o.scenarios, test.fields.scenarios)

			got := o.getAtomicScenario()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

//nolint:goconst
func Test_operator_getAtomicBenchJob(t *testing.T) {
	t.Parallel()
	type fields struct {
		benchjobs *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
	}
	type want struct {
		want map[string]*v1.ValdBenchmarkJob
	}
	type test struct {
		fields     fields
		want       want
		checkFunc  func(want, map[string]*v1.ValdBenchmarkJob) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		name       string
	}
	defaultCheckFunc := func(w want, got map[string]*v1.ValdBenchmarkJob) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name:   "get nil when atomic has no resource",
			fields: fields{},
			want: want{
				want: nil,
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(t *testing.T) {
				t.Helper()
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
			},
		},
		{
			name: "get benchjobs when job list is stored",
			fields: fields{
				benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
					ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
					m := map[string]*v1.ValdBenchmarkJob{
						"scenario-insert": {
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								JobType: "insert",
								InsertConfig: &config.InsertConfig{
									SkipStrictExistCheck: false,
									Timestamp:            "",
								},
							},
							Status: v1.BenchmarkJobAvailable,
						},
						"scenario-search": {
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								JobType: "search",
								SearchConfig: &config.SearchConfig{
									Epsilon:              0.1,
									Radius:               -1,
									Num:                  10,
									MinNum:               10,
									Timeout:              "10s",
									EnableLinearSearch:   false,
									AggregationAlgorithm: "",
								},
							},
							Status: v1.BenchmarkJobAvailable,
						},
					}
					ap.Store(&m)
					return &ap
				}(),
			},
			want: want{
				want: map[string]*v1.ValdBenchmarkJob{
					"scenario-insert": {
						Spec: v1.BenchmarkJobSpec{
							Target: &v1.BenchmarkTarget{
								Host: "localhost",
								Port: 8080,
							},
							Dataset: &v1.BenchmarkDataset{
								Name:    "fashion-minsit",
								Group:   "train",
								Indexes: 10000,
								Range: &config.BenchmarkDatasetRange{
									Start: 0,
									End:   10000,
								},
								URL: "",
							},
							JobType: "insert",
							InsertConfig: &config.InsertConfig{
								SkipStrictExistCheck: false,
								Timestamp:            "",
							},
						},
						Status: v1.BenchmarkJobAvailable,
					},
					"scenario-search": {
						Spec: v1.BenchmarkJobSpec{
							Target: &v1.BenchmarkTarget{
								Host: "localhost",
								Port: 8080,
							},
							Dataset: &v1.BenchmarkDataset{
								Name:    "fashion-minsit",
								Group:   "train",
								Indexes: 10000,
								Range: &config.BenchmarkDatasetRange{
									Start: 0,
									End:   10000,
								},
								URL: "",
							},
							JobType: "search",
							SearchConfig: &config.SearchConfig{
								Epsilon:              0.1,
								Radius:               -1,
								Num:                  10,
								MinNum:               10,
								Timeout:              "10s",
								EnableLinearSearch:   false,
								AggregationAlgorithm: "",
							},
						},
						Status: v1.BenchmarkJobAvailable,
					},
				},
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(t *testing.T) {
				t.Helper()
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := new(operator)
			storeInto(&o.benchjobs, test.fields.benchjobs)

			got := o.getAtomicBenchJob()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

//nolint:goconst
func Test_operator_getAtomicJob(t *testing.T) {
	t.Parallel()
	type fields struct {
		jobs *atomic.Pointer[map[string]string]
	}
	type want struct {
		want map[string]string
	}
	type test struct {
		fields     fields
		want       want
		checkFunc  func(want, map[string]string) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		name       string
	}
	defaultCheckFunc := func(w want, got map[string]string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name:   "get nil when atomic has no resource",
			fields: fields{},
			want: want{
				want: nil,
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(t *testing.T) {
				t.Helper()
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
			},
		},
		{
			name: "get jobs when jobs has resource",
			fields: fields{
				jobs: func() *atomic.Pointer[map[string]string] {
					ap := atomic.Pointer[map[string]string]{}
					m := map[string]string{
						"scenario-insert": "default",
						"scenario-search": "default",
					}
					ap.Store(&m)
					return &ap
				}(),
			},
			want: want{
				want: map[string]string{
					"scenario-insert": "default",
					"scenario-search": "default",
				},
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(t *testing.T) {
				t.Helper()
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := new(operator)
			storeInto(&o.jobs, test.fields.jobs)

			got := o.getAtomicJob()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

//nolint:goconst,maintidx // repeated fixture literals and table volume are inherent to this table-driven test
func Test_operator_jobReconcile(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context //nolint:containedctx // table-driven test args struct, ctx passed through to the function under test
		jobList map[string][]k8s.Job
	}
	type fields struct {
		ctrl               k8s.Controller
		scenarios          *atomic.Pointer[map[string]*scenario]
		benchjobs          *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
		jobs               *atomic.Pointer[map[string]string]
		jobNamespace       string
		jobImageRepository string
		jobImageTag        string
		jobImagePullPolicy string
	}
	type want struct {
		want map[string]string
	}
	type test struct {
		fields     fields
		args       args
		want       want
		checkFunc  func(want, map[string]string) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
	}
	defaultCheckFunc := func(w want, got map[string]string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success when the length of jobList is 0.",
				args: args{
					ctx:     ctx,
					jobList: map[string][]k8s.Job{},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios:          &atomic.Pointer[map[string]*scenario]{},
					benchjobs:          &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{},
					jobs:               &atomic.Pointer[map[string]string]{},
				},
				want: want{
					want: map[string]string{},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success with new job whose namespace is same as jobNamespace and deleted job by etcd",
				args: args{
					ctx: ctx,
					jobList: map[string][]k8s.Job{
						"scenario-insert": {
							{
								Name:      "scenario-insert",
								Namespace: "default",
								Status: k8s.JobStatus{
									Active: 1,
								},
							},
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios:          &atomic.Pointer[map[string]*scenario]{},
					jobs: func() *atomic.Pointer[map[string]string] {
						ap := atomic.Pointer[map[string]string]{}
						m := map[string]string{
							"scenario-completed-insert": "default",
						}
						ap.Store(&m)
						return &ap
					}(),
					benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
						ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
						m := map[string]*v1.ValdBenchmarkJob{
							"scenario-insert": {
								Spec: v1.BenchmarkJobSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									JobType: "insert",
									InsertConfig: &config.InsertConfig{
										SkipStrictExistCheck: false,
										Timestamp:            "",
									},
								},
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					want: map[string]string{
						"scenario-insert": "default",
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success with completed job whose namespace is same as jobNamespace",
				args: args{
					ctx: ctx,
					jobList: map[string][]k8s.Job{
						"scenario-insert": {
							{
								Name:      "scenario-completed-insert",
								Namespace: "default",
								Status: k8s.JobStatus{
									Active:    0,
									Succeeded: 1,
									CompletionTime: func() *k8s.Time {
										t := &k8s.Time{
											Time: time.Now(),
										}
										return t
									}(),
								},
							},
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios:          &atomic.Pointer[map[string]*scenario]{},
					jobs: func() *atomic.Pointer[map[string]string] {
						ap := atomic.Pointer[map[string]string]{}
						m := map[string]string{
							"scenario-completed-insert": "default",
						}
						ap.Store(&m)
						return &ap
					}(),
					benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
						ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
						m := map[string]*v1.ValdBenchmarkJob{
							"scenario-insert": {
								Spec: v1.BenchmarkJobSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									JobType: "insert",
									InsertConfig: &config.InsertConfig{
										SkipStrictExistCheck: false,
										Timestamp:            "",
									},
								},
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					want: map[string]string{
						"scenario-completed-insert": "default",
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success with job whose namespace is not same as jobNamespace",
				args: args{
					ctx: ctx,
					jobList: map[string][]k8s.Job{
						"scenario-insert": {
							{
								Name:      "scenario-insert",
								Namespace: "benchmark",
								Status: k8s.JobStatus{
									Active: 1,
								},
							},
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios:          &atomic.Pointer[map[string]*scenario]{},
					benchjobs:          &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{},
					jobs:               &atomic.Pointer[map[string]string]{},
					ctrl:               nil,
				},
				want: want{
					want: map[string]string{},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &operator{
				jobNamespace:       test.fields.jobNamespace,
				jobImageRepository: test.fields.jobImageRepository,
				jobImageTag:        test.fields.jobImageTag,
				jobImagePullPolicy: test.fields.jobImagePullPolicy,
				ctrl:               test.fields.ctrl,
			}
			storeInto(&o.benchjobs, test.fields.benchjobs)
			storeInto(&o.jobs, test.fields.jobs)

			o.jobReconcile(test.args.ctx, test.args.jobList)
			got := o.getAtomicJob()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

//nolint:maintidx // table volume is inherent to this table-driven test
func Test_operator_benchJobReconcile(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx          context.Context //nolint:containedctx // table-driven test args struct, ctx passed through to the function under test
		benchJobList map[string]v1.ValdBenchmarkJob
	}
	type fields struct {
		ctrl               k8s.Controller
		scenarios          *atomic.Pointer[map[string]*scenario]
		benchjobs          *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
		jobs               *atomic.Pointer[map[string]string]
		jobNamespace       string
		jobImageRepository string
		jobImageTag        string
		jobImagePullPolicy string
		configMapName      string
	}
	type want struct {
		scenarios map[string]*scenario
		benchjobs map[string]*v1.ValdBenchmarkJob
	}
	type test struct {
		fields     fields
		args       args
		want       want
		checkFunc  func(want, map[string]*scenario, map[string]*v1.ValdBenchmarkJob) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
	}
	defaultCheckFunc := func(w want, gotS map[string]*scenario, gotJ map[string]*v1.ValdBenchmarkJob) error {
		if !reflect.DeepEqual(w.scenarios, gotS) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotS, w.scenarios)
		}
		if !reflect.DeepEqual(w.benchjobs, gotJ) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotJ, w.benchjobs)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success when benchJobList is empty",
				args: args{
					ctx:          ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios:          &atomic.Pointer[map[string]*scenario]{},
					benchjobs:          &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{},
					jobs:               &atomic.Pointer[map[string]string]{},
					ctrl:               nil,
				},
				want: want{
					benchjobs: map[string]*v1.ValdBenchmarkJob{},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success when benchJobList has new benchmark Job with owner reference (reconcile after submitted scenario)",
				args: args{
					ctx: ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: k8s.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []k8s.OwnerReference{
									{
										Name: "scenario",
									},
								},
								Generation: 1,
							},
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								JobType: "insert",
								InsertConfig: &config.InsertConfig{
									SkipStrictExistCheck: false,
									Timestamp:            "",
								},
							},
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					configMapName:      "benchmark-job-test-configmap",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									Name:       "scenario",
									Namespace:  "default",
									Generation: 1,
									Spec: v1.ValdBenchmarkScenarioSpec{
										Target: &v1.BenchmarkTarget{
											Host: "localhost",
											Port: 8080,
										},
										Dataset: &v1.BenchmarkDataset{
											Name:    "fashion-minsit",
											Group:   "train",
											Indexes: 10000,
											Range: &config.BenchmarkDatasetRange{
												Start: 0,
												End:   10000,
											},
											URL: "",
										},
										Jobs: []*v1.BenchmarkJobSpec{
											{
												JobType: "insert",
												InsertConfig: &config.InsertConfig{
													SkipStrictExistCheck: false,
													Timestamp:            "",
												},
											},
										},
									},
									Status: v1.BenchmarkScenarioAvailable,
								},
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					benchjobs: &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{},
					jobs:      &atomic.Pointer[map[string]string]{},
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					scenarios: map[string]*scenario{
						"scenario": {
							Crd: &v1.ValdBenchmarkScenario{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 1,
								Spec: v1.ValdBenchmarkScenarioSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									Jobs: []*v1.BenchmarkJobSpec{
										{
											JobType: "insert",
											InsertConfig: &config.InsertConfig{
												SkipStrictExistCheck: false,
												Timestamp:            "",
											},
										},
									},
								},
								Status: v1.BenchmarkScenarioAvailable,
							},
							BenchJobStatus: map[string]v1.BenchmarkJobStatus{
								"scenario-insert": "",
							},
						},
					},
					benchjobs: map[string]*v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: k8s.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []k8s.OwnerReference{
									{
										Name: "scenario",
									},
								},
								Generation: 1,
							},
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								JobType: "insert",
								InsertConfig: &config.InsertConfig{
									SkipStrictExistCheck: false,
									Timestamp:            "",
								},
							},
							Status: v1.BenchmarkJobHealthy,
						},
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success when benchJobList has updated benchmark Job with owner reference (reconcile after updated scenario)",
				args: args{
					ctx: ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: k8s.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []k8s.OwnerReference{
									{
										Name: "scenario",
									},
								},
								Generation: 2,
							},
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 20000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   20000,
									},
									URL: "",
								},
								JobType: "insert",
								InsertConfig: &config.InsertConfig{
									SkipStrictExistCheck: false,
									Timestamp:            "",
								},
							},
							Status: v1.BenchmarkJobAvailable,
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					configMapName:      "benchmark-job-test-configmap",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									Name:       "scenario",
									Namespace:  "default",
									Generation: 2,
									Spec: v1.ValdBenchmarkScenarioSpec{
										Target: &v1.BenchmarkTarget{
											Host: "localhost",
											Port: 8080,
										},
										Dataset: &v1.BenchmarkDataset{
											Name:    "fashion-minsit",
											Group:   "train",
											Indexes: 20000,
											Range: &config.BenchmarkDatasetRange{
												Start: 0,
												End:   20000,
											},
											URL: "",
										},
										Jobs: []*v1.BenchmarkJobSpec{
											{
												JobType: "insert",
												InsertConfig: &config.InsertConfig{
													SkipStrictExistCheck: false,
													Timestamp:            "",
												},
											},
										},
									},
									Status: v1.BenchmarkScenarioAvailable,
								},
								BenchJobStatus: map[string]v1.BenchmarkJobStatus{
									"scenario-insert": v1.BenchmarkJobAvailable,
								},
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
						ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
						m := map[string]*v1.ValdBenchmarkJob{
							"scenario-insert": {
								ObjectMeta: k8s.ObjectMeta{
									Name:      "scenario-insert",
									Namespace: "default",
									OwnerReferences: []k8s.OwnerReference{
										{
											Name: "scenario",
										},
									},
									Generation: 1,
								},
								Spec: v1.BenchmarkJobSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									JobType: "insert",
									InsertConfig: &config.InsertConfig{
										SkipStrictExistCheck: false,
										Timestamp:            "",
									},
								},
								Status: v1.BenchmarkJobAvailable,
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					jobs: &atomic.Pointer[map[string]string]{},
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					scenarios: map[string]*scenario{
						"scenario": {
							Crd: &v1.ValdBenchmarkScenario{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 2,
								Spec: v1.ValdBenchmarkScenarioSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 20000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   20000,
										},
										URL: "",
									},
									Jobs: []*v1.BenchmarkJobSpec{
										{
											JobType: "insert",
											InsertConfig: &config.InsertConfig{
												SkipStrictExistCheck: false,
												Timestamp:            "",
											},
										},
									},
								},
								Status: v1.BenchmarkScenarioAvailable,
							},
							BenchJobStatus: map[string]v1.BenchmarkJobStatus{
								"scenario-insert": v1.BenchmarkJobAvailable,
							},
						},
					},
					benchjobs: map[string]*v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: k8s.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []k8s.OwnerReference{
									{
										Name: "scenario",
									},
								},
								Generation: 2,
							},
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 20000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   20000,
									},
									URL: "",
								},
								JobType: "insert",
								InsertConfig: &config.InsertConfig{
									SkipStrictExistCheck: false,
									Timestamp:            "",
								},
							},
							Status: v1.BenchmarkJobAvailable,
						},
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success when benchJobList has updated benchmark Job with owner reference (reconcile after updated job status)",
				args: args{
					ctx: ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: k8s.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []k8s.OwnerReference{
									{
										Name: "scenario",
									},
								},
								Generation: 1,
							},
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								JobType: "insert",
								InsertConfig: &config.InsertConfig{
									SkipStrictExistCheck: false,
									Timestamp:            "",
								},
							},
							Status: v1.BenchmarkJobAvailable,
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									Name:       "scenario",
									Namespace:  "default",
									Generation: 1,
									Spec: v1.ValdBenchmarkScenarioSpec{
										Target: &v1.BenchmarkTarget{
											Host: "localhost",
											Port: 8080,
										},
										Dataset: &v1.BenchmarkDataset{
											Name:    "fashion-minsit",
											Group:   "train",
											Indexes: 10000,
											Range: &config.BenchmarkDatasetRange{
												Start: 0,
												End:   10000,
											},
											URL: "",
										},
										Jobs: []*v1.BenchmarkJobSpec{
											{
												JobType: "insert",
												InsertConfig: &config.InsertConfig{
													SkipStrictExistCheck: false,
													Timestamp:            "",
												},
											},
										},
									},
									Status: v1.BenchmarkScenarioAvailable,
								},
								BenchJobStatus: map[string]v1.BenchmarkJobStatus{
									"scenario-insert": "",
								},
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
						ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
						m := map[string]*v1.ValdBenchmarkJob{
							"scenario-insert": {
								ObjectMeta: k8s.ObjectMeta{
									Name:      "scenario-insert",
									Namespace: "default",
									OwnerReferences: []k8s.OwnerReference{
										{
											Name: "scenario",
										},
									},
									Generation: 1,
								},
								Spec: v1.BenchmarkJobSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									JobType: "insert",
									InsertConfig: &config.InsertConfig{
										SkipStrictExistCheck: false,
										Timestamp:            "",
									},
								},
								Status: "",
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					jobs: &atomic.Pointer[map[string]string]{},
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					scenarios: map[string]*scenario{
						"scenario": {
							Crd: &v1.ValdBenchmarkScenario{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 1,
								Spec: v1.ValdBenchmarkScenarioSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									Jobs: []*v1.BenchmarkJobSpec{
										{
											JobType: "insert",
											InsertConfig: &config.InsertConfig{
												SkipStrictExistCheck: false,
												Timestamp:            "",
											},
										},
									},
								},
								Status: v1.BenchmarkScenarioAvailable,
							},
							BenchJobStatus: map[string]v1.BenchmarkJobStatus{
								"scenario-insert": v1.BenchmarkJobAvailable,
							},
						},
					},
					benchjobs: map[string]*v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: k8s.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []k8s.OwnerReference{
									{
										Name: "scenario",
									},
								},
								Generation: 1,
							},
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								JobType: "insert",
								InsertConfig: &config.InsertConfig{
									SkipStrictExistCheck: false,
									Timestamp:            "",
								},
							},
							Status: v1.BenchmarkJobAvailable,
						},
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success when benchJobList has new benchmark Job with owner reference and benchJob has deleted job",
				args: args{
					ctx: ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: k8s.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []k8s.OwnerReference{
									{
										Name: "scenario",
									},
								},
								Generation: 1,
							},
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								JobType: "insert",
								InsertConfig: &config.InsertConfig{
									SkipStrictExistCheck: false,
									Timestamp:            "",
								},
							},
							Status: v1.BenchmarkJobAvailable,
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									Name:       "scenario",
									Namespace:  "default",
									Generation: 1,
									Spec: v1.ValdBenchmarkScenarioSpec{
										Target: &v1.BenchmarkTarget{
											Host: "localhost",
											Port: 8080,
										},
										Dataset: &v1.BenchmarkDataset{
											Name:    "fashion-minsit",
											Group:   "train",
											Indexes: 10000,
											Range: &config.BenchmarkDatasetRange{
												Start: 0,
												End:   10000,
											},
											URL: "",
										},
										Jobs: []*v1.BenchmarkJobSpec{
											{
												JobType: "insert",
												InsertConfig: &config.InsertConfig{
													SkipStrictExistCheck: false,
													Timestamp:            "",
												},
											},
										},
									},
									Status: v1.BenchmarkScenarioAvailable,
								},
								BenchJobStatus: map[string]v1.BenchmarkJobStatus{
									"scenario-insert": "",
								},
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
						ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
						m := map[string]*v1.ValdBenchmarkJob{
							"scenario-insert": {
								ObjectMeta: k8s.ObjectMeta{
									Name:      "scenario-insert",
									Namespace: "default",
									OwnerReferences: []k8s.OwnerReference{
										{
											Name: "scenario",
										},
									},
									Generation: 1,
								},
								Spec: v1.BenchmarkJobSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									JobType: "insert",
									InsertConfig: &config.InsertConfig{
										SkipStrictExistCheck: false,
										Timestamp:            "",
									},
								},
								Status: "",
							},
							"scenario-deleted-insert": {
								ObjectMeta: k8s.ObjectMeta{
									Name:      "scenario-deleted-insert",
									Namespace: "default",
									OwnerReferences: []k8s.OwnerReference{
										{
											Name: "scenario-deleted",
										},
									},
									Generation: 1,
								},
								Spec: v1.BenchmarkJobSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									JobType: "insert",
									InsertConfig: &config.InsertConfig{
										SkipStrictExistCheck: false,
										Timestamp:            "",
									},
								},
								Status: v1.BenchmarkJobCompleted,
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					jobs: &atomic.Pointer[map[string]string]{},
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					scenarios: map[string]*scenario{
						"scenario": {
							Crd: &v1.ValdBenchmarkScenario{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 1,
								Spec: v1.ValdBenchmarkScenarioSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									Jobs: []*v1.BenchmarkJobSpec{
										{
											JobType: "insert",
											InsertConfig: &config.InsertConfig{
												SkipStrictExistCheck: false,
												Timestamp:            "",
											},
										},
									},
								},
								Status: v1.BenchmarkScenarioAvailable,
							},
							BenchJobStatus: map[string]v1.BenchmarkJobStatus{
								"scenario-insert": v1.BenchmarkJobAvailable,
							},
						},
					},
					benchjobs: map[string]*v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: k8s.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []k8s.OwnerReference{
									{
										Name: "scenario",
									},
								},
								Generation: 1,
							},
							Spec: v1.BenchmarkJobSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								JobType: "insert",
								InsertConfig: &config.InsertConfig{
									SkipStrictExistCheck: false,
									Timestamp:            "",
								},
							},
							Status: v1.BenchmarkJobAvailable,
						},
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &operator{
				jobNamespace:       test.fields.jobNamespace,
				jobImageRepository: test.fields.jobImageRepository,
				jobImageTag:        test.fields.jobImageTag,
				jobImagePullPolicy: test.fields.jobImagePullPolicy,
				configMapName:      test.fields.configMapName,
				ctrl:               test.fields.ctrl,
			}
			storeInto(&o.scenarios, test.fields.scenarios)
			storeInto(&o.benchjobs, test.fields.benchjobs)
			storeInto(&o.jobs, test.fields.jobs)

			o.benchJobReconcile(test.args.ctx, test.args.benchJobList)
			gotS := o.getAtomicScenario()
			gotJ := o.getAtomicBenchJob()
			if err := checkFunc(test.want, gotS, gotJ); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

//nolint:goconst,maintidx // repeated fixture literals and table volume are inherent to this table-driven test
func Test_operator_benchScenarioReconcile(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx          context.Context //nolint:containedctx // table-driven test args struct, ctx passed through to the function under test
		scenarioList map[string]v1.ValdBenchmarkScenario
	}
	type fields struct {
		ctrl               k8s.Controller
		scenarios          *atomic.Pointer[map[string]*scenario]
		benchjobs          *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
		jobs               *atomic.Pointer[map[string]string]
		jobNamespace       string
		jobImageRepository string
		jobImageTag        string
		jobImagePullPolicy string
	}
	type want struct {
		want map[string]*scenario
	}
	type test struct {
		fields     fields
		args       args
		want       want
		checkFunc  func(want, map[string]*scenario) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
	}
	defaultCheckFunc := func(w want, got map[string]*scenario) error {
		if len(w.want) != len(got) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		for k, ws := range w.want {
			gs, ok := got[k]
			if !ok {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
			}
			// check CRD
			if !reflect.DeepEqual(ws.Crd, gs.Crd) {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
			}
			// check benchJobStatus
			if len(ws.BenchJobStatus) != len(gs.BenchJobStatus) {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
			}
			for k, v := range gs.BenchJobStatus {
				sk := strings.Split(k, "-")
				wk := strings.Join(sk[:len(sk)-1], "-")
				if v != ws.BenchJobStatus[wk] {
					return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
				}
			}
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success with scenarioList is empty",
				args: args{
					ctx:          ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios:          &atomic.Pointer[map[string]*scenario]{},
					benchjobs:          &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{},
					jobs:               &atomic.Pointer[map[string]string]{},
					ctrl:               nil,
				},
				want: want{
					want: map[string]*scenario{},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success with scenarioList has new scenario with no scenario has been applied yet.",
				args: args{
					ctx: ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{
						"scenario": {
							ObjectMeta: k8s.ObjectMeta{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 1,
							},
							Spec: v1.ValdBenchmarkScenarioSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								Jobs: []*v1.BenchmarkJobSpec{
									{
										JobType: "insert",
										InsertConfig: &config.InsertConfig{
											SkipStrictExistCheck: false,
											Timestamp:            "",
										},
									},
									{
										JobType: "search",
										SearchConfig: &config.SearchConfig{
											Epsilon:              0.1,
											Radius:               -1,
											Num:                  10,
											MinNum:               5,
											EnableLinearSearch:   false,
											AggregationAlgorithm: "",
										},
									},
								},
							},
							Status: v1.BenchmarkScenarioAvailable,
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios:          &atomic.Pointer[map[string]*scenario]{},
					benchjobs:          &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{},
					jobs:               &atomic.Pointer[map[string]string]{},
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					want: map[string]*scenario{
						"scenario": {
							Crd: &v1.ValdBenchmarkScenario{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 1,
								Spec: v1.ValdBenchmarkScenarioSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									Jobs: []*v1.BenchmarkJobSpec{
										{
											JobType: "insert",
											InsertConfig: &config.InsertConfig{
												SkipStrictExistCheck: false,
												Timestamp:            "",
											},
										},
										{
											JobType: "search",
											SearchConfig: &config.SearchConfig{
												Epsilon:              0.1,
												Radius:               -1,
												Num:                  10,
												MinNum:               5,
												EnableLinearSearch:   false,
												AggregationAlgorithm: "",
											},
										},
									},
								},
								Status: v1.BenchmarkScenarioHealthy,
							},
							BenchJobStatus: map[string]v1.BenchmarkJobStatus{
								"scenario-insert": v1.BenchmarkJobNotReady,
								"scenario-search": v1.BenchmarkJobNotReady,
							},
						},
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success with scenarioList has only status updated scenario.",
				args: args{
					ctx: ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{
						"scenario": {
							ObjectMeta: k8s.ObjectMeta{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 1,
							},
							Spec: v1.ValdBenchmarkScenarioSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								Jobs: []*v1.BenchmarkJobSpec{
									{
										JobType: "insert",
										InsertConfig: &config.InsertConfig{
											SkipStrictExistCheck: false,
											Timestamp:            "",
										},
									},
									{
										JobType: "search",
										SearchConfig: &config.SearchConfig{
											Epsilon:              0.1,
											Radius:               -1,
											Num:                  10,
											MinNum:               5,
											EnableLinearSearch:   false,
											AggregationAlgorithm: "",
										},
									},
								},
							},
							Status: v1.BenchmarkScenarioAvailable,
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									Name:       "scenario",
									Namespace:  "default",
									Generation: 1,
									Spec: v1.ValdBenchmarkScenarioSpec{
										Target: &v1.BenchmarkTarget{
											Host: "localhost",
											Port: 8080,
										},
										Dataset: &v1.BenchmarkDataset{
											Name:    "fashion-minsit",
											Group:   "train",
											Indexes: 10000,
											Range: &config.BenchmarkDatasetRange{
												Start: 0,
												End:   10000,
											},
											URL: "",
										},
										Jobs: []*v1.BenchmarkJobSpec{
											{
												JobType: "insert",
												InsertConfig: &config.InsertConfig{
													SkipStrictExistCheck: false,
													Timestamp:            "",
												},
											},
											{
												JobType: "search",
												SearchConfig: &config.SearchConfig{
													Epsilon:              0.1,
													Radius:               -1,
													Num:                  10,
													MinNum:               5,
													EnableLinearSearch:   false,
													AggregationAlgorithm: "",
												},
											},
										},
									},
									Status: v1.BenchmarkScenarioHealthy,
								},
								BenchJobStatus: map[string]v1.BenchmarkJobStatus{
									"scenario-insert-1234567890": v1.BenchmarkJobNotReady,
									"scenario-search-1234567891": v1.BenchmarkJobNotReady,
								},
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					benchjobs: &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{},
					jobs:      &atomic.Pointer[map[string]string]{},
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					want: map[string]*scenario{
						"scenario": {
							Crd: &v1.ValdBenchmarkScenario{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 1,
								Spec: v1.ValdBenchmarkScenarioSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									Jobs: []*v1.BenchmarkJobSpec{
										{
											JobType: "insert",
											InsertConfig: &config.InsertConfig{
												SkipStrictExistCheck: false,
												Timestamp:            "",
											},
										},
										{
											JobType: "search",
											SearchConfig: &config.SearchConfig{
												Epsilon:              0.1,
												Radius:               -1,
												Num:                  10,
												MinNum:               5,
												EnableLinearSearch:   false,
												AggregationAlgorithm: "",
											},
										},
									},
								},
								Status: v1.BenchmarkScenarioAvailable,
							},
							BenchJobStatus: map[string]v1.BenchmarkJobStatus{
								"scenario-insert": v1.BenchmarkJobNotReady,
								"scenario-search": v1.BenchmarkJobNotReady,
							},
						},
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success with scenarioList has updated scenario when job is already running",
				args: args{
					ctx: ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{
						"scenario": {
							ObjectMeta: k8s.ObjectMeta{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 2,
							},
							Spec: v1.ValdBenchmarkScenarioSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 20000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   20000,
									},
									URL: "",
								},
								Jobs: []*v1.BenchmarkJobSpec{
									{
										JobType: "insert",
										InsertConfig: &config.InsertConfig{
											SkipStrictExistCheck: false,
											Timestamp:            "",
										},
									},
									{
										JobType: "search",
										SearchConfig: &config.SearchConfig{
											Epsilon:              0.1,
											Radius:               -1,
											Num:                  10,
											MinNum:               5,
											EnableLinearSearch:   false,
											AggregationAlgorithm: "",
										},
									},
								},
							},
							Status: v1.BenchmarkScenarioAvailable,
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									Name:       "scenario",
									Namespace:  "default",
									Generation: 1,
									Spec: v1.ValdBenchmarkScenarioSpec{
										Target: &v1.BenchmarkTarget{
											Host: "localhost",
											Port: 8080,
										},
										Dataset: &v1.BenchmarkDataset{
											Name:    "fashion-minsit",
											Group:   "train",
											Indexes: 10000,
											Range: &config.BenchmarkDatasetRange{
												Start: 0,
												End:   10000,
											},
											URL: "",
										},
										Jobs: []*v1.BenchmarkJobSpec{
											{
												JobType: "insert",
												InsertConfig: &config.InsertConfig{
													SkipStrictExistCheck: false,
													Timestamp:            "",
												},
											},
											{
												JobType: "search",
												SearchConfig: &config.SearchConfig{
													Epsilon:              0.1,
													Radius:               -1,
													Num:                  10,
													MinNum:               5,
													EnableLinearSearch:   false,
													AggregationAlgorithm: "",
												},
											},
										},
									},
									Status: v1.BenchmarkScenarioAvailable,
								},
								BenchJobStatus: map[string]v1.BenchmarkJobStatus{
									"scenario-insert": v1.BenchmarkJobAvailable,
									"scenario-search": v1.BenchmarkJobAvailable,
								},
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					benchjobs: &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{},
					jobs:      &atomic.Pointer[map[string]string]{},
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					want: map[string]*scenario{
						"scenario": {
							Crd: &v1.ValdBenchmarkScenario{
								Name:       "scenario",
								Namespace:  "default",
								Generation: 2,
								Spec: v1.ValdBenchmarkScenarioSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 20000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   20000,
										},
										URL: "",
									},
									Jobs: []*v1.BenchmarkJobSpec{
										{
											JobType: "insert",
											InsertConfig: &config.InsertConfig{
												SkipStrictExistCheck: false,
												Timestamp:            "",
											},
										},
										{
											JobType: "search",
											SearchConfig: &config.SearchConfig{
												Epsilon:              0.1,
												Radius:               -1,
												Num:                  10,
												MinNum:               5,
												EnableLinearSearch:   false,
												AggregationAlgorithm: "",
											},
										},
									},
								},
								Status: v1.BenchmarkScenarioAvailable,
							},
							BenchJobStatus: map[string]v1.BenchmarkJobStatus{
								"scenario-insert": v1.BenchmarkJobNotReady,
								"scenario-search": v1.BenchmarkJobNotReady,
							},
						},
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
			return test{
				name: "success with scenarioList has another scenario when scenario is already running",
				args: args{
					ctx: ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{
						"scenario-v2": {
							ObjectMeta: k8s.ObjectMeta{
								Name:       "scenario-v2",
								Namespace:  "default",
								Generation: 1,
							},
							Spec: v1.ValdBenchmarkScenarioSpec{
								Target: &v1.BenchmarkTarget{
									Host: "localhost",
									Port: 8080,
								},
								Dataset: &v1.BenchmarkDataset{
									Name:    "fashion-minsit",
									Group:   "train",
									Indexes: 10000,
									Range: &config.BenchmarkDatasetRange{
										Start: 0,
										End:   10000,
									},
									URL: "",
								},
								Jobs: []*v1.BenchmarkJobSpec{
									{
										JobType: "insert",
										InsertConfig: &config.InsertConfig{
											SkipStrictExistCheck: false,
											Timestamp:            "",
										},
									},
									{
										JobType: "search",
										SearchConfig: &config.SearchConfig{
											Epsilon:              0.1,
											Radius:               -1,
											Num:                  10,
											MinNum:               5,
											EnableLinearSearch:   false,
											AggregationAlgorithm: "",
										},
									},
								},
							},
							Status: v1.BenchmarkScenarioAvailable,
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImageRepository: "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									Name:       "scenario",
									Namespace:  "default",
									Generation: 1,
									Spec: v1.ValdBenchmarkScenarioSpec{
										Target: &v1.BenchmarkTarget{
											Host: "localhost",
											Port: 8080,
										},
										Dataset: &v1.BenchmarkDataset{
											Name:    "fashion-minsit",
											Group:   "train",
											Indexes: 10000,
											Range: &config.BenchmarkDatasetRange{
												Start: 0,
												End:   10000,
											},
											URL: "",
										},
										Jobs: []*v1.BenchmarkJobSpec{
											{
												JobType: "insert",
												InsertConfig: &config.InsertConfig{
													SkipStrictExistCheck: false,
													Timestamp:            "",
												},
											},
											{
												JobType: "search",
												SearchConfig: &config.SearchConfig{
													Epsilon:              0.1,
													Radius:               -1,
													Num:                  10,
													MinNum:               5,
													EnableLinearSearch:   false,
													AggregationAlgorithm: "",
												},
											},
										},
									},
									Status: v1.BenchmarkScenarioAvailable,
								},
								BenchJobStatus: map[string]v1.BenchmarkJobStatus{
									"scenario-insert": v1.BenchmarkJobAvailable,
									"scenario-search": v1.BenchmarkJobAvailable,
								},
							},
						}
						ap.Store(&m)
						return &ap
					}(),
					benchjobs: &atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{},
					jobs:      &atomic.Pointer[map[string]string]{},
					ctrl: &mockCtrl{
						StartFunc: func(ctx context.Context) (<-chan error, error) {
							return nil, nil
						},
						GetManagerFunc: func() k8s.Manager {
							m := &mock.MockManager{
								Manager: &mock.MockManager{},
							}
							return m
						},
					},
				},
				want: want{
					want: map[string]*scenario{
						"scenario-v2": {
							Crd: &v1.ValdBenchmarkScenario{
								Name:       "scenario-v2",
								Namespace:  "default",
								Generation: 1,
								Spec: v1.ValdBenchmarkScenarioSpec{
									Target: &v1.BenchmarkTarget{
										Host: "localhost",
										Port: 8080,
									},
									Dataset: &v1.BenchmarkDataset{
										Name:    "fashion-minsit",
										Group:   "train",
										Indexes: 10000,
										Range: &config.BenchmarkDatasetRange{
											Start: 0,
											End:   10000,
										},
										URL: "",
									},
									Jobs: []*v1.BenchmarkJobSpec{
										{
											JobType: "insert",
											InsertConfig: &config.InsertConfig{
												SkipStrictExistCheck: false,
												Timestamp:            "",
											},
										},
										{
											JobType: "search",
											SearchConfig: &config.SearchConfig{
												Epsilon:              0.1,
												Radius:               -1,
												Num:                  10,
												MinNum:               5,
												EnableLinearSearch:   false,
												AggregationAlgorithm: "",
											},
										},
									},
								},
								Status: v1.BenchmarkScenarioHealthy,
							},
							BenchJobStatus: map[string]v1.BenchmarkJobStatus{
								"scenario-v2-insert": v1.BenchmarkJobNotReady,
								"scenario-v2-search": v1.BenchmarkJobNotReady,
							},
						},
					},
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T, args args) {
					t.Helper()
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					cancel()
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &operator{
				jobNamespace:       test.fields.jobNamespace,
				jobImageRepository: test.fields.jobImageRepository,
				jobImageTag:        test.fields.jobImageTag,
				jobImagePullPolicy: test.fields.jobImagePullPolicy,
				ctrl:               test.fields.ctrl,
			}
			storeInto(&o.scenarios, test.fields.scenarios)
			storeInto(&o.benchjobs, test.fields.benchjobs)
			storeInto(&o.jobs, test.fields.jobs)

			o.benchScenarioReconcile(test.args.ctx, test.args.scenarioList)
			got := o.getAtomicScenario()
			t.Log(got["scenario"])
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

//nolint:goconst,maintidx // repeated fixture literals and table volume are inherent to this table-driven test
func Test_operator_checkAtomics(t *testing.T) {
	t.Parallel()
	type fields struct {
		scenarios *atomic.Pointer[map[string]*scenario]
		benchjobs *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
		jobs      *atomic.Pointer[map[string]string]
	}
	type want struct {
		err error
	}
	type test struct {
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		name       string
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	defaultScenarioMap := map[string]*scenario{
		"scenario": {
			Crd: &v1.ValdBenchmarkScenario{
				Name:      "scenario",
				Namespace: "default",
				Spec: v1.ValdBenchmarkScenarioSpec{
					Target: &v1.BenchmarkTarget{
						Host: "localhost",
						Port: 8080,
					},
					Dataset: &v1.BenchmarkDataset{
						Name:    "fashion-minsit",
						Group:   "train",
						Indexes: 10000,
						Range: &config.BenchmarkDatasetRange{
							Start: 0,
							End:   10000,
						},
						URL: "",
					},
					Jobs: []*v1.BenchmarkJobSpec{
						{
							JobType: "insert",
							InsertConfig: &config.InsertConfig{
								SkipStrictExistCheck: false,
								Timestamp:            "",
							},
						},
						{
							JobType: "search",
							SearchConfig: &config.SearchConfig{
								Epsilon:              0.1,
								Radius:               -1,
								Num:                  10,
								MinNum:               10,
								Timeout:              "10s",
								EnableLinearSearch:   false,
								AggregationAlgorithm: "",
							},
						},
						{
							JobType: "update",
							UpdateConfig: &config.UpdateConfig{
								SkipStrictExistCheck: false,
								Timestamp:            "",
							},
						},
					},
				},
				Status: v1.BenchmarkScenarioHealthy,
			},
			BenchJobStatus: map[string]v1.BenchmarkJobStatus{
				"scenario-insert": v1.BenchmarkJobCompleted,
				"scenario-search": v1.BenchmarkJobAvailable,
				"scenario-update": v1.BenchmarkJobAvailable,
			},
		},
	}
	defaultBenchJobMap := map[string]*v1.ValdBenchmarkJob{
		"scenario-insert": {
			ObjectMeta: k8s.ObjectMeta{
				Name:      "scenario-insert",
				Namespace: "default",
				OwnerReferences: []k8s.OwnerReference{
					{
						Kind: ScenarioKind,
						Name: "scenario",
					},
				},
			},
			Spec: v1.BenchmarkJobSpec{
				Target: &v1.BenchmarkTarget{
					Host: "localhost",
					Port: 8080,
				},
				Dataset: &v1.BenchmarkDataset{
					Name:    "fashion-minsit",
					Group:   "train",
					Indexes: 10000,
					Range: &config.BenchmarkDatasetRange{
						Start: 0,
						End:   10000,
					},
					URL: "",
				},
				JobType: "insert",
				InsertConfig: &config.InsertConfig{
					SkipStrictExistCheck: false,
					Timestamp:            "",
				},
			},
			Status: v1.BenchmarkJobCompleted,
		},
		"scenario-search": {
			ObjectMeta: k8s.ObjectMeta{
				Name:      "scenario-search",
				Namespace: "default",
				OwnerReferences: []k8s.OwnerReference{
					{
						Kind: ScenarioKind,
						Name: "scenario",
					},
				},
			},
			Spec: v1.BenchmarkJobSpec{
				Target: &v1.BenchmarkTarget{
					Host: "localhost",
					Port: 8080,
				},
				Dataset: &v1.BenchmarkDataset{
					Name:    "fashion-minsit",
					Group:   "train",
					Indexes: 10000,
					Range: &config.BenchmarkDatasetRange{
						Start: 0,
						End:   10000,
					},
					URL: "",
				},
				JobType: "search",
				SearchConfig: &config.SearchConfig{
					Epsilon:              0.1,
					Radius:               -1,
					Num:                  10,
					MinNum:               10,
					Timeout:              "10s",
					EnableLinearSearch:   false,
					AggregationAlgorithm: "",
				},
			},
			Status: v1.BenchmarkJobAvailable,
		},
		"scenario-update": {
			ObjectMeta: k8s.ObjectMeta{
				Name:      "scenario-update",
				Namespace: "default",
				OwnerReferences: []k8s.OwnerReference{
					{
						Kind: ScenarioKind,
						Name: "scenario",
					},
				},
			},
			Spec: v1.BenchmarkJobSpec{
				Target: &v1.BenchmarkTarget{
					Host: "localhost",
					Port: 8080,
				},
				Dataset: &v1.BenchmarkDataset{
					Name:    "fashion-minsit",
					Group:   "train",
					Indexes: 10000,
					Range: &config.BenchmarkDatasetRange{
						Start: 0,
						End:   10000,
					},
					URL: "",
				},
				JobType: "update",
				UpdateConfig: &config.UpdateConfig{
					SkipStrictExistCheck: false,
					Timestamp:            "",
				},
			},
			Status: v1.BenchmarkJobAvailable,
		},
	}
	defaultJobMap := map[string]string{
		"scenario-insert": "default",
		"scenario-search": "default",
		"scenario-update": "default",
	}
	tests := []test{
		func() test {
			return test{
				name: "return nil with no mismatch atomics",
				fields: fields{
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						ap.Store(&defaultScenarioMap)
						return &ap
					}(),
					benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
						ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
						ap.Store(&defaultBenchJobMap)
						return &ap
					}(),
					jobs: func() *atomic.Pointer[map[string]string] {
						ap := atomic.Pointer[map[string]string]{}
						ap.Store(&defaultJobMap)
						return &ap
					}(),
				},
				want:      want{},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T) {
					t.Helper()
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
				},
			}
		}(),
		func() test {
			return test{
				name: "return mismatch error when scenario and job has atomic and benchJob has no atomic",
				fields: fields{
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						ap.Store(&defaultScenarioMap)
						return &ap
					}(),
					jobs: func() *atomic.Pointer[map[string]string] {
						ap := atomic.Pointer[map[string]string]{}
						ap.Store(&defaultJobMap)
						return &ap
					}(),
				},
				want: want{
					err: errors.ErrMismatchBenchmarkAtomics(defaultJobMap, map[string]*v1.ValdBenchmarkJob{}, defaultScenarioMap),
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T) {
					t.Helper()
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
				},
			}
		}(),
		func() test {
			benchJobMap := map[string]*v1.ValdBenchmarkJob{}
			for k, v := range defaultBenchJobMap {
				val := *v
				benchJobMap[k] = &val
			}
			benchJobMap["scenario-search"].SetNamespace("benchmark")
			return test{
				name: "return mismatch error when benchJob with different namespace",
				fields: fields{
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						ap.Store(&defaultScenarioMap)
						return &ap
					}(),
					benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
						ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
						ap.Store(&benchJobMap)
						return &ap
					}(),
					jobs: func() *atomic.Pointer[map[string]string] {
						ap := atomic.Pointer[map[string]string]{}
						ap.Store(&defaultJobMap)
						return &ap
					}(),
				},
				want: want{
					err: errors.ErrMismatchBenchmarkAtomics(defaultJobMap, benchJobMap, defaultScenarioMap),
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T) {
					t.Helper()
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
				},
			}
		}(),
		func() test {
			benchJobMap := map[string]*v1.ValdBenchmarkJob{}
			for k, v := range defaultBenchJobMap {
				val := *v
				benchJobMap[k] = &val
			}
			benchJobMap["scenario-search"].Status = v1.BenchmarkJobNotReady
			return test{
				name: "return mismatch error when status is not same between benchJob and scenario.BenchJobStatus",
				fields: fields{
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						ap.Store(&defaultScenarioMap)
						return &ap
					}(),
					benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
						ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
						ap.Store(&benchJobMap)
						return &ap
					}(),
					jobs: func() *atomic.Pointer[map[string]string] {
						ap := atomic.Pointer[map[string]string]{}
						ap.Store(&defaultJobMap)
						return &ap
					}(),
				},
				want: want{
					err: errors.ErrMismatchBenchmarkAtomics(defaultJobMap, benchJobMap, defaultScenarioMap),
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T) {
					t.Helper()
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
				},
			}
		}(),
		func() test {
			benchJobMap := map[string]*v1.ValdBenchmarkJob{}
			for k, v := range defaultBenchJobMap {
				val := *v
				benchJobMap[k] = &val
			}
			ors := make([]k8s.OwnerReference, 0, len(benchJobMap["scenario-search"].OwnerReferences))
			for _, v := range benchJobMap["scenario-search"].OwnerReferences {
				or := v.DeepCopy()
				or.Name = "incorrectName"
				ors = append(ors, *or)
			}
			benchJobMap["scenario-search"].OwnerReferences = ors
			return test{
				name: "return mismatch error when scenario does not have key of bench job owners scenario",
				fields: fields{
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						ap.Store(&defaultScenarioMap)
						return &ap
					}(),
					benchjobs: func() *atomic.Pointer[map[string]*v1.ValdBenchmarkJob] {
						ap := atomic.Pointer[map[string]*v1.ValdBenchmarkJob]{}
						ap.Store(&benchJobMap)
						return &ap
					}(),
					jobs: func() *atomic.Pointer[map[string]string] {
						ap := atomic.Pointer[map[string]string]{}
						ap.Store(&defaultJobMap)
						return &ap
					}(),
				},
				want: want{
					err: errors.ErrMismatchBenchmarkAtomics(defaultJobMap, benchJobMap, defaultScenarioMap),
				},
				checkFunc: defaultCheckFunc,
				beforeFunc: func(t *testing.T) {
					t.Helper()
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
				},
			}
		}(),
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := new(operator)
			storeInto(&o.scenarios, test.fields.scenarios)
			storeInto(&o.benchjobs, test.fields.benchjobs)
			storeInto(&o.jobs, test.fields.jobs)

			err := o.checkAtomics()
			if err = checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

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

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Operator
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Operator, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Operator, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_initCtrl(t *testing.T) {
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			err := o.initCtrl()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_loadAtomic(t *testing.T) {
// 	type args struct {
// 		p *atomic.Pointer[M]
// 	}
// 	type want struct {
// 		want M
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, M) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got M) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           p:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           p:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := loadAtomic(test.args.p)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_k8sClient(t *testing.T) {
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		want k8s.Client
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, k8s.Client) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.Client) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			got := o.k8sClient()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_singleOwnerRef(t *testing.T) {
// 	type args struct {
// 		apiVersion string
// 		kind       string
// 		name       string
// 		uid        k8s.UID
// 	}
// 	type want struct {
// 		want []k8s.OwnerReference
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []k8s.OwnerReference) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []k8s.OwnerReference) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           apiVersion:"",
// 		           kind:"",
// 		           name:"",
// 		           uid:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           apiVersion:"",
// 		           kind:"",
// 		           name:"",
// 		           uid:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := singleOwnerRef(test.args.apiVersion, test.args.kind, test.args.name, test.args.uid)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_deleteBenchmarkJob(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		name       string
// 		generation int64
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           name:"",
// 		           generation:0,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           name:"",
// 		           generation:0,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			err := o.deleteBenchmarkJob(test.args.ctx, test.args.name, test.args.generation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_deleteJob(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		name string
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           name:"",
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           name:"",
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			err := o.deleteJob(test.args.ctx, test.args.name)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_createBenchmarkJob(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		scenario v1.ValdBenchmarkScenario
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		want []string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           scenario:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           scenario:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			got, err := o.createBenchmarkJob(test.args.ctx, test.args.scenario)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_createJob(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		bjr v1.ValdBenchmarkJob
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           bjr:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           bjr:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			err := o.createJob(test.args.ctx, test.args.bjr)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_updateCachedStatuses(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		c       k8s.Client
// 		desired map[string]S
// 		resolve func(name string, status S) (obj k8s.Object, ok bool)
// 	}
// 	type want struct {
// 		wantUpdated []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotUpdated []string) error {
// 		if !reflect.DeepEqual(gotUpdated, w.wantUpdated) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUpdated, w.wantUpdated)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           c:nil,
// 		           desired:nil,
// 		           resolve:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           c:nil,
// 		           desired:nil,
// 		           resolve:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotUpdated := updateCachedStatuses(test.args.ctx, test.args.c, test.args.desired, test.args.resolve)
// 			if err := checkFunc(test.want, gotUpdated); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_updateBenchmarkScenarioStatus(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		ss  map[string]v1.ValdBenchmarkScenarioStatus
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		want []string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           ss:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           ss:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			got, err := o.updateBenchmarkScenarioStatus(test.args.ctx, test.args.ss)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_updateBenchmarkJobStatus(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		js  map[string]v1.BenchmarkJobStatus
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		want []string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           js:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           js:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			got, err := o.updateBenchmarkJobStatus(test.args.ctx, test.args.js)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_checkJobsStatus(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		jobs map[string]string
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           jobs:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           jobs:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			err := o.checkJobsStatus(test.args.ctx, test.args.jobs)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_GetScenarioStatus(t *testing.T) {
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		want map[v1.ValdBenchmarkScenarioStatus]int64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, map[v1.ValdBenchmarkScenarioStatus]int64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got map[v1.ValdBenchmarkScenarioStatus]int64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			got := o.GetScenarioStatus()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_GetBenchmarkJobStatus(t *testing.T) {
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		want map[v1.BenchmarkJobStatus]int64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, map[v1.BenchmarkJobStatus]int64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got map[v1.BenchmarkJobStatus]int64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			got := o.GetBenchmarkJobStatus()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_PreStart(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			err := o.PreStart(test.args.in0)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		jobNamespace       string
// 		jobImageRepository string
// 		jobImageTag        string
// 		jobImagePullPolicy string
// 		configMapName      string
// 		rcd                time.Duration
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           ctrl:nil,
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           jobNamespace:"",
// 		           jobImageRepository:"",
// 		           jobImageTag:"",
// 		           jobImagePullPolicy:"",
// 		           configMapName:"",
// 		           rcd:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			o := &operator{
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImageRepository: test.fields.jobImageRepository,
// 				jobImageTag:        test.fields.jobImageTag,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				configMapName:      test.fields.configMapName,
// 				rcd:                test.fields.rcd,
// 			}
//
// 			got, err := o.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// TestOperator_Start_ErrorRelayStopsOnCancelWhenConsumerGone pins that the relay
// goroutine does not block forever on its error channel once the consumer has
// stopped reading: after the buffer is full, cancelling ctx must end it.
func TestOperator_Start_ErrorRelayStopsOnCancelWhenConsumerGone(t *testing.T) {
	t.Parallel()

	const sendTimeout = 5 * time.Second

	scch := make(chan error)
	eg, egctx := errgroup.New(t.Context())
	ctx, cancel := context.WithCancel(egctx)
	defer cancel()

	o := &operator{
		eg:  eg,
		rcd: time.Hour,
		ctrl: &mockCtrl{
			StartFunc: func(context.Context) (<-chan error, error) { return scch, nil },
		},
	}
	// the returned channel is intentionally never read.
	if _, err := o.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	// one more error than the channel buffer: the last one cannot be buffered.
	for range startErrChanBufferSize + 1 {
		select {
		case scch <- errors.New("controller error"):
		case <-time.After(sendTimeout):
			t.Fatal("relay goroutine did not receive the controller error")
		}
	}

	cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			t.Logf("eg.Wait() = %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(sendTimeout):
		t.Fatal("relay goroutine stayed blocked on the error channel after ctx cancel")
	}
}

// doneCountCtx counts Done() calls. A select loop evaluates ctx.Done() on every
// iteration, so the count is a deterministic measure of how often it spins.
type doneCountCtx struct {
	context.Context //nolint:containedctx // test wrapper that counts Done() calls
	n               atomic.Int64
}

func (c *doneCountCtx) Done() <-chan struct{} {
	c.n.Add(1)
	return c.Context.Done()
}

// TestOperator_Start_ClosedControllerChannelDoesNotSpin pins that the relay
// loop stops selecting the controller error channel once it is closed (the
// controller's manager returned while ctx is still live). A closed channel is
// always ready, so without nil-ing it the loop busy-spins; the spin shows up
// as an unbounded number of ctx.Done() evaluations.
func TestOperator_Start_ClosedControllerChannelDoesNotSpin(t *testing.T) {
	t.Parallel()

	ch := make(chan error)
	close(ch)

	eg, egctx := errgroup.New(t.Context())
	cctx, cancel := context.WithCancel(egctx)
	defer cancel()
	ctx := &doneCountCtx{Context: cctx}

	o := &operator{
		eg:  eg,
		rcd: time.Hour,
		ctrl: &mockCtrl{
			StartFunc: func(context.Context) (<-chan error, error) { return ch, nil },
		},
	}
	if _, err := o.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	time.Sleep(100 * time.Millisecond)
	if n := ctx.n.Load(); n > 100 {
		t.Errorf("ctx.Done() evaluated %d times with a closed controller channel; the relay loop is busy-spinning", n)
	}

	cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			t.Logf("eg.Wait() = %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("relay goroutine did not exit on context cancel")
	}
}

// newClientMockCtrl returns a controller whose manager routes Create, Get and
// Status().Update to the given hooks (nil means success).
func newClientMockCtrl(
	create func(context.Context, k8s.Object) error,
	get func(context.Context, k8s.ObjectKey, k8s.Object) error,
	statusUpdate func(context.Context, k8s.Object) error,
) *mockCtrl {
	return &mockCtrl{
		StartFunc: func(context.Context) (<-chan error, error) { return nil, nil },
		GetManagerFunc: func() k8s.Manager {
			return k8smock.NewHookManagerMock(create, get, statusUpdate)
		},
	}
}

// TestOperator_benchJobReconcile_CreateJobFailureIsRetried pins that a job
// whose creation failed is not cached as created (neither as a new job nor as
// the new generation of an already cached job) and is not marked healthy, so
// that the next reconcile retries the creation.
func TestOperator_benchJobReconcile_CreateJobFailureIsRetried(t *testing.T) {
	t.Parallel()

	const name = "bench-job"
	tests := []struct {
		// cached is the job already held in cbjl before the first reconcile (nil: none).
		cached *v1.ValdBenchmarkJob
		name   string
		// wantStatusAfterRetry is the cached status expected after the retry succeeded.
		wantStatusAfterRetry v1.BenchmarkJobStatus
		// listed is the job reported by the reconciler.
		listed v1.ValdBenchmarkJob
		// wantCachedGenAfterFailure is the generation expected in cbjl after the
		// failed creation; -1 means the job must be absent.
		wantCachedGenAfterFailure int64
	}{
		{
			name:                      "new job",
			listed:                    v1.ValdBenchmarkJob{ObjectMeta: k8s.ObjectMeta{Name: name, Namespace: "default"}},
			wantCachedGenAfterFailure: -1,
			wantStatusAfterRetry:      v1.BenchmarkJobHealthy,
		},
		{
			name: "new generation of a cached job",
			cached: &v1.ValdBenchmarkJob{
				Name: name, Namespace: "default", Generation: 1,
				Status: v1.BenchmarkJobHealthy,
			},
			listed: v1.ValdBenchmarkJob{
				ObjectMeta: k8s.ObjectMeta{Name: name, Namespace: "default", Generation: 2},
				Status:     v1.BenchmarkJobAvailable,
			},
			wantCachedGenAfterFailure: 1,
			wantStatusAfterRetry:      v1.BenchmarkJobAvailable,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			createErr := errors.New("create failed")
			creates := 0
			var statusUpdated []string
			o := &operator{
				jobNamespace:       "default",
				jobImageRepository: "vdaas/vald-benchmark-job",
				jobImageTag:        "latest",
				jobImagePullPolicy: "Always",
				configMapName:      "vald-benchmark-operator-config",
				ctrl: newClientMockCtrl(
					func(context.Context, k8s.Object) error {
						creates++
						if creates == 1 {
							return createErr
						}
						return nil
					},
					nil,
					func(_ context.Context, obj k8s.Object) error {
						statusUpdated = append(statusUpdated, obj.GetName())
						return nil
					},
				),
			}
			if tc.cached != nil {
				cached := map[string]*v1.ValdBenchmarkJob{name: tc.cached}
				o.benchjobs.Store(&cached)
			}
			list := map[string]v1.ValdBenchmarkJob{name: tc.listed}

			// first reconcile: createJob fails.
			o.benchJobReconcile(t.Context(), list)
			if creates != 1 {
				t.Fatalf("Create called %d times, want 1", creates)
			}
			got, ok := o.getAtomicBenchJob()[name]
			switch {
			case tc.wantCachedGenAfterFailure < 0 && ok:
				t.Errorf("job %q was cached although its creation failed", name)
			case tc.wantCachedGenAfterFailure >= 0 && !ok:
				t.Fatalf("cached job %q disappeared after a failed creation", name)
			case ok && got.GetGeneration() != tc.wantCachedGenAfterFailure:
				t.Errorf("cached generation = %d after a failed creation, want %d", got.GetGeneration(), tc.wantCachedGenAfterFailure)
			}
			if len(statusUpdated) != 0 {
				t.Errorf("status updated for %v although the creation failed", statusUpdated)
			}

			// second reconcile: the job is retried and, on success, cached.
			o.benchJobReconcile(t.Context(), list)
			if creates != 2 {
				t.Fatalf("Create called %d times after the second reconcile, want 2", creates)
			}
			got, ok = o.getAtomicBenchJob()[name]
			if !ok {
				t.Fatalf("job %q is not cached after a successful retry", name)
			}
			if got.GetGeneration() != tc.listed.GetGeneration() {
				t.Errorf("cached generation = %d after the retry, want %d", got.GetGeneration(), tc.listed.GetGeneration())
			}
			if got.Status != tc.wantStatusAfterRetry {
				t.Errorf("cached job status = %q, want %q", got.Status, tc.wantStatusAfterRetry)
			}
		})
	}
}

// TestOperator_benchJobReconcile_CreateJobAlreadyExistsIsCreated pins that an
// AlreadyExists error from Create (the job was created before, e.g. before an
// operator restart) counts as created: the job is cached on the first
// reconcile instead of failing on every reconcile and never being tracked.
// createJob wraps the Create error, so this also covers that the wrapped error
// is still recognized.
func TestOperator_benchJobReconcile_CreateJobAlreadyExistsIsCreated(t *testing.T) {
	t.Parallel()

	const name = "bench-job"
	tests := []struct {
		// cached is the job already held in cbjl before the reconcile (nil: none).
		cached *v1.ValdBenchmarkJob
		name   string
		// wantStatus is the cached status expected after the reconcile.
		wantStatus v1.BenchmarkJobStatus
		// listed is the job reported by the reconciler.
		listed v1.ValdBenchmarkJob
		// wantStatusUpdated is whether the CR status is expected to be updated.
		wantStatusUpdated bool
		// wantGen is the cached generation expected after the reconcile.
		wantGen int64
	}{
		{
			name:              "new job",
			listed:            v1.ValdBenchmarkJob{ObjectMeta: k8s.ObjectMeta{Name: name, Namespace: "default"}},
			wantGen:           0,
			wantStatus:        v1.BenchmarkJobHealthy,
			wantStatusUpdated: true,
		},
		{
			// the job name is shared across generations, so AlreadyExists here means the
			// old version job is not gone yet: keep the old generation cached so the next
			// reconcile retries delete+create instead of never applying the new spec.
			name: "new generation of a cached job keeps the old generation",
			cached: &v1.ValdBenchmarkJob{
				Name: name, Namespace: "default", Generation: 1,
				Status: v1.BenchmarkJobHealthy,
			},
			listed: v1.ValdBenchmarkJob{
				ObjectMeta: k8s.ObjectMeta{Name: name, Namespace: "default", Generation: 2},
				Status:     v1.BenchmarkJobAvailable,
			},
			wantGen:    1,
			wantStatus: v1.BenchmarkJobHealthy,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var statusUpdated []string
			o := &operator{
				jobNamespace:       "default",
				jobImageRepository: "vdaas/vald-benchmark-job",
				jobImageTag:        "latest",
				jobImagePullPolicy: "Always",
				configMapName:      "vald-benchmark-operator-config",
				ctrl: newClientMockCtrl(
					func(context.Context, k8s.Object) error {
						return k8smock.NewAlreadyExistsError(name)
					},
					nil,
					func(_ context.Context, obj k8s.Object) error {
						statusUpdated = append(statusUpdated, obj.GetName())
						return nil
					},
				),
			}
			if tc.cached != nil {
				cached := map[string]*v1.ValdBenchmarkJob{name: tc.cached}
				o.benchjobs.Store(&cached)
			}

			o.benchJobReconcile(t.Context(), map[string]v1.ValdBenchmarkJob{name: tc.listed})

			got, ok := o.getAtomicBenchJob()[name]
			if !ok {
				t.Fatalf("job %q is not cached although it already exists", name)
			}
			if got.GetGeneration() != tc.wantGen {
				t.Errorf("cached generation = %d, want %d", got.GetGeneration(), tc.wantGen)
			}
			if got.Status != tc.wantStatus {
				t.Errorf("cached status = %q, want %q", got.Status, tc.wantStatus)
			}
			if tc.wantStatusUpdated != (len(statusUpdated) != 0) {
				t.Errorf("status updated for %v, want updated = %t", statusUpdated, tc.wantStatusUpdated)
			}
		})
	}
}

// TestOperator_checkJobsStatus_GetFailureDoesNotSkipOthers pins that a failing
// Get for one job neither aborts the status update of the other jobs nor is
// swallowed: the remaining job is updated and the Get error is returned.
func TestOperator_checkJobsStatus_GetFailureDoesNotSkipOthers(t *testing.T) {
	t.Parallel()

	const (
		failing = "job-get-fails"
		healthy = "job-succeeded"
	)
	getErr := errors.New("get failed")
	var updated []string
	o := &operator{
		ctrl: newClientMockCtrl(
			nil,
			func(_ context.Context, key k8s.ObjectKey, obj k8s.Object) error {
				if key.Name == failing {
					return getErr
				}
				job, ok := obj.(*k8s.Job)
				if !ok {
					return errors.New("unexpected object type")
				}
				job.Status.Succeeded = 1
				return nil
			},
			func(_ context.Context, obj k8s.Object) error {
				updated = append(updated, obj.GetName())
				return nil
			},
		),
	}
	cached := map[string]*v1.ValdBenchmarkJob{}
	for _, n := range []string{failing, healthy} {
		cached[n] = &v1.ValdBenchmarkJob{
			Name: n, Namespace: "default",
			Status: v1.BenchmarkJobHealthy,
		}
	}
	o.benchjobs.Store(&cached)

	err := o.checkJobsStatus(t.Context(), map[string]string{
		failing: "default",
		healthy: "default",
	})
	if !errors.Is(err, getErr) {
		t.Errorf("checkJobsStatus() error = %v, want it to wrap %v", err, getErr)
	}
	if len(updated) != 1 || updated[0] != healthy {
		t.Errorf("status updated for %v, want only [%s]", updated, healthy)
	}
	got := o.getAtomicBenchJob()
	if st := got[healthy].Status; st != v1.BenchmarkJobCompleted {
		t.Errorf("cached status of %s = %q, want %q", healthy, st, v1.BenchmarkJobCompleted)
	}
	if st := got[failing].Status; st != v1.BenchmarkJobHealthy {
		t.Errorf("cached status of %s = %q, want unchanged %q", failing, st, v1.BenchmarkJobHealthy)
	}
}

// TestOperator_refreshStatuses_DoesNotLoseConcurrentReconcile pins that the
// periodic status refresh and a benchmark job reconcile do not overwrite each
// other's cache updates. The refresh is started while the reconcile is inside
// its createJob call; without serialization the reconcile then stores a copy
// taken before the refresh and reverts the refreshed status.
func TestOperator_refreshStatuses_DoesNotLoseConcurrentReconcile(t *testing.T) {
	t.Parallel()

	const (
		existing = "job-existing"
		added    = "job-added"
	)
	inCreate := make(chan struct{})
	release := make(chan struct{})
	o := &operator{
		jobImageRepository: "vdaas/vald-benchmark-job",
		jobImagePullPolicy: "Always",
		configMapName:      "refresh-test-configmap",
		ctrl: newClientMockCtrl(
			func(ctx context.Context, obj k8s.Object) error {
				if obj.GetName() != added {
					return nil
				}
				close(inCreate)
				select {
				case <-release:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			},
			func(_ context.Context, _ k8s.ObjectKey, obj k8s.Object) error {
				if job, ok := obj.(*k8s.Job); ok {
					job.Status.Succeeded = 1
				}
				return nil
			},
			nil,
		),
	}
	existingJob := v1.ValdBenchmarkJob{
		Name: existing, Namespace: "default",
		Status: v1.BenchmarkJobHealthy,
	}
	cached := map[string]*v1.ValdBenchmarkJob{existing: &existingJob}
	o.benchjobs.Store(&cached)
	jobs := map[string]string{existing: "default"}
	o.jobs.Store(&jobs)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	reconcileDone := make(chan struct{})
	go func() {
		defer close(reconcileDone)
		o.benchJobReconcile(ctx, map[string]v1.ValdBenchmarkJob{
			existing: existingJob,
			added:    {ObjectMeta: k8s.ObjectMeta{Name: added, Namespace: "default"}},
		})
	}()
	select {
	case <-inCreate:
	case <-ctx.Done():
		t.Fatal("benchJobReconcile did not reach createJob")
	}
	refreshDone := make(chan struct{})
	go func() {
		defer close(refreshDone)
		o.refreshStatuses(ctx)
	}()
	// Give an unserialized refresh the chance to finish inside the createJob
	// window; a serialized one stays blocked until the reconcile returns.
	select {
	case <-refreshDone:
	case <-time.After(200 * time.Millisecond):
	}
	close(release)
	for _, done := range []chan struct{}{reconcileDone, refreshDone} {
		select {
		case <-done:
		case <-ctx.Done():
			t.Fatal("reconcile or refresh did not finish")
		}
	}

	got := o.getAtomicBenchJob()
	if _, ok := got[added]; !ok {
		t.Errorf("job %s added by the reconcile was lost", added)
	}
	if bj, ok := got[existing]; !ok || bj.Status != v1.BenchmarkJobCompleted {
		t.Errorf("job %s = %v, want status %q from the refresh", existing, bj, v1.BenchmarkJobCompleted)
	}
}

// TestOperator_Start_CheckAtomicsWaitsForReconcile pins that the periodic
// atomics check runs under mu. A reconcile publishes the job, benchmark job
// and scenario atomics one at a time while holding mu, so an unlocked check
// can observe its half-published state and report a spurious mismatch.
func TestOperator_Start_CheckAtomicsWaitsForReconcile(t *testing.T) {
	t.Parallel()

	const (
		rcd     = 5 * time.Millisecond
		observe = 20 * rcd
		waitFor = 5 * time.Second
	)
	eg, egctx := errgroup.New(t.Context())
	ctx, cancel := context.WithCancel(egctx)
	defer cancel()
	o := &operator{
		eg:   eg,
		rcd:  rcd,
		ctrl: newClientMockCtrl(nil, nil, nil),
	}

	// a reconcile in progress: it holds mu and has published only part of its
	// update (a job and a scenario without the matching benchmark job).
	o.mu.Lock()
	jobs := map[string]string{"job": "default"}
	o.jobs.Store(&jobs)
	scenarios := map[string]*scenario{"scenario": {
		Crd: &v1.ValdBenchmarkScenario{Name: "scenario", Namespace: "default"},
	}}
	o.scenarios.Store(&scenarios)
	ech, err := o.Start(ctx)
	if err != nil {
		o.mu.Unlock()
		t.Fatalf("Start() error = %v, want nil", err)
	}
	select {
	case err := <-ech:
		if err != nil {
			t.Errorf("atomics check reported %v for a state published while a reconcile held mu", err)
		}
	case <-time.After(observe):
	}
	// the reconcile finishes with a consistent state and releases mu.
	emptyJobs := map[string]string{}
	o.jobs.Store(&emptyJobs)
	emptyScenarios := map[string]*scenario{}
	o.scenarios.Store(&emptyScenarios)
	o.mu.Unlock()
	select {
	case err := <-ech:
		if err != nil {
			t.Errorf("atomics check reported %v for a consistent state", err)
		}
	case <-time.After(observe):
	}

	cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			t.Logf("eg.Wait() = %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(waitFor):
		t.Fatal("Start goroutine did not exit on context cancel")
	}
	for err := range ech {
		if err != nil {
			t.Errorf("atomics check reported %v", err)
		}
	}
}

// TestOperator_updateStatus_FailedUpdateKeepsCache pins that a status whose
// Status().Update failed is not committed to the cache: the cached object
// keeps its old status so the next refresh retries the transition, while a
// successfully updated name in the same batch is committed.
func TestOperator_updateStatus_FailedUpdateKeepsCache(t *testing.T) {
	t.Parallel()

	const (
		failing = "fails"
		passing = "passes"
	)
	updateErr := errors.New("update failed")
	statusUpdate := func(_ context.Context, obj k8s.Object) error {
		if obj.GetName() == failing {
			return updateErr
		}
		return nil
	}
	tests := []struct {
		run  func(t *testing.T, o *operator) (got map[string]string, err error)
		want map[string]string
		name string
	}{
		{
			name: "benchmark job",
			want: map[string]string{
				failing: string(v1.BenchmarkJobHealthy),
				passing: string(v1.BenchmarkJobCompleted),
			},
			run: func(t *testing.T, o *operator) (map[string]string, error) {
				t.Helper()
				cached := map[string]*v1.ValdBenchmarkJob{}
				for _, n := range []string{failing, passing} {
					cached[n] = &v1.ValdBenchmarkJob{
						Name: n, Namespace: "default",
						Status: v1.BenchmarkJobHealthy,
					}
				}
				o.benchjobs.Store(&cached)
				_, err := o.updateBenchmarkJobStatus(t.Context(), map[string]v1.BenchmarkJobStatus{
					failing: v1.BenchmarkJobCompleted,
					passing: v1.BenchmarkJobCompleted,
				})
				got := map[string]string{}
				for n, bj := range o.getAtomicBenchJob() {
					got[n] = string(bj.Status)
				}
				return got, err
			},
		},
		{
			name: "benchmark scenario",
			want: map[string]string{
				failing: string(v1.BenchmarkScenarioHealthy),
				passing: string(v1.BenchmarkScenarioCompleted),
			},
			run: func(t *testing.T, o *operator) (map[string]string, error) {
				t.Helper()
				cached := map[string]*scenario{}
				for _, n := range []string{failing, passing} {
					cached[n] = &scenario{Crd: &v1.ValdBenchmarkScenario{
						Name: n, Namespace: "default",
						Status: v1.BenchmarkScenarioHealthy,
					}}
				}
				o.scenarios.Store(&cached)
				_, err := o.updateBenchmarkScenarioStatus(t.Context(), map[string]v1.ValdBenchmarkScenarioStatus{
					failing: v1.BenchmarkScenarioCompleted,
					passing: v1.BenchmarkScenarioCompleted,
				})
				got := map[string]string{}
				for n, sc := range o.getAtomicScenario() {
					got[n] = string(sc.Crd.Status)
				}
				return got, err
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			o := &operator{ctrl: newClientMockCtrl(nil, nil, statusUpdate)}
			got, err := tc.run(t, o)
			if err != nil {
				t.Fatalf("update status error = %v, want nil", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("cached statuses = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestOperator_benchJobReconcile_PrunesStaleScenarioJobStatus pins that a
// scenario touched by the reconcile drops the job statuses of benchmark jobs
// no longer listed (so they cannot keep it from completing), but keeps the
// NotReady entries of jobs just created that the list may not contain yet.
func TestOperator_benchJobReconcile_PrunesStaleScenarioJobStatus(t *testing.T) {
	t.Parallel()

	const (
		owner   = "scenario"
		listed  = "scenario-listed"
		fresh   = "scenario-fresh"
		healthy = "scenario-stale-healthy"
		done    = "scenario-stale-completed"
	)
	o := &operator{
		jobImageRepository: "vdaas/vald-benchmark-job",
		jobImagePullPolicy: "Always",
		configMapName:      "prune-test-configmap",
		ctrl:               newClientMockCtrl(nil, nil, nil),
	}
	scenarios := map[string]*scenario{owner: {
		Crd: &v1.ValdBenchmarkScenario{Name: owner, Namespace: "default"},
		BenchJobStatus: map[string]v1.BenchmarkJobStatus{
			listed:  v1.BenchmarkJobNotReady,
			fresh:   v1.BenchmarkJobNotReady,
			healthy: v1.BenchmarkJobHealthy,
			done:    v1.BenchmarkJobCompleted,
		},
	}}
	o.scenarios.Store(&scenarios)

	o.benchJobReconcile(t.Context(), map[string]v1.ValdBenchmarkJob{
		listed: {
			ObjectMeta: k8s.ObjectMeta{
				Name:            listed,
				Namespace:       "default",
				OwnerReferences: []k8s.OwnerReference{{Name: owner}},
			},
		},
	})

	want := map[string]v1.BenchmarkJobStatus{
		listed: "",
		fresh:  v1.BenchmarkJobNotReady,
	}
	if got := o.getAtomicScenario()[owner].BenchJobStatus; !reflect.DeepEqual(got, want) {
		t.Errorf("scenario BenchJobStatus = %v, want %v", got, want)
	}
	// the published snapshot taken before the reconcile must stay untouched.
	if got := len(scenarios[owner].BenchJobStatus); got != 4 {
		t.Errorf("previous snapshot BenchJobStatus has %d entries, want 4", got)
	}
}

// TestOperator_benchScenarioReconcile_JobCreationFailureIsResumed pins that a
// scenario whose benchmark job creation failed is marked NotReady instead of
// Healthy, is not completed by the status refresh, and that the next
// reconcile — or, without any scenario event, the next status refresh —
// creates only the jobs still missing (no duplicates, chained after
// the ones already created) and then marks the scenario Healthy.
func TestOperator_benchScenarioReconcile_JobCreationFailureIsResumed(t *testing.T) {
	t.Parallel()

	const name = "scenario"
	jobTypes := []string{"insert", "search", "upsert"}
	newScenario := func(gen int64, status v1.ValdBenchmarkScenarioStatus) v1.ValdBenchmarkScenario {
		jobs := make([]*v1.BenchmarkJobSpec, 0, len(jobTypes))
		for _, jt := range jobTypes {
			jobs = append(jobs, &v1.BenchmarkJobSpec{JobType: jt})
		}
		return v1.ValdBenchmarkScenario{
			Name: name, Namespace: "default", Generation: gen,
			Spec:   v1.ValdBenchmarkScenarioSpec{Jobs: jobs},
			Status: status,
		}
	}
	tests := []struct {
		// cached is the scenario held before the first reconcile (nil: none).
		cached *scenario
		name   string
		listed v1.ValdBenchmarkScenario
		// failOn is the Create call that fails.
		failOn int
	}{
		{
			name:   "new scenario, first job fails",
			listed: newScenario(1, v1.BenchmarkScenarioAvailable),
			failOn: 1,
		},
		{
			name:   "new scenario, second job fails",
			listed: newScenario(1, v1.BenchmarkScenarioAvailable),
			failOn: 2,
		},
		{
			name: "new generation, second job fails",
			cached: func() *scenario {
				old := newScenario(1, v1.BenchmarkScenarioHealthy)
				return &scenario{
					Crd:            &old,
					BenchJobStatus: map[string]v1.BenchmarkJobStatus{"old-job": v1.BenchmarkJobCompleted},
				}
			}(),
			listed: newScenario(2, v1.BenchmarkScenarioHealthy),
			failOn: 2,
		},
	}
	for _, tc := range tests {
		for _, viaRefresh := range []bool{false, true} {
			resumer := "reconcile"
			if viaRefresh {
				resumer = "status refresh"
			}
			t.Run(tc.name+", resumed by "+resumer, func(t *testing.T) {
				t.Parallel()

				type createdJob struct {
					name, jobType, before string
				}
				var (
					creates  int
					created  []createdJob
					statuses []v1.ValdBenchmarkScenarioStatus
				)
				o := &operator{
					jobNamespace: "default",
					ctrl: newClientMockCtrl(
						func(_ context.Context, obj k8s.Object) error {
							creates++
							if creates == tc.failOn {
								return errors.New("create failed")
							}
							bj, ok := obj.(*v1.ValdBenchmarkJob)
							if !ok {
								return errors.Errorf("unexpected object %T", obj)
							}
							created = append(created, createdJob{
								name:    bj.GetName(),
								jobType: bj.Spec.JobType,
								before:  bj.GetAnnotations()[BeforeJobName],
							})
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
				if tc.cached != nil {
					cached := map[string]*scenario{name: tc.cached}
					o.scenarios.Store(&cached)
				}
				list := map[string]v1.ValdBenchmarkScenario{name: tc.listed}

				// first reconcile: creation fails at call failOn.
				o.benchScenarioReconcile(t.Context(), list)
				if want := []v1.ValdBenchmarkScenarioStatus{v1.BenchmarkScenarioNotReady}; !reflect.DeepEqual(statuses, want) {
					t.Fatalf("scenario status writes after the failed creation = %v, want %v", statuses, want)
				}
				if got, want := len(created), tc.failOn-1; got != want {
					t.Fatalf("%d jobs created before the failure, want %d", got, want)
				}
				bjs := o.getAtomicScenario()[name].BenchJobStatus
				if len(bjs) != len(created) {
					t.Errorf("scenario tracks %v, want the %d jobs created before the failure", bjs, len(created))
				}
				// the jobs created before the failure run and complete.
				done := map[string]v1.ValdBenchmarkJob{}
				for _, c := range created {
					done[c.name] = v1.ValdBenchmarkJob{
						Name:            c.name,
						Namespace:       "default",
						OwnerReferences: []k8s.OwnerReference{{Name: name, Kind: ScenarioKind}},
						Status:          v1.BenchmarkJobCompleted,
					}
				}
				o.benchJobReconcile(t.Context(), done)
				if viaRefresh {
					// the refresh resumes the creation itself; it must not
					// complete the scenario.
					o.refreshStatuses(t.Context())
					if want := []v1.ValdBenchmarkScenarioStatus{v1.BenchmarkScenarioNotReady, v1.BenchmarkScenarioHealthy}; !reflect.DeepEqual(statuses, want) {
						t.Fatalf("scenario status writes after the status refresh = %v, want %v", statuses, want)
					}
				} else {
					// second reconcile: only the missing jobs are created.
					o.benchScenarioReconcile(t.Context(), list)
				}
				if len(created) != len(jobTypes) {
					t.Fatalf("created jobs = %+v, want one per spec job %v", created, jobTypes)
				}
				for i, c := range created {
					if c.jobType != jobTypes[i] {
						t.Errorf("created job %d has type %q, want %q", i, c.jobType, jobTypes[i])
					}
					wantBefore := ""
					if i > 0 {
						wantBefore = created[i-1].name
					}
					if c.before != wantBefore {
						t.Errorf("created job %d waits for %q, want %q", i, c.before, wantBefore)
					}
				}
				if got := statuses[len(statuses)-1]; got != v1.BenchmarkScenarioHealthy {
					t.Errorf("scenario status after the resumed creation = %q, want %q", got, v1.BenchmarkScenarioHealthy)
				}
				bjs = o.getAtomicScenario()[name].BenchJobStatus
				if len(bjs) != len(jobTypes) {
					t.Errorf("scenario tracks %v, want %d jobs", bjs, len(jobTypes))
				}
				for _, c := range created[:tc.failOn-1] {
					if bjs[c.name] != v1.BenchmarkJobCompleted {
						t.Errorf("job %s created before the failure has status %q, want %q kept", c.name, bjs[c.name], v1.BenchmarkJobCompleted)
					}
				}

				// a further reconcile creates nothing.
				o.benchScenarioReconcile(t.Context(), list)
				if len(created) != len(jobTypes) {
					t.Errorf("reconcile after a complete creation created %+v", created[len(jobTypes):])
				}
			})
		}
	}
}

// TestOperator_refreshStatuses_CompletesOnlyFullyTrackedScenario pins that a
// scenario is completed only when a job of every spec job is tracked and all
// tracked jobs are completed: a scenario with untracked or not-yet-created
// jobs must not be reported Completed.
func TestOperator_refreshStatuses_CompletesOnlyFullyTrackedScenario(t *testing.T) {
	t.Parallel()

	twoJobs := []*v1.BenchmarkJobSpec{{JobType: "insert"}, {JobType: "search"}}
	tests := []struct {
		name          string
		sc            scenario
		wantCompleted bool
	}{
		{
			name: "no job tracked",
			sc:   scenario{Crd: &v1.ValdBenchmarkScenario{Spec: v1.ValdBenchmarkScenarioSpec{Jobs: twoJobs}}},
		},
		{
			// all spec jobs are created (jobsIncomplete unset): the other
			// job's status was pruned because its benchmark job was deleted.
			name: "one of two spec jobs tracked, the other deleted",
			sc: scenario{
				Crd:            &v1.ValdBenchmarkScenario{Spec: v1.ValdBenchmarkScenarioSpec{Jobs: twoJobs}},
				BenchJobStatus: map[string]v1.BenchmarkJobStatus{"a": v1.BenchmarkJobCompleted},
			},
			wantCompleted: true,
		},
		{
			name: "every job deleted",
			sc: scenario{
				Crd:        &v1.ValdBenchmarkScenario{Spec: v1.ValdBenchmarkScenarioSpec{Jobs: twoJobs}},
				jobsPruned: true,
			},
			wantCompleted: true,
		},
		{
			name: "job creation incomplete",
			sc: scenario{
				Crd: &v1.ValdBenchmarkScenario{Spec: v1.ValdBenchmarkScenarioSpec{Jobs: twoJobs}},
				BenchJobStatus: map[string]v1.BenchmarkJobStatus{
					"a": v1.BenchmarkJobCompleted,
					"b": v1.BenchmarkJobCompleted,
				},
				createdJobs:    []string{"a"},
				jobsIncomplete: true,
			},
		},
		{
			name: "tracked job not completed",
			sc: scenario{
				Crd: &v1.ValdBenchmarkScenario{Spec: v1.ValdBenchmarkScenarioSpec{Jobs: twoJobs}},
				BenchJobStatus: map[string]v1.BenchmarkJobStatus{
					"a": v1.BenchmarkJobCompleted,
					"b": v1.BenchmarkJobHealthy,
				},
			},
		},
		{
			name: "all spec jobs completed",
			sc: scenario{
				Crd: &v1.ValdBenchmarkScenario{Spec: v1.ValdBenchmarkScenarioSpec{Jobs: twoJobs}},
				BenchJobStatus: map[string]v1.BenchmarkJobStatus{
					"a": v1.BenchmarkJobCompleted,
					"b": v1.BenchmarkJobCompleted,
				},
			},
			wantCompleted: true,
		},
		{
			name:          "spec without jobs",
			sc:            scenario{Crd: &v1.ValdBenchmarkScenario{}},
			wantCompleted: true,
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
			sc := tc.sc
			crd := *sc.Crd
			crd.Name = "scenario"
			crd.Status = v1.BenchmarkScenarioHealthy
			sc.Crd = &crd
			scenarios := map[string]*scenario{crd.Name: &sc}
			o.scenarios.Store(&scenarios)

			o.refreshStatuses(t.Context())

			var want []v1.ValdBenchmarkScenarioStatus
			if tc.wantCompleted {
				want = []v1.ValdBenchmarkScenarioStatus{v1.BenchmarkScenarioCompleted}
			}
			if !reflect.DeepEqual(statuses, want) {
				t.Errorf("scenario status writes = %v, want %v", statuses, want)
			}
		})
	}
}

// TestOperator_benchJobReconcile_PrunesDeletedJobsOfEveryScenario pins that
// the job statuses of deleted benchmark jobs are dropped from every scenario,
// including when no benchmark job is listed at all and for a scenario none of
// whose jobs is listed, while NotReady entries of just-created jobs are kept.
func TestOperator_benchJobReconcile_PrunesDeletedJobsOfEveryScenario(t *testing.T) {
	t.Parallel()

	const (
		owner   = "scenario"
		other   = "other"
		fresh   = "scenario-fresh"
		healthy = "scenario-deleted-healthy"
		done    = "scenario-deleted-completed"
	)
	tests := []struct {
		list map[string]v1.ValdBenchmarkJob
		name string
	}{
		{
			name: "no benchmark job listed",
			list: map[string]v1.ValdBenchmarkJob{},
		},
		{
			name: "only another scenario's job listed",
			list: map[string]v1.ValdBenchmarkJob{
				"other-job": {
					ObjectMeta: k8s.ObjectMeta{
						Name:            "other-job",
						Namespace:       "default",
						OwnerReferences: []k8s.OwnerReference{{Name: other}},
					},
					Status: v1.BenchmarkJobCompleted,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			o := &operator{ctrl: newClientMockCtrl(nil, nil, nil)}
			scenarios := map[string]*scenario{
				owner: {
					Crd: &v1.ValdBenchmarkScenario{Name: owner, Namespace: "default"},
					BenchJobStatus: map[string]v1.BenchmarkJobStatus{
						fresh:   v1.BenchmarkJobNotReady,
						healthy: v1.BenchmarkJobHealthy,
						done:    v1.BenchmarkJobCompleted,
					},
				},
				other: {
					Crd: &v1.ValdBenchmarkScenario{Name: other, Namespace: "default"},
				},
			}
			o.scenarios.Store(&scenarios)

			o.benchJobReconcile(t.Context(), tc.list)

			want := map[string]v1.BenchmarkJobStatus{fresh: v1.BenchmarkJobNotReady}
			if got := o.getAtomicScenario()[owner].BenchJobStatus; !reflect.DeepEqual(got, want) {
				t.Errorf("scenario BenchJobStatus = %v, want %v", got, want)
			}
			// the published snapshot taken before the reconcile must stay untouched.
			if got := len(scenarios[owner].BenchJobStatus); got != 3 {
				t.Errorf("previous snapshot BenchJobStatus has %d entries, want 3", got)
			}
		})
	}
}
