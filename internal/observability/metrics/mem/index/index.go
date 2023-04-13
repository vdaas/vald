//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package index

import (
	"context"
	"runtime"
	"time"

	"github.com/vdaas/vald/internal/observability/metrics"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	allocMetricsName        = "alloc_bytes"
	allocMetricsDescription = "Currently allocated number of bytes on the heap"

	totalAllocMetricsName        = "alloc_bytes_total"
	totalAllocMetricsDescription = "Cumulative bytes allocated for heap objects"

	sysMetricsName        = "sys_bytes"
	sysMetricsDescription = "Total bytes of memory obtained from the OS"

	mallocsMetricsName        = "mallocs_total"
	mallocsMetricsDescription = "The cumulative count of heap objects allocated"

	freesMetricsName        = "frees_total"
	freesMetricsDescription = "The cumulative count of heap objects freed"

	heapAllocMetricsName        = "heap_alloc_bytes"
	heapAllocMetricsDescription = "Bytes of allocated heap object"

	heapSysMetricsName        = "heap_sys_bytes"
	heapSysMetricsDescription = "Bytes of heap memory obtained from the OS"

	heapIdleMetricsName        = "heap_idle_bytes"
	heapIdleMetricsDescription = "Bytes in idle (unused) spans"

	heapInuseMetricsName        = "heap_inuse_bytes"
	heapInuseMetricsDescription = "Bytes in in-use spans"

	heapReleasedMetricsName        = "heap_released_bytes"
	heapReleasedMetricsDescription = "Bytes of physical memory returned to the OS"

	stackInuseMetricsName        = "stack_inuse_bytes"
	stackInuseMetricsDescription = "Bytes in stack spans"

	stackSysMetricsName        = "stack_sys_bytes"
	stackSysMetricsDescription = "Bytes of stack memory obtained from the OS"

	pauseTotalMsMetricsName        = "pause_ms_total"
	pauseTotalMsMetricsDescription = "The cumulative milliseconds in GC"

	numGCMetricsName        = "gc_count"
	numGCMetricsDescription = "The number of completed GC cycles"
)

type memoryMetrics struct{}

func New() metrics.Metric {
	return &memoryMetrics{}
}

