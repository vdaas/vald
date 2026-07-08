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

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease"
	"github.com/vdaas/vald/pkg/operator/vald/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// managerMock is a minimal manager.Manager stub for NewReconciler: only
// GetClient and GetScheme are exercised by the resource controller.
type managerMock struct {
	manager.Manager

	client k8s.Client
	scheme *runtime.Scheme
}

func (m *managerMock) GetClient() k8s.Client {
	return m.client
}

func (m *managerMock) GetScheme() *runtime.Scheme {
	return m.scheme
}

func newTestScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add ValdOperatorRelease scheme: %v", err)
	}
	if err := valdrelease.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add ValdRelease scheme: %v", err)
	}
	return scheme
}

func TestResourceController(t *testing.T) {
	t.Parallel()

	rc := newResourceController(&config.Config{})

	t.Run("GetName returns valdoperatorrelease", func(t *testing.T) {
		t.Parallel()
		if got, want := rc.GetName(), Name; got != want {
			t.Errorf("GetName() = %q, want %q", got, want)
		}
	})

	t.Run("For returns ValdOperatorRelease without options", func(t *testing.T) {
		t.Parallel()
		obj, opts := rc.For()
		if _, ok := obj.(*v1.ValdOperatorRelease); !ok {
			t.Errorf("For() object = %T, want *v1.ValdOperatorRelease", obj)
		}
		if opts != nil {
			t.Errorf("For() options = %v, want nil", opts)
		}
	})

	t.Run("Owns returns ValdRelease without options", func(t *testing.T) {
		t.Parallel()
		obj, opts := rc.Owns()
		if _, ok := obj.(*valdrelease.ValdRelease); !ok {
			t.Errorf("Owns() object = %T, want *valdrelease.ValdRelease", obj)
		}
		if opts != nil {
			t.Errorf("Owns() options = %v, want nil", opts)
		}
	})

	t.Run("Watches returns nil", func(t *testing.T) {
		t.Parallel()
		obj, h, opts := rc.Watches()
		if obj != nil || h != nil || opts != nil {
			t.Errorf("Watches() = (%v, %v, %v), want (nil, nil, nil)", obj, h, opts)
		}
	})
}

func TestResourceController_NewReconciler(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	c := fake.NewClientBuilder().WithScheme(newTestScheme(t)).Build()
	cfg := &config.Config{NodePoolLabelPrefix: "vald.vdaas.org"}

	rc := newResourceController(cfg)
	rec := rc.NewReconciler(context.Background(), &managerMock{client: c, scheme: scheme})

	r, ok := rec.(*reconciler)
	if !ok {
		t.Fatalf("NewReconciler() = %T, want *reconciler", rec)
	}
	if r.cfg != cfg {
		t.Error("NewReconciler() did not wire the runtime config")
	}
	if r.syncer == nil {
		t.Error("NewReconciler() did not construct the resource syncer")
	}
	if !scheme.Recognizes(v1.GroupVersion.WithKind("ValdOperatorRelease")) {
		t.Error("NewReconciler() did not register the ValdOperatorRelease scheme")
	}
	if !scheme.Recognizes(valdrelease.GVK) {
		t.Error("NewReconciler() did not register the ValdRelease scheme")
	}
}

//nolint:goconst
func TestReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	type test struct {
		name    string
		request ctrl.Request
	}

	tests := []test{
		{
			name: "returns empty result without error when the resource is not found",
			request: ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "missing",
					Namespace: "default",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			scheme := newTestScheme(t)
			c := fake.NewClientBuilder().WithScheme(scheme).Build()
			// Build the reconciler through NewReconciler so that the
			// type-bound ObjectClient is wired exactly as in production.
			r := newResourceController(&config.Config{}).NewReconciler(
				context.Background(), &managerMock{client: c, scheme: scheme},
			)

			res, err := r.Reconcile(context.Background(), tc.request)
			if err != nil {
				t.Errorf("Reconcile() error = %v, want nil", err)
			}
			if res != (ctrl.Result{}) {
				t.Errorf("Reconcile() result = %v, want empty", res)
			}
		})
	}
}

