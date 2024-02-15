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
package circuitbreaker

import (
	"context"

	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
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

func (*breakerMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        metricsName,
				Description: metricsDescription,
			},
			view.Stream{
<<<<<<< HEAD
				Aggregation: view.AggregationLastValue{},
=======
				Aggregation: meric.AggregationLastValue{},
>>>>>>> feature/agent/qbg
			},
		),
	}, nil
}

func (bm *breakerMetrics) Register(m metrics.Meter) error {
	breakerState, err := m.Int64ObservableGauge(
		metricsName,
		metrics.WithDescription(metricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(
		func(ctx context.Context, o api.Observer) error {
			ms := circuitbreaker.Metrics(ctx)
			if len(ms) != 0 {
				for name, sts := range ms {
					if len(sts) != 0 {
						for st, cnt := range sts {
							o.ObserveInt64(breakerState, cnt,
<<<<<<< HEAD
								api.WithAttributes(
									attribute.String(bm.breakerNameKey, name),
									attribute.String(bm.stateKey, st.String())))
=======
								attribute.String(bm.breakerNameKey, name),
								attribute.String(bm.stateKey, st.String()))
>>>>>>> feature/agent/qbg
						}
					}
				}
			}
			return nil
		}, breakerState,
	)
	return err
}
