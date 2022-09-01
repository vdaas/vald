package backoff

import (
	"context"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/observability-v2/attribute"
	"github.com/vdaas/vald/internal/observability-v2/metrics"
)

type backoffMetrics struct {
	backoffNameKey string
}

func New() metrics.Metric {
	return &backoffMetrics{
		backoffNameKey: "backoff_name",
	}
}

func (bm *backoffMetrics) Register(m metrics.Meter) error {
	retryCount, err := m.AsyncInt64().Gauge(
		"backoff_retry_count",
		metrics.WithDescription("Backoff retry count"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			retryCount,
		},
		func(ctx context.Context) {
			ms := backoff.Metrics(ctx)
			if len(ms) == 0 {
				return
			}
			for name, cnt := range ms {
				retryCount.Observe(ctx, cnt, attribute.String(bm.backoffNameKey, name))
			}
		},
	)
}
