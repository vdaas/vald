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
	"sync"
	"time"
)

// Exemplar holds a sample of high-latency requests.
type Exemplar struct {
	mu sync.Mutex
	pq priorityQueue
	k  int
}

// NewExemplar creates a new Exemplar with a capacity of k.
func NewExemplar(k int) *Exemplar {
	k = max(k, 1)
	return &Exemplar{
		pq: make(priorityQueue, 0, k),
		k:  k,
	}
}

// Offer adds a request to the exemplar.
func (e *Exemplar) Offer(latency time.Duration, requestID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.pq) < e.k {
		heap.Push(&e.pq, &item{
			latency:   latency,
			requestID: requestID,
		})
	} else if latency > e.pq[0].latency {
		heap.Pop(&e.pq)
		heap.Push(&e.pq, &item{
			latency:   latency,
			requestID: requestID,
		})
	}
}

// Snapshot returns a snapshot of the exemplars.
func (e *Exemplar) Snapshot() []*item {
	e.mu.Lock()
	defer e.mu.Unlock()

	items := make([]*item, len(e.pq))
	copy(items, e.pq)
	return items
}

// item is an item in the priority queue.
type item struct {
	latency   time.Duration
	requestID string
}

// priorityQueue implements heap.Interface.
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
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
