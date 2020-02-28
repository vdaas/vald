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
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal/client/ngtd/rest"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
	"github.com/vdaas/vald/internal/log"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd"
	"github.com/yahoojapan/ngtd/kvs"
)

const (
	baseDir = "/tmp/ngtd/"
	port    = 8200
)

var (
	targets []string
	dataset string
)

var searchConfig = &payload.Search_Config{
	Num:     10,
	Radius:  -1,
	Epsilon: 0.01,
}

func init() {
	testing.Init()

	log.Init()

	if err := os.RemoveAll(baseDir); err != nil {
		log.Error(err)
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Error(err)
	}

	var dataset string
	flag.StringVar(&dataset, "dataset", "", "set available dataset list (choice with comma)")
	flag.Parse()

	targets = strings.Split(strings.TrimSpace(dataset), ",")
}

func StartNGTD(tb testing.TB, t ngtd.ServerType, dim int) func() {
	tb.Helper()

	gongt.SetDimension(dim)

	db, err := kvs.NewGoLevel(baseDir + "meta")
	if err != nil {
		tb.Error(err)
	}

	n, err := ngtd.NewNGTD(baseDir+"ngt", db, port)
	if err != nil {
		tb.Error(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		wg.Done()

		if err := n.ListenAndServe(t); err != nil {
			tb.Errorf("ngtd returned error: %s", err.Error())
		}
	}()

	wg.Wait()

	return func() {
		n.Stop()

		if err := os.RemoveAll(baseDir + "meta"); err != nil {
			tb.Error(err)
		}

		if err := os.RemoveAll(baseDir + "ngt"); err != nil {
			tb.Error(err)
		}
	}
}

func BenchmarkNGTD_REST_Sequential(b *testing.B) {
	client, err := rest.New(context.Background(), rest.WithAddr("127.0.0.1:"+strconv.Itoa(port)))
	if err != nil {
		b.Error(err)
	}

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithServerStarter(func(tb testing.TB, d assets.Dataset) func() {
				return StartNGTD(tb, ngtd.HTTP, d.Dimension())
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
				strategy.NewRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}

func BenchmarkNGTD_gRPC_Sequential(b *testing.B) {
	client, err := rest.New(context.Background(), rest.WithAddr("127.0.0.1:"+strconv.Itoa(port)))
	if err != nil {
		b.Error(err)
	}

	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithServerStarter(func(tb testing.TB, d assets.Dataset) func() {
				return StartNGTD(tb, ngtd.GRPC, d.Dimension())
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
				strategy.NewRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}
