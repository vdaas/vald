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
	"testing"
	"time"
)

// TestScale_RingBuffer_WrapAndReset verifies that the ring buffer correctly resets slots
// when the window index wraps around and overrides an old slot.
func TestScale_RingBuffer_WrapAndReset(t *testing.T) {
	t.Parallel()

	// Configuration:
	// Width = 1 second
	// Capacity = 3 slots
	// Ring buffer indices: 0, 1, 2
	width := uint64(1)
	capacity := uint64(3)
	s, err := newScale("test_wrap", width, capacity, 0, TimeScale, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create scale: %v", err)
	}

	ctx := context.Background()

	// Helper to record a hit at a specific timestamp (seconds)
	recordAt := func(sec int64) {
		s.Record(ctx, &RequestResult{
			EndedAt: time.Unix(sec, 0),
		})
	}

	// 1. Fill the buffer: Seconds 0, 1, 2 -> Slots 0, 1, 2
	recordAt(0) // Slot 0
	recordAt(1) // Slot 1
	recordAt(2) // Slot 2

	snap := s.Snapshot()
	if snap.Slots[0].Total != 1 {
		t.Errorf("Slot 0: expected 1, got %d", snap.Slots[0].Total)
	}
	if snap.Slots[1].Total != 1 {
		t.Errorf("Slot 1: expected 1, got %d", snap.Slots[1].Total)
	}
	if snap.Slots[2].Total != 1 {
		t.Errorf("Slot 2: expected 1, got %d", snap.Slots[2].Total)
	}

	// 2. Wrap around: Second 3 -> Slot 0 (3 % 3 == 0)
	// This should RESET Slot 0 (which had data from Second 0)
	recordAt(3)

	snap = s.Snapshot()
	// Slot 0 should now have 1 hit (from Second 3), NOT 2 hits (0 + 1)
	if snap.Slots[0].Total != 1 {
		t.Errorf("Slot 0 (after wrap): expected 1, got %d. Indicates failure to reset.", snap.Slots[0].Total)
	}
	// Check timestamp to confirm it's from Second 3
	if snap.Slots[0].LastUpdated != time.Unix(3, 0).UnixNano() {
		t.Errorf("Slot 0 timestamp: expected %d, got %d", time.Unix(3, 0).UnixNano(), snap.Slots[0].LastUpdated)
	}

	// Slots 1 and 2 should remain untouched (containing data from Seconds 1 and 2)
	if snap.Slots[1].Total != 1 {
		t.Errorf("Slot 1: expected 1, got %d", snap.Slots[1].Total)
	}
	if snap.Slots[2].Total != 1 {
		t.Errorf("Slot 2: expected 1, got %d", snap.Slots[2].Total)
	}

	// 3. Wrap again: Second 5 -> Slot 2 (5 % 3 == 2)
	// Skip Second 4 (Slot 1) to show we don't need contiguous updates
	recordAt(5)

	snap = s.Snapshot()
	if snap.Slots[2].Total != 1 {
		t.Errorf("Slot 2 (after wrap): expected 1, got %d", snap.Slots[2].Total)
	}
	if snap.Slots[2].LastUpdated != time.Unix(5, 0).UnixNano() {
		t.Errorf("Slot 2 timestamp: expected %d, got %d", time.Unix(5, 0).UnixNano(), snap.Slots[2].LastUpdated)
	}
}
