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
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type memory struct {
	alloc        metrics.Int64Measure
	totalAlloc   metrics.Int64Measure
	sys          metrics.Int64Measure
	mallocs      metrics.Int64Measure
	frees        metrics.Int64Measure
	pauseTotalMs metrics.Float64Measure
	numGC        metrics.Int64Measure
}

func NewMetric() metrics.Metric {
	return &memory{
		alloc:        *metrics.Int64("vdaas.org/vald/memory/alloc", "currently allocated number of bytes on the heap", metrics.UnitBytes),
		totalAlloc:   *metrics.Int64("vdaas.org/vald/memory/alloc_total", "cumulative bytes allocated for heap objects", metrics.UnitBytes),
		sys:          *metrics.Int64("vdaas.org/vald/memory/sys", "total bytes of memory obtained from the OS", metrics.UnitBytes),
		mallocs:      *metrics.Int64("vdaas.org/vald/memory/mallocs_total", "the cumulative count of heap objects allocated", metrics.UnitDimensionless),
		frees:        *metrics.Int64("vdaas.org/vald/memory/frees_total", "the cumulative count of heap objects freed", metrics.UnitDimensionless),
		pauseTotalMs: *metrics.Float64("vdaas.org/vald/memory/pause_ms_total", "the cumulative milliseconds in GC", metrics.UnitMilliseconds),
		numGC:        *metrics.Int64("vdaas.org/vald/memory/gc_count", "the number of completed GC cycles", metrics.UnitDimensionless),
	}
}

func (m *memory) Measurement() ([]metrics.Measurement, error) {
	var mstats runtime.MemStats
	runtime.ReadMemStats(&mstats)

	pauseTotalMs := 0.0
	if mstats.PauseTotalNs > 0 {
		pauseTotalMs = float64(mstats.PauseTotalNs) / 1000000.0
	}

	return []metrics.Measurement{
		m.alloc.M(int64(mstats.Alloc)),
		m.totalAlloc.M(int64(mstats.TotalAlloc)),
		m.sys.M(int64(mstats.Sys)),
		m.mallocs.M(int64(mstats.Mallocs)),
		m.frees.M(int64(mstats.Frees)),
		m.pauseTotalMs.M(float64(pauseTotalMs)),
		m.numGC.M(int64(mstats.NumGC)),
	}, nil
}

func (m *memory) MeasurementWithTags() ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (m *memory) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "alloc_bytes",
			Description: "currently allocated number of bytes on the heap",
			Measure:     &m.alloc,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "alloc_bytes_total",
			Description: "cumulative bytes allocated for heap objects",
			Measure:     &m.totalAlloc,
			Aggregation: metrics.Count(),
		},
		&metrics.View{
			Name:        "sys_bytes",
			Description: "total bytes of memory obtained from the OS",
			Measure:     &m.sys,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "mallocs_total",
			Description: "the cumulative count of heap objects allocated",
			Measure:     &m.mallocs,
			Aggregation: metrics.Count(),
		},
		&metrics.View{
			Name:        "frees_total",
			Description: "the cumulative count of heap objects freed",
			Measure:     &m.frees,
			Aggregation: metrics.Count(),
		},
		&metrics.View{
			Name:        "pause_ms_total",
			Description: "the cumulative milliseconds in GC",
			Measure:     &m.pauseTotalMs,
			Aggregation: metrics.Count(),
		},
		&metrics.View{
			Name:        "gc_count",
			Description: "the number of completed GC cycles",
			Measure:     &m.numGC,
			Aggregation: metrics.LastValue(),
		},
	}
}
