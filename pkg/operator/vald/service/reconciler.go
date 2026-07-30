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

package service

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/resource"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/pkg/operator/vald/config"
)

const name = "valdoperatorrelease"

// defaultRequeueAfterWaiting is the requeue interval for phases that are
// waiting on a normal transition (Progressing/Pending conditions) when
// requeue.success is not configured. It is intentionally longer than the
// chart-default requeue.on_error (100ms), which is meant for genuine failures.
const defaultRequeueAfterWaiting = 3 * time.Second

// errPhaseWaiting marks reconcile pipeline stops that are normal waits
// (Progressing/Pending conditions, e.g. VRS objects not ready yet) rather
// than failures. Reconcile requeues them at the waiting interval instead of
// the on_error interval.
var errPhaseWaiting = errors.New("phase is waiting")

// resourceController wires the ValdOperatorRelease reconciler into the shared
// controller manager as a k8s.ResourceController.
type resourceController struct {
	cfg *config.Config
}

func newResourceController(cfg *config.Config) k8s.ResourceController {
	return &resourceController{cfg: cfg}
}

func (*resourceController) GetName() string {
	return name
}

// NewReconciler returns the reconciler for the ValdOperatorRelease. It registers the
// ValdOperatorRelease/ValdRelease schemes and the client-go native schemes to the
// manager's scheme before constructing the reconciler. Building the vor
// resource client should not fail in practice (mgr.GetConfig() is already
// resolved by the time the manager builds this controller), but
// k8s.ResourceController.NewReconciler has no error return, so a failure is
// recorded in initErr and surfaced by Reconcile instead of panicking here.
func (rc *resourceController) NewReconciler(_ context.Context, mgr k8s.Manager) k8s.Reconciler {
	scheme := mgr.GetScheme()
	for _, reg := range []struct {
		add  func(*k8s.Scheme) error
		name string
	}{
		{k8s.AddClientGoScheme, "client-go"},
		{v1.AddToScheme, "ValdOperatorRelease"},
		{valdrelease.AddToScheme, "ValdRelease"},
	} {
		if err := reg.add(scheme); err != nil {
			log.Errorf("failed to register %s scheme: %v", reg.name, err)
		}
	}
	r := &reconciler{
		client: mgr.GetClient(),
		cfg:    rc.cfg,
		syncer: resource.NewSyncer(mgr.GetClient(), scheme, rc.cfg.ManagedGenerationLabel),
	}
	cl, err := client.NewFromManager(mgr)
	if err != nil {
		r.initErr = errors.Wrap(err, "failed to build k8s client for ValdOperatorRelease reconciler")
		log.Error(r.initErr)
		return r
	}
	r.vor = resource.NewClient(cl, new(v1.ValdOperatorRelease), new(v1.ValdOperatorReleaseList))
	return r
}

// MaxConcurrentReconciles implements the optional k8s.ConcurrentReconciler
// interface so that the configured worker count is applied when the
// ValdOperatorRelease controller is built.
func (rc *resourceController) MaxConcurrentReconciles() int {
	return rc.cfg.MaxConcurrentReconciles
}

func (*resourceController) For() (k8s.Object, []k8s.ForOption) {
	return new(v1.ValdOperatorRelease), nil
}

func (*resourceController) Owns() (k8s.Object, []k8s.OwnsOption) {
	return new(valdrelease.ValdRelease), nil
}

// Watches returns the kind of the ValdOperatorRelease and the event handler.
// It will always return nil.
func (*resourceController) Watches() (k8s.Object, k8s.EventHandler, []k8s.WatchesOption) {
	return nil, nil, nil
}

type reconciler struct {
	client k8s.Client
	cfg    *config.Config
	syncer *resource.Syncer
	vor    *resource.Client[*v1.ValdOperatorRelease, *v1.ValdOperatorReleaseList]
	// initErr records a NewReconciler construction failure (building the vor
	// resource client) so Reconcile can surface it, since
	// k8s.ResourceController.NewReconciler has no error return.
	initErr error
}

