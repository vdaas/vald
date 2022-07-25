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
package backoff

import (
	"context"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type backoffMetrics struct {
	nameKey    metrics.Key
	retryCount metrics.Int64Measure
}

func New() (metrics.Metric, error) {
	key, err := metrics.NewKey("backoff_name")
	if err != nil {
		return nil, err
	}

	return &backoffMetrics{
		nameKey: key,
		retryCount: *metrics.Int64(
			metrics.ValdOrg+"/backoff/retry_count",
			"Backoff retry count",
			metrics.UnitDimensionless),
	}, nil
}

func (*backoffMetrics) Measurement(_ context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (bm *backoffMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	ms := backoff.Metrics(ctx)
	if len(ms) == 0 {
		return []metrics.MeasurementWithTags{}, nil
	}

	mts := make([]metrics.MeasurementWithTags, 0, len(ms))
	for name, cnt := range ms {
		mts = append(mts, metrics.MeasurementWithTags{
			Measurement: bm.retryCount.M(cnt),
			Tags: map[metrics.Key]string{
				bm.nameKey: name,
			},
		})
	}
	return mts, nil
}

func (bm *backoffMetrics) View() []*metrics.View {
	return []*metrics.View{
		{
			Name:        "backoff_retry_count",
			Description: bm.retryCount.Description(),
			Measure:     &bm.retryCount,
			TagKeys: []metrics.Key{
				bm.nameKey,
			},
			Aggregation: metrics.LastValue(),
		},
	}
}
