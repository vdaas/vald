// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
package faiss

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/core/faiss/service"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	IndexCountMetricsName        = "agent_core_faiss_index_count"
	IndexCountMetricsDescription = "Agent Faiss index count"

	UncommittedIndexCountMetricsName        = "agent_core_faiss_uncommitted_index_count"
	UncommittedIndexCountMetricsDescription = "Agent Faiss index count"

	InsertVQueueCountMetricsName        = "agent_core_faiss_insert_vqueue_count"
	InsertVQueueCountMetricsDescription = "Agent Faiss insert vqueue count"

	DeleteVQueueCountMetricsName        = "agent_core_faiss_delete_vqueue_count"
	DeleteVQueueCountMetricsDescription = "Agent Faiss delete vqueue count"

	CompletedCreateIndexTotalMetricsName        = "agent_core_faiss_completed_create_index_total"
	CompletedCreateIndexTotalMetricsDescription = "The cumulative count of completed create index execution"

	ExecutedProactiveGCTotalMetricsName        = "agent_core_faiss_executed_proactive_gc_total"
	ExecutedProactiveGCTotalMetricsDescription = "The cumulative count of proactive GC execution"

	IsIndexingMetricsName        = "agent_core_faiss_is_indexing"
	IsIndexingMetricsDescription = "Currently indexing or no"

	IsSavingMetricsName        = "agent_core_faiss_is_saving"
	IsSavingMetricsDescription = "Currently saving or not"

	TrainCountMetricsName        = "agent_core_faiss_train_count"
	TrainCountMetricsDescription = "Agent Faiss train count"
)

type faissMetrics struct {
	faiss service.Faiss
}

func New(f service.Faiss) metrics.Metric {
	return &faissMetrics{
		faiss: f,
	}
}

func (f *faissMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        IndexCountMetricsName,
				Description: IndexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        UncommittedIndexCountMetricsName,
				Description: UncommittedIndexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        InsertVQueueCountMetricsName,
				Description: InsertVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        DeleteVQueueCountMetricsName,
				Description: DeleteVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        CompletedCreateIndexTotalMetricsName,
				Description: CompletedCreateIndexTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        ExecutedProactiveGCTotalMetricsName,
				Description: ExecutedProactiveGCTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        IsIndexingMetricsName,
				Description: IsIndexingMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        IsSavingMetricsName,
				Description: IsSavingMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        TrainCountMetricsName,
				Description: TrainCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

func (f *faissMetrics) Register(m metrics.Meter) error {
	indexCount, err := m.Int64ObservableGauge(
		IndexCountMetricsName,
		metrics.WithDescription(IndexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedIndexCount, err := m.Int64ObservableGauge(
		UncommittedIndexCountMetricsName,
		metrics.WithDescription(UncommittedIndexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	insertVQueueCount, err := m.Int64ObservableGauge(
		InsertVQueueCountMetricsName,
		metrics.WithDescription(InsertVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	deleteVQueueCount, err := m.Int64ObservableGauge(
		DeleteVQueueCountMetricsName,
		metrics.WithDescription(DeleteVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	completedCreateIndexTotal, err := m.Int64ObservableGauge(
		CompletedCreateIndexTotalMetricsName,
		metrics.WithDescription(CompletedCreateIndexTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	executedProactiveGCTotal, err := m.Int64ObservableGauge(
		ExecutedProactiveGCTotalMetricsName,
		metrics.WithDescription(ExecutedProactiveGCTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isIndexing, err := m.Int64ObservableGauge(
		IsIndexingMetricsName,
		metrics.WithDescription(IsIndexingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isSaving, err := m.Int64ObservableGauge(
		IsSavingMetricsName,
		metrics.WithDescription(IsSavingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	trainCount, err := m.Int64ObservableGauge(
		TrainCountMetricsName,
		metrics.WithDescription(TrainCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			var indexing int64
			if f.faiss.IsIndexing() {
				indexing = 1
			}
			var saving int64
			if f.faiss.IsSaving() {
				saving = 1
			}

			o.ObserveInt64(indexCount, int64(f.faiss.Len()))
			o.ObserveInt64(uncommittedIndexCount, int64(f.faiss.InsertVQueueBufferLen()+f.faiss.DeleteVQueueBufferLen()))
			o.ObserveInt64(insertVQueueCount, int64(f.faiss.InsertVQueueBufferLen()))
			o.ObserveInt64(deleteVQueueCount, int64(int64(f.faiss.DeleteVQueueBufferLen())))
			o.ObserveInt64(completedCreateIndexTotal, int64(f.faiss.NumberOfCreateIndexExecution()))
			o.ObserveInt64(executedProactiveGCTotal, int64(f.faiss.NumberOfProactiveGCExecution()))
			o.ObserveInt64(isIndexing, int64(indexing))
			o.ObserveInt64(isSaving, int64(saving))
			o.ObserveInt64(trainCount, int64(f.faiss.GetTrainSize()))

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
		trainCount,
	)
	return err
}
