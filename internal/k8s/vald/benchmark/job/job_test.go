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

package job

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/test/goleak"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

// testControllerName is the controller name used across the tests, matching
// the name pkg/operator/benchmark/service registers in production.
const testControllerName = "benchmark job resource"

// managerMock is a minimal manager.Manager stub: New()'s reconciler only
// exercises GetClient and GetScheme.
type managerMock struct {
	manager.Manager

	client k8s.Client
	scheme *runtime.Scheme
}

func (m *managerMock) GetClient() k8s.Client      { return m.client }
func (m *managerMock) GetScheme() *runtime.Scheme { return m.scheme }

// GetConfig satisfies client.NewFromManager's k8s.Manager requirement: the
// embedded manager.Manager is nil, so without this override GetConfig would
// panic instead of returning a config client.NewFromManager can build a
// (non-dialing) clientset from.
func (*managerMock) GetConfig() *rest.Config {
	return &rest.Config{}
}

func newManagerMock(t *testing.T, c k8s.Client) *managerMock {
	t.Helper()
	return &managerMock{client: c, scheme: runtime.NewScheme()}
}

// errClient fails every List call with the configured error, regardless of
// the ListOptions passed to it.
type errClient struct {
	k8s.Client

	err error
}

func (c *errClient) List(_ context.Context, _ k8s.ObjectList, _ ...k8s.ListOption) error {
	return c.err
}

