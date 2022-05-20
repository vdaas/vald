package backoff

import (
	"context"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type backoffMetrics struct {
	bo         backoff.Backoff
	retryCount metrics.Int64Measure
}

func New(bo backoff.Backoff) metrics.Metric {
	return &backoffMetrics{
		bo: bo,
		retryCount: *metrics.Int64(
			metrics.ValdOrg+"/grpc/backoff/retry_count",
			"Backoff retry count",
			metrics.UnitDimensionless),
	}
}

func (bm *backoffMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return nil, nil
}

func (bm *backoffMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return nil, nil
}

func (bm *backoffMetrics) View() []*metrics.View {
	return nil
}
