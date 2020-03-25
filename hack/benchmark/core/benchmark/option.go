package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

type Option func(*benchmark)

var (
	defaultOptions = []Option{}
)

func WithName(name string) Option {
	return func(b *benchmark) {
		if len(name) != 0 {
			b.name = name
		}
	}
}

func WithFloat32(
	fn func(context.Context, *testing.B, assets.Dataset) (interface{}, func(), error),
	strategies ...Strategy,
) Option {
	return func(b *benchmark) {
		if len(strategies) != 0 {
			b.typ = Float32
			b.initializer = fn
			b.strategies = strategies
		}
	}
}

func WithFloat64(
	fn func(context.Context, *testing.B, assets.Dataset) (interface{}, func(), error),
	strategies ...Strategy,
) Option {
	return func(b *benchmark) {
		if len(strategies) != 0 {
			b.typ = Float64
			b.initializer = fn
			b.strategies = strategies
		}
	}
}
