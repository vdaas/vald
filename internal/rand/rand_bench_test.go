// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package rand

import (
	"math"
	stdrand "math/rand"
	randv2 "math/rand/v2"
	"testing"
)

// =============================================================================
// Uint32 Benchmarks
// =============================================================================

func Benchmark_Uint32_MathRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = stdrand.Uint32()
		}
	})
}

func Benchmark_Uint32_MathRandV2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = randv2.Uint32()
		}
	})
}

func Benchmark_Uint32_InternalRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Uint32()
		}
	})
}

// =============================================================================
// Uint64 Benchmarks
// =============================================================================

func Benchmark_Uint64_MathRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = stdrand.Uint64()
		}
	})
}

func Benchmark_Uint64_MathRandV2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = randv2.Uint64()
		}
	})
}

func Benchmark_Uint64_InternalRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Uint64()
		}
	})
}

// =============================================================================
// Float32 Benchmarks
// =============================================================================

func Benchmark_Float32_MathRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = stdrand.Float32()
		}
	})
}

func Benchmark_Float32_MathRandV2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = randv2.Float32()
		}
	})
}

func Benchmark_Float32_InternalRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Float32()
		}
	})
}

// =============================================================================
// Float64 Benchmarks
// =============================================================================

func Benchmark_Float64_MathRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = stdrand.Float64()
		}
	})
}

func Benchmark_Float64_MathRandV2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = randv2.Float64()
		}
	})
}

func Benchmark_Float64_InternalRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Float64()
		}
	})
}

// =============================================================================
// LimitedUint32 Benchmarks (Range)
// =============================================================================

func Benchmark_LimitedUint32_MathRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// math/rand doesn't have a direct Uint32n, so we use Int63n or similar
			_ = stdrand.Int31n(100)
		}
	})
}

func Benchmark_LimitedUint32_MathRandV2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = randv2.N(uint32(100))
		}
	})
}

func Benchmark_LimitedUint32_InternalRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = LimitedUint32(100)
		}
	})
}

// =============================================================================
// LimitedUint64 Benchmarks (Range)
// =============================================================================

func Benchmark_LimitedUint64_MathRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = stdrand.Int63n(100)
		}
	})
}

func Benchmark_LimitedUint64_MathRandV2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = randv2.N(uint64(100))
		}
	})
}

func Benchmark_LimitedUint64_InternalRand(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = LimitedUint64(100)
		}
	})
}

// =============================================================================
// Internal Implementation Details
// =============================================================================

func Benchmark_rand_Uint32_Parallel(b *testing.B) {
	var x uint32 = 0
	r := &rng[uint32]{
		x: &x,
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = r.Value()
		}
	})
}

func Benchmark_rand_init(b *testing.B) {
	r := &rng[uint32]{}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand := r.init(); rand.x == nil {
				b.Errorf("r.init() returns invalid object: %#v", rand)
			}
		}
	})
}

func Benchmark_rand64_Uint64_Parallel(b *testing.B) {
	var x uint64 = 0
	r := &rng[uint64]{
		x: &x,
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = r.Value()
		}
	})
}

func Benchmark_rand64_init(b *testing.B) {
	r := &rng[uint64]{}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand := r.init(); rand.x == nil {
				b.Errorf("r.init() returns invalid object: %#v", rand)
			}
		}
	})
}

func BenchmarkLimitedUint32_MaxUint64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = LimitedUint32(math.MaxUint32)
		}
	})
}

func BenchmarkLimitedUint64_MaxUint64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = LimitedUint64(math.MaxUint64)
		}
	})
}
