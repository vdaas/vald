// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package target

import (
	"context"
	"reflect"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	mirrv1 "github.com/vdaas/vald/internal/k8s/vald/mirror/api/v1"
	"github.com/vdaas/vald/internal/log"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type (
	MirrorTargetWatcher k8s.ResourceController
	MirrorTarget        = mirrv1.ValdMirrorTarget
	MirrorTargetStatus  = mirrv1.MirrorTargetStatus
	MirrorTargetPhase   = mirrv1.MirrorTargetPhase
)

const (
	MirrorTargetPhasePending      = mirrv1.MirrorTargetPending
	MirrorTargetPhaseConnected    = mirrv1.MirrorTargetConnected
	MirrorTargetPhaseDisconnected = mirrv1.MirrorTargetDisconnected
	MirrorTargetPhaseUnknown      = mirrv1.MirrorTargetUnknown
)

type reconciler struct {
	mgr         manager.Manager
	name        string
	onError     func(err error)
	onReconcile func(ctx context.Context, mm map[string]Target)
	lopts       []client.ListOption
}

type Target struct {
	Colocation string
	Host       string
	Port       int
	Phase      MirrorTargetPhase
}

func New(opts ...Option) (MirrorTargetWatcher, error) {
	r := new(reconciler)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(r); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(err, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return r, nil
}

func (r *reconciler) addListOpts(opt client.ListOption) {
	if opt == nil {
		return
	}
	r.lopts = append(r.lopts, opt)
}

func (r *reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	ml := &mirrv1.ValdMirrorTargetList{}
	if len(r.lopts) != 0 {
		err = r.mgr.GetClient().List(ctx, ml, r.lopts...)
	} else {
		err = r.mgr.GetClient().List(ctx, ml)
	}
	if err != nil {
		r.onError(err)
		if apierr.IsNotFound(err) {
			log.Errorf("not found: %s", err)
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, nil
		}
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 100 * time.Millisecond,
		}, err
	}
	tm := make(map[string]Target)
	for _, m := range ml.Items {
		name := m.GetObjectMeta().GetName()
		tm[name] = Target{
			Colocation: m.Spec.Colocation,
			Host:       m.Spec.Target.Host,
			Port:       m.Spec.Target.Port,
			Phase:      m.Status.Phase,
		}
	}
	r.onReconcile(ctx, tm)
	return res, nil
}

func (r *reconciler) GetName() string {
	return r.name
}

func (r *reconciler) NewReconciler(ctx context.Context, mgr manager.Manager) reconcile.Reconciler {
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}
	mirrv1.AddToScheme(r.mgr.GetScheme())
	return r
}

func (r *reconciler) For() (client.Object, []builder.ForOption) {
	return new(mirrv1.ValdMirrorTarget), nil
}

func (r *reconciler) Owns() (client.Object, []builder.OwnsOption) {
	return nil, nil
}

func (r *reconciler) Watches() (*source.Kind, handler.EventHandler, []builder.WatchesOption) {
	return nil, nil, nil
}
