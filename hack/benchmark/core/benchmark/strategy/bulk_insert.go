package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type bulkInsert struct{}

func NewBulkInsert(opts ...BulkInsertOption) benchmark.Strategy {
	bi := new(bulkInsert)
	for _, opt := range append(defaultBulkInsertOptions, opts...) {
		opt(bi)
	}
	return bi
}

func (bi *bulkInsert) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("BulkInsert", func(bb *testing.B) {
		train := dataset.Train()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			_, errs := ngt.BulkInsert(train)
			for _, err := range errs {
				if err != nil {
					bb.Error(err)
				}
			}
		}
		bb.StopTimer()
	})
}
