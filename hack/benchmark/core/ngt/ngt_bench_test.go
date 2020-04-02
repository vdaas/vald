//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"strings"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/core/benchmark/strategy"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
	"github.com/vdaas/vald/hack/benchmark/internal/core/ngt"
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

func initCore(ctx context.Context, b *testing.B, dataset assets.Dataset) (core.Core32, core.Closer, error) {
	ngt, err := ngt.New(
		ngt.WithDimension(dataset.Dimension()),
		ngt.WithObjectType(dataset.ObjectType()),
	)
	if err != nil {
		return nil, nil, err
	}
	return ngt, ngt, nil
}

func BenchmarkNGT_Insert(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsert(
					strategy.WithCore32(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkNGT_BulkInsert(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewBulkInsert(
					strategy.WithCore32(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkNGT_InsertCommit(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewInsertCommit(
					10,
					strategy.WithCore32(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkNGT_BulkInsertCommit(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewBulkInsertCommit(
					10,
					strategy.WithCore32(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkNGT_Search(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewSearch(
					size, radius, epsilon,
					strategy.WithCore32(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkNGT_Remove(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewRemove(
					strategy.WithCore32(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}

func BenchmarkNGT_GetVector(b *testing.B) {
	for _, target := range targets {
		benchmark.New(b,
			benchmark.WithName(target),
			benchmark.WithStrategy(
				strategy.NewGetVector(
					strategy.WithCore32(initCore),
				),
			),
		).Run(context.Background(), b)
	}
}
