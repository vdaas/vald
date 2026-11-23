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

//go:build race

package sync

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// Pool is a set of temporary objects that may be individually saved and
// retrieved, specialized generically on T with zero interface boxing or
// reflection.
//
// Under -race, Pool synchronizes accesses with a mutex so TSAN observes exact
// happens-before edges across Put and Get while preserving two-GC victim aging.
//
// A Pool is safe for use by multiple goroutines simultaneously and must not be
// copied after first use.
type Pool[T any] struct {
	New    func() T
	header poolHeader
}

type racePoolStore[T any] struct {
	items []T
	mu    sync.Mutex
}

// Put adds x to the pool.
func (p *Pool[T]) Put(x T) {
	if p == nil {
		return
	}
	s := p.ensureLocal()
	s.mu.Lock()
	s.items = append(s.items, x)
	s.mu.Unlock()
}

// Get selects an arbitrary item from the Pool, removes it from the Pool, and
// returns it to the caller.
func (p *Pool[T]) Get() T {
	if p == nil {
		var zero T
		return zero
	}
	if s := (*racePoolStore[T])(atomic.LoadPointer(&p.header.local)); s != nil {
		if x, ok := s.pop(); ok {
			return x
		}
	}
	if s := (*racePoolStore[T])(atomic.LoadPointer(&p.header.victim)); s != nil {
		if x, ok := s.pop(); ok {
			return x
		}
	}
	if p.New != nil {
		return p.New()
	}
	var zero T
	return zero
}

func (p *Pool[T]) ensureLocal() *racePoolStore[T] {
	lp := atomic.LoadPointer(&p.header.local)
	if lp != nil {
		return (*racePoolStore[T])(lp)
	}
	allPoolsMu.Lock()
	defer allPoolsMu.Unlock()
	lp = atomic.LoadPointer(&p.header.local)
	if lp != nil {
		return (*racePoolStore[T])(lp)
	}
	store := new(racePoolStore[T])
	atomic.StorePointer(&p.header.local, unsafe.Pointer(store))
	registerPool(&p.header)
	return store
}

func (s *racePoolStore[T]) pop() (T, bool) {
	s.mu.Lock()
	n := len(s.items)
	if n == 0 {
		s.mu.Unlock()
		var zero T
		return zero, false
	}
	idx := n - 1
	x := s.items[idx]
	var zero T
	s.items[idx] = zero
	s.items = s.items[:idx]
	s.mu.Unlock()
	return x, true
}
