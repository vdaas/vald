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
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/e2e/benchmark/internal"
	"github.com/vdaas/vald/internal/log"
	"google.golang.org/grpc"
)

var (
	train        [][]float64
	test         [][]float64
	ids          []string
	client       agent.AgentClient
	searchConfig = &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.01,
	}
)

func init() {
	log.Init(log.DefaultGlg())

	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client = agent.NewAgentClient(conn)

	datasetName := "../../assets/dataset/fashion-mnist-784-euclidean.hdf5"
	train, test, err = internal.Load(datasetName)
	if err != nil {
		log.Fatal(err)
	}

	ids = internal.CreateIDs(len(train))
}

func BenchmarkAgentNGTSequentialInsert(b *testing.B) {
	ids = internal.Insert(b, ids, train, func(id string, vector []float64) error {
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
	internal.Remove(b, ids, func(id string) error {
		_, err := client.Remove(context.Background(), &payload.Object_ID{
			Id: id,
		})
		return err
	})
}

func BenchmarkAgentNGTStreamInsert(b *testing.B) {
	st, err := client.StreamInsert(context.Background())
	if err != nil {
		b.Error(err)
	}
	go func(st agent.Agent_StreamInsertClient) {
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
	}(st)
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
	internal.CreateIndex(b, func() error {
		_, err := client.CreateIndex(context.Background(), &payload.Controll_CreateIndexRequest{
			PoolSize: 10000,
		})
		return err
	})
}

func BenchmarkAgentNGTStreamSearch(b *testing.B) {
	sti, err := client.StreamInsert(context.Background())
	if err != nil {
		b.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer log.Info("finish Insert")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, err := sti.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					if !strings.Contains(err.Error(), "already exists") {
						b.Error(err)
					}
				}
			}
		}
	}()

	log.Info("start Insert")
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
			// if !strings.Contains(err.Error(), "already exists") {
			// }
		}
	}

	if err := sti.CloseSend(); err != nil {
		b.Error(err)
	}
	cancel()

	log.Info("start Indexing")
	_, err = client.CreateIndex(context.Background(), &payload.Controll_CreateIndexRequest{
		PoolSize: 10000,
	})
	if err != nil {
		b.Error(err)
	}
	log.Info("finish Indexing")

	ctx, cancel = context.WithCancel(context.Background())
	st, err := client.StreamSearch(ctx)
	if err != nil {
		b.Error(err)
	}

	go func() {
		defer log.Info("finish Search")
		for {
			select {
			case <-ctx.Done():
				return
			default:
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
		}
	}()

	log.Info("start StreamSearch Benchmark")
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		flg := uint64(0)
		idx := 0
		for pb.Next() {
			if atomic.LoadUint64(&flg) > 0 {
				continue
			}
			err := st.Send(&payload.Search_Request{
				Vector: &payload.Object_Vector{
					Vector: train[idx],
				},
				Config: searchConfig,
			})
			if err != nil {
				if err == io.EOF {
					atomic.StoreUint64(&flg, 1)
				}
				// if !strings.Contains(err.Error(), "already exists") {
				b.Error(err)
				// }
			}
			if idx >= len(train)-1 {
				idx = 0
			} else {
				idx++
			}
		}
	})

	if err := st.CloseSend(); err != nil {
		b.Error(err)
	}
	cancel()
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
