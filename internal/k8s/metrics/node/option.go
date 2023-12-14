//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package node provides kubernetes node information and preriodically update
package node

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
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

func WithOnReconcileFunc(f func(nodes map[string]Node)) Option {
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
