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
package backoff

import (
	"context"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
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

func (*backoffMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        metricsName,
				Description: metricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

func (bm *backoffMetrics) Register(m metrics.Meter) (err error) {
	retryCount, err := m.Int64ObservableGauge(
		metricsName,
		metrics.WithDescription(metricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	_, err = m.RegisterCallback(
		func(ctx context.Context, o api.Observer) error {
			ms := backoff.Metrics(ctx)
			if len(ms) == 0 {
				return nil
			}
			for name, cnt := range ms {
				o.ObserveInt64(retryCount, cnt, api.WithAttributes(attribute.String(bm.backoffNameKey, name)))
			}
			return nil
		}, retryCount,
	)
	return err
}
