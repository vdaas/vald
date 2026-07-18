//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package resource

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newFakeAPI(t *testing.T, objs ...Object) kclient.WithWatch {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add corev1 scheme: %v", err)
	}
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func newConfigMapClient(
	t *testing.T, objs ...Object,
) *Client[*corev1.ConfigMap, *corev1.ConfigMapList] {
	t.Helper()
	return NewClient(newFakeAPI(t, objs...), new(corev1.ConfigMap), new(corev1.ConfigMapList))
}

func newConfigMap(name, namespace string, data map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Data:       data,
	}
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	type test struct {
		name      string
		getName   string
		objs      []Object
		wantErr   bool
		wantFound bool
	}

	tests := []test{
		{
			name:      "returns the object when it exists",
			objs:      []Object{newConfigMap("exists", "default", map[string]string{"k": "v"})},
			getName:   "exists",
			wantFound: true,
		},
		{
			name:    "returns NotFound when the object is missing",
			getName: "missing",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			oc := newConfigMapClient(t, tc.objs...)
			obj, err := oc.Get(context.Background(), tc.getName, "default")
			if tc.wantErr {
				if err == nil {
					t.Fatal("Get() error = nil, want error")
				}
				if !apierrors.IsNotFound(err) {
					t.Errorf("Get() error = %v, want NotFound", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Get() error = %v, want nil", err)
			}
			if tc.wantFound && obj.GetName() != tc.getName {
				t.Errorf("Get() name = %q, want %q", obj.GetName(), tc.getName)
			}
		})
	}
}

func TestClient_CreateOrUpdate(t *testing.T) {
	t.Parallel()

	oc := newConfigMapClient(t)
	ctx := context.Background()

	obj := newConfigMap("target", "default", map[string]string{"k": "v1"})
	ope, err := ctrl.CreateOrUpdate(ctx, oc.api, obj, func() error { return nil })
	if err != nil {
		t.Fatalf("CreateOrUpdate(create) error = %v", err)
	}
	if ope != OperationResultCreated {
		t.Errorf("CreateOrUpdate(create) = %v, want %v", ope, OperationResultCreated)
	}

	// mutate on the second pass: the object exists, so it must be updated
	ope, err = ctrl.CreateOrUpdate(ctx, oc.api, obj, func() error {
		obj.Data["k"] = "v2"
		return nil
	})
	if err != nil {
		t.Fatalf("CreateOrUpdate(update) error = %v", err)
	}
	if ope != OperationResultUpdated {
		t.Errorf("CreateOrUpdate(update) = %v, want %v", ope, OperationResultUpdated)
	}

	got, err := oc.Get(ctx, "target", "default")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Data["k"] != "v2" {
		t.Errorf("data = %q, want %q", got.Data["k"], "v2")
	}
}

func TestClient_Wait(t *testing.T) {
	t.Parallel()

	t.Run("returns the context error when canceled", func(t *testing.T) {
		t.Parallel()

		oc := newConfigMapClient(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		done, err := oc.Wait(ctx, "missing", "default", func(*corev1.ConfigMap) (bool, error) {
			return true, nil
		})
		if done || err == nil {
			t.Errorf("Wait() = (%v, %v), want (false, context error)", done, err)
		}
	})

	t.Run("returns true once eval reports done", func(t *testing.T) {
		t.Parallel()

		oc := newConfigMapClient(t, newConfigMap("ready", "default", nil))
		done, err := oc.Wait(
			context.Background(), "ready", "default",
			func(cm *corev1.ConfigMap) (bool, error) {
				return cm.GetName() == "ready", nil
			},
		)
		if err != nil {
			t.Fatalf("Wait() error = %v", err)
		}
		if !done {
			t.Error("Wait() done = false, want true")
		}
	})
}

func TestClient_List(t *testing.T) {
	t.Parallel()

	lc := newConfigMapClient(t,
		newConfigMap("a", "default", nil),
		newConfigMap("b", "default", nil),
		newConfigMap("c", "other", nil),
	)
	list, err := lc.List(context.Background(), kclient.InNamespace("default"))
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list.Items) != 2 {
		t.Errorf("List() len = %d, want 2", len(list.Items))
	}
}