// --- requeue rules ---

// TestReconciler_Reconcile_RequeueRules verifies the requeue interval per
// reconcile outcome: waiting phases (Progressing/Pending conditions) use
// requeue.success (or the waiting default) instead of requeue.on_error, which
// remains reserved for genuine failures.
func TestReconciler_Reconcile_RequeueRules(t *testing.T) {
	t.Parallel()

	seedFirstCondition := func(cr *v1.ValdOperatorRelease) *v1.ValdOperatorRelease {
		cr.Status.Phase = phaseWaitForClusterCreate
		cr.Status.Conditions = []metav1.Condition{{
			Type:               phaseWaitForClusterCreate,
			Status:             metav1.ConditionUnknown,
			Reason:             "Progressing",
			Message:            "Waiting for Cluster Creation.",
			LastTransitionTime: metav1.Now(),
		}}
		return cr
	}
	// The cluster ID is not provisioned yet: checkClusters reports Pending.
	waitingCR := func() *v1.ValdOperatorRelease {
		return seedFirstCondition(newCRWithInfra([]v1.ValdOperatorReleaseInfra{
			{Role: "green", Clusters: []v1.DestClusters{{ID: "", Name: "cluster-a"}}},
		}))
	}
	// No infrastructure at all: checkClusters reports Failed.
	failingCR := func() *v1.ValdOperatorRelease {
		return seedFirstCondition(newCRWithInfra(nil))
	}

	tests := []struct {
		cr          *v1.ValdOperatorRelease
		cfg         *config.Config
		name        string
		wantRequeue time.Duration
		wantErr     bool
	}{
		{
			name:        "waiting phase uses the waiting default instead of on_error",
			cr:          waitingCR(),
			cfg:         &config.Config{RequeueAfterError: 100 * time.Millisecond},
			wantRequeue: defaultRequeueAfterWaiting,
		},
		{
			name: "waiting phase uses requeue.success when configured",
			cr:   waitingCR(),
			cfg: &config.Config{
				RequeueAfterSuccess: 30 * time.Second,
				RequeueAfterError:   100 * time.Millisecond,
			},
			wantRequeue: 30 * time.Second,
		},
		{
			name:        "failed phase uses on_error",
			cr:          failingCR(),
			cfg:         &config.Config{RequeueAfterError: 100 * time.Millisecond},
			wantRequeue: 100 * time.Millisecond,
		},
		{
			name:    "failed phase without on_error returns the error",
			cr:      failingCR(),
			cfg:     &config.Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scheme := newTestScheme(t)
			c := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tt.cr).
				WithStatusSubresource(tt.cr).
				Build()
			r := &reconciler{
				client: c,
				cfg:    tt.cfg,
				vor:    resource.NewObjectClient[v1.ValdOperatorRelease](c),
			}

			res, err := r.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: types.NamespacedName{Name: tt.cr.Name, Namespace: tt.cr.Namespace},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantRequeue, res.RequeueAfter)
		})
	}
}

// --- result constructors ---

//nolint:goconst
func TestResultConstructors(t *testing.T) {
	tests := []struct {
		name        string
		got         result
		wantStatus  metav1.ConditionStatus
		wantReason  string
		wantMessage string
	}{
		{
			name:        "progressing carries the message",
			got:         progressing("building"),
			wantStatus:  metav1.ConditionUnknown,
			wantReason:  "Progressing",
			wantMessage: "building",
		},
		{
			name:        "pending carries the message",
			got:         pending("waiting for external system"),
			wantStatus:  metav1.ConditionUnknown,
			wantReason:  "Pending",
			wantMessage: "waiting for external system",
		},
		{
			name:       "succeeded has an empty message",
			got:        succeeded(),
			wantStatus: metav1.ConditionTrue,
			wantReason: "Succeeded",
		},
		{
			name:        "failed with error includes the error text",
			got:         failed(errors.New("something broke")),
			wantStatus:  metav1.ConditionFalse,
			wantReason:  "Failed",
			wantMessage: "failed: something broke",
		},
		{
			name:        "failed with nil error uses the fallback message",
			got:         failed(nil),
			wantStatus:  metav1.ConditionFalse,
			wantReason:  "Failed",
			wantMessage: "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantStatus, tt.got.status)
			assert.Equal(t, tt.wantReason, tt.got.reason)
			assert.Equal(t, tt.wantMessage, tt.got.message)
		})
	}
}

