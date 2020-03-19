package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/gongt"
	"github.com/vdaas/vald/internal/core/ngt"
)

type bulkInsertCommit struct {
	poolSize uint32
}

func NewBulkInsertCommit(opts ...BulkInsertCommitOption) benchmark.Strategy {
	bc := new(bulkInsertCommit)
	for _, opt := range append(defaultBulkInsertCommitOptions, opts...) {
		opt(bc)
	}
	return bc
}

func (bc *bulkInsertCommit) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, gongt gongt.NGT, dataset assets.Dataset) {
	b.Run("BulkInsertCommit", func(bb *testing.B) {
		train := dataset.Train()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			_, errs := ngt.BulkInsertCommit(train, bc.poolSize)
			for _, err := range errs {
				if err != nil {
					bb.Error(err)
				}
			}
		}
		bb.StopTimer()
	})
}
