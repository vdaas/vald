// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

const (
	benchDim       = 128
	benchBaseNum   = 20000
	benchBatchNum  = 1000
	benchPoolSize  = 16
	benchPollEvery = 20 * time.Microsecond
	// benchNoAuto keeps the auto index/save loops from firing during a run.
	benchNoAuto = "96h"
)

func newBenchNGT(b *testing.B, stats, gc bool) (service.NGT, context.Context) {
	b.Helper()
	c := &config.NGT{
		Dimension:              benchDim,
		DistanceType:           "l2",
		ObjectType:             "float",
		CreationEdgeSize:       50,
		SearchEdgeSize:         50,
		EnableInMemoryMode:     true,
		EnableStatistics:       stats,
		EnableProactiveGC:      gc,
		AutoIndexDurationLimit: benchNoAuto,
		AutoIndexCheckDuration: benchNoAuto,
		AutoSaveIndexDuration:  benchNoAuto,
		AutoIndexLength:        10000000000,
		KVSDB:                  &config.KVSDB{Concurrency: 1},
		ErrorBufferLimit:       100,
	}
	// Statistics/proactive GC are applied as options (agentd.go), not from cfg.
	n, err := service.New(c.Bind(), service.WithEnableStatistics(stats), service.WithProactiveGC(gc))
	if err != nil {
		b.Fatalf("failed to create ngt service: %v", err)
	}
	ctx, cancel := context.WithCancel(b.Context())
	b.Cleanup(cancel)
	n.Start(ctx)
	for i, v := range vector.GaussianDistributedFloat32VectorGenerator(benchBaseNum, benchDim) {
		if err := n.Insert("base-"+strconv.Itoa(i), v); err != nil {
			b.Fatalf("failed to insert base vector: %v", err)
		}
	}
	if err := n.CreateIndex(ctx, benchPoolSize); err != nil {
		b.Fatalf("failed to create base index: %v", err)
	}
	return n, ctx
}

// BenchmarkCreateIndexSearchUnavailable measures, per CreateIndex of a
// benchBatchNum-vector batch, how long Search is refused (IsIndexing() is
// true, the exact condition Search checks) next to the whole CreateIndex call.
func BenchmarkCreateIndexSearchUnavailable(b *testing.B) {
	log.Init(log.WithLoggerType("nop"))
	for _, bc := range []struct {
		name      string
		stats, gc bool
	}{
		{name: "default"},
		{name: "statistics+proactive_gc", stats: true, gc: true},
	} {
		b.Run(bc.name, func(b *testing.B) {
			n, ctx := newBenchNGT(b, bc.stats, bc.gc)
			var unavailable, total time.Duration
			b.ResetTimer()
			for iter := range b.N {
				b.StopTimer()
				for j, v := range vector.GaussianDistributedFloat32VectorGenerator(benchBatchNum, benchDim) {
					if err := n.Insert("batch-"+strconv.Itoa(iter)+"-"+strconv.Itoa(j), v); err != nil {
						b.Fatalf("failed to insert batch vector: %v", err)
					}
				}
				stop := make(chan struct{})
				refused := make(chan time.Duration, 1)
				go func() {
					var sum time.Duration
					var since time.Time
					for {
						select {
						case <-stop:
							if !since.IsZero() {
								sum += time.Since(since)
							}
							refused <- sum
							return
						default:
						}
						if n.IsIndexing() {
							if since.IsZero() {
								since = time.Now()
							}
						} else if !since.IsZero() {
							sum += time.Since(since)
							since = time.Time{}
						}
						time.Sleep(benchPollEvery)
					}
				}()
				b.StartTimer()
				start := time.Now()
				err := n.CreateIndex(ctx, benchPoolSize)
				total += time.Since(start)
				b.StopTimer()
				close(stop)
				unavailable += <-refused
				if err != nil {
					b.Fatalf("failed to create index: %v", err)
				}
			}
			b.ReportMetric(float64(unavailable.Microseconds())/float64(b.N)/1e3, "unavailable-ms/op")
			b.ReportMetric(float64(total.Microseconds())/float64(b.N)/1e3, "createindex-ms/op")
		})
	}
}