func (*memoryMetrics) View() ([]*metrics.View, error) {
	alloc, err := view.New(
		view.MatchInstrumentName(allocMetricsDescription),
		view.WithSetDescription(allocMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	totalAlloc, err := view.New(
		view.MatchInstrumentName(totalAllocMetricsDescription),
		view.WithSetDescription(totalAllocMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	sys, err := view.New(
		view.MatchInstrumentName(sysMetricsName),
		view.WithSetDescription(sysMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	mallocs, err := view.New(
		view.MatchInstrumentName(mallocsMetricsName),
		view.WithSetDescription(mallocsMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	frees, err := view.New(
		view.MatchInstrumentName(freesMetricsName),
		view.WithSetDescription(freesMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	heapAlloc, err := view.New(
		view.MatchInstrumentName(heapAllocMetricsName),
		view.WithSetDescription(heapAllocMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	heapSys, err := view.New(
		view.MatchInstrumentName(heapSysMetricsName),
		view.WithSetDescription(heapSysMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	heapIdle, err := view.New(
		view.MatchInstrumentName(heapIdleMetricsName),
		view.WithSetDescription(heapIdleMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	heapInuse, err := view.New(
		view.MatchInstrumentName(heapInuseMetricsName),
		view.WithSetDescription(heapInuseMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	heapReleased, err := view.New(
		view.MatchInstrumentName(heapReleasedMetricsName),
		view.WithSetDescription(heapReleasedMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	stackInuse, err := view.New(
		view.MatchInstrumentName(stackInuseMetricsName),
		view.WithSetDescription(stackInuseMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	stackSys, err := view.New(
		view.MatchInstrumentName(stackSysMetricsName),
		view.WithSetDescription(stackSysMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	pauseTotalMs, err := view.New(
		view.MatchInstrumentName(pauseTotalMsMetricsName),
		view.WithSetDescription(pauseTotalMsMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	numGC, err := view.New(
		view.MatchInstrumentName(numGCMetricsName),
		view.WithSetDescription(numGCMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	return []*metrics.View{
		&alloc,
		&totalAlloc,
		&sys,
		&mallocs,
		&frees,
		&heapAlloc,
		&heapSys,
		&heapIdle,
		&heapInuse,
		&heapReleased,
		&stackInuse,
		&stackSys,
		&pauseTotalMs,
		&numGC,
	}, nil
}

func (*memoryMetrics) Register(m metrics.Meter) error {
	alloc, err := m.AsyncInt64().Gauge(
		allocMetricsName,
		metrics.WithDescription(allocMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	totalAlloc, err := m.AsyncInt64().Gauge(
		totalAllocMetricsDescription,
		metrics.WithDescription(totalAllocMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	sys, err := m.AsyncInt64().Gauge(
		sysMetricsName,
		metrics.WithDescription(sysMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	mallocs, err := m.AsyncInt64().Gauge(
		mallocsMetricsName,
		metrics.WithDescription(mallocsMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	frees, err := m.AsyncInt64().Gauge(
		freesMetricsName,
		metrics.WithDescription(freesMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	heapAlloc, err := m.AsyncInt64().Gauge(
		heapAllocMetricsName,
		metrics.WithDescription(heapAllocMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapSys, err := m.AsyncInt64().Gauge(
		heapSysMetricsName,
		metrics.WithDescription(heapSysMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapIdle, err := m.AsyncInt64().Gauge(
		heapIdleMetricsName,
		metrics.WithDescription(heapIdleMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapInuse, err := m.AsyncInt64().Gauge(
		heapInuseMetricsName,
		metrics.WithDescription(heapInuseMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapReleased, err := m.AsyncInt64().Gauge(
		heapReleasedMetricsName,
		metrics.WithDescription(heapReleasedMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	stackInuse, err := m.AsyncInt64().Gauge(
		stackInuseMetricsName,
		metrics.WithDescription(stackInuseMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	stackSys, err := m.AsyncInt64().Gauge(
		stackSysMetricsName,
		metrics.WithDescription(stackSysMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	pauseTotalMs, err := m.AsyncInt64().Gauge(
		pauseTotalMsMetricsName,
		metrics.WithDescription(pauseTotalMsMetricsDescription),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	numGC, err := m.AsyncInt64().Gauge(
		numGCMetricsName,
		metrics.WithDescription(numGCMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			alloc,
			totalAlloc,
			sys,
			mallocs,
			frees,
			heapAlloc,
			heapSys,
			heapIdle,
			heapInuse,
			heapReleased,
			stackInuse,
			stackSys,
			pauseTotalMs,
			numGC,
		},
		func(ctx context.Context) {
			var mstats runtime.MemStats
			runtime.ReadMemStats(&mstats)

			alloc.Observe(ctx, int64(mstats.Alloc))
			totalAlloc.Observe(ctx, int64(mstats.TotalAlloc))
			sys.Observe(ctx, int64(mstats.Sys))
			mallocs.Observe(ctx, int64(mstats.Mallocs))
			frees.Observe(ctx, int64(mstats.Frees))
			heapAlloc.Observe(ctx, int64(mstats.HeapAlloc))
			heapSys.Observe(ctx, int64(mstats.HeapSys))
			heapIdle.Observe(ctx, int64(mstats.HeapIdle))
			heapInuse.Observe(ctx, int64(mstats.HeapInuse))
			heapReleased.Observe(ctx, int64(mstats.HeapReleased))
			stackInuse.Observe(ctx, int64(mstats.StackInuse))
			stackSys.Observe(ctx, int64(mstats.StackSys))

			ptMs := int64(0)
			if mstats.PauseTotalNs > 0 {
				ptMs = int64(mstats.PauseTotalNs / uint64(time.Millisecond))
			}
			pauseTotalMs.Observe(ctx, ptMs)
			numGC.Observe(ctx, int64(mstats.NextGC))
		},
	)
}
