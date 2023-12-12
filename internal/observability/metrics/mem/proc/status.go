//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package proc

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
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

	kilo = 1024
)

type procStatusMetrics struct {
	pfile string
}

func New() metrics.Metric {
	return &procStatusMetrics{
		pfile: fmt.Sprintf("/proc/%d/status", os.Getpid()),
	}
}

func (*procStatusMetrics) View() ([]*metrics.View, error) {
	vmpeak, err := view.New(
		view.MatchInstrumentName(vmpeakMetricsName),
		view.WithSetDescription(vmpeakMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmsize, err := view.New(
		view.MatchInstrumentName(vmsizeMetricsName),
		view.WithSetDescription(vmsizeMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmdata, err := view.New(
		view.MatchInstrumentName(vmdataMetricsName),
		view.WithSetDescription(vmdataMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmrss, err := view.New(
		view.MatchInstrumentName(vmrssMetricsName),
		view.WithSetDescription(vmrssMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmhwm, err := view.New(
		view.MatchInstrumentName(vmhwmMetricsName),
		view.WithSetDescription(vmhwmMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmstk, err := view.New(
		view.MatchInstrumentName(vmstkMetricsName),
		view.WithSetDescription(vmstkMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmswap, err := view.New(
		view.MatchInstrumentName(vmswapMetricsName),
		view.WithSetDescription(vmswapMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmexe, err := view.New(
		view.MatchInstrumentName(vmexeMetricsName),
		view.WithSetDescription(vmexeMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmlib, err := view.New(
		view.MatchInstrumentName(vmlibMetricsName),
		view.WithSetDescription(vmlibMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmlck, err := view.New(
		view.MatchInstrumentName(vmlckMetricsName),
		view.WithSetDescription(vmlckMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmpin, err := view.New(
		view.MatchInstrumentName(vmpinMetricsName),
		view.WithSetDescription(vmpinMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	vmpte, err := view.New(
		view.MatchInstrumentName(vmpteMetricsName),
		view.WithSetDescription(vmpteMetricsDescription),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	return []*metrics.View{
		&vmpeak,
		&vmsize,
		&vmdata,
		&vmrss,
		&vmhwm,
		&vmstk,
		&vmswap,
		&vmexe,
		&vmlib,
		&vmlck,
		&vmpin,
		&vmpte,
	}, nil
}

func (p *procStatusMetrics) Register(m metrics.Meter) error {
	vmpeak, err := m.AsyncInt64().Gauge(
		vmpeakMetricsName,
		metrics.WithDescription(vmpeakMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmsize, err := m.AsyncInt64().Gauge(
		vmsizeMetricsName,
		metrics.WithDescription(vmsizeMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmdata, err := m.AsyncInt64().Gauge(
		vmdataMetricsName,
		metrics.WithDescription(vmdataMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmrss, err := m.AsyncInt64().Gauge(
		vmrssMetricsName,
		metrics.WithDescription(vmrssMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmhwm, err := m.AsyncInt64().Gauge(
		vmhwmMetricsName,
		metrics.WithDescription(vmhwmMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmstk, err := m.AsyncInt64().Gauge(
		vmstkMetricsName,
		metrics.WithDescription(vmstkMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmswap, err := m.AsyncInt64().Gauge(
		vmswapMetricsName,
		metrics.WithDescription(vmswapMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmexe, err := m.AsyncInt64().Gauge(
		vmexeMetricsName,
		metrics.WithDescription(vmexeMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmlib, err := m.AsyncInt64().Gauge(
		vmlibMetricsName,
		metrics.WithDescription(vmlibMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmlck, err := m.AsyncInt64().Gauge(
		vmlckMetricsName,
		metrics.WithDescription(vmlckMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmpin, err := m.AsyncInt64().Gauge(
		vmpinMetricsName,
		metrics.WithDescription(vmpinMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}

	vmpte, err := m.AsyncInt64().Gauge(
		vmpteMetricsName,
		metrics.WithDescription(vmpteMetricsDescription),
		metrics.WithUnit(metrics.Bytes),
	)
	if err != nil {
		return err
	}
	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			vmpeak,
			vmsize,
			vmdata,
			vmrss,
			vmhwm,
			vmstk,
			vmswap,
			vmexe,
			vmlib,
			vmlck,
			vmpin,
			vmpte,
		},
		func(ctx context.Context) {
			buf, err := os.ReadFile(p.pfile)
			if err != nil {
				return
			}
			lines := strings.Split(conv.Btoa(buf), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				switch {
				case strings.HasPrefix(line, "VmPeak"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmpeak.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmSize"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmsize.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmHWM"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmhwm.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmRSS"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmrss.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmData"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmdata.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmStk"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmstk.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmExe"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmexe.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmLck"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmlck.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmLib"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmlib.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmPTE"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmpte.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmSwap"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmswap.Observe(ctx, f*kilo)
					}
				case strings.HasPrefix(line, "VmPin"):
					f, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						vmpin.Observe(ctx, f*kilo)
					}
				}
			}
		},
	)
}
