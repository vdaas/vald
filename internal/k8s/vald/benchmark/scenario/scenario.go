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

package scenario

import (
	"context"
	"strconv"
	"time"

	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type BenchmarkScenarioWatcher k8s.ResourceController

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "vald.benchmark.scenario", Version: "v1"}
	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespace   string
	onError     func(err error)
	onReconcile func(ctx context.Context, operatorList map[string]BenchmarkScenarioSpec)
	lopts       []client.ListOption
}

func New(opts ...Option) (BenchmarkScenarioWatcher, error) {
	r := new(reconciler)
	for _, opt := range append(defaultOpts, opts...) {
		// TODO: impl error handling after implement functional option
		opt(r)
	}
	return r, nil
}

func (r *reconciler) AddListOpts(opt client.ListOption) {}

func (r *reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	bs := new(BenchmarkScenarioList)

	if r.lopts == nil {
		err = r.mgr.GetClient().List(ctx, bs, r.lopts...)
	} else {
		err = r.mgr.GetClient().List(ctx, bs)
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
			log.Errorf("not found: %s", err)
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, nil
		}
		return
	}

	var scenarios = make(map[string]BenchmarkScenarioSpec, 0)
	for _, item := range bs.Items {
		name := strconv.FormatInt(time.Now().UnixNano(), 10)
		scenarios[name] = item.Spec
	}

	if r.onReconcile != nil {
		r.onReconcile(ctx, scenarios)
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

	AddToScheme(r.mgr.GetScheme())

	return r
}

func (r *reconciler) For() (client.Object, []builder.ForOption) {
	return nil, nil
}

func (r *reconciler) Owns() (client.Object, []builder.OwnsOption) {
	return nil, nil
}

func (r *reconciler) Watches() (*source.Kind, handler.EventHandler, []builder.WatchesOption) {
	// return &source.Kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil, nil
}
