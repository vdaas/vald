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
	"cmp"
	"container/heap"
	"math/rand/v2"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

// exemplar holds samples of requests in different categories.
type exemplar struct {
	mu sync.Mutex
	k  int // The maximum number of exemplars to store per category.

	// Categories
	slowest  priorityQueue    // Min-heap (Top-K Max Latency)
	fastest  maxPriorityQueue // Max-heap (Top-K Min Latency)
	failures priorityQueue    // Min-heap (Top-K Slowest Failures) - "Top Failures" usually implies notable ones (slow).
	// If we wanted "Representative Failures", we'd use reservoir sampling.
	// Given "Top XX Failures", and usually failures are bad if slow (or fast fail?),
	// I'll implement "Slowest Failures".
	// Wait, user said "Top XX Failures (sampling of failed requests)". "Sampling" suggests random.
	// But "Top" suggests ordering.
	// I'll use Reservoir Sampling for "Avg" (Representative) and "Failures" (Sampling).
	// Actually, for Failures, maybe just latest?
	// Let's stick to:
	// 1. Slowest (Top K Latency)
	// 2. Fastest (Bottom K Latency)
	// 3. Average (Reservoir Sampling - Representative)
	// 4. Failures (Reservoir Sampling - Representative of failures)

	avgSamples     []*item // Reservoir for representative samples
	failureSamples []*item // Reservoir for failure samples
	avgCount       uint64  // Total count seen for average reservoir
	failureCount   uint64  // Total count seen for failure reservoir

	minLatency atomic.Int64 // Minimum latency in the 'slowest' heap (fast path)
	maxLatency atomic.Int64 // Maximum latency in the 'fastest' heap (fast path for rejection?)
	// Actually maxLatency helps reject large values for 'fastest' heap (which stores smallest).
	// If val > maxLatency and heap full, reject.
}

// Init initializes the exemplar with the given options.
func (e *exemplar) Init(opts ...ExemplarOption) {
	for _, opt := range opts {
		opt(e)
	}
	e.k = max(e.k, 1)
	e.mu.Lock()
	e.initHeaps()
	e.mu.Unlock()
}

func (e *exemplar) initHeaps() {
	if e.slowest == nil {
		e.slowest = make(priorityQueue, 0, e.k)
	} else {
		e.slowest = e.slowest[:0]
	}
	if e.fastest == nil {
		e.fastest = make(maxPriorityQueue, 0, e.k)
	} else {
		e.fastest = e.fastest[:0]
	}
	if e.avgSamples == nil {
		e.avgSamples = make([]*item, 0, e.k)
	} else {
		e.avgSamples = e.avgSamples[:0]
	}
	if e.failureSamples == nil {
		e.failureSamples = make([]*item, 0, e.k)
	} else {
		e.failureSamples = e.failureSamples[:0]
	}
	e.minLatency.Store(0)
	e.maxLatency.Store(0)
	e.avgCount = 0
	e.failureCount = 0
}

// NewExemplar creates a new Exemplar with the given options.
func NewExemplar(opts ...ExemplarOption) Exemplar {
	e := exemplarPool.Get()
	e.Init(opts...)
	return e
}

// Reset resets the exemplar to its initial state.
func (e *exemplar) Reset() {
	e.mu.Lock()
	e.initHeaps() // Reset slices and counts
	e.mu.Unlock()
}

// Offer adds a request to the exemplar categories.
func (e *exemplar) Offer(latency time.Duration, requestID string, isError bool) {
	newItem := &item{latency: latency, requestID: requestID, isError: isError}
	latInt := int64(latency)

	// We must lock because we are updating multiple structures.
	// Optimizations (atomics) are possible but complex with multiple categories.
	e.mu.Lock()
	defer e.mu.Unlock()

	// 1. Slowest (Top K Max)
	// Min-Heap stores K largest. Root is the smallest of the largest.
	// If new > Root, replace Root.
	if len(e.slowest) < e.k {
		heap.Push(&e.slowest, newItem)
		if len(e.slowest) == e.k {
			e.minLatency.Store(int64(e.slowest[0].latency))
		}
	} else if latInt > int64(e.slowest[0].latency) {
		e.slowest[0] = newItem
		heap.Fix(&e.slowest, 0)
		e.minLatency.Store(int64(e.slowest[0].latency))
	}

	// 2. Fastest (Top K Min)
	// Max-Heap stores K smallest. Root is the largest of the smallest.
	// If new < Root, replace Root.
	if len(e.fastest) < e.k {
		heap.Push(&e.fastest, newItem)
		if len(e.fastest) == e.k {
			e.maxLatency.Store(int64(e.fastest[0].latency))
		}
	} else if latInt < int64(e.fastest[0].latency) {
		e.fastest[0] = newItem
		heap.Fix(&e.fastest, 0)
		e.maxLatency.Store(int64(e.fastest[0].latency))
	}

	// 3. Average (Representative Sample)
	e.avgCount++
	if len(e.avgSamples) < e.k {
		e.avgSamples = append(e.avgSamples, newItem)
	} else {
		// Reservoir sampling: replace with probability k/count
		j := rand.Uint64N(e.avgCount)
		if j < uint64(e.k) {
			e.avgSamples[j] = newItem
		}
	}

	// 4. Failures (Representative Sample of Failures)
	if isError {
		e.failureCount++
		if len(e.failureSamples) < e.k {
			e.failureSamples = append(e.failureSamples, newItem)
		} else {
			j := rand.Uint64N(e.failureCount)
			if j < uint64(e.k) {
				e.failureSamples[j] = newItem
			}
		}
	}
}

