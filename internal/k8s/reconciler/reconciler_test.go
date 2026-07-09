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

package reconciler

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/test/goleak"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

// managerMock is a minimal manager.Manager stub: the reconcilers only
// exercise GetClient, GetScheme and GetFieldIndexer.
type managerMock struct {
	manager.Manager

	client  k8s.Client
	scheme  *runtime.Scheme
	indexer *fieldIndexerMock
}

func (m *managerMock) GetClient() k8s.Client {
	return m.client
}

func (m *managerMock) GetScheme() *runtime.Scheme {
	return m.scheme
}

func (m *managerMock) GetFieldIndexer() kclient.FieldIndexer {
	return m.indexer
}

type fieldIndexerMock struct {
	fields []string
	err    error
}

func (f *fieldIndexerMock) IndexField(
	_ context.Context, _ kclient.Object, field string, _ kclient.IndexerFunc,
) error {
	f.fields = append(f.fields, field)
	return f.err
}

// errClient fails every List/Get call with the configured error.
type errClient struct {
	k8s.Client

	err error
}

func (c *errClient) List(_ context.Context, _ k8s.ObjectList, _ ...k8s.ListOption) error {
	return c.err
}

func (c *errClient) Get(
	_ context.Context, _ k8s.ObjectKey, _ k8s.Object, _ ...k8s.GetOption,
) error {
	return c.err
}

func newManagerMock(t *testing.T, c k8s.Client) *managerMock {
	t.Helper()
	return &managerMock{
		client:  c,
		scheme:  runtime.NewScheme(),
		indexer: new(fieldIndexerMock),
	}
}

func newPod(name, namespace string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": "vald-agent"},
		},
	}
}

func newFakeClient(t *testing.T, objs ...kclient.Object) k8s.Client {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add corev1 scheme: %v", err)
	}
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func TestListReconciler_Metadata(t *testing.T) {
	t.Parallel()

	rc := NewListReconciler(
		"pod watcher",
		new(corev1.Pod),
		func() *corev1.PodList { return new(corev1.PodList) },
	)

	if got, want := rc.GetName(), "pod watcher"; got != want {
		t.Errorf("GetName() = %q, want %q", got, want)
	}
	// The batch reconciler watches through Watches() with a fixed-request
	// handler instead of For(), so the workqueue can deduplicate event waves.
	if obj, fopts := rc.For(); obj != nil || fopts != nil {
		t.Errorf("For() = (%v, %v), want (nil, nil)", obj, fopts)
	}
	if o, oopts := rc.Owns(); o != nil || oopts != nil {
		t.Errorf("Owns() = (%v, %v), want (nil, nil)", o, oopts)
	}
	w, h, wopts := rc.Watches()
	if _, ok := w.(*corev1.Pod); !ok {
		t.Errorf("Watches() object = %T, want *corev1.Pod", w)
	}
	if h == nil {
		t.Error("Watches() handler = nil, want fixed-request handler")
	}
	if wopts != nil {
		t.Errorf("Watches() options = %v, want nil", wopts)
	}
}

// TestListReconciler_Watches_SingleFixedRequest verifies that every event is
// mapped to one fixed reconcile request so the workqueue collapses an event
// wave of N objects into a single full-list reconcile instead of N.
func TestListReconciler_Watches_SingleFixedRequest(t *testing.T) {
	t.Parallel()

	rc := NewListReconciler(
		"pod watcher",
		new(corev1.Pod),
		func() *corev1.PodList { return new(corev1.PodList) },
	)
	lr, ok := rc.(*listReconciler[*corev1.PodList])
	if !ok {
		t.Fatalf("NewListReconciler() = %T, want *listReconciler", rc)
	}
	_, h, _ := lr.Watches()
	if h == nil {
		t.Fatal("Watches() handler = nil")
	}

	q := workqueue.NewTypedRateLimitingQueue(
		workqueue.DefaultTypedControllerRateLimiter[reconcile.Request](),
	)
	defer q.ShutDown()

	ctx := context.Background()
	h.Create(ctx, event.CreateEvent{Object: newPod("a", "default")}, q)
	h.Create(ctx, event.CreateEvent{Object: newPod("b", "default")}, q)
	h.Update(ctx, event.UpdateEvent{
		ObjectOld: newPod("a", "default"),
		ObjectNew: newPod("a", "default"),
	}, q)
	h.Delete(ctx, event.DeleteEvent{Object: newPod("b", "default")}, q)
	h.Generic(ctx, event.GenericEvent{Object: newPod("c", "other")}, q)

	if got := q.Len(); got != 1 {
		t.Fatalf("workqueue length = %d, want 1 (all events must dedup to the fixed request)", got)
	}
	req, _ := q.Get()
	want := reconcile.Request{NamespacedName: types.NamespacedName{Name: "pod watcher"}}
	if req != want {
		t.Errorf("queued request = %+v, want %+v", req, want)
	}
}

func TestListReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	type want struct {
		result        reconcile.Result
		itemCount     int
		err           bool
		callbackFired bool
		onErrorFired  bool
	}
	type test struct {
		name   string
		client func(t *testing.T) k8s.Client
		opts   func(reconciled *int, errored *int) []ListOption[*corev1.PodList]
		want   want
	}

	tests := []test{
		{
			name: "lists objects filtered by namespace and invokes the callback",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t, newPod("a", "default"), newPod("b", "other"))
			},
			opts: func(reconciled, _ *int) []ListOption[*corev1.PodList] {
				return []ListOption[*corev1.PodList]{
					WithNamespace[*corev1.PodList]("default"),
					WithOnReconcile(func(_ context.Context, list *corev1.PodList) {
						*reconciled = len(list.Items)
					}),
				}
			},
			want: want{
				result:        reconcile.Result{},
				itemCount:     1,
				callbackFired: true,
			},
		},
		{
			name: "requeues after the configured success duration",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t, newPod("a", "default"))
			},
			opts: func(reconciled, _ *int) []ListOption[*corev1.PodList] {
				return []ListOption[*corev1.PodList]{
					WithOnReconcile(func(_ context.Context, list *corev1.PodList) {
						*reconciled = len(list.Items)
					}),
					WithRequeueDurations[*corev1.PodList](5*time.Second, 0, 0),
				}
			},
			want: want{
				result:        reconcile.Result{RequeueAfter: 5 * time.Second},
				itemCount:     1,
				callbackFired: true,
			},
		},
		{
			// The error is consumed when errorRequeue is configured because
			// controller-runtime ignores the Result on non-nil errors, which
			// previously made the configured interval ineffective.
			name: "consumes the error and requeues after the configured error duration",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return &errClient{err: errors.New("boom")}
			},
			opts: func(_, errored *int) []ListOption[*corev1.PodList] {
				return []ListOption[*corev1.PodList]{
					WithOnError[*corev1.PodList](func(error) { *errored++ }),
				}
			},
			want: want{
				result:       reconcile.Result{RequeueAfter: defaultErrorRequeueDuration},
				onErrorFired: true,
			},
		},
		{
			name: "swallows NotFound and requeues after the notFound duration",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return &errClient{err: apierrors.NewNotFound(schema.GroupResource{Resource: "pods"}, "missing")}
			},
			opts: func(_, errored *int) []ListOption[*corev1.PodList] {
				return []ListOption[*corev1.PodList]{
					WithOnError[*corev1.PodList](func(error) { *errored++ }),
				}
			},
			want: want{
				result:       reconcile.Result{RequeueAfter: defaultNotFoundRequeueDuration},
				onErrorFired: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var reconciled, errored int
			rc := NewListReconciler(
				"pod watcher",
				new(corev1.Pod),
				func() *corev1.PodList { return new(corev1.PodList) },
				tc.opts(&reconciled, &errored)...,
			)
			mgr := newManagerMock(t, tc.client(t))
			rec := rc.NewReconciler(context.Background(), mgr)

			res, err := rec.Reconcile(context.Background(), reconcile.Request{})
			if tc.want.err != (err != nil) {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tc.want.err)
			}
			if res != tc.want.result {
				t.Errorf("Reconcile() result = %+v, want %+v", res, tc.want.result)
			}
			if tc.want.callbackFired != (reconciled > 0) {
				t.Errorf("callback fired = %v, want %v", reconciled > 0, tc.want.callbackFired)
			}
			if tc.want.callbackFired && reconciled != tc.want.itemCount {
				t.Errorf("callback item count = %d, want %d", reconciled, tc.want.itemCount)
			}
			if tc.want.onErrorFired != (errored > 0) {
				t.Errorf("onError fired = %v, want %v", errored > 0, tc.want.onErrorFired)
			}
		})
	}
}

