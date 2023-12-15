// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mem

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	// metrics from runtime.Memstats
	allocMetricsName        = "alloc_bytes"
	allocMetricsDescription = "Currently allocated number of bytes on the heap"

	totalAllocMetricsName        = "alloc_bytes_total"
	totalAllocMetricsDescription = "Cumulative bytes allocated for heap objects"

	sysMetricsName        = "sys_bytes"
	sysMetricsDescription = "Total bytes of memory obtained from the OS"

	lookupsMetricsName        = "lookups_count"
	lookupsMetricsDescription = "The number of pointers"

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

	heapObjectsMetricsName        = "heap_objects_count"
	heapObjectsMetricsDescription = "The number of allocated heap objects"

	stackInuseMetricsName        = "stack_inuse_bytes"
	stackInuseMetricsDescription = "Bytes in stack spans"

	stackSysMetricsName        = "stack_sys_bytes"
	stackSysMetricsDescription = "Bytes of stack memory obtained from the OS"

	mspanInuseMetricsName        = "mspan_inuse_bytes"
	mspanInuseMetricsDescription = "Bytes of allocated mspan structures"

	mspanSysMetricsName        = "mspan_sys_bytes"
	mspanSysMetricsDescription = "Bytes of memory obtained from the OS for mspan structures"

	mcacheInuseMetricsName        = "mcache_inuse_bytes"
	mcacheInuseMetricsDescription = "Bytes of allocated mcache structures"

	mcacheSysMetricsName        = "mcache_sys_bytes"
	mcacheSysMetricsDescription = "Bytes of memory obtained from the OS mcache structures"

	buckHashSysMetricsName        = "buckhash_sys_bytes"
	buckHashSysMetricsDescription = "Bytes of memory in profiling bucket hash tables"

	gcSysMetricsName        = "gc_sys_bytes"
	gcSysMetricsDescription = "Bytes of memory in GC metadata"

	otherSysMetricsName        = "other_sys_bytes"
	otherSysMetricsDescription = "Bytes of memory in misc off-heap runtime allocations"

	nextGcSysMetricsName        = "next_gc_bytes"
	nextGcSysMetricsDescription = "Target heap size of the next GC"

	pauseTotalMsMetricsName        = "pause_ms_total"
	pauseTotalMsMetricsDescription = "The cumulative milliseconds in GC"

	numGCMetricsName        = "gc_count"
	numGCMetricsDescription = "The number of completed GC cycles"

	numForcedGCMetricsName        = "forced_gc_count"
	numForcedGCMetricsDescription = "The number of GC cycles called by the application"

	heapWillReturnMetricsName        = "heap_will_return_bytes"
	heapWillReturnMetricsDescription = "Bytes of returning to OS. It contains the two following parts (heapWillReturn = heapIdle - heapReleased)"

	liveObjectsMetricsName        = "live_objects_count"
	liveObjectsMetricsDescription = "The cumulative count of living heap objects allocated. It contains the two following parts (liveObjects = mallocs - frees)"

	// metrics from /proc/<pid>/status
	vmpeakMetricsName        = "vmpeak_bytes"
	vmpeakMetricsDescription = "peak virtual memory size"

	vmsizeMetricsName        = "vmsize_bytes"
	vmsizeMetricsDescription = "toal program size"

	vmdataMetricsName        = "vmdata_bytes"
	vmdataMetricsDescription = "size of private data segments"

	vmrssMetricsName        = "vmrss_bytes"
	vmrssMetricsDescription = "size of memory portions. It contains the three following parts (VmRSS = RssAnon + RssFile + RssShmem)"

	vmhwmMetricsName        = "vmhwm_bytes"
	vmhwmMetricsDescription = "peak resident set size (\"high water mark\")"

	vmstkMetricsName        = "vmstk_bytes"
	vmstkMetricsDescription = "size of stack segments"

	vmswapMetricsName        = "vmswap_bytes"
	vmswapMetricsDescription = "amount of swap used by anonymous private data (shmem swap usage is not included)"

	vmexeMetricsName        = "vmexe_bytes"
	vmexeMetricsDescription = "size of text segment"

	vmlibMetricsName        = "vmlib_bytes"
	vmlibMetricsDescription = "size of shared library code"

	vmlckMetricsName        = "vmlck_bytes"
	vmlckMetricsDescription = "locked memory size"

	vmpinMetricsName        = "vmpin_bytes"
	vmpinMetricsDescription = "pinned memory size"

	vmpteMetricsName        = "vmpte_bytes"
	vmpteMetricsDescription = "size of page table entries"

	k = 1024
)

type metricsInfo struct {
	Name  string
	Desc  string
	Unit  metrics.Unit
	Value func() int64
}

