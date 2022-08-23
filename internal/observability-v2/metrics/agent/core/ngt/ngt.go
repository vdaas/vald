package ngt

import (
	"context"

	"github.com/vdaas/vald/internal/observability-v2/metrics"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
	"go.opentelemetry.io/otel/metric/instrument"
)

type ngtMetrics struct {
	ngt service.NGT
}

func New(n service.NGT) metrics.Metric {
	return &ngtMetrics{
		ngt: n,
	}
}

func (n *ngtMetrics) Register(m metrics.Meter) error {
	indexCount, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_index_count",
		instrument.WithDescription("Agent NGT index count"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedIndexCount, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_uncommitted_index_count",
		instrument.WithDescription("Agent NGT uncommitted index count"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	insertVQueueCount, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_insert_vqueue_count",
		instrument.WithDescription("Agent NGT insert vqueue count"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	deleteVQueueCount, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_delete_vqueue_count",
		instrument.WithDescription("Agent NGT delete vqueue count"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	completedCreateIndexTotal, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_completed_create_index_total",
		instrument.WithDescription("the cumulative count of completed create index execution"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	executedProactiveGCTotal, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_executed_proactive_gc_total",
		instrument.WithDescription("the cumulative count of proactive GC execution"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isIndexing, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_is_indexing",
		instrument.WithDescription("currently indexing or no"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isSaving, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_is_saving",
		instrument.WithDescription("currently saving or not"),
		instrument.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]instrument.Asynchronous{
			indexCount,
			uncommittedIndexCount,
			insertVQueueCount,
			deleteVQueueCount,
			completedCreateIndexTotal,
			executedProactiveGCTotal,
			isIndexing,
			isSaving,
		},
		func(ctx context.Context) {
			var indexing int64
			if n.ngt.IsIndexing() {
				indexing = 1
			}

			var saving int64
			if n.ngt.IsSaving() {
				saving = 1
			}

			indexCount.Observe(ctx, int64(n.ngt.Len()))
			uncommittedIndexCount.Observe(ctx, int64(n.ngt.InsertVQueueBufferLen()+n.ngt.DeleteVQueueBufferLen()))
			insertVQueueCount.Observe(ctx, int64(n.ngt.InsertVQueueBufferLen()))
			deleteVQueueCount.Observe(ctx, int64(int64(n.ngt.DeleteVQueueBufferLen())))
			completedCreateIndexTotal.Observe(ctx, int64(n.ngt.NumberOfCreateIndexExecution()))
			executedProactiveGCTotal.Observe(ctx, int64(n.ngt.NumberOfProactiveGCExecution()))
			isIndexing.Observe(ctx, int64(indexing))
			isSaving.Observe(ctx, int64(saving))
		},
	)
}
