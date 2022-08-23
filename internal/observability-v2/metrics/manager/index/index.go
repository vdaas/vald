package index

import (
	"context"

	"github.com/vdaas/vald/internal/observability-v2/metrics"
	"github.com/vdaas/vald/pkg/manager/index/service"
	"go.opentelemetry.io/otel/metric/instrument"
)

type indexerMetrics struct {
	indexer service.Indexer
}

func New(i service.Indexer) metrics.Metric {
	return &indexerMetrics{
		indexer: i,
	}
}

func (im *indexerMetrics) Register(m metrics.Meter) error {
	uuidCount, err := m.AsyncInt64().Gauge(
		"indexer_uuid_count",
		instrument.WithDescription("UUID count"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedUUIDCount, err := m.AsyncInt64().Gauge(
		"indexer_uncommitted_uuid_count",
		instrument.WithDescription("uncommitted UUID count"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isIndexing, err := m.AsyncInt64().Gauge(
		"indexer_is_indexing",
		instrument.WithDescription("currently indexing or not"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]instrument.Asynchronous{
			uuidCount,
			uncommittedUUIDCount,
			isIndexing,
		},
		func(ctx context.Context) {
			var indexing int64
			if im.indexer.IsIndexing() {
				indexing = 1
			}
			uuidCount.Observe(ctx, int64(im.indexer.NumberOfUUIDs()))
			uncommittedUUIDCount.Observe(ctx, int64(im.indexer.NumberOfUncommittedUUIDs()))
			isIndexing.Observe(ctx, int64(indexing))
		},
	)
}
