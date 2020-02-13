package helper

import (
	"testing"
)

type OperationHelper interface {
	Insert() func(b *testing.B)
	InsertParallel() func(b *testing.B)
	CreateIndex() func(b *testing.B)
	CreateIndexParallel() func(b *testing.B)
	Search() func(b *testing.B)
	SearchParallel() func(b *testing.B)
}

type operationHelper struct {
	// gongt instance
	// ngt instance
	// dataset
}

func NewOperationHelper(opts ...OperationHelperOption) OperationHelper {
	oh := new(operationHelper)

	for _, opt := range append(defaultOperationHelperOption, opts...) {
		opt(oh)
	}

	return oh
}

func (oh *operationHelper) Insert() func(b *testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		for i := 0; i < b.N; i++ {

		}
		b.StopTimer()
	}
}

func (oh *operationHelper) InsertParallel() func(b *testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
			}
		})
		b.StopTimer()
	}
}

func (oh *operationHelper) CreateIndex() func(b *testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		for i := 0; i < b.N; i++ {

		}
		b.StopTimer()
	}
}

func (oh *operationHelper) CreateIndexParallel() func(b *testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
			}
		})
		b.StopTimer()
	}
}

func (oh *operationHelper) Search() func(b *testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		for i := 0; i < b.N; i++ {

		}
		b.StopTimer()
	}
}

func (oh *operationHelper) SearchParallel() func(b *testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
			}
		})
		b.StopTimer()
	}
}
