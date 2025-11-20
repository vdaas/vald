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
	"sync"
	"testing"
	"time"
)

func TestScale_RingBuffer_Reset(t *testing.T) {
	t.Parallel()
	// Create a time-based scale with width 1 (second) and capacity 2.
	// This means slots 0 and 1.
	// Time 0 -> Slot 0
	// Time 1 -> Slot 1
	// Time 2 -> Slot 0 (Should reset)

	s, err := newScale("test", 1, 2, 0, TimeScale, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create scale: %v", err)
	}

	ctx := context.Background()

	// Record at Time 0
	t0 := time.Unix(0, 0)
	s.Record(ctx, &RequestResult{
		EndedAt: t0,
	})

	// Verify Slot 0 has 1 request
	snap := s.Snapshot()
	slot0 := snap.Slots[0]
	if slot0.Total != 1 {
		t.Errorf("expected slot 0 total 1, got %d", slot0.Total)
	}

	// Record at Time 1
	t1 := time.Unix(1, 0)
	s.Record(ctx, &RequestResult{
		EndedAt: t1,
	})

	// Verify Slot 1 has 1 request
	snap = s.Snapshot()
	slot1 := snap.Slots[1]
	if slot1.Total != 1 {
		t.Errorf("expected slot 1 total 1, got %d", slot1.Total)
	}

	// Record at Time 2 (Should wrap to Slot 0 and reset)
	t2 := time.Unix(2, 0)
	s.Record(ctx, &RequestResult{
		EndedAt: t2,
	})

	// Verify Slot 0 has 1 request (reset happened)
	snap = s.Snapshot()
	slot0 = snap.Slots[0]
	if slot0.Total != 1 {
		t.Errorf("expected slot 0 total 1 (after reset), got %d", slot0.Total)
	}
	// Verify LastUpdated matches t2
	if slot0.LastUpdated != t2.UnixNano() {
		t.Errorf("expected slot 0 updated %d, got %d", t2.UnixNano(), slot0.LastUpdated)
	}
}

func TestScale_Concurrency_Reset(t *testing.T) {
	t.Parallel()
	// High concurrency test to ensure race conditions in reset are handled.
	s, err := newScale("test_concurrent", 1, 5, 0, TimeScale, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create scale: %v", err)
	}

	ctx := context.Background()
	var wg sync.WaitGroup

	start := time.Now()

	// Spawn goroutines writing to consecutive seconds.
	// This forces rapid wrapping and resetting.
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				// Use fake times that increment to force wrapping
				now := start.Add(time.Duration(offset+j) * time.Second)
				s.Record(ctx, &RequestResult{
					EndedAt: now,
				})
			}
		}(i)
	}

	wg.Wait()
	// If no panic or race detected, pass.
}

func TestScale_Merge(t *testing.T) {
	t.Parallel()
	s1, _ := newScale("test_merge", 1, 2, 0, TimeScale, nil, nil, nil)
	s2, _ := newScale("test_merge", 1, 2, 0, TimeScale, nil, nil, nil)

	// Manually record data since we can't access slots directly
	ctx := context.Background()
	t0 := time.Unix(0, 0)
	t1 := time.Unix(1, 0)

	s1.Record(ctx, &RequestResult{EndedAt: t0})
	s2.Record(ctx, &RequestResult{EndedAt: t0})
	s2.Record(ctx, &RequestResult{EndedAt: t1})

	if err := s1.Merge(s2); err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	snap := s1.Snapshot()
	if snap.Slots[0].Total != 2 { // 1 from s1 + 1 from s2
		t.Errorf("expected slot 0 total 2, got %d", snap.Slots[0].Total)
	}
	if snap.Slots[1].Total != 1 { // 0 from s1 + 1 from s2
		t.Errorf("expected slot 1 total 1, got %d", snap.Slots[1].Total)
	}
}

func TestScale_Clone(t *testing.T) {
	t.Parallel()
	s1, _ := newScale("test_clone", 1, 2, 0, TimeScale, nil, nil, nil)
	ctx := context.Background()
	s1.Record(ctx, &RequestResult{EndedAt: time.Unix(0, 0)})

	s2 := s1.Clone()

	snap1 := s1.Snapshot()
	snap2 := s2.Snapshot()

	if snap2.Name != snap1.Name {
		t.Errorf("expected cloned name %s, got %s", snap1.Name, snap2.Name)
	}
	if snap2.Slots[0].Total != snap1.Slots[0].Total {
		t.Errorf("expected cloned slot 0 total %d, got %d", snap1.Slots[0].Total, snap2.Slots[0].Total)
	}

	// Verify independence
	s1.Record(ctx, &RequestResult{EndedAt: time.Unix(0, 0)})
	snap2 = s2.Snapshot()
	if snap2.Slots[0].Total != 1 {
		t.Errorf("expected cloned slot 0 total to remain 1, got %d", snap2.Slots[0].Total)
	}
}
