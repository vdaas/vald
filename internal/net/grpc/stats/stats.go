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

package stats

import (
	"context"
	"strconv"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	statspb "github.com/vdaas/vald/apis/grpc/v1/rpc/stats"
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
	CpuLimitCores    float64
	CpuUsageCores    float64
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
) (stats *payload.Info_ResourceStats, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	ip := net.LoadLocalIP()
	if ip == "" {
		ip = "unknown"
	}


	stats = &payload.Info_ResourceStats{
		Name: hostname,
		Ip:   ip,
	}
	cgroupStats, err := measureCgroupStats(ctx)
	if err == nil && cgroupStats != nil {
		stats.CgroupStats = &payload.Info_CgroupStats{
			CpuLimitCores:    cgroupStats.CpuLimitCores,
			CpuUsageCores:    cgroupStats.CpuUsageCores,
			MemoryLimitBytes: cgroupStats.MemoryLimitBytes,
			MemoryUsageBytes: cgroupStats.MemoryUsageBytes,
		}
	}
	return stats, err
}

// measureCgroupStats orchestrates the process of sampling and calculating cgroup statistics.
func measureCgroupStats() (stats *CgroupStats, err error) {
	// First sample: Read initial metrics from cgroup files (includes cumulative CPU usage)
	m1, err := readCgroupMetrics()
	if err != nil {
		return nil, errors.ErrCgroupFirstSampleFailed(err)
	}
	t1 := time.Now()

	// Wait 100ms to allow meaningful CPU usage accumulation
	time.Sleep(100 * time.Millisecond)

	// Second sample: Read metrics again to calculate CPU usage rate
	m2, err := readCgroupMetrics()
	if err != nil {
		return nil, errors.ErrCgroupSecondSampleFailed(err)
	}
	t2 := time.Now()

	// Calculate CPU usage rate from cumulative values: (usage2 - usage1) / time_delta
	stats = calculateCpuUsageCores(m1, m2, t2.Sub(t1))

	return stats, nil
}

// readCgroupMetrics reads raw memory & CPU metrics depending on cgroup mode
func readCgroupMetrics() (metrics *CgroupMetrics, err error) {
	mode := detectCgroupMode()
	switch mode {
	case CGV2:
		metrics, err = readCgroupV2Metrics()
	case CGV1:
		metrics, err = readCgroupV1Metrics()
	default:
		err = errors.ErrCgroupModeDetectionFailed()
	}

	return metrics, err
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
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
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
	data, err := file.ReadFile(file.Join(cgroupBasePath, "memory.current"))
	if err != nil {
		return nil, errors.ErrCgroupV2MemoryCurrentRead(err)
	}
	memCur, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return nil, errors.ErrCgroupV2MemoryCurrentParse(err)
	}

	data, err = file.ReadFile(file.Join(cgroupBasePath, "memory.max"))
	if err != nil {
		return nil, errors.ErrCgroupV2MemoryMaxRead(err)
	}
	memMaxStr := strings.TrimSpace(string(data))
	var memMax uint64
	if memMaxStr == "max" {
		memMax = 0
	} else {
		memMax, err = strconv.ParseUint(memMaxStr, 10, 64)
		if err != nil {
			return nil, errors.ErrCgroupV2MemoryMaxParse(err)
		}
	}

	data, err = file.ReadFile(file.Join(cgroupBasePath, "cpu.stat"))
	if err != nil {
		return nil, errors.ErrCgroupV2CPUStatRead(err)
	}
	var usageUS uint64
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
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
	if usageUS == 0 {
		return nil, errors.ErrCgroupV2CPUStatMissingUsage()
	}
	usageNS := usageUS * 1000

	data, err = file.ReadFile(file.Join(cgroupBasePath, "cpu.max"))
	if err != nil {
		return nil, errors.ErrCgroupV2CPUMaxRead(err)
	}
	val := strings.TrimSpace(string(data))
	parts := strings.Fields(val)
	if len(parts) != 2 {
		return nil, errors.ErrCgroupV2CPUMaxMalformed(val)
	}
	var quotaUs uint64
	if parts[0] == "max" {
		quotaUs = 0
	} else {
		quotaUsInt, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return nil, errors.ErrCgroupV2CPUMaxParseQuota(err)
		}
		quotaUs = quotaUsInt
	}
	periodUs, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, errors.ErrCgroupV2CPUMaxParsePeriod(err)
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
	var memErr error
	memoryPaths := []string{
		file.Join(cgroupBasePath, "memory", "memory.usage_in_bytes"),
	}
	for _, path := range memoryPaths {
		data, err := file.ReadFile(path)
		if err == nil {
			memUsage, memErr = strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
			if memErr == nil {
				break
			}
		}
		memErr = err
	}
	if memErr != nil {
		return nil, errors.ErrCgroupV1MemoryUsageReadFailed(memErr)
	}

	var memLimit uint64
	limitPaths := []string{
		file.Join(cgroupBasePath, "memory", "memory.limit_in_bytes"),
	}
	for _, path := range limitPaths {
		data, err := file.ReadFile(path)
		if err == nil {
			memLimit, _ = strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
			break
		}
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
			cpuUsage, cpuErr = strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
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
	periodPaths := []string{
		file.Join(cgroupBasePath, "cpu", "cpu.cfs_period_us"),
		file.Join(cgroupBasePath, "cpu,cpuacct", "cpu.cfs_period_us"),
	}
	for _, path := range quotaPaths {
		data, err := file.ReadFile(path)
		if err == nil {
			quotaInt, parseErr := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
			if parseErr == nil {
				if quotaInt == -1 {
					quota = 0
				} else {
					quota = quotaInt
				}
			}
			break
		}
	}
	for _, path := range periodPaths {
		data, err := file.ReadFile(path)
		if err == nil {
			period, _ = strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
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

// calculateCpuUsageCores computes CPU usage cores and other statistics from two raw metric samples.
func calculateCpuUsageCores(
	m1, m2 *CgroupMetrics, deltaTime time.Duration,
) (calculatedStats *CgroupStats) {
	calculatedStats = &CgroupStats{}

	calculatedStats.MemoryLimitBytes = m2.MemLimitBytes
	calculatedStats.MemoryUsageBytes = m2.MemUsageBytes

	if m2.CPUQuotaUs > 0 && m2.CPUPeriodUs > 0 {
		calculatedStats.CpuLimitCores = float64(m2.CPUQuotaUs) / float64(m2.CPUPeriodUs)
	} else {
		calculatedStats.CpuLimitCores = 0
	}

	dtNano := deltaTime.Nanoseconds()
	if dtNano > 0 {
		dtUsage := int64(m2.CPUUsageNano) - int64(m1.CPUUsageNano)
		if dtUsage < 0 {
			dtUsage = 0
		}

		calculatedStats.CpuUsageCores = float64(dtUsage) / float64(dtNano)
	}

	return calculatedStats
}
