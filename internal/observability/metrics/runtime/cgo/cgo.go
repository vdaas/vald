package cgo

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type cgo struct {
}

func New() metrics.Metric {
	return &cgo{}
}

func (c *cgo) Register(m metrics.Meter) error {
	conter, err := m.AsyncInt64().UpDownCounter(
		"cgo_call_count",
		metrics.WithDescription("number of cgo call"),
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
