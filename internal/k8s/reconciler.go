//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package k8s provides kubernetes control functionality
package k8s

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/safety"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type (
	Manager        = manager.Manager
	OwnerReference = v1.OwnerReference
)

type Manager = manager.Manager

type Controller interface {
	Start(ctx context.Context) (<-chan error, error)
	GetManager() Manager
}

type ResourceController interface {
	GetName() string
	NewReconciler(ctx context.Context, mgr manager.Manager) reconcile.Reconciler
	For() (client.Object, []builder.ForOption)
	Owns() (client.Object, []builder.OwnsOption)
	Watches() (client.Object, handler.EventHandler, []builder.WatchesOption)
}

type controller struct {
	eg             errgroup.Group
	name           string
	merticsAddr    string
	leaderElection bool
	mgr            manager.Manager
	rcs            []ResourceController
	der            net.Dialer
}

func New(opts ...Option) (cl Controller, err error) {
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
		c.mgr, err = manager.New(
			cfg,
			manager.Options{
				Scheme:             runtime.NewScheme(),
				LeaderElection:     c.leaderElection,
				MetricsBindAddress: c.merticsAddr,
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
