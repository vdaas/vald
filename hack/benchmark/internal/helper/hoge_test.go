package helper

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/helper/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/helper/e2e/strategy"
)

func BenchmarkHoge(b *testing.B) {
	for _, target := range []string{} {
		bench := e2e.New(
			e2e.WithName(target),
			e2e.WithDataset(nil),
			e2e.WithClient(nil),
			e2e.WithStrategy(
				strategy.NewStreamInsert(),
				strategy.NewInsert(
					strategy.WithParallelInsert(),
				),
			),
		)
		bench.Run(context.Background(), b)
	}
}
