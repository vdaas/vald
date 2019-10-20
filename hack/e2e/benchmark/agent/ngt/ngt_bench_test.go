//
// Copyright (C) 2019 kpango (Yusuke Kato)
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
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/e2e/benchmark/internal"
	"github.com/vdaas/vald/internal/log"
)

const (
	assetDir   = "../../assets"
	configDir  = assetDir + "/config/"
	datasetDir = assetDir + "/dataset/"

	fashionMnistConfig  = configDir + "fashion-mnist-784-euclidean.yaml"
	fashionMnistDataSet = datasetDir + "fashion-mnist-784-euclidean.hdf5"
)

var (
	searchConfig = &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.01,
	}
)

func init() {
	log.Init(log.DefaultGlg())

}

func BenchmarkAgentNGTSequentialInsert(b *testing.B) {

	ids, train, test := internal.LoadDataAndIDs(b, fashionMnistDataSet)
	internal.Insert(b, ids, train, func(id string, vector []float64) error {
		_, err := client.Insert(context.Background(), &payload.Object_Vector{
			Id: &payload.Object_ID{
				Id: id,
			},
			Vector: vector,
		})
		return err
	})
}

func BenchmarkAgentNGTSequentialCreateIndex(b *testing.B) {
	internal.CreateIndex(b, func() error {
		_, err := client.CreateIndex(context.Background(), &payload.Controll_CreateIndexRequest{
			PoolSize: 10000,
		})
		return err
	})
}

func BenchmarkAgentNGTSequentialSearch(b *testing.B) {
	ids, train, test := internal.LoadDataAndIDs(b, fashionMnistDataSet)
	internal.Search(b, test, func(vector []float64) error {
		_, err := client.Search(context.Background(), &payload.Search_Request{
			Vector: &payload.Object_Vector{
				Vector: vector,
			},
			Config: searchConfig,
		})
		return err
	})
}

func BenchmarkAgentNGTSequentialRemove(b *testing.B) {
	ids, train, test := internal.LoadDataAndIDs(b, fashionMnistDataSet)
	internal.Remove(b, ids, func(id string) error {
		_, err := client.Remove(context.Background(), &payload.Object_ID{
			Id: id,
		})
		return err
	})
}

func BenchmarkAgentNGTStreamInsert(b *testing.B) {
	ids, train, test := internal.LoadDataAndIDs(b, fashionMnistDataSet)
	rctx, rcancel := context.WithCancel(context.Background())
	client := internal.NewAgentClient(b, rctx, "localhost", 8082)
	st, err := client.StreamInsert(context.Background())
	if err != nil {
		b.Error(err)
	}
	go func() {
		for {
			_, err := st.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				if !strings.Contains(err.Error(), "already exists") {
					b.Error(err)
				}
			}
		}
	}()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		idx := 0
		for pb.Next() {
			err := st.Send(&payload.Object_Vector{
				Id: &payload.Object_ID{
					Id: ids[idx],
				},
				Vector: train[idx],
			})
			if err != nil && err != io.EOF {
				if !strings.Contains(err.Error(), "already exists") {
					b.Error(err)
				}
			}
			if idx >= len(ids)-1 {
				idx = 0
			} else {
				idx++
			}
		}
	})
	if err := st.CloseSend(); err != nil {
		b.Error(err)
	}
}

func BenchmarkAgentNGTStreamCreateIndex(b *testing.B) {
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	client := internal.NewAgentClient(b, rctx, "localhost", 8082)
	internal.CreateIndex(b, func() error {
		_, err := client.CreateIndex(context.Background(), &payload.Controll_CreateIndexRequest{
			PoolSize: 10000,
		})
		return err
	})
}

func BenchmarkAgentNGTStreamSearch(b *testing.B) {
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()

	internal.StartAgentNGTServer(b, rctx, fashionMnistConfig)

	client := internal.NewAgentClient(b, rctx, "localhost", 8082)

	sti, err := client.StreamInsert(context.Background())
	if err != nil {
		b.Error(err)
	}

	ctx, cancel := context.WithCancel(rctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				ctx, cancel = context.WithCancel(rctx)
				_, err = client.CreateIndex(ctx, &payload.Controll_CreateIndexRequest{
					PoolSize: 10000,
				})
				if err != nil {
					b.Error(err)
				}
				cancel()
				return
			default:
				_, err := sti.Recv()
				if err != nil {
					if err == io.EOF {
						cancel()
						continue
					}
					if !strings.Contains(err.Error(), "already exists") {
						b.Error(err)
					}
				}
			}
		}
	}()

	ids, train, test := internal.LoadDataAndIDs(b, fashionMnistDataSet)

	for i, data := range train {
		err := sti.Send(&payload.Object_Vector{
			Id: &payload.Object_ID{
				Id: ids[i],
			},
			Vector: data,
		})
		if err != nil {
			if err != io.EOF {
				b.Error(err)
			}
		}
	}
	if err := sti.CloseSend(); err != nil {
		b.Error(err)
	}
	cancel()

	st, err := client.StreamSearch(context.Background())
	if err != nil {
		b.Error(err)
	}

	wg.Wait()
	b.ReportAllocs()
	b.ResetTimer()
	b.Run("StreamSearch", func(bb *testing.B) {
		for _, data := range test {
			err := st.Send(&payload.Search_Request{
				Vector: &payload.Object_Vector{
					Vector: data,
				},
				Config: searchConfig,
			})
			if err != nil {
				if err == io.EOF {
					return
				}
				b.Error(err)
			}
			_, err = st.Recv()
			if err != nil {
				if err != io.EOF {
					b.Error(err)
				}
			}
		}
	})
	if err := st.CloseSend(); err != nil {
		b.Error(err)
	}
}

func BenchmarkAgentNGTStreamRemove(b *testing.B) {
	st, err := client.StreamRemove(context.Background())
	if err != nil {
		b.Error(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(st agent.Agent_StreamRemoveClient, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			_, err := st.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				b.Error(err)
			}
		}
	}(st, wg)
	internal.Remove(b, ids, func(id string) error {
		err := st.Send(&payload.Object_ID{
			Id: id,
		})
		if err == io.EOF {
			return nil
		}
		return err
	})
	if err := st.CloseSend(); err != nil {
		b.Error(err)
	}

	wg.Wait()
}
