package job

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/index/job/correction/service"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	correctedOldIndexCount = "index_job_correction_corrected_old_index_count"
	correctedOldIndexCountDesc = "The number of corrected old indexes while index correction job"

	correctedReplicationCount = "index_job_correction_corrected_replication_count"
	correctedReplicationCountDesc = "The number of operation happend to correct replication number while index correction job"
)

type correctionMetrics struct {
	correction service.Corrector
}

func New(c service.Corrector) metrics.Metric {
	return &correctionMetrics{
		correction: c,
	}
}

func (c *correctionMetrics) View() ([]*metrics.View, error) {
	oldIndexCount, err := view.New(
		view.MatchInstrumentName(correctedOldIndexCount),
		view.WithSetDescription(correctedOldIndexCountDesc),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	replicationCount, err := view.New(
		view.MatchInstrumentName(correctedReplicationCount),
		view.WithSetDescription(correctedReplicationCountDesc),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	return []*metrics.View{
		&oldIndexCount,
		&replicationCount,
	}, nil
}

func (c *correctionMetrics) Register(m metrics.Meter) error {
	oldIndexCount, err := m.AsyncInt64().Gauge(
		correctedOldIndexCount,
		metrics.WithDescription(correctedOldIndexCountDesc),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	replicationCount, err := m.AsyncInt64().Gauge(
		correctedReplicationCount,
		metrics.WithDescription(correctedReplicationCountDesc),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			oldIndexCount,
			replicationCount,
		},
		func(ctx context.Context) {
			oldIndexCount.Observe(ctx, int64(c.correction.NumberOfCorrectedOldIndex()))
			replicationCount.Observe(ctx, int64(c.correction.NumberOfCorrectedReplication()))
		},
	)
}
