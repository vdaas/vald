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
package metrics

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test"
)

// newBenchmarkCollector creates a collector with full features enabled for benchmarking.
func newBenchmarkCollector(b *testing.B) Collector {
	b.Helper()
	c, err := NewCollector(
		WithTimeScale("1m_window", time.Minute, 60),
		WithRangeScale("request_range", 100, 100),
		WithLatencyHistogram(),
		WithQueueWaitHistogram(),
		WithLatencyTDigest(),
		WithQueueWaitTDigest(),
		WithExemplar(),
	)
	if err != nil {
		b.Fatalf("failed to create collector: %v", err)
	}
	return c
}

// collectorBench is the per-case benchmark body driven by BenchmarkCollector
// below through the generic test framework (test.BenchmarkCase instantiates
// internal/test's table-driven runner with *testing.B instead of *testing.T).
type collectorBench func(b *testing.B, c Collector)

// recordParallel simulates multiple concurrent writers recording request
// results, measuring the write throughput of the collector.
func recordParallel(b *testing.B, c Collector) {
	b.Helper()
	ctx := b.Context()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rr := GetRequestResult()
			// Random latency between 1ms and 101ms
			rr.Latency = time.Millisecond + time.Duration(rand.N(int64(100*time.Millisecond))) // skipcq: GSC-G404
			// Random status code (0-19)
			rr.Status = codes.Code(rand.N(uint32(MaxGRPCCodes))) // skipcq: GSC-G404

			c.Record(ctx, 0, rr)

			PutRequestResult(rr)
		}
	})
	// Exclude the framework teardown (goroutine-leak verification) from the
	// measured window; unlike b.Loop, RunParallel does not confine the timer.
	b.StopTimer()
}

// snapshot measures the read performance of generating a global snapshot.
// The collector is pre-filled with data to ensure the snapshot calculation
// is non-trivial.
func snapshot(b *testing.B, c Collector) {
	b.Helper()
	ctx := b.Context()

	// Pre-fill with significant data to simulate a running state
	preFillCount := 100_000
	for range preFillCount {
		rr := GetRequestResult()
		rr.Latency = time.Millisecond + time.Duration(rand.N(int64(100*time.Millisecond))) // skipcq: GSC-G404
		rr.Status = codes.Code(rand.N(uint32(MaxGRPCCodes)))                               // skipcq: GSC-G404
		c.Record(ctx, 0, rr)
		PutRequestResult(rr)
	}

	b.ReportAllocs()
	b.ResetTimer()

	// b.Loop stops the timer itself after the final iteration, so the
	// framework teardown is already excluded here (no StopTimer needed).
	for b.Loop() {
		_ = c.GlobalSnapshot()
	}
}

// recordWithBackgroundSnapshot measures write performance while heavy read
// operations (Snapshots) are occurring in the background. This tests lock
// contention between Record and Snapshot.
func recordWithBackgroundSnapshot(b *testing.B, c Collector) {
	b.Helper()
	// The snapshotter is cancelled and awaited before returning so the
	// framework's goroutine-leak verification sees a clean state.
	ctx, cancel := context.WithCancel(b.Context())
	var wg sync.WaitGroup
	wg.Go(func() {
		// High frequency snapshotting
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = c.GlobalSnapshot()
			}
		}
	})

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rr := GetRequestResult()
			// Random latency between 1ms and 101ms
			rr.Latency = time.Millisecond + time.Duration(rand.N(int64(100*time.Millisecond))) // skipcq: GSC-G404
			// Random status code (0-19)
			rr.Status = codes.Code(rand.N(uint32(MaxGRPCCodes))) // skipcq: GSC-G404

			c.Record(ctx, 0, rr)

			PutRequestResult(rr)
		}
	})
	b.StopTimer()
	cancel()
	wg.Wait()
}

// BenchmarkCollector drives the collector benchmarks through the same
// table-driven framework as the unit tests: test.Run is instantiated with
// X = *testing.B (via the test.BenchmarkCase alias), so each case becomes a
// b.Run sub-benchmark with the shared collector setup living in do.
func BenchmarkCollector(b *testing.B) {
	if err := test.Run(b.Context(), b, func(b *testing.B, run collectorBench) (struct{}, error) {
		b.Helper()
		run(b, newBenchmarkCollector(b))
		return struct{}{}, nil
	}, []test.BenchmarkCase[struct{}, collectorBench]{
		{Name: "record_parallel", Args: recordParallel},
		{Name: "snapshot", Args: snapshot},
		{Name: "record_with_background_snapshot", Args: recordWithBackgroundSnapshot},
	}...); err != nil {
		b.Error(err)
	}
}
