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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ObjectOption configures the event-type object reconciler.
type ObjectOption[T k8s.Object] func(*objectReconciler[T])

// objectReconciler is an event-type k8s.ResourceController: each reconcile
// fetches the single object named by the request and hands it to the
// callback. It generalizes the former internal/k8s/v2/pod reconciler.
type objectReconciler[T k8s.Object] struct {
	mgr           manager.Manager
	newObj        func() T
	onReconcile   func(ctx context.Context, obj T) (k8s.Result, error)
	onError       func(err error)
	addToScheme   func(s *runtime.Scheme) error
	fieldIndexes  map[string]func(o k8s.Object) []string
	name          string
	forOpts       []builder.ForOption
	maxConcurrent int
	// initErr records a NewReconciler initialization failure (missing
	// manager, scheme or index registration) so Reconcile can surface it
	// instead of silently continuing with a broken setup.
	initErr error
}

// NewObjectReconciler returns an event-type ResourceController named name
// that watches the object created by newObj and dispatches each reconciled
// object to the WithOnObjectReconcile callback.
func NewObjectReconciler[T k8s.Object](
	name string, newObj func() T, opts ...ObjectOption[T],
) k8s.ResourceController {
	r := &objectReconciler[T]{
		name:   name,
		newObj: newObj,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	return r
}

// WithOnObjectReconcile sets the callback invoked with the fetched object.
func WithOnObjectReconcile[T k8s.Object](
	f func(ctx context.Context, obj T) (k8s.Result, error),
) ObjectOption[T] {
	return func(r *objectReconciler[T]) {
		r.onReconcile = f
	}
}

// WithObjectOnError sets the callback invoked when fetching fails with an
// error other than NotFound.
func WithObjectOnError[T k8s.Object](f func(err error)) ObjectOption[T] {
	return func(r *objectReconciler[T]) {
		r.onError = f
	}
}

// WithObjectForOptions appends builder options applied to For().
func WithObjectForOptions[T k8s.Object](opts ...builder.ForOption) ObjectOption[T] {
	return func(r *objectReconciler[T]) {
		for _, opt := range opts {
			if opt != nil {
				r.forOpts = append(r.forOpts, opt)
			}
		}
	}
}

// WithObjectAddToScheme sets the scheme registration function executed in
// NewReconciler. When unset, the client-go native scheme is registered.
func WithObjectAddToScheme[T k8s.Object](f func(s *runtime.Scheme) error) ObjectOption[T] {
	return func(r *objectReconciler[T]) {
		r.addToScheme = f
	}
}

// WithObjectFieldIndex registers a cache index for the given field.
func WithObjectFieldIndex[T k8s.Object](
	field string, indexer func(o k8s.Object) []string,
) ObjectOption[T] {
	return func(r *objectReconciler[T]) {
		if field == "" || indexer == nil {
			return
		}
		if r.fieldIndexes == nil {
			r.fieldIndexes = make(map[string]func(o k8s.Object) []string, 1)
		}
		r.fieldIndexes[field] = indexer
	}
}

// WithObjectMaxConcurrentReconciles sets the number of concurrent reconcile
// workers requested for this controller. Non-positive values keep the
// controller-runtime default (1).
func WithObjectMaxConcurrentReconciles[T k8s.Object](n int) ObjectOption[T] {
	return func(r *objectReconciler[T]) {
		if n > 0 {
			r.maxConcurrent = n
		}
	}
}

// Reconcile fetches the requested object and dispatches it to the callback.
// NotFound errors are ignored so that deleted objects do not requeue.
func (r *objectReconciler[T]) Reconcile(
	ctx context.Context, req reconcile.Request,
) (reconcile.Result, error) {
	if r.initErr != nil {
		return reconcile.Result{}, r.initErr
	}
	if r.mgr == nil {
		return reconcile.Result{}, errors.Errorf("manager is not registered for %s reconciler", r.name)
	}
	obj, err := resource.GetObject(ctx, r.mgr.GetClient(), req.Name, req.Namespace, r.newObj())
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		if r.onError != nil {
			r.onError(err)
		}
		return reconcile.Result{}, err
	}
	if r.onReconcile != nil {
		return r.onReconcile(ctx, obj)
	}
	return reconcile.Result{}, nil
}

// GetName returns the name of resource controller.
func (r *objectReconciler[T]) GetName() string {
	return r.name
}

// MaxConcurrentReconciles implements the optional k8s.ConcurrentReconciler
// interface: values greater than zero request that many reconcile workers.
func (r *objectReconciler[T]) MaxConcurrentReconciles() int {
	return r.maxConcurrent
}

// NewReconciler registers the scheme and the configured field indexes to the
// manager and returns the reconciler itself. Registration failures are
// recorded and surfaced by Reconcile because this method cannot return an
// error.
func (r *objectReconciler[T]) NewReconciler(
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
		if err := r.mgr.GetFieldIndexer().IndexField(ctx, r.newObj(), field, indexer); err != nil {
			r.initErr = errors.Wrapf(err, "failed to register field index %s for %s reconciler", field, r.name)
			log.Error(r.initErr)
			return r
		}
	}
	return r
}

// For returns the watched runtime object.
func (r *objectReconciler[T]) For() (k8s.Object, []builder.ForOption) {
	return r.newObj(), r.forOpts
}

// Owns returns the owned object. It will always return nil.
func (*objectReconciler[T]) Owns() (k8s.Object, []builder.OwnsOption) {
	return nil, nil
}

// Watches returns the watched object and the event handler.
// It will always return nil.
func (*objectReconciler[T]) Watches() (k8s.Object, handler.EventHandler, []builder.WatchesOption) {
	return nil, nil, nil
}
