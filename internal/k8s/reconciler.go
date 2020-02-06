//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"github.com/vdaas/vald/internal/safety"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type Controller interface {
	Start(ctx context.Context) <-chan error
}

type ResourceController interface {
	GetName() string
	NewReconciler(mgr manager.Manager) reconcile.Reconciler
	For() runtime.Object
	Owns() runtime.Object
	Watches() (*source.Kind, handler.EventHandler)
}

type controller struct {
	eg             errgroup.Group
	name           string
	merticsAddr    string
	leaderElection bool
	mgr            manager.Manager
	rcs            []ResourceController
}

func New(opts ...Option) (cl Controller, err error) {
	c := new(controller)

	for _, opt := range append(defaultOpts, opts...) {
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

	for _, rc := range c.rcs {
		if rc != nil {
			bc := builder.ControllerManagedBy(c.mgr).Named(rc.GetName())
			f := rc.For()
			if f != nil {
				bc = bc.For(f)
			}
			o := rc.Owns()
			if o != nil {
				bc = bc.Owns(o)
			}
			src, h := rc.Watches()
			if src != nil {
				if h == nil {
					h = &handler.EnqueueRequestForObject{}
				}
				bc = bc.Watches(src, h)
			}
			_, err = bc.Build(rc.NewReconciler(c.mgr))
			if err != nil {
				return nil, err
			}
		}
	}

	return c, nil
}

func (c *controller) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	c.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		err := c.mgr.Start(ctx.Done())
		if err != nil {
			ech <- err
		}
		return nil
	}))

	return ech
}
