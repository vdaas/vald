// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/manager/index/service"
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

func (im *indexerMetrics) View() ([]*metrics.View, error) {
	uuidCount, err := view.New(
		view.MatchInstrumentName(uuidCountMetricsName),
		view.WithSetDescription(uuidCountMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	uncommittedUUIDCount, err := view.New(
		view.MatchInstrumentName(uncommittedUUIDCountMetricsName),
		view.WithSetDescription(uncommittedUUIDCountMetricsDescription),
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

	return []*metrics.View{
		&uuidCount,
		&uncommittedUUIDCount,
		&isIndexing,
	}, nil
}

func (im *indexerMetrics) Register(m metrics.Meter) error {
	uuidCount, err := m.AsyncInt64().Gauge(
		uuidCountMetricsName,
		metrics.WithDescription(uuidCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedUUIDCount, err := m.AsyncInt64().Gauge(
		uncommittedUUIDCountMetricsName,
		metrics.WithDescription(uncommittedUUIDCountMetricsDescription),
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

	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
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
