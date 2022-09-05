package goroutine

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type goroutine struct {
}

func New() metrics.Metric {
	return &goroutine{}
}

func (g *goroutine) Register(m metrics.Meter) error {
	conter, err := m.AsyncInt64().Gauge(
		"goroutine_count",
		metrics.WithDescription("number of goroutines"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			conter,
		},
		func(ctx context.Context) {
			conter.Observe(ctx, int64(runtime.NumGoroutine()))
		},
	)
}
