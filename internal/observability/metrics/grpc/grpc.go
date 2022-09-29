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
