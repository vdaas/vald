// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package k8s provides kubernetes control functionality
package k8s

import (
	"context"
	"reflect"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	ctrlcontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	mserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Controller interface {
	Start(ctx context.Context) (<-chan error, error)
	GetManager() Manager
}

//nolint:gochecknoglobals // immutable function alias confining k8s.io imports to this package
var Now = metav1.Now

type ResourceController interface {
	GetName() string
	NewReconciler(ctx context.Context, mgr manager.Manager) reconcile.Reconciler
	For() (client.Object, []builder.ForOption)
	Owns() (client.Object, []builder.OwnsOption)
	Watches() (client.Object, handler.EventHandler, []builder.WatchesOption)
}

// ConcurrentReconciler is an optional interface a ResourceController can
// implement to request a per-controller worker count. When the returned value
// is greater than zero it is applied as MaxConcurrentReconciles on the built
// controller; existing ResourceController implementations are unaffected.
type ConcurrentReconciler interface {
	MaxConcurrentReconciles() int
}

type controller struct {
	eg                      errgroup.Group
	mgr                     manager.Manager
	der                     net.Dialer
	name                    string
	metricsAddr             string
	leaderElectionID        string
	leaderElectionNamespace string
	leaseDuration           *time.Duration
	renewDeadline           *time.Duration
	retryPeriod             *time.Duration
	syncPeriod              *time.Duration
	cacheNamespaces         []string
	rcs                     []ResourceController
	leaderElection          bool
}

func New(opts ...Option) (cl Controller, err error) {
	setControllerRuntimeLogger()
	c := new(controller)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	if c.mgr == nil {
		cfg, err := config.GetConfig()
		if err != nil {
			return nil, err
		}
		if cfg == nil {
			return nil, errors.ErrInvalidReconcilerConfig
		}
		if c.der != nil {
			cfg.Dial = c.der.GetDialer()
		}
		copts := cache.Options{
			SyncPeriod: c.syncPeriod,
		}
		if len(c.cacheNamespaces) > 0 {
			copts.DefaultNamespaces = make(map[string]cache.Config, len(c.cacheNamespaces))
			for _, ns := range c.cacheNamespaces {
				copts.DefaultNamespaces[ns] = cache.Config{}
			}
		}
		c.mgr, err = manager.New(
			cfg,
			manager.Options{
				Scheme:                  runtime.NewScheme(),
				LeaderElection:          c.leaderElection,
				LeaderElectionID:        c.leaderElectionID,
				LeaderElectionNamespace: c.leaderElectionNamespace,
				LeaseDuration:           c.leaseDuration,
				RenewDeadline:           c.renewDeadline,
				RetryPeriod:             c.retryPeriod,
				Cache:                   copts,
				Metrics:                 mserver.Options{BindAddress: c.metricsAddr},
			},
		)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *controller) Start(ctx context.Context) (<-chan error, error) {
	if c.der != nil {
		c.der.StartDialerCache(ctx)
	}
	for _, rc := range c.rcs {
		if rc != nil {
			bc := builder.ControllerManagedBy(c.mgr).Named(rc.GetName())
			f, fopts := rc.For()
			if f != nil {
				bc = bc.For(f, fopts...)
			}
			o, oopts := rc.Owns()
			if o != nil {
				bc = bc.Owns(o, oopts...)
			}
			src, h, wopts := rc.Watches()
			if src != nil {
				if h == nil {
					h = &handler.EnqueueRequestForObject{}
				}
				bc = bc.Watches(src, h, wopts...)
			}
			if cr, ok := rc.(ConcurrentReconciler); ok {
				if n := cr.MaxConcurrentReconciles(); n > 0 {
					bc = bc.WithOptions(ctrlcontroller.Options{
						MaxConcurrentReconciles: n,
					})
				}
			}
			_, err := bc.Build(rc.NewReconciler(ctx, c.mgr))
			if err != nil {
				return nil, err
			}
		}
	}
	ech := make(chan error, 1)
	c.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		err := c.mgr.Start(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
			case ech <- err:
			}
		}
		return nil
	}))

	return ech, nil
}

func (c *controller) GetManager() Manager {
	return c.mgr
}

//nolint:gochecknoglobals // immutable function alias confining k8s.io imports to this package
var AddClientGoScheme = clientgoscheme.AddToScheme
