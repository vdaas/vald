//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

func TestMetrics(rt *testing.T) {
	parseArgs(rt)
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()

	for N, name := range targets {
		address := addresses[N]
		if address == "" {
			address = "localhost:8082"
		}

		if name == "" {
			continue
		}
		rt.Run(name, func(t *testing.T) {
			data := assets.Data(name)(t)
			if data == nil {
				t.Logf("assets %s is nil", name)
				return
			}
			ids := data.IDs()

			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			if strings.Contains(address, "localhost") {
				internal.StartAgentNGTServer(t, ctx, data)
			}

			client := internal.NewAgentClient(t, ctx, address)

			for i, v := range data.Train() {
				_, err := client.Insert(ctx, &payload.Object_Vector{
					Id:     ids[i],
					Vector: v,
				})
				if err != nil {
					t.Error(err)
				}
			}

			_, err := client.CreateIndex(ctx, &payload.Controll_CreateIndexRequest{
				PoolSize: 10000,
			})
			if err != nil {
				if err == io.EOF {
					return
				}
				t.Error(err)
			}

			datasetNeighbors := make([][]string, len(data.Neighbors()))
			for i, ns := range data.Neighbors() {
				datasetNeighbors[i] = make([]string, 0, len(ns))
				for _, n := range ns {
					datasetNeighbors[i] = append(datasetNeighbors[i], strconv.Itoa(n))
				}
			}
			querySize := len(data.Query())
			for k := 1; k <= int(searchConfig.Num); k++ {
				t.Run(fmt.Sprintf("Recall@%d", k), func(tt *testing.T) {
					results := make([][]string, querySize)
					var qps time.Duration = 0.0
					for i, v := range data.Query() {
						config := *searchConfig
						config.Num = uint32(k)
						start := time.Now()
						resp, err := client.Search(ctx, &payload.Search_Request{
							Vector: v,
							Config: &config,
						})
						qps += time.Since(start)
						if err != nil {
							tt.Error(err)
						}
						results[i] = make([]string, len(resp.Results))
						for j, r := range resp.Results {
							results[i][j] = r.Id
						}
					}
					m, std, _ := internal.MeanStdRecalls(datasetNeighbors, results, k)
					tt.Logf("mean: %f, std: %f, qps: %f", m, std, qps.Seconds()/float64(querySize))
				})
			}
		})
	}
}
