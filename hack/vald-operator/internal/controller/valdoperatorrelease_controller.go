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

package controller

import (
	"context"
	"fmt"

	"github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/config"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/desired"

	"k8s.io/apimachinery/pkg/api/equality"

	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/domain/valdoperatorrelease"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle"
	builder "github.com/vdaas/vald/hack/vald-operator/internal/pkg/lifecycle/builder/vald"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	controllerv1 "github.com/vdaas/vald/hack/vald-operator/api/v1"
)

// ValdOperatorReleaseReconciler reconciles a ValdOperatorRelease object. The reconciler is
// the orchestrator: it fetches the CR, builds the Domain and lifecycle flow,
// and delegates the actual write side (CreateOrUpdate / prune / key) to
// ResourceSyncer.
type ValdOperatorReleaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Config *config.Config
	Syncer *ResourceSyncer
}

// +kubebuilder:rbac:groups=vald.vdaas.org,resources=valdoperatorreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vald.vdaas.org,resources=valdoperatorreleases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vald.vdaas.org,resources=valdoperatorreleases/finalizers,verbs=update

// +kubebuilder:rbac:groups=vald.vdaas.org,resources=valdreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vald.vdaas.org,resources=valdreleases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vald.vdaas.org,resources=valdreleases/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch

// Reconcile implements the main Kubernetes reconciliation loop for ValdOperatorRelease.
// It fetches the object and delegates to reconcileRoutine.
func (r *ValdOperatorReleaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	vor := &controllerv1.ValdOperatorRelease{}
	logger.Info("Reconcile", "type", "ValdOperatorRelease")

	if err := r.Get(ctx, req.NamespacedName, vor); err != nil {
		logger.Info("Notice: unable to fetch", "type", "ValdOperatorRelease", "name", req.Name, "namespace", req.Namespace, "error", err.Error())
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("Fetched", "type", "ValdOperatorRelease", "name", vor.Name, "namespace", vor.Namespace)
	d := &valdoperatorrelease.Domain{ValdOperatorRelease: vor}

	if err := r.reconcileRoutine(ctx, d); err != nil {
		logger.Error(err, err.Error())
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileRoutine reconciles ValdOperatorRelease by iterating all accumulated conditions on every
// reconcile. This ensures self-healing: if a previously-True condition breaks (e.g. a VRS
// resource is deleted or the spec is modified), it is detected and processing restarts from
// that point.
func (r *ValdOperatorReleaseReconciler) reconcileRoutine(ctx context.Context, d *valdoperatorrelease.Domain) (err error) {
	originalStatus := d.Status.DeepCopy()
	defer func() {
		if !equality.Semantic.DeepEqual(originalStatus, d.Status) {
			if uerr := r.Status().Update(ctx, d.ValdOperatorRelease); uerr != nil {
				err = uerr
			}
		}
	}()

	logger := logf.FromContext(ctx)

	// Resolve node-pool availability up front so the VrsBuilder can be a pure
	// function of (CR, Config, Capability). Skip the node listing entirely
	// when REQUIRE_NODEPOOL_MATCH is disabled — the Builder will then treat
	// every infra entry as schedulable.
	capability := builder.AlwaysAvailable()
	if r.Config.RequireNodePoolMatch {
		capability, err = builder.ResolveNodePoolCapability(ctx, r.Client, d.Namespace, r.Config.NodePoolLabelPrefix)
		if err != nil {
			return fmt.Errorf("failed to resolve node pool capability: %w", err)
		}
	}

	lc := d.NewLifeCycleFlow(r.Client, r.Config, capability)
	d.InitProgress(lc.LifeCycles)

	if len(d.Status.Conditions) == 0 {
		first := &lc.LifeCycles[lc.Current]
		d.Status.Phase = first.Condition.Type
		util.UpdateStatus(&d.Status.Conditions, first.Condition.MakeCondition(desired.Progressing("")))
		return nil
	}

	for _, cond := range d.Status.Conditions {
		d.Status.Phase = cond.Type

		ri := lc.LifeCycles.GetIndex(cond.Type)
		if ri < 0 {
			logger.Info("Skipping unrecognized condition", "condition", cond.Type, "name", d.Name, "namespace", d.Namespace)
			continue
		}

		cf := &lc.LifeCycles[ri]
		if cf.Builder == nil && cf.Checker == nil {
			// Terminal phase (e.g. Completed): nothing to build, nothing to check.
			continue
		}

		if err := r.reconcileCondition(ctx, cf, d); err != nil {
			return err
		}
	}

	if next := d.AdvanceToNextPhase(lc); next >= 0 {
		logger.Info("Advanced to next phase", "phase", lc.LifeCycles[next].Condition.Type, "name", d.Name, "namespace", d.Namespace)
	}

	return nil
}

func (r *ValdOperatorReleaseReconciler) reconcileCondition(ctx context.Context, cf *lifecycle.LifeCycle, d *valdoperatorrelease.Domain) error {
	if cf.Builder != nil {
		ope, err := r.Syncer.Sync(ctx, cf.Builder, d.ValdOperatorRelease)
		if err != nil {
			err = fmt.Errorf("failed to sync resources for condition %s: %w", cf.Condition.Type, err)
			util.UpdateStatus(&d.Status.Conditions, cf.Condition.MakeCondition(desired.Failed(err)))
			return err
		}
		logf.FromContext(ctx).Info("Synced resources for condition",
			"condition", cf.Condition.Type, "name", d.Name, "namespace", d.Namespace, "operation", ope)
	}

	if cf.Checker == nil {
		return nil
	}
	result := cf.Checker.IsReady(ctx)
	cn := cf.Condition.MakeCondition(result)
	util.UpdateStatus(&d.Status.Conditions, cn)
	if cn.Status != metav1.ConditionTrue {
		return fmt.Errorf("condition %s is not ready: %s", cn.Type, cn.Message)
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ValdOperatorReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Syncer == nil {
		r.Syncer = NewResourceSyncer(r.Client, r.Scheme)
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&controllerv1.ValdOperatorRelease{}).
		Owns(&valdrelease.ValdRelease{}).
		Named("valdoperatorrelease").
		Complete(r)
}
