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
package target

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	mirrv1 "github.com/vdaas/vald/internal/k8s/vald/mirror/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// mockManager is a minimal k8s.Manager stub: the reconciler built by New()
// only ever exercises GetClient and GetScheme.
type mockManager struct {
	k8s.Manager

	client k8s.Client
	scheme *k8s.Scheme
}

func (m *mockManager) GetClient() k8s.Client  { return m.client }
func (m *mockManager) GetScheme() *k8s.Scheme { return m.scheme }

// GetConfig satisfies client.NewFromManager's k8s.Manager requirement: the
// embedded k8s.Manager is nil, so without this override GetConfig would
// panic instead of returning a config client.NewFromManager can build a
// (non-dialing) clientset from.
func (*mockManager) GetConfig() *rest.Config {
	return &rest.Config{}
}

// errListClient fails every List call with the configured error. Every other
// method panics if called, which no path exercised by these tests reaches.
type errListClient struct {
	k8s.Client

	err error
}

func (c *errListClient) List(context.Context, k8s.ObjectList, ...k8s.ListOption) error {
	return c.err
}

const (
	testRoleLabelKey  = "role"
	testColocationAZA = "az-a"
	testColocationAZB = "az-b"
	testHostA         = "host-a"
	testHostB         = "host-b"
)

func newScheme(t *testing.T) *k8s.Scheme {
	t.Helper()
	scheme := k8s.NewScheme()
	if err := mirrv1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add ValdMirrorTarget scheme: %v", err)
	}
	return scheme
}

func newFakeClient(t *testing.T, objs ...k8s.Object) k8s.Client {
	t.Helper()
	return fake.NewClientBuilder().WithScheme(newScheme(t)).WithObjects(objs...).Build()
}

func newMirrorTarget(
	name, namespace, host string, port int, colocation string, phase MirrorTargetPhase,
) *mirrv1.ValdMirrorTarget {
	return &mirrv1.ValdMirrorTarget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{testRoleLabelKey: "mirror"},
		},
		Spec: mirrv1.MirrorTargetSpec{
			Colocation: colocation,
			Target:     mirrv1.MirrorTarget{Host: host, Port: port},
		},
		Status: mirrv1.MirrorTargetStatus{Phase: phase},
	}
}

func TestNew_ControllerName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opt  string
	}{
		{name: "reflects the configured controller name", opt: "mirror target watcher"},
		{name: "reflects a different controller name", opt: "another watcher"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			watcher, err := New(
				WithControllerName(tc.opt),
				WithOnReconcileFunc(func(context.Context, map[string]Target) {}),
			)
			if err != nil {
				t.Fatalf("New() error = %v, want nil", err)
			}
			if got := watcher.GetName(); got != tc.opt {
				t.Errorf("GetName() = %q, want %q", got, tc.opt)
			}
		})
	}
}

// TestNew_InvalidOptionsAreNonFatal pins down New()'s existing severity
// handling: options that fail validation with a plain (non-critical) error
// are logged and skipped rather than aborting construction.
func TestNew_InvalidOptionsAreNonFatal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts []Option
	}{
		{name: "empty controller name", opts: []Option{WithControllerName("")}},
		{name: "nil onError func", opts: []Option{WithOnErrorFunc(nil)}},
		{name: "nil onReconcile func", opts: []Option{WithOnReconcileFunc(nil)}},
		{name: "no namespaces", opts: []Option{WithNamespaces()}},
		{name: "empty namespace", opts: []Option{WithNamespaces("")}},
		{name: "empty labels", opts: []Option{WithLabels(nil)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			watcher, err := New(tc.opts...)
			if err != nil {
				t.Fatalf("New() error = %v, want nil (invalid non-critical options must be warned and skipped)", err)
			}
			if watcher == nil {
				t.Fatal("New() watcher = nil, want non-nil")
			}
		})
	}
}