// --- phase.condition ---

func TestPhaseCondition(t *testing.T) {
	tests := []struct {
		name       string
		result     result
		wantStatus metav1.ConditionStatus
		wantReason string
	}{
		{
			name:       "Succeeded",
			result:     succeeded(),
			wantStatus: metav1.ConditionTrue,
			wantReason: "Succeeded",
		},
		{
			name:       "Progressing",
			result:     progressing(""),
			wantStatus: metav1.ConditionUnknown,
			wantReason: "Progressing",
		},
		{
			name:       "Pending",
			result:     pending("waiting for VPC"),
			wantStatus: metav1.ConditionUnknown,
			wantReason: "Pending",
		},
		{
			name:       "Failed",
			result:     failed(assert.AnError),
			wantStatus: metav1.ConditionFalse,
			wantReason: "Failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := phase{name: phaseWaitForClusterCreate, message: "default message"}
			got := p.condition(tt.result)
			assert.Equal(t, phaseWaitForClusterCreate, got.Type)
			assert.Equal(t, tt.wantStatus, got.Status)
			assert.Equal(t, tt.wantReason, got.Reason)
			assert.False(t, got.LastTransitionTime.IsZero())
		})
	}
}

func TestPhaseCondition_MessageFallback(t *testing.T) {
	p := phase{name: phaseWaitForClusterCreate, message: "default message"}

	// When the result carries no message, the phase's default message is used.
	assert.Equal(t, "default message", p.condition(progressing("")).Message)
	assert.Equal(t, "default message", p.condition(succeeded()).Message)

	// When the result carries a message, it takes precedence.
	assert.Equal(t, "waiting for VPC", p.condition(pending("waiting for VPC")).Message)
	assert.Contains(t, p.condition(failed(assert.AnError)).Message, assert.AnError.Error())
}

func TestPhaseCondition_StatusesAreDistinct(t *testing.T) {
	p := phase{name: phaseWaitForClusterCreate, message: "msg"}
	assert.NotEqual(t, p.condition(progressing("")).Status, p.condition(succeeded()).Status)
	assert.NotEqual(t, p.condition(progressing("")).Status, p.condition(failed(assert.AnError)).Status)
	assert.NotEqual(t, p.condition(succeeded()).Status, p.condition(failed(assert.AnError)).Status)
}

// --- phases.index ---

func newTestPhases() phases {
	return phases{
		{name: phaseWaitForClusterCreate, check: func(context.Context) result { return succeeded() }},
		{name: phaseWaitForCreateVrs, check: func(context.Context) result { return succeeded() }},
		{name: phaseCompleted},
	}
}

func TestPhases_Index(t *testing.T) {
	ps := newTestPhases()

	tests := []struct {
		name string
		ct   string
		want int
	}{
		{"empty string returns 0", "", 0},
		{"first condition", phaseWaitForClusterCreate, 0},
		{"middle condition", phaseWaitForCreateVrs, 1},
		{"last condition", phaseCompleted, 2},
		{"not found returns -1", "NonExistent", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ps.index(tt.ct))
		})
	}
}

