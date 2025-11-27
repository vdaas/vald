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
	"github.com/zeebo/xxh3"
)

// exemplar is a sharded exemplar storage.
// It implements the Exemplar interface.
type exemplar struct {
	shards    []exemplarShard
	numShards int
	k         int // Capacity per shard
}

// exemplarShard holds samples of requests in different categories.
// It uses mutexes to protect its state.
type exemplarShard struct {
	mu sync.Mutex

	// Categories
	slowest  priorityQueue       // Min-heap (Top-K Max Latency)
	fastest  smallestLatencyHeap // Max-heap (Top-K Min Latency) to evict largest, keeping smallest.
	failures priorityQueue       // Min-heap (Top-K Slowest Failures)

	avgSamples     []*ExemplarItem // Reservoir for representative samples
	failureSamples []*ExemplarItem // Reservoir for failure samples
	avgCount       uint64          // Total count seen for average reservoir
	failureCount   uint64          // Total count seen for failure reservoir

	minLatency atomic.Int64 // Minimum latency in the 'slowest' heap (fast path)
	maxLatency atomic.Int64 // Maximum latency in the 'fastest' heap (fast path)
}

// Init initializes the exemplar with the given options.
func (e *exemplar) Init(opts ...ExemplarOption) {
	for _, opt := range opts {
		opt(e)
	}
	e.k = max(e.k, 1)
	e.numShards = max(e.numShards, 1)

	if len(e.shards) < e.numShards {
		e.shards = make([]exemplarShard, e.numShards)
	} else {
		e.shards = e.shards[:e.numShards]
	}

	for i := range e.shards {
		e.shards[i].init(e.k)
	}
}

// NewExemplar creates a new Exemplar with the given options.
func NewExemplar(opts ...ExemplarOption) Exemplar {
	e := new(exemplar)
	e.Init(opts...)
	return e
}

func (e *exemplarShard) init(k int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.slowest == nil {
		e.slowest = make(priorityQueue, 0, k)
	} else {
		e.slowest = e.slowest[:0]
	}
	if e.fastest == nil {
		e.fastest = make(smallestLatencyHeap, 0, k)
	} else {
		e.fastest = e.fastest[:0]
	}
	if e.avgSamples == nil {
		e.avgSamples = make([]*ExemplarItem, 0, k)
	} else {
		e.avgSamples = e.avgSamples[:0]
	}
	if e.failureSamples == nil {
		e.failureSamples = make([]*ExemplarItem, 0, k)
	} else {
		e.failureSamples = e.failureSamples[:0]
	}
	e.minLatency.Store(0)
	e.maxLatency.Store(0)
	e.avgCount = 0
	e.failureCount = 0
}

// Reset resets the exemplar to its initial state.
func (e *exemplar) Reset() {
	for i := range e.shards {
		e.shards[i].init(e.k)
	}
}

// Offer adds a request to the exemplar categories.
func (e *exemplar) Offer(latency time.Duration, requestID string, err error, msg string) {
	shardIdx := 0
	if e.numShards > 1 {
		shardIdx = int(xxh3.HashString(requestID) % uint64(e.numShards))
	}
	e.shards[shardIdx].offer(e.k, latency, requestID, err, msg)
}

