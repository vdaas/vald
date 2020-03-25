package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
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

func (s *search) Run(ctx context.Context, b *testing.B, c interface{}, typ benchmark.Type, dataset assets.Dataset) {
	cnt := 0
	b.Run("Search", func(bb *testing.B) {
		_, err := s.preStart(ctx, bb, c, dataset)
		if err != nil {
			b.Fatal(err)
		}

		switch typ {
		case benchmark.Float32:
			s.float32(ctx, bb, c.(core.Core32), dataset, &cnt)
		case benchmark.Float64:
			s.float64(ctx, bb, c.(core.Core64), dataset, &cnt)
		default:
			bb.Fatal("invalid data type")
		}
	})
}
func (s *search) float32(ctx context.Context, b *testing.B, core core.Core32, dataset assets.Dataset, cnt *int) {
	query := dataset.Query()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.Search(query[*cnt%len(query)], s.size, s.epsilon, s.radius)
		if err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}

func (s *search) float64(ctx context.Context, b *testing.B, core core.Core64, dataset assets.Dataset, cnt *int) {
	query := dataset.QueryAsFloat64()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.Search(query[*cnt%len(query)], s.size, s.epsilon, s.radius)
		if err != nil {
			b.Error(err)
		}
		*cnt++
	}
	b.StopTimer()
}
