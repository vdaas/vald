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
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/metrics"
	"gopkg.in/yaml.v2"
)

var (
	configuration string
	outputPath    string
)

type Params struct {
	Host             string    `json:"host" yaml:"host"`
	Port             uint      `json:"port" yaml:"port"`
	DatasetName      []string  `json:"dataset_name" yaml:"dataset_name"`
	SearchEdgeSize   []int     `json:"search_edge_size" yaml:"search_edge_size"`
	CreationEdgeSize []int     `json:"creation_edge_size" yaml:"creation_edge_size"`
	Epsilon          []float32 `json:"epsilon" yaml:"epsilon"`
	K                []int     `json:"k" yaml:"k"`
}

func (p *Params) Address() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

func init() {
	flag.StringVar(&configuration, "conf", "metrics.yaml", "set metrics configuration file path")
	flag.StringVar(&outputPath, "output", "metrics.gob", "set result output path")
}

func TestMetrics(rt *testing.T) {
	flag.Parse()

	input, err := os.Open(configuration)
	if err != nil {
		rt.Fatal(err)
	}
	defer func() {
		if err := input.Close(); err != nil {
			rt.Error(err)
		}
	}()
	var p Params
	if err := yaml.NewDecoder(input).Decode(&p); err != nil {
		rt.Fatal(err)
	}

	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()

	m := make([]*metrics.Metrics, 0, len(p.SearchEdgeSize)*len(p.CreationEdgeSize)*len(p.DatasetName))
	for _, name := range p.DatasetName {
		rt.Run(name, func(t *testing.T) {
			data := assets.Data(name)(t)
			datasetNeighbors := make([][]string, len(data.Neighbors()))
			for i, ns := range data.Neighbors() {
				datasetNeighbors[i] = make([]string, 0, len(ns))
				for _, n := range ns {
					datasetNeighbors[i] = append(datasetNeighbors[i], strconv.Itoa(n))
				}
			}
			if data == nil {
				t.Logf("assets %s is nil", name)
				return
			}

			for _, searchEdgeSize := range p.SearchEdgeSize {
				for _, creationEdgeSize := range p.CreationEdgeSize {
					t.Run(fmt.Sprintf("%d-%d", searchEdgeSize, creationEdgeSize), func(tt *testing.T) {
						ctx, cancel := context.WithCancel(rctx)
						defer cancel()

						internal.StartAgentNGTServer(t, ctx, data,
							internal.WithCreationEdgeSize(creationEdgeSize),
							internal.WithSearchEdgeSize(searchEdgeSize),
							internal.WithGRPCHost(p.Host),
							internal.WithGRPCPort(p.Port))

						client := internal.NewAgentClient(t, ctx, p.Address())

						buildStart := time.Now()
						for i, v := range data.Train() {
							_, err := client.Insert(ctx, &payload.Object_Vector{
								Id:     data.IDs()[i],
								Vector: v,
							})
							if err != nil {
								tt.Error(err)
							}
						}

						_, err := client.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
							PoolSize: 10000,
						})
						if err != nil {
							t.Error(err)
						}
						buildEnd := time.Since(buildStart)

						for _, k := range p.K {
							searchMetrics := make([]*metrics.SearchMetrics, 0, len(p.Epsilon))
							for _, eps := range p.Epsilon {
								querySize := len(data.Query())
								tt.Run(fmt.Sprintf("Recall@%d with %f", k, eps), func(ttt *testing.T) {
									results := make([][]string, querySize)
									config := *searchConfig
									config.Epsilon = eps
									config.Num = uint32(k)
									var elapsed time.Duration = 0.0
									for i, v := range data.Query() {
										start := time.Now()
										resp, err := client.Search(ctx, &payload.Search_Request{
											Vector: v,
											Config: &config,
										})
										elapsed += time.Since(start)
										if err != nil {
											ttt.Error(err)
										}
										results[i] = make([]string, len(resp.Results))
										for j, r := range resp.Results {
											results[i][j] = r.Id
										}
									}
									recall, _, _ := internal.MeanStdRecalls(datasetNeighbors, results, k)
									qps := 1 / (elapsed.Seconds() / float64(querySize))
									searchMetrics = append(searchMetrics, &metrics.SearchMetrics{
										Recall:  recall,
										Qps:     qps,
										Epsilon: eps,
									})
									ttt.Logf("recall: %f, qps: %f", recall, qps)
								})
							}
							m = append(m, &metrics.Metrics{
								BuildTime:        int64(buildEnd),
								DatasetName:      name,
								SearchEdgeSize:   searchEdgeSize,
								CreationEdgeSize: creationEdgeSize,
								K:                k,
								Search:           searchMetrics,
							})
						}
					})
				}
			}
		})
	}
	output, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		rt.Error(err)
	}
	defer func() {
		if err := output.Close(); err != nil {
			rt.Error(err)
		}
	}()
	if err := gob.NewEncoder(output).Encode(m); err != nil {
		rt.Error(err)
	}
}
