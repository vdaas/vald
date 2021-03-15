package location

import (
	"testing"
	"time"
)

var loc *time.Location

func BenchmarkGMT(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GMT()
		}
	})
}

func BenchmarkGMT1(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var l *time.Location
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l = GMT()
		}
	})
	loc = l
}

func BenchmarkUTC(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			UTC()
		}
	})
}

func BenchmarkUTC1(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var l *time.Location
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l = UTC()
		}
	})
	loc = l
}

func BenchmarkJST(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			JST()
		}
	})
}

func BenchmarkJST1(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var l *time.Location
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l = JST()
		}
	})
	loc = l
}
