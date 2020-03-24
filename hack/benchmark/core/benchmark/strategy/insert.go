package stratedy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type insert struct{}

func NewInsert(opts ...InsertOption) benchmark.Strategy {
	isrt := new(insert)
	for _, opt := range append(defaultInsertOptions, opts...) {
		opt(isrt)
	}
	return isrt
}

func (isrt *insert) Run(ctx context.Context, b *testing.B, c interface{}, typ benchmark.Type, dataset assets.Dataset) {
	cnt := 0
	b.Run("Insert", func(bb *testing.B) {
		switch typ {
		case benchmark.Float32:
			isrt.float32(ctx, bb, c.(core.Core32), dataset, &cnt)
		case benchmark.Float64:
			isrt.float64(ctx, bb, c.(core.Core64), dataset, &cnt)
		default:
			bb.Fatal("invalid data type")
		}
	})
}

func (isrt *insert) float32(ctx context.Context, b *testing.B, core core.Core32, dataset assets.Dataset, cnt *int) {
	train := dataset.Train()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		core.Insert(train[*cnt%len(train)])
		*cnt++
	}
	b.StopTimer()
}

func (isrt *insert) float64(ctx context.Context, b *testing.B, core core.Core64, dataset assets.Dataset, cnt *int) {
	train := dataset.TrainAsFloat64()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		core.Insert(train[*cnt%len(train)])
		*cnt++
	}
	b.StopTimer()
}
