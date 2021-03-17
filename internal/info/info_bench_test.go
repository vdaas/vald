package info

import (
	"testing"
)

var stringResult string
var getResult Detail

func initBench() {
	Init("benchmark_test")
}

func BenchmarkString(b *testing.B) {
	initBench()

	b.ReportAllocs()
	b.ResetTimer()
	var r string
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r = String()
		}
	})
	stringResult = r
}

func BenchmarkGet(b *testing.B) {
	initBench()

	b.ReportAllocs()
	b.ResetTimer()
	var r Detail
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r = Get()
		}
	})
	getResult = r
}

func Benchmark_info_String(b *testing.B) {
	i, err := New(WithServerName("benchmark"))
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var r string
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r = i.String()
		}
	})
	stringResult = r
}

func Benchmark_Detail_String(b *testing.B) {
	i, err := New(WithServerName("benchmark"))
	if err != nil {
		b.Fatal(err)
	}
	d := i.Get()

	b.ReportAllocs()
	b.ResetTimer()
	var r string
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r = d.String()
		}
	})
	stringResult = r
}

func Benchmark_info_Get(b *testing.B) {
	i, err := New(WithServerName("benchmark"))
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var r Detail
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r = i.Get()
		}
	})
	getResult = r
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
	var s string
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s = st.String()
		}
	})
	stringResult = s
}
