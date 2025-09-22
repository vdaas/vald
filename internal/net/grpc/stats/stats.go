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
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	statspb "github.com/vdaas/vald/apis/grpc/v1/rpc/stats"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	ios "github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
)

// CgroupMode represents the cgroup version
type CgroupMode int

const (
	Unknown CgroupMode = iota
	CGV1
	CGV2
)

// CgroupMetrics holds raw values directly read from cgroup files
type CgroupMetrics struct {
	Mode CgroupMode

	// Raw values from cgroup files
	MemUsageBytes uint64
	MemLimitBytes uint64
	CPUUsageNano  uint64
	CPUQuotaUs    int64 // -1 means "max" / unlimited
	CPUPeriodUs   int64 // 0 if unknown
}

// CgroupStats holds calculated resource usage statistics ready for use
type CgroupStats struct {
	// Calculated values for protobuf (ready to use)
	CpuLimit         float64 // CPU cores available
	CpuUsage         float64 // CPU usage in cores (not percentage)
	MemoryLimitBytes uint64  // Memory limit in bytes
	MemoryUsageBytes uint64  // Memory usage in bytes
}

func Register(srv *grpc.Server) {
	ssrv := &server{}
	statspb.RegisterStatsServer(srv, ssrv)
}

type server struct {
	statspb.UnimplementedStatsServer
}

func (s *server) ResourceStats(
	ctx context.Context, _ *payload.Empty,
) (stats *payload.Info_ResourceStats, err error) {
	hostname, err := ios.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	log.Debugf("hostname: %s", hostname)

	ip := net.LoadLocalIP()
	if ip == "" {
		ip = "unknown"
	}
	log.Debugf("ip: %s", ip)

	cgroupStats, err := measureCgroupStats()
	if err != nil {
		log.Debugf("failed to get cgroup stats: %v", err)
		stats = &payload.Info_ResourceStats{
			Name: hostname,
			Ip:   ip,
			CgroupStats: &payload.Info_CgroupStats{
				CpuLimit:         0,
				CpuUsage:         0,
				MemoryLimitBytes: 0,
				MemoryUsageBytes: 0,
			},
		}
		err = nil
		return
	}

	stats = &payload.Info_ResourceStats{
		Name: hostname,
		Ip:   ip,
		CgroupStats: &payload.Info_CgroupStats{
			CpuLimit:         cgroupStats.CpuLimit,
			CpuUsage:         cgroupStats.CpuUsage,
			MemoryLimitBytes: cgroupStats.MemoryLimitBytes,
			MemoryUsageBytes: cgroupStats.MemoryUsageBytes,
		},
	}
	return
}

// measureCgroupStats orchestrates the process of sampling and calculating cgroup statistics.
func measureCgroupStats() (stats *CgroupStats, err error) {
	// First sample
	m1, err := readCgroupMetrics()
	if err != nil {
		err = fmt.Errorf("failed to get first sample: %w", err)
		return
	}
	t1 := time.Now()

	time.Sleep(100 * time.Millisecond)

	// Second sample
	m2, err := readCgroupMetrics()
	if err != nil {
		err = fmt.Errorf("failed to get second sample: %w", err)
		return
	}
	t2 := time.Now()

	// Calculate final values from the two samples.
	stats = calculateCgroupStats(m1, m2, t2.Sub(t1))
	return
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
		err = errors.New("unable to detect cgroups mode")
	}
	return
}

