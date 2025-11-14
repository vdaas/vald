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

package metrics

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

func TestCollector_WithCustomCountersAndScales(t *testing.T) {
	counterNames := []string{"c1", "c2", "c3"}
	c, err := NewCollector(
		WithCustomCounters(counterNames...),
		WithRangeScale("test-range", 10, 10),
		WithTimeScale("test-time", 5, 10),
	)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	if len(c.rangeScales) != 1 {
		t.Fatal("expected 1 range scale")
	}
	if len(c.timeScales) != 1 {
		t.Fatal("expected 1 time scale")
	}

	// Verify that the slot has the correct number of counters
	numCounters := len(counterNames)
	if len(c.rangeScales[0].slots[0].Counters) != numCounters {
		t.Errorf("range scale slot has %d counters, want %d", len(c.rangeScales[0].slots[0].Counters), numCounters)
	}
	if len(c.timeScales[0].slots[0].Counters) != numCounters {
		t.Errorf("time scale slot has %d counters, want %d", len(c.timeScales[0].slots[0].Counters), numCounters)
	}

	// Record a request to populate a slot
	rr := &RequestResult{
		EndedAt: time.Now(),
		Latency: 100 * time.Millisecond,
	}
	ctx := WithRequestID(context.Background(), 5)
	c.Record(ctx, rr)

	// Check snapshot
	rangeSnap := c.RangeScalesSnapshot()["test-range"]
	if len(rangeSnap.Slots[0].Counters) != numCounters {
		t.Errorf("range scale snapshot slot has %d counters, want %d", len(rangeSnap.Slots[0].Counters), numCounters)
	}

	timeSnap := c.TimeScalesSnapshot()["test-time"]
	if len(timeSnap.Slots[0].Counters) != numCounters {
		t.Errorf("time scale snapshot slot has %d counters, want %d", len(timeSnap.Slots[0].Counters), numCounters)
	}
}

func TestMergeCollectors_Compatibility(t *testing.T) {
	// Compatible collectors
	c1, err := NewCollector(WithExemplar(WithExemplarCapacity(10)))
	if err != nil {
		t.Fatalf("Failed to create c1: %v", err)
	}
	c2, err := NewCollector(WithExemplar(WithExemplarCapacity(10)))
	if err != nil {
		t.Fatalf("Failed to create c2: %v", err)
	}
	_, err = MergeCollectors(c1, c2)
	if err != nil {
		t.Errorf("MergeCollectors with compatible configs failed: %v", err)
	}

	// Incompatible collectors
	c3, err := NewCollector(WithExemplar(WithExemplarCapacity(20)))
	if err != nil {
		t.Fatalf("Failed to create c3: %v", err)
	}
	_, err = MergeCollectors(c1, c3)
	if err == nil {
		t.Error("MergeCollectors with incompatible configs should have failed")
	}
}

func TestMergeSnapshots_Compatibility(t *testing.T) {
	s1 := &GlobalSnapshot{BoundsCRC32: 123, SketchKind: "tdigest"}
	s2 := &GlobalSnapshot{BoundsCRC32: 123, SketchKind: "tdigest"}
	_, err := MergeSnapshots(s1, s2)
	if err != nil {
		t.Errorf("MergeSnapshots with compatible snapshots failed: %v", err)
	}

	s3 := &GlobalSnapshot{BoundsCRC32: 456, SketchKind: "tdigest"}
	_, err = MergeSnapshots(s1, s3)
	if err == nil {
		t.Error("MergeSnapshots with incompatible BoundsCRC32 should have failed")
	}

	s4 := &GlobalSnapshot{BoundsCRC32: 123, SketchKind: "other"}
	_, err = MergeSnapshots(s1, s4)
	if err == nil {
		t.Error("MergeSnapshots with incompatible SketchKind should have failed")
	}
}

