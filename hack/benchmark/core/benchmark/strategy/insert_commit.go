package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/gongt"
	"github.com/vdaas/vald/internal/core/ngt"
)

type insertCommit struct {
	poolSize uint32
}

func NewInsertCommit(opts ...InsertCommitOption) benchmark.Strategy {
	ic := new(insertCommit)
	for _, opt := range append(defaultInsertCommitOptions, opts...) {
		opt(ic)
	}
	return ic
}

func (ic *insertCommit) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, gongt gongt.NGT, dataset assets.Dataset) {
	cnt := 0
	b.Run("InsertCommit", func(bb *testing.B) {
		train := dataset.Train()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			ngt.InsertCommit(train[cnt%len(train)], ic.poolSize)
			cnt++
		}
		bb.StopTimer()
	})
}
