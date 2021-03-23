package location

import (
	"testing"
)

func BenchmarkGMT(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if l := GMT(); l == nil {
				b.Error("GMT return nil")
			}
		}
	})
}

func BenchmarkUTC(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if l := UTC(); l == nil {
				b.Error("UTC return nil")
			}
		}
	})
}

func BenchmarkJST(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if l := JST(); l == nil {
				b.Error("JST return nil")
			}
		}
	})
}
