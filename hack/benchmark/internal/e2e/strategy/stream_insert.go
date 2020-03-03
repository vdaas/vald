package strategy

import (
	"context"
	"io"
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

func (sisrt *streamInsert) dataProvider(b *testing.B, dataset assets.Dataset) func() *client.ObjectVector {
	ids, trains := dataset.IDs(), dataset.Train()

	var cnt uint32

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	defer b.StopTimer()

	return func() *client.ObjectVector {
		n := int(atomic.AddUint32(&cnt, 1))
		if n > b.N {
			return nil
		}

		return &client.ObjectVector{
			Id:     ids[n%len(ids)],
			Vector: trains[n%len(trains)],
		}
	}
}

func (sisrt *streamInsert) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("StreamInsert", func(bb *testing.B) {
		c.StreamInsert(ctx, sisrt.dataProvider(b, dataset), func(err error) {
			if err != nil {
				if err != io.EOF {
					bb.Error(err)
				}
			}
		})
	})
}
