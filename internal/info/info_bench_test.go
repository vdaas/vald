package info

import (
	"testing"
)

var getResult Detail

func BenchmarkString(b *testing.B) {
	Init("benchmark")
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if s := String(); s == "" {
				b.Error("String return empty string")
			}
		}
	})
}

func BenchmarkGet(b *testing.B) {
	sn := "benchmark"
	Init(sn)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if d := Get(); d.ServerName != sn {
				b.Errorf("Get server name is not match, result: %s", d.ServerName)
			}
		}
	})
}

func Benchmark_info_String(b *testing.B) {
	i, err := New(WithServerName("benchmark"))
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if s := i.String(); s == "" {
				b.Error("String return empty string")
			}
		}
	})
}

func Benchmark_Detail_String(b *testing.B) {
	i, err := New(WithServerName("benchmark"))
	if err != nil {
		b.Fatal(err)
	}
	d := i.Get()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if s := d.String(); s == "" {
				b.Error("String return empty string")
			}
		}
	})
}

func Benchmark_info_Get(b *testing.B) {
	sn := "benchmark"
	i, err := New(WithServerName(sn))
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if d := i.Get(); d.ServerName != sn {
				b.Errorf("Get server name is not match, result: %s", d.ServerName)
			}
		}
	})
}

func Benchmark_info_prepare(b *testing.B) {
	i, err := New(WithServerName("benchmark"))
	if err != nil {
		b.Fatal(err)
	}
	in := i.(*info)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			in.prepare()
		}
	})
}

func Benchmark_StackTrace_String(b *testing.B) {
	i, err := New(WithServerName("benchmark"))
	if err != nil {
		b.Fatal(err)
	}
	st := i.Get().StackTrace[0]

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if s := st.String(); s == "" {
				b.Error("String return empty string")
			}
		}
	})
}
