package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/helper/e2e"
	"github.com/vdaas/vald/internal/client"
)

type insert struct {
	parallel bool
}

func NewInsert(opts ...InsertOption) e2e.Strategy {
	i := new(insert)

	for _, opt := range append(defaultInsertOption, opts...) {
		opt(i)
	}

	return i
}

func (isrt *insert) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	if isrt.parallel {
		isrt.runParallel(ctx, b, c, dataset)
		return
	}
	isrt.run(ctx, b, c, dataset)
}

func (isrt *insert) run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			isrt.do(ctx, b, c, dataset)
		}
	})
}

func (isrt *insert) runParallel(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				isrt.do(ctx, b, c, dataset)
			}
		})
	})
}

func (isrt *insert) do(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	if err := c.Insert(ctx, &client.ObjectVector{
		Vector: []float32{},
	}); err != nil {
		b.Error(err)
	}
}
