package backoff

import (
	"context"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/observability/metrics"
)

var (
	serviceNameKey = metrics.MustNewKey("grpc_service")
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
	return []metrics.Measurement{}, nil
}

func (bm *backoffMetrics) MeasurementWithTags(ctx context.Context) (mts []metrics.MeasurementWithTags, err error) {
	ms := bm.bo.Metrics(ctx)
	mts = make([]metrics.MeasurementWithTags, 0, len(ms))
	for svc, val := range ms {
		mts = append(mts, metrics.MeasurementWithTags{
			Measurement: bm.retryCount.M(int64(val)),
			Tags: map[metrics.Key]string{
				serviceNameKey: svc,
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
				serviceNameKey,
			},
			Aggregation: metrics.Count(),
		},
	}
}
