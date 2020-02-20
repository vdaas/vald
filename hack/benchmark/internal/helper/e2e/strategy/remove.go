package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/helper/e2e"
	"github.com/vdaas/vald/internal/client"
)

type remove struct {
	parallel bool
}

func NewRemove(opts ...RemoveOption) e2e.Strategy {
	r := new(remove)
	for _, opt := range append(defaultRemoveOptions, opts...) {
		opt(r)
	}
	return r
}

func (r *remove) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	if r.parallel {
		r.runParallel(ctx, b, c, dataset)
		return
	}
	r.run(ctx, b, c, dataset)
}

func (r *remove) run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("Remove", func(bb *testing.B) {
		ids := dataset.IDs()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			r.do(ctx, bb, c, ids[i%len(ids)])
		}
		bb.StopTimer()
	})
}

func (r *remove) runParallel(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("ParallelRemove", func(bb *testing.B) {
		ids := dataset.IDs()

		bb.StartTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		bb.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				r.do(ctx, bb, c, ids[i%len(ids)])
				i++
			}
		})
		bb.StopTimer()
	})
}

func (r *remove) do(ctx context.Context, b *testing.B, c client.Client, id string) {
	if err := c.Remove(ctx, &client.ObjectID{
		Id: id,
	}); err != nil {
		b.Error(err)
	}
}
