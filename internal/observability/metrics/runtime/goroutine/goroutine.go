// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package goroutine

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type goroutine struct{}

func New() metrics.Metric {
	return &goroutine{}
}

func (g *goroutine) Register(m metrics.Meter) error {
	conter, err := m.AsyncInt64().Gauge(
		"goroutine_count",
		metrics.WithDescription("number of goroutines"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			conter,
		},
		func(ctx context.Context) {
			conter.Observe(ctx, int64(runtime.NumGoroutine()))
		},
	)
}
