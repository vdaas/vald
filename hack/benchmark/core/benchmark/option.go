package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type Option func(*benchmark)

var (
	defaultOptions = []Option{
		WithPreStart(
			func(context.Context, *testing.B, assets.Dataset, ngt.NGT) {},
		),
	}
)

func WithName(name string) Option {
	return func(b *benchmark) {
		if len(name) != 0 {
			b.name = name
		}
	}
}

func WithPreStart(
	f func(context.Context, *testing.B, assets.Dataset, ngt.NGT),
) Option {
	return func(b *benchmark) {
		if f != nil {
			b.prestart = f
		}
	}
}

func WithNGT(ngt ngt.NGT) Option {
	return func(b *benchmark) {
		if ngt != nil {
			b.ngt = ngt
		}
	}
}

func WithStrategy(strategies ...Strategy) Option {
	return func(b *benchmark) {
		if len(strategies) != 0 {
			b.strategies = strategies
		}
	}
}