// TestNew_Reconcile drives the reconciler returned by New() end-to-end
// through a fake controller-runtime client, characterizing the List ->
// map[string]Target conversion and the namespace/label filtering options.
func TestNew_Reconcile(t *testing.T) {
	t.Parallel()

	type test struct {
		client  func(t *testing.T) k8s.Client
		targets map[string]Target
		name    string
		opts    []Option
	}

	tests := []test{
		{
			name: "converts the listed targets into a name-keyed map",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t,
					newMirrorTarget("t1", "default", testHostA, 8081, testColocationAZA, MirrorTargetPhaseConnected),
					newMirrorTarget("t2", "default", testHostB, 8082, testColocationAZB, MirrorTargetPhasePending),
				)
			},
			targets: map[string]Target{
				"t1": {Colocation: testColocationAZA, Host: testHostA, Port: 8081, Phase: MirrorTargetPhaseConnected},
				"t2": {Colocation: testColocationAZB, Host: testHostB, Port: 8082, Phase: MirrorTargetPhasePending},
			},
		},
		{
			name: "filters targets by namespace",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t,
					newMirrorTarget("t1", "ns-a", testHostA, 8081, testColocationAZA, MirrorTargetPhaseConnected),
					newMirrorTarget("t2", "ns-b", testHostB, 8082, testColocationAZB, MirrorTargetPhasePending),
				)
			},
			opts: []Option{WithNamespaces("ns-a")},
			targets: map[string]Target{
				"t1": {Colocation: testColocationAZA, Host: testHostA, Port: 8081, Phase: MirrorTargetPhaseConnected},
			},
		},
		{
			name: "filters targets by label",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				a := newMirrorTarget("t1", "default", testHostA, 8081, testColocationAZA, MirrorTargetPhaseConnected)
				a.Labels = map[string]string{testRoleLabelKey: "primary"}
				b := newMirrorTarget("t2", "default", testHostB, 8082, testColocationAZB, MirrorTargetPhasePending)
				b.Labels = map[string]string{testRoleLabelKey: "replica"}
				return newFakeClient(t, a, b)
			},
			opts: []Option{WithLabels(map[string]string{testRoleLabelKey: "primary"})},
			targets: map[string]Target{
				"t1": {Colocation: testColocationAZA, Host: testHostA, Port: 8081, Phase: MirrorTargetPhaseConnected},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var errored int
			var got map[string]Target
			watcher, err := New(append([]Option{
				WithControllerName("mirror target watcher"),
				WithOnReconcileFunc(func(_ context.Context, m map[string]Target) { got = m }),
				WithOnErrorFunc(func(error) { errored++ }),
			}, tc.opts...)...)
			if err != nil {
				t.Fatalf("New() error = %v, want nil", err)
			}

			mgr := &mockManager{client: tc.client(t), scheme: k8s.NewScheme()}
			rec := watcher.NewReconciler(context.Background(), mgr)

			res, err := rec.Reconcile(context.Background(), k8s.Request{})
			if err != nil {
				t.Fatalf("Reconcile() error = %v, want nil", err)
			}
			if res.RequeueAfter != 0 {
				t.Errorf("Reconcile() RequeueAfter = %v, want 0", res.RequeueAfter)
			}
			if !reflect.DeepEqual(got, tc.targets) {
				t.Errorf("onReconcile map = %+v, want %+v", got, tc.targets)
			}
			if errored != 0 {
				t.Errorf("onError fired %d times, want 0", errored)
			}
		})
	}
}

// TestNew_Reconcile_NilOnReconcile pins down that Reconcile does not panic
// when no WithOnReconcileFunc callback was ever configured (regression test:
// the callback used to be invoked unconditionally).
func TestNew_Reconcile_NilOnReconcile(t *testing.T) {
	t.Parallel()

	watcher, err := New(WithControllerName("mirror target watcher"))
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}

	mgr := &mockManager{
		client: newFakeClient(t, newMirrorTarget("t1", "default", testHostA, 8081, testColocationAZA, MirrorTargetPhaseConnected)),
		scheme: k8s.NewScheme(),
	}
	rec := watcher.NewReconciler(context.Background(), mgr)

	if _, err := rec.Reconcile(context.Background(), k8s.Request{}); err != nil {
		t.Fatalf("Reconcile() error = %v, want nil", err)
	}
}

// TestNew_Reconcile_ListErrors characterizes the requeue duration and
// onError notification for List failures. The raw returned error is
// intentionally not asserted: the hand-rolled reconciler this delegates to
// used to return the error alongside RequeueAfter (which controller-runtime
// then discards because Result is ignored on a non-nil error), while the
// generic list reconciler consumes the error itself to make the configured
// RequeueAfter actually effective. Both behave identically from
// controller-runtime's perspective; only RequeueAfter and onError firing are
// part of the observable contract pinned down here.
func TestNew_Reconcile_ListErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err     error
		name    string
		requeue time.Duration
	}{
		{
			name:    "generic list error requeues after the invalid-error duration and notifies onError",
			err:     errors.New("boom"),
			requeue: 100 * time.Millisecond,
		},
		{
			name:    "NotFound list error requeues after the not-found duration and notifies onError",
			err:     apierrors.NewNotFound(schema.GroupResource{Resource: "valdmirrortargets"}, "missing"),
			requeue: time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var errored int
			watcher, err := New(
				WithControllerName("mirror target watcher"),
				WithOnReconcileFunc(func(context.Context, map[string]Target) {}),
				WithOnErrorFunc(func(error) { errored++ }),
			)
			if err != nil {
				t.Fatalf("New() error = %v, want nil", err)
			}

			mgr := &mockManager{
				client: &errListClient{err: tc.err},
				scheme: k8s.NewScheme(),
			}
			rec := watcher.NewReconciler(context.Background(), mgr)

			res, _ := rec.Reconcile(context.Background(), k8s.Request{})
			if res.RequeueAfter != tc.requeue {
				t.Errorf("Reconcile() RequeueAfter = %v, want %v", res.RequeueAfter, tc.requeue)
			}
			if errored == 0 {
				t.Error("onError callback was not invoked")
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
// 		want MirrorTargetWatcher
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, MirrorTargetWatcher, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got MirrorTargetWatcher, err error) error {
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
