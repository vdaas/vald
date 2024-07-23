//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package pod provides kubernetes pod information and preriodically update
package pod

import (
	"context"

	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PodWatcher k8s.ResourceController

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespace   string
	onError     func(err error)
	onReconcile func(ctx context.Context, pod *corev1.Pod) (reconcile.Result, error)
	lopts       []client.ListOption
	forOpts     []builder.ForOption
}

type Pod struct {
	Name        string
	NodeName    string
	Namespace   string
	IP          string
	CPULimit    float64
	CPURequest  float64
	MemLimit    float64
	MemRequest  float64
	Labels      map[string]string
	Annotations map[string]string
}

func New(opts ...Option) PodWatcher {
	r := new(reconciler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}
	return r
}

func (r *reconciler) addListOpts(opt client.ListOption) {
	if opt == nil {
		return
	}
	if r.lopts == nil {
		r.lopts = make([]client.ListOption, 0, 1)
	}
	r.lopts = append(r.lopts, opt)
}

func (r *reconciler) Reconcile(
	ctx context.Context, req reconcile.Request,
) (reconcile.Result, error) {
	var pod corev1.Pod
	r.mgr.GetClient().Get(ctx, req.NamespacedName, &pod)
	if r.onReconcile != nil {
		return r.onReconcile(ctx, &pod)
	}
	return reconcile.Result{}, nil
}

func (r *reconciler) GetName() string {
	return r.name
}

func (r *reconciler) NewReconciler(ctx context.Context, mgr manager.Manager) reconcile.Reconciler {
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	corev1.AddToScheme(r.mgr.GetScheme())
	if err := r.mgr.GetFieldIndexer().IndexField(ctx, &corev1.Pod{}, "status.phase", func(obj client.Object) []string {
		pod, ok := obj.(*corev1.Pod)
		if !ok || pod.GetDeletionTimestamp() != nil {
			return nil
		}
		return []string{string(pod.Status.Phase)}
	}); err != nil {
		log.Error(err)
	}
	return r
}

func (r *reconciler) For() (client.Object, []builder.ForOption) {
	return new(corev1.Pod), r.forOpts
}

func (*reconciler) Owns() (client.Object, []builder.OwnsOption) {
	return nil, nil
}

func (*reconciler) Watches() (client.Object, handler.EventHandler, []builder.WatchesOption) {
	return nil, nil, nil
}
