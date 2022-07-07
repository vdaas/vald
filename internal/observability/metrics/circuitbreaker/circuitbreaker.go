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
	nameKey   metrics.Key
	isOpening metrics.Int64Measure
}

func New() (metrics.Metric, error) {
	key, err := metrics.NewKey("breaker_name")
	if err != nil {
		return nil, err
	}

	return &breakerMetrics{
		nameKey: key,
		isOpening: *metrics.Int64(
			metrics.ValdOrg+"/circuitbreaker/is_opening",
			"currently breaker state is open or not",
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
	for name, st := range ms {
		opening := 0
		if st == circuitbreaker.StateOpen {
			opening = 1
		}
		mts = append(mts, metrics.MeasurementWithTags{
			Measurement: bm.isOpening.M(int64(opening)),
			Tags: map[metrics.Key]string{
				bm.nameKey: name,
			},
		})
	}
	return mts, nil
}

func (bm *breakerMetrics) View() []*metrics.View {
	return []*metrics.View{
		{
			Name:        "breaker_is_opening",
			Description: bm.isOpening.Description(),
			Measure:     &bm.isOpening,
			TagKeys: []metrics.Key{
				bm.nameKey,
			},
			Aggregation: metrics.LastValue(),
		},
	}
}
