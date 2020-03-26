package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type strategy struct {
	initCore32 func(context.Context, assets.Dataset) (core.Core32, func(), error)
	initCore64 func(context.Context, assets.Dataset) (core.Core64, func(), error)
	propName   string
	preProp32  func(context.Context, *testing.B, core.Core32, assets.Dataset) (interface{}, error)
	preProp64  func(context.Context, *testing.B, core.Core64, assets.Dataset) (interface{}, error)
	mode       core.Mode
	prop32     func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error)
	prop64     func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error)
}

func newStrategy(opts ...StrategyOption) benchmark.Strategy {
	s := &strategy{
		// invalid mode.
		mode: core.Mode(100),
	}
	for _, opt := range append(defaultStrategyOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *strategy) Run(ctx context.Context, b *testing.B, dataset assets.Dataset) {
	var cnt uint64
	switch s.mode {
	case core.Float32:
		c, close, err := s.initCore32(ctx, dataset)
		if err != nil {
			b.Fatal(err)
		}
		defer close()

		b.Run(s.propName, func(bb *testing.B) {
			s.float32(ctx, bb, c, dataset, nil, &cnt)
		})
	case core.Float64:
		c, close, err := s.initCore64(ctx, dataset)
		if err != nil {
			b.Fatal(err)
		}
		defer close()

		b.Run(s.propName, func(bb *testing.B) {
			s.float64(ctx, bb, c, dataset, nil, &cnt)
		})
	default:
		b.Fatalf("invalid mode: %v", s.mode)
	}
}

func toUint(in interface{}) (out []uint) {
	if in != nil {
		var ok bool
		if out, ok = in.([]uint); ok {
			return
		}
	}
	return
}

func (s *strategy) float32(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset, ids []uint, cnt *uint64) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.prop32(ctx, b, c, dataset, ids, cnt)
		atomic.AddUint64(cnt, 1)
	}
	b.StopTimer()
}

func (s *strategy) float64(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset, ids []uint, cnt *uint64) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.prop64(ctx, b, c, dataset, ids, cnt)
		atomic.AddUint64(cnt, 1)
	}
	b.StopTimer()
}
