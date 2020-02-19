package insert

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/helper/e2e"
	"github.com/vdaas/vald/internal/client"
)

const (
	serial   = "Insert Serial"
	parallel = "Insert Parallel"
)

type insert struct {
	parallel bool
}

func New(opts ...Option) e2e.Strategy {
	i := new(insert)

	for _, opt := range append(defaultOption, opts...) {
		opt(i)
	}

	return i
}

func (isrt *insert) Run(ctx context.Context, b *testing.B, client client.Client, dataset assets.Dataset) error {
	if isrt.parallel {
		return isrt.runParallel(ctx, b, client, dataset)
	}
	return isrt.run(ctx, b, client, dataset)
}

func (isrt *insert) run(ctx context.Context, b *testing.B, client client.Client, dataset assets.Dataset) error {
	b.Run(serial, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := client.Insert(ctx, nil); err != nil {

			}
		}
	})

	return nil
}

func (isrt *insert) runParallel(ctx context.Context, b *testing.B, client client.Client, dataset assets.Dataset) error {
	b.Run(parallel, func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				if err := client.Insert(ctx, nil); err != nil {
				}
			}
		})
	})
	return nil
}
