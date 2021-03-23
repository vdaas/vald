//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
package rand

import (
	"math"
	"testing"
)

func BenchmarkUint32(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Uint32()
		}
	})
}

func BenchmarkLimitedUint32_0(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			LimitedUint32(0)
		}
	})
}

func BenchmarkLimitedUint32_10(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			LimitedUint32(10)
		}
	})
}

func BenchmarkLimitedUint32_100(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			LimitedUint32(100)
		}
	})
}

func BenchmarkLimitedUint32_MaxUint64(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			LimitedUint32(math.MaxUint64)
		}
	})
}

func Benchmark_rand_Uint32_0(b *testing.B) {
	var x uint32 = 0
	r := &rand{
		x: &x,
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_Uint32_10(b *testing.B) {
	var x uint32 = 10
	r := &rand{
		x: &x,
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_Uint32_100(b *testing.B) {
	var x uint32 = 100
	r := &rand{
		x: &x,
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_Uint32_MaxUint32(b *testing.B) {
	var x uint32 = math.MaxUint32
	r := &rand{
		x: &x,
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_init_0(b *testing.B) {
	r := &rand{}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.init()
		}
	})
}

func Benchmark_rand_init_10(b *testing.B) {
	var x uint32 = 100
	r := &rand{
		x: &x,
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_init_100(b *testing.B) {
	var x uint32 = 100
	r := &rand{
		x: &x,
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_init_MaxUint32(b *testing.B) {
	var x uint32 = math.MaxUint32
	r := &rand{
		x: &x,
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}
