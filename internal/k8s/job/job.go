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
package job

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
)

// JobWatcher is a type alias for k8s resource controller.
type JobWatcher k8s.ResourceController

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespaces  []string
	onError     func(err error)
	onReconcile func(jobList map[string][]Job)
	pool        sync.Pool
}

// Job is a type alias for the k8s job definition.
type Job = batchv1.Job

// New returns the JobWatcher that implements reconciliation loop, or any errors occurred.
func New(opts ...Option) (JobWatcher, error) {
	r := &reconciler{
		pool: sync.Pool{
			New: func() interface{} {
				return make(map[string][]Job)
			},
		},
	}

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return r, nil
}

// Reconcile implements k8s reconciliation loop to retrieve the Job information from k8s.
func (r *reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	js := new(batchv1.JobList)

	listOpts := make([]client.ListOption, 0, len(r.namespaces))
	for _, ns := range r.namespaces {
		listOpts = append(listOpts, client.InNamespace(ns))
	}

	err = r.mgr.GetClient().List(ctx, js, listOpts...)
	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100,
		}
		if k8serrors.IsNotFound(err) {
			log.Error("not found", err)
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, nil
		}
		return
	}

	jobs := r.pool.Get().(map[string][]Job)
	for _, job := range js.Items {
		name, ok := job.GetObjectMeta().GetLabels()["app"]
		if !ok {
			jns := strings.Split(job.GetName(), "-")
			name = strings.Join(jns[:len(jns)-1], "-")
		}

		if _, ok := jobs[name]; !ok {
			jobs[name] = make([]Job, 0, len(js.Items))
		}
		jobs[name] = append(jobs[name], job)
	}

	if r.onReconcile != nil {
		r.onReconcile(jobs)
	}

	for name := range jobs {
		jobs[name] = jobs[name][:0:len(jobs[name])]
	}

	r.pool.Put(jobs)

	return
}

// GetName returns the name of resource controller.
func (r *reconciler) GetName() string {
	return r.name
}

// NewReconciler returns the reconciler for the Job.
func (r *reconciler) NewReconciler(mgr manager.Manager) reconcile.Reconciler {
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}

	batchv1.AddToScheme(r.mgr.GetScheme())
	return r
}

// For returns the runtime.Object which is job.
func (r *reconciler) For() (client.Object, []builder.ForOption) {
	return new(batchv1.Job), nil
}

// Owns returns the owner of the job watcher.
// It will always return nil.
func (r *reconciler) Owns() (client.Object, []builder.OwnsOption) {
	return nil, nil
}

// Watches returns the kind of the job and the event handler.
// It will always return nil.
func (r *reconciler) Watches() (*source.Kind, handler.EventHandler, []builder.WatchesOption) {
	// return &source.Kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil, nil
}

// Create creates batchjob object along with job template and given context.
func Create(ctx context.Context, job *Job) error {
	// TODO: get config from cache.
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	scheme := runtime.NewScheme()
	if err = batchv1.AddToScheme(scheme); err != nil {
		return err
	}

	c, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return err
	}

	if err := c.Create(ctx, job); err != nil {
		return err
	}
	return nil
}