func newBenchJob(name, namespace string) *v1.ValdBenchmarkJob {
	return &v1.ValdBenchmarkJob{
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
		t.Fatalf("failed to add v1 scheme: %v", err)
	}
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

// TestOptions verifies that each public Option mutates the expected field on
// settings, and that the registered callbacks are invoked with the documented
// signatures.
func TestOptions(t *testing.T) {
	t.Parallel()

	t.Run("WithControllerName sets name", func(t *testing.T) {
		t.Parallel()
		s := new(settings)
		if err := WithControllerName("bench job watcher")(s); err != nil {
			t.Fatalf("WithControllerName() error = %v", err)
		}
		if got, want := s.Name, "bench job watcher"; got != want {
			t.Errorf("name = %q, want %q", got, want)
		}
	})

	t.Run("WithNamespaces sets namespaces", func(t *testing.T) {
		t.Parallel()
		s := new(settings)
		if err := WithNamespaces("ns-a", "ns-b")(s); err != nil {
			t.Fatalf("WithNamespaces() error = %v", err)
		}
		want := []string{"ns-a", "ns-b"}
		if len(s.Namespaces) != len(want) {
			t.Fatalf("namespaces = %v, want %v", s.Namespaces, want)
		}
		for i := range want {
			if s.Namespaces[i] != want[i] {
				t.Errorf("namespaces[%d] = %q, want %q", i, s.Namespaces[i], want[i])
			}
		}
	})

	t.Run("WithOnErrorFunc registers a callback invoked on error", func(t *testing.T) {
		t.Parallel()
		s := new(settings)
		var got error
		if err := WithOnErrorFunc(func(err error) { got = err })(s); err != nil {
			t.Fatalf("WithOnErrorFunc() error = %v", err)
		}
		if s.OnError == nil {
			t.Fatal("onError = nil, want a registered callback")
		}
		want := errors.New("boom")
		s.OnError(want)
		if !errors.Is(got, want) {
			t.Errorf("onError callback received = %v, want %v", got, want)
		}
	})

	t.Run("WithOnReconcileFunc registers a callback receiving map[string]v1.ValdBenchmarkJob", func(t *testing.T) {
		t.Parallel()
		s := new(settings)
		var got map[string]v1.ValdBenchmarkJob
		if err := WithOnReconcileFunc(func(_ context.Context, jobList map[string]v1.ValdBenchmarkJob) {
			got = jobList
		})(s); err != nil {
			t.Fatalf("WithOnReconcileFunc() error = %v", err)
		}
		if s.OnReconcile == nil {
			t.Fatal("onReconcile = nil, want a registered callback")
		}
		want := map[string]v1.ValdBenchmarkJob{"a": *newBenchJob("a", "default")}
		s.OnReconcile(context.Background(), want)
		if len(got) != 1 {
			t.Fatalf("onReconcile callback received %d items, want 1", len(got))
		}
		if _, ok := got["a"]; !ok {
			t.Errorf("onReconcile callback received = %v, want key %q present", got, "a")
		}
	})
}

// TestNew_OptionErrorSeverity pins down the unified severity policy shared
// with the mirror target watcher: a non-critical option error is warned and
// skipped so construction still succeeds, while an errors.ErrCriticalOption
// aborts New with the wrapped failure.
func TestNew_OptionErrorSeverity(t *testing.T) {
	t.Parallel()

	t.Run("non-critical option error is skipped and construction succeeds", func(t *testing.T) {
		t.Parallel()

		failing := Option(func(*settings) error { return errors.New("option boom") })
		got, err := New(failing, WithControllerName(testControllerName))
		if err != nil {
			t.Fatalf("New() error = %v, want nil (non-critical option errors must be skipped)", err)
		}
		if got == nil {
			t.Fatal("New() watcher = nil, want non-nil")
		}
		if name := got.GetName(); name != testControllerName {
			t.Errorf("GetName() = %q, want %q (options after the failing one must still apply)", name, testControllerName)
		}
	})

	t.Run("critical option error aborts construction", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.NewErrCriticalOption("name", "value")
		failing := Option(func(*settings) error { return wantErr })
		got, err := New(failing)
		if err == nil {
			t.Fatal("New() error = nil, want a wrapped critical option failure")
		}
		if !errors.Is(err, wantErr) {
			t.Errorf("New() error = %v, want it to wrap %v", err, wantErr)
		}
		if got != nil {
			t.Errorf("New() watcher = %v, want nil on critical option failure", got)
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

			s := new(settings)
			e := &errors.ErrInvalidOption{}
			if err := tc.opts[0](s); !errors.As(err, &e) {
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

	rc, err := New(
		WithControllerName(testControllerName),
		WithOnErrorFunc(func(err error) {}),
		WithNamespaces("default"),
		WithOnErrorFunc(func(err error) {}),
		WithOnReconcileFunc(func(context.Context, map[string]v1.ValdBenchmarkJob) {}),
	)
	if err != nil {
		t.Fatalf("New() error = %v, want nil for the production option pattern", err)
	}
	if rc == nil {
		t.Fatal("New() watcher = nil, want non-nil")
	}
	if got, want := rc.GetName(), testControllerName; got != want {
		t.Errorf("GetName() = %q, want %q", got, want)
	}
}

// TestNew_Metadata verifies New() returns a BenchmarkWatcher whose
// GetName() reflects WithControllerName, and that it plugs into the batch
// watch/list contract (For unused, Watches installed) rather than the
// hand-rolled For()-based contract the old reconciler used.
func TestNew_Metadata(t *testing.T) {
	t.Parallel()

	rc, err := New(WithControllerName(testControllerName))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got, want := rc.GetName(), testControllerName; got != want {
		t.Errorf("GetName() = %q, want %q", got, want)
	}
	if obj, fopts := rc.For(); obj != nil || fopts != nil {
		t.Errorf("For() = (%v, %v), want (nil, nil)", obj, fopts)
	}
	w, h, wopts := rc.Watches()
	if _, ok := w.(*v1.ValdBenchmarkJob); !ok {
		t.Errorf("Watches() object = %T, want *v1.ValdBenchmarkJob", w)
	}
	if h == nil {
		t.Error("Watches() handler = nil, want fixed-request handler")
	}
	if wopts != nil {
		t.Errorf("Watches() options = %v, want nil", wopts)
	}
}

// TestNew_Reconcile drives New()'s watcher through NewReconciler+Reconcile
// against a fake controller-runtime client, exercising the externally
// observable behavior of options end to end.
func TestNew_Reconcile(t *testing.T) {
	t.Parallel()

	type want struct {
		result        reconcile.Result
		jobNames      []string
		err           bool
		callbackFired bool
		onErrorFired  bool
	}
	type test struct {
		name   string
		client func(t *testing.T) k8s.Client
		opts   func(reconciled *map[string]v1.ValdBenchmarkJob, errored *error) []Option
		want   want
	}

	tests := []test{
		{
			// Baseline: no namespace restriction lists every namespace and
			// invokes the callback with a map keyed by job name.
			name: "no namespace filter lists every job and invokes the callback",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t, newBenchJob("a", "ns-a"), newBenchJob("b", "ns-b"))
			},
			opts: func(reconciled *map[string]v1.ValdBenchmarkJob, _ *error) []Option {
				return []Option{
					WithOnReconcileFunc(func(_ context.Context, jobList map[string]v1.ValdBenchmarkJob) {
						*reconciled = jobList
					}),
				}
			},
			want: want{
				result:        reconcile.Result{},
				jobNames:      []string{"a", "b"},
				callbackFired: true,
			},
		},
		{
			// Regression test for the WithNamespaces bug: before the fix,
			// WithNamespaces only recorded r.namespaces and never translated
			// it into a k8s.InNamespace ListOption, so List() always
			// traversed every namespace. This must now observe only the job
			// in the requested namespace.
			name: "WithNamespaces restricts the callback to the given namespace",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t, newBenchJob("a", "ns-a"), newBenchJob("b", "ns-b"))
			},
			opts: func(reconciled *map[string]v1.ValdBenchmarkJob, _ *error) []Option {
				return []Option{
					WithNamespaces("ns-a"),
					WithOnReconcileFunc(func(_ context.Context, jobList map[string]v1.ValdBenchmarkJob) {
						*reconciled = jobList
					}),
				}
			},
			want: want{
				result:        reconcile.Result{},
				jobNames:      []string{"a"},
				callbackFired: true,
			},
		},
		{
			name: "list error invokes WithOnErrorFunc and requeues after the error duration",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return &errClient{err: errors.New("boom")}
			},
			opts: func(_ *map[string]v1.ValdBenchmarkJob, errored *error) []Option {
				return []Option{
					WithOnErrorFunc(func(err error) { *errored = err }),
				}
			},
			want: want{
				result:       reconcile.Result{RequeueAfter: 100 * time.Millisecond},
				onErrorFired: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var reconciled map[string]v1.ValdBenchmarkJob
			var errored error
			rc, err := New(append([]Option{WithControllerName(testControllerName)},
				tc.opts(&reconciled, &errored)...)...)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			mgr := newManagerMock(t, tc.client(t))
			rec := rc.NewReconciler(context.Background(), mgr)

			res, err := rec.Reconcile(context.Background(), reconcile.Request{})
			if tc.want.err != (err != nil) {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tc.want.err)
			}
			if res != tc.want.result {
				t.Errorf("Reconcile() result = %+v, want %+v", res, tc.want.result)
			}
			if tc.want.callbackFired != (reconciled != nil) {
				t.Errorf("callback fired = %v, want %v", reconciled != nil, tc.want.callbackFired)
			}
			if tc.want.callbackFired {
				if len(reconciled) != len(tc.want.jobNames) {
					t.Fatalf("callback map = %v, want keys %v", reconciled, tc.want.jobNames)
				}
				for _, name := range tc.want.jobNames {
					if _, ok := reconciled[name]; !ok {
						t.Errorf("callback map = %v, want key %q present", reconciled, name)
					}
				}
			}
			if tc.want.onErrorFired != (errored != nil) {
				t.Errorf("onError fired = %v, want %v", errored != nil, tc.want.onErrorFired)
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
// 		want BenchmarkWatcher
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, BenchmarkWatcher, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got BenchmarkWatcher, err error) error {
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
