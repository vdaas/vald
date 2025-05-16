// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"time"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/strings"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	// metrics from runtime.Memstats.
	AllocMetricsName        = "alloc_bytes"
	AllocMetricsDescription = "Currently allocated number of bytes on the heap"

	TotalAllocMetricsName        = "alloc_bytes_total"
	TotalAllocMetricsDescription = "Cumulative bytes allocated for heap objects"

	SysMetricsName        = "sys_bytes"
	SysMetricsDescription = "Total bytes of memory obtained from the OS"

	LookupsMetricsName        = "lookups_count"
	LookupsMetricsDescription = "The number of pointers"

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

	HeapObjectsMetricsName        = "heap_objects_count"
	HeapObjectsMetricsDescription = "The number of allocated heap objects"

	StackInuseMetricsName        = "stack_inuse_bytes"
	StackInuseMetricsDescription = "Bytes in stack spans"

	StackSysMetricsName        = "stack_sys_bytes"
	StackSysMetricsDescription = "Bytes of stack memory obtained from the OS"

	MspanInuseMetricsName        = "mspan_inuse_bytes"
	MspanInuseMetricsDescription = "Bytes of allocated mspan structures"

	MspanSysMetricsName        = "mspan_sys_bytes"
	MspanSysMetricsDescription = "Bytes of memory obtained from the OS for mspan structures"

	McacheInuseMetricsName        = "mcache_inuse_bytes"
	McacheInuseMetricsDescription = "Bytes of allocated mcache structures"

	McacheSysMetricsName        = "mcache_sys_bytes"
	McacheSysMetricsDescription = "Bytes of memory obtained from the OS mcache structures"

	BuckHashSysMetricsName        = "buckhash_sys_bytes"
	BuckHashSysMetricsDescription = "Bytes of memory in profiling bucket hash tables"

	GcSysMetricsName        = "gc_sys_bytes"
	GcSysMetricsDescription = "Bytes of memory in GC metadata"

	OtherSysMetricsName        = "other_sys_bytes"
	OtherSysMetricsDescription = "Bytes of memory in misc off-heap runtime allocations"

	NextGcSysMetricsName        = "next_gc_bytes"
	NextGcSysMetricsDescription = "Target heap size of the next GC"

	PauseTotalMsMetricsName        = "pause_ms_total"
	PauseTotalMsMetricsDescription = "The cumulative milliseconds in GC"

	NumGCMetricsName        = "gc_count"
	NumGCMetricsDescription = "The number of completed GC cycles"

	NumForcedGCMetricsName        = "forced_gc_count"
	NumForcedGCMetricsDescription = "The number of GC cycles called by the application"

	HeapWillReturnMetricsName        = "heap_will_return_bytes"
	HeapWillReturnMetricsDescription = "Bytes of returning to OS. It contains the two following parts (heapWillReturn = heapIdle - heapReleased)"

	LiveObjectsMetricsName        = "live_objects_count"
	LiveObjectsMetricsDescription = "The cumulative count of living heap objects allocated. It contains the two following parts (liveObjects = mallocs - frees)"

	// metrics from /proc/<pid>/status.
	VmpeakMetricsName        = "vmpeak_bytes"
	VmpeakMetricsDescription = "peak virtual memory size"

	VmsizeMetricsName        = "vmsize_bytes"
	VmsizeMetricsDescription = "toal program size"

	VmdataMetricsName        = "vmdata_bytes"
	VmdataMetricsDescription = "size of private data segments"

	VmrssMetricsName        = "vmrss_bytes"
	VmrssMetricsDescription = "size of memory portions. It contains the three following parts (VmRSS = RssAnon + RssFile + RssShmem)"

	VmhwmMetricsName        = "vmhwm_bytes"
	VmhwmMetricsDescription = "peak resident set size (\"high water mark\")"

	VmstkMetricsName        = "vmstk_bytes"
	VmstkMetricsDescription = "size of stack segments"

	VmswapMetricsName        = "vmswap_bytes"
	VmswapMetricsDescription = "amount of swap used by anonymous private data (shmem swap usage is not included)"

	VmexeMetricsName        = "vmexe_bytes"
	VmexeMetricsDescription = "size of text segment"

	VmlibMetricsName        = "vmlib_bytes"
	VmlibMetricsDescription = "size of shared library code"

	VmlckMetricsName        = "vmlck_bytes"
	VmlckMetricsDescription = "locked memory size"

	VmpinMetricsName        = "vmpin_bytes"
	VmpinMetricsDescription = "pinned memory size"

	VmpteMetricsName        = "vmpte_bytes"
	VmpteMetricsDescription = "size of page table entries"

	k = 1024
)

type metricsInfo struct {
	Name  string
	Desc  string
	Unit  string
	Value func() int64
}

