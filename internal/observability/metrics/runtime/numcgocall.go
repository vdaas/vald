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
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type numcgocall struct {
	count metrics.Int64Measure
}

func NewNumberOfCGOCall() metrics.Metric {
	return &numcgocall{
		count: *metrics.Int64("vdaas.org/runtime/num_cgo_call", "number of cgo call", metrics.UnitDimensionless),
	}
}

func (n *numcgocall) Measurement() ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		n.count.M(int64(runtime.NumCgoCall())),
	}, nil
}

func (n *numcgocall) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "num_cgo_call",
			Description: "number of cgo call",
			Measure:     &n.count,
			Aggregation: metrics.LastValue(),
		},
	}
}
