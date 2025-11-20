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
	"slices"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/sync"
)

// exemplar holds a sample of high-latency requests.
// It uses a mutex-protected priority queue (min-heap) to store the top k requests
// with the highest latencies.
type exemplar struct {
	pq         priorityQueue // Min-heap of exemplars.
	minLatency atomic.Int64  // Minimum latency in the heap, for lock-free check.
	mu         sync.Mutex
	k          int // The maximum number of exemplars to store.
}

// Init initializes the exemplar with the given options.
func (e *exemplar) Init(opts ...ExemplarOption) {
	for _, opt := range opts {
		opt(e)
	}
	e.k = max(e.k, 1)
	e.mu.Lock()
	if e.pq == nil {
		e.pq = make(priorityQueue, 0, e.k)
	} else {
		e.pq = e.pq[:0]
	}
	e.minLatency.Store(0)
	e.mu.Unlock()
}

// NewExemplar creates a new Exemplar with the given options.
func NewExemplar(opts ...ExemplarOption) Exemplar {
	e := exemplarPool.Get()
	e.Init(opts...)
	return e
}

// Reset resets the exemplar to its initial state, clearing all data but keeping capacity.
func (e *exemplar) Reset() {
	e.mu.Lock()
	// Clear the slice but keep capacity
	for i := range e.pq {
		e.pq[i] = nil
	}
	e.pq = e.pq[:0]
	e.minLatency.Store(0)
	e.mu.Unlock()
}

// Offer adds a request to the exemplar.
// If the priority queue is not full, the new item is added.
// If the priority queue is full and the new item's latency is greater than the minimum latency in the queue,
// the new item replaces the minimum latency item.
func (e *exemplar) Offer(latency time.Duration, requestID string) {
	// Fast-path: if the heap is full and the new latency is smaller than the current minimum,
	// we can return immediately without locking.
	minLat := e.minLatency.Load()
	if minLat > 0 && int64(latency) <= minLat {
		return
	}

	newItem := &item{
		latency:   latency,
		requestID: requestID,
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.pq) < e.k {
		heap.Push(&e.pq, newItem)
		// If we reached capacity, set the minLatency
		if len(e.pq) == e.k {
			e.minLatency.Store(int64(e.pq[0].latency))
		}
	} else if latency > e.pq[0].latency {
		// Replace the smallest item (at index 0) with the new item
		// and fix the heap invariant.
		e.pq[0] = newItem
		heap.Fix(&e.pq, 0)
		e.minLatency.Store(int64(e.pq[0].latency))
	}
}

// Snapshot returns a snapshot of the exemplars.
func (e *exemplar) Snapshot() []*item {
	e.mu.Lock()
	defer e.mu.Unlock()

	items := slices.Clone(e.pq)

	// Sort items by latency in descending order.
	slices.SortFunc(items, func(a, b *item) int {
		return cmp.Compare(b.latency, a.latency)
	})

	return items
}

// Clone returns a deep copy of the exemplar.
func (e *exemplar) Clone() Exemplar {
	newE := exemplarPool.Get()
	newE.Reset()
	newE.k = e.k

	e.mu.Lock()
	defer e.mu.Unlock()

	if cap(newE.pq) < len(e.pq) {
		newE.pq = make(priorityQueue, len(e.pq), cap(e.pq))
	} else {
		newE.pq = newE.pq[:len(e.pq)]
	}

	// Deep copy items
	for i, item := range e.pq {
		if item != nil {
			val := *item
			newE.pq[i] = &val
		}
	}
	newE.minLatency.Store(e.minLatency.Load())

	return newE
}

// item is an item in the priority queue, representing a single request exemplar.
// It is unexported to encapsulate the implementation details of the priority queue.
type item struct {
	requestID string
	latency   time.Duration
}

// priorityQueue implements heap.Interface and is a min-heap of items.
// It is unexported to encapsulate the implementation details of the Exemplar.
type priorityQueue []*item

// Len returns the number of items in the priority queue.
func (pq priorityQueue) Len() int { return len(pq) }

// Less returns true if the item at index i has a smaller latency than the item at index j.
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].latency < pq[j].latency
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	item, ok := x.(*item)
	if !ok {
		return
	}
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
