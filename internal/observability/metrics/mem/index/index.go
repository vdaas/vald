//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	AllocMetricsName        = "alloc_bytes"
	AllocMetricsDescription = "Currently allocated number of bytes on the heap"

	TotalAllocMetricsName        = "alloc_bytes_total"
	TotalAllocMetricsDescription = "Cumulative bytes allocated for heap objects"

	SysMetricsName        = "sys_bytes"
	SysMetricsDescription = "Total bytes of memory obtained from the OS"

	MallocsMetricsName        = "mallocs_total"
	MallocsMetricsDescription = "The cumulative count of heap objects allocated"

	FreesMetricsName        = "frees_total"
	FreesMetricsDescription = "The cumulative count of heap objects freed"

	HeapAllocMetricsName        = "heap_alloc_bytes"
	HeapAllocMetricsDescription = "Bytes of allocated heap object"

	HeapSysMetricsName        = "heap_sys_bytes"
	HeapSysMetricsDescription = "Bytes of heap memory obtained from the OS"

	HeapIdleMetricsName        = "heap_idle_bytes"
	HeapIdleMetricsDescription = "Bytes in idle (unused) spans"

	HeapInuseMetricsName        = "heap_inuse_bytes"
	HeapInuseMetricsDescription = "Bytes in in-use spans"

	HeapReleasedMetricsName        = "heap_released_bytes"
	HeapReleasedMetricsDescription = "Bytes of physical memory returned to the OS"

	StackInuseMetricsName        = "stack_inuse_bytes"
	StackInuseMetricsDescription = "Bytes in stack spans"

	StackSysMetricsName        = "stack_sys_bytes"
	StackSysMetricsDescription = "Bytes of stack memory obtained from the OS"

	PauseTotalMsMetricsName        = "pause_ms_total"
	pauseTotalMsMetricsDescription = "The cumulative milliseconds in GC"

	numGCMetricsName        = "gc_count"
	numGCMetricsDescription = "The number of completed GC cycles"
)

type memoryMetrics struct{}

func New() metrics.Metric {
	return &memoryMetrics{}
}

func (*memoryMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        AllocMetricsName,
				Description: AllocMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        TotalAllocMetricsName,
				Description: TotalAllocMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        SysMetricsName,
				Description: SysMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        MallocsMetricsName,
				Description: MallocsMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        FreesMetricsName,
				Description: FreesMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        HeapAllocMetricsName,
				Description: HeapAllocMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        HeapSysMetricsName,
				Description: HeapSysMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        HeapIdleMetricsName,
				Description: HeapIdleMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        HeapInuseMetricsName,
				Description: HeapInuseMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        HeapReleasedMetricsName,
				Description: HeapReleasedMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        StackInuseMetricsName,
				Description: StackInuseMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        StackSysMetricsName,
				Description: StackSysMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        PauseTotalMsMetricsName,
				Description: pauseTotalMsMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        numGCMetricsName,
				Description: numGCMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

// skipcq: GO-R1005
func (*memoryMetrics) Register(m metrics.Meter) error {
	alloc, err := m.Int64ObservableGauge(
		AllocMetricsName,
		metrics.WithDescription(AllocMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	totalAlloc, err := m.Int64ObservableGauge(
		TotalAllocMetricsName,
		metrics.WithDescription(TotalAllocMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	sys, err := m.Int64ObservableGauge(
		SysMetricsName,
		metrics.WithDescription(SysMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	mallocs, err := m.Int64ObservableGauge(
		MallocsMetricsName,
		metrics.WithDescription(MallocsMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	frees, err := m.Int64ObservableGauge(
		FreesMetricsName,
		metrics.WithDescription(FreesMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	heapAlloc, err := m.Int64ObservableGauge(
		HeapAllocMetricsName,
		metrics.WithDescription(HeapAllocMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapSys, err := m.Int64ObservableGauge(
		HeapSysMetricsName,
		metrics.WithDescription(HeapSysMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapIdle, err := m.Int64ObservableGauge(
		HeapIdleMetricsName,
		metrics.WithDescription(HeapIdleMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapInuse, err := m.Int64ObservableGauge(
		HeapInuseMetricsName,
		metrics.WithDescription(HeapInuseMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapReleased, err := m.Int64ObservableGauge(
		HeapReleasedMetricsName,
		metrics.WithDescription(HeapReleasedMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	stackInuse, err := m.Int64ObservableGauge(
		StackInuseMetricsName,
		metrics.WithDescription(StackInuseMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	stackSys, err := m.Int64ObservableGauge(
		StackSysMetricsName,
		metrics.WithDescription(StackSysMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	pauseTotalMs, err := m.Int64ObservableGauge(
		PauseTotalMsMetricsName,
		metrics.WithDescription(pauseTotalMsMetricsDescription),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	numGC, err := m.Int64ObservableGauge(
		numGCMetricsName,
		metrics.WithDescription(numGCMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			var mstats runtime.MemStats
			runtime.ReadMemStats(&mstats)
			o.ObserveInt64(alloc, int64(mstats.Alloc))
			o.ObserveInt64(frees, int64(mstats.Frees))
			o.ObserveInt64(heapAlloc, int64(mstats.HeapAlloc))
			o.ObserveInt64(heapIdle, int64(mstats.HeapIdle))
			o.ObserveInt64(heapInuse, int64(mstats.HeapInuse))
			o.ObserveInt64(heapReleased, int64(mstats.HeapReleased))
			o.ObserveInt64(heapSys, int64(mstats.HeapSys))
			o.ObserveInt64(mallocs, int64(mstats.Mallocs))
			o.ObserveInt64(stackInuse, int64(mstats.StackInuse))
			o.ObserveInt64(stackSys, int64(mstats.StackSys))
			o.ObserveInt64(sys, int64(mstats.Sys))
			o.ObserveInt64(totalAlloc, int64(mstats.TotalAlloc))
			var ptMs int64
			if mstats.PauseTotalNs > 0 {
				ptMs = int64(float64(mstats.PauseTotalNs) / float64(time.Millisecond))
			}
			o.ObserveInt64(pauseTotalMs, ptMs)
			o.ObserveInt64(numGC, int64(mstats.NextGC))
			return nil
		},
		alloc,
		frees,
		heapAlloc,
		heapIdle,
		heapInuse,
		heapReleased,
		heapSys,
		mallocs,
		stackInuse,
		stackSys,
		sys,
		totalAlloc,
		pauseTotalMs,
		numGC,
	)
	return err
}
