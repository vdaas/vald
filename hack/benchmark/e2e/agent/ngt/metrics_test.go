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
	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/internal/config"
	"io"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

type Metrics struct {
	Recall    []time.Duration
	Qps       []time.Duration
	BuildTime time.Duration
}

func (m Metrics) Len() int {
	return len(m.Recall)
}

func (m Metrics) Less(i, j int) bool {
	return m.Recall[i] < m.Recall[j]
}

func (m Metrics) Swap(i, j int) {
	m.Recall[i], m.Recall[j] = m.Recall[j], m.Recall[i]
	m.Qps[i], m.Qps[j] = m.Qps[j], m.Qps[i]
}

type Params struct {
	Name []string
	SearchEdgeSize []int
	CreationEdgeSize []int
	Epsilon []float32
	K []int
	Output string
}

type Param struct {
	Data assets.Dataset
	SearchEdgeSize int
	CreationEdgeSize int
	Epsilon float32
	K int
}

func (p *Param) buildIndex(tb testing.TB, client agent.AgentClient, c context.Context) time.Duration {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	client := internal.NewAgentClient(tb, ctx, address)

	ids := p.Data.IDs()
	buildStart := time.Now()
	for i, v := range p.Data.Train() {
		_, err := client.Insert(ctx, &payload.Object_Vector{
			Id:     ids[i],
			Vector: v,
		})
		if err != nil {
			tb.Error(err)
		}
	}

	_, err := client.CreateIndex(ctx, &payload.Controll_CreateIndexRequest{
		PoolSize: 10000,
	})
	if err != nil {
		if err == io.EOF {
			return 0
		}
		tb.Error(err)
	}
	buildEnd := time.Since(buildStart)
	return buildEnd
}

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

			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			for i := range searchEdgeSizes {
				searchEdgeSize := searchEdgeSizes[i]
				creationEdgeSize := creationEdgeSizes[i]

				svrs := internal.StartAgentNGTServer(t, ctx, data,
					internal.WithCreationEdgeSize(creationEdgeSize),
					internal.WithSearchEdgeSize(searchEdgeSize))
				var svr *config.Server = nil
				for _, s := range svrs {
					if s.Mode == "GRPC" {
						svr = s
						break
					}
				}
				if svr == nil {
					t.Error("grpc server should running")
				}

				buildIndex()

				datasetNeighbors := make([][]string, len(data.Neighbors()))
				for i, ns := range data.Neighbors() {
					datasetNeighbors[i] = make([]string, 0, len(ns))
					for _, n := range ns {
						datasetNeighbors[i] = append(datasetNeighbors[i], strconv.Itoa(n))
					}
				}
				querySize := len(data.Query())
				recalls := make([]float64, 0, kVar)
				qpss := make([]float64, 0, kVar)
				t.Run(fmt.Sprintf("Recall@%d", kVar), func(tt *testing.T) {
					results := make([][]string, querySize)
					for _, eps := range epsilons {
						config := *searchConfig
						config.Epsilon = eps
						config.Num = uint32(kVar)
						var elapsed time.Duration = 0.0
						for i, v := range data.Query() {
							start := time.Now()
							resp, err := client.Search(ctx, &payload.Search_Request{
								Vector: v,
								Config: searchConfig,
							})
							elapsed += time.Since(start)
							if err != nil {
								tt.Error(err)
							}
							results[i] = make([]string, len(resp.Results))
							for j, r := range resp.Results {
								results[i][j] = r.Id
							}
						}
						m, std, _ := internal.MeanStdRecalls(datasetNeighbors, results, kVar)
						recalls = append(recalls, m)
						qps := 1 / (elapsed.Seconds() / float64(querySize))
						qpss = append(qpss, qps)
						tt.Logf("mean: %f, std: %f, qps: %f", m, std, qps)
					}
				})

			})
			m := &Metrics{
				Recall: recalls,
				Qps:    qpss,
			}
			sort.Sort(m)


		}
	}
}
