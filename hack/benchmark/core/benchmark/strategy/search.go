package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type search struct {
	size     int
	epsilon  float32
	radius   float32
	preStart PreStart
}

func NewSearch(opts ...SearchOption) benchmark.Strategy {
	s := new(search)
	for _, opt := range append(defaultSearchOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *search) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	cnt := 0
	b.Run("Search", func(bb *testing.B) {
		query := dataset.Query()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			_, err := ngt.Search(query[cnt%len(query)], s.size, s.epsilon, s.radius)
			if err != nil {
				bb.Error(err)
			}
			cnt++
		}
		bb.StopTimer()
	})
}
