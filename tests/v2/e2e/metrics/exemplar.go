//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/zeebo/xxh3"
)

const (
	// defaultAvgSamplingRate is the default sampling rate for average exemplars (1/16).
	defaultAvgSamplingRate = 16

	// defaultCapacity is the default capacity for exemplar heaps/reservoirs.
	defaultCapacity = 10
)

// exemplar is a thread-safe exemplar storage.
// It implements the Exemplar interface.
type exemplar struct {
	slowest        priorityQueue
	fastest        smallestLatencyHeap
	avgSamples     []*ExemplarItem
	failureSamples []*ExemplarItem
	k              int
	samplingRate   uint64
	storedCount    atomic.Uint64
	avgCount       uint64
	failureCount   uint64
	minLatency     atomic.Int64
	maxLatency     atomic.Int64
	mu             sync.Mutex
}

// shardedExemplar is a sharded wrapper around exemplar.
type shardedExemplar struct {
	shards []*exemplar
}

// exemplarConfig holds configuration for Exemplar.
type exemplarConfig struct {
	Capacity     int `json:"capacity"      yaml:"capacity"`
	NumShards    int `json:"num_shards"    yaml:"num_shards"`
	SamplingRate int `json:"sampling_rate" yaml:"sampling_rate"`
}

// Init initializes the exemplar with the given options.
func (e *exemplar) Init(opts ...ExemplarOption) error {
	cfg := exemplarConfig{
		Capacity:     defaultCapacity,
		SamplingRate: defaultAvgSamplingRate,
	}
	// Apply user options
	for _, opt := range opts {
		opt(&cfg)
	}
	e.k = max(cfg.Capacity, 1)

	rate := max(cfg.SamplingRate, 1)
	// Ensure samplingRate is a power of 2 for efficient bitwise sampling
	if rate&(rate-1) != 0 {
		return errors.New("samplingRate must be a power of 2")
	}
	e.samplingRate = uint64(rate)
	e.initHeaps()
	return nil
}

func (e *exemplar) initHeaps() {
	e.mu.Lock()
	defer e.mu.Unlock()
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
	e.storedCount.Store(0)
	e.avgCount = 0
	e.failureCount = 0
}

// NewExemplar creates a new Exemplar with the given options.
func NewExemplar(opts ...ExemplarOption) (Exemplar, error) {
	cfg := exemplarConfig{
		Capacity:     defaultCapacity,
		SamplingRate: defaultAvgSamplingRate,
	}
	for _, opt := range append(defaultExemplarOpts, opts...) {
		opt(&cfg)
	}

	if cfg.NumShards <= 1 {
		e := new(exemplar)
		// We need to pass capacity via options to Init, but options now take *exemplarConfig.
		applyConfig := func(c *exemplarConfig) {
			*c = cfg
		}
		err := e.Init(applyConfig)
		if err != nil {
			return nil, err
		}
		return e, nil
	}

	se := &shardedExemplar{
		shards: make([]*exemplar, cfg.NumShards),
	}
	applyConfig := func(c *exemplarConfig) {
		*c = cfg
	}
	for i := range se.shards {
		e := new(exemplar)
		err := e.Init(applyConfig)
		if err != nil {
			return nil, err
		}
		se.shards[i] = e
	}
	return se, nil
}

// Reset resets the sharded exemplar.
func (se *shardedExemplar) Reset() {
	for _, e := range se.shards {
		e.Reset()
	}
}

// Reset resets the exemplar to its initial state.
func (e *exemplar) Reset() {
	e.initHeaps()
}

// Offer adds a request to the sharded exemplar.
func (se *shardedExemplar) Offer(latency time.Duration, requestID string, err error, msg string) {
	shardIdx := shardIndex(xxh3.HashString(requestID), len(se.shards))
	se.shards[shardIdx].Offer(latency, requestID, err, msg)
}

