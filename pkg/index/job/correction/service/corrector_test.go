// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Corrector
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Corrector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Corrector, err error) error {
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
// func Test_correct_StartClient(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			got, err := c.StartClient(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			err := c.Start(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_PreStop(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			err := c.PreStop(test.args.in0)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_NumberOfCheckedIndex(t *testing.T) {
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			got := c.NumberOfCheckedIndex()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_NumberOfCorrectedOldIndex(t *testing.T) {
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			got := c.NumberOfCorrectedOldIndex()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_NumberOfCorrectedReplication(t *testing.T) {
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			got := c.NumberOfCorrectedReplication()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_loadReplicaInfo(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		originAddr string
// 		id         string
// 		replicas   []string
// 		counts     map[string]*payload.Info_Index_Count
// 		ts         int64
// 		start      time.Time
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
// 	}
// 	type want struct {
// 		wantFound       map[string]*payload.Object_Timestamp
// 		wantSkipped     []string
// 		wantLatest      int64
// 		wantLatestAgent string
// 		err             error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, map[string]*payload.Object_Timestamp, []string, int64, string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotFound map[string]*payload.Object_Timestamp, gotSkipped []string, gotLatest int64, gotLatestAgent string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotFound, w.wantFound) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotFound, w.wantFound)
// 		}
// 		if !reflect.DeepEqual(gotSkipped, w.wantSkipped) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSkipped, w.wantSkipped)
// 		}
// 		if !reflect.DeepEqual(gotLatest, w.wantLatest) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLatest, w.wantLatest)
// 		}
// 		if !reflect.DeepEqual(gotLatestAgent, w.wantLatestAgent) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLatestAgent, w.wantLatestAgent)
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
// 		           originAddr:"",
// 		           id:"",
// 		           replicas:nil,
// 		           counts:nil,
// 		           ts:0,
// 		           start:time.Time{},
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           originAddr:"",
// 		           id:"",
// 		           replicas:nil,
// 		           counts:nil,
// 		           ts:0,
// 		           start:time.Time{},
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			gotFound, gotSkipped, gotLatest, gotLatestAgent, err := c.loadReplicaInfo(
// 				test.args.ctx,
// 				test.args.originAddr,
// 				test.args.id,
// 				test.args.replicas,
// 				test.args.counts,
// 				test.args.ts,
// 				test.args.start,
// 			)
// 			if err := checkFunc(test.want, gotFound, gotSkipped, gotLatest, gotLatestAgent, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_getLatestObject(t *testing.T) {
// 	type args struct {
// 		ctx         context.Context
// 		id          string
// 		addr        string
// 		latestAgent string
// 		latest      int64
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
// 	}
// 	type want struct {
// 		wantLatestObject *payload.Object_Vector
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Vector) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLatestObject *payload.Object_Vector) error {
// 		if !reflect.DeepEqual(gotLatestObject, w.wantLatestObject) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLatestObject, w.wantLatestObject)
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
// 		           id:"",
// 		           addr:"",
// 		           latestAgent:"",
// 		           latest:0,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           id:"",
// 		           addr:"",
// 		           latestAgent:"",
// 		           latest:0,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			gotLatestObject := c.getLatestObject(test.args.ctx, test.args.id, test.args.addr, test.args.latestAgent, test.args.latest)
// 			if err := checkFunc(test.want, gotLatestObject); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_correctTimestamp(t *testing.T) {
// 	type args struct {
// 		ctx          context.Context
// 		id           string
// 		latestObject *payload.Object_Vector
// 		found        map[string]*payload.Object_Timestamp
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           id:"",
// 		           latestObject:nil,
// 		           found:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           id:"",
// 		           latestObject:nil,
// 		           found:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			c.correctTimestamp(test.args.ctx, test.args.id, test.args.latestObject, test.args.found)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_correctOversupply(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		id       string
// 		selfAddr string
// 		debugMsg string
// 		found    map[string]*payload.Object_Timestamp
// 		diff     int
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
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
// 		           id:"",
// 		           selfAddr:"",
// 		           debugMsg:"",
// 		           found:nil,
// 		           diff:0,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           id:"",
// 		           selfAddr:"",
// 		           debugMsg:"",
// 		           found:nil,
// 		           diff:0,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			err := c.correctOversupply(test.args.ctx, test.args.id, test.args.selfAddr, test.args.debugMsg, test.args.found, test.args.diff)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_correct_correctShortage(t *testing.T) {
// 	type args struct {
// 		ctx          context.Context
// 		id           string
// 		selfAddr     string
// 		debugMsg     string
// 		latestObject *payload.Object_Vector
// 		found        map[string]*payload.Object_Timestamp
// 		diff         int
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		discoverer                   discoverer.Client
// 		gateway                      vc.Client
// 		checkedList                  pogreb.DB
// 		checkedIndexCount            atomic.Uint64
// 		correctedOldIndexCount       atomic.Uint64
// 		correctedReplicationCount    atomic.Uint64
// 		indexReplica                 int
// 		streamListConcurrency        int
// 		backgroundSyncInterval       time.Duration
// 		backgroundCompactionInterval time.Duration
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
// 		           id:"",
// 		           selfAddr:"",
// 		           debugMsg:"",
// 		           latestObject:nil,
// 		           found:nil,
// 		           diff:0,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 		           id:"",
// 		           selfAddr:"",
// 		           debugMsg:"",
// 		           latestObject:nil,
// 		           found:nil,
// 		           diff:0,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           discoverer:nil,
// 		           gateway:nil,
// 		           checkedList:nil,
// 		           checkedIndexCount:nil,
// 		           correctedOldIndexCount:nil,
// 		           correctedReplicationCount:nil,
// 		           indexReplica:0,
// 		           streamListConcurrency:0,
// 		           backgroundSyncInterval:nil,
// 		           backgroundCompactionInterval:nil,
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
// 			c := &correct{
// 				eg:                           test.fields.eg,
// 				discoverer:                   test.fields.discoverer,
// 				gateway:                      test.fields.gateway,
// 				checkedList:                  test.fields.checkedList,
// 				checkedIndexCount:            test.fields.checkedIndexCount,
// 				correctedOldIndexCount:       test.fields.correctedOldIndexCount,
// 				correctedReplicationCount:    test.fields.correctedReplicationCount,
// 				indexReplica:                 test.fields.indexReplica,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			err := c.correctShortage(test.args.ctx, test.args.id, test.args.selfAddr, test.args.debugMsg, test.args.latestObject, test.args.found, test.args.diff)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
