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

type streamSearch struct {
	cfg *client.SearchConfig
}

func NewStreamSearch(opts ...StreamSearchOption) e2e.Strategy {
	s := new(streamSearch)
	for _, opt := range append(defaultStreamSearchOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *streamSearch) dataProvider(b *testing.B, dataset assets.Dataset) func() *client.SearchRequest {
	queries := dataset.Query()

	var cnt uint32

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()

	return func() *client.SearchRequest {
		n := int(atomic.AddUint32(&cnt, 1))
		if n > b.N {
			return nil
		}

		return &client.SearchRequest{
			Vector: queries[n%len(queries)],
			Config: s.cfg,
		}
	}
}

func (s *streamSearch) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("StreamSearch", func(bb *testing.B) {
		c.StreamSearch(ctx, s.dataProvider(b, dataset), func(_ *client.SearchResponse, err error) {
			if err != nil {
				if err != io.EOF {
					bb.Error(err)
				}
			}
		})
	})
}
