//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package circuitbreaker

import (
	"context"

	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	metricsName        = "circuit_breaker_state"
	metricsDescription = "Current circuit breaker state"
)

type breakerMetrics struct {
	breakerNameKey string
	stateKey       string
}

func New() metrics.Metric {
	return &breakerMetrics{
		breakerNameKey: "name",
		stateKey:       "state",
	}
}

func (*breakerMetrics) View() ([]*metrics.View, error) {
	breakerState, err := view.New(
		view.MatchInstrumentName(metricsName),
		view.WithSetDescription(metricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	return []*metrics.View{
		&breakerState,
	}, nil
}

func (bm *breakerMetrics) Register(m metrics.Meter) error {
	breakerState, err := m.AsyncInt64().Gauge(
		metricsName,
		metrics.WithDescription(metricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			breakerState,
		},
		func(ctx context.Context) {
			ms := circuitbreaker.Metrics(ctx)
			if len(ms) != 0 {
				for name, sts := range ms {
					if len(sts) != 0 {
						for st, cnt := range sts {
							breakerState.Observe(ctx, cnt,
								attribute.String(bm.breakerNameKey, name),
								attribute.String(bm.stateKey, st.String()))
						}
					}
				}
			}
		},
	)
}
