// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ngt

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	indexCountMetricsName        = "agent_core_ngt_index_count"
	indexCountMetricsDescription = "Agent NGT index count"

	uncommittedIndexCountMetricsName        = "agent_core_ngt_uncommitted_index_count"
	uncommittedIndexCountMetricsDescription = "Agent NGT index count"

	insertVQueueCountMetricsName        = "agent_core_ngt_insert_vqueue_count"
	insertVQueueCountMetricsDescription = "Agent NGT insert vqueue count"

	deleteVQueueCountMetricsName        = "agent_core_ngt_delete_vqueue_count"
	deleteVQueueCountMetricsDescription = "Agent NGT delete vqueue count"

	completedCreateIndexTotalMetricsName        = "agent_core_ngt_completed_create_index_total"
	completedCreateIndexTotalMetricsDescription = "The cumulative count of completed create index execution"

	executedProactiveGCTotalMetricsName        = "agent_core_ngt_executed_proactive_gc_total"
	executedProactiveGCTotalMetricsDescription = "The cumulative count of proactive GC execution"

	isIndexingMetricsName        = "agent_core_ngt_is_indexing"
	isIndexingMetricsDescription = "Currently indexing or no"

	isSavingMetricsName        = "agent_core_ngt_is_saving"
	isSavingMetricsDescription = "Currently saving or not"

	brokenIndexStoreCountMetricsName        = "agent_core_ngt_broken_index_store_count"
	brokenIndexStoreCountMetricsDescription = "How many broken index generations have been stored"

	kvsRangeDurationMetricsName        = "agent_core_ngt_kvs_range_duration"
	kvsRangeDurationMetricsDescription = "The duration of the kvs range method"
)

type ngtMetrics struct {
	ngt service.NGT
}

func New(n service.NGT) metrics.Metric {
	return &ngtMetrics{
		ngt: n,
	}
}

func (n *ngtMetrics) View() ([]*metrics.View, error) {
	indexCount, err := view.New(
		view.MatchInstrumentName(indexCountMetricsName),
		view.WithSetDescription(indexCountMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	uncommittedIndexCount, err := view.New(
		view.MatchInstrumentName(uncommittedIndexCountMetricsName),
		view.WithSetDescription(uncommittedIndexCountMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	insertVQueueCount, err := view.New(
		view.MatchInstrumentName(insertVQueueCountMetricsName),
		view.WithSetDescription(insertVQueueCountMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	deleteVQueueCount, err := view.New(
		view.MatchInstrumentName(deleteVQueueCountMetricsName),
		view.WithSetDescription(deleteVQueueCountMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	completedCreateIndexTotal, err := view.New(
		view.MatchInstrumentName(completedCreateIndexTotalMetricsName),
		view.WithSetDescription(completedCreateIndexTotalMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	executedProactiveGCTotal, err := view.New(
		view.MatchInstrumentName(executedProactiveGCTotalMetricsName),
		view.WithSetDescription(executedProactiveGCTotalMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	isIndexing, err := view.New(
		view.MatchInstrumentName(isIndexingMetricsName),
		view.WithSetDescription(isIndexingMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	isSaving, err := view.New(
		view.MatchInstrumentName(isSavingMetricsName),
		view.WithSetDescription(isSavingMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	brokenIndexCount, err := view.New(
		view.MatchInstrumentName(brokenIndexStoreCountMetricsName),
		view.WithSetDescription(brokenIndexStoreCountMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	kvsRangeDuration, err := view.New(
		view.MatchInstrumentName(kvsRangeDurationMetricsName),
		view.WithSetDescription(kvsRangeDurationMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	return []*metrics.View{
		&indexCount,
		&uncommittedIndexCount,
		&insertVQueueCount,
		&deleteVQueueCount,
		&completedCreateIndexTotal,
		&executedProactiveGCTotal,
		&isIndexing,
		&isSaving,
		&brokenIndexCount,
		&kvsRangeDuration,
	}, nil
}

func (n *ngtMetrics) Register(m metrics.Meter) error {
	indexCount, err := m.AsyncInt64().Gauge(
		indexCountMetricsName,
		metrics.WithDescription(indexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedIndexCount, err := m.AsyncInt64().Gauge(
		uncommittedIndexCountMetricsName,
		metrics.WithDescription(uncommittedIndexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	insertVQueueCount, err := m.AsyncInt64().Gauge(
		insertVQueueCountMetricsName,
		metrics.WithDescription(insertVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	deleteVQueueCount, err := m.AsyncInt64().Gauge(
		deleteVQueueCountMetricsName,
		metrics.WithDescription(deleteVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	completedCreateIndexTotal, err := m.AsyncInt64().Gauge(
		completedCreateIndexTotalMetricsName,
		metrics.WithDescription(completedCreateIndexTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	executedProactiveGCTotal, err := m.AsyncInt64().Gauge(
		executedProactiveGCTotalMetricsName,
		metrics.WithDescription(executedProactiveGCTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isIndexing, err := m.AsyncInt64().Gauge(
		isIndexingMetricsName,
		metrics.WithDescription(isIndexingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isSaving, err := m.AsyncInt64().Gauge(
		isSavingMetricsName,
		metrics.WithDescription(isSavingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	brokenIndexCount, err := m.AsyncInt64().Gauge(
		brokenIndexStoreCountMetricsName,
		metrics.WithDescription(brokenIndexStoreCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	kvsRangeDuration, err := m.AsyncInt64().Gauge(
		kvsRangeDurationMetricsName,
		metrics.WithDescription(kvsRangeDurationMetricsDescription),
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
			brokenIndexCount,
			kvsRangeDuration,
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
			brokenIndexCount.Observe(ctx, int64(n.ngt.BrokenIndexCount()))
			kvsRangeDuration.Observe(ctx, n.ngt.KvsRangeDuration())
		},
	)
}
