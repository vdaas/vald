package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type bulkRemove struct{}

func NewBulkRemove(opts ...BulkRemoveOption) benchmark.Strategy {
	br := new(bulkRemove)
	for _, opt := range append(defaultBulkRemoveOptions, opts...) {
		opt(br)
	}
	return br
}

func (br *bulkRemove) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("BulkRemove", func(bb *testing.B) {
		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
		}
		bb.StopTimer()
	})
}