func getMemstatsMetrics() []*metricsInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return []*metricsInfo{
		{
			Name: AllocMetricsName,
			Desc: AllocMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.Alloc)
			},
		},
		{
			Name: TotalAllocMetricsName,
			Desc: TotalAllocMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.TotalAlloc)
			},
		},
		{
			Name: SysMetricsName,
			Desc: SysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.Sys)
			},
		},
		{
			Name: LookupsMetricsName,
			Desc: LookupsMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.Lookups)
			},
		},
		{
			Name: MallocsMetricsName,
			Desc: MallocsMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.Mallocs)
			},
		},
		{
			Name: FreesMetricsName,
			Desc: FreesMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.Frees)
			},
		},
		{
			Name: HeapAllocMetricsName,
			Desc: HeapAllocMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapAlloc)
			},
		},
		{
			Name: HeapSysMetricsName,
			Desc: HeapSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapSys)
			},
		},
		{
			Name: HeapIdleMetricsName,
			Desc: HeapIdleMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapIdle)
			},
		},
		{
			Name: HeapInuseMetricsName,
			Desc: HeapInuseMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapInuse)
			},
		},
		{
			Name: HeapReleasedMetricsName,
			Desc: HeapReleasedMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapReleased)
			},
		},
		{
			Name: HeapObjectsMetricsName,
			Desc: HeapObjectsMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.HeapObjects)
			},
		},
		{
			Name: StackInuseMetricsName,
			Desc: StackInuseMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.StackInuse)
			},
		},
		{
			Name: StackSysMetricsName,
			Desc: StackSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.StackSys)
			},
		},
		{
			Name: MspanInuseMetricsName,
			Desc: MspanInuseMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.MSpanInuse)
			},
		},
		{
			Name: MspanSysMetricsName,
			Desc: MspanSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.MSpanSys)
			},
		},
		{
			Name: McacheInuseMetricsName,
			Desc: McacheInuseMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.MCacheInuse)
			},
		},
		{
			Name: McacheSysMetricsName,
			Desc: McacheSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.MCacheSys)
			},
		},
		{
			Name: BuckHashSysMetricsName,
			Desc: BuckHashSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.BuckHashSys)
			},
		},
		{
			Name: GcSysMetricsName,
			Desc: GcSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.GCSys)
			},
		},
		{
			Name: OtherSysMetricsName,
			Desc: OtherSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.OtherSys)
			},
		},
		{
			Name: NextGcSysMetricsName,
			Desc: NextGcSysMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.NextGC)
			},
		},
		{
			Name: PauseTotalMsMetricsName,
			Desc: PauseTotalMsMetricsDescription,
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
			Name: NumGCMetricsName,
			Desc: NumGCMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.NumGC)
			},
		},
		{
			Name: NumForcedGCMetricsName,
			Desc: NumForcedGCMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.NumForcedGC)
			},
		},
		{
			Name: HeapWillReturnMetricsName,
			Desc: HeapWillReturnMetricsDescription,
			Unit: metrics.Bytes,
			Value: func() int64 {
				return int64(m.HeapIdle - m.HeapReleased)
			},
		},
		{
			Name: LiveObjectsMetricsName,
			Desc: LiveObjectsMetricsDescription,
			Unit: metrics.Dimensionless,
			Value: func() int64 {
				return int64(m.Mallocs - m.Frees)
			},
		},
	}
}

// skipcq: GO-R1005
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
					Name: VmpeakMetricsName,
					Desc: VmpeakMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmSize"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmsizeMetricsName,
					Desc: VmsizeMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmHWM"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmhwmMetricsName,
					Desc: VmhwmMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmRSS"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmrssMetricsName,
					Desc: VmrssMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmData"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmdataMetricsName,
					Desc: VmdataMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmStk"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmstkMetricsName,
					Desc: VmstkMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmExe"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmexeMetricsName,
					Desc: VmexeMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmLck"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmlckMetricsName,
					Desc: VmlckMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmLib"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmlibMetricsName,
					Desc: VmlibMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmPTE"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmpteMetricsName,
					Desc: VmpteMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmSwap"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmswapMetricsName,
					Desc: VmswapMetricsDescription,
					Unit: metrics.Bytes,
					Value: func() int64 {
						return f * k
					},
				})
			}
		case strings.HasPrefix(line, "VmPin"):
			f, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				m = append(m, &metricsInfo{
					Name: VmpinMetricsName,
					Desc: VmpinMetricsDescription,
					Unit: metrics.Bytes,
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

func (mm *memMetrics) View() ([]metrics.View, error) {
	mInfo := getMemstatsMetrics()
	if m, err := getProcStatusMetrics(mm.pid); err == nil {
		mInfo = append(mInfo, m...)
	}

	views := make([]metrics.View, 0, len(mInfo))
	for _, m := range mInfo {
		views = append(views, view.NewView(
			view.Instrument{
				Name:        m.Name,
				Description: m.Desc,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		))
	}
	return views, nil
}

func (mm *memMetrics) Register(m metrics.Meter) error {
	mInfo := getMemstatsMetrics()
	if metrics, err := getProcStatusMetrics(mm.pid); err == nil {
		mInfo = append(mInfo, metrics...)
	}

	instruments := make([]api.Int64Observable, 0, len(mInfo))
	for _, info := range mInfo {
		instrument, err := m.Int64ObservableGauge(
			info.Name,
			metrics.WithDescription(info.Desc),
			metrics.WithUnit(info.Unit),
		)
		if err == nil {
			instruments = append(instruments, instrument)
		}
	}
	oinsts := make([]api.Observable, 0, len(instruments))
	for _, instrument := range instruments {
		oinsts = append(oinsts, instrument)
	}
	_, err := m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			metrics := getMemstatsMetrics()
			if m, err := getProcStatusMetrics(mm.pid); err == nil {
				metrics = append(metrics, m...)
			}

			for i, instrument := range instruments {
				o.ObserveInt64(instrument, metrics[i].Value())
			}
			return nil
		},
		oinsts...,
	)
	return err
}
