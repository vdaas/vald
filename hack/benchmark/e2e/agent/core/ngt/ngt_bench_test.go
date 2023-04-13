// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ngt

import (
	"context"
	"flag"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/operation"
	"github.com/vdaas/vald/hack/benchmark/internal/starter/agent/core/ngt"
	"github.com/vdaas/vald/internal/client/v1/client/agent/core"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/strings"
)

var (
	datasets []string
	grpcAddr string
)

func init() {
	testing.Init()
	log.Init(log.WithLoggerType(logger.NOP.String()))

	var dataset string

	flag.StringVar(&dataset, "dataset", "identity-128", "set available dataset list (choice with comma)")
	flag.StringVar(&grpcAddr, "grpc_address", "127.0.0.1:8081", "set vald agent address for gRPC")
	flag.Parse()

	datasets = strings.Split(strings.TrimSpace(dataset), ",")
}

func BenchmarkAgentNGT_gRPC_Sequential(b *testing.B) {
	for _, dname := range datasets {
		b.Run(dname, func(b *testing.B) {
			ctx := context.Background()

			dataset := assets.Data(dname)(b)

			c, err := core.New(
				core.WithAddrs(grpcAddr),
				core.WithGRPCClient(
					grpc.New(
						grpc.WithAddrs(grpcAddr),
						grpc.WithInsecure(true),
					),
				),
			)
			if err != nil {
				b.Fatal(err)
			}

			defer ngt.New(
				ngt.WithDimension(dataset.Dimension()),
				ngt.WithDistanceType(dataset.DistanceType()),
				ngt.WithObjectType(dataset.ObjectType()),
				ngt.WithClient(c),
			).Run(ctx, b)()

			op := operation.New(
				operation.WithClient(c),
				operation.WithIndexer(c),
			)

			insertedNum := op.Insert(ctx, b, dataset)
			op.CreateIndex(ctx, b)
			op.Search(ctx, b, dataset)
			op.SearchByID(ctx, b, insertedNum)
			op.Remove(ctx, b, insertedNum)
		})
	}
}

func BenchmarkAgentNGT_gRPC_Stream(b *testing.B) {
	for _, dname := range datasets {
		b.Run(dname, func(b *testing.B) {
			ctx := context.Background()

			dataset := assets.Data(dname)(b)

			c, err := core.New(
				core.WithAddrs(grpcAddr),
				core.WithGRPCClient(
					grpc.New(
						grpc.WithAddrs(grpcAddr),
						grpc.WithInsecure(true),
					),
				),
			)
			if err != nil {
				b.Fatal(err)
			}

			defer ngt.New(
				ngt.WithDimension(dataset.Dimension()),
				ngt.WithDistanceType(dataset.DistanceType()),
				ngt.WithObjectType(dataset.ObjectType()),
				ngt.WithClient(c),
			).Run(ctx, b)()

			op := operation.New(
				operation.WithClient(c),
				operation.WithIndexer(c),
			)

			insertedNum := op.StreamInsert(ctx, b, dataset)
			op.CreateIndex(ctx, b)
			op.StreamSearch(ctx, b, dataset)
			op.StreamSearchByID(ctx, b, insertedNum)
			op.StreamRemove(ctx, b, insertedNum)
		})
	}
}
