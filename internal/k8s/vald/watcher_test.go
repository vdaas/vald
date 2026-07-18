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

package vald

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	mirrv1 "github.com/vdaas/vald/internal/k8s/vald/mirror/api/v1"
	"github.com/vdaas/vald/internal/test/goleak"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

// mockManager is a minimal k8s.Manager stub: the reconciler built by
// NewListWatcher only ever exercises GetClient and GetScheme.
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

// recordingClient wraps a k8s.Client and records the k8s.ListOptions applied
// on every List call, so tests can assert on what the reconciler actually
// asked the client for (namespace/label scoping) instead of only on the
// returned items.
type recordingClient struct {
	k8s.Client

	lastOpts k8s.ListOptions
}

func (c *recordingClient) List(
	ctx context.Context, list k8s.ObjectList, opts ...k8s.ListOption,
) error {
	c.lastOpts = k8s.ListOptions{}
	for _, o := range opts {
		o.ApplyToList(&c.lastOpts)
	}
	return c.Client.List(ctx, list, opts...)
}

// errListClient fails every List call with the configured error.
type errListClient struct {
	k8s.Client

	err error
}

func (c *errListClient) List(context.Context, k8s.ObjectList, ...k8s.ListOption) error {
	return c.err
}

const (
	testRoleLabelKey = "role"
	testHostA        = "host-a"
	testHostB        = "host-b"
)

func newMirrorTarget(name, namespace, host string) *mirrv1.ValdMirrorTarget {
	return &mirrv1.ValdMirrorTarget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{testRoleLabelKey: "mirror"},
		},
		Spec: mirrv1.MirrorTargetSpec{
			Target: mirrv1.MirrorTarget{Host: host, Port: 8081},
		},
	}
}

func newFakeClient(t *testing.T, objs ...kclient.Object) k8s.Client {
	t.Helper()
	scheme := k8s.NewScheme()
	if err := mirrv1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add ValdMirrorTarget scheme: %v", err)
	}
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

// hostOf is the convert function used by the tests: it maps every listed
// object to its spec host, so the callback map is easy to assert on.
func hostOf(m *mirrv1.ValdMirrorTarget) string { return m.Spec.Target.Host }

// TestNewListWatcher_OptionErrorPolicy pins down the unified option error
// policy shared by every watcher package: SkipNonCriticalOptionError warns
// and continues on any non-critical option error (including the validation
// errors the shared options emit) and aborts construction only on
// errors.ErrCriticalOption.
func TestNewListWatcher_OptionErrorPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		optErr  error
		name    string
		wantErr bool
	}{
		{
			name:   "warns and continues on a plain option error",
			optErr: errors.New("boom"),
		},
		{
			name:   "warns and continues on an invalid option error",
			optErr: errors.NewErrInvalidOption("namespaces", nil),
		},
		{
			name:    "aborts construction on a critical option error",
			optErr:  errors.NewErrCriticalOption("name", "value"),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			failing := WatcherOption[string](func(*WatcherConfig[string]) error { return tc.optErr })
			watcher, err := NewListWatcher(
				mirrv1.AddToScheme,
				hostOf,
				SkipNonCriticalOptionError,
				WithWatcherControllerName[string]("list watcher"),
				failing,
			)
			if tc.wantErr {
				if err == nil {
					t.Fatal("NewListWatcher() error = nil, want a wrapped option failure")
				}
				if !errors.Is(err, tc.optErr) {
					t.Errorf("NewListWatcher() error = %v, want it to wrap %v", err, tc.optErr)
				}
				if watcher != nil {
					t.Errorf("NewListWatcher() watcher = %v, want nil on abort", watcher)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewListWatcher() error = %v, want nil (non-critical option errors must be skipped)", err)
			}
			if watcher == nil {
				t.Fatal("NewListWatcher() watcher = nil, want non-nil")
			}
			if got, want := watcher.GetName(), "list watcher"; got != want {
				t.Errorf("GetName() = %q, want %q (options after the failing one must still apply)", got, want)
			}
		})
	}
}

// TestWatcherOptions_Validation pins down the per-option validation shared by
// every watcher package: empty or nil values are rejected with
// errors.ErrInvalidOption and leave the config untouched (New then warns and
// skips them under SkipNonCriticalOptionError).
func TestWatcherOptions_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		opt  WatcherOption[string]
		name string
	}{
		{name: "empty controller name", opt: WithWatcherControllerName[string]("")},
		{name: "no namespaces", opt: WithWatcherNamespaces[string]()},
		{name: "empty namespace element", opt: WithWatcherNamespaces[string]("ns-a", "")},
		{name: "nil labels", opt: WithWatcherLabels[string](nil)},
		{name: "empty labels", opt: WithWatcherLabels[string](map[string]string{})},
		{name: "nil onError func", opt: WithWatcherOnErrorFunc[string](nil)},
		{name: "nil onReconcile func", opt: WithWatcherOnReconcileFunc[string](nil)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := new(WatcherConfig[string])
			err := tc.opt(c)
			if err == nil {
				t.Fatal("option error = nil, want errors.ErrInvalidOption")
			}
			e := &errors.ErrInvalidOption{}
			if !errors.As(err, &e) {
				t.Errorf("option error = %v, want errors.ErrInvalidOption", err)
			}
			if !reflect.DeepEqual(c, new(WatcherConfig[string])) {
				t.Errorf("config = %+v, want untouched zero value on invalid option", c)
			}
		})
	}
}

