package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type remove struct {
	ids      []uint
	preStart PreStart
}

func NewRemove(opts ...RemoveOption) benchmark.Strategy {
	r := new(remove)
	for _, opt := range append(defaultRemoveOptions, opts...) {
		opt(r)
	}
	return r
}

func (r *remove) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	cnt := 0
	b.Run("Remove", func(bb *testing.B) {
		r.ids = append(r.ids, r.preStart(ctx, bb, ngt, dataset)...)

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			if err := ngt.Remove(r.ids[cnt%len(r.ids)]); err != nil {
				b.Error(err)
			}
			cnt++
		}
		bb.StopTimer()
	})
}
