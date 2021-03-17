package rand

import (
	"math"
	"testing"
)

func Benchmark_Uint32(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Uint32()
		}
	})
}

func Benchmark_LimitedUint32_0(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			LimitedUint32(0)
		}
	})
}

func Benchmark_LimitedUint32_10(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			LimitedUint32(10)
		}
	})
}

func Benchmark_LimitedUint32_100(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			LimitedUint32(100)
		}
	})
}

func Benchmark_LimitedUint32_MaxUint64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			LimitedUint32(math.MaxUint64)
		}
	})
}

func Benchmark_rand_Uint32_0(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var x uint32 = 0
	r := &rand{
		x: &x,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_Uint32_10(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var x uint32 = 10
	r := &rand{
		x: &x,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_Uint32_100(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var x uint32 = 100
	r := &rand{
		x: &x,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_Uint32_MaxUint32(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var x uint32 = math.MaxUint32
	r := &rand{
		x: &x,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_init_0(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	r := &rand{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.init()
		}
	})
}

func Benchmark_rand_init_10(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var x uint32 = 100
	r := &rand{
		x: &x,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_init_100(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var x uint32 = 100
	r := &rand{
		x: &x,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}

func Benchmark_rand_init_MaxUint32(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var x uint32 = math.MaxUint32
	r := &rand{
		x: &x,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint32()
		}
	})
}
