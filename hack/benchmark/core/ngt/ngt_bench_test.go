//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package ngt provides benchmark program
package ngt

import (
	"context"
	"flag"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/core/benchmark/strategy"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm"
	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/strings"
)

const (
	size    int     = 1
	radius  float32 = -1
	epsilon float32 = 0.1
)

var targets []string

func init() {
	testing.Init()
	log.Init(log.WithLoggerType(logger.NOP.String()))

	var dataset string
	flag.StringVar(&dataset, "dataset", "", "available dataset(choice with comma)")
	flag.Parse()
	targets = strings.Split(strings.TrimSpace(dataset), ",")
}

func initCore(ctx context.Context, b *testing.B, dataset assets.Dataset) (algorithm.Bit32, algorithm.Closer, error) {
	ngt, err := ngt.New(
		ngt.WithDimension(dataset.Dimension()),
		ngt.WithObjectType(dataset.ObjectType()),
	)
	if err != nil {
		return nil, nil, err
	}
	return ngt, ngt, nil
}

func BenchmarkNGTSequential_Insert(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsert(
					strategy.WithBit32(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTParallel_Insert(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsert(
					strategy.WithBit32(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTSequential_BulkInsert(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewBulkInsert(
					strategy.WithBit32(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTParallel_BulkInsert(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewBulkInsert(
					strategy.WithBit32(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTSequential_InsertCommit(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsertCommit(
					10,
					strategy.WithBit32(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTParallel_InsertCommit(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsertCommit(
					10,
					strategy.WithBit32(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTSequential_BulkInsertCommit(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewBulkInsertCommit(
					10,
					strategy.WithBit32(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTParallel_BulkInsertCommit(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewBulkInsertCommit(
					10,
					strategy.WithBit32(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTSequential_Search(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewSearch(
					size, radius, epsilon,
					strategy.WithBit32(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTParallel_Search(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewSearch(
					size, radius, epsilon,
					strategy.WithBit32(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTSequential_Remove(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewRemove(
					strategy.WithBit32(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTParallel_Remove(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewRemove(
					strategy.WithBit32(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTSequential_GetVector(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewGetVector(
					strategy.WithBit32(initCore),
				),
			),
		).Run(ctx, b)
	}
}

func BenchmarkNGTParallel_GetVector(b *testing.B) {
	ctx := context.Background()
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewGetVector(
					strategy.WithBit32(initCore),
					strategy.WithParallel(),
				),
			),
		).Run(ctx, b)
	}
}
