//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package gongt provides benchmark program
package gongt

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/core/benchmark/strategy"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm"
	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm/gongt"
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

func initCore(ctx context.Context, b *testing.B, dataset assets.Dataset) (algorithm.Bit64, algorithm.Closer, error) {
	ngt, err := gongt.New(
		gongt.WithDimension(dataset.Dimension()),
		gongt.WithObjectType(dataset.ObjectType()),
	)
	if err != nil {
		return nil, nil, err
	}
	return ngt, ngt, nil
}

func BenchmarkGoNGTSequential_Insert(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsert(
					strategy.WithBit64(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTParallel_Insert(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsert(
					strategy.WithBit64(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTSequential_BulkInsert(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewBulkInsert(
					strategy.WithBit64(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTParallel_BulkInsert(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewBulkInsert(
					strategy.WithBit64(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTSequential_InsertCommit(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsertCommit(
					10,
					strategy.WithBit64(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTParallel_InsertCommit(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsertCommit(
					10,
					strategy.WithBit64(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTSequential_Search(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewSearch(
					size, radius, epsilon,
					strategy.WithBit64(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTParallel_Search(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewSearch(
					size, radius, epsilon,
					strategy.WithBit64(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTSequential_Remove(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewRemove(
					strategy.WithBit64(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTParallel_Remove(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewRemove(
					strategy.WithBit64(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTSequential_GetVector(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewGetVector(
					strategy.WithBit64(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkGoNGTParallel_GetVector(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewGetVector(
					strategy.WithBit64(initCore),
				),
			),
		).Run(ctx, b)
	}
}