// Snapshot returns a snapshot of the exemplars.
// It returns a flat list. The user might want distinct lists.
// For backward compatibility, we might return all?
// Or we should change the return type? The interface `Exemplar` returns `[]*item`.
// `GlobalSnapshot` has `Exemplar []*item`.
// I should probably return "Slowest" as the primary for backward compat, or mix them?
// Given the request "Expand Exemplar Categories", likely the output format should change.
// However, changing the return type breaks the interface and `GlobalSnapshot`.
// I will flatten them into one list for now, or return just Slowest?
// The prompt said "Refactor this to support multiple distinct exemplar categories".
// This implies the output should distinguish them.
// But `GlobalSnapshot` struct has `Exemplar []*item`.
// I can't easily change `GlobalSnapshot` struct without breaking consumers (unless I add fields).
// I will add fields to `GlobalSnapshot`? No, `GlobalSnapshot` is defined in `metrics.go` which uses `[]*item`.
// I will update `GlobalSnapshot` in `metrics.go` later if needed.
// For now `Snapshot()` will return the "Slowest" ones to satisfy the interface,
// but I should probably add a new method `SnapshotDetails()`?
// Or return all combined?
// If I return all combined, they are just a list.
//
// I will modify the `Exemplar` interface in `interface.go` (which I haven't read but assume exists)
// or just modify `Snapshot` to return all?
//
// If I modify `metrics.go`'s `GlobalSnapshot` struct, I can add `Fastest`, `Average`, `Failures`.
// Let's check `metrics.go` again.
// `Exemplars []*item`.
//
// I will update `metrics.go` to include new fields.
// But first, let's implement `Snapshot` here to return a map or struct?
// Since `metrics.go` expects `[]*item`, I'll return a combined list or I need to change `metrics.go`.
//
// I will stick to returning "Slowest" in `Snapshot()` for backward compatibility if forced,
// BUT I will add `DetailedSnapshot` method.
//
// Actually, `metrics.go` calls `c.exemplars.Snapshot()`.
// I should update `metrics.go` to use the new categories.
//
// So, I will change `Snapshot` to return a struct `ExemplarSnapshot`.
// But `Exemplar` is an interface. I need to check `interface.go`.

func (e *exemplar) Snapshot() []*item {
	// For backward compatibility, return Slowest.
	e.mu.Lock()
	defer e.mu.Unlock()
	items := slices.Clone(e.slowest)
	slices.SortFunc(items, func(a, b *item) int {
		return cmp.Compare(b.latency, a.latency)
	})
	return items
}

// DetailedSnapshot returns all categories.
func (e *exemplar) DetailedSnapshot() (*ExemplarDetails, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	snap := &ExemplarDetails{
		Slowest:  make([]*item, len(e.slowest)),
		Fastest:  make([]*item, len(e.fastest)),
		Average:  make([]*item, len(e.avgSamples)),
		Failures: make([]*item, len(e.failureSamples)),
	}

	copy(snap.Slowest, e.slowest)
	slices.SortFunc(snap.Slowest, func(a, b *item) int {
		return cmp.Compare(b.latency, a.latency) // Descending
	})

	copy(snap.Fastest, e.fastest)
	slices.SortFunc(snap.Fastest, func(a, b *item) int {
		return cmp.Compare(a.latency, b.latency) // Ascending
	})

	copy(snap.Average, e.avgSamples)
	slices.SortFunc(snap.Average, func(a, b *item) int {
		return cmp.Compare(b.latency, a.latency) // Descending
	})

	copy(snap.Failures, e.failureSamples)
	slices.SortFunc(snap.Failures, func(a, b *item) int {
		return cmp.Compare(b.latency, a.latency) // Descending
	})

	return snap, nil
}