// TestNewListWatcher_Reconcile drives the reconciler end-to-end through a
// fake controller-runtime client, characterizing the per-item conversion into
// a name-keyed map and the namespace/label scoping of the List call.
func TestNewListWatcher_Reconcile(t *testing.T) {
	t.Parallel()

	type want struct {
		hosts     map[string]string
		namespace string
	}
	tests := []struct {
		client func(t *testing.T) k8s.Client
		want   want
		name   string
		opts   []WatcherOption[string]
	}{
		{
			name: "converts every listed item via convert, keyed by object name",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t,
					newMirrorTarget("t1", "default", testHostA),
					newMirrorTarget("t2", "default", testHostB),
				)
			},
			want: want{hosts: map[string]string{"t1": testHostA, "t2": testHostB}},
		},
		{
			name: "a single namespace scopes the List call",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t,
					newMirrorTarget("t1", "ns-a", testHostA),
					newMirrorTarget("t2", "ns-b", testHostB),
				)
			},
			opts: []WatcherOption[string]{WithWatcherNamespaces[string]("ns-a")},
			want: want{hosts: map[string]string{"t1": testHostA}, namespace: "ns-a"},
		},
		{
			name: "multiple namespaces keep the known last-wins limitation",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t,
					newMirrorTarget("t1", "ns-a", testHostA),
					newMirrorTarget("t2", "ns-b", testHostB),
				)
			},
			opts: []WatcherOption[string]{WithWatcherNamespaces[string]("ns-a", "ns-b")},
			want: want{hosts: map[string]string{"t2": testHostB}, namespace: "ns-b"},
		},
		{
			name: "labels restrict the List call",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				a := newMirrorTarget("t1", "default", testHostA)
				a.Labels = map[string]string{testRoleLabelKey: "primary"}
				b := newMirrorTarget("t2", "default", testHostB)
				b.Labels = map[string]string{testRoleLabelKey: "replica"}
				return newFakeClient(t, a, b)
			},
			opts: []WatcherOption[string]{WithWatcherLabels[string](map[string]string{testRoleLabelKey: "primary"})},
			want: want{hosts: map[string]string{"t1": testHostA}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got map[string]string
			watcher, err := NewListWatcher(
				mirrv1.AddToScheme,
				hostOf,
				SkipNonCriticalOptionError,
				append([]WatcherOption[string]{
					WithWatcherControllerName[string]("list watcher"),
					WithWatcherOnReconcileFunc(func(_ context.Context, m map[string]string) { got = m }),
				}, tc.opts...)...,
			)
			if err != nil {
				t.Fatalf("NewListWatcher() error = %v", err)
			}

			client := &recordingClient{Client: tc.client(t)}
			rec := watcher.NewReconciler(context.Background(), &mockManager{client: client, scheme: k8s.NewScheme()})

			if _, err := rec.Reconcile(context.Background(), k8s.Request{}); err != nil {
				t.Fatalf("Reconcile() error = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want.hosts) {
				t.Errorf("onReconcile map = %v, want %v", got, tc.want.hosts)
			}
			if tc.want.namespace != "" && client.lastOpts.Namespace != tc.want.namespace {
				t.Errorf("List() ListOptions.Namespace = %q, want %q", client.lastOpts.Namespace, tc.want.namespace)
			}
		})
	}
}

// TestNewListWatcher_NilOnReconcile verifies the nil callback guard: without
// a configured OnReconcile, Reconcile neither panics nor spends work on
// converting the listed items.
func TestNewListWatcher_NilOnReconcile(t *testing.T) {
	t.Parallel()

	converted := 0
	watcher, err := NewListWatcher(
		mirrv1.AddToScheme,
		func(m *mirrv1.ValdMirrorTarget) string {
			converted++
			return m.Spec.Target.Host
		},
		SkipNonCriticalOptionError,
		WithWatcherControllerName[string]("list watcher"),
	)
	if err != nil {
		t.Fatalf("NewListWatcher() error = %v", err)
	}

	mgr := &mockManager{
		client: newFakeClient(t, newMirrorTarget("t1", "default", testHostA)),
		scheme: k8s.NewScheme(),
	}
	rec := watcher.NewReconciler(context.Background(), mgr)

	if _, err := rec.Reconcile(context.Background(), k8s.Request{}); err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}
	if converted != 0 {
		t.Errorf("convert ran %d times without an OnReconcile callback, want 0", converted)
	}
}

