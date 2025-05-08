// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package index

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/manager/index/service"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	UuidCountMetricsName        = "indexer_uuid_count"
	UuidCountMetricsDescription = "UUID count"

	UncommittedUUIDCountMetricsName        = "indexer_uncommitted_uuid_count"
	UncommittedUUIDCountMetricsDescription = "Uncommitted UUID count"

	IsIndexingMetricsName        = "indexer_is_indexing"
	IsIndexingMetricsDescription = "Currently indexing or not"
)

type indexerMetrics struct {
	indexer service.Indexer
}

func New(i service.Indexer) metrics.Metric {
	return &indexerMetrics{
		indexer: i,
	}
}

func (*indexerMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        UuidCountMetricsName,
				Description: UuidCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        UncommittedUUIDCountMetricsName,
				Description: UncommittedUUIDCountMetricsDescription,
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
	}, nil
}

func (im *indexerMetrics) Register(m metrics.Meter) error {
	uuidCount, err := m.Int64ObservableGauge(
		UuidCountMetricsName,
		metrics.WithDescription(UuidCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedUUIDCount, err := m.Int64ObservableGauge(
		UncommittedUUIDCountMetricsName,
		metrics.WithDescription(UncommittedUUIDCountMetricsDescription),
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

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			var indexing int64
			if im.indexer.IsIndexing() {
				indexing = 1
			}
			o.ObserveInt64(uuidCount, int64(im.indexer.NumberOfUUIDs()))
			o.ObserveInt64(uncommittedUUIDCount, int64(im.indexer.NumberOfUncommittedUUIDs()))
			o.ObserveInt64(isIndexing, int64(indexing))
			return nil
		},
		uuidCount,
		uncommittedUUIDCount,
		isIndexing,
	)
	return err
}
