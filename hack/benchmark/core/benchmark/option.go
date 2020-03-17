package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

// Option is benchmark configure.
type Option func(*benchmark)

var (
	defaultOptions = []Option{
		WithPreStart(
			func(context.Context, *testing.B, assets.Dataset, ngt.NGT) {},
		),
	}
)

// WithName returns Option that sets name.
func WithName(name string) Option {
	return func(b *benchmark) {
		if len(name) != 0 {
			b.name = name
		}
	}
}

// WithPreStart returns Option that sets prestart.
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

// WithStrategy returns Option that sets benchmark strategy.
func WithStrategy() Option {
	return func(b *benchmark) {}
}