//nolint:goconst
func TestPhases_Pipeline(t *testing.T) {
	r := &reconciler{cfg: &config.Config{}}
	cr := &v1.ValdOperatorRelease{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"}}

	ps := r.phases(cr, alwaysAvailable())

	assert.Len(t, ps, 3)
	assert.Equal(t, phaseWaitForClusterCreate, ps[0].name)
	assert.Nil(t, ps[0].build, "cluster-create phase produces no resources")
	assert.NotNil(t, ps[0].check)

	assert.Equal(t, phaseWaitForCreateVrs, ps[1].name)
	assert.NotNil(t, ps[1].build)
	assert.NotNil(t, ps[1].fetch)
	assert.NotNil(t, ps[1].check)

	assert.Equal(t, phaseCompleted, ps[2].name)
	assert.True(t, ps[2].terminal(), "Completed is a terminal phase")
}

// --- advance ---

func makeAdvancePhases() phases {
	return phases{
		{name: "phase-a", check: func(context.Context) result { return succeeded() }},
		{name: "phase-b", check: func(context.Context) result { return succeeded() }},
		{name: "phase-terminal"}, // no build, no check
	}
}

func newCR() *v1.ValdOperatorRelease {
	return &v1.ValdOperatorRelease{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
	}
}

func TestReconciler_Advance_ToNonTerminal(t *testing.T) {
	r := &reconciler{}
	cr := newCR()

	next := r.advance(makeAdvancePhases(), 0, cr)

	assert.Equal(t, 1, next)
	assert.Equal(t, "phase-b", cr.Status.Phase)
	assert.Equal(t, 1, cr.Status.Progress.Completed)
	assert.Len(t, cr.Status.Conditions, 1)
	assert.Equal(t, "phase-b", cr.Status.Conditions[0].Type)
	assert.Equal(t, metav1.ConditionUnknown, cr.Status.Conditions[0].Status, "non-terminal phases seed Progressing")
}

func TestReconciler_Advance_ToTerminal(t *testing.T) {
	r := &reconciler{}
	cr := newCR()

	next := r.advance(makeAdvancePhases(), 1, cr) // currently at phase-b; next is phase-terminal

	assert.Equal(t, 2, next)
	assert.Equal(t, "phase-terminal", cr.Status.Phase)
	assert.Equal(t, 2, cr.Status.Progress.Completed)
	assert.Equal(t, metav1.ConditionTrue, cr.Status.Conditions[0].Status, "terminal phases seed Succeeded")
}

func TestReconciler_Advance_NoNext(t *testing.T) {
	r := &reconciler{}
	cr := newCR()

	next := r.advance(makeAdvancePhases(), 2, cr) // already at the last phase

	assert.Equal(t, -1, next)
	assert.Empty(t, cr.Status.Phase)
	assert.Empty(t, cr.Status.Conditions)
}

// --- checkClusters ---

func newCRWithInfra(infras []v1.ValdOperatorReleaseInfra) *v1.ValdOperatorRelease {
	return &v1.ValdOperatorRelease{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
		Spec: v1.ValdOperatorReleaseSpec{
			Infrastructure: infras,
		},
	}
}

//nolint:goconst
func TestCheckClusters(t *testing.T) {
	tests := []struct {
		name       string
		cr         *v1.ValdOperatorRelease
		wantStatus metav1.ConditionStatus
		wantReason string
	}{
		{
			name:       "empty infrastructure",
			cr:         newCRWithInfra(nil),
			wantStatus: metav1.ConditionFalse,
			wantReason: "Failed",
		},
		{
			name: "infra with no clusters",
			cr: newCRWithInfra([]v1.ValdOperatorReleaseInfra{
				{Role: "green", Clusters: []v1.DestClusters{}},
			}),
			wantStatus: metav1.ConditionFalse,
			wantReason: "Failed",
		},
		{
			name: "cluster with empty Id",
			cr: newCRWithInfra([]v1.ValdOperatorReleaseInfra{
				{Role: "green", Clusters: []v1.DestClusters{{ID: "", Name: "cluster-a"}}},
			}),
			wantStatus: metav1.ConditionUnknown,
			wantReason: "Pending",
		},
		{
			name: "cluster with empty Name",
			cr: newCRWithInfra([]v1.ValdOperatorReleaseInfra{
				{Role: "green", Clusters: []v1.DestClusters{{ID: "abc-123", Name: ""}}},
			}),
			wantStatus: metav1.ConditionUnknown,
			wantReason: "Pending",
		},
		{
			name: "valid clusters",
			cr: newCRWithInfra([]v1.ValdOperatorReleaseInfra{
				{Role: "green", Clusters: []v1.DestClusters{{ID: "abc-123", Name: "cluster-a"}}},
			}),
			wantStatus: metav1.ConditionTrue,
			wantReason: "Succeeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkClusters(tt.cr)
			assert.Equal(t, tt.wantStatus, got.status)
			assert.Equal(t, tt.wantReason, got.reason)
		})
	}
}

