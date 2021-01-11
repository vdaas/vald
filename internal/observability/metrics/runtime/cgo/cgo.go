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

// Package cgo provides functions for runtime cgo stats
package cgo

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type cgo struct {
	count metrics.Int64Measure
}

func New() metrics.Metric {
	return &cgo{
		count: *metrics.Int64(metrics.ValdOrg+"/runtime/cgo_call_count", "number of cgo call", metrics.UnitDimensionless),
	}
}

func (c *cgo) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		c.count.M(int64(runtime.NumCgoCall())),
	}, nil
}

func (c *cgo) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (c *cgo) View() []*metrics.View {
	return []*metrics.View{
		{
			Name:        "cgo_call_count",
			Description: c.count.Description(),
			Measure:     &c.count,
			Aggregation: metrics.LastValue(),
		},
	}
}
