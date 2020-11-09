package main

import (
	"context"
	"reflect"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type ReplicaSetWatcher k8s.ResourceController

type reconciler struct {
	ctx         context.Context
	mgr         manager.Manager
	name        string
	namespace   string
	onError     func(err error)
	onReconcile func(rs []ReplicaSet)
}

type ReplicaSet = appsv1.ReplicaSet

func New(opts ...Option) (ReplicaSetWatcher, error) {
	r := new(reconciler)

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return erros.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return r
}

func (r *reconciler) Reconcile(req reconcile.Request) (res reconcile.Result, err error) {
	rsl := new(appsv1.ReplicaSetList)

	err = r.mgr.GetClient().List(r.ctx, rsl)
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

	rs := make(map[string][]ReplicaSet, len(rsl.Items))

	for _, replicaset := range rsl.Items {
		name, ok := replicaset.GetObjectMeta().GetLabels()["app"]
		if !ok {
			pns := strings.Split(pod.GetName(), "-")
			name = strings.Join(pns[:len(pns)-1], "-")
		}

		if _, ok := rs[name]; !ok {
			rs[name] = make([]ReplicaSet, 0, len(rsl.Items))
		}

		rs[name] = append(rs[name], replicaset)
	}

	if r.onReconcile != nil {
		r.onReconcile(rs)
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
	appv1.AddToScheme(r.mgr.GetScheme())
	return r
}

func (r *reconciler) For() runtime.Object {
	return new(appv1.ReplicaSet)
}

func (r *reconciler) Owns() runtime.Object {
	return nil
}

func (r *reconciler) Watches() (*source.Kind, handler.EventHandler) {
	// return &source.Kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil
}
