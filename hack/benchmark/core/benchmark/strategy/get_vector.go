package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
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

func (g *getVector) Run(ctx context.Context, b *testing.B, c interface{}, typ benchmark.Type, dataset assets.Dataset) {
	cnt := 0
	b.Run("GetVector", func(bb *testing.B) {
		obj, err := g.preStart(ctx, bb, c, dataset)
		if err != nil {
			b.Fatal(err)
		}

		g.ids = append(g.ids, obj.([]uint)...)

		switch typ {
		case benchmark.Float32:
			g.float32(ctx, bb, c.(core.Core32), dataset, &cnt)
		case benchmark.Float64:
			g.float64(ctx, bb, c.(core.Core64), dataset, &cnt)
		default:
			bb.Fatal("invalid data type")
		}
	})
}
func (g *getVector) float32(ctx context.Context, b *testing.B, core core.Core32, dataset assets.Dataset, cnt *int) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.GetVector(g.ids[*cnt%len(g.ids)])
		if err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}

func (g *getVector) float64(ctx context.Context, b *testing.B, core core.Core64, dataset assets.Dataset, cnt *int) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.GetVector(g.ids[*cnt%len(g.ids)])
		if err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}
