// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// NOT IMPLEMENTED BELOW
//
// func Test_jobType_String(t *testing.T) {
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		jt         jobType
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
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
//
// 			got := test.jt.String()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Job
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Job, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Job, err error) error {
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
// func Test_job_PreStart(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		dataset            *config.BenchmarkDataset
// 		jobType            jobType
// 		jobFunc            func(context.Context, chan error) error
// 		insertConfig       *config.InsertConfig
// 		updateConfig       *config.UpdateConfig
// 		upsertConfig       *config.UpsertConfig
// 		searchConfig       *config.SearchConfig
// 		removeConfig       *config.RemoveConfig
// 		objectConfig       *config.ObjectConfig
// 		client             vald.Client
// 		hdf5               hdf5.Data
// 		beforeJobName      string
// 		beforeJobNamespace string
// 		k8sClient          client.Client
// 		beforeJobDur       time.Duration
// 		limiter            rate.Limiter
// 		rps                int
// 		concurrencyLimit   int
// 		timeout            time.Duration
// 		timestamp          int64
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
// 		           eg:nil,
// 		           dataset:nil,
// 		           jobType:nil,
// 		           jobFunc:nil,
// 		           insertConfig:nil,
// 		           updateConfig:nil,
// 		           upsertConfig:nil,
// 		           searchConfig:nil,
// 		           removeConfig:nil,
// 		           objectConfig:nil,
// 		           client:nil,
// 		           hdf5:nil,
// 		           beforeJobName:"",
// 		           beforeJobNamespace:"",
// 		           k8sClient:nil,
// 		           beforeJobDur:nil,
// 		           limiter:nil,
// 		           rps:0,
// 		           concurrencyLimit:0,
// 		           timeout:nil,
// 		           timestamp:0,
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
// 		           dataset:nil,
// 		           jobType:nil,
// 		           jobFunc:nil,
// 		           insertConfig:nil,
// 		           updateConfig:nil,
// 		           upsertConfig:nil,
// 		           searchConfig:nil,
// 		           removeConfig:nil,
// 		           objectConfig:nil,
// 		           client:nil,
// 		           hdf5:nil,
// 		           beforeJobName:"",
// 		           beforeJobNamespace:"",
// 		           k8sClient:nil,
// 		           beforeJobDur:nil,
// 		           limiter:nil,
// 		           rps:0,
// 		           concurrencyLimit:0,
// 		           timeout:nil,
// 		           timestamp:0,
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
// 			j := &job{
// 				eg:                 test.fields.eg,
// 				dataset:            test.fields.dataset,
// 				jobType:            test.fields.jobType,
// 				jobFunc:            test.fields.jobFunc,
// 				insertConfig:       test.fields.insertConfig,
// 				updateConfig:       test.fields.updateConfig,
// 				upsertConfig:       test.fields.upsertConfig,
// 				searchConfig:       test.fields.searchConfig,
// 				removeConfig:       test.fields.removeConfig,
// 				objectConfig:       test.fields.objectConfig,
// 				client:             test.fields.client,
// 				hdf5:               test.fields.hdf5,
// 				beforeJobName:      test.fields.beforeJobName,
// 				beforeJobNamespace: test.fields.beforeJobNamespace,
// 				k8sClient:          test.fields.k8sClient,
// 				beforeJobDur:       test.fields.beforeJobDur,
// 				limiter:            test.fields.limiter,
// 				rps:                test.fields.rps,
// 				concurrencyLimit:   test.fields.concurrencyLimit,
// 				timeout:            test.fields.timeout,
// 				timestamp:          test.fields.timestamp,
// 			}
//
// 			err := j.PreStart(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_job_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		dataset            *config.BenchmarkDataset
// 		jobType            jobType
// 		jobFunc            func(context.Context, chan error) error
// 		insertConfig       *config.InsertConfig
// 		updateConfig       *config.UpdateConfig
// 		upsertConfig       *config.UpsertConfig
// 		searchConfig       *config.SearchConfig
// 		removeConfig       *config.RemoveConfig
// 		objectConfig       *config.ObjectConfig
// 		client             vald.Client
// 		hdf5               hdf5.Data
// 		beforeJobName      string
// 		beforeJobNamespace string
// 		k8sClient          client.Client
// 		beforeJobDur       time.Duration
// 		limiter            rate.Limiter
// 		rps                int
// 		concurrencyLimit   int
// 		timeout            time.Duration
// 		timestamp          int64
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
// 		           dataset:nil,
// 		           jobType:nil,
// 		           jobFunc:nil,
// 		           insertConfig:nil,
// 		           updateConfig:nil,
// 		           upsertConfig:nil,
// 		           searchConfig:nil,
// 		           removeConfig:nil,
// 		           objectConfig:nil,
// 		           client:nil,
// 		           hdf5:nil,
// 		           beforeJobName:"",
// 		           beforeJobNamespace:"",
// 		           k8sClient:nil,
// 		           beforeJobDur:nil,
// 		           limiter:nil,
// 		           rps:0,
// 		           concurrencyLimit:0,
// 		           timeout:nil,
// 		           timestamp:0,
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
// 		           dataset:nil,
// 		           jobType:nil,
// 		           jobFunc:nil,
// 		           insertConfig:nil,
// 		           updateConfig:nil,
// 		           upsertConfig:nil,
// 		           searchConfig:nil,
// 		           removeConfig:nil,
// 		           objectConfig:nil,
// 		           client:nil,
// 		           hdf5:nil,
// 		           beforeJobName:"",
// 		           beforeJobNamespace:"",
// 		           k8sClient:nil,
// 		           beforeJobDur:nil,
// 		           limiter:nil,
// 		           rps:0,
// 		           concurrencyLimit:0,
// 		           timeout:nil,
// 		           timestamp:0,
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
// 			j := &job{
// 				eg:                 test.fields.eg,
// 				dataset:            test.fields.dataset,
// 				jobType:            test.fields.jobType,
// 				jobFunc:            test.fields.jobFunc,
// 				insertConfig:       test.fields.insertConfig,
// 				updateConfig:       test.fields.updateConfig,
// 				upsertConfig:       test.fields.upsertConfig,
// 				searchConfig:       test.fields.searchConfig,
// 				removeConfig:       test.fields.removeConfig,
// 				objectConfig:       test.fields.objectConfig,
// 				client:             test.fields.client,
// 				hdf5:               test.fields.hdf5,
// 				beforeJobName:      test.fields.beforeJobName,
// 				beforeJobNamespace: test.fields.beforeJobNamespace,
// 				k8sClient:          test.fields.k8sClient,
// 				beforeJobDur:       test.fields.beforeJobDur,
// 				limiter:            test.fields.limiter,
// 				rps:                test.fields.rps,
// 				concurrencyLimit:   test.fields.concurrencyLimit,
// 				timeout:            test.fields.timeout,
// 				timestamp:          test.fields.timestamp,
// 			}
//
// 			got, err := j.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_job_Stop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		dataset            *config.BenchmarkDataset
// 		jobType            jobType
// 		jobFunc            func(context.Context, chan error) error
// 		insertConfig       *config.InsertConfig
// 		updateConfig       *config.UpdateConfig
// 		upsertConfig       *config.UpsertConfig
// 		searchConfig       *config.SearchConfig
// 		removeConfig       *config.RemoveConfig
// 		objectConfig       *config.ObjectConfig
// 		client             vald.Client
// 		hdf5               hdf5.Data
// 		beforeJobName      string
// 		beforeJobNamespace string
// 		k8sClient          client.Client
// 		beforeJobDur       time.Duration
// 		limiter            rate.Limiter
// 		rps                int
// 		concurrencyLimit   int
// 		timeout            time.Duration
// 		timestamp          int64
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
// 		           eg:nil,
// 		           dataset:nil,
// 		           jobType:nil,
// 		           jobFunc:nil,
// 		           insertConfig:nil,
// 		           updateConfig:nil,
// 		           upsertConfig:nil,
// 		           searchConfig:nil,
// 		           removeConfig:nil,
// 		           objectConfig:nil,
// 		           client:nil,
// 		           hdf5:nil,
// 		           beforeJobName:"",
// 		           beforeJobNamespace:"",
// 		           k8sClient:nil,
// 		           beforeJobDur:nil,
// 		           limiter:nil,
// 		           rps:0,
// 		           concurrencyLimit:0,
// 		           timeout:nil,
// 		           timestamp:0,
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
// 		           dataset:nil,
// 		           jobType:nil,
// 		           jobFunc:nil,
// 		           insertConfig:nil,
// 		           updateConfig:nil,
// 		           upsertConfig:nil,
// 		           searchConfig:nil,
// 		           removeConfig:nil,
// 		           objectConfig:nil,
// 		           client:nil,
// 		           hdf5:nil,
// 		           beforeJobName:"",
// 		           beforeJobNamespace:"",
// 		           k8sClient:nil,
// 		           beforeJobDur:nil,
// 		           limiter:nil,
// 		           rps:0,
// 		           concurrencyLimit:0,
// 		           timeout:nil,
// 		           timestamp:0,
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
// 			j := &job{
// 				eg:                 test.fields.eg,
// 				dataset:            test.fields.dataset,
// 				jobType:            test.fields.jobType,
// 				jobFunc:            test.fields.jobFunc,
// 				insertConfig:       test.fields.insertConfig,
// 				updateConfig:       test.fields.updateConfig,
// 				upsertConfig:       test.fields.upsertConfig,
// 				searchConfig:       test.fields.searchConfig,
// 				removeConfig:       test.fields.removeConfig,
// 				objectConfig:       test.fields.objectConfig,
// 				client:             test.fields.client,
// 				hdf5:               test.fields.hdf5,
// 				beforeJobName:      test.fields.beforeJobName,
// 				beforeJobNamespace: test.fields.beforeJobNamespace,
// 				k8sClient:          test.fields.k8sClient,
// 				beforeJobDur:       test.fields.beforeJobDur,
// 				limiter:            test.fields.limiter,
// 				rps:                test.fields.rps,
// 				concurrencyLimit:   test.fields.concurrencyLimit,
// 				timeout:            test.fields.timeout,
// 				timestamp:          test.fields.timestamp,
// 			}
//
// 			err := j.Stop(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_calcRecall(t *testing.T) {
// 	type args struct {
// 		linearRes *payload.Search_Response
// 		searchRes *payload.Search_Response
// 	}
// 	type want struct {
// 		wantRecall float64
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, float64) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRecall float64) error {
// 		if !reflect.DeepEqual(gotRecall, w.wantRecall) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRecall, w.wantRecall)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           linearRes:nil,
// 		           searchRes:nil,
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
// 		           linearRes:nil,
// 		           searchRes:nil,
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
// 			gotRecall := calcRecall(test.args.linearRes, test.args.searchRes)
// 			if err := checkFunc(test.want, gotRecall); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_addNoiseToVec(t *testing.T) {
// 	type args struct {
// 		oVec []float32
// 	}
// 	type want struct {
// 		want []float32
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []float32) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []float32) error {
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
// 		           oVec:nil,
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
// 		           oVec:nil,
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
// 			got := addNoiseToVec(test.args.oVec)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }