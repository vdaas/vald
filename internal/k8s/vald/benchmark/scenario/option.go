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

package scenario

import (
	"context"

	"github.com/vdaas/vald/internal/k8s/vald"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
)

type (
	// Option configures New; it is the ValdBenchmarkScenario instantiation of
	// the shared watcher factory option.
	Option = vald.WatcherOption[v1.ValdBenchmarkScenario]

	// config is the ValdBenchmarkScenario instantiation of the shared watcher
	// factory config that the options mutate.
	config = vald.WatcherConfig[v1.ValdBenchmarkScenario]
)

// WithControllerName returns Option that sets the reconciler name.
func WithControllerName(name string) Option {
	return vald.WithWatcherControllerName[v1.ValdBenchmarkScenario](name)
}

// WithNamespaces returns Option to restrict List to the given namespaces.
func WithNamespaces(nss ...string) Option {
	return vald.WithWatcherNamespaces[v1.ValdBenchmarkScenario](nss...)
}

// WithOnErrorFunc returns Option that sets the reconcile error callback.
func WithOnErrorFunc(f func(err error)) Option {
	return vald.WithWatcherOnErrorFunc[v1.ValdBenchmarkScenario](f)
}

// WithOnReconcileFunc returns Option that sets the reconcile callback.
func WithOnReconcileFunc(
	f func(ctx context.Context, scenarioList map[string]v1.ValdBenchmarkScenario),
) Option {
	return vald.WithWatcherOnReconcileFunc(f)
}
