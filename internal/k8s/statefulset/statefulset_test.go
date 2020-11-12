package statefulset

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want StatefulSetWatcher
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, StatefulSetWatcher, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got StatefulSetWatcher, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reconciler_Reconcile(t *testing.T) {
	type args struct {
		req reconcile.Request
	}
	type fields struct {
		ctx         context.Context
		mgr         manager.Manager
		name        string
		namespace   string
		onError     func(err error)
		onReconcile func(rs map[string][]StatefulSet)
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
		           req: nil,
		       },
		       fields: fields {
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
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
		           req: nil,
		           },
		           fields: fields {
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				ctx:         test.fields.ctx,
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				namespace:   test.fields.namespace,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			gotRes, err := r.Reconcile(test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reconciler_GetName(t *testing.T) {
	type fields struct {
		ctx         context.Context
		mgr         manager.Manager
		name        string
		namespace   string
		onError     func(err error)
		onReconcile func(rs map[string][]StatefulSet)
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
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
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
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				ctx:         test.fields.ctx,
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				namespace:   test.fields.namespace,
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
	type args struct {
		ctx context.Context
		mgr manager.Manager
	}
	type fields struct {
		ctx         context.Context
		mgr         manager.Manager
		name        string
		namespace   string
		onError     func(err error)
		onReconcile func(rs map[string][]StatefulSet)
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
		           ctx: nil,
		           mgr: nil,
		       },
		       fields: fields {
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
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
		           mgr: nil,
		           },
		           fields: fields {
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				ctx:         test.fields.ctx,
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				namespace:   test.fields.namespace,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			got := r.NewReconciler(test.args.ctx, test.args.mgr)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reconciler_For(t *testing.T) {
	type fields struct {
		ctx         context.Context
		mgr         manager.Manager
		name        string
		namespace   string
		onError     func(err error)
		onReconcile func(rs map[string][]StatefulSet)
	}
	type want struct {
		want runtime.Object
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, runtime.Object) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got runtime.Object) error {
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
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
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
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				ctx:         test.fields.ctx,
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				namespace:   test.fields.namespace,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			got := r.For()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reconciler_Owns(t *testing.T) {
	type fields struct {
		ctx         context.Context
		mgr         manager.Manager
		name        string
		namespace   string
		onError     func(err error)
		onReconcile func(rs map[string][]StatefulSet)
	}
	type want struct {
		want runtime.Object
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, runtime.Object) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got runtime.Object) error {
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
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
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
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				ctx:         test.fields.ctx,
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				namespace:   test.fields.namespace,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			got := r.Owns()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_reconciler_Watches(t *testing.T) {
	type fields struct {
		ctx         context.Context
		mgr         manager.Manager
		name        string
		namespace   string
		onError     func(err error)
		onReconcile func(rs map[string][]StatefulSet)
	}
	type want struct {
		want  *source.Kind
		want1 handler.EventHandler
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *source.Kind, handler.EventHandler) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *source.Kind, got1 handler.EventHandler) error {
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
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
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
		           ctx: nil,
		           mgr: nil,
		           name: "",
		           namespace: "",
		           onError: nil,
		           onReconcile: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				ctx:         test.fields.ctx,
				mgr:         test.fields.mgr,
				name:        test.fields.name,
				namespace:   test.fields.namespace,
				onError:     test.fields.onError,
				onReconcile: test.fields.onReconcile,
			}

			got, got1 := r.Watches()
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
