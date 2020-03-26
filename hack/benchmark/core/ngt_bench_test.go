package ngt

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/core/benchmark/strategy"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
	"github.com/vdaas/vald/hack/benchmark/internal/core/ngt"
)

var targets []string

func init() {
	testing.Init()

	var dataset string
	flag.StringVar(&dataset, "dataset", "", "available dataset(choice with comma)")
	flag.Parse()
	targets = strings.Split(strings.TrimSpace(dataset), ",")
}

func BenchmarkNGT_Insert(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsert(
					strategy.WithCore32(
						func(ctx context.Context, b *testing.B, dataset assets.Dataset) (core.Core32, core.Closer, error) {
							ngt, err := ngt.New(
								ngt.WithDimension(dataset.Dimension()),
								ngt.WithObjectType(dataset.ObjectType()),
							)
							if err != nil {
								return nil, nil, err
							}
							return ngt, ngt, nil
						}),
				),
			),
		).Run(context.Background(), b)
	}
}
