package strategy

import (
	"context"
	"io"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/helper/e2e"
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

func (sr *streamRemove) dataProvider(b *testing.B, dataset assets.Dataset) func() *client.ObjectID {
	ids := dataset.IDs()

	var cnt uint32

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()

	return func() *client.ObjectID {
		n := int(atomic.AddUint32(&cnt, 1))
		if n > b.N {
			return nil
		}

		return &client.ObjectID{
			Id: ids[n%len(ids)],
		}
	}
}

func (sr *streamRemove) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("StreamRemove", func(bb *testing.B) {
		c.StreamRemove(ctx, sr.dataProvider(b, dataset), func(err error) {
			if err != nil {
				if err != io.EOF {
					b.Error(err)
				}
			}
		})
	})
}
