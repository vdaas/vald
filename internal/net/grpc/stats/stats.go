//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package stats

import (
	"context"
	"strconv"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	statspb "github.com/vdaas/vald/apis/grpc/v1/rpc/stats"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
)

// CgroupMode represents the cgroup version
type CgroupMode int

const (
	Unknown CgroupMode = iota
	CGV1
	CGV2
)

const (
	cgroupBasePath = "/sys/fs/cgroup"
)

// CgroupMetrics holds raw values directly read from cgroup files
type CgroupMetrics struct {
	Mode CgroupMode

	// Raw values from cgroup files
	MemUsageBytes uint64
	MemLimitBytes uint64 // 0 means unlimited
	CPUUsageNano  uint64
	CPUQuotaUs    uint64 // 0 means unlimited
	CPUPeriodUs   uint64 // 0 if unknown
}

// CgroupStats holds calculated resource usage statistics ready for use
type CgroupStats struct {
	CPULimitCores    float64
	CPUUsageCores    float64
	MemoryLimitBytes uint64
	MemoryUsageBytes uint64
}

func Register(srv *grpc.Server) {
	statspb.RegisterStatsServer(srv, new(server))
}

type server struct {
	statspb.UnimplementedStatsServer
}

func (s *server) ResourceStats(
	ctx context.Context, _ *payload.Empty,
) (stats *payload.Info_Stats_ResourceStats, err error) {
	return GetResourceStats(ctx)
}

// GetResourceStats returns local resource stats measured from cgroup metrics.
func GetResourceStats(ctx context.Context) (stats *payload.Info_Stats_ResourceStats, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	ip := net.LoadLocalIP()
	if ip == "" {
		ip = "unknown"
	}

	stats = &payload.Info_Stats_ResourceStats{
		Name: hostname,
		Ip:   ip,
	}
	cgroupStats, err := measureCgroupStats(ctx)
	if err != nil {
		log.Warn("failed to measure cgroup stats", err)
	}
	if cgroupStats != nil {
		stats.CgroupStats = &payload.Info_Stats_CgroupStats{
			CpuLimitCores:    cgroupStats.CPULimitCores,
			CpuUsageCores:    cgroupStats.CPUUsageCores,
			MemoryLimitBytes: cgroupStats.MemoryLimitBytes,
			MemoryUsageBytes: cgroupStats.MemoryUsageBytes,
		}
	}
	return stats, nil
}

// measureCgroupStats orchestrates the process of sampling and calculating cgroup statistics.
func measureCgroupStats(ctx context.Context) (*CgroupStats, error) {
	// First sample: Read initial metrics from cgroup files (includes cumulative CPU usage)
	m1, err := readCgroupMetrics()
	if err != nil {
		return nil, errors.ErrCgroupFirstSampleFailed(err)
	}
	t1 := time.Now()

	// Wait 100ms to allow meaningful CPU usage accumulation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	// Second sample: Read metrics again to calculate CPU usage rate
	m2, err := readCgroupMetrics()
	if err != nil {
		return nil, errors.ErrCgroupSecondSampleFailed(err)
	}
	t2 := time.Now()

	// Calculate CPU usage rate from cumulative values: (usage2 - usage1) / time_delta
	cgroupStats := calculateCPUUsageCores(m1, m2, t2.Sub(t1))

	return &cgroupStats, nil
}

// readCgroupMetrics reads raw memory & CPU metrics depending on cgroup mode
func readCgroupMetrics() (metrics *CgroupMetrics, err error) {
	switch detectCgroupMode() {
	case CGV2:
		return readCgroupV2Metrics()
	case CGV1:
		return readCgroupV1Metrics()
	default:
		return nil, errors.ErrCgroupModeDetectionFailed
	}
}

// detectCgroupMode inspects /sys/fs/cgroup to detect cgroups mode
func detectCgroupMode() CgroupMode {
	// cgroups v2 unified mount has cgroup.controllers
	if file.Exists(file.Join(cgroupBasePath, "cgroup.controllers")) {
		return CGV2
	}

	data, err := file.ReadFile("/proc/self/cgroup")
	if err != nil {
		return Unknown
	}
	for _, line := range strings.Split(conv.Btoa(data), "\n") {
		// v2 line looks like: "0::/kubepods.slice/..."
		parts := strings.Split(line, ":")
		if len(parts) >= 3 && parts[1] == "" {
			return CGV2
		}
	}

	// If not v2, assume v1 on most distros
	return CGV1
}

