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
	mock "github.com/vdaas/vald/internal/test/mock/k8s"
	"github.com/vdaas/vald/pkg/operator/vald/config"
)

const (
	testClusterNameA = "cluster-a"
	testClusterIDABC = "abc-123"
	testRoleGreen    = "green"
)

// managerMock is a minimal k8s.Manager stub for NewReconciler: GetClient and
// GetScheme are exercised by the resource controller directly, and GetConfig
// is exercised indirectly through client.NewFromManager, which builds the vor
// resource client.
type managerMock struct {
	k8s.Manager

	client k8s.Client
	scheme *k8s.Scheme
}

func (m *managerMock) GetClient() k8s.Client {
	return m.client
}

func (m *managerMock) GetScheme() *k8s.Scheme {
	return m.scheme
}

func (*managerMock) GetConfig() *k8s.RESTConfig {
	return &k8s.RESTConfig{}
}

func newTestScheme(t *testing.T) *k8s.Scheme {
	t.Helper()
	scheme := k8s.NewScheme()
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
		if got, want := rc.GetName(), name; got != want {
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

	scheme := k8s.NewScheme()
	c := mock.NewFakeClientBuilder().WithScheme(newTestScheme(t)).Build()
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
		request k8s.Request
	}

	tests := []test{
		{
			name: "returns empty result without error when the resource is not found",
			request: k8s.Request{
				NamespacedName: k8s.NamespacedName{
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
			c := mock.NewFakeClientBuilder().WithScheme(scheme).Build()
			// Build the reconciler through NewReconciler so that the
			// type-bound ObjectClient is wired exactly as in production.
			r := newResourceController(&config.Config{}).NewReconciler(
				context.Background(), &managerMock{client: c, scheme: scheme},
			)

			res, err := r.Reconcile(context.Background(), tc.request)
			if err != nil {
				t.Errorf("Reconcile() error = %v, want nil", err)
			}
			if res != (k8s.Result{}) {
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
		cr.Status.Conditions = []k8s.Condition{{
			Type:               phaseWaitForClusterCreate,
			Status:             k8s.ConditionUnknown,
			Reason:             conditionReasonProgressing,
			Message:            "Waiting for Cluster Creation.",
			LastTransitionTime: k8s.Now(),
		}}
		return cr
	}
	// The cluster ID is not provisioned yet: checkClusters reports Pending.
	waitingCR := func() *v1.ValdOperatorRelease {
		return seedFirstCondition(newCRWithInfra([]v1.ValdOperatorReleaseInfra{
			{Role: testRoleGreen, Clusters: []v1.DestClusters{{ID: "", Name: testClusterNameA}}},
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
			c := mock.NewFakeClientBuilder().
				WithScheme(scheme).
				WithObjects(tt.cr).
				WithStatusSubresource(tt.cr).
				Build()
			r := &reconciler{
				client: c,
				cfg:    tt.cfg,
				vor:    resource.NewClient(c, new(v1.ValdOperatorRelease), new(v1.ValdOperatorReleaseList)),
			}

			res, err := r.Reconcile(context.Background(), k8s.Request{
				NamespacedName: k8s.NamespacedName{Name: tt.cr.Name, Namespace: tt.cr.Namespace},
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

func TestResultConstructors(t *testing.T) {
	tests := []struct {
		name        string
		got         result
		wantStatus  k8s.ConditionStatus
		wantReason  string
		wantMessage string
	}{
		{
			name:        "progressing carries the message",
			got:         progressing("building"),
			wantStatus:  k8s.ConditionUnknown,
			wantReason:  conditionReasonProgressing,
			wantMessage: "building",
		},
		{
			name:        "pending carries the message",
			got:         pending("waiting for external system"),
			wantStatus:  k8s.ConditionUnknown,
			wantReason:  conditionReasonPending,
			wantMessage: "waiting for external system",
		},
		{
			name:       "succeeded has an empty message",
			got:        succeeded(),
			wantStatus: k8s.ConditionTrue,
			wantReason: conditionReasonSucceeded,
		},
		{
			name:        "failed with error includes the error text",
			got:         failed(errors.New("something broke")),
			wantStatus:  k8s.ConditionFalse,
			wantReason:  conditionReasonFailed,
			wantMessage: "failed: something broke",
		},
		{
			name:        "failed with nil error uses the fallback message",
			got:         failed(nil),
			wantStatus:  k8s.ConditionFalse,
			wantReason:  conditionReasonFailed,
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
		wantStatus k8s.ConditionStatus
		wantReason string
	}{
		{
			name:       "Succeeded",
			result:     succeeded(),
			wantStatus: k8s.ConditionTrue,
			wantReason: "Succeeded",
		},
		{
			name:       "Progressing",
			result:     progressing(""),
			wantStatus: k8s.ConditionUnknown,
			wantReason: "Progressing",
		},
		{
			name:       "Pending",
			result:     pending("waiting for VPC"),
			wantStatus: k8s.ConditionUnknown,
			wantReason: "Pending",
		},
		{
			name:       "Failed",
			result:     failed(assert.AnError),
			wantStatus: k8s.ConditionFalse,
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
	cr := &v1.ValdOperatorRelease{ObjectMeta: k8s.ObjectMeta{Name: "test", Namespace: "default"}}

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
		ObjectMeta: k8s.ObjectMeta{Name: "test", Namespace: "default"},
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
	assert.Equal(t, k8s.ConditionUnknown, cr.Status.Conditions[0].Status, "non-terminal phases seed Progressing")
}

func TestReconciler_Advance_ToTerminal(t *testing.T) {
	r := &reconciler{}
	cr := newCR()

	next := r.advance(makeAdvancePhases(), 1, cr) // currently at phase-b; next is phase-terminal

	assert.Equal(t, 2, next)
	assert.Equal(t, "phase-terminal", cr.Status.Phase)
	assert.Equal(t, 2, cr.Status.Progress.Completed)
	assert.Equal(t, k8s.ConditionTrue, cr.Status.Conditions[0].Status, "terminal phases seed Succeeded")
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
		ObjectMeta: k8s.ObjectMeta{Name: "test", Namespace: "default"},
		Spec: v1.ValdOperatorReleaseSpec{
			Infrastructure: infras,
		},
	}
}

func TestCheckClusters(t *testing.T) {
	tests := []struct {
		name       string
		cr         *v1.ValdOperatorRelease
		wantStatus k8s.ConditionStatus
		wantReason string
	}{
		{
			name:       "empty infrastructure",
			cr:         newCRWithInfra(nil),
			wantStatus: k8s.ConditionFalse,
			wantReason: conditionReasonFailed,
		},
		{
			name: "infra with no clusters",
			cr: newCRWithInfra([]v1.ValdOperatorReleaseInfra{
				{Role: testRoleGreen, Clusters: []v1.DestClusters{}},
			}),
			wantStatus: k8s.ConditionFalse,
			wantReason: conditionReasonFailed,
		},
		{
			name: "cluster with empty Id",
			cr: newCRWithInfra([]v1.ValdOperatorReleaseInfra{
				{Role: testRoleGreen, Clusters: []v1.DestClusters{{ID: "", Name: testClusterNameA}}},
			}),
			wantStatus: k8s.ConditionUnknown,
			wantReason: conditionReasonPending,
		},
		{
			name: "cluster with empty Name",
			cr: newCRWithInfra([]v1.ValdOperatorReleaseInfra{
				{Role: testRoleGreen, Clusters: []v1.DestClusters{{ID: testClusterIDABC, Name: ""}}},
			}),
			wantStatus: k8s.ConditionUnknown,
			wantReason: conditionReasonPending,
		},
		{
			name: "valid clusters",
			cr: newCRWithInfra([]v1.ValdOperatorReleaseInfra{
				{Role: testRoleGreen, Clusters: []v1.DestClusters{{ID: testClusterIDABC, Name: testClusterNameA}}},
			}),
			wantStatus: k8s.ConditionTrue,
			wantReason: conditionReasonSucceeded,
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
		{Role: testRoleGreen, Clusters: []v1.DestClusters{{ID: "", Name: testClusterNameA}}},
	})
	got := checkClusters(cr)
	assert.Equal(t, pending("").status, got.status)
	assert.NotEqual(t, failed(nil).status, got.status)
}

// --- vrsReady ---

func newUnstructuredConfigMap(name, ns string) *k8s.Unstructured {
	u := &k8s.Unstructured{}
	u.SetGroupVersionKind(k8s.GroupVersionKind{Version: "v1", Kind: "ConfigMap"})
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

	scheme := k8s.NewScheme()
	assert.NoError(t, k8s.AddClientGoScheme(scheme))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := newUnstructuredConfigMap("test", "default")
			cb := mock.NewFakeClientBuilder().WithScheme(scheme)
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
			c := mock.NewFakeClientBuilder().
				WithScheme(scheme).
				WithObjects(cr).
				WithStatusSubresource(cr).
				Build()
			r := &reconciler{
				client: c,
				cfg:    &config.Config{},
				vor:    resource.NewClient(c, new(v1.ValdOperatorRelease), new(v1.ValdOperatorReleaseList)),
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
//
// func Test_newResourceController(t *testing.T) {
// 	type args struct {
// 		cfg *config.Config
// 	}
// 	type want struct {
// 		want k8s.ResourceController
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, k8s.ResourceController) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.ResourceController) error {
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
// 		           cfg:nil,
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
// 		           cfg:nil,
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
// 			got := newResourceController(test.args.cfg)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_resourceController_GetName(t *testing.T) {
// 	type fields struct {
// 		cfg *config.Config
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
// 		           cfg:nil,
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
// 		           cfg:nil,
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
// 			r := &resourceController{
// 				cfg: test.fields.cfg,
// 			}
//
// 			got := r.GetName()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_resourceController_NewReconciler(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 		mgr k8s.Manager
// 	}
// 	type fields struct {
// 		cfg *config.Config
// 	}
// 	type want struct {
// 		want k8s.Reconciler
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, k8s.Reconciler) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.Reconciler) error {
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
// 		           cfg:nil,
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
// 		           cfg:nil,
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
// 			rc := &resourceController{
// 				cfg: test.fields.cfg,
// 			}
//
// 			got := rc.NewReconciler(test.args.in0, test.args.mgr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_resourceController_MaxConcurrentReconciles(t *testing.T) {
// 	type fields struct {
// 		cfg *config.Config
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
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
// 		           cfg:nil,
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
// 		           cfg:nil,
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
// 			rc := &resourceController{
// 				cfg: test.fields.cfg,
// 			}
//
// 			got := rc.MaxConcurrentReconciles()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_resourceController_For(t *testing.T) {
// 	type fields struct {
// 		cfg *config.Config
// 	}
// 	type want struct {
// 		want  k8s.Object
// 		want1 []k8s.ForOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, k8s.Object, []k8s.ForOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.Object, got1 []k8s.ForOption) error {
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
// 		           cfg:nil,
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
// 		           cfg:nil,
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
// 			r := &resourceController{
// 				cfg: test.fields.cfg,
// 			}
//
// 			got, got1 := r.For()
// 			if err := checkFunc(test.want, got, got1); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_resourceController_Owns(t *testing.T) {
// 	type fields struct {
// 		cfg *config.Config
// 	}
// 	type want struct {
// 		want  k8s.Object
// 		want1 []k8s.OwnsOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, k8s.Object, []k8s.OwnsOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.Object, got1 []k8s.OwnsOption) error {
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
// 		           cfg:nil,
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
// 		           cfg:nil,
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
// 			r := &resourceController{
// 				cfg: test.fields.cfg,
// 			}
//
// 			got, got1 := r.Owns()
// 			if err := checkFunc(test.want, got, got1); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_resourceController_Watches(t *testing.T) {
// 	type fields struct {
// 		cfg *config.Config
// 	}
// 	type want struct {
// 		want  k8s.Object
// 		want1 k8s.EventHandler
// 		want2 []k8s.WatchesOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, k8s.Object, k8s.EventHandler, []k8s.WatchesOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.Object, got1 k8s.EventHandler, got2 []k8s.WatchesOption) error {
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
// 		           cfg:nil,
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
// 		           cfg:nil,
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
// 			r := &resourceController{
// 				cfg: test.fields.cfg,
// 			}
//
// 			got, got1, got2 := r.Watches()
// 			if err := checkFunc(test.want, got, got1, got2); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_reconciler_Reconcile(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req k8s.Request
// 	}
// 	type fields struct {
// 		client  k8s.Client
// 		cfg     *config.Config
// 		syncer  *resource.Syncer
// 		vor     *resource.Client[*v1.ValdOperatorRelease, *v1.ValdOperatorReleaseList]
// 		initErr error
// 	}
// 	type want struct {
// 		want k8s.Result
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, k8s.Result, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.Result, err error) error {
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 				client:  test.fields.client,
// 				cfg:     test.fields.cfg,
// 				syncer:  test.fields.syncer,
// 				vor:     test.fields.vor,
// 				initErr: test.fields.initErr,
// 			}
//
// 			got, err := r.Reconcile(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_reconciler_reconcileValdOperatorRelease(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		cr  *v1.ValdOperatorRelease
// 	}
// 	type fields struct {
// 		client  k8s.Client
// 		cfg     *config.Config
// 		syncer  *resource.Syncer
// 		vor     *resource.Client[*v1.ValdOperatorRelease, *v1.ValdOperatorReleaseList]
// 		initErr error
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
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
// 		           cr:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 		           cr:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 				client:  test.fields.client,
// 				cfg:     test.fields.cfg,
// 				syncer:  test.fields.syncer,
// 				vor:     test.fields.vor,
// 				initErr: test.fields.initErr,
// 			}
//
// 			err := r.reconcileValdOperatorRelease(test.args.ctx, test.args.cr)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_reconciler_reconcilePhase(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		p   *phase
// 		cr  *v1.ValdOperatorRelease
// 	}
// 	type fields struct {
// 		client  k8s.Client
// 		cfg     *config.Config
// 		syncer  *resource.Syncer
// 		vor     *resource.Client[*v1.ValdOperatorRelease, *v1.ValdOperatorReleaseList]
// 		initErr error
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
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
// 		           p:phase{},
// 		           cr:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 		           p:phase{},
// 		           cr:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 				client:  test.fields.client,
// 				cfg:     test.fields.cfg,
// 				syncer:  test.fields.syncer,
// 				vor:     test.fields.vor,
// 				initErr: test.fields.initErr,
// 			}
//
// 			err := r.reconcilePhase(test.args.ctx, test.args.p, test.args.cr)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_reconciler_waitingRequeueInterval(t *testing.T) {
// 	type fields struct {
// 		client  k8s.Client
// 		cfg     *config.Config
// 		syncer  *resource.Syncer
// 		vor     *resource.Client[*v1.ValdOperatorRelease, *v1.ValdOperatorReleaseList]
// 		initErr error
// 	}
// 	type want struct {
// 		want time.Duration
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, time.Duration) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got time.Duration) error {
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
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 				client:  test.fields.client,
// 				cfg:     test.fields.cfg,
// 				syncer:  test.fields.syncer,
// 				vor:     test.fields.vor,
// 				initErr: test.fields.initErr,
// 			}
//
// 			got := r.waitingRequeueInterval()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_reconciler_syncPhase(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		p   *phase
// 		cr  *v1.ValdOperatorRelease
// 	}
// 	type fields struct {
// 		client  k8s.Client
// 		cfg     *config.Config
// 		syncer  *resource.Syncer
// 		vor     *resource.Client[*v1.ValdOperatorRelease, *v1.ValdOperatorReleaseList]
// 		initErr error
// 	}
// 	type want struct {
// 		want resource.SyncResults
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, resource.SyncResults, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got resource.SyncResults, err error) error {
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
// 		           p:phase{},
// 		           cr:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 		           p:phase{},
// 		           cr:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 				client:  test.fields.client,
// 				cfg:     test.fields.cfg,
// 				syncer:  test.fields.syncer,
// 				vor:     test.fields.vor,
// 				initErr: test.fields.initErr,
// 			}
//
// 			got, err := r.syncPhase(test.args.ctx, test.args.p, test.args.cr)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_reconciler_advance(t *testing.T) {
// 	type args struct {
// 		ps      phases
// 		current int
// 		cr      *v1.ValdOperatorRelease
// 	}
// 	type fields struct {
// 		client  k8s.Client
// 		cfg     *config.Config
// 		syncer  *resource.Syncer
// 		vor     *resource.Client[*v1.ValdOperatorRelease, *v1.ValdOperatorReleaseList]
// 		initErr error
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
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
// 		           ps:nil,
// 		           current:0,
// 		           cr:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 		           ps:nil,
// 		           current:0,
// 		           cr:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 				client:  test.fields.client,
// 				cfg:     test.fields.cfg,
// 				syncer:  test.fields.syncer,
// 				vor:     test.fields.vor,
// 				initErr: test.fields.initErr,
// 			}
//
// 			got := r.advance(test.args.ps, test.args.current, test.args.cr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_reconciler_phases(t *testing.T) {
// 	type args struct {
// 		cr         *v1.ValdOperatorRelease
// 		capability nodePoolCapability
// 	}
// 	type fields struct {
// 		client  k8s.Client
// 		cfg     *config.Config
// 		syncer  *resource.Syncer
// 		vor     *resource.Client[*v1.ValdOperatorRelease, *v1.ValdOperatorReleaseList]
// 		initErr error
// 	}
// 	type want struct {
// 		want phases
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, phases) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got phases) error {
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
// 		           cr:nil,
// 		           capability:nodePoolCapability{},
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 		           cr:nil,
// 		           capability:nodePoolCapability{},
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 				client:  test.fields.client,
// 				cfg:     test.fields.cfg,
// 				syncer:  test.fields.syncer,
// 				vor:     test.fields.vor,
// 				initErr: test.fields.initErr,
// 			}
//
// 			got := r.phases(test.args.cr, test.args.capability)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_checkClusters(t *testing.T) {
// 	type args struct {
// 		cr *v1.ValdOperatorRelease
// 	}
// 	type want struct {
// 		want result
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, result) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got result) error {
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
// 		           cr:nil,
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
// 		           cr:nil,
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
// 			got := checkClusters(test.args.cr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_reconciler_vrsReady(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		built []k8s.Object
// 	}
// 	type fields struct {
// 		client  k8s.Client
// 		cfg     *config.Config
// 		syncer  *resource.Syncer
// 		vor     *resource.Client[*v1.ValdOperatorRelease, *v1.ValdOperatorReleaseList]
// 		initErr error
// 	}
// 	type want struct {
// 		want result
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, result) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got result) error {
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
// 		           built:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 		           built:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           cfg:nil,
// 		           syncer:nil,
// 		           vor:nil,
// 		           initErr:nil,
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
// 				client:  test.fields.client,
// 				cfg:     test.fields.cfg,
// 				syncer:  test.fields.syncer,
// 				vor:     test.fields.vor,
// 				initErr: test.fields.initErr,
// 			}
//
// 			got := r.vrsReady(test.args.ctx, test.args.built)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
