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
	"sync/atomic"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type ngt struct {
	insertP      *uint64
	insert       metrics.Int64Measure
	insertVCache metrics.Int64Measure
	deleteVCache metrics.Int64Measure
}

func NewNGTMetrics(i *uint64) metrics.Metric {
	return &ngt{
		insertP:      i,
		insert:       *metrics.Int64("vdaas.org/ngt/insert_count", "number of inserted vectors", metrics.UnitDimensionless),
		insertVCache: *metrics.Int64("vdaas.org/ngt/vcache_insert_count", "number of insert vcache", metrics.UnitDimensionless),
		deleteVCache: *metrics.Int64("vdaas.org/ngt/vcache_delete_count", "number of delete vcache", metrics.UnitDimensionless),
	}
}

func (n *ngt) Measurement() ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		n.insert.M(int64(atomic.LoadUint64(n.insertP))),
		// TODO: implement vcache metrics
		n.insertVCache.M(int64(0)),
		n.deleteVCache.M(int64(0)),
	}, nil
}

func (n *ngt) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "insert_count",
			Description: "number of inserted vectors",
			Measure:     &n.insert,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "vcache_insert_count",
			Description: "number of insert vcache",
			Measure:     &n.insertVCache,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "vcache_delete_count",
			Description: "number of delete vcache",
			Measure:     &n.deleteVCache,
			Aggregation: metrics.LastValue(),
		},
	}
}
