//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
package replicaset

import (
	"context"
	"reflect"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
)

// ReplicaSetWatcher is a type alias for k8s resource controller.
type ReplicaSetWatcher k8s.ResourceController

type reconciler struct {
	ctx                  context.Context
	mgr                  manager.Manager
	name                 string
	namespace            string
	onError              func(err error)
	onReconcile          func(rs map[string][]ReplicaSet)
	lastReconciledResult map[string][]ReplicaSet
}

// ReplicaSet is a type alias for the k8s replica set definition.
type ReplicaSet = appsv1.ReplicaSet

// New returns the ReplicaSetWatcher that implements reconciliation loop, or any error occurred.
func New(opts ...Option) (ReplicaSetWatcher, error) {
	r := &reconciler{
		lastReconciledResult: make(map[string][]ReplicaSet),
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return r, nil
}

// Reconcile implements k8s reconciliation loop to retrieve the ReplicaSet information from k8s.
func (r *reconciler) Reconcile(req reconcile.Request) (res reconcile.Result, err error) {
	rsl := new(appsv1.ReplicaSetList)

	err = r.mgr.GetClient().List(r.ctx, rsl)
	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100,
		}
		if k8serrors.IsNotFound(err) {
			log.Warn(errors.ErrK8sResourceNotFound(err))
			res.RequeueAfter = time.Second
			return res, nil
		}
		return
	}

	rs := make(map[string][]ReplicaSet)

	for _, replicaset := range rsl.Items {
		name, ok := replicaset.GetObjectMeta().GetLabels()["app"]
		if !ok {
			pns := strings.Split(replicaset.GetName(), "-")
			name = strings.Join(pns[:len(pns)-1], "-")
		}

		if _, ok := rs[name]; !ok {
			lrr, ok := r.lastReconciledResult[name]
			if !ok {
				rs[name] = make([]ReplicaSet, 0, 0)
			} else {
				rs[name] = make([]ReplicaSet, 0, len(lrr))
			}
		}

		rs[name] = append(rs[name], replicaset)
	}

	r.lastReconciledResult = rs

	if r.onReconcile != nil {
		r.onReconcile(rs)
	}

	return
}

// GetName returns the name of resource controller.
func (r *reconciler) GetName() string {
	return r.name
}

// NewReconciler returns the reconciler for the ReplicaSet.
func (r *reconciler) NewReconciler(ctx context.Context, mgr manager.Manager) reconcile.Reconciler {
	if r.ctx == nil && ctx != nil {
		r.ctx = ctx
	}
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	appsv1.AddToScheme(r.mgr.GetScheme())
	return r
}

// For returns the runtime.Object which is replica set.
func (r *reconciler) For() runtime.Object {
	return new(appsv1.ReplicaSet)
}

// Owns returns the owner of the replica set watcher.
// It will always return nil.
func (r *reconciler) Owns() runtime.Object {
	return nil
}

// Watches returns the kind of the replica set and the event handler.
// It will always return nil.
func (r *reconciler) Watches() (*source.Kind, handler.EventHandler) {
	// return &source.Kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil
}
