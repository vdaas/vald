// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package mirror

import (
	"context"

	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	metricsName        = "gateway_mirror_connecting_target"
	metricsDescription = "Target to which the mirror gateway is connecting"

	targetAddrKey = "addr"
)

type mirrorMetrics struct {
	mirr service.Mirror
}

func New(mirr service.Mirror) metrics.Metric {
	return &mirrorMetrics{
		mirr: mirr,
	}
}

func (*mirrorMetrics) View() ([]*metrics.View, error) {
	target, err := view.New(
		view.MatchInstrumentName(metricsName),
		view.WithSetDescription(metricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}
	return []*metrics.View{
		&target,
	}, nil
}

func (mm *mirrorMetrics) Register(m metrics.Meter) error {
	targetGauge, err := m.AsyncInt64().Gauge(
		metricsName,
		metrics.WithDescription(metricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			targetGauge,
		},
		func(ctx context.Context) {
			mm.mirr.RangeAllMirrorAddr(func(addr string, _ any) bool {
				targetGauge.Observe(ctx, 1, attribute.String(targetAddrKey, addr))
				return true
			})
		},
	)
}
