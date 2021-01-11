//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package index provides functions for indexer stats
package index

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/manager/index/service"
)

type indexerMetrics struct {
	indexer              service.Indexer
	uuidCount            metrics.Int64Measure
	uncommittedUUIDCount metrics.Int64Measure
	isIndexing           metrics.Int64Measure
}

func New(i service.Indexer) metrics.Metric {
	return &indexerMetrics{
		indexer: i,
		uuidCount: *metrics.Int64(
			metrics.ValdOrg+"/manager/index/uuid_count",
			"UUID count",
			metrics.UnitDimensionless),
		uncommittedUUIDCount: *metrics.Int64(
			metrics.ValdOrg+"/manager/index/uncommitted_uuid_count",
			"uncommitted UUID count",
			metrics.UnitDimensionless),
		isIndexing: *metrics.Int64(
			metrics.ValdOrg+"/manager/index/is_indexing",
			"currently indexing or not",
			metrics.UnitDimensionless),
	}
}

func (i *indexerMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	var isIndexing int64
	if i.indexer.IsIndexing() {
		isIndexing = 1
	}

	return []metrics.Measurement{
		i.uuidCount.M(int64(i.indexer.NumberOfUUIDs())),
		i.uncommittedUUIDCount.M(int64(i.indexer.NumberOfUncommittedUUIDs())),
		i.isIndexing.M(isIndexing),
	}, nil
}

func (i *indexerMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (i *indexerMetrics) View() []*metrics.View {
	return []*metrics.View{
		{
			Name:        "indexer_uuid_count",
			Description: i.uuidCount.Description(),
			Measure:     &i.uuidCount,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "indexer_uncommitted_uuid_count",
			Description: i.uncommittedUUIDCount.Description(),
			Measure:     &i.uncommittedUUIDCount,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "indexer_is_indexing",
			Description: i.isIndexing.Description(),
			Measure:     &i.isIndexing,
			Aggregation: metrics.LastValue(),
		},
	}
}
