package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

// Benchmark is an interface for NGT benchmark.
type Benchmark interface {
	Run(context.Context, *testing.B)
}

type benchmark struct {
	name     string
	ngt      ngt.NGT
	dataset  assets.Dataset
	prestart func(context.Context, *testing.B, assets.Dataset, ngt.NGT)
}

func New(b *testing.B, opts ...Option) Benchmark {
	bm := new(benchmark)
	for _, opt := range append(defaultOptions, opts...) {
		opt(bm)
	}

	fn := assets.Data(bm.name)
	if fn == nil {
		b.Fatalf("dataset provider is nil: %v", bm.name)
	}

	bm.dataset = fn(b)
	if bm.dataset == nil {
		b.Fatalf("dataset is nil: %v", bm.name)
	}

	return bm
}

func (bm *benchmark) Run(ctx context.Context, b *testing.B) {
	func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		bm.prestart(ctx, b, bm.dataset, bm.ngt)

		b.StopTimer()
		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		b.Run(bm.name, func(bb *testing.B) {

		})
		b.StopTimer()
	}()
}
