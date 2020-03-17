package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type createIndex struct{}

func NewCreateIndex(opts ...CreateIndexOption) benchmark.Strategy {
	c := new(createIndex)
	for _, opt := range append(defaultCreateIndexOptions, opts...) {
		opt(c)
	}
	return c
}

func (c *createIndex) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("CreateIndex", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {

		}
	})
}
