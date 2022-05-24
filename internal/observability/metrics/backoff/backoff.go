package backoff

import (
	"context"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type backoffMetrics struct {
	bo         backoff.Backoff
	nameKey    metrics.Key
	retryCount metrics.Int64Measure
}

func New(bo backoff.Backoff) (metrics.Metric, error) {
	key, err := metrics.NewKey("backoff_name")
	if err != nil {
		return nil, err
	}

	return &backoffMetrics{
		bo:      bo,
		nameKey: key,
		retryCount: *metrics.Int64(
			metrics.ValdOrg+"/backoff/retry_count",
			"Backoff retry count",
			metrics.UnitDimensionless),
	}, nil
}

func (*backoffMetrics) Measurement(_ context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{}, nil
}

func (bm *backoffMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	ms := bm.bo.Metrics(ctx)
	mts := make([]metrics.MeasurementWithTags, 0, len(ms))
	for name, cnt := range ms {
		mts = append(mts, metrics.MeasurementWithTags{
			Measurement: bm.retryCount.M(int64(cnt)),
			Tags: map[metrics.Key]string{
				bm.nameKey: name,
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
				bm.nameKey,
			},
			Aggregation: metrics.LastValue(),
		},
	}
}
