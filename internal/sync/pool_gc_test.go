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
	"runtime/debug"
	stdsync "sync"
	"testing"
)

func TestPoolChainedStdlibGC(t *testing.T) {
	// Not parallel: verifies our registered poolCleanup preserves stdlib sync.Pool draining.
	prevGC := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(prevGC)
	prevProcs := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(prevProcs)

	var gen Pool[string]
	var std stdsync.Pool
	gen.Put("hello")
	std.Put(&testPoolItem{val: 101})

	runtime.GC()
	if got, _ := std.Get().(*testPoolItem); got == nil || got.val != 101 {
		t.Fatalf("stdlib Get() after 1st GC = %+v, want val=101", got)
	}
	if got := gen.Get(); got != "hello" {
		t.Fatalf("generic Get() after 1st GC = %q, want %q", got, "hello")
	}

	gen.Put("hello2")
	std.Put(&testPoolItem{val: 202})
	runtime.GC()
	runtime.GC()

	if got := std.Get(); got != nil {
		t.Fatalf("stdlib Get() after 2 GCs = %v, want nil", got)
	}
	if got := gen.Get(); got != "" {
		t.Fatalf("generic Get() after 2 GCs = %q, want empty", got)
	}
}

func TestPoolAllocs(t *testing.T) {
	var ptrPool Pool[*testPoolItem]
	item := &testPoolItem{val: 1}
	ptrPool.Put(item)
	if allocs := testing.AllocsPerRun(1000, func() {
		v := ptrPool.Get()
		ptrPool.Put(v)
	}); allocs != 0 {
		t.Fatalf("Pool[*testPoolItem] allocs = %v, want 0", allocs)
	}

	// Unlike sync.Pool (which boxes non-pointer values into any), Pool[value] must not allocate.
	var valPool Pool[testValueStruct]
	valPool.Put(testValueStruct{a: 1, b: 2, c: 3, d: 4})
	if allocs := testing.AllocsPerRun(1000, func() {
		v := valPool.Get()
		valPool.Put(v)
	}); allocs != 0 {
		t.Fatalf("Pool[testValueStruct] allocs = %v, want 0", allocs)
	}
}

func BenchmarkPool(b *testing.B) {
	b.Run("GenericPool/Pointer/Serial", func(b *testing.B) {
		p := Pool[*testPoolItem]{
			New: func() *testPoolItem { return new(testPoolItem) },
		}
		p.Put(new(testPoolItem))
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			v := p.Get()
			p.Put(v)
		}
	})

	b.Run("StdlibPool/Pointer/Serial", func(b *testing.B) {
		p := stdsync.Pool{
			New: func() any { return new(testPoolItem) },
		}
		p.Put(new(testPoolItem))
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			v, _ := p.Get().(*testPoolItem)
			p.Put(v)
		}
	})

	b.Run("GenericPool/Pointer/Parallel", func(b *testing.B) {
		p := Pool[*testPoolItem]{
			New: func() *testPoolItem { return new(testPoolItem) },
		}
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				v := p.Get()
				p.Put(v)
			}
		})
	})

	b.Run("StdlibPool/Pointer/Parallel", func(b *testing.B) {
		p := stdsync.Pool{
			New: func() any { return new(testPoolItem) },
		}
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				v, _ := p.Get().(*testPoolItem)
				p.Put(v)
			}
		})
	})

	b.Run("GenericPool/ValueStruct/Serial", func(b *testing.B) {
		var p Pool[testValueStruct]
		p.Put(testValueStruct{a: 1, b: 2, c: 3, d: 4})
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			v := p.Get()
			p.Put(v)
		}
	})

	b.Run("StdlibPool/ValueStruct/Serial", func(b *testing.B) {
		var p stdsync.Pool
		//nolint:staticcheck // intentionally benchmarks interface boxing allocation in stdlib sync.Pool.
		p.Put(testValueStruct{a: 1, b: 2, c: 3, d: 4})
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			v, _ := p.Get().(testValueStruct)
			//nolint:staticcheck // intentionally benchmarks interface boxing allocation in stdlib sync.Pool.
			p.Put(v)
		}
	})

	b.Run("GenericPool/Burst16/Serial", func(b *testing.B) {
		p := Pool[*testPoolItem]{
			New: func() *testPoolItem { return new(testPoolItem) },
		}
		var items [16]*testPoolItem
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			for i := range items {
				items[i] = p.Get()
			}
			for i := range items {
				p.Put(items[i])
			}
		}
	})

	b.Run("StdlibPool/Burst16/Serial", func(b *testing.B) {
		p := stdsync.Pool{
			New: func() any { return new(testPoolItem) },
		}
		var items [16]*testPoolItem
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			for i := range items {
				items[i], _ = p.Get().(*testPoolItem)
			}
			for i := range items {
				p.Put(items[i])
			}
		}
	})

	b.Run("GenericPool/Burst16/Parallel", func(b *testing.B) {
		p := Pool[*testPoolItem]{
			New: func() *testPoolItem { return new(testPoolItem) },
		}
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var items [16]*testPoolItem
			for pb.Next() {
				for i := range items {
					items[i] = p.Get()
				}
				for i := range items {
					p.Put(items[i])
				}
			}
		})
	})

	b.Run("StdlibPool/Burst16/Parallel", func(b *testing.B) {
		p := stdsync.Pool{
			New: func() any { return new(testPoolItem) },
		}
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var items [16]*testPoolItem
			for pb.Next() {
				for i := range items {
					items[i], _ = p.Get().(*testPoolItem)
				}
				for i := range items {
					p.Put(items[i])
				}
			}
		})
	})

	b.Run("GenericPool/Slice/Serial", func(b *testing.B) {
		p := Pool[[]float32]{
			New: func() []float32 { return make([]float32, 0, 128) },
		}
		p.Put(make([]float32, 0, 128))
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			s := p.Get()
			p.Put(s)
		}
	})

	b.Run("StdlibPool/Slice/Serial", func(b *testing.B) {
		p := stdsync.Pool{
			New: func() any { return make([]float32, 0, 128) },
		}
		//nolint:staticcheck // intentionally benchmarks interface boxing allocation in stdlib sync.Pool.
		p.Put(make([]float32, 0, 128))
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			s, _ := p.Get().([]float32)
			//nolint:staticcheck // intentionally benchmarks interface boxing allocation in stdlib sync.Pool.
			p.Put(s)
		}
	})
}
