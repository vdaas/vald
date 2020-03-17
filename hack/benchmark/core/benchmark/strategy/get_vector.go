package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type getVector struct{}

func NewGetVector(opts ...GetVectorOption) benchmark.Strategy {
	g := new(getVector)
	for _, opt := range append(defaultGetVectorOptions, opts...) {
		opt(g)
	}
	return g
}

func (g *getVector) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("GetVector", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {

		}
	})
}
