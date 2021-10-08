//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"strings"
	"time"

	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
	onError     func(err error)
	onReconcile func(podList map[string][]Pod)

	// list options
	namespace string
	labels    map[string]string
}

type Pod struct {
	Name       string
	NodeName   string
	Namespace  string
	IP         string
	CPULimit   float64
	CPURequest float64
	MemLimit   float64
	MemRequest float64
}

func New(opts ...Option) PodWatcher {
	r := new(reconciler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}

	return r
}

func (r *reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	ps := &corev1.PodList{}

	lo := make([]client.ListOption, 3)
	if r.namespace != "" {
		lo = append(lo, client.InNamespace(r.namespace))
	}
	if r.labels != nil || len(r.labels) > 0 {
		lo = append(lo, client.MatchingLabels(r.labels))
	}

	err = r.mgr.GetClient().List(ctx, ps, lo...)
	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100,
		}
		if errors.IsNotFound(err) {
			log.Errorf("not found: %s", err)
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, nil
		}
		return
	}

	var (
		cpuLimit   float64
		cpuRequest float64
		memLimit   float64
		memRequest float64
		pods       = make(map[string][]Pod, len(ps.Items))
	)

	for _, pod := range ps.Items {
		if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodRunning {
			cpuLimit = 0.0
			memLimit = 0.0
			cpuRequest = 0.0
			memRequest = 0.0
			for _, container := range pod.Spec.Containers {
				request := container.Resources.Requests
				limit := container.Resources.Limits
				cpuLimit += float64(limit.Cpu().Value())
				memLimit += float64(limit.Memory().Value())
				cpuRequest += float64(request.Cpu().Value())
				memRequest += float64(request.Memory().Value())
			}
			cpuLimit /= float64(len(pod.Spec.Containers))
			memLimit /= float64(len(pod.Spec.Containers))
			cpuRequest /= float64(len(pod.Spec.Containers))
			memRequest /= float64(len(pod.Spec.Containers))
			podName, ok := pod.GetObjectMeta().GetLabels()["app"]
			if !ok {
				pns := strings.Split(pod.GetName(), "-")
				podName = strings.Join(pns[:len(pns)-1], "-")
			}

			if _, ok := pods[podName]; !ok {
				pods[podName] = make([]Pod, 0, len(ps.Items))
			}

			pods[podName] = append(pods[podName], Pod{
				Name:       pod.GetName(),
				NodeName:   pod.Spec.NodeName,
				Namespace:  pod.GetNamespace(),
				IP:         pod.Status.PodIP,
				CPULimit:   cpuLimit,
				CPURequest: cpuRequest,
				MemLimit:   memLimit,
				MemRequest: memRequest,
			})
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

func (r *reconciler) NewReconciler(mgr manager.Manager) reconcile.Reconciler {
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	corev1.AddToScheme(r.mgr.GetScheme())
	return r
}

func (r *reconciler) For() (client.Object, []builder.ForOption) {
	return new(corev1.Pod), nil
}

func (r *reconciler) Owns() (client.Object, []builder.OwnsOption) {
	return nil, nil
}

func (r *reconciler) Watches() (*source.Kind, handler.EventHandler, []builder.WatchesOption) {
	// return &source.Kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil, nil
}
