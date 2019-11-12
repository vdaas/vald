//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"sync"
	"time"

	"github.com/vdaas/vald/internal/k8s"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PodWatcher interface {
	k8s.ResourceController
	GetPods(name string) (pods []Pod, ok bool)
}

type reconciler struct {
	mu          sync.RWMutex
	podList     map[string][]Pod
	mgr         manager.Manager
	name        string
	onError     func(err error)
	onReconcile func(podList map[string][]Pod)
}

type Pod struct {
	Name     string
	NodeName string
	IP       string
	CPU      float64
	Mem      float64
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

	err = r.mgr.GetClient().List(context.TODO(), ps, client.InNamespace(req.Namespace), client.MatchingFields{
		"status.phase": string(corev1.PodRunning),
	})

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
		pods     = make(map[string][]Pod, len(ps.Items))
	)

	for _, pod := range ps.Items {
		cpuUsage = 0.0
		memUsage = 0.0
		for _, container := range pod.Spec.Containers {
			request := container.Resources.Requests
			limit := container.Resources.Limits
			cpuUsage += (float64(request.Cpu().Value()) /
				float64(limit.Cpu().Value())) * 100.0
			memUsage += (float64(request.Memory().Value()) /
				float64(limit.Memory().Value())) * 100.0
		}

		cpuUsage = cpuUsage / float64(len(pod.Spec.Containers))
		memUsage = memUsage / float64(len(pod.Spec.Containers))

		// pod.GetObjectMeta().GetLabels()["app"]
		podMetaName := pod.GetObjectMeta().GetName()

		if _, ok := pods[podMetaName]; !ok {
			pods[podMetaName] = make([]Pod, 0, len(ps.Items))
		}

		pods[podMetaName] = append(pods[podMetaName], Pod{
			Name:     pod.GetName(),
			NodeName: pod.Spec.NodeName,
			IP:       pod.Status.PodIP,
			CPU:      cpuUsage,
			Mem:      memUsage,
		})
	}

	if r.onReconcile != nil {
		r.onReconcile(pods)
	}

	r.mu.Lock()
	r.podList = pods
	r.mu.Lock()

	return
}

func (r *reconciler) GetPods(name string) (pods []Pod, ok bool) {
	r.mu.RLock()
	pods, ok = r.podList[name]
	r.mu.RUnlock()
	return
}

func (r *reconciler) GetName() string {
	return r.name
}

func (r *reconciler) NewReconciler(mgr manager.Manager) reconcile.Reconciler {
	if r.mgr == nil {
		r.mgr = mgr
	}
	return r
}

func (r *reconciler) For() runtime.Object {
	return new(appsv1.ReplicaSet)
}

// func (r *reconciler) Owns() runtime.Object {
// 	return new(corev1.Pod)
// }