// Offer adds a request to the exemplar.
func (e *exemplar) Offer(latency time.Duration, requestID string, err error, msg string) {
	latInt := int64(latency)
	isError := err != nil

	// Check if we have enough samples in the slowest heap.
	// If not, force slow path to populate the heap.
	// storedCount tracks the number of items in the slowest heap.
	if e.storedCount.Load() >= uint64(e.k) {
		// Fast path check to avoid locking for requests that are neither slowest nor fastest.
		// This is an optimistic check. Race conditions are handled by the lock below.
		minLat := e.minLatency.Load()
		maxLat := e.maxLatency.Load()

		// Fast Path (Check 1): avoid locking if it's not an outlier candidate.
		// minLatency is the lower bound of the Slowest (Top-K) bucket.
		// maxLatency is the upper bound of the Fastest (Bottom-K) bucket.
		if !isError && latInt <= minLat && latInt >= maxLat {
			// Probabilistic sampling: only record 1 out of samplingRate samples.
			// We use a fast bitwise check assuming samplingRate is a power of 2.
			if rand.Uint64()&(e.samplingRate-1) != 0 { //nolint:gosec // use math/rand/v2 for performance
				return
			}

			// Must lock before updating reservoir state
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
	}

	newItem := &ExemplarItem{
		Latency:   latency,
		RequestID: requestID,
		Err:       err,
		Msg:       msg,
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Double Check (Check 2): Re-read atomic values after acquiring lock
	if !isError && e.storedCount.Load() >= uint64(e.k) {
		minLat := e.minLatency.Load()
		maxLat := e.maxLatency.Load()
		if latInt <= minLat && latInt >= maxLat {
			// Race condition: another goroutine updated min/max while we waited for lock.
			// This request is no longer an outlier. Treat as average.
			e.updateAverageSample(newItem)
			return
		}
	}

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
		e.storedCount.Add(1)
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
		j := rand.Uint64N(e.avgCount) //nolint:gosec // use math/rand/v2 for performance
		// Ensure e.k is non-negative before casting
		if e.k > 0 && j < uint64(e.k) {
			e.avgSamples[j] = item
		}
	}
}

func (e *exemplar) updateFailureSample(item *ExemplarItem) {
	e.failureCount++
	if len(e.failureSamples) < e.k {
		e.failureSamples = append(e.failureSamples, item)
	} else {
		j := rand.Uint64N(e.failureCount) //nolint:gosec // use math/rand/v2 for performance
		// Ensure e.k is non-negative before casting
		if e.k > 0 && j < uint64(e.k) {
			e.failureSamples[j] = item
		}
	}
}

// Snapshot returns a snapshot of the exemplars (Slowest only for backward compatibility).
func (se *shardedExemplar) Snapshot() []*ExemplarItem {
	details, _ := se.DetailedSnapshot()
	if details == nil {
		return nil
	}
	return details.Slowest
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
func (se *shardedExemplar) DetailedSnapshot() (*ExemplarDetails, error) {
	if len(se.shards) == 0 {
		return nil, errors.New("exemplar: no shards available")
	}

	// Create a temporary exemplar to merge all shards
	k := se.shards[0].k
	rate := se.shards[0].samplingRate
	merged := new(exemplar)
	// Single shard, preserve config
	merged.Init(WithExemplarCapacity(k), WithExemplarSamplingRate(int(rate)))

	for _, shard := range se.shards {
		shard.mu.Lock()

		slowest := slices.Clone(shard.slowest)
		fastest := slices.Clone(shard.fastest)
		avg := slices.Clone(shard.avgSamples)
		failures := slices.Clone(shard.failureSamples)
		avgCount := shard.avgCount
		failCount := shard.failureCount

		shard.mu.Unlock()

		for _, item := range slowest {
			merged.updateSlowest(item)
		}
		for _, item := range fastest {
			merged.updateFastest(item)
		}

		merged.avgSamples = mergeReservoir(merged.avgSamples, avg, merged.avgCount, avgCount, k)
		merged.avgCount += merged.avgCount

		merged.failureSamples = mergeReservoir(merged.failureSamples, failures, merged.failureCount, failCount, k)
		merged.failureCount += failCount
	}

	return merged.DetailedSnapshot()
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
		return cmp.Compare(b.Latency, a.Latency)
	})
	slices.SortFunc(snap.Fastest, func(a, b *ExemplarItem) int {
		return cmp.Compare(a.Latency, b.Latency)
	})
	slices.SortFunc(snap.Average, func(a, b *ExemplarItem) int {
		return cmp.Compare(b.Latency, a.Latency)
	})
	slices.SortFunc(snap.Failures, func(a, b *ExemplarItem) int {
		return cmp.Compare(b.Latency, a.Latency)
	})

	return snap, nil
}

// Merge merges another exemplar into this one.
func (se *shardedExemplar) Merge(other Exemplar) error {
	if other == nil {
		return nil
	}
	if o, ok := other.(*shardedExemplar); ok {
		return mergeShards(se.shards, o.shards)
	}
	if _, ok := other.(*exemplar); ok {
		return errors.New("cannot merge single exemplar into sharded exemplar")
	}
	return errors.New("unknown exemplar type")
}

// Merge merges another exemplar into this one.
func (e *exemplar) Merge(other Exemplar) error {
	if other == nil {
		return nil
	}
	if o, ok := other.(*exemplar); ok {
		return e.mergeExemplar(o)
	}
	if _, ok := other.(*shardedExemplar); ok {
		return errors.New("cannot merge sharded exemplar into single exemplar")
	}
	return errors.New("unknown exemplar type")
}

