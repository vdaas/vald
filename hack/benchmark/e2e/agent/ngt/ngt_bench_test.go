//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
package ngt

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
	"github.com/vdaas/vald/hack/benchmark/internal/starter/agent"
	"github.com/vdaas/vald/internal/client/agent/grpc"
	"github.com/vdaas/vald/internal/client/agent/rest"
	"github.com/vdaas/vald/internal/log"
)

var (
	targets []string
)

var searchConfig = &payload.Search_Config{
	Num:     1,
	Radius:  -1,
	Epsilon: 0.01,
}

func init() {
	testing.Init()
	log.Init()

	var (
		dataset string
		num     uint
		radius  float64
		epsilon float64
	)

	flag.StringVar(&dataset, "dataset", "", "set available dataset list (choice with comma)")
	flag.UintVar(&num, "num", uint(searchConfig.Num), "search response size")
	flag.Float64Var(&radius, "radius", float64(searchConfig.Radius), "search radius size")
	flag.Float64Var(&epsilon, "epsilon", float64(searchConfig.Epsilon), "search epsilon size")
	flag.Parse()

	searchConfig.Num = uint32(num)
	searchConfig.Radius = float32(radius)
	searchConfig.Epsilon = float32(epsilon)
	targets = strings.Split(strings.TrimSpace(dataset), ",")
}

func BenchmarkAgentNGT_REST_Sequential(b *testing.B) {
	ctx := context.Background()

	client := rest.New(ctx)

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithServerStarter(func(ctx context.Context, tb testing.TB, d assets.Dataset) func() {
				return agent.New(
					agent.WithDimentaion(d.Dimension()),
					agent.WithDistanceType(d.DistanceType()),
					agent.WithObjectType(d.ObjectType()),
				).Run(ctx, tb)
			}),
			e2e.WithClient(client),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(client),
				),
				strategy.NewSearch(
					strategy.WithSearchConfig(searchConfig),
				),
			),
		)
		bench.Run(ctx, b)
	}
}

func BenchmarkAgentNGT_gRPC_Sequential(b *testing.B) {
	ctx := context.Background()

	client, err := grpc.New(ctx)
	if err != nil {
		b.Fatal(err)
	}

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithServerStarter(func(ctx context.Context, tb testing.TB, d assets.Dataset) func() {
				return agent.New(
					agent.WithDimentaion(d.Dimension()),
					agent.WithDistanceType(d.DistanceType()),
					agent.WithObjectType(d.ObjectType()),
				).Run(ctx, tb)
			}),
			e2e.WithClient(client),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(client),
				),
				strategy.NewSearch(
					strategy.WithSearchConfig(searchConfig),
				),
			),
		)
		bench.Run(ctx, b)
	}
}

func BenchmarkAgentNGT_gRPC_Stream(b *testing.B) {
	ctx := context.Background()

	client, err := grpc.New(ctx)
	if err != nil {
		b.Fatal(err)
	}

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithServerStarter(func(ctx context.Context, tb testing.TB, d assets.Dataset) func() {
				return agent.New(
					agent.WithDimentaion(d.Dimension()),
					agent.WithDistanceType(d.DistanceType()),
					agent.WithObjectType(d.ObjectType()),
				).Run(ctx, tb)
			}),
			e2e.WithClient(client),
			e2e.WithStrategy(
				strategy.NewStreamInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(client),
				),
				strategy.NewStreamSearch(
					strategy.WithStreamSearchConfig(searchConfig),
				),
			),
		)
		bench.Run(ctx, b)
	}
}
