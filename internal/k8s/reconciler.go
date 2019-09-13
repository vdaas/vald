//
// Copyright (C) 2019 kpango(Yusuke Kato)
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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
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
	Owns() runtime.Object
}

type controller struct {
	name string
	mgr  manager.Manager
	rcs  []ResourceController
}

func New(opts ...Option) (cl Controller, err error) {
	c := new(controller)

	for _, opt := range opts {
		opt(c)
	}

	if c.mgr == nil {
		c.mgr, err = manager.New(
			config.GetConfigOrDie(),
			manager.Options{})
		if err != nil {
			return nil, err
		}
	}

	for _, rc := range c.rcs {
		err = ctrl.NewControllerManagedBy(c.mgr).
			Named(rc.GetName()).
			For(rc.For()).
			Owns(rc.Owns()).
			Complete(rc.NewReconciler(c.mgr))
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *controller) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	go func() {
		defer close(ech)
		err := c.mgr.Start(ctx.Done())
		if err != nil {
			ech <- err
		}
	}()

	return ech
}
