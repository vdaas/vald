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
package vald

import (
	"context"
	"flag"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
)

var (
	targets  []string
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
	flag.StringVar(&grpcAddr, "grpc_address", "127.0.0.1:8081", "set vald gateway address for gRPC")
	flag.UintVar(&waitSeconds, "wait", 30, "indexing wait time (secs)")
	flag.Parse()

	targets = strings.Split(strings.TrimSpace(dataset), ",")
	wait = time.Duration(time.Duration(waitSeconds) * time.Second)
}

func BenchmarkValdGateway_gRPC_Sequential(b *testing.B) {
	client, err := vald.New(
		vald.WithClient(grpc.New(
			grpc.WithAddrs(grpcAddr),
			grpc.WithInsecure(true),
		)),
		vald.WithAddrs(
			grpcAddr,
		),
	)
	if err != nil {
		b.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.TODO())
	_, err = client.Start(ctx)
	if err != nil {
		b.Fatal(err)
	}
	defer client.Stop(ctx)
	defer cancel()
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

func BenchmarkValdGateway_gRPC_Stream(b *testing.B) {
	client, err := vald.New(
		vald.WithClient(grpc.New(
			grpc.WithAddrs(grpcAddr),
			grpc.WithInsecure(true),
		)),
		vald.WithAddrs(
			grpcAddr,
		),
	)
	if err != nil {
		b.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.TODO())
	_, err = client.Start(ctx)
	if err != nil {
		b.Fatal(err)
	}
	defer client.Stop(ctx)
	defer cancel()
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
