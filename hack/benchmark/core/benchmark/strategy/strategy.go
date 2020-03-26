package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type strategy struct {
	core32    core.Core32
	core64    core.Core64
	propName  string
	preProp32 func(context.Context, *testing.B, core.Core32, assets.Dataset) (interface{}, error)
	preProp64 func(context.Context, *testing.B, core.Core64, assets.Dataset) (interface{}, error)
	mode      core.Mode
	prop32    func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error)
	prop64    func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error)
}

func newStrategy(opts ...StrategyOption) *strategy {
	s := new(strategy)
	for _, opt := range append(defaultStrategyOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *strategy) PreProcess(ctx context.Context, b *testing.B, dataset assets.Dataset) (interface{}, error) {
	switch s.mode {
	case core.Float32:
		return s.preProp32(ctx, b, s.core32, dataset)
	case core.Float64:
		return s.preProp64(ctx, b, s.core64, dataset)
	default:
		b.Fatalf("invalid mode: %v", s.mode)
		return nil, nil
	}
}

func (s *strategy) Run(ctx context.Context, b *testing.B, dataset assets.Dataset) {
	var cnt uint64
	switch s.mode {
	case core.Float32:
		obj, err := s.preProp32(ctx, b, s.core32, dataset)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(s.propName, func(bb *testing.B) {
			s.float32(ctx, b, dataset, toUint(obj), &cnt)
		})
	case core.Float64:
		obj, err := s.preProp32(ctx, b, s.core32, dataset)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(s.propName, func(bb *testing.B) {
			s.float64(ctx, b, dataset, toUint(obj), &cnt)
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

func (s *strategy) float32(ctx context.Context, b *testing.B, dataset assets.Dataset, ids []uint, cnt *uint64) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.prop32(ctx, b, s.core32, dataset, ids, cnt)
		atomic.AddUint64(cnt, 1)
	}
	b.StopTimer()
}

func (s *strategy) float64(ctx context.Context, b *testing.B, dataset assets.Dataset, ids []uint, cnt *uint64) {
	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.prop64(ctx, b, s.core64, dataset, ids, cnt)
		atomic.AddUint64(cnt, 1)
	}
	b.StopTimer()
}
