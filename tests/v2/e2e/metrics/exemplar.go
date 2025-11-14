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
	"container/heap"
	"sync/atomic"
	"time"
)

// Exemplar holds a sample of high-latency requests.
// It uses a lock-free priority queue (min-heap) to store the top k requests
// with the highest latencies. This allows for efficient and concurrent updates
// without blocking.
type Exemplar struct {
	k  int // The maximum number of exemplars to store.
	pq atomic.Pointer[priorityQueue]
}

// NewExemplar creates a new Exemplar with a capacity of k.
// It initializes a lock-free priority queue to store the exemplars.
func NewExemplar(k int) *Exemplar {
	k = max(k, 1)
	initialPQ := make(priorityQueue, 0, k)
	e := &Exemplar{
		k: k,
	}
	e.pq.Store(&initialPQ)
	return e
}

// Offer adds a request to the exemplar using a lock-free compare-and-swap (CAS) loop.
// This ensures that updates to the priority queue are atomic and thread-safe.
// If the priority queue is not full, the new item is added.
// If the priority queue is full and the new item's latency is greater than the minimum latency in the queue,
// the new item replaces the minimum latency item.
func (e *Exemplar) Offer(latency time.Duration, requestID string) {
	newItem := &item{
		latency:   latency,
		requestID: requestID,
	}

	for {
		oldPQPtr := e.pq.Load()
		oldPQ := *oldPQPtr

		if len(oldPQ) < e.k {
			newPQ := make(priorityQueue, len(oldPQ), len(oldPQ)+1)
			copy(newPQ, oldPQ)
			heap.Push(&newPQ, newItem)
			if e.pq.CompareAndSwap(oldPQPtr, &newPQ) {
				return
			}
		} else if latency > oldPQ[0].latency {
			newPQ := make(priorityQueue, len(oldPQ), e.k)
			copy(newPQ, oldPQ)
			newPQ[0] = newItem
			heap.Fix(&newPQ, 0)
			if e.pq.CompareAndSwap(oldPQPtr, &newPQ) {
				return
			}
		} else {
			// The new item is not larger than the smallest in a full queue.
			return
		}
	}
}

// Snapshot returns a snapshot of the exemplars. It is lock-free and returns a copy of the current exemplars.
func (e *Exemplar) Snapshot() []*item {
	pqPtr := e.pq.Load()
	pq := *pqPtr
	items := make([]*item, len(pq))
	copy(items, pq)
	return items
}

// item is an item in the priority queue, representing a single request exemplar.
// It is unexported to encapsulate the implementation details of the priority queue.
type item struct {
	latency   time.Duration
	requestID string
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
	item := x.(*item)
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
