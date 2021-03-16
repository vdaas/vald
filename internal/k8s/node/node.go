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

// Package node provides kubernetes node information and preriodically update
package node

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type NodeWatcher k8s.ResourceController

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespace   string
	onError     func(err error)
	onReconcile func(nodes []Node)
}

type Node struct {
	Name         string
	InternalAddr string
	ExternalAddr string
	CPUCapacity  float64
	CPURemain    float64
	MemCapacity  float64
	MemRemain    float64
}

func New(opts ...Option) NodeWatcher {
	r := new(reconciler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}

	return r
}

func (r *reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	ns := &corev1.NodeList{}

	err = r.mgr.GetClient().List(ctx, ns)

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

	nodes := make([]Node, 0, len(ns.Items))

	for _, node := range ns.Items {
		if node.GetObjectMeta().GetDeletionTimestamp() != nil {
			continue
		}
		remain := node.Status.Allocatable
		limit := node.Status.Capacity
		var eip, iip string
		for _, addr := range node.Status.Addresses {
			switch addr.Type {
			case corev1.NodeInternalIP:
				iip = addr.Address
			case corev1.NodeInternalDNS:
				if iip == "" {
					iip = addr.Address
				}
			case corev1.NodeExternalIP:
				eip = addr.Address
			case corev1.NodeExternalDNS:
				if eip == "" {
					eip = addr.Address
				}
			}
		}
		nodes = append(nodes, Node{
			Name:         node.GetName(),
			ExternalAddr: eip,
			InternalAddr: iip,
			CPUCapacity:  float64(limit.Cpu().Value()),
			CPURemain:    float64(remain.Cpu().Value()),
			MemCapacity:  float64(limit.Memory().Value()),
			MemRemain:    float64(remain.Memory().Value()),
		})
	}
	if r.onReconcile != nil {
		r.onReconcile(nodes)
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
	return new(corev1.Node), nil
}

func (r *reconciler) Owns() (client.Object, []builder.OwnsOption) {
	return nil, nil
}

func (r *reconciler) Watches() (*source.Kind, handler.EventHandler, []builder.WatchesOption) {
	// return &source.Kind{Type: new(corev1.Node)}, &handler.EnqueueRequestForObject{}
	return nil, nil, nil
}
