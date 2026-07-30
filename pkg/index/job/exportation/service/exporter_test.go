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
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/test/goleak"
)

// NOTE: StartClient/Start/doExportIndex paths that require a real vc.Client
// (github.com/vdaas/vald/internal/client/v1/client/vald.Client) gateway are
// intentionally not covered here: that interface embeds the entire Vald API
// surface (Flush/Index/Insert/Object/Remove/Search/Update/Upsert clients,
// dozens of RPC and streaming methods) and no reusable mock exists in
// internal/test/mock. Hand-writing a full mock only to exercise these paths
// would be scaffolding disproportionate to this task; New is covered below
// since it does not require a gateway.
func TestNew(t *testing.T) {
	t.Parallel()
	type want struct {
		err error
	}
	type test struct {
		want      want
		checkFunc func(*testing.T, want, Exporter, error)
		name      string
		opts      []Option
	}
	tests := []test{
		func() test {
			tmpDir := t.TempDir()
			return test{
				name: "returns Exporter and nil when the index directory and db file are created",
				opts: []Option{
					WithIndexPath(tmpDir),
				},
				want: want{
					err: nil,
				},
				checkFunc: func(t *testing.T, w want, e Exporter, err error) {
					t.Helper()
					if !errors.Is(err, w.err) {
						t.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%v\"", err, w.err)
					}
					if e == nil {
						t.Fatal("returned Exporter must not be nil")
					}
					files, err := file.ListInDir(tmpDir)
					if err != nil {
						t.Fatal(err)
					}
					if len(files) == 0 {
						t.Error("New must create a db file under the configured index path")
					}
					exp, ok := e.(*export)
					if !ok {
						t.Fatal("Exporter must be backed by *export")
					}
					if err := exp.storedVector.Close(true); err != nil {
						t.Errorf("failed to close stored vector db: %v", err)
					}
				},
			}
		}(),
		func() test {
			return test{
				name: "returns error when a critical option fails",
				opts: []Option{
					WithGateway(nil),
				},
				checkFunc: func(t *testing.T, _ want, e Exporter, err error) {
					t.Helper()
					var critOpt *errors.ErrCriticalOption
					if !errors.As(err, &critOpt) {
						t.Errorf("got_error: \"%v\", want a wrapped *errors.ErrCriticalOption", err)
					}
					if e != nil {
						t.Errorf("returned Exporter must be nil, got: %#v", e)
					}
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			blocker := file.Join(tmpDir, "blocker")
			f, err := os.Create(blocker)
			if err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
			return test{
				name: "returns error when the index path cannot be created",
				opts: []Option{
					WithIndexPath(file.Join(blocker, "sub")),
				},
				checkFunc: func(t *testing.T, _ want, e Exporter, err error) {
					t.Helper()
					if err == nil {
						t.Error("expected an error when the index path parent is a regular file")
					}
					if e != nil {
						t.Errorf("returned Exporter must be nil, got: %#v", e)
					}
				},
				want: want{},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			got, err := New(test.opts...)
			test.checkFunc(tt, test.want, got, err)
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_export_StartClient(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		gateway                      vc.Client
// 		storedVector                 pogreb.DB
// 		indexPath                    string
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
// 		           gateway:nil,
// 		           storedVector:nil,
// 		           indexPath:"",
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
// 		           gateway:nil,
// 		           storedVector:nil,
// 		           indexPath:"",
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
// 			e := &export{
// 				eg:                           test.fields.eg,
// 				gateway:                      test.fields.gateway,
// 				storedVector:                 test.fields.storedVector,
// 				indexPath:                    test.fields.indexPath,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			got, err := e.StartClient(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_export_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		gateway                      vc.Client
// 		storedVector                 pogreb.DB
// 		indexPath                    string
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
// 		           gateway:nil,
// 		           storedVector:nil,
// 		           indexPath:"",
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
// 		           gateway:nil,
// 		           storedVector:nil,
// 		           indexPath:"",
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
// 			e := &export{
// 				eg:                           test.fields.eg,
// 				gateway:                      test.fields.gateway,
// 				storedVector:                 test.fields.storedVector,
// 				indexPath:                    test.fields.indexPath,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			err := e.Start(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_export_doExportIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		gateway                      vc.Client
// 		storedVector                 pogreb.DB
// 		indexPath                    string
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
// 		           gateway:nil,
// 		           storedVector:nil,
// 		           indexPath:"",
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
// 		           gateway:nil,
// 		           storedVector:nil,
// 		           indexPath:"",
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
// 			e := &export{
// 				eg:                           test.fields.eg,
// 				gateway:                      test.fields.gateway,
// 				storedVector:                 test.fields.storedVector,
// 				indexPath:                    test.fields.indexPath,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			err := e.doExportIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_export_PreStop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg                           errgroup.Group
// 		gateway                      vc.Client
// 		storedVector                 pogreb.DB
// 		indexPath                    string
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
// 		           gateway:nil,
// 		           storedVector:nil,
// 		           indexPath:"",
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
// 		           gateway:nil,
// 		           storedVector:nil,
// 		           indexPath:"",
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
// 			e := &export{
// 				eg:                           test.fields.eg,
// 				gateway:                      test.fields.gateway,
// 				storedVector:                 test.fields.storedVector,
// 				indexPath:                    test.fields.indexPath,
// 				streamListConcurrency:        test.fields.streamListConcurrency,
// 				backgroundSyncInterval:       test.fields.backgroundSyncInterval,
// 				backgroundCompactionInterval: test.fields.backgroundCompactionInterval,
// 			}
//
// 			err := e.PreStop(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
