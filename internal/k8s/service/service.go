//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package service

import (
	"context"
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

// SvcWatcher is a type alias of k8s.ResourceController for service resources.
type SvcWatcher k8s.ResourceController

// Service represents a kubernetes service information.
type Service struct {
	Name        string
	ClusterIP   string
	ClusterIPs  []string
	Ports       []servicePort
	Labels      map[string]string
	Annotations map[string]string
}

type servicePort struct {
	Name string
	Port int32
}

type reconciler struct {
	mgr         manager.Manager
	name        string
	namespace   string
	onError     func(err error)
	onReconcile func(svcs []Service)
	lopts       []client.ListOption
}

// New returns a new SvcWatcher instance.
func New(opts ...Option) SvcWatcher {
	r := new(reconciler)
	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}

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

func extractAPIPorts(ports []corev1.ServicePort) []servicePort {
	var apiPorts []servicePort
	for _, port := range ports {
		if port.Name == "grpc" || port.Name == "rest" {
			apiPorts = append(apiPorts, servicePort{
				Name: port.Name,
				Port: port.Port,
			})
		}
	}
	return apiPorts
}

// Reconcile reconciles the service resources and put the information into the Service struct.
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

	svcs := make([]Service, 0, len(svcList.Items))
	for i := range svcList.Items {
		svc := &svcList.Items[i]
		if svc.GetDeletionTimestamp() != nil {
			log.Debugf("reconcile process will be skipped for svc: %s, status: %v, deletion timestamp: %s",
				svc.GetName(),
				svc.Status,
				svc.GetDeletionTimestamp())
			continue
		}

		ports := extractAPIPorts(svc.Spec.Ports)
		svcs = append(svcs, Service{
			Name:        svc.GetName(),
			ClusterIP:   svc.Spec.ClusterIP,
			ClusterIPs:  svc.Spec.ClusterIPs,
			Ports:       ports,
			Labels:      svc.GetLabels(),
			Annotations: svc.GetAnnotations(),
		})
	}
	if r.onReconcile != nil {
		r.onReconcile(svcs)
	}

	return res, nil
}

// GetName returns the reconciler name.
func (r *reconciler) GetName() string {
	return r.name
}

// NewReconciler returns a new reconciler instance with corev1 scheme added.
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
