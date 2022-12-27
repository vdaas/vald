//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

package job

import (
	"context"

	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)


type Option func(*reconciler) error

var defaultOpts = []Option{}

// WithControllerName returns Option that sets r.name.
func WithControllerName(name string) Option {
	return func(r *reconciler) error {
		r.name = name
		return nil
	}
}

// WithManager returns Option that sets r.mgr.
func WithManager(mgr manager.Manager) Option {
	return func(r *reconciler) error {
		r.mgr = mgr
		return nil
	}
}

// WithNamespaces returns Option to set the namespace.
func WithNamespaces(nss ...string) Option {
	return func(r *reconciler) error {
		r.namespaces = nss
		return nil
	}
}

// WithOnErrorFunc returns Option that sets r.onError.
func WithOnErrorFunc(f func(err error)) Option {
	return func(r *reconciler) error {
		r.onError = f
		return nil
	}
}

// WithOnReconcileFunc returns Option that sets r.onReconcile.
func WithOnReconcileFunc(f func(ctx context.Context, jobList map[string]v1.BenchmarkJobSpec)) Option {
	return func(r *reconciler) error {
		r.onReconcile = f
		return nil
	}
}
