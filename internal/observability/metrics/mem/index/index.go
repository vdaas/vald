// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package index

import (
	"context"
	"runtime"
	"time"

	"github.com/vdaas/vald/internal/observability/metrics"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
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

func (*memoryMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        allocMetricsName,
				Description: allocMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        totalAllocMetricsDescription,
				Description: totalAllocMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        sysMetricsName,
				Description: sysMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        mallocsMetricsName,
				Description: mallocsMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        freesMetricsName,
				Description: freesMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        heapAllocMetricsName,
				Description: heapAllocMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        heapSysMetricsName,
				Description: heapSysMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        heapIdleMetricsName,
				Description: heapIdleMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        heapInuseMetricsName,
				Description: heapInuseMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        heapReleasedMetricsName,
				Description: heapReleasedMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        stackInuseMetricsName,
				Description: stackInuseMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        stackSysMetricsName,
				Description: stackSysMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        pauseTotalMsMetricsName,
				Description: pauseTotalMsMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        numGCMetricsName,
				Description: numGCMetricsDescription,
			},
			view.Stream{
				Aggregation: aggregation.LastValue{},
			},
		),
	}, nil
}

func (*memoryMetrics) Register(m metrics.Meter) error {
	alloc, err := m.Int64ObservableGauge(
		allocMetricsName,
		metrics.WithDescription(allocMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	totalAlloc, err := m.Int64ObservableGauge(
		totalAllocMetricsDescription,
		metrics.WithDescription(totalAllocMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	sys, err := m.Int64ObservableGauge(
		sysMetricsName,
		metrics.WithDescription(sysMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	mallocs, err := m.Int64ObservableGauge(
		mallocsMetricsName,
		metrics.WithDescription(mallocsMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	frees, err := m.Int64ObservableGauge(
		freesMetricsName,
		metrics.WithDescription(freesMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	heapAlloc, err := m.Int64ObservableGauge(
		heapAllocMetricsName,
		metrics.WithDescription(heapAllocMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapSys, err := m.Int64ObservableGauge(
		heapSysMetricsName,
		metrics.WithDescription(heapSysMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapIdle, err := m.Int64ObservableGauge(
		heapIdleMetricsName,
		metrics.WithDescription(heapIdleMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapInuse, err := m.Int64ObservableGauge(
		heapInuseMetricsName,
		metrics.WithDescription(heapInuseMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapReleased, err := m.Int64ObservableGauge(
		heapReleasedMetricsName,
		metrics.WithDescription(heapReleasedMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	stackInuse, err := m.Int64ObservableGauge(
		stackInuseMetricsName,
		metrics.WithDescription(stackInuseMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	stackSys, err := m.Int64ObservableGauge(
		stackSysMetricsName,
		metrics.WithDescription(stackSysMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	pauseTotalMs, err := m.Int64ObservableGauge(
		pauseTotalMsMetricsName,
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
