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
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func Test_extractAPIPorts(t *testing.T) {
	type args struct {
		ports []corev1.ServicePort
	}
	type want struct {
		ports []servicePort
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []servicePort) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got []servicePort) error {
		if !reflect.DeepEqual(got, w.ports) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.ports)
		}
		return nil
	}
	tests := []test{
		{
			name: "success to parse grpc",
			args: args{
				ports: []corev1.ServicePort{
					{
						Name: "grpc",
						Port: 8081,
					},
					{
						Name: "readiness",
						Port: 3000,
					},
					{
						Name: "liveness",
						Port: 3001,
					},
					{
						Name: "pprof",
						Port: 6060,
					},
				},
			},
			want: want{
				ports: []servicePort{
					{
						Name: "grpc",
						Port: 8081,
					},
				},
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "success to parse rest",
			args: args{
				ports: []corev1.ServicePort{
					{
						Name: "rest",
						Port: 8082,
					},
					{
						Name: "readiness",
						Port: 3000,
					},
					{
						Name: "liveness",
						Port: 3001,
					},
					{
						Name: "pprof",
						Port: 6060,
					},
				},
			},
			want: want{
				ports: []servicePort{
					{
						Name: "rest",
						Port: 8082,
					},
				},
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "success to parse grpc and rest",
			args: args{
				ports: []corev1.ServicePort{
					{
						Name: "rest",
						Port: 8082,
					},
					{
						Name: "grpc",
						Port: 8081,
					},
					{
						Name: "readiness",
						Port: 3000,
					},
					{
						Name: "liveness",
						Port: 3001,
					},
					{
						Name: "pprof",
						Port: 6060,
					},
				},
			},
			want: want{
				ports: []servicePort{
					{
						Name: "rest",
						Port: 8082,
					},
					{
						Name: "grpc",
						Port: 8081,
					},
				},
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "empty input returns empty",
			args: args{
				ports: nil,
			},
			want: want{
				ports: nil,
			},
			checkFunc: defaultCheckFunc,
		},
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

			got := extractAPIPorts(test.args.ports)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want SvcWatcher
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, SvcWatcher) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got SvcWatcher) error {
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
// 			got := New(test.args.opts...)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_reconciler_addListOpts(t *testing.T) {
// 	type args struct {
// 		opt client.ListOption
// 	}
// 	type fields struct {
// 		mgr         manager.Manager
// 		name        string
// 		namespace   string
// 		onError     func(err error)
// 		onReconcile func(svcs []Service)
// 		lopts       []client.ListOption
// 	}
// 	type want struct {
// 	}
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
// 		           opt:nil,
// 		       },
// 		       fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 		           opt:nil,
// 		           },
// 		           fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 			r := &reconciler{
// 				mgr:         test.fields.mgr,
// 				name:        test.fields.name,
// 				namespace:   test.fields.namespace,
// 				onError:     test.fields.onError,
// 				onReconcile: test.fields.onReconcile,
// 				lopts:       test.fields.lopts,
// 			}
//
// 			r.addListOpts(test.args.opt)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_reconciler_Reconcile(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		in1 reconcile.Request
// 	}
// 	type fields struct {
// 		mgr         manager.Manager
// 		name        string
// 		namespace   string
// 		onError     func(err error)
// 		onReconcile func(svcs []Service)
// 		lopts       []client.ListOption
// 	}
// 	type want struct {
// 		wantRes reconcile.Result
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, reconcile.Result, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes reconcile.Result, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           in1:nil,
// 		       },
// 		       fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 		           in1:nil,
// 		           },
// 		           fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 			r := &reconciler{
// 				mgr:         test.fields.mgr,
// 				name:        test.fields.name,
// 				namespace:   test.fields.namespace,
// 				onError:     test.fields.onError,
// 				onReconcile: test.fields.onReconcile,
// 				lopts:       test.fields.lopts,
// 			}
//
// 			gotRes, err := r.Reconcile(test.args.ctx, test.args.in1)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
//
// func Test_reconciler_GetName(t *testing.T) {
// 	type fields struct {
// 		mgr         manager.Manager
// 		name        string
// 		namespace   string
// 		onError     func(err error)
// 		onReconcile func(svcs []Service)
// 		lopts       []client.ListOption
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
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
// 		       fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 			r := &reconciler{
// 				mgr:         test.fields.mgr,
// 				name:        test.fields.name,
// 				namespace:   test.fields.namespace,
// 				onError:     test.fields.onError,
// 				onReconcile: test.fields.onReconcile,
// 				lopts:       test.fields.lopts,
// 			}
//
// 			got := r.GetName()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_reconciler_NewReconciler(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 		mgr manager.Manager
// 	}
// 	type fields struct {
// 		mgr         manager.Manager
// 		name        string
// 		namespace   string
// 		onError     func(err error)
// 		onReconcile func(svcs []Service)
// 		lopts       []client.ListOption
// 	}
// 	type want struct {
// 		want reconcile.Reconciler
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, reconcile.Reconciler) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got reconcile.Reconciler) error {
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
// 		           in0:nil,
// 		           mgr:nil,
// 		       },
// 		       fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 		           mgr:nil,
// 		           },
// 		           fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 			r := &reconciler{
// 				mgr:         test.fields.mgr,
// 				name:        test.fields.name,
// 				namespace:   test.fields.namespace,
// 				onError:     test.fields.onError,
// 				onReconcile: test.fields.onReconcile,
// 				lopts:       test.fields.lopts,
// 			}
//
// 			got := r.NewReconciler(test.args.in0, test.args.mgr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_reconciler_For(t *testing.T) {
// 	type fields struct {
// 		mgr         manager.Manager
// 		name        string
// 		namespace   string
// 		onError     func(err error)
// 		onReconcile func(svcs []Service)
// 		lopts       []client.ListOption
// 	}
// 	type want struct {
// 		want  client.Object
// 		want1 []builder.ForOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, client.Object, []builder.ForOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got client.Object, got1 []builder.ForOption) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		if !reflect.DeepEqual(got1, w.want1) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 			r := &reconciler{
// 				mgr:         test.fields.mgr,
// 				name:        test.fields.name,
// 				namespace:   test.fields.namespace,
// 				onError:     test.fields.onError,
// 				onReconcile: test.fields.onReconcile,
// 				lopts:       test.fields.lopts,
// 			}
//
// 			got, got1 := r.For()
// 			if err := checkFunc(test.want, got, got1); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_reconciler_Owns(t *testing.T) {
// 	type fields struct {
// 		mgr         manager.Manager
// 		name        string
// 		namespace   string
// 		onError     func(err error)
// 		onReconcile func(svcs []Service)
// 		lopts       []client.ListOption
// 	}
// 	type want struct {
// 		want  client.Object
// 		want1 []builder.OwnsOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, client.Object, []builder.OwnsOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got client.Object, got1 []builder.OwnsOption) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		if !reflect.DeepEqual(got1, w.want1) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 			r := &reconciler{
// 				mgr:         test.fields.mgr,
// 				name:        test.fields.name,
// 				namespace:   test.fields.namespace,
// 				onError:     test.fields.onError,
// 				onReconcile: test.fields.onReconcile,
// 				lopts:       test.fields.lopts,
// 			}
//
// 			got, got1 := r.Owns()
// 			if err := checkFunc(test.want, got, got1); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_reconciler_Watches(t *testing.T) {
// 	type fields struct {
// 		mgr         manager.Manager
// 		name        string
// 		namespace   string
// 		onError     func(err error)
// 		onReconcile func(svcs []Service)
// 		lopts       []client.ListOption
// 	}
// 	type want struct {
// 		want  client.Object
// 		want1 handler.EventHandler
// 		want2 []builder.WatchesOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, client.Object, handler.EventHandler, []builder.WatchesOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got client.Object, got1 handler.EventHandler, got2 []builder.WatchesOption) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		if !reflect.DeepEqual(got1, w.want1) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
// 		}
// 		if !reflect.DeepEqual(got2, w.want2) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got2, w.want2)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 		           mgr:nil,
// 		           name:"",
// 		           namespace:"",
// 		           onError:nil,
// 		           onReconcile:nil,
// 		           lopts:nil,
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
// 			r := &reconciler{
// 				mgr:         test.fields.mgr,
// 				name:        test.fields.name,
// 				namespace:   test.fields.namespace,
// 				onError:     test.fields.onError,
// 				onReconcile: test.fields.onReconcile,
// 				lopts:       test.fields.lopts,
// 			}
//
// 			got, got1, got2 := r.Watches()
// 			if err := checkFunc(test.want, got, got1, got2); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
