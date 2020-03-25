package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

type SearchOption func(*search)

var (
	defaultSearchOptions = []SearchOption{
		WithSearchSize(10),
		WithSearchEpsilon(0.01),
		WithSearchRadius(-1),
		WithSearchPreStart(
			func(ctx context.Context, b *testing.B, c interface{}, dataset assets.Dataset) (interface{}, error) {
				ids, err := (new(defaultInsert)).PreStart(ctx, b, c, dataset)
				if err != nil {
					return nil, err
				}

				_, err = (new(defaultCreateIndex)).PreStart(ctx, b, c, dataset)
				if err != nil {
					return nil, err
				}

				return ids, nil
			},
		),
	}
)

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

func WithSearchPreStart(fn PreStart) SearchOption {
	return func(s *search) {
		if fn != nil {
			s.preStart = fn
		}
	}
}
