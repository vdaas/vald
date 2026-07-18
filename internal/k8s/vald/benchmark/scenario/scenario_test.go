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

package scenario

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/test/goleak"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

// testControllerName is the controller name used across the tests, matching
// the name pkg/operator/benchmark/service registers in production.
const testControllerName = "benchmark scenario resource"

// managerMock is a minimal k8s.Manager stub: New()'s reconciler only
// exercises GetClient and GetScheme.
type managerMock struct {
	k8s.Manager

	client k8s.Client
	scheme *k8s.Scheme
}

func (m *managerMock) GetClient() k8s.Client { return m.client }

func (m *managerMock) GetScheme() *k8s.Scheme { return m.scheme }

// GetConfig satisfies client.NewFromManager's k8s.Manager requirement: the
// embedded k8s.Manager is nil, so without this override GetConfig would
// panic instead of returning a config client.NewFromManager can build a
// (non-dialing) clientset from.
func (*managerMock) GetConfig() *rest.Config {
	return &rest.Config{}
}

// recordingClient wraps a k8s.Client and records the k8s.ListOptions applied
// on every List call, so tests can assert on what the reconciler actually
// asked the client for (namespace scoping in particular) instead of only on
// the returned items.
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

func newScenario(name, namespace string) *v1.ValdBenchmarkScenario {
	return &v1.ValdBenchmarkScenario{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func newFakeClient(t *testing.T, objs ...kclient.Object) k8s.Client {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add scenario scheme: %v", err)
	}
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

// newReconciler drives New(opts...) all the way down to a k8s.Reconciler
// backed by client, exercising the exact same NewListReconciler wiring New()
// uses in production instead of re-implementing the reconcile loop in the
// test.
func newReconciler(t *testing.T, client k8s.Client, opts ...Option) k8s.Reconciler {
	t.Helper()
	watcher, err := New(opts...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	mgr := &managerMock{client: client, scheme: runtime.NewScheme()}
	return watcher.NewReconciler(context.Background(), mgr)
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a BenchmarkWatcher with the configured name", func(t *testing.T) {
		t.Parallel()

		watcher, err := New(WithControllerName(testControllerName))
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		if watcher == nil {
			t.Fatal("New() watcher = nil")
		}
		if got, want := watcher.GetName(), testControllerName; got != want {
			t.Errorf("GetName() = %q, want %q", got, want)
		}
	})

	t.Run("skips non-critical option errors and construction succeeds", func(t *testing.T) {
		t.Parallel()

		watcher, err := New(
			func(*config) error { return errors.New("option boom") },
			WithControllerName(testControllerName),
		)
		if err != nil {
			t.Fatalf("New() error = %v, want nil (non-critical option errors must be skipped)", err)
		}
		if watcher == nil {
			t.Fatal("New() watcher = nil, want non-nil")
		}
		if got := watcher.GetName(); got != testControllerName {
			t.Errorf("GetName() = %q, want %q (options after the failing one must still apply)", got, testControllerName)
		}
	})

	t.Run("aborts construction on a critical option error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.NewErrCriticalOption("name", "value")
		watcher, err := New(func(*config) error { return wantErr })
		if err == nil {
			t.Fatal("New() error = nil, want a wrapped critical option failure")
		}
		if !errors.Is(err, wantErr) {
			t.Errorf("New() error = %v, want it to wrap %v", err, wantErr)
		}
		if watcher != nil {
			t.Errorf("New() watcher = %v, want nil on critical option failure", watcher)
		}
	})
}

// TestNew_InvalidOptionsAreNonFatal pins down the validation unified with the
// mirror target watcher: options rejecting empty or nil values fail with
// errors.ErrInvalidOption, which New logs and skips instead of aborting.
func TestNew_InvalidOptionsAreNonFatal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts []Option
	}{
		{name: "empty controller name", opts: []Option{WithControllerName("")}},
		{name: "no namespaces", opts: []Option{WithNamespaces()}},
		{name: "empty namespace", opts: []Option{WithNamespaces("")}},
		{name: "nil onError func", opts: []Option{WithOnErrorFunc(nil)}},
		{name: "nil onReconcile func", opts: []Option{WithOnReconcileFunc(nil)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := new(config)
			e := &errors.ErrInvalidOption{}
			if err := tc.opts[0](c); !errors.As(err, &e) {
				t.Errorf("option error = %v, want errors.ErrInvalidOption", err)
			}

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

// TestNew_ProductionCallPattern mirrors the exact option usage of
// pkg/operator/benchmark/service initCtrl to guarantee the severity policy
// change keeps the production construction path succeeding.
func TestNew_ProductionCallPattern(t *testing.T) {
	t.Parallel()

	watcher, err := New(
		WithControllerName(testControllerName),
		WithNamespaces("default"),
		WithOnErrorFunc(func(err error) {}),
		WithOnReconcileFunc(func(context.Context, map[string]v1.ValdBenchmarkScenario) {}),
	)
	if err != nil {
		t.Fatalf("New() error = %v, want nil for the production option pattern", err)
	}
	if watcher == nil {
		t.Fatal("New() watcher = nil, want non-nil")
	}
	if got := watcher.GetName(); got != testControllerName {
		t.Errorf("GetName() = %q, want %q", got, testControllerName)
	}
}

func TestOnReconcileFunc(t *testing.T) {
	t.Parallel()

	client := newFakeClient(t, newScenario("scenario-a", "ns-a"))

	var got map[string]v1.ValdBenchmarkScenario
	rec := newReconciler(t, client,
		WithControllerName(testControllerName),
		WithOnReconcileFunc(func(_ context.Context, scenarioList map[string]v1.ValdBenchmarkScenario) {
			got = scenarioList
		}),
	)

	if _, err := rec.Reconcile(context.Background(), k8s.Request{}); err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}

	if got == nil {
		t.Fatal("OnReconcileFunc callback was not invoked")
	}
	if _, ok := got["scenario-a"]; !ok {
		t.Errorf("callback map = %v, want key %q", got, "scenario-a")
	}
	if len(got) != 1 {
		t.Errorf("callback map length = %d, want 1", len(got))
	}
}

func TestOnReconcileFunc_NilCallbackDoesNotPanic(t *testing.T) {
	t.Parallel()

	client := newFakeClient(t, newScenario("scenario-a", "ns-a"))
	rec := newReconciler(t, client, WithControllerName(testControllerName))

	if _, err := rec.Reconcile(context.Background(), k8s.Request{}); err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}
}

func TestOnErrorFunc(t *testing.T) {
	t.Parallel()

	var errored int
	rec := newReconciler(t, &recordingClient{Client: erroringClient{}},
		WithControllerName(testControllerName),
		WithOnErrorFunc(func(error) { errored++ }),
	)

	if _, err := rec.Reconcile(context.Background(), k8s.Request{}); err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}
	if errored != 1 {
		t.Errorf("onError fired %d times, want 1", errored)
	}
}

