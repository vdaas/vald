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
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	metricsName        = "backoff_retry_count"
	metricsDescription = "Backoff retry count"
)

type backoffMetrics struct {
	backoffNameKey string
}

func New() metrics.Metric {
	return &backoffMetrics{
		backoffNameKey: "backoff_name",
	}
}

func (bm *backoffMetrics) View() ([]*metrics.View, error) {
	retryCount, err := view.New(
		view.MatchInstrumentName(metricsName),
		view.WithSetDescription(metricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}
	return []*metrics.View{
		&retryCount,
	}, nil
}

func (bm *backoffMetrics) Register(m metrics.Meter) error {
	retryCount, err := m.AsyncInt64().Gauge(
		metricsName,
		metrics.WithDescription(metricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			retryCount,
		},
		func(ctx context.Context) {
			ms := backoff.Metrics(ctx)
			if len(ms) == 0 {
				return
			}
			for name, cnt := range ms {
				retryCount.Observe(ctx, cnt, attribute.String(bm.backoffNameKey, name))
			}
		},
	)
}
