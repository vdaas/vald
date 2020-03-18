package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type getVector struct {
	ids      []uint
	preStart PreStart
}

func NewGetVector(opts ...GetVectorOption) benchmark.Strategy {
	g := new(getVector)
	for _, opt := range append(defaultGetVectorOptions, opts...) {
		opt(g)
	}
	return g
}

func (g *getVector) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	cnt := 0
	b.Run("GetVector", func(bb *testing.B) {
		g.ids = append(g.ids, g.preStart(ctx, bb, ngt, dataset)...)

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			_, err := ngt.GetVector(g.ids[cnt%len(g.ids)])
			if err != nil {
				bb.Error(err)
			}
			cnt++
		}
		bb.StopTimer()
	})
}
