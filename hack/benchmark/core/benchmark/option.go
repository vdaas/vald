package benchmark

import (
	"github.com/vdaas/vald/internal/core/ngt"
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
