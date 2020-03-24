package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

// Benchmark is an interface for NGT benchmark.
type Benchmark interface {
	Run(context.Context, *testing.B)
}

type Type uint32

const (
	Float32 Type = iota
	Float64
)

type benchmark struct {
	name       string
	core       interface{}
	dataset    assets.Dataset
	typ        Type
	strategies []Strategy
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

		b.StopTimer()
		b.ReportAllocs()
		b.ResetTimer()
		b.StartTimer()
		b.Run(bm.name, func(bb *testing.B) {
			for _, strategy := range bm.strategies {
				strategy.Run(ctx, bb, bm.core, bm.typ, bm.dataset)
			}
		})
		b.StopTimer()
	}()
}
