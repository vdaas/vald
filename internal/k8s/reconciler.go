//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package k8s provides kubernetes controll functionallity
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
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Controller interface {
	Start(ctx context.Context) <-chan error
}

type ResourceController interface {
	GetName() string
	NewReconciler(mgr manager.Manager) reconcile.Reconciler
	For() runtime.Object
	// Owns() runtime.Object
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
		c.mgr, err = manager.New(
			cfg,
			manager.Options{
				Scheme:             runtime.NewScheme(),
				LeaderElection:     c.leaderElection,
				MetricsBindAddress: c.merticsAddr,
			})
		if err != nil {
			return nil, err
		}
	}

	for _, rc := range c.rcs {
		if rc != nil {
			err = builder.ControllerManagedBy(c.mgr).
				Named(rc.GetName()).
				For(rc.For()).
				Complete(rc.NewReconciler(c.mgr))
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
