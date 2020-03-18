package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type createAndSaveIndex struct {
	poolSize uint32
	preStart func(context.Context, *testing.B, ngt.NGT, assets.Dataset)
}

func NewCreateAndSaveIndex(opts ...CreateAndSaveIndexOption) benchmark.Strategy {
	c := new(createAndSaveIndex)
	for _, opt := range append(defaultCreateAndSaveIndexOptions, opts...) {
		opt(c)
	}
	return c
}

func (c *createAndSaveIndex) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("CreateAndSaveIndex", func(bb *testing.B) {
		c.preStart(ctx, b, ngt, dataset)

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			if err := ngt.CreateAndSaveIndex(c.poolSize); err != nil {
				bb.Error(err)
			}
		}
		bb.StopTimer()
	})
}
