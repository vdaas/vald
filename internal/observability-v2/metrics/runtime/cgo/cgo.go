package cgo

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability-v2/metrics"
	"go.opentelemetry.io/otel/metric/instrument"
)

type cgo struct {
}

func New() metrics.Metric {
	return &cgo{}
}

func (c *cgo) Register(m metrics.Meter) error {
	conter, err := m.AsyncInt64().UpDownCounter(
		"cgo_call_count",
		instrument.WithDescription("number of cgo call"),
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
