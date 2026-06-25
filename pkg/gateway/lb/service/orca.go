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

package service

import (
	"cmp"
	"context"
	"slices"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

type ORCAPolicy struct {
	MinFanout       int
	MaxFanout       int
	CPUThreshold    float64
	MemoryThreshold float64
}

type resourceLoad struct {
	cpu    float64
	memory float64
}

type resourceLoadSnapshot struct {
	collectedAt time.Time
	loads       map[string]resourceLoad
}

type orca struct {
	enabled         bool
	refreshInterval time.Duration
	reportTTL       time.Duration
	read            ORCAPolicy
	write           ORCAPolicy
	snapshot        atomic.Pointer[resourceLoadSnapshot]
}

func newORCA(refreshInterval, reportTTL time.Duration, read, write ORCAPolicy) *orca {
	return &orca{
		enabled:         true,
		refreshInterval: refreshInterval,
		reportTTL:       reportTTL,
		read:            normalizeORCAPolicy(read),
		write:           normalizeORCAPolicy(write),
	}
}

func normalizeORCAPolicy(policy ORCAPolicy) ORCAPolicy {
	if policy.MinFanout < 1 {
		policy.MinFanout = 1
	}
	if policy.MaxFanout > 0 && policy.MaxFanout < policy.MinFanout {
		policy.MaxFanout = policy.MinFanout
	}
	return policy
}

func (o *orca) start(ctx context.Context, client discoverer.Client) {
	o.collect(ctx, client)
	ticker := time.NewTicker(o.refreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			o.collect(ctx, client)
		}
	}
}

func (o *orca) collect(ctx context.Context, client discoverer.Client) {
	loads := make(map[string]resourceLoad)
	var mu sync.Mutex
	collect := func(gc grpc.Client) {
		err := gc.RangeConcurrent(ctx, -1, func(
			ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
		) error {
			stats, err := vc.NewValdClient(conn).ResourceStats(ctx, new(payload.Empty), copts...)
			if err != nil || stats == nil || stats.GetCgroupStats() == nil {
				return nil
			}
			cgroup := stats.GetCgroupStats()
			load := resourceLoad{
				cpu:    utilization(cgroup.GetCpuUsageCores(), cgroup.GetCpuLimitCores()),
				memory: utilization(float64(cgroup.GetMemoryUsageBytes()), float64(cgroup.GetMemoryLimitBytes())),
			}
			mu.Lock()
			loads[addr] = load
			mu.Unlock()
			return nil
		})
		if err != nil {
			log.Warnf("ORCA resource stats collection failed: %v", err)
		}
	}
	writeClient := client.GetClient()
	collect(writeClient)
	readClient := client.GetReadReplicaClient()
	if readClient != nil && readClient != writeClient {
		collect(readClient)
	}
	o.snapshot.Store(&resourceLoadSnapshot{
		collectedAt: time.Now(),
		loads:       loads,
	})
}

func utilization(usage, limit float64) float64 {
	if limit <= 0 {
		return 0
	}
	return usage / limit
}

func (o *orca) selectAddrs(kind BroadCastKind, addrs []string) []string {
	snapshot := o.snapshot.Load()
	if snapshot == nil || time.Since(snapshot.collectedAt) > o.reportTTL {
		return nil
	}
	policy := o.read
	if kind == WRITE {
		policy = o.write
	}
	type candidate struct {
		addr     string
		load     resourceLoad
		known    bool
		accepted bool
	}
	candidates := make([]candidate, 0, len(addrs))
	for _, addr := range addrs {
		load, known := snapshot.loads[addr]
		accepted := !known ||
			(policy.CPUThreshold <= 0 || load.cpu <= policy.CPUThreshold) &&
			(policy.MemoryThreshold <= 0 || load.memory <= policy.MemoryThreshold)
		candidates = append(candidates, candidate{
			addr:     addr,
			load:     load,
			known:    known,
			accepted: accepted,
		})
	}
	slices.SortFunc(candidates, func(left, right candidate) int {
		if left.accepted != right.accepted {
			if left.accepted {
				return -1
			}
			return 1
		}
		if left.known != right.known {
			if left.known {
				return -1
			}
			return 1
		}
		leftScore := max(left.load.cpu, left.load.memory)
		rightScore := max(right.load.cpu, right.load.memory)
		if scoreOrder := cmp.Compare(leftScore, rightScore); scoreOrder != 0 {
			return scoreOrder
		}
		return cmp.Compare(left.addr, right.addr)
	})
	selected := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if !candidate.accepted && len(selected) >= policy.MinFanout {
			continue
		}
		selected = append(selected, candidate.addr)
		if policy.MaxFanout > 0 && len(selected) >= policy.MaxFanout {
			break
		}
	}
	return selected
}
