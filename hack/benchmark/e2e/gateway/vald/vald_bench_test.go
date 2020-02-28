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

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
	"github.com/vdaas/vald/internal/log"
)

var (
	targets   []string
	addresses []string
	wait      time.Duration
)

var searchConfig = &payload.Search_Config{
	Num:     10,
	Radius:  -1,
	Epsilon: 0.01,
}

func init() {
	testing.Init()

	log.Init()

	var (
		dataset     string
		address     string
		waitSeconds uint
	)

	flag.StringVar(&dataset, "dataset", "", "set dataset (choice with comma)")
	flag.StringVar(&address, "address", "0.0.0.0:5001", "set vald gateway address")
	flag.UintVar(&waitSeconds, "wait", 30, "indexing wait time(secs)")
	flag.Parse()

	targets = strings.Split(strings.TrimSpace(dataset), ",")
	addresses = strings.Split(strings.TrimSpace(address), ",")
	wait = time.Duration(time.Duration(waitSeconds) * time.Second)
}

func BenchmarkValdGateway_Sequential(b *testing.B) {
	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(nil),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewSearch(
					strategy.WithSearchConfig(searchConfig),
				),
				strategy.NewRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}

func BenchmarkValdGateway_Stream(b *testing.B) {
	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(nil),
			e2e.WithStrategy(
				strategy.NewStreamInsert(),
				strategy.NewStreamSearch(
					strategy.WithStreamSearchConfig(searchConfig),
				),
				strategy.NewStreamRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}
