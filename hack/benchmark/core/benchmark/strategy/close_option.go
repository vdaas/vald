package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type close struct{}

func NewClose() benchmark.Strategy {
	return nil
}

func (c *close) Run(ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset) {
	b.Run("Close", func(bb *testing.B) {
		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
		}
		bb.StopTimer()
	})
}
