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

// Package node provides kubernetes node information and preriodically update
package node

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

type NodeWatcher k8s.ResourceController

type reconciler struct {
	ctx         context.Context
	mgr         manager.Manager
	name        string
	onError     func(err error)
	onReconcile func(nodeList map[string]Node)
}

type Node struct {
	Name    string
	CPU     float64
	Mem     float64
	Pods    int64
	Storage int64
}

func New(opts ...Option) NodeWatcher {
	r := new(reconciler)

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *reconciler) Reconcile(req reconcile.Request) (res reconcile.Result, err error) {
	m := &metrics.NodeMetricsList{}

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

	nodes := make(map[string]Node, len(m.Items))

	for _, node := range m.Items {
		nodeName := node.GetName()
		nodes[nodeName] = Node{
			Name:    nodeName,
			CPU:     float64(node.Usage.Cpu().Value()),
			Mem:     float64(node.Usage.Memory().Value()),
			Storage: node.Usage.StorageEphemeral().Value(),
			Pods:    node.Usage.Pods().Value(),
		}
	}

	if r.onReconcile != nil {
		r.onReconcile(nodes)
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
	return new(metrics.NodeMetrics)
}

func (r *reconciler) Owns() runtime.Object {
	// return new(metrics.PodMetrics)
	return nil
}

func (r *reconciler) Watches() (*source.Kind, handler.EventHandler) {
	// return &source.Kind{Type: new(metrics.NodeMetrics)}, &handler.EnqueueRequestForObject{}
	return nil, nil
}
