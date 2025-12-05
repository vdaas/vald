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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

// slot holds the metrics for a single window in a scale.
// It is an implementation of the Slot interface.
// paddedCounter is a padded atomic counter to prevent false sharing.
type paddedCounter struct {
	val atomic.Uint64
	_   [paddingSize]byte
}

type slot struct {
	Latency     Histogram
	QueueWait   Histogram
	Exemplars   Exemplar
	Counters    []paddedCounter
	id          uint64
	WindowStart uint64
	Total       atomic.Uint64
	Errors      atomic.Uint64
	updatedNS   atomic.Int64
	// padding to prevent false sharing
	_  [paddingSize]byte
	mu sync.RWMutex
}

// newSlot creates a new Slot with the given configuration.
func newSlot(numCounters int, latencies, queueWaits Histogram, exemplars Exemplar) Slot {
	return &slot{
		id:        collectorIDCounter.Add(1),
		Latency:   latencies,
		QueueWait: queueWaits,
		Counters:  make([]paddedCounter, numCounters),
		Exemplars: exemplars,
	}
}

// Reset resets the slot data to its initial state.
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
		s.Counters[i].val.Store(0)
	}
}

// Record processes a single RequestResult for the slot.
// Uses RLock for high concurrency in the hot path, upgrading to Lock only for window rotation.
func (s *slot) Record(rr *RequestResult, windowIdx uint64) {
	if rr == nil {
		return
	}

	// --- 1. Fast Path (Optimistic Read) ---
	s.mu.RLock()
	if s.WindowStart == windowIdx {
		s.recordInternal(rr)
		s.mu.RUnlock()
		return
	}
	s.mu.RUnlock()

	// --- 2. Slow Path (Write Lock for Transition) ---
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check: Another goroutine might have reset the window while we waited for the lock.
	if windowIdx < s.WindowStart {
		// Data is for an old window that has already passed. Drop it to preserve integrity.
		return
	}

	// If the window is still mismatched (implied s.WindowStart < windowIdx), reset to the new window.
	if s.WindowStart != windowIdx {
		s.reset()
		s.WindowStart = windowIdx
	}

	s.recordInternal(rr)
}

// recordInternal updates the metrics in the slot.
// It must be called while holding s.mu (Read or Write).
func (s *slot) recordInternal(rr *RequestResult) {
	s.Total.Add(1)
	if rr.Err != nil {
		s.Errors.Add(1)
	}

	// Update updatedNS (max) with CAS loop for correctness,
	// because multiple readers (fast path) can execute this concurrently.
	if e := rr.EndedAt.UnixNano(); e > 0 {
		for {
			curr := s.updatedNS.Load()
			if e <= curr {
				break
			}
			if s.updatedNS.CompareAndSwap(curr, e) {
				break
			}
		}
	}

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
func (s *slot) Clone() Slot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counters := make([]paddedCounter, len(s.Counters))
	for k := range s.Counters {
		counters[k].val.Store(s.Counters[k].val.Load())
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
		id:          collectorIDCounter.Add(1),
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
func (s *slot) Merge(other Slot) error {
	if s == other {
		return nil
	}

	os, ok := other.(*slot)
	if !ok {
		return errors.New("incompatible slot implementation")
	}

	// To prevent deadlocks, always lock in a consistent order.
	if s.id < os.id {
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

	// Merge updatedNS using CAS loop or simple check-and-set if locked
	// But since we hold lock, check-and-set is safe against other writers to `s`.
	// But `os` is also locked.
	if t := os.updatedNS.Load(); t > s.updatedNS.Load() {
		s.updatedNS.Store(t)
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
func (s *slot) Snapshot() *SlotSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counters := make([]uint64, len(s.Counters))
	for j := range counters {
		counters[j] = s.Counters[j].val.Load()
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
