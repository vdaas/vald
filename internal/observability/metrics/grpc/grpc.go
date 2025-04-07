// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package grpc

import (
	"context"

	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	latencyMetricsName        = "server_latency"
	latencyMetricsDesctiption = "Server latency in milliseconds, by method"

	completedRPCsMetricsName        = "server_completed_rpcs"
	completedRPCsMetricsDescription = "Count of RPCs by method and status"

	poolConnMetricsName        = "server_pool_conn"
	poolConnMetricsDescription = "Count of healthy pool connections by target address"
)

type grpcServerMetrics struct {
	poolTargetAddrKey string
}

func New() metrics.Metric {
	return &grpcServerMetrics{
		"target_address",
	}
}

func (*grpcServerMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        latencyMetricsName,
				Description: latencyMetricsDesctiption,
			},
			view.Stream{
				Aggregation: view.AggregationExplicitBucketHistogram{
					Boundaries: metrics.DefaultMillisecondsDistribution,
				},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        completedRPCsMetricsName,
				Description: completedRPCsMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationSum{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        poolConnMetricsName,
				Description: poolConnMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationSum{},
			},
		),
	}, nil
}

func (gm *grpcServerMetrics) Register(m metrics.Meter) error {
	// The metrics are dynamically registered at the grpc server interceptor package,
	healthyConn, err := m.Int64ObservableGauge(
		poolConnMetricsName,
		metrics.WithDescription(poolConnMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	_, err = m.RegisterCallback(
		func(ctx context.Context, o api.Observer) error {
			ms := pool.Metrics(ctx)
			if len(ms) == 0 {
				return nil
			}
			for name, cnt := range ms {
				o.ObserveInt64(healthyConn, cnt, api.WithAttributes(attribute.String(gm.poolTargetAddrKey, name)))
			}
			return nil
		}, healthyConn,
	)
	return err
}
