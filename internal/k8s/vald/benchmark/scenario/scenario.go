//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package scenario

import (
	"context"
	"reflect"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/log"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type BenchmarkScenarioWatcher k8s.ResourceController

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespaces  []string
	onError     func(err error)
	onReconcile func(ctx context.Context, operatorList map[string]v1.ValdBenchmarkScenario)
	lopts       []client.ListOption
}

func New(opts ...Option) (BenchmarkScenarioWatcher, error) {
	r := new(reconciler)
	for _, opt := range append(defaultOpts, opts...) {
		err := opt(r)
		if err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return r, nil
}

func (r *reconciler) AddListOpts(opt client.ListOption) {
	if opt == nil {
		return
	}
	if r.lopts == nil {
		r.lopts = make([]client.ListOption, 0, 1)
	}
	r.lopts = append(r.lopts, opt)
}

func (r *reconciler) Reconcile(
	ctx context.Context, _ reconcile.Request,
) (res reconcile.Result, err error) {
	bs := new(v1.ValdBenchmarkScenarioList)

	err = r.mgr.GetClient().List(ctx, bs, r.lopts...)
	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100,
		}
		if k8serrors.IsNotFound(err) {
			log.Errorf("not found: %s", err)
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, nil
		}
		return
	}
	scenarios := make(map[string]v1.ValdBenchmarkScenario, 0)
	for _, item := range bs.Items {
		name := item.Name
		scenarios[name] = item
	}

	if r.onReconcile != nil {
		r.onReconcile(ctx, scenarios)
	}

	return
}

func (r *reconciler) GetName() string {
	return r.name
}

func (r *reconciler) NewReconciler(_ context.Context, mgr manager.Manager) reconcile.Reconciler {
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	v1.AddToScheme(r.mgr.GetScheme())

	return r
}

func (*reconciler) For() (client.Object, []builder.ForOption) {
	return new(v1.ValdBenchmarkScenario), nil
}

func (*reconciler) Owns() (client.Object, []builder.OwnsOption) {
	return nil, nil
}

func (*reconciler) Watches() (client.Object, handler.EventHandler, []builder.WatchesOption) {
	// return &source.Kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	// return &source.Kind{Type: new(v1.ValdBenchmarkScenario)}, &handler.EnqueueRequestForObject{}, nil
	return nil, &handler.EnqueueRequestForObject{}, nil
}
