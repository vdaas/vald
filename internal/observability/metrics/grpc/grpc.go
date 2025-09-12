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
	"github.com/vdaas/vald/internal/observability/metrics"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	LatencyMetricsName        = "server_latency"
	LatencyMetricsDescription = "Server latency in milliseconds, by method"

	CompletedRPCsMetricsName        = "server_completed_rpcs"
	CompletedRPCsMetricsDescription = "Count of RPCs by method and status"

	PoolConnMetricsName        = "server_pool_conn"
	PoolConnMetricsDescription = "Count of healthy pool connections by target address"
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
				Name:        LatencyMetricsName,
				Description: LatencyMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationExplicitBucketHistogram{
					Boundaries: metrics.DefaultMillisecondsDistribution,
				},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        CompletedRPCsMetricsName,
				Description: CompletedRPCsMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationSum{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        PoolConnMetricsName,
				Description: PoolConnMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationSum{},
			},
		),
	}, nil
}

func (gm *grpcServerMetrics) Register(m metrics.Meter) error {
	// The metrics are dynamically registered at the grpc server interceptor package,
	// so do nothing in this part
	return nil
}
