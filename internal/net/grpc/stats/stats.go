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
	"strings"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	statspb "github.com/vdaas/vald/apis/grpc/v1/rpc/stats"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/os"
)

func Register(srv *grpc.Server) {
	ssrv := &server{}
	statspb.RegisterStatsServer(srv, ssrv)
}

type server struct {
	statspb.UnimplementedStatsServer
	lastCPUTotal   atomic.Uint64
	lastCPUIdle    atomic.Uint64
	lastCPUTime    atomic.Int64
	statsLoadedCnt atomic.Uint64
}

var (
	cpuStatsPath           = "/proc/stat"
	memoryStatsPath        = "/proc/meminfo"
	processMemoryStatsPath = "/proc/self/status"
)

func (s *server) ResourceStats(
	ctx context.Context, _ *payload.Empty,
) (*payload.Info_ResourceStats, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	log.Debugf("hostname: %s", hostname)

	ip := net.LoadLocalIP()
	if ip == "" {
		ip = "unknown"
	}
	log.Debugf("ip: %s", ip)

	cpuUsage, err := s.getCPUStats(cpuStatsPath)
	if err != nil {
		cpuUsage = 0.0
	}
	log.Debugf("cpuUsage: %f", cpuUsage)

	memoryUsage, err := getMemoryStats()
	if err != nil {
		memoryUsage = 0.0
	}
	log.Debugf("memoryUsage: %f", memoryUsage)

	return &payload.Info_ResourceStats{
		Name:        hostname,
		Ip:          ip,
		CpuUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
	}, nil
}

func (s *server) getCPUStats(path string) (usage float64, err error) {
	var data []byte
	data, err = file.ReadFile(path)
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
		return 0, errors.ErrCPULineNotFound()
	}

	fields := strings.Fields(cpuLine)
	if len(fields) < 8 {
		return 0, errors.ErrInvalidCPULineFormat()
	}

	var total, idle uint64
	for i := 1; i < len(fields) && i < 8; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return 0, errors.ErrCPUFieldParseFailed(i, err)
		}
		total += val
		if i == 4 {
			idle = val
		}
	}

	now := time.Now()

	s.lastCPUTotal.Store(total)
	s.lastCPUIdle.Store(idle)
	s.lastCPUTime.Store(now.UnixNano())

	if s.statsLoadedCnt.Add(1) <= 1 {
		return s.getCPUStats(path)
	}

	deltaTotal := total - s.lastCPUTotal.Load()

	if deltaTotal > 0 {
		deltaIdle := idle - s.lastCPUIdle.Load()
		usage = float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100.0
	}

	return usage, nil
}

func getMemoryStats() (usage float64, err error) {
	totalMem, err := getTotalMemory(memoryStatsPath)
	if err != nil {
		return 0, err
	}

	processMemory, err := getProcessMemory(processMemoryStatsPath)
	if err != nil {
		return 0, err
	}

	usage = float64(processMemory) / float64(totalMem) * 100.0

	return usage, nil
}

func getTotalMemory(path string) (usage uint64, err error) {
	var data []byte
	data, err = file.ReadFile(path)
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
					return 0, errors.ErrMemTotalParseFailed(err)
				}
				usage = memTotal * 1024
				return usage, nil
			}
		}
	}

	return 0, errors.ErrTotalMemoryNotFound()
}

func getProcessMemory(path string) (usage uint64, err error) {
	var data []byte
	data, err = file.ReadFile(path)
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
					return 0, errors.ErrVmRSSParseFailed(err)
				}
				usage = vmRSS * 1024
				return usage, nil
			}
		}
	}

	return 0, errors.ErrProcessMemoryNotFound()
}
