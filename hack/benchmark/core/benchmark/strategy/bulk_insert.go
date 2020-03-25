package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type bulkInsert struct{}

func NewBulkInsert(opts ...BulkInsertOption) benchmark.Strategy {
	bi := new(bulkInsert)
	for _, opt := range append(defaultBulkInsertOptions, opts...) {
		opt(bi)
	}
	return bi
}

func (bi *bulkInsert) Run(ctx context.Context, b *testing.B, c interface{}, typ benchmark.Type, dataset assets.Dataset) {
	cnt := 0
	b.Run("BulkInsert", func(bb *testing.B) {
		switch typ {
		case benchmark.Float32:
			bi.float32(ctx, bb, c.(core.Core32), dataset, &cnt)
		case benchmark.Float64:
			bi.float64(ctx, bb, c.(core.Core64), dataset, &cnt)
		default:
			bb.Fatal("invalid data type")
		}
	})
}

func (bi *bulkInsert) float32(ctx context.Context, b *testing.B, core core.Core32, dataset assets.Dataset, cnt *int) {
	train := dataset.Train()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, errs := core.BulkInsert(train)
		if err := wrapErrors(errs); err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}

func (bi *bulkInsert) float64(ctx context.Context, b *testing.B, core core.Core64, dataset assets.Dataset, cnt *int) {
	train := dataset.TrainAsFloat64()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, errs := core.BulkInsert(train)
		if err := wrapErrors(errs); err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}
