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
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

type resourceLoad struct {
	cpu    float64
	memory float64
}

type resourceLoadSnapshot struct {
	collectedAt time.Time
	loads       map[string]resourceLoad
}

type orca struct {
	refreshInterval time.Duration
	reportTTL       time.Duration
	replica         int
	snapshot        atomic.Pointer[resourceLoadSnapshot]
	readLogOnce     sync.Once
	writeLogOnce    sync.Once
}

func newORCA(refreshInterval, reportTTL time.Duration, replica int) *orca {
	if replica < 1 {
		replica = 1
	}
	return &orca{
		refreshInterval: refreshInterval,
		reportTTL:       reportTTL,
		replica:         replica,
	}
}

func (o *orca) start(ctx context.Context, client grpc.Client) {
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

func (o *orca) collect(ctx context.Context, client grpc.Client) {
	loads := make(map[string]resourceLoad)
	var mu sync.Mutex
	err := client.RangeConcurrent(ctx, -1, func(
		ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
	) error {
		stats, err := vc.NewValdClient(conn).ResourceStats(ctx, new(payload.Empty), copts...)
		if err != nil {
			log.Warnf("ORCA resource stats collection failed for %s: %v", addr, err)
			return nil
		}
		if stats == nil || stats.GetCgroupStats() == nil {
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
	if len(loads) == 0 {
		return
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

func (o *orca) order(addrs []string) []string {
	snapshot := o.snapshot.Load()
	if snapshot == nil || time.Since(snapshot.collectedAt) > o.reportTTL {
		return nil
	}
	type candidate struct {
		addr  string
		load  resourceLoad
		known bool
	}
	candidates := make([]candidate, 0, len(addrs))
	for _, addr := range addrs {
		load, known := snapshot.loads[addr]
		candidates = append(candidates, candidate{
			addr:  addr,
			load:  load,
			known: known,
		})
	}
	slices.SortStableFunc(candidates, func(left, right candidate) int {
		if left.known != right.known {
			if left.known {
				return -1
			}
			return 1
		}
		leftScore := max(left.load.cpu, left.load.memory)
		rightScore := max(right.load.cpu, right.load.memory)
		return cmp.Compare(leftScore, rightScore)
	})
	ordered := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		ordered = append(ordered, candidate.addr)
	}
	return ordered
}

func (o *orca) read(addrs []string) []string {
	ordered := o.order(addrs)
	if len(ordered) == 0 {
		return nil
	}
	limit := len(ordered) - o.replica + 1
	if limit < 1 {
		limit = 1
	}
	selected := ordered[:limit]
	o.readLogOnce.Do(func() {
		log.Infof(
			"ORCA READ fanout selected %d of %d agents with replica %d: %v",
			len(selected),
			len(ordered),
			o.replica,
			selected,
		)
	})
	return selected
}

func (o *orca) logWrite(selected []string, total, requested int) {
	o.writeLogOnce.Do(func() {
		log.Infof(
			"ORCA WRITE placement selected %d of %d agents for replica %d: %v",
			len(selected),
			total,
			requested,
			selected,
		)
	})
}
