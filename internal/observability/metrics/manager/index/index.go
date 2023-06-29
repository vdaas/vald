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
package index

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/manager/index/service"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

const (
	uuidCountMetricsName        = "indexer_uuid_count"
	uuidCountMetricsDescription = "UUID count"

	uncommittedUUIDCountMetricsName        = "indexer_uncommitted_uuid_count"
	uncommittedUUIDCountMetricsDescription = "Uncommitted UUID count"

	isIndexingMetricsName        = "indexer_is_indexing"
	isIndexingMetricsDescription = "Currently indexing or not"
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
				Name:        uuidCountMetricsName,
				Description: uuidCountMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        uncommittedUUIDCountMetricsName,
				Description: uncommittedUUIDCountMetricsDescription,
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
	}, nil
}

func (im *indexerMetrics) Register(m metrics.Meter) error {
	uuidCount, err := m.Int64ObservableGauge(
		uuidCountMetricsName,
		metrics.WithDescription(uuidCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedUUIDCount, err := m.Int64ObservableGauge(
		uncommittedUUIDCountMetricsName,
		metrics.WithDescription(uncommittedUUIDCountMetricsDescription),
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
	return nil
}
