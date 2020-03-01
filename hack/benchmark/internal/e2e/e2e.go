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
	name          string
	strategies    []Strategy
	dataset       assets.Dataset
	client        client.Client
	serverStarter func(context.Context, testing.TB, assets.Dataset) func()
}

func New(b *testing.B, opts ...Option) Runner {
	e := new(e2e)
	for _, opt := range append(defaultOptions, opts...) {
		opt(e)
	}

	fn := assets.Data(e.name)
	if fn == nil {
		b.Fatalf("dataset provider is nil: %v", e.name)
	}

	e.dataset = fn(b)
	if e.dataset == nil {
		b.Fatalf("dataset is nil: %v", e.name)
	}

	return e
}

func (e *e2e) Run(ctx context.Context, b *testing.B) {
	func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		defer e.serverStarter(ctx, b, assets.Data(e.name)(b))()

		b.StopTimer()
		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()

		b.Run(e.name, func(bb *testing.B) {
			for _, strategy := range e.strategies {
				strategy.Run(ctx, bb, e.client, e.dataset)
			}
		})

		b.StopTimer()
	}()
}