// detectCgroupMode inspects /sys/fs/cgroup to detect cgroups mode
func detectCgroupMode() CgroupMode {
	// cgroups v2 unified mount has cgroup.controllers
	if file.Exists("/sys/fs/cgroup/cgroup.controllers") {
		return CGV2
	}
	// Fallback: parse /proc/self/cgroup
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
	base := "/sys/fs/cgroup"

	data, err := file.ReadFile(filepath.Join(base, "memory.current"))
	if err != nil {
		err = fmt.Errorf("v2 read memory.current: %w", err)
		return
	}
	memCur, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		err = fmt.Errorf("v2 parse memory.current: %w", err)
		return
	}

	data, err = file.ReadFile(filepath.Join(base, "memory.max"))
	if err != nil {
		err = fmt.Errorf("v2 read memory.max: %w", err)
		return
	}
	memMaxStr := strings.TrimSpace(string(data))
	var memMax uint64
	if memMaxStr == "max" {
		memMax = 0
	} else {
		memMax, err = strconv.ParseUint(memMaxStr, 10, 64)
		if err != nil {
			err = fmt.Errorf("v2 parse memory.max: %w", err)
			return
		}
	}

	data, err = file.ReadFile(filepath.Join(base, "cpu.stat"))
	if err != nil {
		err = fmt.Errorf("v2 read cpu.stat: %w", err)
		return
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
		err = fmt.Errorf("v2 cpu.stat missing usage_usec")
		return
	}
	usageNS := usageUS * 1000

	data, err = file.ReadFile(filepath.Join(base, "cpu.max"))
	if err != nil {
		err = fmt.Errorf("v2 read cpu.max: %w", err)
		return
	}
	val := strings.TrimSpace(string(data))
	parts := strings.Fields(val)
	if len(parts) != 2 {
		err = fmt.Errorf("v2 cpu.max malformed: %q", val)
		return
	}
	var quotaUs int64
	if parts[0] == "max" {
		quotaUs = -1
	} else {
		quotaUs, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			err = fmt.Errorf("v2 cpu.max parse quota: %w", err)
			return
		}
	}
	periodUs, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("v2 cpu.max parse period: %w", err)
		return
	}

	metrics = &CgroupMetrics{
		Mode:          CGV2,
		MemUsageBytes: memCur,
		MemLimitBytes: memMax,
		CPUUsageNano:  usageNS,
		CPUQuotaUs:    quotaUs,
		CPUPeriodUs:   periodUs,
	}
	return
}

// readCgroupV1Metrics reads cgroups v1 raw metrics
func readCgroupV1Metrics() (metrics *CgroupMetrics, err error) {
	base := "/sys/fs/cgroup"

	// Memory usage - try different paths
	memoryPaths := []string{
		filepath.Join(base, "memory", "memory.usage_in_bytes"),
	}
	var memUsage uint64
	var memErr error
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
		err = fmt.Errorf("v1 memory usage read failed: %w", memErr)
		return
	}

	// Memory limit
	var memLimit uint64
	limitPaths := []string{
		filepath.Join(base, "memory", "memory.limit_in_bytes"),
	}
	for _, path := range limitPaths {
		data, err := file.ReadFile(path)
		if err == nil {
			memLimit, _ = strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
			break
		}
	}

	// CPU usage - try different paths
	cpuPaths := []string{
		filepath.Join(base, "cpuacct", "cpuacct.usage"),
		filepath.Join(base, "cpu,cpuacct", "cpuacct.usage"),
	}
	var cpuUsage uint64
	var cpuErr error
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
		err = fmt.Errorf("v1 cpuacct.usage read failed: %w", cpuErr)
		return
	}

	// CPU quota/period
	var quota, period int64 = -1, 0
	quotaPaths := []string{
		filepath.Join(base, "cpu", "cpu.cfs_quota_us"),
		filepath.Join(base, "cpu,cpuacct", "cpu.cfs_quota_us"),
	}
	periodPaths := []string{
		filepath.Join(base, "cpu", "cpu.cfs_period_us"),
		filepath.Join(base, "cpu,cpuacct", "cpu.cfs_period_us"),
	}

	for _, path := range quotaPaths {
		data, err := file.ReadFile(path)
		if err == nil {
			quota, _ = strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
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
		CPUQuotaUs:    quota,
		CPUPeriodUs:   period,
	}
	return
}

// calculateCgroupStats computes the final, user-facing statistics from two raw metric samples.
func calculateCgroupStats(
	m1, m2 *CgroupMetrics, deltaTime time.Duration,
) (finalStats *CgroupStats) {
	finalStats = &CgroupStats{}

	// 1. Calculate Memory stats (using the latest sample, m2)
	finalStats.MemoryLimitBytes = m2.MemLimitBytes
	finalStats.MemoryUsageBytes = m2.MemUsageBytes

	// 2. Calculate CPU limit (using the latest sample, m2)
	if m2.CPUQuotaUs > 0 && m2.CPUPeriodUs > 0 {
		finalStats.CpuLimit = float64(m2.CPUQuotaUs) / float64(m2.CPUPeriodUs)
	} else {
		finalStats.CpuLimit = float64(runtime.NumCPU())
	}
	// Avoid division by zero
	if finalStats.CpuLimit <= 0 {
		finalStats.CpuLimit = 1.0
	}

	// 3. Calculate CPU usage in cores (using the delta between m1 and m2)
	dtNano := deltaTime.Nanoseconds()
	if dtNano > 0 {
		deltaUsage := int64(m2.CPUUsageNano) - int64(m1.CPUUsageNano)
		// Handle counter reset
		if deltaUsage < 0 {
			deltaUsage = 0
		}

		// CPU usage in cores: (CPU time used per nanosecond)
		finalStats.CpuUsage = float64(deltaUsage) / float64(dtNano)
	}

	return
}
