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

// Package pod provides kubernetes pod information and preriodically update
package pod

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/k8s"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	metrics "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type PodWatcher k8s.ResourceController

type reconciler struct {
	ctx         context.Context
	mgr         manager.Manager
	name        string
	onError     func(err error)
	onReconcile func(podList map[string]Pod)
}

type Pod struct {
	Name      string
	Namespace string
	CPU       float64
	Mem       float64
}

func New(opts ...Option) PodWatcher {
	r := new(reconciler)

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *reconciler) Reconcile(req reconcile.Request) (res reconcile.Result, err error) {
	m := &metrics.PodMetricsList{}

	err = r.mgr.GetClient().List(r.ctx, m)

	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100,
		}
		if errors.IsNotFound(err) {
			res = reconcile.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}
		}
		return
	}

	var (
		cpuUsage float64
		memUsage float64
		pods     = make(map[string]Pod, len(m.Items))
	)

	for _, pod := range m.Items {
		cpuUsage = 0.0
		memUsage = 0.0
		for _, container := range pod.Containers {
			cpuUsage += float64(container.Usage.Cpu().Value())
			memUsage += float64(container.Usage.Memory().Value())
		}

		cpuUsage /= float64(len(pod.Containers))
		memUsage /= float64(len(pod.Containers))
		podMetaName := pod.GetObjectMeta().GetName()

		pods[podMetaName] = Pod{
			Name:      pod.GetName(),
			Namespace: pod.GetNamespace(),
			CPU:       cpuUsage,
			Mem:       memUsage,
		}
	}

	if r.onReconcile != nil {
		r.onReconcile(pods)
	}

	return
}

func (r *reconciler) GetName() string {
	return r.name
}

func (r *reconciler) NewReconciler(ctx context.Context, mgr manager.Manager) reconcile.Reconciler {
	if r.ctx == nil && ctx != nil {
		r.ctx = ctx
	}
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	metrics.AddToScheme(r.mgr.GetScheme())
	return r
}

func (r *reconciler) For() runtime.Object {
	// WARN: metrics should be renew
	// https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/resource-metrics-api.md#further-improvements
	return new(metrics.PodMetrics)
}

func (r *reconciler) Owns() runtime.Object {
	// return new(metrics.PodMetrics)
	return nil
}

func (r *reconciler) Watches() (*source.Kind, handler.EventHandler) {
	// return &source.Kind{Type: new(metrics.PodMetrics)}, &handler.EnqueueRequestForObject{}
	return nil, nil
}
