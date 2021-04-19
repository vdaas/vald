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
package configmap

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ConfigMapWatcher is a type alias for k8s resource controller
type ConfigMapWatcher k8s.ResourceController

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespaces  []string
	onError     func(err error)
	onReconcile func(rs map[string][]ConfigMap) // map[namespace][]configmap
	listOpts    []client.ListOption
	nsConfigmapsPool        sync.Pool
}

// ConfigMap is a type alias for the k8s configmap definition.
type ConfigMap = corev1.ConfigMap

// New returns the ConfigMapWather that implements reconciliation loop, or any error occured.
func New(opts ...Option) (ConfigMapWatcher, error) {
	r := &reconciler{
		nsConfigmapsPool: sync.Pool{
			New: func() interface{} {
				return make(map[string][]ConfigMap)
			},
		},
	}

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	if len(r.namespaces) != 0 {
		r.listOpts = make([]client.ListOption, 0, len(r.namespaces))
		for _, ns := range r.namespaces {
			r.listOpts = append(r.listOpts, client.InNamespace(ns))
		}
	}

	return r, nil
}

// Reconcile implements k8s reconciliation loop to retrive the ConfigMap information from k8s.
func (r *reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	cml := new(corev1.ConfigMapList)

	// TODO: add option for config map name.

	err = r.mgr.GetClient().List(ctx, cml, r.listOpts...)
	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100,
		}
		if k8serrors.IsNotFound(err) {
			log.Errorf("not found: %v", err)
			res.RequeueAfter = time.Second
			return res, nil
		}
		return
	}

	cmm := r.nsConfigmapsPool.Get().(map[string][]ConfigMap)

	for _, configmap := range cml.Items {
		if _, ok := cmm[configmap.Namespace]; !ok {
			cmm[configmap.Namespace] = make([]ConfigMap, 0)
		}
		cmm[configmap.Namespace] = append(cmm[configmap.Namespace], configmap)
	}

	if r.onReconcile != nil {
		r.onReconcile(cmm)
	}

	for name := range cmm {
		cmm[name] = cmm[name][:0:len(cmm[name])]
	}
	r.nsConfigmapsPool.Put(cmm)

	return
}

// GetName returns the name of resource controller.
func (r *reconciler) GetName() string {
	return r.name
}

// NewReconciler returns the reconciler for the ConfigMap.
func (r *reconciler) NewReconciler(mgr manager.Manager) reconcile.Reconciler {
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}

	corev1.AddToScheme(r.mgr.GetScheme())
	return r
}

// For returns the runtime.Object which is ConfigMap.
func (r *reconciler) For() (client.Object, []builder.ForOption) {
	return new(corev1.ConfigMap), nil
}

// Owns returns the owner of the ConfigMap wathcer.
// It will always return nil.
func (r *reconciler) Owns() (client.Object, []builder.OwnsOption) {
	return nil, nil
}

// Watches returns the kind of the ConfigMap and the event handler.
// It will always retrun nil.
func (r *reconciler) Watches() (*source.Kind, handler.EventHandler, []builder.WatchesOption) {
	// return &source.kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil, nil
}
