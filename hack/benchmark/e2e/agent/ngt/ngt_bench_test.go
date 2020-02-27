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
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
)

func BenchmarkAgentNGT_REST_Sequential(b *testing.B) {
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
					strategy.WithSearchConfig(searchConfig),
				),
				strategy.NewRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}

func BenchmarkAgentNGT_gRPC_Sequential(b *testing.B) {
	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(nil),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(nil),
					strategy.WithCreateIndexPoolSize(10000),
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

func BenchmarkAgentNGT_gRPC_Stream(b *testing.B) {
	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(nil),
			e2e.WithStrategy(
				strategy.NewStreamInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(nil),
					strategy.WithCreateIndexPoolSize(10000),
				),
				strategy.NewStreamSearch(
					strategy.WithStreamSearchConfig(searchConfig),
				),
				strategy.NewStreamRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}
