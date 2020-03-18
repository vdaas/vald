package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type saveIndex struct {
	poolSize uint32
	preStart PreStart
}

func NewSaveIndex(opts ...SaveIndexOption) benchmark.Strategy {
	s := new(saveIndex)
	for _, opt := range append(defaultSaveIndexOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *saveIndex) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("SaveIndex", func(bb *testing.B) {
		s.preStart(ctx, bb, ngt, dataset)

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			if err := ngt.CreateIndex(s.poolSize); err != nil {
				bb.Error(err)
			}
		}
		bb.StopTimer()
	})
}
