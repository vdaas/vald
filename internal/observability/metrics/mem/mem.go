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

// Package mem provides memory metrics functions
package mem

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type memory struct {
	total     metrics.Int64Measure
	available metrics.Int64Measure
	used      metrics.Int64Measure
	free      metrics.Int64Measure
}

func NewMetric() metrics.Metric {
	return &memory{
		total:     *metrics.Int64("vdaas.org/memory/total", "size of total memory", metrics.UnitBytes),
		available: *metrics.Int64("vdaas.org/memory/available", "size of available memory", metrics.UnitBytes),
		used:      *metrics.Int64("vdaas.org/memory/used", "size of used memory", metrics.UnitBytes),
		free:      *metrics.Int64("vdaas.org/memory/free", "size of free memory", metrics.UnitBytes),
	}
}

func (m *memory) Measurement() ([]metrics.Measurement, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return []metrics.Measurement{
		m.total.M(int64(v.Total)),
		m.available.M(int64(v.Available)),
		m.used.M(int64(v.Used)),
		m.free.M(int64(v.Free)),
	}, nil
}

func (m *memory) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "memory_total",
			Description: "size of total memory",
			Measure:     &m.total,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "memory_available",
			Description: "size of available memory",
			Measure:     &m.available,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "memory_used",
			Description: "size of used memory",
			Measure:     &m.used,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "memory_free",
			Description: "size of free memory",
			Measure:     &m.free,
			Aggregation: metrics.LastValue(),
		},
	}
}
