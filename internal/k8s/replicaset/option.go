//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
package replicaset

import "sigs.k8s.io/controller-runtime/pkg/manager"

// Option represents functional option for reconciler.
type Option func(*reconciler) error

var defaultOpts = []Option{}

// WithControllerName returns Option to set the controller name.
func WithControllerName(name string) Option {
	return func(r *reconciler) error {
		r.name = name
		return nil
	}
}

// WithNamespace returns Option to set the namespace.
func WithNamespace(ns string) Option {
	return func(r *reconciler) error {
		r.namespace = ns
		return nil
	}
}

// WithManager returns Option to set the controller manager.
func WithManager(mgr manager.Manager) Option {
	return func(r *reconciler) error {
		r.mgr = mgr
		return nil
	}
}

// WithOnErrorFunc returns Option to set the onError hook.
func WithOnErrorFunc(f func(err error)) Option {
	return func(r *reconciler) error {
		r.onError = f
		return nil
	}
}

// WithOnReconcileFunc returns Option to set the onReconcile hook.
func WithOnReconcileFunc(f func(replicasetList map[string][]ReplicaSet)) Option {
	return func(r *reconciler) error {
		r.onReconcile = f
		return nil
	}
}