func (e *exemplar) mergeExemplar(src *exemplar) error {
	src.mu.Lock()
	slowest := slices.Clone(src.slowest)
	fastest := slices.Clone(src.fastest)
	avg := slices.Clone(src.avgSamples)
	failures := slices.Clone(src.failureSamples)
	avgCount := src.avgCount
	failCount := src.failureCount
	src.mu.Unlock()

	e.mu.Lock()
	defer e.mu.Unlock()

	for _, item := range slowest {
		e.updateSlowest(item)
	}
	for _, item := range fastest {
		e.updateFastest(item)
	}

	e.avgSamples = mergeReservoir(e.avgSamples, avg, e.avgCount, avgCount, e.k)
	e.avgCount += avgCount

	e.failureSamples = mergeReservoir(e.failureSamples, failures, e.failureCount, failCount, e.k)
	e.failureCount += failCount

	return nil
}

// mergeReservoir merges two reservoir samples (dst and src) into a new reservoir.
// It uses a weighted selection algorithm to ensure that the merged reservoir is a
// statistically valid sample of the combined population.
//
// NOTE: dst and src slices may be modified in-place to avoid allocations.
func mergeReservoir(dst, src []*ExemplarItem, n1, n2 uint64, k int) []*ExemplarItem {
	if n1 == 0 {
		return src
	}
	if n2 == 0 {
		return dst
	}

	// If the combined size is less than or equal to k, just return the combined slice.
	if len(dst)+len(src) <= k {
		// Ensure result fits in k if possible, but here we just append.
		// If capacity allows, append is allocation-free.
		return append(dst, src...)
	}

	// Create a new reservoir by probabilistically selecting items from both reservoirs.
	newReservoir := make([]*ExemplarItem, 0, k)

	// Weighted reservoir sampling merge
	// Algorithm:
	// We want to select k items from the union of two reservoirs representing populations of size n1 and n2.
	// Total population N = n1 + n2.
	// We iterate k times. In each step, we pick an item from reservoir 1 with probability n1/N, and from reservoir 2 with probability n2/N.
	// Once an item is picked, we remove it from the source reservoir and decrement the corresponding population count (n1 or n2) and N.

	n := n1 + n2
	currN1 := n1
	currN2 := n2

	// We use dst and src directly, modifying them (Swap and Remove).
	// This avoids allocating clones.
	// Callers must pass copies or disposable slices.
	d := dst
	s := src

	for range k {
		if rand.Uint64N(n) < currN1 { //nolint:gosec // use math/rand/v2 for performance
			// Draw from the first reservoir.
			if len(d) > 0 {
				idx := rand.IntN(len(d)) //nolint:gosec // use math/rand/v2 for performance
				newReservoir = append(newReservoir, d[idx])
				// Remove the selected item.
				// Use "Swap and Remove" pattern for O(1) performance.
				d[idx] = d[len(d)-1]
				d = d[:len(d)-1]
				currN1--
			}
		} else {
			// Draw from the second reservoir.
			if len(s) > 0 {
				idx := rand.IntN(len(s)) //nolint:gosec // use math/rand/v2 for performance
				newReservoir = append(newReservoir, s[idx])
				// Remove the selected item.
				// Use "Swap and Remove" pattern for O(1) performance.
				s[idx] = s[len(s)-1]
				s = s[:len(s)-1]
				currN2--
			}
		}
		n--
	}
	return newReservoir
}

// Clone returns a deep copy.
func (se *shardedExemplar) Clone() Exemplar {
	newSE := &shardedExemplar{
		shards: make([]*exemplar, len(se.shards)),
	}
	for i, e := range se.shards {
		cloned := e.Clone()
		c, ok := cloned.(*exemplar)
		if ok {
			newSE.shards[i] = c
		}
	}
	return newSE
}

// Clone returns a deep copy.
func (e *exemplar) Clone() Exemplar {
	newE := new(exemplar)
	newE.k = e.k
	newE.samplingRate = e.samplingRate

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
	newE.storedCount.Store(e.storedCount.Load())

	return newE
}

// ExemplarItem is an item in the priority queue.
type ExemplarItem struct {
	Err       error         `json:"err,omitempty"`
	RequestID string        `json:"request_id,omitempty"`
	Msg       string        `json:"msg,omitempty"`
	Latency   time.Duration `json:"latency,omitempty"`
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
	item, ok := x.(*ExemplarItem)
	if ok {
		*pq = append(*pq, item)
	}
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
	item, ok := x.(*ExemplarItem)
	if ok {
		*pq = append(*pq, item)
	}
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
	Slowest  []*ExemplarItem `json:"slowest,omitempty"`
	Fastest  []*ExemplarItem `json:"fastest,omitempty"`
	Average  []*ExemplarItem `json:"average,omitempty"`
	Failures []*ExemplarItem `json:"failures,omitempty"`
}
