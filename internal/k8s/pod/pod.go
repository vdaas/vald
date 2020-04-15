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

package pod

import (
	"context"
	"strings"
	"time"

	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
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
	namespace   string
	onError     func(err error)
	onReconcile func(podList map[string][]Pod)
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

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *reconciler) Reconcile(req reconcile.Request) (res reconcile.Result, err error) {
	ps := &corev1.PodList{}

	err = r.mgr.GetClient().List(r.ctx, ps)

	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100,
		}
		if errors.IsNotFound(err) {
			log.Error("not found", err)
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

func (r *reconciler) NewReconciler(ctx context.Context, mgr manager.Manager) reconcile.Reconciler {
	if r.ctx == nil && ctx != nil {
		r.ctx = ctx
	}
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	corev1.AddToScheme(r.mgr.GetScheme())
	return r
}

func (r *reconciler) For() runtime.Object {
	return new(corev1.Pod)
}

func (r *reconciler) Owns() runtime.Object {
	return nil
}

func (r *reconciler) Watches() (*source.Kind, handler.EventHandler) {
	// return &source.Kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil
}
