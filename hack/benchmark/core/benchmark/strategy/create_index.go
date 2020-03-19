package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/gongt"
	"github.com/vdaas/vald/internal/core/ngt"
)

type createIndex struct {
	poolSize uint32
	preStart PreStart
}

func NewCreateIndex(opts ...CreateIndexOption) benchmark.Strategy {
	c := new(createIndex)
	for _, opt := range append(defaultCreateIndexOptions, opts...) {
		opt(c)
	}
	return c
}

func (c *createIndex) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, gongt gongt.NGT, dataset assets.Dataset) {
	b.Run("CreateIndex", func(bb *testing.B) {
		c.preStart(ctx, bb, ngt, dataset)

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			if err := ngt.CreateIndex(c.poolSize); err != nil {
				bb.Error(err)
			}
		}
		bb.StopTimer()
	})
}
