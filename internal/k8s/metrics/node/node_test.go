//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package node provides kubernetes node information and preriodically update
package node

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}
	type want struct {
		want NodeWatcher
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, NodeWatcher) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got NodeWatcher) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           opts: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := New(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_reconciler_Reconcile(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req reconcile.Request
	}
	type fields struct {
		mgr         manager.Manager
		name        string
		onError     func(err error)
		onReconcile func(nodeList map[string]Node)
	}
	type want struct {
		wantRes reconcile.Result
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, reconcile.Result, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes reconcile.Result, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx: nil,
		           req: nil,
		           },
		           fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &reconciler{
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			gotRes, err := r.Reconcile(test.args.ctx, test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_reconciler_GetName(t *testing.T) {
	t.Parallel()
	type fields struct {
		mgr         manager.Manager
		name        string
		onError     func(err error)
		onReconcile func(nodeList map[string]Node)
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &reconciler{
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			got := r.GetName()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_reconciler_NewReconciler(t *testing.T) {
	t.Parallel()
	type args struct {
		mgr manager.Manager
	}
	type fields struct {
		mgr         manager.Manager
		name        string
		onError     func(err error)
		onReconcile func(nodeList map[string]Node)
	}
	type want struct {
		want reconcile.Reconciler
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, reconcile.Reconciler) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got reconcile.Reconciler) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           mgr: nil,
		       },
		       fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           mgr: nil,
		           },
		           fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &reconciler{
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			got := r.NewReconciler(test.args.mgr)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_reconciler_For(t *testing.T) {
	t.Parallel()
	type fields struct {
		mgr         manager.Manager
		name        string
		onError     func(err error)
		onReconcile func(nodeList map[string]Node)
	}
	type want struct {
		want  client.Object
		want1 []builder.ForOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, client.Object, []builder.ForOption) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got client.Object, got1 []builder.ForOption) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &reconciler{
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			got, got1 := r.For()
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_reconciler_Owns(t *testing.T) {
	t.Parallel()
	type fields struct {
		mgr         manager.Manager
		name        string
		onError     func(err error)
		onReconcile func(nodeList map[string]Node)
	}
	type want struct {
		want  client.Object
		want1 []builder.OwnsOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, client.Object, []builder.OwnsOption) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got client.Object, got1 []builder.OwnsOption) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &reconciler{
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			got, got1 := r.Owns()
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_reconciler_Watches(t *testing.T) {
	t.Parallel()
	type fields struct {
		mgr         manager.Manager
		name        string
		onError     func(err error)
		onReconcile func(nodeList map[string]Node)
	}
	type want struct {
		want  *source.Kind
		want1 handler.EventHandler
		want2 []builder.WatchesOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *source.Kind, handler.EventHandler, []builder.WatchesOption) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *source.Kind, got1 handler.EventHandler, got2 []builder.WatchesOption) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		if !reflect.DeepEqual(got2, w.want2) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got2, w.want2)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           mgr: nil,
		           name: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &reconciler{
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			got, got1, got2 := r.Watches()
			if err := test.checkFunc(test.want, got, got1, got2); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
