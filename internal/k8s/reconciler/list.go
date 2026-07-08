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
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	defaultErrorRequeueDuration    = 100 * time.Millisecond
	defaultNotFoundRequeueDuration = time.Second
)

// ListOption configures the batch-type list reconciler.
type ListOption[L k8s.ObjectList] func(*listReconciler[L])

// listReconciler is a batch-type k8s.ResourceController: every reconcile it
// lists all objects matching the configured options and hands the whole list
// to the callback. It generalizes the former internal/k8s/{pod,node,service,
// job,metrics} reconcilers.
type listReconciler[L k8s.ObjectList] struct {
	mgr             manager.Manager
	newList         func() L
	onReconcile     func(ctx context.Context, list L)
	onError         func(err error)
	addToScheme     func(s *runtime.Scheme) error
	fieldIndexes    map[string]func(o k8s.Object) []string
	obj             k8s.Object
	name            string
	lopts           []k8s.ListOption
	forOpts         []builder.ForOption
	successRequeue  time.Duration
	notFoundRequeue time.Duration
	errorRequeue    time.Duration
	maxConcurrent   int
	// initErr records a NewReconciler initialization failure (missing
	// manager, scheme or index registration) so Reconcile can surface it
	// instead of silently continuing with a broken setup.
	initErr error
}

// NewListReconciler returns a batch-type ResourceController named name that
// watches obj and, on every reconcile, lists objects into the list created by
// newList and passes it to the WithOnReconcile callback.
func NewListReconciler[L k8s.ObjectList](
	name string, obj k8s.Object, newList func() L, opts ...ListOption[L],
) k8s.ResourceController {
	r := &listReconciler[L]{
		name:            name,
		obj:             obj,
		newList:         newList,
		errorRequeue:    defaultErrorRequeueDuration,
		notFoundRequeue: defaultNotFoundRequeueDuration,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	return r
}

// WithOnReconcile sets the callback invoked with the freshly listed objects.
func WithOnReconcile[L k8s.ObjectList](f func(ctx context.Context, list L)) ListOption[L] {
	return func(r *listReconciler[L]) {
		r.onReconcile = f
	}
}

// WithOnError sets the callback invoked when listing fails.
func WithOnError[L k8s.ObjectList](f func(err error)) ListOption[L] {
	return func(r *listReconciler[L]) {
		r.onError = f
	}
}

// WithNamespace restricts the list to the given namespace when non-empty.
func WithNamespace[L k8s.ObjectList](ns string) ListOption[L] {
	return func(r *listReconciler[L]) {
		if ns != "" {
			r.lopts = append(r.lopts, k8s.InNamespace(ns))
		}
	}
}

// WithLabels restricts the list to objects matching the given labels.
func WithLabels[L k8s.ObjectList](ls map[string]string) ListOption[L] {
	return func(r *listReconciler[L]) {
		if len(ls) > 0 {
			r.lopts = append(r.lopts, k8s.MatchingLabels(ls))
		}
	}
}

// WithFields restricts the list to objects matching the given field selectors.
// The referenced fields must be indexed via WithFieldIndex.
func WithFields[L k8s.ObjectList](fs map[string]string) ListOption[L] {
	return func(r *listReconciler[L]) {
		if len(fs) > 0 {
			r.lopts = append(r.lopts, k8s.MatchingFields(fs))
		}
	}
}

// WithListOptions appends raw list options.
func WithListOptions[L k8s.ObjectList](opts ...k8s.ListOption) ListOption[L] {
	return func(r *listReconciler[L]) {
		for _, opt := range opts {
			if opt != nil {
				r.lopts = append(r.lopts, opt)
			}
		}
	}
}

// WithForOptions appends builder options applied to the watch. Options that
// also implement builder.WatchesOption (e.g. builder.WithPredicates) carry
// over; others are ignored because the batch reconciler watches through
// Watches() rather than For().
func WithForOptions[L k8s.ObjectList](opts ...builder.ForOption) ListOption[L] {
	return func(r *listReconciler[L]) {
		for _, opt := range opts {
			if opt != nil {
				r.forOpts = append(r.forOpts, opt)
			}
		}
	}
}

// WithAddToScheme sets the scheme registration function executed in
// NewReconciler. When unset, the client-go native scheme is registered.
func WithAddToScheme[L k8s.ObjectList](f func(s *runtime.Scheme) error) ListOption[L] {
	return func(r *listReconciler[L]) {
		r.addToScheme = f
	}
}

// WithFieldIndex registers a cache index for the given field so that
// WithFields selectors on that field work against the informer cache.
func WithFieldIndex[L k8s.ObjectList](
	field string, indexer func(o k8s.Object) []string,
) ListOption[L] {
	return func(r *listReconciler[L]) {
		if field == "" || indexer == nil {
			return
		}
		if r.fieldIndexes == nil {
			r.fieldIndexes = make(map[string]func(o k8s.Object) []string, 1)
		}
		r.fieldIndexes[field] = indexer
	}
}

// WithRequeueDurations adjusts the requeue intervals: success is the periodic
// requeue after a successful reconcile (0 disables it), notFound and onError
// are the retry intervals for NotFound and other list errors.
func WithRequeueDurations[L k8s.ObjectList](
	success, notFound, onError time.Duration,
) ListOption[L] {
	return func(r *listReconciler[L]) {
		r.successRequeue = success
		if notFound > 0 {
			r.notFoundRequeue = notFound
		}
		if onError > 0 {
			r.errorRequeue = onError
		}
	}
}

// WithMaxConcurrentReconciles sets the number of concurrent reconcile workers
// requested for this controller. Non-positive values keep the
// controller-runtime default (1).
func WithMaxConcurrentReconciles[L k8s.ObjectList](n int) ListOption[L] {
	return func(r *listReconciler[L]) {
		if n > 0 {
			r.maxConcurrent = n
		}
	}
}

// Reconcile lists the configured objects and dispatches them to the callback.
func (r *listReconciler[L]) Reconcile(
	ctx context.Context, _ reconcile.Request,
) (reconcile.Result, error) {
	if r.initErr != nil {
		return reconcile.Result{}, r.initErr
	}
	if r.mgr == nil {
		return reconcile.Result{}, errors.Errorf("manager is not registered for %s reconciler", r.name)
	}
	list, err := resource.ListObjects(ctx, r.mgr.GetClient(), r.newList(), r.lopts...)
	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		if apierrors.IsNotFound(err) {
			log.Errorf("not found: %s", err)
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: r.notFoundRequeue,
			}, nil
		}
		// controller-runtime ignores the Result when a non-nil error is
		// returned, so consume the error here to make the configured
		// errorRequeue interval effective.
		if r.errorRequeue > 0 {
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: r.errorRequeue,
			}, nil
		}
		return reconcile.Result{}, err
	}

	if r.onReconcile != nil {
		r.onReconcile(ctx, list)
	}

	if r.successRequeue > 0 {
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: r.successRequeue,
		}, nil
	}
	return reconcile.Result{}, nil
}

