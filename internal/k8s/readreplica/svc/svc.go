//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package svc provides kubernetes svc information and preriodically update
package svc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type SvcWatcher k8s.ResourceController

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespace   string
	idKey       string // readreplica.label_key
	onError     func(err error)
	onReconcile func(svcs []ReadReplicaSvc)
	lopts       []client.ListOption
}

type ReadReplicaSvc struct {
	Name      string
	Addr      string
	ReplicaID uint64
}

func New(readreplicaLabel map[string]string, idKey string, opts ...Option) SvcWatcher {
	r := new(reconciler)

	opts = append(opts, WithLabels(readreplicaLabel))
	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}

	r.idKey = idKey
	return r
}

func (r *reconciler) addListOpts(opt client.ListOption) {
	if opt == nil {
		return
	}
	if r.lopts == nil {
		r.lopts = make([]client.ListOption, 0, 1)
	}
	r.lopts = append(r.lopts, opt)
}

func (r *reconciler) Reconcile(ctx context.Context, _ reconcile.Request) (res reconcile.Result, err error) {
	svcList := &corev1.ServiceList{}

	if r.lopts != nil {
		err = r.mgr.GetClient().List(ctx, svcList, r.lopts...)
	} else {
		err = r.mgr.GetClient().List(ctx, svcList)
	}

	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100, //nolint:gomnd
		}
		if errors.IsNotFound(err) {
			res = reconcile.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}
		}
		return res, err
	}

	svcs := make([]ReadReplicaSvc, 0, len(svcList.Items))
	for i := range svcList.Items {
		svc := &svcList.Items[i]
		if svc.GetDeletionTimestamp() != nil {
			log.Debugf("reconcile process will be skipped for node: %s, status: %v, deletion timestamp: %s",
				svc.GetName(),
				svc.Status,
				svc.GetDeletionTimestamp())
			continue
		}
		labels := svc.GetLabels()
		v, ok := labels[r.idKey]
		if !ok {
			log.Errorf("this svc(%s) does not have readreplica id label(%s)", svc.GetName(), r.idKey)
			return reconcile.Result{}, fmt.Errorf("no valid label is put in the svc")
		}
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			log.Error(err)
			return reconcile.Result{}, err
		}
		svcs = append(svcs, ReadReplicaSvc{
			Name:      svc.GetName(),
			Addr:      svc.Spec.ClusterIP,
			ReplicaID: id,
		})
	}
	if r.onReconcile != nil {
		r.onReconcile(svcs)
	}

	return res, nil
}

func (r *reconciler) GetName() string {
	return r.name
}

func (r *reconciler) NewReconciler(_ context.Context, mgr manager.Manager) reconcile.Reconciler {
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	corev1.AddToScheme(r.mgr.GetScheme())
	return r
}

func (*reconciler) For() (client.Object, []builder.ForOption) {
	return new(corev1.Service), nil
}

func (*reconciler) Owns() (client.Object, []builder.OwnsOption) {
	return nil, nil
}

func (*reconciler) Watches() (client.Object, handler.EventHandler, []builder.WatchesOption) {
	return nil, nil, nil
}
