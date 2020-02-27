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
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
	"github.com/vdaas/vald/internal/log"
)

var (
	searchConfig = &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.01,
	}
	targets    []string
	addresses  []string
	wait       time.Duration
	datasetVar string
	addressVar string
	once       sync.Once
	waitVar    int64
)

func init() {
	log.Init()

	flag.StringVar(&datasetVar, "dataset", "", "available dataset(choice with comma)")
	flag.StringVar(&addressVar, "address", "", "vald gateway address")
	flag.Int64Var(&waitVar, "wait", 30, "indexing wait time(secs)")
}

func parseArgs(tb testing.TB) {
	tb.Helper()
	once.Do(func() {
		flag.Parse()
		targets = strings.Split(strings.TrimSpace(datasetVar), ",")
		addresses = strings.Split(strings.TrimSpace(addressVar), ",")
		if len(targets) != len(addresses) {
			tb.Fatal("address and dataset must have same length.")
		}
		wait = time.Duration(waitVar) * time.Second
	})
}

func BenchmarkValdGateway_Sequential(b *testing.B) {
	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			// TODO: input vald client.
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
			// TODO: input vald client.
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
