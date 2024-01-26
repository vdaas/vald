// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package target

import (
	"context"

	"github.com/vdaas/vald/internal/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Option represents the functional option for reconciler.
type Option func(r *reconciler) error

var defaultOptions = []Option{}

// WithControllerName returns the option to set the name of controller.
func WithControllerName(name string) Option {
	return func(r *reconciler) error {
		if name == "" {
			return errors.NewErrInvalidOption("controllerName", name)
		}
		r.name = name
		return nil
	}
}

// WithManager returns the option to set the controller manager.
func WithManager(mgr manager.Manager) Option {
	return func(r *reconciler) error {
		if mgr == nil {
			return errors.NewErrInvalidOption("manager", mgr)
		}
		r.mgr = mgr
		return nil
	}
}

// WithOnErrorFunc returns the option to set the function to notify an error.
func WithOnErrorFunc(f func(error)) Option {
	return func(r *reconciler) error {
		if f == nil {
			return errors.NewErrInvalidOption("onErrorFunc", f)
		}
		r.onError = f
		return nil
	}
}

// WithOnReconcileFunc returns the option to set the function to get the reconciled result.
func WithOnReconcileFunc(f func(context.Context, map[string]Target)) Option {
	return func(r *reconciler) error {
		if f == nil {
			return errors.NewErrInvalidOption("onReconcileFunc", f)
		}
		r.onReconcile = f
		return nil
	}
}

// WithNamespace returns the option to set the namespace to get resources matching the given namespace..
func WithNamespace(ns string) Option {
	return func(r *reconciler) error {
		if ns == "" {
			return errors.NewErrInvalidOption("namespace", ns)
		}
		r.addListOpts(client.InNamespace(ns))
		return nil
	}
}

// WithLabels returns the option to set the label selector to get resources matching the given label.
func WithLabels(labels map[string]string) Option {
	return func(r *reconciler) error {
		if len(labels) == 0 {
			return errors.NewErrInvalidOption("labels", labels)
		}
		r.addListOpts(client.MatchingLabels(labels))
		return nil
	}
}
