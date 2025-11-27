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
	"sync/atomic"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// slot holds the metrics for a single window in a scale.
// It is an implementation of the Slot interface.
type slot struct {
	mu          sync.RWMutex    // protects WindowStart and coordinates Reset
	WindowStart uint64          // The window index (idx / width) this slot represents
	Total       atomic.Uint64   // total number of requests in this slot
	Errors      atomic.Uint64   // number of errored requests in this slot
	updatedNS   atomic.Int64    // UnixNano timestamp of the last update
	Latency     Histogram       // latency histogram
	QueueWait   Histogram       // queue wait histogram
	Counters    []atomic.Uint64 // custom counters
	Exemplars   Exemplar        // exemplar heap
}

// newSlot creates a new Slot with the given configuration.
// It initializes a new slot with histograms, exemplars, and custom counters.
func newSlot(numCounters int, latencies, queueWaits Histogram, exemplars Exemplar) Slot {
	return &slot{
		Latency:   latencies,
		QueueWait: queueWaits,
		Counters:  make([]atomic.Uint64, numCounters),
		Exemplars: exemplars,
	}
}

// Reset resets the slot data to its initial state.
// It acquires the write lock to ensure exclusive access during reset.
func (s *slot) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reset()
}

// reset clears the slot data.
// It assumes that the write lock (s.mu) is already held by the caller.
func (s *slot) reset() {
	s.Total.Store(0)
	s.Errors.Store(0)
	s.updatedNS.Store(0)
	if s.Latency != nil {
		s.Latency.Reset()
	}
	if s.QueueWait != nil {
		s.QueueWait.Reset()
	}
	if s.Exemplars != nil {
		s.Exemplars.Reset()
	}
	for i := range s.Counters {
		s.Counters[i].Store(0)
	}
}

// Record processes a single RequestResult for the slot.
// It implements a ring buffer logic: if the given windowIdx is different from the slot's current WindowStart,
// it resets the slot to start a new window. This check is done using double-checked locking for performance and safety.
func (s *slot) Record(rr *RequestResult, windowIdx uint64) {
	// Optimistic read lock to check if the slot is valid for the current window.
	s.mu.RLock()

	// Ignore late arrivals for older windows to preserve data integrity.
	if windowIdx < s.WindowStart {
		s.mu.RUnlock()
		return
	}

	if s.WindowStart != windowIdx {
		s.mu.RUnlock()
		// Slot is stale or uninitialized. Upgrade to write lock to reset.
		s.mu.Lock()
		// Double-check under write lock.
		if windowIdx < s.WindowStart {
			s.mu.Unlock()
			return
		}
		if s.WindowStart != windowIdx {
			s.reset()
			s.WindowStart = windowIdx
		}
		s.mu.Unlock()
		// Re-acquire read lock to proceed with recording.
		s.mu.RLock()
	}

	// Verify again that the window matches (it could have changed if we lost the race
	// after unlocking and before re-locking, though unlikely for typical usage).
	if s.WindowStart != windowIdx {
		s.mu.RUnlock()
		return // Slot changed again, ignore this metric (or retry)
	}

	defer s.mu.RUnlock()

	s.Total.Add(1)
	if rr.Err != nil {
		s.Errors.Add(1)
	}
	s.updatedNS.Store(rr.EndedAt.UnixNano())
	if s.Latency != nil {
		s.Latency.Record(float64(rr.Latency.Nanoseconds()))
	}
	if s.QueueWait != nil {
		s.QueueWait.Record(float64(rr.QueueWait.Nanoseconds()))
	}
	if s.Exemplars != nil {
		s.Exemplars.Offer(rr.Latency, rr.RequestID, rr.Err, rr.Msg)
	}
}

// Clone returns a deep copy of the slot.
// It creates a new independent slot with copied data.
func (s *slot) Clone() Slot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counters := make([]atomic.Uint64, len(s.Counters))
	for k := range s.Counters {
		counters[k].Store(s.Counters[k].Load())
	}
	var l, q Histogram
	var e Exemplar
	if s.Latency != nil {
		l = s.Latency.Clone()
	}
	if s.QueueWait != nil {
		q = s.QueueWait.Clone()
	}
	if s.Exemplars != nil {
		e = s.Exemplars.Clone()
	}

	newS := &slot{
		Latency:     l,
		QueueWait:   q,
		Counters:    counters,
		Exemplars:   e,
		WindowStart: s.WindowStart,
	}
	newS.Total.Store(s.Total.Load())
	newS.Errors.Store(s.Errors.Load())
	newS.updatedNS.Store(s.updatedNS.Load())

	return newS
}

// Merge merges another slot into this one.
// It aggregates metrics from the other slot into the current slot.
// It uses ordered locking based on pointer addresses to prevent deadlocks.
func (s *slot) Merge(other Slot) error {
	if s == other {
		return nil
	}

	os, ok := other.(*slot)
	if !ok {
		return errors.New("incompatible slot implementation")
	}

	// To prevent deadlocks, always lock in a consistent order.
	if uintptr(unsafe.Pointer(s)) < uintptr(unsafe.Pointer(os)) {
		s.mu.Lock()
		os.mu.Lock()
	} else {
		os.mu.Lock()
		s.mu.Lock()
	}
	defer s.mu.Unlock()
	defer os.mu.Unlock()

	s.Total.Add(os.Total.Load())
	s.Errors.Add(os.Errors.Load())
	if os.updatedNS.Load() > s.updatedNS.Load() {
		s.updatedNS.Store(os.updatedNS.Load())
	}
	if s.Latency != nil && os.Latency != nil {
		if err := s.Latency.Merge(os.Latency); err != nil {
			return err
		}
	}
	if s.QueueWait != nil && os.QueueWait != nil {
		if err := s.QueueWait.Merge(os.QueueWait); err != nil {
			return err
		}
	}
	if s.Exemplars != nil && os.Exemplars != nil {
		if err := s.Exemplars.Merge(os.Exemplars); err != nil {
			return err
		}
	}
	return nil
}

// Snapshot returns a snapshot of the slot's current state.
// The returned snapshot is a static view and safe for concurrent access.
func (s *slot) Snapshot() *SlotSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counters := make([]uint64, len(s.Counters))
	for j := range counters {
		counters[j] = s.Counters[j].Load()
	}
	var latSnap, qwSnap *HistogramSnapshot
	var exSnap []*ExemplarItem
	if s.Latency != nil {
		latSnap = s.Latency.Snapshot()
	}
	if s.QueueWait != nil {
		qwSnap = s.QueueWait.Snapshot()
	}
	if s.Exemplars != nil {
		exSnap = s.Exemplars.Snapshot()
	}
	return &SlotSnapshot{
		Total:       s.Total.Load(),
		Errors:      s.Errors.Load(),
		LastUpdated: s.updatedNS.Load(),
		Latencies:   latSnap,
		QueueWaits:  qwSnap,
		Counters:    counters,
		Exemplars:   exSnap,
	}
}
