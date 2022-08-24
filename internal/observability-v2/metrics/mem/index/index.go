package index

import (
	"context"
	"runtime"
	"time"

	"github.com/vdaas/vald/internal/observability-v2/metrics"
)

type memoryMetrics struct {
}

func New() metrics.Metric {
	return &memoryMetrics{}
}

func (mm *memoryMetrics) Register(m metrics.Meter) error {
	alloc, err := m.AsyncInt64().Gauge(
		"alloc_bytes",
		metrics.WithDescription("currently allocated number of bytes on the heap"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	totalAlloc, err := m.AsyncInt64().Gauge(
		"alloc_bytes_total",
		metrics.WithDescription("cumulative bytes allocated for heap objects"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	sys, err := m.AsyncInt64().Gauge(
		"sys_bytes",
		metrics.WithDescription("total bytes of memory obtained from the OS"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	mallocs, err := m.AsyncInt64().Gauge(
		"mallocs_total",
		metrics.WithDescription("the cumulative count of heap objects allocated"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	frees, err := m.AsyncInt64().Gauge(
		"frees_total",
		metrics.WithDescription("the cumulative count of heap objects freed"),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	heapAlloc, err := m.AsyncInt64().Gauge(
		"heap_alloc_bytes",
		metrics.WithDescription("bytes of allocated heap object"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapSys, err := m.AsyncInt64().Gauge(
		"heap_sys_bytes",
		metrics.WithDescription("bytes of heap memory obtained from the OS"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapIdle, err := m.AsyncInt64().Gauge(
		"heap_idle_bytes",
		metrics.WithDescription("bytes in idle (unused) spans"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapInuse, err := m.AsyncInt64().Gauge(
		"heap_inuse_bytes",
		metrics.WithDescription("bytes in in-use spans"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	heapReleased, err := m.AsyncInt64().Gauge(
		"heap_released_bytes",
		metrics.WithDescription("bytes of physical memory returned to the OS"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	stackInuse, err := m.AsyncInt64().Gauge(
		"stack_inuse_bytes",
		metrics.WithDescription("bytes in stack spans"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	stackSys, err := m.AsyncInt64().Gauge(
		"stack_sys_bytes",
		metrics.WithDescription("bytes of stack memory obtained from the OS"),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	pauseTotalMs, err := m.AsyncInt64().Gauge( // TODO
		"pause_ms_total",
		metrics.WithDescription("the cumulative milliseconds in GC"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return err
	}

	numGC, err := m.AsyncInt64().Gauge(
		"gc_count",
		metrics.WithDescription("the number of completed GC cycles"),
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
