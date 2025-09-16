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
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	statspb "github.com/vdaas/vald/apis/grpc/v1/rpc/stats"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	internalOS "github.com/vdaas/vald/internal/os"
)

var (
	cpuStatsMu   sync.RWMutex
	lastCPUTotal uint64
	lastCPUIdle  uint64
	lastCPUTime  time.Time
)

func Register(srv *grpc.Server) {
	ssrv := &server{}
	statspb.RegisterStatsServer(srv, ssrv)
}

type server struct {
	statspb.UnimplementedStatsServer
}

func (s *server) ResourceStats(
	ctx context.Context, _ *payload.Empty,
) (*payload.Info_ResourceStats, error) {
	hostname, err := internalOS.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	ip := net.LoadLocalIP()
	if ip == "" {
		ip = "unknown"
	}

	cpuUsage, err := getCPUStats()
	if err != nil {
		cpuUsage = 0.0
	}

	memoryUsage, err := getMemoryStats()
	if err != nil {
		memoryUsage = 0.0
	}

	return &payload.Info_ResourceStats{
		Name:        hostname,
		Ip:          ip,
		CpuUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
	}, nil
}

func getCPUStats() (float64, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	var cpuLine string
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			cpuLine = line
			break
		}
	}

	if cpuLine == "" {
		return 0, errors.New("cpu line not found in /proc/stat")
	}

	fields := strings.Fields(cpuLine)
	if len(fields) < 8 {
		return 0, errors.New("invalid cpu line format in /proc/stat")
	}

	var total, idle uint64
	for i := 1; i < len(fields) && i < 8; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to parse CPU field %d", i)
		}
		total += val
		if i == 4 {
			idle = val
		}
	}

	cpuStatsMu.Lock()
	defer cpuStatsMu.Unlock()

	var usage float64
	now := time.Now()

	if lastCPUTime.IsZero() || now.Sub(lastCPUTime) > 5*time.Minute {
		lastCPUTotal = total
		lastCPUIdle = idle
		lastCPUTime = now

		time.Sleep(100 * time.Millisecond)
		return getCPUStats()
	}

	deltaTotal := total - lastCPUTotal
	deltaIdle := idle - lastCPUIdle

	if deltaTotal > 0 {
		usage = float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100.0
	}

	if now.Sub(lastCPUTime) >= time.Second {
		lastCPUTotal = total
		lastCPUIdle = idle
		lastCPUTime = now
	}

	return usage, nil
}

func getMemoryStats() (float64, error) {
	totalMem, err := getTotalMemory()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get total memory")
	}

	processMemory, err := getProcessMemory()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get process memory")
	}

	usage := float64(processMemory) / float64(totalMem) * 100.0

	return usage, nil
}

func getTotalMemory() (uint64, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memTotal, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					return 0, errors.Wrap(err, "failed to parse MemTotal")
				}
				return memTotal * 1024, nil
			}
		}
	}

	return 0, errors.New("MemTotal not found in /proc/meminfo")
}

func getProcessMemory() (uint64, error) {
	data, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "VmRSS:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				vmRSS, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					return 0, errors.Wrap(err, "failed to parse VmRSS")
				}
				return vmRSS * 1024, nil
			}
		}
	}

	return 0, errors.New("VmRSS not found in /proc/self/status")
}
