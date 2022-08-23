package circuitbreaker

import (
	"context"

	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/observability-v2/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
)

type breakerMetrics struct {
	breakerNameKey string
	stateKey       string
}

func New() metrics.Metric {
	return &breakerMetrics{
		breakerNameKey: "name",
		stateKey:       "state",
	}
}

func (bm *breakerMetrics) Register(m metrics.Meter) error {
	breakerState, err := m.AsyncInt64().Gauge(
		"circuit_breaker_state",
		instrument.WithDescription("current circuit breaker state"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]instrument.Asynchronous{
			breakerState,
		},
		func(ctx context.Context) {
			ms := circuitbreaker.Metrics(ctx)
			if len(ms) == 0 {
				return
			}
			for name, sts := range ms {
				for st, cnt := range sts {
					breakerState.Observe(ctx, cnt,
						attribute.String(bm.breakerNameKey, name),
						attribute.String(bm.stateKey, st.String()),
					)
				}
			}
		},
	)
}
