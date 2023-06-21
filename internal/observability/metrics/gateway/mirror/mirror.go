package mirror

import (
	"context"

	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	metricsName        = "gateway_mirror_connecting_target"
	metricsDescription = "Target to which the mirror gateway is connecting"

	targetAddrKey = "addr"
)

type mirrorMetrics struct {
	mirr service.Mirror
}

func New(mirr service.Mirror) metrics.Metric {
	return &mirrorMetrics{
		mirr: mirr,
	}
}

func (*mirrorMetrics) View() ([]*metrics.View, error) {
	target, err := view.New(
		view.MatchInstrumentName(metricsName),
		view.WithSetDescription(metricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}
	return []*metrics.View{
		&target,
	}, nil
}

func (mm *mirrorMetrics) Register(m metrics.Meter) error {
	targetGauge, err := m.AsyncInt64().Gauge(
		metricsName,
		metrics.WithDescription(metricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			targetGauge,
		},
		func(ctx context.Context) {
			mm.mirr.RangeAllMirrorAddr(func(addr string, _ any) bool {
				targetGauge.Observe(ctx, 1, attribute.String(targetAddrKey, addr))
				return true
			})
		},
	)
}
