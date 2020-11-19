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
package statefulset

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// StatefulSetWatcher is a type alias for k8s resource controller
type StatefulSetWatcher k8s.ResourceController

type reconciler struct {
	ctx         context.Context
	mgr         manager.Manager
	name        string
	namespace   string
	onError     func(err error)
	onReconcile func(rs map[string][]StatefulSet)
	pool        sync.Pool
}

// Statefulset is a type alias for the k8s statefulset definition.
type StatefulSet = appsv1.StatefulSet

// New returns the StatefulSetWather that implements reconciliation loop, or any error occured.
func New(opts ...Option) (StatefulSetWatcher, error) {
	r := &reconciler{
		pool: sync.Pool{
			New: func() interface{} {
				return make(map[string][]StatefulSet)
			},
		},
	}

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return r, nil
}

// Reconcile implements k8s reconciliation loop to retrive the StatefulSet information from k8s.
func (r *reconciler) Reconcile(req reconcile.Request) (res reconcile.Result, err error) {
	ssl := new(appsv1.StatefulSetList)

	err = r.mgr.GetClient().List(r.ctx, ssl)
	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100,
		}
		if k8serrors.IsNotFound(err) {
			log.Errorf("not found: %s", err)
			res.RequeueAfter = time.Second
			return res, nil
		}
		return
	}

	ssm := make(map[string][]StatefulSet)
	appList := make(map[string]bool)

	for _, statefulset := range ssl.Items {
		name, ok := statefulset.GetObjectMeta().GetLabels()["app"]
		if !ok {
			pns := strings.Split(statefulset.GetName(), "-")
			name = strings.Join(pns[:len(pns)-1], "-")
		}
		if _, ok := ssm[name]; !ok {
			ssm[name] = make([]StatefulSet, 1)
		}
		if !appList[name] {
			appList[name] = true
		}
		ssm[name] = append(ssm[name], statefulset)
	}

	if r.onReconcile != nil {
		r.onReconcile(ssm)
	}

	for name := range ssm {
		if !appList[name] {
			delete(ssm, name)
		} else {
			ssm[name] = ssm[name][:0:len(ssm[name])]
		}
	}
	r.pool.Put(ssm)

	return
}

// GetName returns the name of resource controller.
func (r *reconciler) GetName() string {
	return r.name
}

// NewReconciler returns the reconciler for the StatefulSet.
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

// For returns the runtime.Object which is StatefulSet.
func (r *reconciler) For() runtime.Object {
	return new(appsv1.StatefulSet)
}

// Owns returns the owner of the StatefulSet wathcer.
// It will always return nil.
func (r *reconciler) Owns() runtime.Object {
	return nil
}

// Watches returns the kind of the StatefulSet and the event handler.
// It will always retrun nil.
func (r *reconciler) Watches() (*source.Kind, handler.EventHandler) {
	// return &source.kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil
}
