package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type remove struct {
	preStart PreStart
}

func NewRemove(opts ...RemoveOption) benchmark.Strategy {
	r := new(remove)
	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}
	return r
}

func (d *remove) Run(ctx context.Context, b *testing.B, c interface{}, typ benchmark.Type, dataset assets.Dataset) {
	cnt := 0
	b.Run("Remove", func(bb *testing.B) {
		obj, err := d.preStart(ctx, b, c, dataset)
		if err != nil {
			b.Error(err)
		}

		ids := obj.([]uint)

		switch typ {
		case benchmark.Float32:
			d.float32(ctx, bb, c.(core.Core32), ids, &cnt)
		case benchmark.Float64:
			d.float64(ctx, bb, c.(core.Core64), ids, &cnt)
		default:
			bb.Fatal("invalid data type")
		}
	})
}

func (d *remove) float32(ctx context.Context, b *testing.B, core core.Core32, ids []uint, cnt *int) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		core.Remove(ids[*cnt%len(ids)])
		*cnt++
	}

	b.StopTimer()
}

func (d *remove) float64(ctx context.Context, b *testing.B, core core.Core64, ids []uint, cnt *int) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		core.Remove(ids[*cnt%len(ids)])
		*cnt++
	}
	b.StopTimer()
}