// readCgroupV2Metrics reads cgroups v2 raw metrics
func readCgroupV2Metrics() (metrics *CgroupMetrics, err error) {
	// TODO: The current implementation directly uses /sys/fs/cgroup, but in some environments,
	// the cgroup namespace may not be separated per pod, resulting in reading values for the
	// entire node rather than per-pod values. Add functionality to specify appropriate paths
	// to ensure better isolation.
	data, err := file.ReadFile(file.Join(cgroupBasePath, "memory.current"))
	if err != nil {
		return nil, errors.ErrCgroupV2MemoryCurrentReadFailed(err)
	}
	memCur, err := strconv.ParseUint(strings.TrimSpace(conv.Btoa(data)), 10, 64)
	if err != nil {
		return nil, errors.ErrCgroupV2MemoryCurrentParseFailed(err)
	}

	data, err = file.ReadFile(file.Join(cgroupBasePath, "memory.max"))
	if err != nil {
		return nil, errors.ErrCgroupV2MemoryMaxReadFailed(err)
	}
	var memMax uint64
	memMaxData := strings.TrimSpace(conv.Btoa(data))
	if memMaxData == "max" {
		memMax = 0
	} else {
		memMax, err = strconv.ParseUint(memMaxData, 10, 64)
		if err != nil {
			return nil, errors.ErrCgroupV2MemoryMaxParseFailed(err)
		}
	}

	data, err = file.ReadFile(file.Join(cgroupBasePath, "cpu.stat"))
	if err != nil {
		return nil, errors.ErrCgroupV2CPUStatReadFailed(err)
	}
	var usageUS uint64
	for _, line := range strings.Split(conv.Btoa(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[0] == "usage_usec" {
			usageUS, err = strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				continue
			}
			break
		}
	}
	usageNS := usageUS * 1000

	data, err = file.ReadFile(file.Join(cgroupBasePath, "cpu.max"))
	if err != nil {
		return nil, errors.ErrCgroupV2CPUMaxReadFailed(err)
	}
	parts := strings.Fields(strings.TrimSpace(conv.Btoa(data)))
	if len(parts) != 2 {
		return nil, errors.ErrCgroupV2CPUMaxMalformed(strings.TrimSpace(conv.Btoa(data)))
	}
	var quotaUs uint64
	if parts[0] == "max" {
		quotaUs = 0
	} else {
		quotaUs, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return nil, errors.ErrCgroupV2CPUMaxParseQuotaFailed(err)
		}
	}
	periodUs, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, errors.ErrCgroupV2CPUMaxParsePeriodFailed(err)
	}

	metrics = &CgroupMetrics{
		Mode:          CGV2,
		MemUsageBytes: memCur,
		MemLimitBytes: memMax,
		CPUUsageNano:  usageNS,
		CPUQuotaUs:    quotaUs,
		CPUPeriodUs:   periodUs,
	}

	return metrics, nil
}

// readCgroupV1Metrics reads cgroups v1 raw metrics
func readCgroupV1Metrics() (metrics *CgroupMetrics, err error) {
	var memUsage uint64
	data, err := file.ReadFile(file.Join(cgroupBasePath, "memory", "memory.usage_in_bytes"))
	if err != nil {
		return nil, errors.ErrCgroupV1MemoryUsageReadFailed(err)
	}
	memUsage, err = strconv.ParseUint(strings.TrimSpace(conv.Btoa(data)), 10, 64)
	if err != nil {
		return nil, errors.ErrCgroupV1MemoryUsageParseFailed(err)
	}

	var memLimit uint64
	data, err = file.ReadFile(file.Join(cgroupBasePath, "memory", "memory.limit_in_bytes"))
	if err == nil {
		memLimit, _ = strconv.ParseUint(strings.TrimSpace(conv.Btoa(data)), 10, 64)
	}

	var cpuUsage uint64
	var cpuErr error
	cpuPaths := []string{
		file.Join(cgroupBasePath, "cpuacct", "cpuacct.usage"),
		file.Join(cgroupBasePath, "cpu,cpuacct", "cpuacct.usage"),
	}
	for _, path := range cpuPaths {
		data, err := file.ReadFile(path)
		if err == nil {
			cpuUsage, cpuErr = strconv.ParseUint(strings.TrimSpace(conv.Btoa(data)), 10, 64)
			if cpuErr == nil {
				break
			}
		}
		cpuErr = err
	}
	if cpuErr != nil {
		return nil, errors.ErrCgroupV1CPUUsageReadFailed(cpuErr)
	}

	var quota, period int64 = 0, 0
	quotaPaths := []string{
		file.Join(cgroupBasePath, "cpu", "cpu.cfs_quota_us"),
		file.Join(cgroupBasePath, "cpu,cpuacct", "cpu.cfs_quota_us"),
	}
	for _, path := range quotaPaths {
		data, err := file.ReadFile(path)
		if err == nil {
			quota, err = strconv.ParseInt(strings.TrimSpace(conv.Btoa(data)), 10, 64)
			if err == nil {
				if quota == -1 {
					quota = 0
				}
			}
			break
		}
	}
	periodPaths := []string{
		file.Join(cgroupBasePath, "cpu", "cpu.cfs_period_us"),
		file.Join(cgroupBasePath, "cpu,cpuacct", "cpu.cfs_period_us"),
	}
	for _, path := range periodPaths {
		data, err := file.ReadFile(path)
		if err == nil {
			period, _ = strconv.ParseInt(strings.TrimSpace(conv.Btoa(data)), 10, 64)
			break
		}
	}

	metrics = &CgroupMetrics{
		Mode:          CGV1,
		MemUsageBytes: memUsage,
		MemLimitBytes: memLimit,
		CPUUsageNano:  cpuUsage,
		CPUQuotaUs:    uint64(quota),
		CPUPeriodUs:   uint64(period),
	}

	return metrics, nil
}

// calculateCPUUsageCores computes CPU usage cores and other statistics from two raw metric samples.
func calculateCPUUsageCores(
	m1, m2 *CgroupMetrics, deltaTime time.Duration,
) (calculatedStats CgroupStats) {
	calculatedStats.MemoryLimitBytes = m2.MemLimitBytes
	calculatedStats.MemoryUsageBytes = m2.MemUsageBytes

	if m2.CPUQuotaUs > 0 && m2.CPUPeriodUs > 0 {
		calculatedStats.CPULimitCores = float64(m2.CPUQuotaUs) / float64(m2.CPUPeriodUs)
	} else {
		calculatedStats.CPULimitCores = 0
	}

	dtNano := deltaTime.Nanoseconds()
	if dtNano > 0 {
		dtUsage := int64(m2.CPUUsageNano) - int64(m1.CPUUsageNano)
		if dtUsage < 0 {
			dtUsage = 0
		}

		calculatedStats.CPUUsageCores = float64(dtUsage) / float64(dtNano)
	}

	return calculatedStats
}
