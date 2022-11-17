// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
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

func (gm *grpcServerMetrics) View() ([]*metrics.View, error) {
	latencyHistgram, err := view.New(
		view.MatchInstrumentName(latencyMetricsName),
		view.WithSetDescription(latencyMetricsDesctiption),
		view.WithSetAggregation(aggregation.ExplicitBucketHistogram{
			Boundaries: metrics.DefaultMillisecondsDistribution,
		}),
	)
	if err != nil {
		return nil, err
	}

	completedRPCCnt, err := view.New(
		view.MatchInstrumentName(completedRPCsMetricsName),
		view.WithSetDescription(completedRPCsMetricsDescription),
		view.WithSetAggregation(aggregation.Sum{}),
	)
	if err != nil {
		return nil, err
	}
	return []*metrics.View{
		&latencyHistgram,
		&completedRPCCnt,
	}, nil
}

func (gm *grpcServerMetrics) Register(m metrics.Meter) error {
	// The metrics are dynamically registered at the grpc server interceptor package,
	// so do nothing in this part
	return nil
}