// erroringClient fails every List call, used to exercise WithOnErrorFunc.
type erroringClient struct {
	k8s.Client
}

func (erroringClient) List(context.Context, k8s.ObjectList, ...k8s.ListOption) error {
	return errors.New("boom")
}

// TestWithNamespaces_ScopesListCall pins down the WithNamespaces bug fix:
// before the fix, WithNamespaces only recorded r.namespaces and never turned
// it into a client.ListOption, so Reconcile listed every namespace
// regardless of what was configured. This test fails (RED) against that
// behavior — both ns-a and ns-b scenarios would show up in the callback
// and/or the recorded ListOptions.Namespace would be empty — and passes
// (GREEN) once WithNamespaces is reflected into the List call.
func TestWithNamespaces_ScopesListCall(t *testing.T) {
	t.Parallel()

	client := &recordingClient{
		Client: newFakeClient(t, newScenario("scenario-a", "ns-a"), newScenario("scenario-b", "ns-b")),
	}

	var got map[string]v1.ValdBenchmarkScenario
	rec := newReconciler(t, client,
		WithControllerName(testControllerName),
		WithNamespaces("ns-a"),
		WithOnReconcileFunc(func(_ context.Context, scenarioList map[string]v1.ValdBenchmarkScenario) {
			got = scenarioList
		}),
	)

	if _, err := rec.Reconcile(context.Background(), k8s.Request{}); err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}

	if got := client.lastOpts.Namespace; got != "ns-a" {
		t.Errorf("List() was called with ListOptions.Namespace = %q, want %q (WithNamespaces was not applied to the List call)", got, "ns-a")
	}

	if len(got) != 1 {
		t.Fatalf("callback map = %v, want exactly 1 item scoped to ns-a", got)
	}
	if _, ok := got["scenario-a"]; !ok {
		t.Errorf("callback map = %v, want key %q", got, "scenario-a")
	}
	if _, ok := got["scenario-b"]; ok {
		t.Errorf("callback map = %v, want ns-b scenario to be filtered out", got)
	}
}

// TestWithNamespaces_Unset verifies the converse: when WithNamespaces is not
// given, every namespace is listed (the pre-existing, still-desired
// behavior for the zero-value case).
func TestWithNamespaces_Unset(t *testing.T) {
	t.Parallel()

	client := newFakeClient(t, newScenario("scenario-a", "ns-a"), newScenario("scenario-b", "ns-b"))

	var got map[string]v1.ValdBenchmarkScenario
	rec := newReconciler(t, client,
		WithControllerName(testControllerName),
		WithOnReconcileFunc(func(_ context.Context, scenarioList map[string]v1.ValdBenchmarkScenario) {
			got = scenarioList
		}),
	)

	if _, err := rec.Reconcile(context.Background(), k8s.Request{}); err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("callback map = %v, want both namespaces' scenarios", got)
	}
}

// TestNewReconciler_RegistersScenarioScheme verifies New()'s reconciler
// registers the ValdBenchmarkScenario scheme (via
// reconciler.WithAddToScheme(v1.AddToScheme)) instead of falling back to the
// listReconciler default (the client-go native scheme), which would leave
// the manager's scheme unable to recognize ValdBenchmarkScenario.
func TestNewReconciler_RegistersScenarioScheme(t *testing.T) {
	t.Parallel()

	watcher, err := New(WithControllerName(testControllerName))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	mgr := &managerMock{client: newFakeClient(t), scheme: runtime.NewScheme()}
	watcher.NewReconciler(context.Background(), mgr)

	if !mgr.scheme.Recognizes(v1.GroupVersion.WithKind("ValdBenchmarkScenario")) {
		t.Error("NewReconciler() did not register the ValdBenchmarkScenario scheme")
	}
}

// NOT IMPLEMENTED BELOW
