// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package sync

import "sync/atomic"

const (
	// dequeueBits holds half of the 64-bit headTail field in poolDequeue.
	dequeueBits = 32

	// dequeueMask masks out one 32-bit head or tail index.
	dequeueMask = 1<<dequeueBits - 1

	// dequeueLimitDivisor bounds each backing ring below half the index space.
	dequeueLimitDivisor = 4

	// dequeueLimit is the maximum size of a single poolDequeue backing slice.
	dequeueLimit = (1 << dequeueBits) / dequeueLimitDivisor

	poolChainInitSize     = 8
	poolChainGrowthFactor = 2

	poolSlotEmpty    uint32 = 0
	poolSlotOccupied uint32 = 1
)

type poolSlot[T any] struct {
	val   T
	state uint32
}

// poolDequeue is a generic lock-free fixed-size single-producer, multi-consumer queue
// mirroring Go's sync.poolDequeue without interface boxing.
type poolDequeue[T any] struct {
	vals []poolSlot[T]
	// headTail packs a 32-bit head index (upper bits) and 32-bit tail index (lower bits).
	headTail atomic.Uint64
}

// poolChain is a dynamically-sized lock-free queue backed by a doubly-linked list of rings.
type poolChain[T any] struct {
	head *poolChainElt[T]
	tail atomic.Pointer[poolChainElt[T]]
}

type poolChainElt[T any] struct {
	next atomic.Pointer[poolChainElt[T]]
	prev atomic.Pointer[poolChainElt[T]]
	poolDequeue[T]
}

func (d *poolDequeue[T]) cap32() uint32 {
	//nolint:gosec // len(d.vals) is a power of two bounded by dequeueLimit (1<<30).
	return uint32(len(d.vals))
}

func (d *poolDequeue[T]) unpack(ptrs uint64) (head, tail uint32) {
	return uint32((ptrs >> dequeueBits) & dequeueMask), uint32(ptrs & dequeueMask)
}

func (d *poolDequeue[T]) pack(head, tail uint32) uint64 {
	return (uint64(head) << dequeueBits) | uint64(tail&dequeueMask)
}

func (d *poolDequeue[T]) pushHead(val T) bool {
	ptrs := d.headTail.Load()
	head, tail := d.unpack(ptrs)
	n := d.cap32()
	if (tail+n)&dequeueMask == head {
		// Queue is full.
		return false
	}
	slot := &d.vals[head&(n-1)]
	// Ensure any concurrent popTail has finished clearing this slot.
	if atomic.LoadUint32(&slot.state) != poolSlotEmpty {
		return false
	}
	// Single-producer slot write: headTail.Add below provides the release barrier.
	slot.val = val
	slot.state = poolSlotOccupied
	d.headTail.Add(1 << dequeueBits)
	return true
}

func (d *poolDequeue[T]) popHead() (T, bool) {
	var (
		ptrs uint64
		head uint32
		tail uint32
	)
	for {
		ptrs = d.headTail.Load()
		head, tail = d.unpack(ptrs)
		if tail == head {
			var zero T
			return zero, false
		}
		head--
		ptrs2 := d.pack(head, tail)
		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
			slot, val := d.takeSlot(head)
			// Owner P is the sole producer, so plain store releases this slot without LOCK overhead.
			slot.state = poolSlotEmpty
			return val, true
		}
	}
}

func (d *poolDequeue[T]) popTail() (T, bool) {
	var (
		ptrs uint64
		head uint32
		tail uint32
	)
	for {
		ptrs = d.headTail.Load()
		head, tail = d.unpack(ptrs)
		if tail == head {
			var zero T
			return zero, false
		}
		ptrs2 := d.pack(head, tail+1)
		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
			slot, val := d.takeSlot(tail)
			// Atomic release synchronizes with pushHead's atomic.LoadUint32 on slot.state.
			atomic.StoreUint32(&slot.state, poolSlotEmpty)
			return val, true
		}
	}
}

func (d *poolDequeue[T]) takeSlot(idx uint32) (*poolSlot[T], T) {
	slot := &d.vals[idx&(d.cap32()-1)]
	val := slot.val
	var zero T
	slot.val = zero
	return slot, val
}

func (c *poolChain[T]) pushHead(val T) {
	d := c.head
	if d != nil && d.pushHead(val) {
		return
	}
	c.pushHeadGrow(d, val)
}

func (c *poolChain[T]) pushHeadGrow(d *poolChainElt[T], val T) {
	if d == nil {
		d = new(poolChainElt[T])
		d.vals = make([]poolSlot[T], poolChainInitSize)
		c.head = d
		c.tail.Store(d)
		if d.pushHead(val) {
			return
		}
	}

	newSize := min(len(d.vals)*poolChainGrowthFactor, dequeueLimit)
	d2 := new(poolChainElt[T])
	d2.prev.Store(d)
	d2.vals = make([]poolSlot[T], newSize)
	c.head = d2
	d.next.Store(d2)
	d2.pushHead(val)
}

func (c *poolChain[T]) popHead() (T, bool) {
	d := c.head
	for d != nil {
		if val, ok := d.popHead(); ok {
			return val, true
		}
		d = d.prev.Load()
	}
	var zero T
	return zero, false
}

func (c *poolChain[T]) popTail() (T, bool) {
	d := c.tail.Load()
	if d == nil {
		var zero T
		return zero, false
	}

	for {
		// Load next before popping tail so non-nil d2 proves d is permanently frozen.
		d2 := d.next.Load()

		if val, ok := d.popTail(); ok {
			return val, true
		}

		if d2 == nil {
			var zero T
			return zero, false
		}

		if c.tail.CompareAndSwap(d, d2) {
			d2.prev.Store(nil)
		}
		d = d2
	}
}
