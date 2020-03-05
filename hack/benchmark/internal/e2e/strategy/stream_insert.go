package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client"
)

type streamInsert struct{}

func NewStreamInsert(opts ...StreamInsertOption) e2e.Strategy {
	s := new(streamInsert)
	for _, opt := range append(defaultStreamInsertOptions, opts...) {
		opt(s)
	}
	return s
}

func (sisrt *streamInsert) dataProvider(total *uint32, b *testing.B, dataset assets.Dataset) func() *client.ObjectVector {
	ids, trains := dataset.IDs(), dataset.Train()

	var cnt uint32

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()

	return func() *client.ObjectVector {
		n := int(atomic.AddUint32(&cnt, 1)) - 1
		if n >= b.N {
			return nil
		}

		total := int(atomic.AddUint32(total, 1)) - 1
		return &client.ObjectVector{
			Id:     ids[total%len(ids)],
			Vector: trains[total%len(trains)],
		}
	}
}

func (sisrt *streamInsert) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	var total uint32
	b.Run("StreamInsert", func(bb *testing.B) {
		c.StreamInsert(ctx, sisrt.dataProvider(&total, bb, dataset), func(err error) {
			if err != nil {
				bb.Error(err)
			}
		})
	})
}
