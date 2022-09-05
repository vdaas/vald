package ngt

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
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
		metrics.WithDescription("Agent NGT index count"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedIndexCount, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_uncommitted_index_count",
		metrics.WithDescription("Agent NGT uncommitted index count"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	insertVQueueCount, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_insert_vqueue_count",
		metrics.WithDescription("Agent NGT insert vqueue count"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	deleteVQueueCount, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_delete_vqueue_count",
		metrics.WithDescription("Agent NGT delete vqueue count"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	completedCreateIndexTotal, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_completed_create_index_total",
		metrics.WithDescription("the cumulative count of completed create index execution"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	executedProactiveGCTotal, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_executed_proactive_gc_total",
		metrics.WithDescription("the cumulative count of proactive GC execution"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isIndexing, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_is_indexing",
		metrics.WithDescription("currently indexing or no"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isSaving, err := m.AsyncInt64().UpDownCounter(
		"agent_core_ngt_is_saving",
		metrics.WithDescription("currently saving or not"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
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
