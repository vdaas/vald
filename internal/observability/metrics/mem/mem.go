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

// Package mem provides memory metrics functions
package mem

import (
	"context"
	"runtime"
	"time"

	"github.com/vdaas/vald/internal/observability/metrics"
)

type memory struct {
	alloc        metrics.Int64Measure
	totalAlloc   metrics.Int64Measure
	sys          metrics.Int64Measure
	mallocs      metrics.Int64Measure
	frees        metrics.Int64Measure
	heapAlloc    metrics.Int64Measure
	heapSys      metrics.Int64Measure
	heapIdle     metrics.Int64Measure
	heapInuse    metrics.Int64Measure
	heapReleased metrics.Int64Measure
	stackInuse   metrics.Int64Measure
	stackSys     metrics.Int64Measure
	pauseTotalMs metrics.Int64Measure
	numGC        metrics.Int64Measure
}

func New() metrics.Metric {
	return &memory{
		alloc:        *metrics.Int64(metrics.ValdOrg+"/memory/alloc", "currently allocated number of bytes on the heap", metrics.UnitBytes),
		totalAlloc:   *metrics.Int64(metrics.ValdOrg+"/memory/alloc_total", "cumulative bytes allocated for heap objects", metrics.UnitBytes),
		sys:          *metrics.Int64(metrics.ValdOrg+"/memory/sys", "total bytes of memory obtained from the OS", metrics.UnitBytes),
		mallocs:      *metrics.Int64(metrics.ValdOrg+"/memory/mallocs_total", "the cumulative count of heap objects allocated", metrics.UnitDimensionless),
		frees:        *metrics.Int64(metrics.ValdOrg+"/memory/frees_total", "the cumulative count of heap objects freed", metrics.UnitDimensionless),
		heapAlloc:    *metrics.Int64(metrics.ValdOrg+"/memory/heap_alloc", "bytes of allocated heap objects", metrics.UnitBytes),
		heapSys:      *metrics.Int64(metrics.ValdOrg+"/memory/heap_sys", "bytes of heap memory obtained from the OS", metrics.UnitBytes),
		heapIdle:     *metrics.Int64(metrics.ValdOrg+"/memory/heap_idle", "bytes in idle (unused) spans", metrics.UnitBytes),
		heapInuse:    *metrics.Int64(metrics.ValdOrg+"/memory/heap_inuse", "bytes in in-use spans", metrics.UnitBytes),
		heapReleased: *metrics.Int64(metrics.ValdOrg+"/memory/heap_released", "bytes of physical memory returned to the OS", metrics.UnitBytes),
		stackInuse:   *metrics.Int64(metrics.ValdOrg+"/memory/stack_inuse", "bytes in stack spans", metrics.UnitBytes),
		stackSys:     *metrics.Int64(metrics.ValdOrg+"/memory/stack_sys", "bytes of stack memory obtained from the OS", metrics.UnitBytes),
		pauseTotalMs: *metrics.Int64(metrics.ValdOrg+"/memory/pause_ms_total", "the cumulative milliseconds in GC", metrics.UnitMilliseconds),
		numGC:        *metrics.Int64(metrics.ValdOrg+"/memory/gc_count", "the number of completed GC cycles", metrics.UnitDimensionless),
	}
}

func (m *memory) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	var mstats runtime.MemStats
	runtime.ReadMemStats(&mstats)

	pauseTotalMs := int64(0)
	if mstats.PauseTotalNs > 0 {
		pauseTotalMs = int64(mstats.PauseTotalNs / uint64(time.Millisecond))
	}

	return []metrics.Measurement{
		m.alloc.M(int64(mstats.Alloc)),
		m.totalAlloc.M(int64(mstats.TotalAlloc)),
		m.sys.M(int64(mstats.Sys)),
		m.mallocs.M(int64(mstats.Mallocs)),
		m.frees.M(int64(mstats.Frees)),
		m.heapAlloc.M(int64(mstats.HeapAlloc)),
		m.heapSys.M(int64(mstats.HeapSys)),
		m.heapIdle.M(int64(mstats.HeapIdle)),
		m.heapInuse.M(int64(mstats.HeapInuse)),
		m.heapReleased.M(int64(mstats.HeapReleased)),
		m.stackInuse.M(int64(mstats.StackInuse)),
		m.stackSys.M(int64(mstats.StackSys)),
		m.pauseTotalMs.M(pauseTotalMs),
		m.numGC.M(int64(mstats.NumGC)),
	}, nil
}

func (m *memory) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (m *memory) View() []*metrics.View {
	return []*metrics.View{
		{
			Name:        "alloc_bytes",
			Description: m.alloc.Description(),
			Measure:     &m.alloc,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "alloc_bytes_total",
			Description: m.totalAlloc.Description(),
			Measure:     &m.totalAlloc,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "sys_bytes",
			Description: m.sys.Description(),
			Measure:     &m.sys,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "mallocs_total",
			Description: m.mallocs.Description(),
			Measure:     &m.mallocs,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "frees_total",
			Description: m.frees.Description(),
			Measure:     &m.frees,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "heap_alloc_bytes",
			Description: m.heapAlloc.Description(),
			Measure:     &m.heapAlloc,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "heap_sys_bytes",
			Description: m.heapSys.Description(),
			Measure:     &m.heapSys,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "heap_idle_bytes",
			Description: m.heapIdle.Description(),
			Measure:     &m.heapIdle,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "heap_inuse_bytes",
			Description: m.heapInuse.Description(),
			Measure:     &m.heapInuse,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "heap_released_bytes",
			Description: m.heapReleased.Description(),
			Measure:     &m.heapReleased,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "stack_inuse_bytes",
			Description: m.stackInuse.Description(),
			Measure:     &m.stackInuse,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "stack_sys_bytes",
			Description: m.stackSys.Description(),
			Measure:     &m.stackSys,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "pause_ms_total",
			Description: m.pauseTotalMs.Description(),
			Measure:     &m.pauseTotalMs,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "gc_count",
			Description: m.numGC.Description(),
			Measure:     &m.numGC,
			Aggregation: metrics.LastValue(),
		},
	}
}
