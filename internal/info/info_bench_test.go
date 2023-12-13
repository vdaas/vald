// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package info

import (
	"testing"
)

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

func BenchmarkDetail_String(b *testing.B) {
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

func BenchmarkStackTrace_String(b *testing.B) {
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
