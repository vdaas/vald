package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type bulkRemove struct {
	ids       []uint
	chunkSize int
	preStart  PreStart
}

func NewBulkRemove(opts ...BulkRemoveOption) benchmark.Strategy {
	br := new(bulkRemove)
	for _, opt := range append(defaultBulkRemoveOptions, opts...) {
		opt(br)
	}
	return br
}

func (br *bulkRemove) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	cnt := 1
	b.Run("BulkRemove", func(bb *testing.B) {
		br.ids = append(br.ids, br.preStart(ctx, bb, ngt, dataset)...)

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			err := ngt.BulkRemove(br.ids[(cnt-1)*br.chunkSize : cnt*br.chunkSize]...)
			if err != nil {
				bb.Error(err)
			}
			cnt++
		}
		bb.StopTimer()
	})
}