// Reconcile implements the main Kubernetes reconciliation loop for
// ValdOperatorRelease. It fetches the object and delegates to reconcileValdOperatorRelease.
// The configured requeue intervals shape the result: not_found retries
// missing objects, on_error replaces the exponential backoff with a fixed
// interval and success enables a periodic re-reconcile. Each defaults to
// disabled (zero), which preserves the standard controller-runtime behavior.
// Waiting phases (Progressing/Pending conditions) are not treated as errors:
// they requeue at requeue.success when configured, otherwise at
// defaultRequeueAfterWaiting.
func (r *reconciler) Reconcile(ctx context.Context, req k8s.Request) (k8s.Result, error) {
	log.Debug("reconciling ValdOperatorRelease")

	if r.initErr != nil {
		return k8s.Result{}, r.initErr
	}

	cr, err := r.vor.Get(ctx, req.Name, req.Namespace)
	if err != nil {
		log.Infof("unable to fetch ValdOperatorRelease %s/%s: %v", req.Namespace, req.Name, err)
		if resource.IgnoreNotFound(err) != nil {
			return k8s.Result{}, err
		}
		if r.cfg.RequeueAfterNotFound > 0 {
			return k8s.Result{RequeueAfter: r.cfg.RequeueAfterNotFound}, nil
		}
		return k8s.Result{}, nil
	}
	log.Debugf("fetched ValdOperatorRelease %s/%s", cr.Namespace, cr.Name)

	if err := r.reconcileValdOperatorRelease(ctx, cr); err != nil {
		if errors.Is(err, errPhaseWaiting) {
			// A waiting (Progressing/Pending) phase is a normal state, not a
			// failure: requeue at the waiting interval so ready-waits do not
			// hot-loop on the short on_error interval.
			log.Debugf("ValdOperatorRelease %s/%s is waiting: %v", cr.Namespace, cr.Name, err)
			return k8s.Result{RequeueAfter: r.waitingRequeueInterval()}, nil
		}
		log.Error(err)
		if r.cfg.RequeueAfterError > 0 {
			// The error is consumed here on purpose: returning it alongside a
			// result would make controller-runtime ignore the fixed interval
			// and fall back to its rate limiter.
			return k8s.Result{RequeueAfter: r.cfg.RequeueAfterError}, nil
		}
		return k8s.Result{}, err
	}

	if r.cfg.RequeueAfterSuccess > 0 {
		return k8s.Result{RequeueAfter: r.cfg.RequeueAfterSuccess}, nil
	}
	return k8s.Result{}, nil
}

// reconcileValdOperatorRelease walks the phase pipeline by iterating all accumulated
// conditions on every reconcile. This ensures self-healing: if a
// previously-True condition breaks (e.g. a VRS resource is deleted or the
// spec is modified), it is detected and processing restarts from that point.
func (r *reconciler) reconcileValdOperatorRelease(
	ctx context.Context, cr *v1.ValdOperatorRelease,
) (err error) {
	originalStatus := cr.Status.DeepCopy()
	defer func() {
		if !resource.SemanticDeepEqual(originalStatus, &cr.Status) {
			if uerr := r.vor.UpdateStatus(ctx, cr); uerr != nil {
				err = errors.Join(err, uerr)
			}
		}
	}()

	// Resolve node-pool availability up front so the VRS builder can be a
	// pure function of (CR, Config, Capability). Skip the node listing
	// entirely when node_pool.require_match is disabled — the builder then
	// treats every infra entry as schedulable.
	capability := alwaysAvailable()
	if r.cfg.RequireNodePoolMatch {
		capability, err = resolveNodePoolCapability(ctx, r.client, cr.Namespace, r.cfg.NodePoolLabelPrefix)
		if err != nil {
			return errors.Wrap(err, "failed to resolve node pool capability")
		}
	}

	ps := r.phases(cr, capability)
	current := ps.index(cr.Status.Phase)
	// The terminal Completed phase has no work, so it is excluded from Total.
	cr.Status.Progress.Total = len(ps) - 1
	cr.Status.Progress.Completed = current

	if len(cr.Status.Conditions) == 0 {
		first := &ps[max(current, 0)]
		cr.Status.Phase = first.name
		resource.UpsertCondition(&cr.Status.Conditions, first.condition(progressing("")))
		return nil
	}

	for _, cond := range cr.Status.Conditions {
		cr.Status.Phase = cond.Type

		i := ps.index(cond.Type)
		if i < 0 {
			log.Warnf("skipping unrecognized condition %s for ValdOperatorRelease %s/%s", cond.Type, cr.Namespace, cr.Name)
			continue
		}

		p := &ps[i]
		if p.terminal() {
			continue
		}

		if err := r.reconcilePhase(ctx, p, cr); err != nil {
			return err
		}
	}

	if next := r.advance(ps, current, cr); next >= 0 {
		log.Infof("advanced ValdOperatorRelease %s/%s to next phase %s", cr.Namespace, cr.Name, ps[next].name)
	}
	return nil
}

// reconcilePhase runs a single phase: build + sync the desired objects when
// the phase produces any, then record the readiness check outcome as a
// condition. A non-True condition aborts the pipeline with an error.
func (r *reconciler) reconcilePhase(
	ctx context.Context, p *phase, cr *v1.ValdOperatorRelease,
) error {
	if p.build != nil {
		ope, err := r.syncPhase(ctx, p, cr)
		if err != nil {
			err = errors.Wrapf(err, "failed to sync resources for condition %s", p.name)
			resource.UpsertCondition(&cr.Status.Conditions, p.condition(failed(err)))
			return err
		}
		log.Infof("synced resources for condition %s of ValdOperatorRelease %s/%s: %v",
			p.name, cr.Namespace, cr.Name, ope)
	}

	if p.check == nil {
		return nil
	}
	cn := p.condition(p.check(ctx))
	resource.UpsertCondition(&cr.Status.Conditions, cn)
	if cn.Status == k8s.ConditionTrue {
		return nil
	}
	if cn.Status == k8s.ConditionUnknown {
		// Progressing/Pending: a normal wait, distinguished from failures so
		// that Reconcile can apply the waiting requeue interval.
		return errors.Wrapf(errPhaseWaiting, "condition %s is not ready: %s", cn.Type, cn.Message)
	}
	return errors.Errorf("condition %s is not ready: %s", cn.Type, cn.Message)
}

