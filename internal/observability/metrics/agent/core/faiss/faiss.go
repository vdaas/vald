// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

const (
	indexCountMetricsName        = "agent_core_faiss_index_count"
	indexCountMetricsDescription = "Agent Faiss index count"

	uncommittedIndexCountMetricsName        = "agent_core_faiss_uncommitted_index_count"
	uncommittedIndexCountMetricsDescription = "Agent Faiss index count"

	insertVQueueCountMetricsName        = "agent_core_faiss_insert_vqueue_count"
	insertVQueueCountMetricsDescription = "Agent Faiss insert vqueue count"

	deleteVQueueCountMetricsName        = "agent_core_faiss_delete_vqueue_count"
	deleteVQueueCountMetricsDescription = "Agent Faiss delete vqueue count"

	completedCreateIndexTotalMetricsName        = "agent_core_faiss_completed_create_index_total"
	completedCreateIndexTotalMetricsDescription = "The cumulative count of completed create index execution"

	executedProactiveGCTotalMetricsName        = "agent_core_faiss_executed_proactive_gc_total"
	executedProactiveGCTotalMetricsDescription = "The cumulative count of proactive GC execution"

	isIndexingMetricsName        = "agent_core_faiss_is_indexing"
	isIndexingMetricsDescription = "Currently indexing or no"

	isSavingMetricsName        = "agent_core_faiss_is_saving"
	isSavingMetricsDescription = "Currently saving or not"

	trainCountMetricsName        = "agent_core_faiss_train_count"
	trainCountMetricsDescription = "Agent Faiss train count"
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
				Name:        indexCountMetricsName,
				Description: indexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        uncommittedIndexCountMetricsName,
				Description: uncommittedIndexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        insertVQueueCountMetricsName,
				Description: insertVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        deleteVQueueCountMetricsName,
				Description: deleteVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        completedCreateIndexTotalMetricsName,
				Description: completedCreateIndexTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        executedProactiveGCTotalMetricsName,
				Description: executedProactiveGCTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        isIndexingMetricsName,
				Description: isIndexingMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        isSavingMetricsName,
				Description: isSavingMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        trainCountMetricsName,
				Description: trainCountMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
	}, nil
}

func (f *faissMetrics) Register(m metrics.Meter) error {
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

	trainCount, err := m.Int64ObservableGauge(
		trainCountMetricsName,
		metrics.WithDescription(trainCountMetricsDescription),
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