func TestCollector_Record(t *testing.T) {
	c, err := NewCollector()
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}
	rr := &RequestResult{
		Latency:   100 * time.Millisecond,
		QueueWait: 20 * time.Millisecond,
		Err:       errors.New("test error"),
	}
	c.Record(context.Background(), rr)

	snap := c.GlobalSnapshot()
	if snap.Total != 1 {
		t.Errorf("snap.Total = %d, want 1", snap.Total)
	}
	if snap.Errors != 1 {
		t.Errorf("snap.Errors = %d, want 1", snap.Errors)
	}
	if snap.Latencies.Total != 1 {
		t.Errorf("snap.Latencies.Total = %d, want 1", snap.Latencies.Total)
	}
	if snap.QueueWaits.Total != 1 {
		t.Errorf("snap.QueueWaits.Total = %d, want 1", snap.QueueWaits.Total)
	}
}

func TestJSONMarshaling(t *testing.T) {
	c, err := NewCollector(
		WithCustomCounters("c1"),
		WithRangeScale("rs", 1, 1),
	)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}
	c.Record(WithRequestID(context.Background(), 0), &RequestResult{
		Latency: 1,
		EndedAt: time.Unix(1, 0),
	})

	gs := c.GlobalSnapshot()
	if _, err := json.Marshal(gs); err != nil {
		t.Errorf("failed to marshal GlobalSnapshot: %v", err)
	}

	rs := c.RangeScalesSnapshot()
	if _, err := json.Marshal(rs); err != nil {
		t.Errorf("failed to marshal RangeScalesSnapshot: %v", err)
	}
}

func TestCollector_Merge(t *testing.T) {
	c1, err := NewCollector(
		WithCustomCounters("c1"),
		WithExemplar(WithExemplarCapacity(5)),
	)
	if err != nil {
		t.Fatalf("Failed to create c1: %v", err)
	}
	c2, err := NewCollector(
		WithCustomCounters("c1", "c2"),
		WithExemplar(WithExemplarCapacity(5)),
	)
	if err != nil {
		t.Fatalf("Failed to create c2: %v", err)
	}

	c1.Record(context.Background(), &RequestResult{Latency: 100})
	c1Handle, _ := c1.CounterHandle("c1")
	c1Handle.Inc()

	c2.Record(context.Background(), &RequestResult{Latency: 200, Err: errors.New("err")})
	c2Handle1, _ := c2.CounterHandle("c1")
	c2Handle1.Inc()
	c2Handle2, _ := c2.CounterHandle("c2")
	c2Handle2.Inc()

	if err := c1.Merge(c2); err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	snap := c1.GlobalSnapshot()
	if snap.Total != 2 {
		t.Errorf("snap.Total = %d, want 2", snap.Total)
	}
	if snap.Errors != 1 {
		t.Errorf("snap.Errors = %d, want 1", snap.Errors)
	}

	c1Handle, _ = c1.CounterHandle("c1")
	if c1Handle.value.Load() != 2 {
		t.Errorf("counter c1 = %d, want 2", c1Handle.value.Load())
	}
	c2Handle, _ := c1.CounterHandle("c2")
	if c2Handle.value.Load() != 1 {
		t.Errorf("counter c2 = %d, want 1", c2Handle.value.Load())
	}
}

func TestScale_Recording(t *testing.T) {
	c, err := NewCollector(
		WithRangeScale("range", 10, 10),
		WithTimeScale("time", 5, 10),
	)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Test RangeScale
	ctx := WithRequestID(context.Background(), 15)
	c.Record(ctx, &RequestResult{EndedAt: time.Now()})
	rangeSnap := c.RangeScalesSnapshot()["range"]
	if rangeSnap.Slots[1].Total != 1 {
		t.Errorf("range scale slot 1 should have 1 request, but has %d", rangeSnap.Slots[1].Total)
	}

	// Test TimeScale
	c2, err := NewCollector(
		WithTimeScale("time", 5, 10),
	)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}
	ts := time.Unix(23, 0)
	c2.Record(context.Background(), &RequestResult{EndedAt: ts})
	timeSnap := c2.TimeScalesSnapshot()["time"]
	if timeSnap.Slots[4].Total != 1 {
		t.Errorf("time scale slot 4 should have 1 request, but has %d", timeSnap.Slots[4].Total)
	}
}

func TestDummy(t *testing.T) {
	// This is a dummy test to ensure that the test file is always valid,
	// even if other tests are commented out.
}
