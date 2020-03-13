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
package vald

import (
	"context"
	"flag"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
	"github.com/vdaas/vald/internal/client/gateway/vald/grpc"
	"github.com/vdaas/vald/internal/client/gateway/vald/rest"
	"github.com/vdaas/vald/internal/log"
)

var (
	targets  []string
	restAddr string
	grpcAddr string
	wait     time.Duration
)

func init() {
	testing.Init()
	log.Init()

	var (
		dataset     string
		waitSeconds uint
	)

	flag.StringVar(&dataset, "dataset", "", "set dataset (choice with comma)")
	flag.StringVar(&restAddr, "rest_address", "http://127.0.0.1:8080", "set vald gateway address for REST")
	flag.StringVar(&grpcAddr, "grpc_address", "127.0.0.1:8081", "set vald gateway address for gRPC")
	flag.UintVar(&waitSeconds, "wait", 30, "indexing wait time (secs)")
	flag.Parse()

	targets = strings.Split(strings.TrimSpace(dataset), ",")
	wait = time.Duration(time.Duration(waitSeconds) * time.Second)
}

func BenchmarkGateway_REST_Sequential(b *testing.B) {
	ctx := context.Background()

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(
				rest.New(
					rest.WithAddr(
						restAddr,
					),
				),
			),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewSearch(),
			),
		)
		bench.Run(ctx, b)
	}
}

func BenchmarkGateway_REST_Stream(b *testing.B) {
	ctx := context.Background()

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(
				rest.New(
					rest.WithAddr(
						restAddr,
					),
				),
			),
			e2e.WithStrategy(
				strategy.NewStreamInsert(),
				strategy.NewStreamSearch(),
			),
		)
		bench.Run(ctx, b)
	}
}

func BenchmarkGateway_gRPC_Sequential(b *testing.B) {
	ctx := context.Background()
	client, err := grpc.New(ctx,
		grpc.WithAddr(
			grpcAddr,
		),
	)
	if err != nil {
		b.Fatal(err)
	}

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(client),
			e2e.WithStrategy(
				strategy.NewStreamInsert(),
				strategy.NewStreamSearch(),
			),
		)
		bench.Run(ctx, b)
	}
}

func BenchmarkGateway_gRPC_Stream(b *testing.B) {
	ctx := context.Background()
	client, err := grpc.New(ctx,
		grpc.WithAddr(
			grpcAddr,
		),
	)
	if err != nil {
		b.Fatal(err)
	}

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(client),
			e2e.WithStrategy(
				strategy.NewStreamInsert(),
				strategy.NewStreamSearch(),
			),
		)
		bench.Run(ctx, b)
	}
}
