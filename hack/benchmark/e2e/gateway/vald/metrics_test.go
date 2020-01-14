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
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/hack/benchmark/e2e/internal"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/metrics"
)

var (
	configuration = flag.String("conf", "metrics.yaml", "set metrics configuration file path")
	outputPath    = flag.String("output", "metrics.gob", "set result output path")
	wait          = flag.Duration("wait-for-indexing", 30 * time.Second, "wait time for automatic indexing")
)

func TestMetrics(rt *testing.T) {
	flag.Parse()

	p, err := internal.FromYaml(*configuration)
	if err != nil {
		rt.Error(err)
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

						client := internal.NewValdClient(t, ctx, p.Address())

						buildTime := internal.InsertMetrics(tt, client, ctx, data.Train(), data.IDs())

						time.Sleep(*wait)

						for _, k := range p.K {
							searchMetrics := make([]*metrics.SearchMetrics, 0, len(p.Epsilon))
							for _, eps := range p.Epsilon {
								tt.Run(fmt.Sprintf("Recall@%d with %f", k, eps), func(ttt *testing.T) {
									recall, qps := internal.SearchMetrics(ttt, client, ctx, data.Query(), datasetNeighbors, k, eps, -1)
									searchMetrics = append(searchMetrics, &metrics.SearchMetrics{
										Recall:  recall,
										Qps:     qps,
										Epsilon: eps,
									})
									ttt.Logf("recall: %f, qps: %f", recall, qps)
								})
							}
							m = append(m, &metrics.Metrics{
								BuildTime:        buildTime,
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
	output, err := os.OpenFile(*outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
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
