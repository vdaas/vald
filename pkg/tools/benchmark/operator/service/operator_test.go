// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/job"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/internal/test/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func Test_operator_getAtomicScenario(t *testing.T) {
	t.Parallel()
	type fields struct {
		scenarios *atomic.Pointer[map[string]*scenario]
	}
	type want struct {
		want map[string]*scenario
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, map[string]*scenario) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
								"scneario-insert": v1.BenchmarkJobAvailable,
								"scneario-search": v1.BenchmarkJobAvailable,
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
							"scneario-insert": v1.BenchmarkJobAvailable,
							"scneario-search": v1.BenchmarkJobAvailable,
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
			o := &operator{
				scenarios: test.fields.scenarios,
			}

			got := o.getAtomicScenario()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_operator_getAtomicBenchJob(t *testing.T) {
	t.Parallel()
	type fields struct {
		benchjobs *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
	}
	type want struct {
		want map[string]*v1.ValdBenchmarkJob
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, map[string]*v1.ValdBenchmarkJob) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
			o := &operator{
				benchjobs: test.fields.benchjobs,
			}

			got := o.getAtomicBenchJob()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_operator_getAtomicJob(t *testing.T) {
	t.Parallel()
	type fields struct {
		jobs *atomic.Pointer[map[string]string]
	}
	type want struct {
		want map[string]string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, map[string]string) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
			o := &operator{
				jobs: test.fields.jobs,
			}

			got := o.getAtomicJob()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_operator_jobReconcile(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context
		jobList map[string][]job.Job
	}
	type fields struct {
		jobNamespace       string
		jobImage           string
		jobImagePullPolicy string
		scenarios          *atomic.Pointer[map[string]*scenario]
		benchjobs          *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
		jobs               *atomic.Pointer[map[string]string]
		ctrl               k8s.Controller
	}
	type want struct {
		want map[string]string
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, map[string]string) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got map[string]string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success when the length of jobList is 0.",
				args: args{
					ctx:     ctx,
					jobList: map[string][]job.Job{},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImage:           "vdaas/vald-benchmark-job",
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success with new job whose namespace is same as jobNamespace and deleted job by etcd",
				args: args{
					ctx: ctx,
					jobList: map[string][]job.Job{
						"scenario-insert": {
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "scenario-insert",
									Namespace: "default",
								},
								Status: job.JobStatus{
									Active: 1,
								},
							},
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImage:           "vdaas/vald-benchmark-job",
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success with completed job whose namespace is same as jobNamespace",
				args: args{
					ctx: ctx,
					jobList: map[string][]job.Job{
						"scenario-insert": {
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "scenario-completed-insert",
									Namespace: "default",
								},
								Status: job.JobStatus{
									Active:    0,
									Succeeded: 1,
									CompletionTime: func() *metav1.Time {
										t := &metav1.Time{
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
					jobImage:           "vdaas/vald-benchmark-job",
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success with job whose namespace is not same as jobNamespace",
				args: args{
					ctx: ctx,
					jobList: map[string][]job.Job{
						"scenario-insert": {
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:      "scenario-insert",
									Namespace: "benchmark",
								},
								Status: job.JobStatus{
									Active: 1,
								},
							},
						},
					},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImage:           "vdaas/vald-benchmark-job",
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
				jobImage:           test.fields.jobImage,
				jobImagePullPolicy: test.fields.jobImagePullPolicy,
				benchjobs:          test.fields.benchjobs,
				jobs:               test.fields.jobs,
				ctrl:               test.fields.ctrl,
			}

			o.jobReconcile(test.args.ctx, test.args.jobList)
			got := o.getAtomicJob()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_operator_benchJobReconcile(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx          context.Context
		benchJobList map[string]v1.ValdBenchmarkJob
	}
	type fields struct {
		jobNamespace       string
		jobImage           string
		jobImagePullPolicy string
		scenarios          *atomic.Pointer[map[string]*scenario]
		benchjobs          *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
		jobs               *atomic.Pointer[map[string]string]
		ctrl               k8s.Controller
	}
	type want struct {
		scenarios map[string]*scenario
		benchjobs map[string]*v1.ValdBenchmarkJob
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, map[string]*scenario, map[string]*v1.ValdBenchmarkJob) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success when benchJobList is empty",
				args: args{
					ctx:          ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImage:           "vdaas/vald-benchmark-job",
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success when benchJobList has new benchmark Job with owner reference (reconcile after submitted scenario)",
				args: args{
					ctx: ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: metav1.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []metav1.OwnerReference{
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
					jobImage:           "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									ObjectMeta: metav1.ObjectMeta{
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
								ObjectMeta: metav1.ObjectMeta{
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
							ObjectMeta: metav1.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []metav1.OwnerReference{
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success when benchJobList has updated benchmark Job with owner reference (reconcile after updated scenario)",
				args: args{
					ctx: ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: metav1.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []metav1.OwnerReference{
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
					jobImage:           "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									ObjectMeta: metav1.ObjectMeta{
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
								ObjectMeta: metav1.ObjectMeta{
									Name:      "scenario-insert",
									Namespace: "default",
									OwnerReferences: []metav1.OwnerReference{
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
								ObjectMeta: metav1.ObjectMeta{
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
							ObjectMeta: metav1.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []metav1.OwnerReference{
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success when benchJobList has updated benchmark Job with owner reference (reconcile after updated job status)",
				args: args{
					ctx: ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: metav1.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []metav1.OwnerReference{
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
					jobImage:           "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									ObjectMeta: metav1.ObjectMeta{
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
								ObjectMeta: metav1.ObjectMeta{
									Name:      "scenario-insert",
									Namespace: "default",
									OwnerReferences: []metav1.OwnerReference{
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
								ObjectMeta: metav1.ObjectMeta{
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
							ObjectMeta: metav1.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []metav1.OwnerReference{
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success when benchJobList has new benchmark Job with owner reference and benchJob has deleted job",
				args: args{
					ctx: ctx,
					benchJobList: map[string]v1.ValdBenchmarkJob{
						"scenario-insert": {
							ObjectMeta: metav1.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []metav1.OwnerReference{
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
					jobImage:           "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									ObjectMeta: metav1.ObjectMeta{
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
								ObjectMeta: metav1.ObjectMeta{
									Name:      "scenario-insert",
									Namespace: "default",
									OwnerReferences: []metav1.OwnerReference{
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
								ObjectMeta: metav1.ObjectMeta{
									Name:      "scenario-deleted-insert",
									Namespace: "default",
									OwnerReferences: []metav1.OwnerReference{
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
								ObjectMeta: metav1.ObjectMeta{
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
							ObjectMeta: metav1.ObjectMeta{
								Name:      "scenario-insert",
								Namespace: "default",
								OwnerReferences: []metav1.OwnerReference{
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
				jobImage:           test.fields.jobImage,
				jobImagePullPolicy: test.fields.jobImagePullPolicy,
				scenarios:          test.fields.scenarios,
				benchjobs:          test.fields.benchjobs,
				jobs:               test.fields.jobs,
				ctrl:               test.fields.ctrl,
			}

			o.benchJobReconcile(test.args.ctx, test.args.benchJobList)
			gotS := o.getAtomicScenario()
			gotJ := o.getAtomicBenchJob()
			if err := checkFunc(test.want, gotS, gotJ); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_operator_benchScenarioReconcile(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx          context.Context
		scenarioList map[string]v1.ValdBenchmarkScenario
	}
	type fields struct {
		jobNamespace       string
		jobImage           string
		jobImagePullPolicy string
		scenarios          *atomic.Pointer[map[string]*scenario]
		benchjobs          *atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
		jobs               *atomic.Pointer[map[string]string]
		ctrl               k8s.Controller
	}
	type want struct {
		want map[string]*scenario
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, map[string]*scenario) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success with scenarioList is empty",
				args: args{
					ctx:          ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{},
				},
				fields: fields{
					jobNamespace:       "default",
					jobImage:           "vdaas/vald-benchmark-job",
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success with scenarioList has new scenario with no scenario has been applied yet.",
				args: args{
					ctx: ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{
						"scenario": {
							ObjectMeta: metav1.ObjectMeta{
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
					jobImage:           "vdaas/vald-benchmark-job",
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
								ObjectMeta: metav1.ObjectMeta{
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success with scenarioList has only status updated scenario.",
				args: args{
					ctx: ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{
						"scenario": {
							ObjectMeta: metav1.ObjectMeta{
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
					jobImage:           "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									ObjectMeta: metav1.ObjectMeta{
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
								ObjectMeta: metav1.ObjectMeta{
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success with scenarioList has updated scenario when job is already running",
				args: args{
					ctx: ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{
						"scenario": {
							ObjectMeta: metav1.ObjectMeta{
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
					jobImage:           "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									ObjectMeta: metav1.ObjectMeta{
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
								ObjectMeta: metav1.ObjectMeta{
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
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "success with scenarioList has another scenario when scenario is already running",
				args: args{
					ctx: ctx,
					scenarioList: map[string]v1.ValdBenchmarkScenario{
						"scenario-v2": {
							ObjectMeta: metav1.ObjectMeta{
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
					jobImage:           "vdaas/vald-benchmark-job",
					jobImagePullPolicy: "Always",
					scenarios: func() *atomic.Pointer[map[string]*scenario] {
						ap := atomic.Pointer[map[string]*scenario]{}
						m := map[string]*scenario{
							"scenario": {
								Crd: &v1.ValdBenchmarkScenario{
									ObjectMeta: metav1.ObjectMeta{
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
								ObjectMeta: metav1.ObjectMeta{
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
				jobImage:           test.fields.jobImage,
				jobImagePullPolicy: test.fields.jobImagePullPolicy,
				scenarios:          test.fields.scenarios,
				benchjobs:          test.fields.benchjobs,
				jobs:               test.fields.jobs,
				ctrl:               test.fields.ctrl,
			}

			o.benchScenarioReconcile(test.args.ctx, test.args.scenarioList)
			got := o.getAtomicScenario()
			t.Log(got["scenario"])
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

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
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
				ObjectMeta: metav1.ObjectMeta{
					Name:      "scenario",
					Namespace: "default",
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
			ObjectMeta: metav1.ObjectMeta{
				Name:      "scenario-insert",
				Namespace: "default",
				OwnerReferences: []metav1.OwnerReference{
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
			ObjectMeta: metav1.ObjectMeta{
				Name:      "scenario-search",
				Namespace: "default",
				OwnerReferences: []metav1.OwnerReference{
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
			ObjectMeta: metav1.ObjectMeta{
				Name:      "scenario-update",
				Namespace: "default",
				OwnerReferences: []metav1.OwnerReference{
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
				name: "return nil with no mismatch atmoics",
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
				name: "return mismatch error when scneario and job has atomic and benchJob has no atomic",
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
				val := v1.ValdBenchmarkJob{}
				val = *v
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
				val := v1.ValdBenchmarkJob{}
				val = *v
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
				val := v1.ValdBenchmarkJob{}
				val = *v
				benchJobMap[k] = &val
			}
			var ors []metav1.OwnerReference
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
			o := &operator{
				scenarios: test.fields.scenarios,
				benchjobs: test.fields.benchjobs,
				jobs:      test.fields.jobs,
			}

			err := o.checkAtomics()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
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
// 		// {
// 		// 	name: "test_case_1",
// 		// 	args: args{
// 		// 		opts: nil,
// 		// 	},
// 		// 	want: want{
// 		// 		want: func() Operator {
// 		// 			o := &operator{
// 		// 				jobNamespace:       "default",
// 		// 				jobImage:           "vdaas/vald-benchmark-job",
// 		// 				jobImagePullPolicy: "Always",
// 		// 				rcd:                10 * time.Second,
// 		// 			}
// 		// 			return o
// 		// 		}(),
// 		// 	},
// 		// 	checkFunc: defaultCheckFunc,
// 		// 	beforeFunc: func(t *testing.T, args args) {
// 		// 		t.Helper()
// 		// 	},
// 		// 	afterFunc: func(t *testing.T, args args) {
// 		// 		t.Helper()
// 		// 	},
// 		// },
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
// func Test_operator_PreStart(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
// 		       },
// 		       fields: fields {
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 			}
//
// 			err := o.PreStart(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_operator_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 			}
//
// 			got, err := o.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }

// func Test_operator_initCtrl(t *testing.T) {
// 	type fields struct {
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
//
// 	tests := []test{
// 		// {
// 		// 	name: "test_case_1",
// 		// 	fields: fields{
// 		// 		jobNamespace:       "default",
// 		// 		jobImage:           "vdaas/vald-benchmark-job",
// 		// 		jobImagePullPolicy: "Always",
// 		// 		// scenarios:nil,
// 		// 		// benchjobs:nil,
// 		// 		// jobs:nil,
// 		// 		// rcd:nil,
// 		// 		eg:   nil,
// 		// 		ctrl: nil,
// 		// 	},
// 		// 	want: want{
// 		// 		err: errors.New("hoge"),
// 		// 	},
// 		// 	checkFunc: defaultCheckFunc,
// 		// 	beforeFunc: func(t *testing.T) {
// 		// 		t.Helper()
// 		// 	},
// 		// 	afterFunc: func(t *testing.T) {
// 		// 		t.Helper()
// 		// 	},
// 		// },
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 			}
//
// 			err := o.initCtrl()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// func Test_operator_deleteBenchmarkJob(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		name       string
// 		generation int64
// 	}
// 	type fields struct {
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
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
// 		ctx        context.Context
// 		name       string
// 		generation int64
// 	}
// 	type fields struct {
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 			}
//
// 			err := o.deleteJob(test.args.ctx, test.args.name, test.args.generation)
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
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
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
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
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
// func Test_operator_updateBenchmarkScenarioStatus(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		ss  map[string]v1.ValdBenchmarkScenarioStatus
// 	}
// 	type fields struct {
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
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
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 		           jobNamespace:"",
// 		           jobImage:"",
// 		           jobImagePullPolicy:"",
// 		           scenarios:nil,
// 		           benchjobs:nil,
// 		           jobs:nil,
// 		           rcd:nil,
// 		           eg:nil,
// 		           ctrl:nil,
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
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
// 		jobNamespace       string
// 		jobImage           string
// 		jobImagePullPolicy string
// 		scenarios          atomic.Pointer[map[string]*scenario]
// 		benchjobs          atomic.Pointer[map[string]*v1.ValdBenchmarkJob]
// 		jobs               atomic.Pointer[map[string]string]
// 		rcd                time.Duration
// 		eg                 errgroup.Group
// 		ctrl               k8s.Controller
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
// 		// func() test {
// 		// 	return test{
// 		// 		name: "test_case_2",
// 		// 		args: args{
// 		// 			ctx:  nil,
// 		// 			jobs: nil,
// 		// 		},
// 		// 		fields: fields{
// 		// 			jobNamespace:       "",
// 		// 			jobImage:           "",
// 		// 			jobImagePullPolicy: "",
// 		// 			scenarios:          nil,
// 		// 			benchjobs:          nil,
// 		// 			jobs:               nil,
// 		// 			rcd:                nil,
// 		// 			eg:                 nil,
// 		// 			ctrl:               nil,
// 		// 		},
// 		// 		want:      want{},
// 		// 		checkFunc: defaultCheckFunc,
// 		// 		beforeFunc: func(t *testing.T, args args) {
// 		// 			t.Helper()
// 		// 		},
// 		// 		afterFunc: func(t *testing.T, args args) {
// 		// 			t.Helper()
// 		// 		},
// 		// 	}
// 		// }(),
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
// 				jobNamespace:       test.fields.jobNamespace,
// 				jobImage:           test.fields.jobImage,
// 				jobImagePullPolicy: test.fields.jobImagePullPolicy,
// 				scenarios:          test.fields.scenarios,
// 				benchjobs:          test.fields.benchjobs,
// 				jobs:               test.fields.jobs,
// 				rcd:                test.fields.rcd,
// 				eg:                 test.fields.eg,
// 				ctrl:               test.fields.ctrl,
// 			}
//
// 			err := o.checkJobsStatus(test.args.ctx, test.args.jobs)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
