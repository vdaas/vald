package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type insertCommit struct {
	poolSize uint32
}

func NewInsertCommit(opts ...InsertCommitOption) benchmark.Strategy {
	isrt := new(insertCommit)
	for _, opt := range append(defaultInsertCommitOptions, opts...) {
		opt(isrt)
	}
	return isrt
}

func (ic *insertCommit) Run(ctx context.Context, b *testing.B, c interface{}, typ benchmark.Type, dataset assets.Dataset) {
	cnt := 0
	b.Run("InsertCommit", func(bb *testing.B) {
		switch typ {
		case benchmark.Float32:
			ic.float32(ctx, bb, c.(core.Core32), dataset, &cnt)
		case benchmark.Float64:
			ic.float64(ctx, bb, c.(core.Core64), dataset, &cnt)
		default:
			bb.Fatal("invalid data type")
		}
	})
}

func (ic *insertCommit) float32(ctx context.Context, b *testing.B, core core.Core32, dataset assets.Dataset, cnt *int) {
	train := dataset.Train()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.InsertCommit(train[*cnt%len(train)], ic.poolSize)
		if err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}

func (ic *insertCommit) float64(ctx context.Context, b *testing.B, core core.Core64, dataset assets.Dataset, cnt *int) {
	train := dataset.TrainAsFloat64()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.InsertCommit(train[*cnt%len(train)], ic.poolSize)
		if err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}
