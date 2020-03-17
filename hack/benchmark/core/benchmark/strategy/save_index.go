package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type saveIndex struct{}

func NewSaveIndex(opts ...SaveIndexOption) benchmark.Strategy {
	s := new(saveIndex)
	for _, opt := range append(defaultSaveIndexOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *saveIndex) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("SaveIndex", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {

		}
	})
}
