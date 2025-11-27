package metrics

import (
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// slot holds the metrics for a single window in a scale.
// It is an implementation of the Slot interface.
type slot struct {
	mu          sync.RWMutex    // protects WindowStart and coordinates Reset
	id          uint64          // Unique ID for lock ordering
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
func newSlot(numCounters int, latencies, queueWaits Histogram, exemplars Exemplar) Slot {
	return &slot{
		id:        collectorIDCounter.Add(1),
		Latency:   latencies,
		QueueWait: queueWaits,
		Counters:  make([]atomic.Uint64, numCounters),
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
		s.Counters[i].Store(0)
	}
}

// Record processes a single RequestResult for the slot.
// Refactored to optimize locking:
// 1. Fast Path: Acquires Read Lock. If window matches, records and returns.
// 2. Slow Path: Acquires Write Lock. Resets window if needed, records immediately, and returns.
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

	// Optimization: Record immediately under the Write Lock.
	// This prevents the data loss race condition where unlocking and re-acquiring
	// a Read Lock would allow another thread to reset the window again before we record.
	s.recordInternal(rr)
}

// recordInternal updates the metrics in the slot.
// It must be called while holding s.mu (Read or Write).
func (s *slot) recordInternal(rr *RequestResult) {
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