// GetName returns the name of resource controller.
func (r *listReconciler[L]) GetName() string {
	return r.name
}

// MaxConcurrentReconciles implements the optional k8s.ConcurrentReconciler
// interface: values greater than zero request that many reconcile workers.
func (r *listReconciler[L]) MaxConcurrentReconciles() int {
	return r.maxConcurrent
}

// NewReconciler registers the scheme and the configured field indexes to the
// manager and returns the reconciler itself. Registration failures are
// recorded and surfaced by Reconcile because this method cannot return an
// error.
func (r *listReconciler[L]) NewReconciler(
	ctx context.Context, mgr manager.Manager,
) reconcile.Reconciler {
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	if r.mgr == nil {
		r.initErr = errors.Errorf("manager is not registered for %s reconciler", r.name)
		log.Error(r.initErr)
		return r
	}
	ats := r.addToScheme
	if ats == nil {
		ats = clientgoscheme.AddToScheme
	}
	if err := ats(r.mgr.GetScheme()); err != nil {
		r.initErr = errors.Wrapf(err, "failed to register scheme for %s reconciler", r.name)
		log.Error(r.initErr)
		return r
	}
	for field, indexer := range r.fieldIndexes {
		if err := r.mgr.GetFieldIndexer().IndexField(ctx, r.obj, field, indexer); err != nil {
			r.initErr = errors.Wrapf(err, "failed to register field index %s for %s reconciler", field, r.name)
			log.Error(r.initErr)
			return r
		}
	}
	return r
}

// For returns nil: the batch reconciler installs its watch through Watches()
// so every event maps to one fixed request instead of one request per object.
func (*listReconciler[L]) For() (k8s.Object, []builder.ForOption) {
	return nil, nil
}

// Owns returns the owned object. It will always return nil.
func (*listReconciler[L]) Owns() (k8s.Object, []builder.OwnsOption) {
	return nil, nil
}

// Watches returns the watched object together with an event handler that maps
// every event to a single fixed request. Reconcile ignores the request and
// lists all objects anyway, so per-object requests only multiply the number
// of full List reconciles (O(N^2) work per event wave); the fixed request
// lets the workqueue deduplicate concurrent events into one reconcile.
func (r *listReconciler[L]) Watches() (k8s.Object, handler.EventHandler, []builder.WatchesOption) {
	var wopts []builder.WatchesOption
	for _, opt := range r.forOpts {
		if wopt, ok := opt.(builder.WatchesOption); ok {
			wopts = append(wopts, wopt)
		}
	}
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: r.name}}
	return r.obj, handler.EnqueueRequestsFromMapFunc(func(context.Context, k8s.Object) []reconcile.Request {
		return []reconcile.Request{req}
	}), wopts
}
