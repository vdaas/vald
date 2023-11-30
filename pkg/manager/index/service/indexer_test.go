//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package service
package service

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		wantIdx Indexer
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Indexer, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotIdx Indexer, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotIdx, w.wantIdx) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIdx, w.wantIdx)
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
// 			gotIdx, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, gotIdx, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_index_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client                 discoverer.Client
// 		eg                     errgroup.Group
// 		creationPoolSize       uint32
// 		indexDuration          time.Duration
// 		indexDurationLimit     time.Duration
// 		saveIndexDurationLimit time.Duration
// 		saveIndexWaitDuration  time.Duration
// 		saveIndexTargetAddrCh  chan string
// 		concurrency            int
// 		indexInfos             indexInfos
// 		indexing               atomic.Value
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
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
// 		           client:nil,
// 		           eg:nil,
// 		           creationPoolSize:0,
// 		           indexDuration:nil,
// 		           indexDurationLimit:nil,
// 		           saveIndexDurationLimit:nil,
// 		           saveIndexWaitDuration:nil,
// 		           saveIndexTargetAddrCh:nil,
// 		           concurrency:0,
// 		           indexInfos:indexInfos{},
// 		           indexing:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
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
// 		           client:nil,
// 		           eg:nil,
// 		           creationPoolSize:0,
// 		           indexDuration:nil,
// 		           indexDurationLimit:nil,
// 		           saveIndexDurationLimit:nil,
// 		           saveIndexWaitDuration:nil,
// 		           saveIndexTargetAddrCh:nil,
// 		           concurrency:0,
// 		           indexInfos:indexInfos{},
// 		           indexing:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
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
// 			idx := &index{
// 				client:                 test.fields.client,
// 				eg:                     test.fields.eg,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexDuration:          test.fields.indexDuration,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				saveIndexWaitDuration:  test.fields.saveIndexWaitDuration,
// 				saveIndexTargetAddrCh:  test.fields.saveIndexTargetAddrCh,
// 				concurrency:            test.fields.concurrency,
// 				indexInfos:             test.fields.indexInfos,
// 				indexing:               test.fields.indexing,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			got, err := idx.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_index_IsIndexing(t *testing.T) {
// 	type fields struct {
// 		client                 discoverer.Client
// 		eg                     errgroup.Group
// 		creationPoolSize       uint32
// 		indexDuration          time.Duration
// 		indexDurationLimit     time.Duration
// 		saveIndexDurationLimit time.Duration
// 		saveIndexWaitDuration  time.Duration
// 		saveIndexTargetAddrCh  chan string
// 		concurrency            int
// 		indexInfos             indexInfos
// 		indexing               atomic.Value
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		           client:nil,
// 		           eg:nil,
// 		           creationPoolSize:0,
// 		           indexDuration:nil,
// 		           indexDurationLimit:nil,
// 		           saveIndexDurationLimit:nil,
// 		           saveIndexWaitDuration:nil,
// 		           saveIndexTargetAddrCh:nil,
// 		           concurrency:0,
// 		           indexInfos:indexInfos{},
// 		           indexing:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
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
// 		           client:nil,
// 		           eg:nil,
// 		           creationPoolSize:0,
// 		           indexDuration:nil,
// 		           indexDurationLimit:nil,
// 		           saveIndexDurationLimit:nil,
// 		           saveIndexWaitDuration:nil,
// 		           saveIndexTargetAddrCh:nil,
// 		           concurrency:0,
// 		           indexInfos:indexInfos{},
// 		           indexing:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
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
// 			idx := &index{
// 				client:                 test.fields.client,
// 				eg:                     test.fields.eg,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexDuration:          test.fields.indexDuration,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				saveIndexWaitDuration:  test.fields.saveIndexWaitDuration,
// 				saveIndexTargetAddrCh:  test.fields.saveIndexTargetAddrCh,
// 				concurrency:            test.fields.concurrency,
// 				indexInfos:             test.fields.indexInfos,
// 				indexing:               test.fields.indexing,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			got := idx.IsIndexing()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_index_NumberOfUUIDs(t *testing.T) {
// 	type fields struct {
// 		client                 discoverer.Client
// 		eg                     errgroup.Group
// 		creationPoolSize       uint32
// 		indexDuration          time.Duration
// 		indexDurationLimit     time.Duration
// 		saveIndexDurationLimit time.Duration
// 		saveIndexWaitDuration  time.Duration
// 		saveIndexTargetAddrCh  chan string
// 		concurrency            int
// 		indexInfos             indexInfos
// 		indexing               atomic.Value
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		want uint32
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint32) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint32) error {
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
// 		           client:nil,
// 		           eg:nil,
// 		           creationPoolSize:0,
// 		           indexDuration:nil,
// 		           indexDurationLimit:nil,
// 		           saveIndexDurationLimit:nil,
// 		           saveIndexWaitDuration:nil,
// 		           saveIndexTargetAddrCh:nil,
// 		           concurrency:0,
// 		           indexInfos:indexInfos{},
// 		           indexing:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
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
// 		           client:nil,
// 		           eg:nil,
// 		           creationPoolSize:0,
// 		           indexDuration:nil,
// 		           indexDurationLimit:nil,
// 		           saveIndexDurationLimit:nil,
// 		           saveIndexWaitDuration:nil,
// 		           saveIndexTargetAddrCh:nil,
// 		           concurrency:0,
// 		           indexInfos:indexInfos{},
// 		           indexing:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
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
// 			idx := &index{
// 				client:                 test.fields.client,
// 				eg:                     test.fields.eg,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexDuration:          test.fields.indexDuration,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				saveIndexWaitDuration:  test.fields.saveIndexWaitDuration,
// 				saveIndexTargetAddrCh:  test.fields.saveIndexTargetAddrCh,
// 				concurrency:            test.fields.concurrency,
// 				indexInfos:             test.fields.indexInfos,
// 				indexing:               test.fields.indexing,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			got := idx.NumberOfUUIDs()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_index_NumberOfUncommittedUUIDs(t *testing.T) {
// 	type fields struct {
// 		client                 discoverer.Client
// 		eg                     errgroup.Group
// 		creationPoolSize       uint32
// 		indexDuration          time.Duration
// 		indexDurationLimit     time.Duration
// 		saveIndexDurationLimit time.Duration
// 		saveIndexWaitDuration  time.Duration
// 		saveIndexTargetAddrCh  chan string
// 		concurrency            int
// 		indexInfos             indexInfos
// 		indexing               atomic.Value
// 		minUncommitted         uint32
// 		uuidsCount             uint32
// 		uncommittedUUIDsCount  uint32
// 	}
// 	type want struct {
// 		want uint32
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint32) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint32) error {
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
// 		           client:nil,
// 		           eg:nil,
// 		           creationPoolSize:0,
// 		           indexDuration:nil,
// 		           indexDurationLimit:nil,
// 		           saveIndexDurationLimit:nil,
// 		           saveIndexWaitDuration:nil,
// 		           saveIndexTargetAddrCh:nil,
// 		           concurrency:0,
// 		           indexInfos:indexInfos{},
// 		           indexing:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
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
// 		           client:nil,
// 		           eg:nil,
// 		           creationPoolSize:0,
// 		           indexDuration:nil,
// 		           indexDurationLimit:nil,
// 		           saveIndexDurationLimit:nil,
// 		           saveIndexWaitDuration:nil,
// 		           saveIndexTargetAddrCh:nil,
// 		           concurrency:0,
// 		           indexInfos:indexInfos{},
// 		           indexing:nil,
// 		           minUncommitted:0,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
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
// 			idx := &index{
// 				client:                 test.fields.client,
// 				eg:                     test.fields.eg,
// 				creationPoolSize:       test.fields.creationPoolSize,
// 				indexDuration:          test.fields.indexDuration,
// 				indexDurationLimit:     test.fields.indexDurationLimit,
// 				saveIndexDurationLimit: test.fields.saveIndexDurationLimit,
// 				saveIndexWaitDuration:  test.fields.saveIndexWaitDuration,
// 				saveIndexTargetAddrCh:  test.fields.saveIndexTargetAddrCh,
// 				concurrency:            test.fields.concurrency,
// 				indexInfos:             test.fields.indexInfos,
// 				indexing:               test.fields.indexing,
// 				minUncommitted:         test.fields.minUncommitted,
// 				uuidsCount:             test.fields.uuidsCount,
// 				uncommittedUUIDsCount:  test.fields.uncommittedUUIDsCount,
// 			}
//
// 			got := idx.NumberOfUncommittedUUIDs()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