func TestListReconciler_NewReconciler(t *testing.T) {
	t.Parallel()

	rc := NewListReconciler(
		"pod watcher",
		new(corev1.Pod),
		func() *corev1.PodList { return new(corev1.PodList) },
		WithFieldIndex[*corev1.PodList]("status.phase", func(k8s.Object) []string { return nil }),
	)
	mgr := newManagerMock(t, newFakeClient(t))
	rc.NewReconciler(context.Background(), mgr)

	if !mgr.scheme.Recognizes(corev1.SchemeGroupVersion.WithKind("Pod")) {
		t.Error("NewReconciler() did not register the default client-go scheme")
	}
	if len(mgr.indexer.fields) != 1 || mgr.indexer.fields[0] != "status.phase" {
		t.Errorf("registered field indexes = %v, want [status.phase]", mgr.indexer.fields)
	}
}

// TestReconciler_InitFailure guards the deferred error propagation: when
// NewReconciler cannot complete its setup (nil manager, failed index
// registration), Reconcile must fail fast with the recorded root cause
// instead of panicking or silently reconciling against a broken setup.
func TestReconciler_InitFailure(t *testing.T) {
	t.Parallel()

	type test struct {
		name    string
		rec     func(t *testing.T) reconcile.Reconciler
		wantMsg string
	}

	tests := []test{
		{
			name: "list reconciler with nil manager fails instead of panicking",
			rec: func(t *testing.T) reconcile.Reconciler {
				t.Helper()
				rc := NewListReconciler(
					"pods",
					new(corev1.Pod),
					func() *corev1.PodList { return new(corev1.PodList) },
				)
				return rc.NewReconciler(context.Background(), nil)
			},
			wantMsg: "manager is not registered",
		},
		{
			name: "list reconciler surfaces field index registration failure",
			rec: func(t *testing.T) reconcile.Reconciler {
				t.Helper()
				rc := NewListReconciler(
					"pods",
					new(corev1.Pod),
					func() *corev1.PodList { return new(corev1.PodList) },
					WithFieldIndex[*corev1.PodList]("status.phase", func(k8s.Object) []string { return nil }),
				)
				mgr := newManagerMock(t, newFakeClient(t))
				mgr.indexer.err = errors.New("index boom")
				return rc.NewReconciler(context.Background(), mgr)
			},
			wantMsg: "failed to register field index",
		},
		{
			name: "object reconciler with nil manager fails instead of panicking",
			rec: func(t *testing.T) reconcile.Reconciler {
				t.Helper()
				rc := NewObjectReconciler(
					"pod",
					func() *corev1.Pod { return new(corev1.Pod) },
				)
				return rc.NewReconciler(context.Background(), nil)
			},
			wantMsg: "manager is not registered",
		},
		{
			name: "object reconciler surfaces field index registration failure",
			rec: func(t *testing.T) reconcile.Reconciler {
				t.Helper()
				rc := NewObjectReconciler(
					"pod",
					func() *corev1.Pod { return new(corev1.Pod) },
					WithObjectFieldIndex[*corev1.Pod]("status.phase", func(k8s.Object) []string { return nil }),
				)
				mgr := newManagerMock(t, newFakeClient(t))
				mgr.indexer.err = errors.New("index boom")
				return rc.NewReconciler(context.Background(), mgr)
			},
			wantMsg: "failed to register field index",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := tc.rec(t)
			_, err := rec.Reconcile(context.Background(), reconcile.Request{})
			if err == nil {
				t.Fatal("Reconcile() error = nil, want init failure")
			}
			if !strings.Contains(err.Error(), tc.wantMsg) {
				t.Errorf("Reconcile() error = %v, want it to contain %q", err, tc.wantMsg)
			}
		})
	}
}

func TestObjectReconciler_Metadata(t *testing.T) {
	t.Parallel()

	rc := NewObjectReconciler(
		"pod object watcher",
		func() *corev1.Pod { return new(corev1.Pod) },
	)

	if got, want := rc.GetName(), "pod object watcher"; got != want {
		t.Errorf("GetName() = %q, want %q", got, want)
	}
	obj, fopts := rc.For()
	if _, ok := obj.(*corev1.Pod); !ok {
		t.Errorf("For() object = %T, want *corev1.Pod", obj)
	}
	if fopts != nil {
		t.Errorf("For() options = %v, want nil", fopts)
	}
	if o, oopts := rc.Owns(); o != nil || oopts != nil {
		t.Errorf("Owns() = (%v, %v), want (nil, nil)", o, oopts)
	}
	if w, h, wopts := rc.Watches(); w != nil || h != nil || wopts != nil {
		t.Errorf("Watches() = (%v, %v, %v), want (nil, nil, nil)", w, h, wopts)
	}
}

func TestObjectReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	type want struct {
		result        reconcile.Result
		err           bool
		callbackFired bool
		onErrorFired  bool
	}
	type test struct {
		name   string
		client func(t *testing.T) k8s.Client
		want   want
	}

	customResult := reconcile.Result{Requeue: true, RequeueAfter: 3 * time.Second}

	tests := []test{
		{
			name: "fetches the object and returns the callback result",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t, newPod("target", "default"))
			},
			want: want{
				result:        customResult,
				callbackFired: true,
			},
		},
		{
			name: "ignores NotFound without invoking the callback",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return newFakeClient(t)
			},
			want: want{
				result: reconcile.Result{},
			},
		},
		{
			name: "returns the error on generic get failure",
			client: func(t *testing.T) k8s.Client {
				t.Helper()
				return &errClient{err: errors.New("boom")}
			},
			want: want{
				result:       reconcile.Result{},
				err:          true,
				onErrorFired: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var gotName string
			var errored int
			rc := NewObjectReconciler(
				"pod object watcher",
				func() *corev1.Pod { return new(corev1.Pod) },
				WithOnObjectReconcile(func(_ context.Context, pod *corev1.Pod) (k8s.Result, error) {
					gotName = pod.GetName()
					return customResult, nil
				}),
				WithObjectOnError[*corev1.Pod](func(error) { errored++ }),
			)
			mgr := newManagerMock(t, tc.client(t))
			rec := rc.NewReconciler(context.Background(), mgr)

			res, err := rec.Reconcile(context.Background(), reconcile.Request{
				NamespacedName: types.NamespacedName{Name: "target", Namespace: "default"},
			})
			if tc.want.err != (err != nil) {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tc.want.err)
			}
			if res != tc.want.result {
				t.Errorf("Reconcile() result = %+v, want %+v", res, tc.want.result)
			}
			if tc.want.callbackFired != (gotName != "") {
				t.Errorf("callback fired = %v, want %v", gotName != "", tc.want.callbackFired)
			}
			if tc.want.callbackFired && gotName != "target" {
				t.Errorf("callback object name = %q, want %q", gotName, "target")
			}
			if tc.want.onErrorFired != (errored > 0) {
				t.Errorf("onError fired = %v, want %v", errored > 0, tc.want.onErrorFired)
			}
		})
	}
}

func TestMaxConcurrentReconciles(t *testing.T) {
	t.Parallel()

	type test struct {
		name string
		rc   k8s.ResourceController
		want int
	}

	tests := []test{
		{
			name: "list reconciler defaults to zero",
			rc: NewListReconciler(
				"pods",
				new(corev1.Pod),
				func() *corev1.PodList { return new(corev1.PodList) },
			),
			want: 0,
		},
		{
			name: "list reconciler exposes the configured worker count",
			rc: NewListReconciler(
				"pods",
				new(corev1.Pod),
				func() *corev1.PodList { return new(corev1.PodList) },
				WithMaxConcurrentReconciles[*corev1.PodList](4),
			),
			want: 4,
		},
		{
			name: "list reconciler ignores non-positive worker counts",
			rc: NewListReconciler(
				"pods",
				new(corev1.Pod),
				func() *corev1.PodList { return new(corev1.PodList) },
				WithMaxConcurrentReconciles[*corev1.PodList](-1),
			),
			want: 0,
		},
		{
			name: "object reconciler defaults to zero",
			rc: NewObjectReconciler(
				"pod",
				func() *corev1.Pod { return new(corev1.Pod) },
			),
			want: 0,
		},
		{
			name: "object reconciler exposes the configured worker count",
			rc: NewObjectReconciler(
				"pod",
				func() *corev1.Pod { return new(corev1.Pod) },
				WithObjectMaxConcurrentReconciles[*corev1.Pod](2),
			),
			want: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cr, ok := tc.rc.(k8s.ConcurrentReconciler)
			if !ok {
				t.Fatalf("%T does not implement k8s.ConcurrentReconciler", tc.rc)
			}
			if got := cr.MaxConcurrentReconciles(); got != tc.want {
				t.Errorf("MaxConcurrentReconciles() = %d, want %d", got, tc.want)
			}
		})
	}
}
