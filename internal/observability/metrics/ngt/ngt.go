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
	ic                    *uint64
	uncommittedIndexCount metrics.Int64Measure
}

func NewNGTMetrics(ic *uint64) metrics.Metric {
	return &ngt{
		ic:                    ic,
		uncommittedIndexCount: *metrics.Int64("vdaas.org/ngt/uncommitted_index_count", "uncommitted index count", metrics.UnitDimensionless),
	}
}

func (n *ngt) Measurement() ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		n.uncommittedIndexCount.M(int64(atomic.LoadUint64(n.ic))),
	}, nil
}

func (n *ngt) MeasurementWithTags() ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (n *ngt) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "uncommitted_index_count",
			Description: "uncommitted index count",
			Measure:     &n.uncommittedIndexCount,
			Aggregation: metrics.LastValue(),
		},
	}
}
