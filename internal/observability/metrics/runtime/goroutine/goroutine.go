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

// Package goroutine provides functions for goroutine runtime stats
package goroutine

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type goroutines struct {
	count metrics.Int64Measure
}

func New() metrics.Metric {
	return &goroutines{
		count: *metrics.Int64(metrics.ValdOrg+"/runtime/goroutine_count", "number of goroutines", metrics.UnitDimensionless),
	}
}

func (g *goroutines) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		g.count.M(int64(runtime.NumGoroutine())),
	}, nil
}

func (g *goroutines) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (g *goroutines) View() []*metrics.View {
	return []*metrics.View{
		{
			Name:        "goroutine_count",
			Description: g.count.Description(),
			Measure:     &g.count,
			Aggregation: metrics.LastValue(),
		},
	}
}
