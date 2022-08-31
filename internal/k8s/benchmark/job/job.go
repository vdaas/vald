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

package job

import (
	"context"

	"github.com/vdaas/vald/internal/k8s"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type BenchmarkJobWatcher k8s.ResourceController

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespace   string
	onError     func(err error)
	onReconcile func(jobList map[string][]BenchmarkJob)
	lopts       []client.ListOption
}

type BenchmarkJob struct {
	Target     *BenchmarkTarget
	Dataset    *BenchmarkDataset
	Replica    int
	Repetition int
	JobType    string
	Dimension  int
	Epsilon    float32
	Radius     float32
	Iter       int
	Num        int32
	MinNUm     int32
	Timeout    string
	Rules      []*BenchmarkJobRule
}

type BenchmarkTarget struct {
	Host string
	Port int
}

type BenchmarkDataset struct {
	Name    string
	Group   string
	Indexes int
	Range   *BenchmarkDatasetRange
}

type BenchmarkDatasetRange struct {
	Start int
	End   int
}

type BenchmarkJobRule struct {
	Name string
	Type string
}

func New(opts ...Option) BenchmarkJobWatcher {
	r := new(reconciler)
	for _, opt := range append(defaultOpts, opts...) {
		opt(r)
	}
	return r
}

func (r *reconciler) AddListOpts(opt client.ListOption) {}

//+kubebuilder:rbac:groups=vald.vdaas.org.vald.vdaas.org,resources=valdbenchmarkjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vald.vdaas.org.vald.vdaas.org,resources=valdbenchmarkjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vald.vdaas.org.vald.vdaas.org,resources=valdbenchmarkjobs/finalizers,verbs=update
func (r *reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	return
}

func (r *reconciler) GetName() string {
	return r.name
}

func (r *reconciler) NewReconciler(ctx context.Context, mgr manager.Manager) reconcile.Reconciler {
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
