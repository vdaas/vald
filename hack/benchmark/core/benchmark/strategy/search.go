package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type search struct{}

func NewSearch(opts ...SearchOption) benchmark.Strategy {
	s := new(search)
	for _, opt := range append(defaultSearchOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *search) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("Search", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {

		}
	})
}