// TestNewListWatcher_ListErrorNotifiesOnError verifies the OnError callback
// wiring: a failing List notifies the callback and requeues instead of
// returning the error.
func TestNewListWatcher_ListErrorNotifiesOnError(t *testing.T) {
	t.Parallel()

	var errored int
	watcher, err := NewListWatcher(
		mirrv1.AddToScheme,
		hostOf,
		SkipNonCriticalOptionError,
		WithWatcherControllerName[string]("list watcher"),
		WithWatcherOnErrorFunc[string](func(error) { errored++ }),
	)
	if err != nil {
		t.Fatalf("NewListWatcher() error = %v", err)
	}

	mgr := &mockManager{
		client: &errListClient{err: errors.New("boom")},
		scheme: k8s.NewScheme(),
	}
	rec := watcher.NewReconciler(context.Background(), mgr)

	res, err := rec.Reconcile(context.Background(), k8s.Request{})
	if err != nil {
		t.Fatalf("Reconcile() error = %v, want nil (list errors are consumed into RequeueAfter)", err)
	}
	if res.RequeueAfter == 0 {
		t.Error("Reconcile() RequeueAfter = 0, want a retry interval on list error")
	}
	if errored != 1 {
		t.Errorf("onError fired %d times, want 1", errored)
	}
}

// TestNewListWatcher_RegistersScheme verifies NewReconciler registers the
// given AddToScheme on the manager's scheme instead of falling back to the
// client-go native scheme.
func TestNewListWatcher_RegistersScheme(t *testing.T) {
	t.Parallel()

	watcher, err := NewListWatcher(
		mirrv1.AddToScheme,
		hostOf,
		SkipNonCriticalOptionError,
		WithWatcherControllerName[string]("list watcher"),
	)
	if err != nil {
		t.Fatalf("NewListWatcher() error = %v", err)
	}

	mgr := &mockManager{client: newFakeClient(t), scheme: k8s.NewScheme()}
	watcher.NewReconciler(context.Background(), mgr)

	if !mgr.scheme.Recognizes(mirrv1.GroupVersion.WithKind("ValdMirrorTarget")) {
		t.Error("NewReconciler() did not register the ValdMirrorTarget scheme")
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestWithWatcherControllerName(t *testing.T) {
// 	type args struct {
// 		name string
// 	}
// 	type want struct {
// 		want WatcherOption[V]
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, WatcherOption[V]) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got WatcherOption[V]) error {
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
// 		           name:"",
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
// 		           name:"",
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
// 			got := WithWatcherControllerName(test.args.name)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithWatcherNamespaces(t *testing.T) {
// 	type args struct {
// 		nss []string
// 	}
// 	type want struct {
// 		want WatcherOption[V]
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, WatcherOption[V]) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got WatcherOption[V]) error {
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
// 		           nss:nil,
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
// 		           nss:nil,
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
// 			got := WithWatcherNamespaces(test.args.nss...)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithWatcherLabels(t *testing.T) {
// 	type args struct {
// 		labels map[string]string
// 	}
// 	type want struct {
// 		want WatcherOption[V]
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, WatcherOption[V]) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got WatcherOption[V]) error {
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
// 		           labels:nil,
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
// 		           labels:nil,
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
// 			got := WithWatcherLabels(test.args.labels)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithWatcherOnErrorFunc(t *testing.T) {
// 	type args struct {
// 		f func(err error)
// 	}
// 	type want struct {
// 		want WatcherOption[V]
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, WatcherOption[V]) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got WatcherOption[V]) error {
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
// 		           f:nil,
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
// 		           f:nil,
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
// 			got := WithWatcherOnErrorFunc(test.args.f)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWithWatcherOnReconcileFunc(t *testing.T) {
// 	type args struct {
// 		f func(ctx context.Context, resources map[string]V)
// 	}
// 	type want struct {
// 		want WatcherOption[V]
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, WatcherOption[V]) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got WatcherOption[V]) error {
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
// 		           f:nil,
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
// 		           f:nil,
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
// 			got := WithWatcherOnReconcileFunc(test.args.f)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestSkipNonCriticalOptionError(t *testing.T) {
// 	type args struct {
// 		err error
// 		opt any
// 	}
// 	type want struct {
// 		wantAbort bool
// 		err       error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotAbort bool, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotAbort, w.wantAbort) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotAbort, w.wantAbort)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           err:nil,
// 		           opt:nil,
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
// 		           err:nil,
// 		           opt:nil,
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
// 			gotAbort, err := SkipNonCriticalOptionError(test.args.err, test.args.opt)
// 			if err := checkFunc(test.want, gotAbort, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestNewListWatcher(t *testing.T) {
// 	type args struct {
// 		addToScheme func(s *k8s.Scheme) error
// 		convert     func(obj PT) V
// 		policy      WatcherOptionErrorPolicy
// 		opts        []WatcherOption[V]
// 	}
// 	type want struct {
// 		want k8s.ResourceController
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, k8s.ResourceController, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.ResourceController, err error) error {
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
// 		           addToScheme:nil,
// 		           convert:nil,
// 		           policy:nil,
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
// 		           addToScheme:nil,
// 		           convert:nil,
// 		           policy:nil,
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
// 			got, err := NewListWatcher(test.args.addToScheme, test.args.convert, test.args.policy, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
