// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// BenchmarkCollector_Record measures the write throughput of the collector.
// It simulates multiple concurrent writers recording request results.
func BenchmarkCollector_Record(b *testing.B) {
	c := newBenchmarkCollector(b)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rr := GetRequestResult()
			// Random latency between 1ms and 101ms
			rr.Latency = time.Millisecond + time.Duration(rand.N(int64(100*time.Millisecond)))
			// Random status code (0-19)
			rr.Status = codes.Code(rand.N(uint32(MaxGRPCCodes)))

			c.Record(ctx, 0, rr)

			PutRequestResult(rr)
		}
	})
}

// BenchmarkCollector_Snapshot measures the read performance of generating a global snapshot.
// The collector is pre-filled with data to ensure the snapshot calculation is non-trivial.
func BenchmarkCollector_Snapshot(b *testing.B) {
	c := newBenchmarkCollector(b)
	ctx := context.Background()

	// Pre-fill with significant data to simulate a running state
	preFillCount := 100_000
	for i := 0; i < preFillCount; i++ {
		rr := GetRequestResult()
		rr.Latency = time.Millisecond + time.Duration(rand.N(int64(100*time.Millisecond)))
		rr.Status = codes.Code(rand.N(uint32(MaxGRPCCodes)))
		c.Record(ctx, 0, rr)
		PutRequestResult(rr)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.GlobalSnapshot()
	}
}

// BenchmarkCollector_Record_WithBackgroundSnapshot measures write performance while
// heavy read operations (Snapshots) are occurring in the background.
// This tests lock contention between Record and Snapshot.
func BenchmarkCollector_Record_WithBackgroundSnapshot(b *testing.B) {
	c := newBenchmarkCollector(b)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a background goroutine that aggressively triggers snapshots
	go func() {
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
	}()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rr := GetRequestResult()
			// Random latency between 1ms and 101ms
			rr.Latency = time.Millisecond + time.Duration(rand.N(int64(100*time.Millisecond)))
			// Random status code (0-19)
			rr.Status = codes.Code(rand.N(uint32(MaxGRPCCodes)))

			c.Record(ctx, 0, rr)

			PutRequestResult(rr)
		}
	})
}
