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

import (
	"runtime"
	"sync"
	"testing"
)

type testPoolItem struct {
	val int
}

type testValueStruct struct {
	a uint64
	b uint64
	c int32
	d int32
}

func TestPool(t *testing.T) {
	t.Parallel()

	t.Run("with allocator", func(t *testing.T) {
		t.Parallel()
		p := Pool[*testPoolItem]{
			New: func() *testPoolItem {
				return &testPoolItem{val: 7}
			},
		}
		item := p.Get()
		if item == nil || item.val != 7 {
			t.Fatalf("Get() = %+v, want val=7", item)
		}
		item.val = 99
		p.Put(item)
		reused := p.Get()
		if reused == nil || (reused.val != 99 && reused.val != 7) {
			t.Fatalf("Get() after Put() = %+v, want val 99 or 7", reused)
		}
	})

	t.Run("without allocator", func(t *testing.T) {
		t.Parallel()
		var p Pool[*testPoolItem]
		if got := p.Get(); got != nil {
			t.Fatalf("Get() on empty unconfigured pool = %+v, want nil", got)
		}
		p.Put(&testPoolItem{val: 42})
		if got := p.Get(); got != nil && got.val != 42 {
			t.Fatalf("Get() = %+v, want nil or val=42", got)
		}
	})

	t.Run("nil receiver safety", func(t *testing.T) {
		t.Parallel()
		var p *Pool[*testPoolItem]
		p.Put(&testPoolItem{val: 1})
		if got := p.Get(); got != nil {
			t.Fatalf("(*Pool)(nil).Get() = %+v, want nil", got)
		}
	})

	t.Run("zero value round-trip for value types", func(t *testing.T) {
		t.Parallel()
		prev := runtime.GOMAXPROCS(1)
		defer runtime.GOMAXPROCS(prev)

		var intPool Pool[int]
		intPool.New = func() int {
			return -1
		}
		// Seed and retrieve explicit zero without triggering New().
		intPool.Put(0)
		if got := intPool.Get(); got != 0 {
			t.Fatalf("Get() after Put(0) = %d, want 0", got)
		}
		if got := intPool.Get(); got != -1 {
			t.Fatalf("Get() on empty pool with New = %d, want -1", got)
		}

		// Exercise struct value types in both private and shared chain slots.
		var structPool Pool[testValueStruct]
		wantA := testValueStruct{}
		wantB := testValueStruct{a: 11, b: 22, c: 33, d: 44}
		structPool.Put(wantA)
		structPool.Put(wantB)
		got1 := structPool.Get()
		got2 := structPool.Get()
		if (got1 != wantA || got2 != wantB) && (got1 != wantB || got2 != wantA) {
			t.Fatalf("Got value structs (%+v, %+v), want (%+v, %+v)", got1, got2, wantA, wantB)
		}
	})

	t.Run("non-comparable reference types", func(t *testing.T) {
		t.Parallel()
		prev := runtime.GOMAXPROCS(1)
		defer runtime.GOMAXPROCS(prev)

		var slicePool Pool[[]byte]
		buf := make([]byte, 8)
		buf[0] = 0xAB
		slicePool.Put(buf)
		got := slicePool.Get()
		if len(got) != 8 || got[0] != 0xAB {
			t.Fatalf("Get() slice = %v, want len=8 first=0xAB", got)
		}
	})

	t.Run("poolChain ring growth and drain", func(t *testing.T) {
		t.Parallel()
		prev := runtime.GOMAXPROCS(1)
		defer runtime.GOMAXPROCS(prev)

		const n = 200
		var p Pool[int]
		for i := range n {
			p.Put(i + 1)
		}
		seen := make(map[int]bool, n)
		for range n {
			v := p.Get()
			if v == 0 || seen[v] {
				t.Fatalf("unexpected value popped from pool: %d", v)
			}
			seen[v] = true
		}
		if got := p.Get(); got != 0 {
			t.Fatalf("Get() after draining %d items = %d, want 0", n, got)
		}
	})

	t.Run("concurrent stress with GC", func(t *testing.T) {
		t.Parallel()
		p := Pool[*testPoolItem]{
			New: func() *testPoolItem {
				return &testPoolItem{val: 1}
			},
		}
		const workers = 16
		const iters = 500
		var wg sync.WaitGroup
		wg.Add(workers)
		for w := range workers {
			go func(id int) {
				defer wg.Done()
				for i := range iters {
					it := p.Get()
					if it == nil {
						t.Errorf("Get() returned nil")
						return
					}
					it.val = id*iters + i
					p.Put(it)
					if i%100 == 0 && id%4 == 0 {
						runtime.GC()
					}
				}
			}(w)
		}
		wg.Wait()
	})
}
