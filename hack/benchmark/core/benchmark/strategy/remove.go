package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type remove struct {
	ids []uint
}

func NewRemove(opts ...RemoveOption) benchmark.Strategy {
	r := new(remove)
	for _, opt := range append(defaultRemoveOptions, opts...) {
		opt(r)
	}
	return r
}

func (r *remove) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("Remove", func(bb *testing.B) {
		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {

		}
		bb.StopTimer()
	})
}
