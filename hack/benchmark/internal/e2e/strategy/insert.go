package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
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
	cnt := 0
	b.Run("Insert", func(bb *testing.B) {
		ids, train := dataset.IDs(), dataset.Train()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			isrt.do(ctx, bb, c, ids[cnt%len(ids)], train[cnt%len(train)])
			cnt++
		}
		bb.StopTimer()
	})
}

func (isrt *insert) runParallel(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	var cnt int64
	b.Run("ParallelInsert", func(bb *testing.B) {
		ids, train := dataset.IDs(), dataset.Train()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				n := int(atomic.AddInt64(&cnt, 1)) - 1
				isrt.do(ctx, bb, c, ids[n%len(ids)], train[n%len(train)])
			}
		})
		bb.StopTimer()
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
