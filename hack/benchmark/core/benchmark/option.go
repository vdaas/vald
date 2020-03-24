package benchmark

import (
	"github.com/vdaas/vald/hack/benchmark/internal/core"
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

func WithCore(core core.Core) Option {
	return func(b *benchmark) {
		b.core = core
	}
}

func WithFloat32(core core.Core32, strategies ...Strategy) Option {
	return func(b *benchmark) {
		if len(strategies) != 0 {
			b.typ = Float32
			b.strategies = strategies
		}
	}
}

func WithFloat64(core core.Core32, strategies ...Strategy) Option {
	return func(b *benchmark) {
		if len(strategies) != 0 {
			b.typ = Float64
			b.strategies = strategies
		}
	}
}
