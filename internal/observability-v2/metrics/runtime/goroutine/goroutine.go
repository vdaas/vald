package goroutine

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability-v2/metrics"
	"go.opentelemetry.io/otel/metric/instrument"
)

type goroutine struct {
	name        string
	description string
	unit        metrics.Unit
}

func New() metrics.Metric {
	return &goroutine{
		name:        metrics.ValdOrg + "/runtime/goroutine_count",
		description: "number of goroutines",
		unit:        metrics.Dimensionless,
	}
}

func (g *goroutine) Name() string {
	return g.name
}

func (g *goroutine) Description() string {
	return g.description
}

func (g *goroutine) Unit() metrics.Unit {
	return g.unit
}

func (g *goroutine) Register(m metrics.Meter) error {
	conter, err := m.AsyncInt64().UpDownCounter(
		g.Name(),
		instrument.WithDescription(g.Description()),
		instrument.WithUnit(g.Unit()),
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
