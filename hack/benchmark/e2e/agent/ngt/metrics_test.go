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
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/wcharczuk/go-chart"
)

var (
	outputPathVar string
)

type Metrics struct {
	Recall []float64
	Qps    []float64
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

func init() {
	flag.StringVar(&outputPathVar, "output", "metrics.svg", "metrics output path")
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
			recalls := make([]float64, 0, searchConfig.Num)
			qpss := make([]float64, 0, searchConfig.Num)
			for k := 1; k <= int(searchConfig.Num); k++ {
				t.Run(fmt.Sprintf("Recall@%d", k), func(tt *testing.T) {
					results := make([][]string, querySize)
					var elapsed time.Duration = 0.0
					for i, v := range data.Query() {
						config := *searchConfig
						config.Num = uint32(k)
						start := time.Now()
						resp, err := client.Search(ctx, &payload.Search_Request{
							Vector: v,
							Config: &config,
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
					m, std, _ := internal.MeanStdRecalls(datasetNeighbors, results, k)
					recalls = append(recalls, m)
					qps := 1 / (elapsed.Seconds() / float64(querySize))
					qpss = append(qpss, math.Log(qps))
					tt.Logf("mean: %f, std: %f, qps: %f", m, std, qps)
				})
			}
			m := &Metrics{
				Recall: recalls,
				Qps: qpss,
			}
			sort.Sort(m)

			graph := chart.Chart{
				Title:          "QPS per Recall",
				TitleStyle:     chart.StyleShow(),
				Width: 1280,
				Height: 960,
				Background: chart.Style{
					Padding: chart.Box{ Top: 50 },
				},
				XAxis: chart.XAxis{
					Name: "Recall",
					Style: chart.StyleShow(),
					Range: &chart.ContinuousRange{
						Min: 0,
						Max: 1,
					},
				},
				YAxis: chart.YAxis{
					Name: "Query per second (1 / s) log scale",
					Style: chart.StyleShow(),
				},
				Series: []chart.Series{
					chart.ContinuousSeries{
						XValues: m.Recall,
						YValues: m.Qps,
					},
				},
			}
			f, err := os.Create(outputPathVar)
			if err != nil {
				t.Error(err)
			}
			defer f.Close()
			if err := graph.Render(chart.SVG, f); err != nil {
				t.Error(err)
			}

		})
	}
}
