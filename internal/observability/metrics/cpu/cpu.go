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

// Package cpu provides cpu metrics functions
package cpu

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
	"github.com/vdaas/vald/internal/observability/metrics"
)

const (
	cpuID         = "cpu_id"
	cpuVendorID   = "cpu_vendor_id"
	cpuFamily     = "cpu_family"
	cpuModel      = "cpu_model"
	cpuStepping   = "cpu_stepping"
	cpuPhysicalID = "cpu_physical_id"
	cpuCoreID     = "cpu_core_id"
	cpuCores      = "cpu_cores"
	cpuModelName  = "cpu_model_name"
	cpuMhz        = "cpu_mhz"
	cpuCacheSize  = "cpu_cache_size"
	cpuFlags      = "cpu_flags"
	cpuMicrocode  = "cpu_microcode"
)

type cpuInfo struct {
	process      *process.Process
	infoStats    []cpu.InfoStat
	infoStatKeys map[string]metrics.Key
	cpuInfo      metrics.Int64Measure
	cpuPercent   metrics.Float64Measure
	numThreads   metrics.Int64Measure
}

func New() (metrics.Metric, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, err
	}

	infoStats, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	infoStatKeys, err := infoStatsLabelKeys()
	if err != nil {
		return nil, err
	}

	return &cpuInfo{
		process:      p,
		infoStats:    infoStats,
		infoStatKeys: infoStatKeys,
		cpuInfo:      *metrics.Int64(metrics.ValdOrg+"/cpu/info", "cpu info", metrics.UnitDimensionless),
		cpuPercent:   *metrics.Float64(metrics.ValdOrg+"/cpu/utilization", "cpu utilization", metrics.UnitDimensionless),
		numThreads:   *metrics.Int64(metrics.ValdOrg+"/thread/count", "number of threads", metrics.UnitDimensionless),
	}, nil
}

func infoStatsLabelKeys() (map[string]metrics.Key, error) {
	keys := []string{
		cpuID,
		cpuVendorID,
		cpuFamily,
		cpuModel,
		cpuStepping,
		cpuPhysicalID,
		cpuCoreID,
		cpuCores,
		cpuModelName,
		cpuMhz,
		cpuCacheSize,
		cpuFlags,
		cpuMicrocode,
	}
	info := make(map[string]metrics.Key, len(keys))
	for _, kstr := range keys {
		k, err := metrics.NewKey(kstr)
		if err != nil {
			return nil, err
		}
		info[kstr] = k
	}
	return info, nil
}

func (c *cpuInfo) MeasurementsCount() int {
	cnt := 0
	rv := reflect.ValueOf(*c)
	for i := 0; i < rv.NumField(); i++ {
		if metrics.IsMeasureType(rv.Field(i).Type()) {
			cnt++
		}
	}
	return cnt
}

func (c *cpuInfo) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	cpuPercent, err := c.process.CPUPercent()
	if err != nil {
		return nil, err
	}
	numThreads, err := c.process.NumThreads()
	if err != nil {
		return nil, err
	}

	return []metrics.Measurement{
		c.cpuPercent.M(cpuPercent),
		c.numThreads.M(int64(numThreads)),
	}, nil
}

func (c *cpuInfo) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	ms := make([]metrics.MeasurementWithTags, 0, len(c.infoStats))
	for i := range c.infoStats {
		infoStat := &c.infoStats[i]
		ms = append(ms, metrics.MeasurementWithTags{
			Measurement: c.cpuInfo.M(int64(1)),
			Tags: map[metrics.Key]string{
				c.infoStatKeys[cpuID]:         strconv.Itoa(int(infoStat.CPU)),
				c.infoStatKeys[cpuVendorID]:   infoStat.VendorID,
				c.infoStatKeys[cpuFamily]:     infoStat.Family,
				c.infoStatKeys[cpuModel]:      infoStat.Model,
				c.infoStatKeys[cpuStepping]:   strconv.Itoa(int(infoStat.Stepping)),
				c.infoStatKeys[cpuPhysicalID]: infoStat.PhysicalID,
				c.infoStatKeys[cpuCoreID]:     infoStat.CoreID,
				c.infoStatKeys[cpuCores]:      strconv.Itoa(int(infoStat.Cores)),
				c.infoStatKeys[cpuModelName]:  infoStat.ModelName,
				c.infoStatKeys[cpuMhz]:        fmt.Sprintf("%g", infoStat.Mhz),
				c.infoStatKeys[cpuCacheSize]:  strconv.Itoa(int(infoStat.CacheSize)),
				// tags must be less than 255 characters
				c.infoStatKeys[cpuFlags]:     fmt.Sprintf("%.255s", fmt.Sprintf("%v", infoStat.Flags)),
				c.infoStatKeys[cpuMicrocode]: infoStat.Microcode,
			},
		})
	}
	return ms, nil
}

func (c *cpuInfo) View() []*metrics.View {
	keys := make([]metrics.Key, 0, len(c.infoStatKeys))
	for _, k := range c.infoStatKeys {
		keys = append(keys, k)
	}

	return []*metrics.View{
		&metrics.View{
			Name:        "cpu_info",
			Description: "cpu info",
			TagKeys:     keys,
			Measure:     &c.cpuInfo,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "cpu_utilization",
			Description: "cpu utilization",
			Measure:     &c.cpuPercent,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "thread_count",
			Description: "number of threads",
			Measure:     &c.numThreads,
			Aggregation: metrics.LastValue(),
		},
	}
}
