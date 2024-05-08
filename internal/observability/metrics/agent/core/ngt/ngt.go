// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
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
)

type ngtMetrics struct {
	ngt service.NGT
}

func New(n service.NGT) metrics.Metric {
	return &ngtMetrics{
		ngt: n,
	}
}

func (n *ngtMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        indexCountMetricsName,
				Description: indexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        uncommittedIndexCountMetricsName,
				Description: uncommittedIndexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        insertVQueueCountMetricsName,
				Description: insertVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        deleteVQueueCountMetricsName,
				Description: deleteVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        completedCreateIndexTotalMetricsName,
				Description: completedCreateIndexTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        executedProactiveGCTotalMetricsName,
				Description: executedProactiveGCTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        isIndexingMetricsName,
				Description: isIndexingMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        isSavingMetricsName,
				Description: isSavingMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        brokenIndexStoreCountMetricsName,
				Description: brokenIndexStoreCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

func (n *ngtMetrics) Register(m metrics.Meter) error {
	indexCount, err := m.Int64ObservableGauge(
		indexCountMetricsName,
		metrics.WithDescription(indexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedIndexCount, err := m.Int64ObservableGauge(
		uncommittedIndexCountMetricsName,
		metrics.WithDescription(uncommittedIndexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	insertVQueueCount, err := m.Int64ObservableGauge(
		insertVQueueCountMetricsName,
		metrics.WithDescription(insertVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	deleteVQueueCount, err := m.Int64ObservableGauge(
		deleteVQueueCountMetricsName,
		metrics.WithDescription(deleteVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	completedCreateIndexTotal, err := m.Int64ObservableGauge(
		completedCreateIndexTotalMetricsName,
		metrics.WithDescription(completedCreateIndexTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	executedProactiveGCTotal, err := m.Int64ObservableGauge(
		executedProactiveGCTotalMetricsName,
		metrics.WithDescription(executedProactiveGCTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isIndexing, err := m.Int64ObservableGauge(
		isIndexingMetricsName,
		metrics.WithDescription(isIndexingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isSaving, err := m.Int64ObservableGauge(
		isSavingMetricsName,
		metrics.WithDescription(isSavingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	brokenIndexCount, err := m.Int64ObservableGauge(
		brokenIndexStoreCountMetricsName,
		metrics.WithDescription(brokenIndexStoreCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			var indexing int64
			if n.ngt.IsIndexing() {
				indexing = 1
			}
			var saving int64
			if n.ngt.IsSaving() {
				saving = 1
			}
			o.ObserveInt64(indexCount, int64(n.ngt.Len()))
			o.ObserveInt64(uncommittedIndexCount, int64(n.ngt.InsertVQueueBufferLen()+n.ngt.DeleteVQueueBufferLen()))
			o.ObserveInt64(insertVQueueCount, int64(n.ngt.InsertVQueueBufferLen()))
			o.ObserveInt64(deleteVQueueCount, int64(int64(n.ngt.DeleteVQueueBufferLen())))
			o.ObserveInt64(completedCreateIndexTotal, int64(n.ngt.NumberOfCreateIndexExecution()))
			o.ObserveInt64(executedProactiveGCTotal, int64(n.ngt.NumberOfProactiveGCExecution()))
			o.ObserveInt64(isIndexing, int64(indexing))
			o.ObserveInt64(isSaving, int64(saving))
			o.ObserveInt64(brokenIndexCount, int64(n.ngt.BrokenIndexCount()))
			return nil
		},
		indexCount,
		uncommittedIndexCount,
		insertVQueueCount,
		deleteVQueueCount,
		completedCreateIndexTotal,
		executedProactiveGCTotal,
		isIndexing,
		isSaving,
		brokenIndexCount,
	)
	return err
}
