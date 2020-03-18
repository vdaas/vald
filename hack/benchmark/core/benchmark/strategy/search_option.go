package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type SearchOption func(*search)

var (
	defaultSearchOptions = []SearchOption{
		WithSearchSize(10),
		WithSearchEpsilon(0.01),
		WithSearchRadius(-1),
		WithSearchPreStart(
			func(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) (ids []uint) {
				ids = (new(preStart)).Func(ctx, b, ngt, dataset)
				ngt.CreateIndex(1000)
				return
			},
		),
	}
)

func WithSearchPreStart(fn PreStart) SearchOption {
	return func(s *search) {
		if fn != nil {
			s.preStart = fn
		}
	}
}

func WithSearchSize(size int) SearchOption {
	return func(s *search) {
		if size > 0 {
			s.size = size
		}
	}
}

func WithSearchEpsilon(epsilon float32) SearchOption {
	return func(s *search) {
		s.epsilon = epsilon
	}
}

func WithSearchRadius(radius float32) SearchOption {
	return func(s *search) {
		s.radius = radius
	}
}
