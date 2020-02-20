package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/helper/e2e"
	"github.com/vdaas/vald/internal/client"
)

type search struct {
	parallel bool
	size     uint32
	epsilon  float32
	radius   float32
	cfg      *client.SearchConfig
}

func NewSearch(opts ...SearchOption) e2e.Strategy {
	s := new(search)

	for _, opt := range append(defaultSearchOptions, opts...) {
		opt(s)
	}

	s.cfg = &client.SearchConfig{
		Num:     s.size,
		Epsilon: s.epsilon,
		Radius:  s.radius,
	}
	return s
}

func (s *search) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	if s.parallel {
		s.runParallel(ctx, b, c, dataset)
		return
	}
	s.run(ctx, b, c, dataset)
}

func (s *search) run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("Search", func(bb *testing.B) {
		queries := dataset.Query()

		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		for i := 0; i < bb.N; i++ {
			s.do(ctx, b, c, queries[i%len(queries)])
		}
		b.StopTimer()
	})
}

func (s *search) runParallel(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("ParallelSearch", func(b *testing.B) {
		queries := dataset.Query()

		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				s.do(ctx, b, c, queries[i%len(queries)])
				i++
			}
		})
		b.StopTimer()
	})
}

func (s *search) do(ctx context.Context, b *testing.B, c client.Client, query []float32) {
	if _, err := c.Search(ctx, &client.SearchRequest{
		Vector: query,
		Config: s.cfg,
	}); err != nil {
		b.Error(err)
	}
}
