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
package circuitbreaker

import (
	"context"

	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type breakerMetrics struct {
	nameKey  metrics.Key
	stateKey metrics.Key
	state    metrics.Int64Measure
}

func New() (metrics.Metric, error) {
	nameKey, err := metrics.NewKey("name")
	if err != nil {
		return nil, err
	}
	stateKey, err := metrics.NewKey("state")
	if err != nil {
		return nil, err
	}

	return &breakerMetrics{
		nameKey:  nameKey,
		stateKey: stateKey,
		state: *metrics.Int64(
			metrics.ValdOrg+"/circuitbreaker/state",
			"current circuit breaker state",
			metrics.UnitDimensionless,
		),
	}, nil
}

func (*breakerMetrics) Measurement(_ context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (bm *breakerMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	ms := circuitbreaker.Metrics(ctx)
	if len(ms) == 0 {
		return []metrics.MeasurementWithTags{}, nil
	}

	mts := make([]metrics.MeasurementWithTags, 0, len(ms))
	for name, sts := range ms {
		for st, count := range sts {
			mts = append(mts, metrics.MeasurementWithTags{
				Measurement: bm.state.M(count),
				Tags: map[metrics.Key]string{
					bm.nameKey:  name,
					bm.stateKey: st.String(),
				},
			})
		}
	}
	return mts, nil
}

func (bm *breakerMetrics) View() []*metrics.View {
	return []*metrics.View{
		{
			Name:        "circuit_breaker_state",
			Description: bm.state.Description(),
			Measure:     &bm.state,
			TagKeys: []metrics.Key{
				bm.nameKey,
				bm.stateKey,
			},
			Aggregation: metrics.LastValue(),
		},
	}
}
