//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package scenario

import (
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/vald"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
)

type BenchmarkWatcher k8s.ResourceController

// New builds a BenchmarkWatcher that lists ValdBenchmarkScenario
// objects on every reconcile and hands them to the configured
// OnReconcileFunc callback as a map keyed by object name. Option application
// errors abort construction only when critical; any other option error is
// logged as a warning and skipped.
func New(opts ...Option) (BenchmarkWatcher, error) {
	return vald.NewListWatcher(
		v1.AddToScheme,
		func(sc *v1.ValdBenchmarkScenario) v1.ValdBenchmarkScenario { return *sc },
		vald.SkipNonCriticalOptionError,
		opts...,
	)
}
