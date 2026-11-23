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

//go:build !race

package sync

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

const (
	// poolLocalPadSize sizes the padding tail so common pointer, map, and
	// scalar instantiations align poolLocal[T] to an exact 128-byte stride
	// (power-of-two shift indexing with zero false sharing across 64B lines).
	poolLocalPadSize = 96
)

// Pool is a set of temporary objects that may be individually saved and
// retrieved, specialized generically on T with zero interface boxing or
// reflection.
//
// Unlike standard library sync.Pool, Pool[T] preserves zero values (including
// zero-valued structs and scalars); callers storing pointer types should avoid
// calling Put with nil.
//
// A Pool is safe for use by multiple goroutines simultaneously and must not be
// copied after first use.
type Pool[T any] struct {
	// New optionally specifies a function to generate a value when Get would
	// otherwise return the zero value of T. It may not be changed concurrently
	// with calls to Get.
	New func() T

	header poolHeader
}

type poolLocalInternal[T any] struct {
	shared     poolChain[T]
	private    T
	privateSet bool
}

type poolLocal[T any] struct {
	poolLocalInternal[T]
	_ [poolLocalPadSize]byte
}

// Put adds x to the pool.
func (p *Pool[T]) Put(x T) {
	if p == nil {
		return
	}
	pid := runtimeProcPin()
	s := atomic.LoadUintptr(&p.header.localSize)
	locals := atomic.LoadPointer(&p.header.local)
	var l *poolLocal[T]
	if uintptr(pid) < s {
		l = indexLocal[T](locals, pid)
	} else {
		l, _ = p.pinSlow()
	}
	if !l.privateSet {
		l.private = x
		l.privateSet = true
	} else {
		l.shared.pushHead(x)
	}
	runtimeProcUnpin()
}

// Get selects an arbitrary item from the Pool, removes it from the Pool, and
// returns it to the caller.
func (p *Pool[T]) Get() T {
	if p == nil {
		var zero T
		return zero
	}
	pid := runtimeProcPin()
	s := atomic.LoadUintptr(&p.header.localSize)
	locals := atomic.LoadPointer(&p.header.local)
	var l *poolLocal[T]
	if uintptr(pid) < s {
		l = indexLocal[T](locals, pid)
	} else {
		l, pid = p.pinSlow()
	}
	if l.privateSet {
		x := l.private
		var zero T
		l.private = zero
		l.privateSet = false
		runtimeProcUnpin()
		return x
	}
	x, ok := l.shared.popHead()
	if !ok {
		x, ok = p.getSlow(pid)
	}
	runtimeProcUnpin()
	if !ok && p.New != nil {
		x = p.New()
	}
	return x
}

func (p *Pool[T]) getSlow(pid int) (T, bool) {
	size := atomic.LoadUintptr(&p.header.localSize)
	locals := atomic.LoadPointer(&p.header.local)
	// Try to steal one element from other procs.
	for i := range int(size) {
		l := indexLocal[T](locals, (pid+i+1)%int(size))
		if x, ok := l.shared.popTail(); ok {
			return x, true
		}
	}

	// Try the victim cache after primary caches so objects can age out on GC.
	size = atomic.LoadUintptr(&p.header.victimSize)
	if uintptr(pid) >= size {
		var zero T
		return zero, false
	}
	locals = atomic.LoadPointer(&p.header.victim)
	l := indexLocal[T](locals, pid)
	if l.privateSet {
		x := l.private
		var zero T
		l.private = zero
		l.privateSet = false
		return x, true
	}
	for i := range int(size) {
		v := indexLocal[T](locals, (pid+i)%int(size))
		if x, ok := v.shared.popTail(); ok {
			return x, true
		}
	}

	// Mark victim cache empty so subsequent misses skip scanning it until next GC.
	atomic.StoreUintptr(&p.header.victimSize, 0)

	var zero T
	return zero, false
}

func (p *Pool[T]) pinSlow() (*poolLocal[T], int) {
	// Mutex acquisition cannot happen while pinned to a P.
	runtimeProcUnpin()
	allPoolsMu.Lock()
	defer allPoolsMu.Unlock()
	pid := runtimeProcPin()
	// poolCleanup cannot run while pinned.
	s := atomic.LoadUintptr(&p.header.localSize)
	l := atomic.LoadPointer(&p.header.local)
	if uintptr(pid) < s {
		return indexLocal[T](l, pid), pid
	}
	size := max(runtime.GOMAXPROCS(0), pid+1)
	local := make([]poolLocal[T], size)
	//nolint:gosec // publishes backing array base pointer for per-P indexLocal arithmetic.
	atomic.StorePointer(&p.header.local, unsafe.Pointer(&local[0]))
	atomic.StoreUintptr(&p.header.localSize, uintptr(size))
	registerPool(&p.header)
	return &local[pid], pid
}

func indexLocal[T any](l unsafe.Pointer, i int) *poolLocal[T] {
	return (*poolLocal[T])(unsafe.Add(l, uintptr(i)*unsafe.Sizeof(poolLocal[T]{})))
}

//go:linkname runtimeProcPin runtime.procPin

func runtimeProcPin() int

//go:linkname runtimeProcUnpin runtime.procUnpin

func runtimeProcUnpin()
