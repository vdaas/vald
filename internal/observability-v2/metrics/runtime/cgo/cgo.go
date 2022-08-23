package cgo

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability-v2/metrics"
	"go.opentelemetry.io/otel/metric/instrument"
)

type cgo struct {
	name        string
	description string
	unit        metrics.Unit
}

func New() metrics.Metric {
	return &cgo{
		name:        metrics.ValdOrg + "/runtime/cgo_call_count",
		description: "number of cgo call",
		unit:        metrics.Dimensionless,
	}
}

func (c *cgo) Register(m metrics.Meter) error {
	conter, err := m.AsyncInt64().UpDownCounter(
		"",
		instrument.WithDescription(""),
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