func TestCheckClusters_PendingIsNotFailed(t *testing.T) {
	// cluster.ID == "" means the external system is still creating the cluster. Should be Pending, not Failed.
	cr := newCRWithInfra([]v1.ValdOperatorReleaseInfra{
		{Role: "green", Clusters: []v1.DestClusters{{ID: "", Name: "cluster-a"}}},
	})
	got := checkClusters(cr)
	assert.Equal(t, pending("").status, got.status)
	assert.NotEqual(t, failed(nil).status, got.status)
}

// --- vrsReady ---

func newUnstructuredConfigMap(name, ns string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
	u.SetName(name)
	u.SetNamespace(ns)
	return u
}

func TestReconciler_VrsReady(t *testing.T) {
	tests := []struct {
		want result
		name string
		nilb bool
		seed bool
	}{
		{
			name: "nil built list: failed",
			nilb: true,
			want: failed(errors.ErrBuiltResourceIsNil),
		},
		{
			name: "object not found: progressing",
			seed: false,
			want: progressing(""),
		},
		{
			name: "all objects found: succeeded",
			seed: true,
			want: succeeded(),
		},
	}

	scheme := runtime.NewScheme()
	assert.NoError(t, corev1.AddToScheme(scheme))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := newUnstructuredConfigMap("test", "default")
			cb := fake.NewClientBuilder().WithScheme(scheme)
			if tt.seed {
				cb = cb.WithObjects(cm)
			}
			r := &reconciler{client: cb.Build()}

			built := []k8s.Object{cm}
			if tt.nilb {
				built = nil
			}
			assert.Equal(t, tt.want, r.vrsReady(context.Background(), built))
		})
	}
}

// --- progress bookkeeping through reconcileValdOperatorRelease ---

func TestReconciler_Progress(t *testing.T) {
	tests := []struct {
		name          string
		phase         string
		wantTotal     int
		wantCompleted int
	}{
		{
			name:          "phase at start",
			phase:         phaseWaitForClusterCreate,
			wantTotal:     2,
			wantCompleted: 0,
		},
		{
			name:          "phase at second step",
			phase:         phaseWaitForCreateVrs,
			wantTotal:     2,
			wantCompleted: 1,
		},
		{
			name:          "phase at completed",
			phase:         phaseCompleted,
			wantTotal:     2,
			wantCompleted: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := newCR()
			cr.Status.Phase = tt.phase

			scheme := newTestScheme(t)
			c := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(cr).
				WithStatusSubresource(cr).
				Build()
			r := &reconciler{
				client: c,
				cfg:    &config.Config{},
				vor:    resource.NewObjectClient[v1.ValdOperatorRelease](c),
			}

			// Conditions are empty, so the reconcile seeds the first
			// condition and records the progress derived from the phase.
			assert.NoError(t, r.reconcileValdOperatorRelease(context.Background(), cr))
			assert.Equal(t, tt.wantTotal, cr.Status.Progress.Total)
			assert.Equal(t, tt.wantCompleted, cr.Status.Progress.Completed)
			assert.Len(t, cr.Status.Conditions, 1)
		})
	}
}

// NOT IMPLEMENTED BELOW
