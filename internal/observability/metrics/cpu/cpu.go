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
	"os"
	"runtime"

	"github.com/shirou/gopsutil/process"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type cpu struct {
	process    *process.Process
	numCPU     metrics.Int64Measure
	percentCPU metrics.Float64Measure
	numThreads metrics.Int64Measure
}

func NewMetric() (metrics.Metric, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, err
	}

	return &cpu{
		process:    p,
		numCPU:     *metrics.Int64("vdaas.org/cpu/num", "number of cpu", metrics.UnitDimensionless),
		percentCPU: *metrics.Float64("vdaas.org/cpu/percent", "cpu usage", metrics.UnitDimensionless),
		numThreads: *metrics.Int64("vdaas.org/thread/num", "number of threads", metrics.UnitDimensionless),
	}, nil
}

func (c *cpu) Measurement() ([]metrics.Measurement, error) {
	cpuPercent, err := c.process.CPUPercent()
	if err != nil {
		return nil, err
	}
	numThreads, err := c.process.NumThreads()
	if err != nil {
		return nil, err
	}

	return []metrics.Measurement{
		c.numCPU.M(int64(runtime.NumCPU())),
		c.percentCPU.M(cpuPercent),
		c.numThreads.M(int64(numThreads)),
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
		&metrics.View{
			Name:        "cpu_percent",
			Description: "cpu usage",
			Measure:     &c.percentCPU,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "num_threads",
			Description: "number of threads",
			Measure:     &c.numThreads,
			Aggregation: metrics.LastValue(),
		},
	}
}