func TestClient_CreateUpdateDelete(t *testing.T) {
	t.Parallel()

	lc := newConfigMapClient(t)
	ctx := context.Background()

	obj := newConfigMap("target", "default", map[string]string{"k": "v1"})
	if err := lc.Create(ctx, obj); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := lc.Get(ctx, "target", "default")
	if err != nil {
		t.Fatalf("Get() after Create error = %v", err)
	}
	if got.Data["k"] != "v1" {
		t.Errorf("data = %q, want %q", got.Data["k"], "v1")
	}

	got.Data["k"] = "v2"
	if err = lc.Update(ctx, got); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	got, err = lc.Get(ctx, "target", "default")
	if err != nil {
		t.Fatalf("Get() after Update error = %v", err)
	}
	if got.Data["k"] != "v2" {
		t.Errorf("data after update = %q, want %q", got.Data["k"], "v2")
	}

	if err := lc.Delete(ctx, got); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := lc.Get(ctx, "target", "default"); !apierrors.IsNotFound(err) {
		t.Errorf("Get() after Delete error = %v, want NotFound", err)
	}
}

func TestClient_Watch(t *testing.T) {
	t.Parallel()

	lc := newConfigMapClient(t)
	ctx := t.Context()

	w, err := lc.Watch(ctx, kclient.InNamespace("default"))
	if err != nil {
		t.Fatalf("Watch() error = %v", err)
	}
	defer w.Stop()

	if err := lc.Create(ctx, newConfigMap("watched", "default", nil)); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	select {
	case event := <-w.ResultChan():
		if event.Type != watch.Added {
			t.Errorf("event.Type = %v, want ADDED", event.Type)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for watch event")
	}
}

// noWatchClient hides the Watch method of the embedded controller-runtime
// client, mimicking watch-incapable clients such as what mgr.GetClient()
// returns.
type noWatchClient struct {
	kclient.Client
}

func TestClient_WatchUnsupported(t *testing.T) {
	t.Parallel()

	t.Run("Watch returns ErrKubernetesClientWatchNotSupported", func(t *testing.T) {
		t.Parallel()

		lc := NewClient(
			noWatchClient{Client: newFakeAPI(t)},
			new(corev1.ConfigMap), new(corev1.ConfigMapList),
		)
		w, err := lc.Watch(context.Background())
		if !errors.Is(err, errors.ErrKubernetesClientWatchNotSupported) {
			t.Errorf("Watch() error = %v, want ErrKubernetesClientWatchNotSupported", err)
		}
		if w != nil {
			t.Errorf("Watch() = %v, want nil", w)
		}
	})

	t.Run("DeleteAndWait wraps ErrKubernetesClientWatchNotSupported and does not delete", func(t *testing.T) {
		t.Parallel()

		obj := newConfigMap("kept", "default", nil)
		lc := NewClient(
			noWatchClient{Client: newFakeAPI(t, obj)},
			new(corev1.ConfigMap), new(corev1.ConfigMapList),
		)
		err := lc.DeleteAndWait(context.Background(), obj)
		if !errors.Is(err, errors.ErrKubernetesClientWatchNotSupported) {
			t.Errorf("DeleteAndWait() error = %v, want wrapped ErrKubernetesClientWatchNotSupported", err)
		}
		// the delete must not have happened: the watch is established first
		if _, err := lc.Get(context.Background(), "kept", "default"); err != nil {
			t.Errorf("Get() after failed DeleteAndWait error = %v, want nil (object kept)", err)
		}
	})
}

// fakeClient adapts a fake controller-runtime client to the client.Client
// interface NewClientOf now expects (client.Client embeds k8s.Client
// directly, so no separate Raw() bridge is needed). GetClientSet and
// GetRESTConfig are never exercised by NewClientOf's Get/List/Watch calls, so
// they return zero values.
type fakeClient struct {
	kclient.WithWatch
}

func (fakeClient) GetClientSet() kubernetes.Interface { return nil }

func (fakeClient) GetRESTConfig() *rest.Config { return nil }

func TestNewClientOf(t *testing.T) {
	t.Parallel()

	raw := newFakeAPI(t, newConfigMap("bridged", "default", map[string]string{"k": "v"}))

	lc := NewClientOf(fakeClient{WithWatch: raw}, new(corev1.ConfigMap), new(corev1.ConfigMapList))

	got, err := lc.Get(context.Background(), "bridged", "default")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Data["k"] != "v" {
		t.Errorf("data = %q, want %q", got.Data["k"], "v")
	}

	list, err := lc.List(context.Background(), kclient.InNamespace("default"))
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list.Items) != 1 {
		t.Errorf("List() len = %d, want 1", len(list.Items))
	}

	// fakeClient embeds a kclient.WithWatch, so the wrapped client keeps the
	// watch capability.
	w, err := lc.Watch(context.Background(), kclient.InNamespace("default"))
	if err != nil {
		t.Fatalf("Watch() error = %v", err)
	}
	w.Stop()
}

