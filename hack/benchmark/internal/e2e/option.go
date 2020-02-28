package e2e

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/client"
)

type Option func(*e2e)

var (
	defaultOptions = []Option{}
)

func WithName(name string) Option {
	return func(e *e2e) {
		if len(name) != 0 {
			e.name = name
		}
	}
}

func WithClient(c client.Client) Option {
	return func(e *e2e) {
		if c != nil {
			e.client = c
		}
	}
}

func WithStrategy(strategis ...Strategy) Option {
	return func(e *e2e) {
		if len(strategis) != 0 {
			e.strategies = strategis
		}
	}
}

func WithServerStarter(f func(context.Context, testing.TB, assets.Dataset) func()) Option {
	return func(e *e2e) {
		if f != nil {
			e.serverStarter = f
		}
	}
}
