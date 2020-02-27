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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kpango/glg"
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
	targets    []string
	datasetVar string
	once       sync.Once
)

func init() {
	log.Init()
	glg.Get().SetMode(glg.NONE)
	if err := os.RemoveAll(baseDir); err != nil {
		log.Error(err)
	}
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Error(err)
	}

	flag.StringVar(&datasetVar, "assets", "", "list available assets(choice with comma)")
}

func parseArgs() {
	once.Do(func() {
		flag.Parse()
		targets = strings.Split(strings.TrimSpace(datasetVar), ",")
	})
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

	go func() {
		err := n.ListenAndServe(t)
		if err != nil {
			tb.Errorf("ngtd returned error: %s", err.Error())
		}
	}()

	time.Sleep(5 * time.Second)

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
	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(nil),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(nil),
				),
				strategy.NewSearch(
					strategy.WithSearchConfig(nil),
				),
				strategy.NewRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}

func BenchmarkNGTD_gRPC_Sequential(b *testing.B) {
	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(nil),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(nil),
				),
				strategy.NewSearch(
					strategy.WithSearchConfig(nil),
				),
				strategy.NewRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}