// NOT IMPLEMENTED BELOW
//
// func Test_cloneSeed(t *testing.T) {
// 	type args struct {
// 		seed S
// 	}
// 	type want struct {
// 		want S
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, S) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got S) error {
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
// 		           seed:nil,
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
// 		           seed:nil,
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
// 			got := cloneSeed(test.args.seed)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestNewClient(t *testing.T) {
// 	type args struct {
// 		api   k8s.Client
// 		tSeed T
// 		lSeed L
// 	}
// 	type want struct {
// 		want *Client[T, L]
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *Client[T, L]) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *Client[T, L]) error {
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
// 		           api:nil,
// 		           tSeed:nil,
// 		           lSeed:nil,
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
// 		           api:nil,
// 		           tSeed:nil,
// 		           lSeed:nil,
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
// 			got := NewClient(test.args.api, test.args.tSeed, test.args.lSeed)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClient_fresh(t *testing.T) {
// 	type want struct {
// 		want T
// 	}
// 	type test struct {
// 		name       string
// 		c          *Client[T, L]
// 		want       want
// 		checkFunc  func(want, T) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got T) error {
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
// 			c := &Client[T, L]{}
//
// 			got := c.fresh()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClient_freshList(t *testing.T) {
// 	type want struct {
// 		want L
// 	}
// 	type test struct {
// 		name       string
// 		c          *Client[T, L]
// 		want       want
// 		checkFunc  func(want, L) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got L) error {
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
// 			c := &Client[T, L]{}
//
// 			got := c.freshList()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClient_UpdateStatus(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		obj T
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		c          *Client[T, L]
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
// 		           obj:nil,
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
// 		           obj:nil,
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
// 			c := &Client[T, L]{}
//
// 			err := c.UpdateStatus(test.args.ctx, test.args.obj)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClient_Create(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		obj  T
// 		opts []kclient.CreateOption
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		c          *Client[T, L]
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
// 		           obj:nil,
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
// 		           ctx:nil,
// 		           obj:nil,
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
// 			c := &Client[T, L]{}
//
// 			err := c.Create(test.args.ctx, test.args.obj, test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClient_Update(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		obj  T
// 		opts []kclient.UpdateOption
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		c          *Client[T, L]
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
// 		           obj:nil,
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
// 		           ctx:nil,
// 		           obj:nil,
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
// 			c := &Client[T, L]{}
//
// 			err := c.Update(test.args.ctx, test.args.obj, test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClient_Delete(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		obj  T
// 		opts []kclient.DeleteOption
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		c          *Client[T, L]
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
// 		           obj:nil,
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
// 		           ctx:nil,
// 		           obj:nil,
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
// 			c := &Client[T, L]{}
//
// 			err := c.Delete(test.args.ctx, test.args.obj, test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_waitLoop(t *testing.T) {
// 	type args struct {
// 		ctx       context.Context
// 		onTimeout error
// 		step      func(context.Context) (done bool, err error)
// 	}
// 	type want struct {
// 		want bool
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool, err error) error {
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
// 		           onTimeout:nil,
// 		           step:nil,
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
// 		           onTimeout:nil,
// 		           step:nil,
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
// 			got, err := waitLoop(test.args.ctx, test.args.onTimeout, test.args.step)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClient_DeleteAndWait(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		obj  T
// 		opts []ListOption
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		c          *Client[T, L]
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
// 		           obj:nil,
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
// 		           ctx:nil,
// 		           obj:nil,
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
// 			c := &Client[T, L]{}
//
// 			err := c.DeleteAndWait(test.args.ctx, test.args.obj, test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
