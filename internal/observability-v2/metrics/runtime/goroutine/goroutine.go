package goroutine

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability-v2/metrics"
	"go.opentelemetry.io/otel/metric/instrument"
)

type goroutine struct {
}

func New() metrics.Metric {
	return &goroutine{}
}

func (g *goroutine) Register(m metrics.Meter) error {
	conter, err := m.AsyncInt64().Gauge(
		"goroutine_count",
		instrument.WithDescription("number of goroutines"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]instrument.Asynchronous{
			conter,
		},
		func(ctx context.Context) {
			conter.Observe(ctx, int64(runtime.NumGoroutine()))
		},
	)
}
