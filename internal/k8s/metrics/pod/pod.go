//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metrics "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type PodWatcher k8s.ResourceController

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespace   string
	onError     func(err error)
	onReconcile func(podList map[string]Pod)
	lopts       []client.ListOption
}

type Pod struct {
	Name      string
	Namespace string
	CPU       float64
	Mem       float64
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

func (r *reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	m := &metrics.PodMetricsList{}

	if r.lopts != nil {
		err = r.mgr.GetClient().List(ctx, m, r.lopts...)
	} else {
		err = r.mgr.GetClient().List(ctx, m)
	}

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
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	metrics.AddToScheme(r.mgr.GetScheme())
	if err := r.mgr.GetFieldIndexer().IndexField(ctx, &metrics.PodMetrics{}, "containers.name", func(obj client.Object) []string {
		pod, ok := obj.(*metrics.PodMetrics)
		if !ok {
			return nil
		}
		res := make([]string, 0, len(pod.Containers))
		for _, pc := range pod.Containers {
			res = append(res, pc.Name)
		}
		return res
	}); err != nil {
		log.Error(err)
	}
	return r
}

func (*reconciler) For() (client.Object, []builder.ForOption) {
	// WARN: metrics should be renew
	// https://github.com/kubernetes/community/blob/main/contributors/design-proposals/instrumentation/resource-metrics-api.md#further-improvements
	return new(metrics.PodMetrics), nil
}

func (*reconciler) Owns() (client.Object, []builder.OwnsOption) {
	// return new(metrics.PodMetrics)
	return nil, nil
}

func (*reconciler) Watches() (*source.Kind, handler.EventHandler, []builder.WatchesOption) {
	// return &source.Kind{Type: new(metrics.PodMetrics)}, &handler.EnqueueRequestForObject{}
	return nil, nil, nil
}
