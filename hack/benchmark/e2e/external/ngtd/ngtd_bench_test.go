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
package ngtd

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal/client/ngtd/grpc"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal/client/ngtd/rest"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal/starter/ngtd"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
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

	var dataset string
	flag.StringVar(&dataset, "dataset", "", "set available dataset list (choice with comma)")
	flag.Parse()

	targets = strings.Split(strings.TrimSpace(dataset), ",")
}

func BenchmarkNGTD_REST_Sequential(b *testing.B) {
	ctx := context.Background()

	client, err := rest.New(ctx, rest.WithAddr("http://127.0.0.1:8200"))
	if err != nil {
		b.Fatal(err)
	}

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithServerStarter(func(ctx context.Context, tb testing.TB, d assets.Dataset) func() {
				return ngtd.New(
					ngtd.WithDimentaion(d.Dimension()),
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

func BenchmarkNGTD_gRPC_Sequential(b *testing.B) {
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
				return ngtd.New(
					ngtd.WithDimentaion(d.Dimension()),
					ngtd.WithServerType(ngtd.ServerType(ngtd.GRPC)),
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

func BenchmarkNGTD_gRPC_Stream(b *testing.B) {
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
				return ngtd.New(
					ngtd.WithDimentaion(d.Dimension()),
					ngtd.WithServerType(ngtd.ServerType(ngtd.GRPC)),
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
