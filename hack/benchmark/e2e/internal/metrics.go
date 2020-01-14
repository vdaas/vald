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
package internal

import (
	"context"
	"fmt"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
	"time"
)

func sum(x []float64) float64 {
	s := 0.0
	for _, a := range x {
		s += a
	}
	return s
}

func mean(x []float64) float64 {
	return sum(x) / float64(len(x))
}

func std(x []float64) float64 {
	x2 := make([]float64, len(x))
	for i, a := range x {
		x2[i] = a * a
	}
	m := mean(x)
	return mean(x2) - m*m
}

func Recall(datasetNeighbors, runNeighbors []string, k int) (r float64) {
	dn := map[string]struct{}{}
	for _, n := range datasetNeighbors[:k] {
		dn[n] = struct{}{}
	}
	for i := 0; i < k; i++ {
		if _, ok := dn[runNeighbors[i]]; ok {
			r++
		}
	}
	return r / float64(k)
}

func Recalls(datasetNeighbors, runNeighbors [][]string, k int) (recalls []float64) {
	recalls = make([]float64, len(runNeighbors))
	for i, d := range datasetNeighbors {
		recalls[i] = Recall(d, runNeighbors[i], k)
	}
	return recalls
}

func MeanStdRecalls(datasetNeighbors, runNeighbors [][]string, k int) (float64, float64, []float64) {
	recalls := Recalls(datasetNeighbors, runNeighbors, k)
	return mean(recalls), std(recalls), recalls
}

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

func FromYaml(filepath string) (p *Params, err error) {
	input, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = input.Close()
	}()
	if err := yaml.NewDecoder(input).Decode(&p); err != nil {
		return nil, err
	}
	return p, nil
}

type Searcher interface {
	Search(context.Context, *payload.Search_Request, ...grpc.CallOption) (*payload.Search_Response, error)
}

func SearchMetrics(t *testing.T, s Searcher, ctx context.Context, query [][]float64, neighbors [][]string, k uint32, eps, radius float32) (recall float64, qps float64) {
	t.Helper()

	c, cancel := context.WithCancel(ctx)
	defer cancel()
	querySize := len(query)
	results := make([][]string, querySize)
	var elapsed time.Duration = 0.0
	for i, v := range query {
		start := time.Now()
		resp, err := s.Search(c, &payload.Search_Request{
			Vector: v,
			Config: &payload.Search_Config{
				Epsilon:eps,
				Num:k,
				Radius:radius,
			},
		})
		elapsed += time.Since(start)
		if err != nil {
			t.Error(err)
		}
		results[i] = make([]string, len(resp.Results))
		for j, r := range resp.Results {
			results[i][j] = r.Id
		}
	}
	recall, _, _ = MeanStdRecalls(neighbors, results, k)
	qps = 1 / (elapsed.Seconds() / float64(querySize))
	return recall, qps
}

type Inserter interface {
	Insert(context.Context, *payload.Object_Vector, ...grpc.CallOption) (*payload.Empty, error)
}

func InsertMetrics(t *testing.T, in Inserter, ctx context.Context, train [][]float64, id []string) time.Duration {
	t.Helper()

	c, cancel := context.WithCancel(ctx)
	defer cancel()

	start := time.Now()
	for i, v := range train {
		_, err := in.Insert(c, &payload.Object_Vector{
			Id:     id[i],
			Vector: v,
		})
		if err != nil {
			t.Error(err)
		}
	}
	return time.Since(start)
}