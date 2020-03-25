package ngt

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/core/benchmark/strategy"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/ngt"
	"github.com/vdaas/vald/internal/log"
)

var targets []string

func init() {
	testing.Init()
	log.Init()

	var dataset string
	flag.StringVar(&dataset, "dataset", "", "available dataset(choice with comma)")
	flag.Parse()
	targets = strings.Split(dataset, ",")
}

func BenchmarkNGT_Insert(b *testing.B) {
	for _, target := range targets {
		benchmark.New(
			b,
			benchmark.WithName(target),
			benchmark.WithFloat32(
				func(ctx context.Context, b *testing.B, dataset assets.Dataset) (interface{}, func(), error) {
					n, err := ngt.New(
						ngt.WithDimension(dataset.Dimension()),
						ngt.WithIndexPath(dataset.ObjectType()),
					)
					if err != nil {
						return nil, nil, err
					}
					return n, n.Close, nil
				},
				strategy.NewInsert(),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkNGT_BulkInsert(b *testing.B) {
	for _, target := range targets {
		benchmark.New(
			b,
			benchmark.WithName(target),
			benchmark.WithFloat32(
				func(ctx context.Context, b *testing.B, dataset assets.Dataset) (interface{}, func(), error) {
					n, err := ngt.New(
						ngt.WithDimension(dataset.Dimension()),
						ngt.WithIndexPath(dataset.ObjectType()),
					)
					if err != nil {
						return nil, nil, err
					}
					return n, n.Close, nil
				},
				strategy.NewBulkInsert(),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkNGT_Remove(b *testing.B) {
	for _, target := range targets {
		benchmark.New(
			b,
			benchmark.WithName(target),
			benchmark.WithFloat32(
				func(ctx context.Context, b *testing.B, dataset assets.Dataset) (interface{}, func(), error) {
					n, err := ngt.New(
						ngt.WithDimension(dataset.Dimension()),
						ngt.WithIndexPath(dataset.ObjectType()),
					)
					if err != nil {
						return nil, nil, err
					}
					return n, n.Close, nil
				},
				strategy.NewRemove(),
			),
		).Run(context.Background(), b)
	}
}
