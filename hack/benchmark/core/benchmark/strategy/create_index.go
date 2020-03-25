package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type createIndex struct {
	poolSize uint32
	preStart PreStart
}

func NewCreateIndex(opts ...CreateIndexOption) benchmark.Strategy {
	ci := new(createIndex)
	for _, opt := range append(defaultCreateIndexOptions, opts...) {
		opt(ci)
	}
	return ci
}

func (ci *createIndex) Run(ctx context.Context, b *testing.B, c interface{}, typ benchmark.Type, dataset assets.Dataset) {
	b.Run("CreateIndex", func(bb *testing.B) {
		switch typ {
		case benchmark.Float32:
			ci.float32(ctx, bb, c.(core.Core32), dataset)
		case benchmark.Float64:
			ci.float64(ctx, bb, c.(core.Core64), dataset)
		default:
			bb.Fatal("invalid data type")
		}
	})
}

func (ci *createIndex) float32(ctx context.Context, b *testing.B, core core.Core32, dataset assets.Dataset) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := core.CreateIndex(ci.poolSize)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func (ci *createIndex) float64(ctx context.Context, b *testing.B, core core.Core64, dataset assets.Dataset) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := core.CreateIndex(ci.poolSize)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}
