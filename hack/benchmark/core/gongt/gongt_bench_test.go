package gongt

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/core/benchmark/strategy"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
	"github.com/vdaas/vald/hack/benchmark/internal/core/gongt"
)

const (
	size    int     = 1
	radius  float32 = -1
	epsilon float32 = 0.01
)

var targets []string

func init() {
	testing.Init()

	var dataset string
	flag.StringVar(&dataset, "dataset", "", "available dataset(choice with comma)")
	flag.Parse()
	targets = strings.Split(strings.TrimSpace(dataset), ",")
}

func initCore(ctx context.Context, b *testing.B, dataset assets.Dataset) (core.Core64, core.Closer, error) {
	ngt, err := gongt.New(
		gongt.WithDimension(dataset.Dimension()),
		gongt.WithObjectType(dataset.ObjectType()),
	)
	if err != nil {
		return nil, nil, err
	}
	return ngt, ngt, nil
}

func BenchmarkGONGT_Insert(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsert(
					strategy.WithCore64(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkGONGT_BulkInsert(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewBulkInsert(
					strategy.WithCore64(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkGONGT_InsertCommit(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsertCommit(
					10,
					strategy.WithCore64(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkGONGT_Search(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewSearch(
					size, radius, epsilon,
					strategy.WithCore64(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkGONGT_Remove(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewRemove(
					strategy.WithCore64(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}
