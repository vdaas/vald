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

package reconciler

import (
	"context"

	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// objectCore holds the type-independent state shared by every
// objectReconciler instantiation, so ObjectOption functions need no type
// parameters and call sites never spell explicit type arguments.
type objectCore struct {
	baseReconciler
	onError func(err error)
	forOpts []builder.ForOption
}

type ObjectOption func(*objectCore)

// objectReconciler is an event-type k8s.ResourceController: each reconcile
// fetches the single object named by the request and hands it to the
// callback. It generalizes the former internal/k8s/v2/pod reconciler.
type objectReconciler[T any, PT resource.Objectable[T]] struct {
	objectCore
	onReconcile func(ctx context.Context, obj PT) (k8s.Result, error)
}

// NewObjectReconciler returns an event-type ResourceController named name
// that watches a freshly constructed PT and dispatches each reconciled
// object to onReconcile — the only type-dependent input, taken positionally
// so T and PT are inferred from it.
func NewObjectReconciler[T any, PT resource.Objectable[T]](
	name string,
	onReconcile func(ctx context.Context, obj PT) (k8s.Result, error),
	opts ...ObjectOption,
) k8s.ResourceController {
	r := &objectReconciler[T, PT]{
		objectCore:  objectCore{baseReconciler: baseReconciler{name: name}},
		onReconcile: onReconcile,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&r.objectCore)
		}
	}
	return r
}

// WithObjectOnError sets the callback invoked when fetching fails with an
// error other than NotFound.
func WithObjectOnError(f func(err error)) ObjectOption {
	return func(c *objectCore) {
		c.onError = f
	}
}

func WithObjectForOptions(opts ...builder.ForOption) ObjectOption {
	return func(c *objectCore) {
		for _, opt := range opts {
			if opt != nil {
				c.forOpts = append(c.forOpts, opt)
			}
		}
	}
}

func WithObjectFieldIndex(field string, indexer func(o k8s.Object) []string) ObjectOption {
	return func(c *objectCore) {
		c.addFieldIndex(field, indexer)
	}
}

// WithObjectMaxConcurrentReconciles sets the number of concurrent reconcile
// workers requested for this controller. Non-positive values keep the
// controller-runtime default (1).
func WithObjectMaxConcurrentReconciles(n int) ObjectOption {
	return func(c *objectCore) {
		if n > 0 {
			c.maxConcurrent = n
		}
	}
}

// Reconcile fetches the requested object and dispatches it to the callback.
// NotFound errors are ignored so that deleted objects do not requeue.
func (r *objectReconciler[T, PT]) Reconcile(
	ctx context.Context, req reconcile.Request,
) (reconcile.Result, error) {
	if err := r.checkReady(); err != nil {
		return reconcile.Result{}, err
	}
	obj, err := resource.GetObject(ctx, r.client, req.Name, req.Namespace, PT(new(T)))
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

// NewReconciler registers the scheme and the configured field indexes to the
// manager and returns the reconciler itself. Registration failures are
// recorded and surfaced by Reconcile because this method cannot return an
// error.
func (r *objectReconciler[T, PT]) NewReconciler(
	ctx context.Context, mgr manager.Manager,
) reconcile.Reconciler {
	r.setup(ctx, mgr, nil, PT(new(T)))
	return r
}

// For returns the watched object; the event-type reconciler installs its
// watch through For() so each object maps to its own request.
func (r *objectReconciler[T, PT]) For() (k8s.Object, []builder.ForOption) {
	return PT(new(T)), r.forOpts
}
