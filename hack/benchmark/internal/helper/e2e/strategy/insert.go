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
	b.Run("Insert", func(b *testing.B) {
		ids, train := dataset.IDs(), dataset.Train()

		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			isrt.do(ctx, b, c, ids[i%len(ids)], train[i%len(train)])
		}
		b.StopTimer()
	})
}

func (isrt *insert) runParallel(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("ParallelInsert", func(b *testing.B) {
		ids, train := dataset.IDs(), dataset.Train()

		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		b.RunParallel(func(p *testing.PB) {
			i := 0
			for p.Next() {
				isrt.do(ctx, b, c, ids[i%len(ids)], train[i%len(train)])
				i++
			}
		})
		b.StopTimer()
	})
}

func (isrt *insert) do(ctx context.Context, b *testing.B, c client.Client, id string, vector []float32) {
	if err := c.Insert(ctx, &client.ObjectVector{
		Id:     id,
		Vector: vector,
	}); err != nil {
		b.Error(err)
	}
}
