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

	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	defaultErrorRequeueDuration    = 100 * time.Millisecond
	defaultNotFoundRequeueDuration = time.Second
)

// listCore holds the type-independent state shared by every listReconciler
// instantiation, so ListOption functions need no type parameters and call
// sites never spell explicit type arguments.
type listCore struct {
	baseReconciler
	onError         func(err error)
	addToScheme     func(s *runtime.Scheme) error
	obj             k8s.Object
	lopts           []k8s.ListOption
	successRequeue  time.Duration
	notFoundRequeue time.Duration
	errorRequeue    time.Duration
}

type ListOption func(*listCore)

// listReconciler is a batch-type k8s.ResourceController: every reconcile it
// lists all objects matching the configured options and hands the whole list
// to the callback. It generalizes the former internal/k8s/{pod,node,service,
// job,metrics} reconcilers.
type listReconciler[L any, PL resource.ListPtr[L]] struct {
	listCore
	onReconcile func(ctx context.Context, list PL)
}

// NewListReconciler returns a batch-type ResourceController named name that
// watches obj and, on every reconcile, lists objects into a freshly
// constructed PL and passes it to onReconcile — the only type-dependent
// input, taken positionally so L and PL are inferred from it.
func NewListReconciler[L any, PL resource.ListPtr[L]](
	name string, obj k8s.Object, onReconcile func(ctx context.Context, list PL), opts ...ListOption,
) k8s.ResourceController {
	r := &listReconciler[L, PL]{
		listCore: listCore{
			baseReconciler:  baseReconciler{name: name},
			obj:             obj,
			errorRequeue:    defaultErrorRequeueDuration,
			notFoundRequeue: defaultNotFoundRequeueDuration,
		},
		onReconcile: onReconcile,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&r.listCore)
		}
	}
	return r
}

func WithOnError(f func(err error)) ListOption {
	return func(c *listCore) {
		c.onError = f
	}
}

func WithNamespace(ns string) ListOption {
	return func(c *listCore) {
		if ns != "" {
			c.lopts = append(c.lopts, k8s.InNamespace(ns))
		}
	}
}

func WithLabels(ls map[string]string) ListOption {
	return func(c *listCore) {
		if len(ls) > 0 {
			c.lopts = append(c.lopts, k8s.MatchingLabels(ls))
		}
	}
}

// WithFields restricts the list to objects matching the given field selectors.
// The referenced fields must be indexed via WithFieldIndex.
func WithFields(fs map[string]string) ListOption {
	return func(c *listCore) {
		if len(fs) > 0 {
			c.lopts = append(c.lopts, k8s.MatchingFields(fs))
		}
	}
}

// WithAddToScheme sets the scheme registration function executed in
// NewReconciler. When unset, the client-go native scheme is registered.
func WithAddToScheme(f func(s *runtime.Scheme) error) ListOption {
	return func(c *listCore) {
		c.addToScheme = f
	}
}

// WithFieldIndex registers a cache index for the given field so that
// WithFields selectors on that field work against the informer cache.
func WithFieldIndex(field string, indexer func(o k8s.Object) []string) ListOption {
	return func(c *listCore) {
		c.addFieldIndex(field, indexer)
	}
}

// WithRequeueDurations adjusts the requeue intervals: success is the periodic
// requeue after a successful reconcile (0 disables it), notFound and onError
// are the retry intervals for NotFound and other list errors.
func WithRequeueDurations(success, notFound, onError time.Duration) ListOption {
	return func(c *listCore) {
		c.successRequeue = success
		if notFound > 0 {
			c.notFoundRequeue = notFound
		}
		if onError > 0 {
			c.errorRequeue = onError
		}
	}
}

// WithMaxConcurrentReconciles sets the number of concurrent reconcile workers
// requested for this controller. Non-positive values keep the
// controller-runtime default (1).
func WithMaxConcurrentReconciles(n int) ListOption {
	return func(c *listCore) {
		if n > 0 {
			c.maxConcurrent = n
		}
	}
}

func (r *listReconciler[L, PL]) Reconcile(
	ctx context.Context, _ reconcile.Request,
) (reconcile.Result, error) {
	if err := r.checkReady(); err != nil {
		return reconcile.Result{}, err
	}
	list, err := resource.ListObjects(ctx, r.client, PL(new(L)), r.lopts...)
	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		if apierrors.IsNotFound(err) {
			log.Errorf("not found: %s", err)
			return reconcile.Result{RequeueAfter: r.notFoundRequeue}, nil
		}
		// controller-runtime ignores the Result when a non-nil error is
		// returned, so consume the error here to make the configured
		// errorRequeue interval effective.
		if r.errorRequeue > 0 {
			return reconcile.Result{RequeueAfter: r.errorRequeue}, nil
		}
		return reconcile.Result{}, err
	}

	if r.onReconcile != nil {
		r.onReconcile(ctx, list)
	}

	if r.successRequeue > 0 {
		return reconcile.Result{RequeueAfter: r.successRequeue}, nil
	}
	return reconcile.Result{}, nil
}

// NewReconciler registers the scheme and the configured field indexes to the
// manager and returns the reconciler itself. Registration failures are
// recorded and surfaced by Reconcile because this method cannot return an
// error.
func (r *listReconciler[L, PL]) NewReconciler(
	ctx context.Context, mgr manager.Manager,
) reconcile.Reconciler {
	r.setup(ctx, mgr, r.addToScheme, r.obj)
	return r
}

// Watches returns the watched object together with an event handler that maps
// every event to a single fixed request; the batch reconciler installs its
// watch here instead of For(), so every event maps to one fixed request
// instead of one request per object. Reconcile ignores the request and lists
// all objects anyway, so per-object requests only multiply the number of full
// List reconciles (O(N^2) work per event wave); the fixed request lets the
// workqueue deduplicate concurrent events into one reconcile.
func (r *listReconciler[L, PL]) Watches() (
	k8s.Object,
	handler.EventHandler,
	[]builder.WatchesOption,
) {
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: r.name}}
	return r.obj, handler.EnqueueRequestsFromMapFunc(func(context.Context, k8s.Object) []reconcile.Request {
		return []reconcile.Request{req}
	}), nil
}