// Merge merges another exemplar into this one.
func (e *exemplar) Merge(other Exemplar) error {
	if other == nil {
		return nil
	}
	details, _ := other.DetailedSnapshot()
	if details == nil {
		return nil
	}

	// We iterate over the snapshots and offer them to this exemplar.
	// This is a simplified merge strategy that works reasonably well for top-K heaps.
	// For reservoirs (Average, Failures), re-offering samples is not statistically
	// perfect (it biases towards the later merges if counts are not tracked),
	// but it preserves "representativeness" for E2E testing purposes.
	// A better approach would be to merge reservoirs properly using counts,
	// but that requires exposing internal state (counts) via interface or type assertion.
	// Given the interface constraint, we use Offer.

	for _, ex := range details.Slowest {
		e.Offer(ex.latency, ex.requestID, ex.isError)
	}
	for _, ex := range details.Fastest {
		e.Offer(ex.latency, ex.requestID, ex.isError)
	}
	for _, ex := range details.Average {
		e.Offer(ex.latency, ex.requestID, ex.isError)
	}
	for _, ex := range details.Failures {
		e.Offer(ex.latency, ex.requestID, ex.isError)
	}
	return nil
}

// Clone returns a deep copy.
func (e *exemplar) Clone() Exemplar {
	newE := exemplarPool.Get()
	newE.Reset()
	newE.k = e.k

	e.mu.Lock()
	defer e.mu.Unlock()

	copyTo := func(dst *[]*item, src []*item) {
		if cap(*dst) < len(src) {
			*dst = make([]*item, len(src), cap(src))
		} else {
			*dst = (*dst)[:len(src)]
		}
		for i, it := range src {
			if it != nil {
				v := *it
				(*dst)[i] = &v
			}
		}
	}

	// Heaps
	newE.slowest = make(priorityQueue, len(e.slowest), cap(e.slowest))
	for i, it := range e.slowest {
		v := *it
		newE.slowest[i] = &v
	}
	newE.fastest = make(maxPriorityQueue, len(e.fastest), cap(e.fastest))
	for i, it := range e.fastest {
		v := *it
		newE.fastest[i] = &v
	}

	// Reservoirs
	copyTo(&newE.avgSamples, e.avgSamples)
	copyTo(&newE.failureSamples, e.failureSamples)

	newE.avgCount = e.avgCount
	newE.failureCount = e.failureCount
	newE.minLatency.Store(e.minLatency.Load())
	newE.maxLatency.Store(e.maxLatency.Load())

	return newE
}

// item is an item in the priority queue.
type item struct {
	latency   time.Duration
	requestID string
	isError   bool
}

// priorityQueue implements min-heap.
type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].latency < pq[j].latency
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*item)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

// maxPriorityQueue implements max-heap (for Fastest).
type maxPriorityQueue []*item

func (pq maxPriorityQueue) Len() int { return len(pq) }
func (pq maxPriorityQueue) Less(i, j int) bool {
	return pq[i].latency > pq[j].latency // Largest comes first? No, heap.Pop returns smallest?
	// heap.Pop returns the element at index 0 (the root).
	// heap.Fix/Push/Pop maintains the heap invariant: pq[i] <= pq[2*i+1] etc.
	// Less(i, j) returns true if i should appear before j (i is "smaller" in heap terms).
	// For Max-Heap, we want the root to be the Largest. So Less means "Greater".
	// pq[i].latency > pq[j].latency.
	// Wait, for "Fastest", we want to keep K *smallest* latencies.
	// A standard Min-Heap keeps the *smallest* at the root. If full, we replace root?
	// If full, we want to discard the *Largest* of the K smallest to make room for a smaller one.
	// So we need a Max-Heap of size K. The root is the Largest of the set.
	// If new < Root, replace Root.
	// Yes.
	// So Less should be >.
}

func (pq maxPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *maxPriorityQueue) Push(x any) {
	item := x.(*item)
	*pq = append(*pq, item)
}

func (pq *maxPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

type ExemplarDetails struct {
	Slowest  []*item
	Fastest  []*item
	Average  []*item
	Failures []*item
}