func (r *reconciler) waitingRequeueInterval() time.Duration {
	if r.cfg.RequeueAfterSuccess > 0 {
		return r.cfg.RequeueAfterSuccess
	}
	return defaultRequeueAfterWaiting
}

func (r *reconciler) syncPhase(
	ctx context.Context, p *phase, cr *v1.ValdOperatorRelease,
) (resource.SyncResults, error) {
	items, err := p.build(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build desired resource")
	}
	return r.syncer.Sync(ctx, cr, items, p.fetch)
}

func (*reconciler) advance(ps phases, current int, cr *v1.ValdOperatorRelease) int {
	next := current + 1
	if next >= len(ps) {
		return -1
	}
	p := &ps[next]
	cr.Status.Phase = p.name
	cr.Status.Progress.Completed = next

	seed := progressing("")
	if p.terminal() {
		seed = succeeded()
	}
	resource.SetStatusCondition(&cr.Status.Conditions, p.condition(seed))
	return next
}

// phases returns the ordered reconcile pipeline for the given CR. The table
// is rebuilt per reconcile because build/check close over the fetched CR and
// the resolved node-pool capability; the VRS phase shares the objects built
// in this pass between its build and check steps.
func (r *reconciler) phases(cr *v1.ValdOperatorRelease, capability nodePoolCapability) phases {
	var built []k8s.Object
	return phases{
		{
			name:    phaseWaitForClusterCreate,
			message: "Waiting for Cluster Creation.",
			check: func(context.Context) result {
				return checkClusters(cr)
			},
		},
		{
			name:    phaseWaitForCreateVrs,
			message: "Waiting for VRS creation.",
			build: func(ctx context.Context) ([]k8s.Object, error) {
				items, err := newVrsBuilder(cr, r.cfg, capability).Build(ctx)
				if err != nil {
					return nil, err
				}
				// Refuse to sync an empty render: Syncer treats an empty desired
				// set as "prune every owned object", so a pass in which the builder
				// skipped everything (all Infrastructure entries inactive, or a
				// require_match node lookup transiently matching nothing) would
				// tear down the whole running deployment. Intentional teardown goes
				// through CR deletion (owner-reference GC), not an empty render.
				if len(items) == 0 {
					return nil, errors.Wrap(errors.ErrBuiltResourceIsNil,
						"builder rendered no ValdRelease objects; refusing to sync an empty desired set")
				}
				built = items
				return items, nil
			},
			fetch: func(ctx context.Context) ([]k8s.Object, error) {
				return fetchExistingVrs(ctx, r.client, cr.Namespace)
			},
			check: func(ctx context.Context) result {
				return r.vrsReady(ctx, built)
			},
		},
		{
			name:    phaseCompleted,
			message: "VRS creation completed successfully.",
			// Terminal phase: no build and no check.
		},
	}
}

func checkClusters(cr *v1.ValdOperatorRelease) result {
	if len(cr.Spec.Infrastructure) == 0 {
		return failed(errors.ErrInfrastructureConfigurationIsMissing)
	}
	// Only Active entries gate the pipeline: vrsBuilder.Build skips inactive
	// entries entirely, so requiring provisioned clusters on an inactive entry
	// would block rendering work that never uses them. Conversely, zero active
	// entries means the builder would render nothing (and the VRS phase refuses
	// an empty render), so hold the pipeline here with a clear message instead.
	active := 0
	for _, infra := range cr.Spec.Infrastructure {
		if !infra.Active {
			continue
		}
		active++
		if len(infra.Clusters) == 0 {
			return failed(errors.ErrNoClustersDefinedInConfiguration)
		}
		for _, cluster := range infra.Clusters {
			if cluster.ID == "" || cluster.Name == "" {
				return pending("waiting for cluster to be provisioned: {id: " + cluster.ID + ", name: " + cluster.Name + "}")
			}
		}
	}
	if active == 0 {
		return pending("waiting for at least one active infrastructure entry")
	}
	return succeeded()
}

// vrsReady re-fetches every object built in this pass and reports readiness:
// a missing object keeps the phase progressing, any other fetch error fails
// it.
func (r *reconciler) vrsReady(ctx context.Context, built []k8s.Object) result {
	// An empty (not just nil) built set must not vacuously report success: the
	// range below would run zero iterations and fall through to succeeded(),
	// advancing the CR to Completed even though nothing was rendered (e.g. every
	// Infrastructure entry is inactive, or a require_match node lookup transiently
	// matched nothing). Treat "no built resources" the same as the nil case.
	if len(built) == 0 {
		return failed(errors.ErrBuiltResourceIsNil)
	}
	for _, obj := range built {
		if _, err := resource.RefreshObject(ctx, r.client, obj); err != nil {
			if resource.IgnoreNotFound(err) == nil {
				return progressing("")
			}
			return failed(err)
		}
	}
	return succeeded()
}