func getMemstatsMetrics() []*metricsInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return []*metricsInfo{
		{
			Name: allocMetricsName,
			Desc: allocMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.Alloc)
			},
		},
		{
			Name: totalAllocMetricsName,
			Desc: totalAllocMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.TotalAlloc)
			},
		},
		{
			Name: sysMetricsName,
			Desc: sysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.Sys)
			},
		},
		{
			Name: lookupsMetricsName,
			Desc: lookupsMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.Lookups)
			},
		},
		{
			Name: mallocsMetricsName,
			Desc: mallocsMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.Mallocs)
			},
		},
		{
			Name: freesMetricsName,
			Desc: freesMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.Frees)
			},
		},
		{
			Name: heapAllocMetricsName,
			Desc: heapAllocMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapAlloc)
			},
		},
		{
			Name: heapSysMetricsName,
			Desc: heapSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapSys)
			},
		},
		{
			Name: heapIdleMetricsName,
			Desc: heapIdleMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapIdle)
			},
		},
		{
			Name: heapInuseMetricsName,
			Desc: heapInuseMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapInuse)
			},
		},
		{
			Name: heapReleasedMetricsName,
			Desc: heapReleasedMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapReleased)
			},
		},
		{
			Name: heapObjectsMetricsName,
			Desc: heapObjectsMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.HeapObjects)
			},
		},
		{
			Name: stackInuseMetricsName,
			Desc: stackInuseMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.StackInuse)
			},
		},
		{
			Name: stackSysMetricsName,
			Desc: stackSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.StackSys)
			},
		},
		{
			Name: mspanInuseMetricsName,
			Desc: mspanInuseMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.MSpanInuse)
			},
		},
		{
			Name: mspanSysMetricsName,
			Desc: mspanSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.MSpanSys)
			},
		},
		{
			Name: mcacheInuseMetricsName,
			Desc: mcacheInuseMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.MCacheInuse)
			},
		},
		{
			Name: mcacheSysMetricsName,
			Desc: mcacheSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.MCacheSys)
			},
		},
		{
			Name: buckHashSysMetricsName,
			Desc: buckHashSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.BuckHashSys)
			},
		},
		{
			Name: gcSysMetricsName,
			Desc: gcSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.GCSys)
			},
		},
		{
			Name: otherSysMetricsName,
			Desc: otherSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.OtherSys)
			},
		},
		{
			Name: nextGcSysMetricsName,
			Desc: nextGcSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.NextGC)
			},
		},
		{
			Name: pauseTotalMsMetricsName,
			Desc: pauseTotalMsMetricsDescription,
			Unit: metrics.Milliseconds,
			Value: func() int64 {
				ptMs := int64(0)
				if m.PauseTotalNs > 0 {
					ptMs = int64(m.PauseTotalNs / uint64(time.Millisecond))
				}
				return ptMs
			},
		},
		{
			Name: numGCMetricsName,
			Desc: numGCMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.NumGC)
			},
		},
		{
			Name: numForcedGCMetricsName,
			Desc: numForcedGCMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.NumForcedGC)
			},
		},
		{
			Name: heapWillReturnMetricsName,
			Desc: heapWillReturnMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapIdle - m.HeapReleased)
			},
		},
		{
			Name: liveObjectsMetricsName,
			Desc: liveObjectsMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.Mallocs - m.Frees)
			},
		},
	}
}

func getProcStatusMetrics(pid int) ([]*metricsInfo, error) {
	buf, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return nil, err
	}
	lines := strings.Split(conv.Btoa(buf), "\n")
	m := make([]*metricsInfo, 0)
	for _, line := range lines {
		fields := strings.Fields(line)
		switch {
		case strings.HasPrefix(line, "VmPeak"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmpeakMetricsName,
					Desc: vmpeakMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmSize"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmsizeMetricsName,
					Desc: vmsizeMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmHWM"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmhwmMetricsName,
					Desc: vmhwmMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmRSS"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmrssMetricsName,
					Desc: vmrssMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmData"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmdataMetricsName,
					Desc: vmdataMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmStk"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmstkMetricsName,
					Desc: vmstkMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmExe"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmexeMetricsName,
					Desc: vmexeMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmLck"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmlckMetricsName,
					Desc: vmlckMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmLib"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmlibMetricsName,
					Desc: vmlibMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmPTE"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmpteMetricsName,
					Desc: vmpteMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmSwap"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmswapMetricsName,
					Desc: vmswapMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmPin"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: vmpinMetricsName,
					Desc: vmpinMetricsDescription,
					Unit: unit.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		}
	}
	return m, nil
}

type memMetrics struct {
	pid int
}

func New() metrics.Metric {
	return &memMetrics{
		pid: os.Getpid(),
	}
}

func (mm *memMetrics) View() ([]*metrics.View, error) {
	mInfo := getMemstatsMetrics()
	if m, err := getProcStatusMetrics(mm.pid); err == nil {
		mInfo = append(mInfo, m...)
	}

	views := make([]*metrics.View, 0, len(mInfo))
	for _, m := range mInfo {
		v, err := view.New(
			view.MatchInstrumentName(m.Name),
			view.WithSetDescription(m.Desc),
			view.WithSetAggregation(aggregation.LastValue{}),
		)
		if err != nil {
			return nil, err
		}
		views = append(views, &v)
	}
	return views, nil
}

func (mm *memMetrics) Register(m metrics.Meter) error {
	mInfo := getMemstatsMetrics()
	if metrics, err := getProcStatusMetrics(mm.pid); err == nil {
		mInfo = append(mInfo, metrics...)
	}

	instruments := make([]metrics.AsynchronousInstrument, 0, len(mInfo))
	for _, info := range mInfo {
		instrument, err := m.AsyncInt64().Gauge(
			info.Name,
			metrics.WithDescription(info.Desc),
			metrics.WithUnit(info.Unit),
		)
		if err != nil {
			return err
		}
		instruments = append(instruments, instrument)
	}

	return m.RegisterCallback(
		instruments,
		func(ctx context.Context) {
			var mstats runtime.MemStats
			runtime.ReadMemStats(&mstats)

			for i, instrument := range instruments {
				instrument.(asyncint64.Gauge).Observe(ctx, mInfo[i].Value())
			}
		},
	)
}
