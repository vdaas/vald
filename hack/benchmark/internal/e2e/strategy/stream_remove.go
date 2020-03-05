package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client"
)

type streamRemove struct{}

func NewStreamRemove(opts ...StreamRemoveOption) e2e.Strategy {
	sr := new(streamRemove)
	for _, opt := range append(defaultStreamRemoveOptions, opts...) {
		opt(sr)
	}
	return sr
}

func (sr *streamRemove) dataProvider(total *uint32, b *testing.B, dataset assets.Dataset) func() *client.ObjectID {
	ids := dataset.IDs()

	var cnt uint32

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	defer b.StopTimer()

	return func() *client.ObjectID {
		n := int(atomic.AddUint32(&cnt, 1)) - 1
		if n >= b.N {
			return nil
		}

		total := int(atomic.AddUint32(total, 1))
		return &client.ObjectID{
			Id: ids[total%len(ids)],
		}
	}
}

func (sr *streamRemove) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	var total uint32
	b.Run("StreamRemove", func(bb *testing.B) {
		c.StreamRemove(ctx, sr.dataProvider(&total, bb, dataset), func(err error) {
			if err != nil {
				b.Error(err)
			}
		})
	})
}