func (s *exemplarShard) offer(k int, latency time.Duration, requestID string, err error, msg string) {
	latInt := int64(latency)
	isError := err != nil

	// Fast path check to avoid locking for requests that are neither slowest nor fastest.
	// This is an optimistic check. Race conditions are handled by the lock below.
	minLat := s.minLatency.Load()
	maxLat := s.maxLatency.Load()

	// If it's not an error, and the latency is between the fastest and slowest top-K,
	// we can skip the heap updates. However, we must still consider it for reservoir sampling.
	if !isError && latInt <= minLat && latInt >= maxLat {
		s.mu.Lock()
		s.updateAverageSample(k, &ExemplarItem{
			Latency:   latency,
			RequestID: requestID,
			Err:       err,
			Msg:       msg,
		})
		s.mu.Unlock()
		return
	}

	newItem := &ExemplarItem{
		Latency:   latency,
		RequestID: requestID,
		Err:       err,
		Msg:       msg,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.updateSlowest(k, newItem)
	s.updateFastest(k, newItem)
	s.updateAverageSample(k, newItem)
	if isError {
		s.updateFailureSample(k, newItem)
	}
}

func (s *exemplarShard) updateSlowest(k int, item *ExemplarItem) {
	latInt := int64(item.Latency)
	if len(s.slowest) < k {
		heap.Push(&s.slowest, item)
		if len(s.slowest) == k {
			s.minLatency.Store(int64(s.slowest[0].Latency))
		}
	} else if latInt > int64(s.slowest[0].Latency) {
		s.slowest[0] = item
		heap.Fix(&s.slowest, 0)
		s.minLatency.Store(int64(s.slowest[0].Latency))
	}
}

func (s *exemplarShard) updateFastest(k int, item *ExemplarItem) {
	latInt := int64(item.Latency)
	if len(s.fastest) < k {
		heap.Push(&s.fastest, item)
		if len(s.fastest) == k {
			s.maxLatency.Store(int64(s.fastest[0].Latency))
		}
	} else if latInt < int64(s.fastest[0].Latency) {
		s.fastest[0] = item
		heap.Fix(&s.fastest, 0)
		s.maxLatency.Store(int64(s.fastest[0].Latency))
	}
}

func (s *exemplarShard) updateAverageSample(k int, item *ExemplarItem) {
	s.avgCount++
	if len(s.avgSamples) < k {
		s.avgSamples = append(s.avgSamples, item)
	} else {
		j := rand.Uint64N(s.avgCount)
		if j < uint64(k) {
			s.avgSamples[j] = item
		}
	}
}

func (s *exemplarShard) updateFailureSample(k int, item *ExemplarItem) {
	s.failureCount++
	if len(s.failureSamples) < k {
		s.failureSamples = append(s.failureSamples, item)
	} else {
		j := rand.Uint64N(s.failureCount)
		if j < uint64(k) {
			s.failureSamples[j] = item
		}
	}
}

// Snapshot returns a snapshot of the exemplars (Slowest only for backward compatibility).
func (e *exemplar) Snapshot() []*ExemplarItem {
	details, _ := e.DetailedSnapshot()
	if details == nil {
		return nil
	}
	return details.Slowest
}

// DetailedSnapshot returns all categories.
func (e *exemplar) DetailedSnapshot() (*ExemplarDetails, error) {
	// Aggregate all shards.
	// Since we might have duplicates across shards (unlikely with RequestID hash but possible if logic changes),
	// or we just have distinct sets.
	// We merge them into a single set of size K (or total size and then pick Top K?).
	// The interface implies return type `ExemplarDetails`.
	// If we return *all* items from all shards, it might be too big if NumShards is large.
	// But usually we want the global Top K.
	// So we should merge heaps.

	merged := new(exemplar)
	merged.Init(WithExemplarCapacity(e.k), WithExemplarNumShards(1))

	// Merge all shards into one temporary exemplar
	for i := range e.shards {
		shard := &e.shards[i]
		shard.mu.Lock()

		// Copy items to avoid holding lock while merging
		slowest := slices.Clone(shard.slowest)
		fastest := slices.Clone(shard.fastest)
		avg := slices.Clone(shard.avgSamples)
		failures := slices.Clone(shard.failureSamples)
		avgCount := shard.avgCount
		failCount := shard.failureCount

		shard.mu.Unlock()

		// Merge into merged (which is a single shard exemplar effectively)
		// We can directly access merged.shards[0]
		mShard := &merged.shards[0]

		for _, item := range slowest {
			mShard.updateSlowest(merged.k, item)
		}
		for _, item := range fastest {
			mShard.updateFastest(merged.k, item)
		}

		// For reservoirs, we use the weighted merge
		mShard.avgSamples = mergeReservoir(mShard.avgSamples, avg, mShard.avgCount, avgCount, merged.k)
		mShard.avgCount += avgCount

		mShard.failureSamples = mergeReservoir(mShard.failureSamples, failures, mShard.failureCount, failCount, merged.k)
		mShard.failureCount += failCount
	}

	mShard := &merged.shards[0]
	snap := &ExemplarDetails{
		Slowest:  mShard.slowest,
		Fastest:  mShard.fastest,
		Average:  mShard.avgSamples,
		Failures: mShard.failureSamples,
	}

	// Sort results
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
		// Fallback
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

	// Merge shard by shard if compatible?
	// Shards are hashed by RequestID. If both use same hash, we can merge corresponding shards.
	if e.numShards == o.numShards {
		for i := range e.shards {
			e.mergeShard(&e.shards[i], &o.shards[i], e.k)
		}
	} else {
		// If incompatible shards, we iterate through all items of other and Offer.
		// Or simpler: get detailed snapshot of other and Offer all.
		details, _ := o.DetailedSnapshot()
		if details != nil {
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
		}
	}
	return nil
}

func (e *exemplar) mergeShard(dst, src *exemplarShard, k int) {
	src.mu.Lock()
	slowest := slices.Clone(src.slowest)
	fastest := slices.Clone(src.fastest)
	avg := slices.Clone(src.avgSamples)
	failures := slices.Clone(src.failureSamples)
	avgCount := src.avgCount
	failCount := src.failureCount
	src.mu.Unlock()

	dst.mu.Lock()
	defer dst.mu.Unlock()

	for _, item := range slowest {
		dst.updateSlowest(k, item)
	}
	for _, item := range fastest {
		dst.updateFastest(k, item)
	}

	dst.avgSamples = mergeReservoir(dst.avgSamples, avg, dst.avgCount, avgCount, k)
	dst.avgCount += avgCount

	dst.failureSamples = mergeReservoir(dst.failureSamples, failures, dst.failureCount, failCount, k)
	dst.failureCount += failCount
}

// mergeReservoir merges two reservoir samples (dst and src) into a new reservoir.
func mergeReservoir(dst, src []*ExemplarItem, n1, n2 uint64, k int) []*ExemplarItem {
	if n1 == 0 {
		return src
	}
	if n2 == 0 {
		return dst
	}

	n := n1 + n2

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
				// Remove the selected item
				dst[idx] = dst[len(dst)-1]
				dst[len(dst)-1] = nil
				dst = slices.Delete(dst, len(dst)-1, len(dst))
				n1--
			}
		} else {
			// Draw from the second reservoir.
			if len(src) > 0 {
				idx := rand.IntN(len(src))
				newReservoir = append(newReservoir, src[idx])
				// Remove the selected item
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
	newE.numShards = e.numShards
	newE.shards = make([]exemplarShard, len(e.shards))

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

	for i := range e.shards {
		src := &e.shards[i]
		dst := &newE.shards[i]

		src.mu.Lock()

		// Copy Heaps
		dst.slowest = make(priorityQueue, len(src.slowest), cap(src.slowest))
		for i, it := range src.slowest {
			v := *it
			dst.slowest[i] = &v
		}
		dst.fastest = make(smallestLatencyHeap, len(src.fastest), cap(src.fastest))
		for i, it := range src.fastest {
			v := *it
			dst.fastest[i] = &v
		}

		// Copy Reservoirs
		copyTo(&dst.avgSamples, src.avgSamples)
		copyTo(&dst.failureSamples, src.failureSamples)

		dst.avgCount = src.avgCount
		dst.failureCount = src.failureCount
		dst.minLatency.Store(src.minLatency.Load())
		dst.maxLatency.Store(src.maxLatency.Load())

		src.mu.Unlock()
	}

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
type smallestLatencyHeap []*ExemplarItem

func (pq smallestLatencyHeap) Len() int { return len(pq) }

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

// Helper needed for clone
// sliceHeader is a stripped-down version of reflect.SliceHeader used for unsafe pointer conversions.
// Defined in interface.go, but we can't access private members across files if we don't need to.
// Here we used standard loop, so no unsafe needed.
