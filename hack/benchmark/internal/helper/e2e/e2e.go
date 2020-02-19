package e2e

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/client"
)

type Runner interface {
	Run(ctx context.Context, b *testing.B)
}

type e2e struct {
	name       string
	strategies []Strategy
	dataset    assets.Dataset
	client     client.Client
}

func New(opts ...Option) Runner {
	e := new(e2e)
	for _, opt := range append(defaultOptions, opts...) {
		opt(e)
	}
	return e
}

func (e *e2e) Run(ctx context.Context, b *testing.B) {
	b.Run(e.name, func(b *testing.B) {
		for _, strategy := range e.strategies {
			strategy.Run(ctx, b, e.client, e.dataset)
		}
	})
}
