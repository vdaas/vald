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

// Package cpu provides cpu metrics functions
package cpu

import (
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type cpu struct {
	numCPU metrics.Int64Measure
}

func NewMetric() metrics.Metric {
	return &cpu{
		numCPU: *metrics.Int64("vdaas.org/cpu/num", "number of cpu", metrics.UnitDimensionless),
	}
}

func (c *cpu) Measurement() ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		c.numCPU.M(int64(runtime.NumCPU())),
	}, nil
}

func (c *cpu) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "num_cpu",
			Description: "number of cpu",
			Measure:     &c.numCPU,
			Aggregation: metrics.LastValue(),
		},
	}
}
