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
package cgo

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	metricsName        = "cgo_call_count"
	metricsDescription = "Number of cgo call"
)

type cgo struct{}

func New() metrics.Metric {
	return &cgo{}
}

func (c *cgo) View() ([]*metrics.View, error) {
	count, err := view.New(
		view.MatchInstrumentName(metricsName),
		view.WithSetDescription(metricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}
	return []*metrics.View{
		&count,
	}, nil
}

func (c *cgo) Register(m metrics.Meter) error {
	count, err := m.AsyncInt64().UpDownCounter(
		metricsName,
		metrics.WithDescription(metricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			count,
		},
		func(ctx context.Context) {
			count.Observe(ctx, int64(runtime.NumGoroutine()))
		},
	)
}
