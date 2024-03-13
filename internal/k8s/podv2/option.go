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

// Package podv2 provides kubernetes pod information and preriodically update
package podv2

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Option func(*reconciler) error

var defaultOptions = []Option{}

func WithControllerName(name string) Option {
	return func(r *reconciler) error {
		r.name = name
		return nil
	}
}

func WithManager(mgr manager.Manager) Option {
	return func(r *reconciler) error {
		r.mgr = mgr
		return nil
	}
}

func WithOnErrorFunc(f func(err error)) Option {
	return func(r *reconciler) error {
		r.onError = f
		return nil
	}
}

func WithOnReconcileFunc(f func(ctx context.Context, pod *corev1.Pod) (reconcile.Result, error)) Option {
	return func(r *reconciler) error {
		r.onReconcile = f
		return nil
	}
}

func WithNamespace(ns string) Option {
	return func(r *reconciler) error {
		if ns == "" {
			return nil
		}
		r.namespace = ns
		r.addListOpts(client.InNamespace(ns))
		return nil
	}
}

func WithLabels(ls map[string]string) Option {
	return func(r *reconciler) error {
		if len(ls) > 0 {
			r.addListOpts(client.MatchingLabels(ls))
		}
		return nil
	}
}

func WithFields(fs map[string]string) Option {
	return func(r *reconciler) error {
		if len(fs) > 0 {
			r.addListOpts(client.MatchingFields(fs))
		}
		return nil
	}
}

func WithForOpts(fopts ...builder.ForOption) Option {
	return func(r *reconciler) error {
		if len(fopts) == 0 {
			return nil
		}
		r.forOpts = append(r.forOpts, fopts...)
		return nil
	}
}
