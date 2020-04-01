//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package ngt provides functions for ngt stats
package ngt

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/ngt/service"
)

type ngtMetrics struct {
	ngt                   service.NGT
	indexCount            metrics.Int64Measure
	uncommittedIndexCount metrics.Int64Measure
	insertVCacheCount     metrics.Int64Measure
	deleteVCacheCount     metrics.Int64Measure
	isIndexing            metrics.Int64Measure
}

func New(n service.NGT) metrics.Metric {
	return &ngtMetrics{
		ngt:                   n,
		indexCount:            *metrics.Int64(metrics.ValdOrg+"/ngt/index_count", "NGT index count", metrics.UnitDimensionless),
		uncommittedIndexCount: *metrics.Int64(metrics.ValdOrg+"/ngt/uncommitted_index_count", "NGT uncommitted index count", metrics.UnitDimensionless),
		insertVCacheCount:     *metrics.Int64(metrics.ValdOrg+"/ngt/insert_vcache_count", "NGT insert vcache count", metrics.UnitDimensionless),
		deleteVCacheCount:     *metrics.Int64(metrics.ValdOrg+"/ngt/delete_vcache_count", "NGT delete vcache count", metrics.UnitDimensionless),
		isIndexing:            *metrics.Int64(metrics.ValdOrg+"/ngt/is_indexing", "currently indexing or not", metrics.UnitDimensionless),
	}
}

func (n *ngtMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	var isIndexing int64
	if n.ngt.IsIndexing() {
		isIndexing = 1
	}

	return []metrics.Measurement{
		n.indexCount.M(int64(n.ngt.Len())),
		n.uncommittedIndexCount.M(int64(n.ngt.InsertVCacheLen() + n.ngt.DeleteVCacheLen())),
		n.insertVCacheCount.M(int64(n.ngt.InsertVCacheLen())),
		n.deleteVCacheCount.M(int64(n.ngt.DeleteVCacheLen())),
		n.isIndexing.M(isIndexing),
	}, nil
}

func (n *ngtMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (n *ngtMetrics) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "ngt_index_count",
			Description: "NGT index count",
			Measure:     &n.indexCount,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "ngt_uncommitted_index_count",
			Description: "NGT uncommitted index count",
			Measure:     &n.uncommittedIndexCount,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "ngt_insert_vcache_count",
			Description: "NGT insert vcache count",
			Measure:     &n.insertVCacheCount,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "ngt_delete_vcache_count",
			Description: "NGT delete vcache count",
			Measure:     &n.deleteVCacheCount,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "ngt_is_indexing",
			Description: "currently indexing or not",
			Measure:     &n.isIndexing,
			Aggregation: metrics.LastValue(),
		},
	}
}
