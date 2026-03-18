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

// NOT IMPLEMENTED BELOW
//
// func Test_job_remove(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		ech chan error
// 	}
// 	type fields struct {
// 		eg                 errgroup.Group
// 		limiter            rate.Limiter
// 		k8sClient          client.Client
// 		hdf5               hdf5.Data
// 		client             vald.Client
// 		insertConfig       *config.InsertConfig
// 		updateConfig       *config.UpdateConfig
// 		upsertConfig       *config.UpsertConfig
// 		searchConfig       *config.SearchConfig
// 		removeConfig       *config.RemoveConfig
// 		objectConfig       *config.ObjectConfig
// 		jobFunc            func(context.Context, chan error) error
// 		dataset            *config.BenchmarkDataset
// 		meta               grpc.MD
// 		beforeJobNamespace string
// 		beforeJobName      string
// 		beforeJobDur       time.Duration
// 		jobType            jobType
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
// 		           ech:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           limiter:nil,
// 		           k8sClient:nil,
// 		           hdf5:nil,
// 		           client:nil,
// 		           insertConfig:nil,
// 		           updateConfig:nil,
// 		           upsertConfig:nil,
// 		           searchConfig:nil,
// 		           removeConfig:nil,
// 		           objectConfig:nil,
// 		           jobFunc:nil,
// 		           dataset:nil,
// 		           meta:nil,
// 		           beforeJobNamespace:"",
// 		           beforeJobName:"",
// 		           beforeJobDur:nil,
// 		           jobType:nil,
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
// 		           ech:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           limiter:nil,
// 		           k8sClient:nil,
// 		           hdf5:nil,
// 		           client:nil,
// 		           insertConfig:nil,
// 		           updateConfig:nil,
// 		           upsertConfig:nil,
// 		           searchConfig:nil,
// 		           removeConfig:nil,
// 		           objectConfig:nil,
// 		           jobFunc:nil,
// 		           dataset:nil,
// 		           meta:nil,
// 		           beforeJobNamespace:"",
// 		           beforeJobName:"",
// 		           beforeJobDur:nil,
// 		           jobType:nil,
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
// 				limiter:            test.fields.limiter,
// 				k8sClient:          test.fields.k8sClient,
// 				hdf5:               test.fields.hdf5,
// 				client:             test.fields.client,
// 				insertConfig:       test.fields.insertConfig,
// 				updateConfig:       test.fields.updateConfig,
// 				upsertConfig:       test.fields.upsertConfig,
// 				searchConfig:       test.fields.searchConfig,
// 				removeConfig:       test.fields.removeConfig,
// 				objectConfig:       test.fields.objectConfig,
// 				jobFunc:            test.fields.jobFunc,
// 				dataset:            test.fields.dataset,
// 				meta:               test.fields.meta,
// 				beforeJobNamespace: test.fields.beforeJobNamespace,
// 				beforeJobName:      test.fields.beforeJobName,
// 				beforeJobDur:       test.fields.beforeJobDur,
// 				jobType:            test.fields.jobType,
// 				rps:                test.fields.rps,
// 				concurrencyLimit:   test.fields.concurrencyLimit,
// 				timeout:            test.fields.timeout,
// 				timestamp:          test.fields.timestamp,
// 			}
//
// 			err := j.remove(test.args.ctx, test.args.ech)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
