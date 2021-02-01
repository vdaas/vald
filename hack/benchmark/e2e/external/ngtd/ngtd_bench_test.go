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
package ngtd

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/client/ngtd/grpc"
	"github.com/vdaas/vald/hack/benchmark/internal/client/ngtd/rest"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
	"github.com/vdaas/vald/hack/benchmark/internal/starter/external/ngtd"
	"github.com/vdaas/vald/internal/log"
)

var targets []string

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
	client, err := rest.New(ctx)
	if err != nil {
		b.Fatal(err)
	}

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithServerStarter(func(ctx context.Context, tb testing.TB, d assets.Dataset) func() {
				return ngtd.New(
					ngtd.WithDimension(d.Dimension()),
				).Run(ctx, tb)
			}),
			e2e.WithClient(client),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(client),
				),
				strategy.NewSearch(),
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
					ngtd.WithDimension(d.Dimension()),
					ngtd.WithServerType(ngtd.ServerType(ngtd.GRPC)),
				).Run(ctx, tb)
			}),
			e2e.WithClient(client),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(client),
				),
				strategy.NewSearch(),
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
					ngtd.WithDimension(d.Dimension()),
					ngtd.WithServerType(ngtd.ServerType(ngtd.GRPC)),
				).Run(ctx, tb)
			}),
			e2e.WithClient(client),
			e2e.WithStrategy(
				strategy.NewStreamInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(client),
				),
				strategy.NewStreamSearch(),
			),
		)
		bench.Run(ctx, b)
	}
}
