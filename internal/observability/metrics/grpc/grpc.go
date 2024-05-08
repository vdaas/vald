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
package grpc

import (
	"github.com/vdaas/vald/internal/observability/metrics"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	latencyMetricsName        = "server_latency"
	latencyMetricsDesctiption = "Server latency in milliseconds, by method"

	completedRPCsMetricsName        = "server_completed_rpcs"
	completedRPCsMetricsDescription = "Count of RPCs by method and status"
)

type grpcServerMetrics struct{}

func New() metrics.Metric {
	return &grpcServerMetrics{}
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
	}, nil
}

func (*grpcServerMetrics) Register(metrics.Meter) error {
	// The metrics are dynamically registered at the grpc server interceptor package,
	// so do nothing in this part
	return nil
}
