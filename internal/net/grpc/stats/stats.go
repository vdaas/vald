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

	"github.com/shirou/gopsutil/v4/docker"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	statspb "github.com/vdaas/vald/apis/grpc/v1/rpc/stats"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	ios "github.com/vdaas/vald/internal/os"
)

const (
	dockerCgroupBasePath = "/sys/fs/cgroup"
	procCgroupPath       = "/proc/1/cgroup"
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

	id := os.Getenv("MY_CONTAINER_ID")
	if err != nil {
		id = "unknown"
	}
	log.Debugf("container ID: %s", id)

	cpuUsage, err := getCPUUsage(id)
	if err != nil {
		log.Debugf("failed to get cpu usage: %v", err)
		cpuUsage = 0.0
	}
	log.Debugf("cpuUsage: %f", cpuUsage)

	memoryUsage, err := getMemoryUsage(id)
	if err != nil {
		log.Debugf("failed to get memory usage: %v", err)
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

func getCPUUsage(id string) (usage float64, err error) {
	stats, err := docker.CgroupCPU(id, dockerCgroupBasePath)
	if err != nil {
		log.Debugf("failed to get cgroup CPU stats: %v", err)
		return 0, nil
	}

	totalCPUTime := stats.User + stats.System
	return float64(totalCPUTime), nil
}

func getMemoryUsage(id string) (usage float64, err error) {
	memStat, err := docker.CgroupMem(id, dockerCgroupBasePath)
	if err != nil {
		return 0, nil
	}

	if memStat.MemLimitInBytes == 0 {
		return 0, nil
	}

	usage = (float64(memStat.MemUsageInBytes) / float64(memStat.MemLimitInBytes)) * 100.0
	return usage, nil
}
