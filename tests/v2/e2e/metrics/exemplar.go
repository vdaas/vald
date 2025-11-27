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
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/sync"
)

// exemplar holds samples of requests in different categories.
type exemplar struct {
	mu sync.Mutex
	k  int // The maximum number of exemplars to store per category.

	// Categories
	slowest  priorityQueue       // Min-heap (Top-K Max Latency)
	fastest  smallestLatencyHeap // Max-heap (Top-K Min Latency) to evict largest, keeping smallest.
	failures priorityQueue       // Min-heap (Top-K Slowest Failures) - "Top Failures" usually implies notable ones (slow).
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

	avgSamples     []*ExemplarItem // Reservoir for representative samples
	failureSamples []*ExemplarItem // Reservoir for failure samples
	avgCount       uint64          // Total count seen for average reservoir
	failureCount   uint64          // Total count seen for failure reservoir

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
		e.fastest = make(smallestLatencyHeap, 0, e.k)
	} else {
		e.fastest = e.fastest[:0]
	}
	if e.avgSamples == nil {
		e.avgSamples = make([]*ExemplarItem, 0, e.k)
	} else {
		e.avgSamples = e.avgSamples[:0]
	}
	if e.failureSamples == nil {
		e.failureSamples = make([]*ExemplarItem, 0, e.k)
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
	e := new(exemplar)
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
func (e *exemplar) Offer(latency time.Duration, requestID string, err error, msg string) {
	latInt := int64(latency)
	isError := err != nil

	// Fast path check to avoid locking for requests that are neither slowest nor fastest.
	// This is an optimistic check. Race conditions are handled by the lock below.
	minLat := e.minLatency.Load()
	maxLat := e.maxLatency.Load()

	// If it's not an error, and the latency is between the fastest and slowest top-K,
	// we can skip the heap updates. However, we must still consider it for reservoir sampling.
	if !isError && latInt <= minLat && latInt >= maxLat {
		e.mu.Lock()
		e.updateAverageSample(&ExemplarItem{
			Latency:   latency,
			RequestID: requestID,
			Err:       err,
			Msg:       msg,
		})
		e.mu.Unlock()
		return
	}

	newItem := &ExemplarItem{
		Latency:   latency,
		RequestID: requestID,
		Err:       err,
		Msg:       msg,
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.updateSlowest(newItem)
	e.updateFastest(newItem)
	e.updateAverageSample(newItem)
	if isError {
		e.updateFailureSample(newItem)
	}
}

func (e *exemplar) updateSlowest(item *ExemplarItem) {
	latInt := int64(item.Latency)
	if len(e.slowest) < e.k {
		heap.Push(&e.slowest, item)
		if len(e.slowest) == e.k {
			e.minLatency.Store(int64(e.slowest[0].Latency))
		}
	} else if latInt > int64(e.slowest[0].Latency) {
		e.slowest[0] = item
		heap.Fix(&e.slowest, 0)
		e.minLatency.Store(int64(e.slowest[0].Latency))
	}
}

func (e *exemplar) updateFastest(item *ExemplarItem) {
	latInt := int64(item.Latency)
	if len(e.fastest) < e.k {
		heap.Push(&e.fastest, item)
		if len(e.fastest) == e.k {
			e.maxLatency.Store(int64(e.fastest[0].Latency))
		}
	} else if latInt < int64(e.fastest[0].Latency) {
		e.fastest[0] = item
		heap.Fix(&e.fastest, 0)
		e.maxLatency.Store(int64(e.fastest[0].Latency))
	}
}

func (e *exemplar) updateAverageSample(item *ExemplarItem) {
	e.avgCount++
	if len(e.avgSamples) < e.k {
		e.avgSamples = append(e.avgSamples, item)
	} else {
		j := rand.Uint64N(e.avgCount)
		if j < uint64(e.k) {
			e.avgSamples[j] = item
		}
	}
}

func (e *exemplar) updateFailureSample(item *ExemplarItem) {
	e.failureCount++
	if len(e.failureSamples) < e.k {
		e.failureSamples = append(e.failureSamples, item)
	} else {
		j := rand.Uint64N(e.failureCount)
		if j < uint64(e.k) {
			e.failureSamples[j] = item
		}
	}
}

// Snapshot returns a snapshot of the exemplars.
func (e *exemplar) Snapshot() []*ExemplarItem {
	// For backward compatibility, return Slowest.
	e.mu.Lock()
	items := slices.Clone(e.slowest)
	e.mu.Unlock()

	slices.SortFunc(items, func(a, b *ExemplarItem) int {
		return cmp.Compare(b.Latency, a.Latency)
	})
	return items
}

// DetailedSnapshot returns all categories.
func (e *exemplar) DetailedSnapshot() (*ExemplarDetails, error) {
	e.mu.Lock()
	snap := &ExemplarDetails{
		Slowest:  slices.Clone(e.slowest),
		Fastest:  slices.Clone(e.fastest),
		Average:  slices.Clone(e.avgSamples),
		Failures: slices.Clone(e.failureSamples),
	}
	e.mu.Unlock()

	slices.SortFunc(snap.Slowest, func(a, b *ExemplarItem) int {
		return cmp.Compare(b.Latency, a.Latency) // Descending
	})

	slices.SortFunc(snap.Fastest, func(a, b *ExemplarItem) int {
		return cmp.Compare(a.Latency, b.Latency) // Ascending
	})

	slices.SortFunc(snap.Average, func(a, b *ExemplarItem) int {
		return cmp.Compare(b.Latency, a.Latency) // Descending
	})

	slices.SortFunc(snap.Failures, func(a, b *ExemplarItem) int {
		return cmp.Compare(b.Latency, a.Latency) // Descending
	})

	return snap, nil
}

// Merge merges another exemplar into this one.
func (e *exemplar) Merge(other Exemplar) error {
	if other == nil {
		return nil
	}
	o, ok := other.(*exemplar)
	if !ok {
		// Fallback to offering items one by one if the type is not the same.
		details, _ := other.DetailedSnapshot()
		if details == nil {
			return nil
		}
		for _, ex := range details.Slowest {
			e.Offer(ex.Latency, ex.RequestID, ex.Err, ex.Msg)
		}
		for _, ex := range details.Fastest {
			e.Offer(ex.Latency, ex.RequestID, ex.Err, ex.Msg)
		}
		for _, ex := range details.Average {
			e.Offer(ex.Latency, ex.RequestID, ex.Err, ex.Msg)
		}
		for _, ex := range details.Failures {
			e.Offer(ex.Latency, ex.RequestID, ex.Err, ex.Msg)
		}
		return nil
	}

	e.mu.Lock()
	o.mu.Lock()
	defer e.mu.Unlock()
	defer o.mu.Unlock()

	// Merge heaps by offering all items from the other heap.
	for _, item := range o.slowest {
		e.updateSlowest(item)
	}
	for _, item := range o.fastest {
		e.updateFastest(item)
	}

	// Merge reservoirs using a weighted algorithm.
	e.avgSamples = e.mergeReservoir(e.avgSamples, o.avgSamples, e.avgCount, o.avgCount)
	e.failureSamples = e.mergeReservoir(e.failureSamples, o.failureSamples, e.failureCount, o.failureCount)

	e.avgCount += o.avgCount
	e.failureCount += o.failureCount

	return nil
}

// mergeReservoir merges two reservoir samples (dst and src) into a new reservoir.
// It uses a weighted selection algorithm to ensure that the merged reservoir is a
// statistically valid sample of the combined population.
func (e *exemplar) mergeReservoir(dst, src []*ExemplarItem, n1, n2 uint64) []*ExemplarItem {
	if n1 == 0 {
		return src
	}
	if n2 == 0 {
		return dst
	}

	n := n1 + n2
	k := e.k

	// If the combined size is less than or equal to k, just return the combined slice.
	if len(dst)+len(src) <= k {
		return append(dst, src...)
	}

	// Create a new reservoir by probabilistically selecting items from both reservoirs.
	newReservoir := make([]*ExemplarItem, 0, k)
	for range k {
		// Decide which reservoir to draw from based on their relative weights.
		if rand.Uint64N(n) < n1 {
			// Draw from the first reservoir.
			if len(dst) > 0 {
				idx := rand.IntN(len(dst))
				newReservoir = append(newReservoir, dst[idx])
				// Remove the selected item to avoid picking it again.
				// Use "Swap and Remove" pattern for O(1) performance.
				// 1. Swap the selected item with the last item.
				dst[idx] = dst[len(dst)-1]
				// 2. Zero out the last item to prevent memory leaks (pointers).
				dst[len(dst)-1] = nil
				// 3. Truncate the slice.
				dst = slices.Delete(dst, len(dst)-1, len(dst))
				n1--
			}
		} else {
			// Draw from the second reservoir.
			if len(src) > 0 {
				idx := rand.IntN(len(src))
				newReservoir = append(newReservoir, src[idx])
				// Remove the selected item.
				// Use "Swap and Remove" pattern for O(1) performance.
				src[idx] = src[len(src)-1]
				src[len(src)-1] = nil
				src = slices.Delete(src, len(src)-1, len(src))
				n2--
			}
		}
		n--
	}
	return newReservoir
}

// Clone returns a deep copy.
func (e *exemplar) Clone() Exemplar {
	newE := new(exemplar)
	newE.k = e.k

	e.mu.Lock()
	defer e.mu.Unlock()

	copyTo := func(dst *[]*ExemplarItem, src []*ExemplarItem) {
		if cap(*dst) < len(src) {
			*dst = make([]*ExemplarItem, len(src), cap(src))
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
	newE.fastest = make(smallestLatencyHeap, len(e.fastest), cap(e.fastest))
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

// ExemplarItem is an item in the priority queue.
type ExemplarItem struct {
	Latency   time.Duration
	RequestID string
	Err       error
	Msg       string
}

// priorityQueue implements min-heap.
type priorityQueue []*ExemplarItem

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Latency < pq[j].Latency
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*ExemplarItem)
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

// smallestLatencyHeap implements a Max-Heap to store the K smallest latencies.
// The root of the heap is the largest value among the smallest K.
// When a new value arrives that is smaller than the root, we pop the root (largest)
// and push the new value, thus maintaining the K smallest values.
type smallestLatencyHeap []*ExemplarItem

func (pq smallestLatencyHeap) Len() int { return len(pq) }

// Less returns true if element i is "smaller" than j.
// Since this is a Max-Heap (to evict the largest), "smaller" means "greater latency".
func (pq smallestLatencyHeap) Less(i, j int) bool {
	return pq[i].Latency > pq[j].Latency
}

func (pq smallestLatencyHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *smallestLatencyHeap) Push(x any) {
	item := x.(*ExemplarItem)
	*pq = append(*pq, item)
}

func (pq *smallestLatencyHeap) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

type ExemplarDetails struct {
	Slowest  []*ExemplarItem
	Fastest  []*ExemplarItem
	Average  []*ExemplarItem
	Failures []*ExemplarItem
}
