package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type bulkInsertCommit struct {
	poolSize uint32
}

func NewBulkInsertCommit(opts ...BulkInsertCommitOption) benchmark.Strategy {
	bi := new(bulkInsertCommit)
	for _, opt := range append(defaultBulkInsertCommitOptions, opts...) {
		opt(bi)
	}
	return bi
}

func (bi *bulkInsertCommit) Run(ctx context.Context, b *testing.B, c interface{}, typ benchmark.Type, dataset assets.Dataset) {
	cnt := 0
	b.Run("BulkInsertCommit", func(bb *testing.B) {
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

func (bi *bulkInsertCommit) float32(ctx context.Context, b *testing.B, core core.Core32, dataset assets.Dataset, cnt *int) {
	train := dataset.Train()

	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, errs := core.BulkInsertCommit(train, bi.poolSize)
		if err := wrapErrors(errs); err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}

func (bi *bulkInsertCommit) float64(ctx context.Context, b *testing.B, core core.Core64, dataset assets.Dataset, cnt *int) {
	train := dataset.TrainAsFloat64()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, errs := core.BulkInsertCommit(train, bi.poolSize)
		if err := wrapErrors(errs); err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}
