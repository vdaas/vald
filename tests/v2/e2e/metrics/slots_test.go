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
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/net/grpc/codes"
)

func TestSlot_Record(t *testing.T) {
	t.Parallel()

	// Create a new slot with mock histograms/exemplars if possible, or real ones.
	// For simplicity, we'll use nil since Record handles nil checks internally for Latency/QueueWait/Exemplars.
	// But to verify logic we might want real ones.
	h, _ := NewHistogram(WithHistogramNumBuckets(10))
	e := NewExemplar(WithExemplarCapacity(10))
	s := newSlot(1, h, h, e).(*slot)

	// 1. Record a simple request
	t0 := time.Now()
	rr := &RequestResult{
		EndedAt: t0,
		Latency: 100 * time.Millisecond,
		Err:     nil,
	}
	s.Record(rr, 0)

	if s.Total.Load() != 1 {
		t.Errorf("expected Total 1, got %d", s.Total.Load())
	}
	if s.WindowStart != 0 {
		t.Errorf("expected WindowStart 0, got %d", s.WindowStart)
	}

	// 2. Record with same window index
	s.Record(rr, 0)
	if s.Total.Load() != 2 {
		t.Errorf("expected Total 2, got %d", s.Total.Load())
	}

	// 3. Record with different window index (should reset)
	s.Record(rr, 1)
	if s.Total.Load() != 1 {
		t.Errorf("expected Total 1 (reset), got %d", s.Total.Load())
	}
	if s.WindowStart != 1 {
		t.Errorf("expected WindowStart 1, got %d", s.WindowStart)
	}
}

func TestSlot_Merge(t *testing.T) {
	t.Parallel()

	s1 := newSlot(1, nil, nil, nil).(*slot)
	s2 := newSlot(1, nil, nil, nil).(*slot)

	s1.Total.Store(10)
	s1.Errors.Store(1)
	s1.updatedNS.Store(100)
	s1.Counters[0].Store(5)

	s2.Total.Store(20)
	s2.Errors.Store(2)
	s2.updatedNS.Store(200)
	s2.Counters[0].Store(10) // Slot merge currently doesn't merge counters? Let's check implementation.
	// Looking at implementation: Merge function does NOT iterate counters slice!
	// It merges Total, Errors, updatedNS, Latency, QueueWait, Exemplars.
	// It seems Counters are NOT merged in Slot.Merge?
	// Let's verify code.

	if err := s1.Merge(s2); err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	if s1.Total.Load() != 30 {
		t.Errorf("expected Total 30, got %d", s1.Total.Load())
	}
	if s1.Errors.Load() != 3 {
		t.Errorf("expected Errors 3, got %d", s1.Errors.Load())
	}
	if s1.updatedNS.Load() != 200 {
		t.Errorf("expected updatedNS 200, got %d", s1.updatedNS.Load())
	}
}

func TestSlot_Clone(t *testing.T) {
	t.Parallel()

	s1 := newSlot(1, nil, nil, nil).(*slot)
	s1.Total.Store(10)
	s1.Counters[0].Store(5)

	s2 := s1.Clone().(*slot)

	if s2.Total.Load() != 10 {
		t.Errorf("expected cloned Total 10, got %d", s2.Total.Load())
	}
	if s2.Counters[0].Load() != 5 {
		t.Errorf("expected cloned Counter 5, got %d", s2.Counters[0].Load())
	}

	// Modify s1, s2 should not change
	s1.Total.Store(20)
	if s2.Total.Load() != 10 {
		t.Errorf("expected independent cloned Total 10, got %d", s2.Total.Load())
	}
}

func TestSlot_Reset(t *testing.T) {
	t.Parallel()
	s := newSlot(1, nil, nil, nil).(*slot)
	s.Total.Store(10)

	s.Reset()

	if s.Total.Load() != 0 {
		t.Errorf("expected Total 0 after Reset, got %d", s.Total.Load())
	}
}

func TestSlot_Concurrent_Record_Reset(t *testing.T) {
	t.Parallel()
	// Verify that concurrent Records and Reset (via window change) don't panic or race.
	s := newSlot(0, nil, nil, nil).(*slot)

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// Switch window index every 100 iterations
				win := uint64((id*1000 + j) / 100)
				s.Record(&RequestResult{
					EndedAt: start,
					Status:  codes.OK,
				}, win)
			}
		}(i)
	}
	wg.Wait()
}
