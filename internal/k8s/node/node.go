//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
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
package node

import (
	"context"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type NodeWatcher interface {
	k8s.ResourceController
	GetNodes() []Node
}

type reconciler struct {
	mu          sync.RWMutex
	nodes       []Node
	mgr         manager.Manager
	name        string
	onError     func(err error)
	onReconcile func(nodes []Node)
}

type Node struct {
	Name string
	IP   string
	CPU  float64
	Mem  float64
}

func New(opts ...Option) NodeWatcher {
	r := new(reconciler)

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *reconciler) Reconcile(req reconcile.Request) (res reconcile.Result, err error) {
	ns := &corev1.NodeList{}

	err = r.mgr.GetClient().List(context.TODO(), ns, client.InNamespace(req.Namespace))

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
		remain := node.Status.Allocatable
		limit := node.Status.Capacity
		nodes = append(nodes, Node{
			Name: node.GetName(),
			IP:   node.Status.Addresses[0].Address,
			CPU: (float64(limit.Cpu().Value()-remain.Cpu().Value()) /
				float64(limit.Cpu().Value())) * 100.0,
			Mem: (float64(limit.Memory().Value()-remain.Memory().Value()) /
				float64(limit.Memory().Value())) * 100.0,
		})
	}

	r.mu.Lock()
	r.nodes = nodes
	r.mu.Unlock()

	if r.onReconcile != nil {
		r.mu.RLock()
		r.onReconcile(r.nodes)
		r.mu.RUnlock()
	}

	return
}

func (r *reconciler) GetNodes() (nodes []Node) {
	r.mu.RLock()
	nodes = r.nodes
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
	corev1.AddToScheme(r.mgr.GetScheme())
	return r
}

func (r *reconciler) For() runtime.Object {
	return new(corev1.NodeList)
}

func (r *reconciler) Owns() runtime.Object {
	return new(corev1.Node)
}
