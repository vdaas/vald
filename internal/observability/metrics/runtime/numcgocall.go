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

// Package runtime provides functions for runtime stats
package runtime

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type numcgocall struct {
	count metrics.Int64Measure
}

func NewNumberOfCGOCall() metrics.Metric {
	return &numcgocall{
		count: *metrics.Int64("vdaas.org/vald/runtime/cgo_call_count", "number of cgo call", metrics.UnitDimensionless),
	}
}

func (n *numcgocall) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		n.count.M(int64(runtime.NumCgoCall())),
	}, nil
}

func (n *numcgocall) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (n *numcgocall) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "cgo_call_count",
			Description: "number of cgo call",
			Measure:     &n.count,
			Aggregation: metrics.Count(),
		},
	}
}
