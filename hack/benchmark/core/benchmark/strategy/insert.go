package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/gongt"
	"github.com/vdaas/vald/internal/core/ngt"
)

type insert struct{}

// NewInsert returns `insert` strategy.
func NewInsert(opts ...InsertOption) benchmark.Strategy {
	isrt := new(insert)
	for _, opt := range append(defaultInsertOption, opts...) {
		opt(isrt)
	}
	return isrt
}

func (isrt *insert) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, gongt gongt.NGT, dataset assets.Dataset) {
	cnt := 0
	b.Run("Insert", func(bb *testing.B) {
		train := dataset.Train()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			if _, err := ngt.Insert(train[cnt%len(train)]); err != nil {
				bb.Error(err)
			}
			cnt++
		}
		bb.StopTimer()
	})
}
