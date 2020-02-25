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
	"flag"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/log"
)

var (
	searchConfig = &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.01,
	}
	targets    []string
	addresses  []string
	wait       time.Duration
	datasetVar string
	addressVar string
	once       sync.Once
	waitVar    int64
)

func init() {
	log.Init()

	flag.StringVar(&datasetVar, "dataset", "", "available dataset(choice with comma)")
	flag.StringVar(&addressVar, "address", "", "vald gateway address")
	flag.Int64Var(&waitVar, "wait", 30, "indexing wait time(secs)")
}

func parseArgs(tb testing.TB) {
	tb.Helper()
	once.Do(func() {
		flag.Parse()
		targets = strings.Split(strings.TrimSpace(datasetVar), ",")
		addresses = strings.Split(strings.TrimSpace(addressVar), ",")
		if len(targets) != len(addresses) {
			tb.Fatal("address and dataset must have same length.")
		}
		wait = time.Duration(waitVar) * time.Second
	})
}

func BenchmarkValdGatewaySequential(rb *testing.B) {
	parseArgs(rb)
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()
	for N, name := range targets {
		if name == "" {
			continue
		}
		rb.Run(name, func(b *testing.B) {
			data := assets.Data(name)(rb)
			if data == nil {
				b.Logf("assets %s is nil", name)
				return
			}
			ids := data.IDs()
			train := data.Train()
			query := data.Query()

			b.ReportAllocs()
			b.ResetTimer()
			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			client := internal.NewValdClient(b, ctx, addresses[N])

			i := 0
			b.Run("Insert objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					_, err := client.Insert(ctx, &payload.Object_Vector{
						Id:     ids[i],
						Vector: train[i],
					})
					if err != nil {
						bb.Error(err)
					}
					i++
				}
			})
			for ; i < len(train); i++ {
				_, err := client.Insert(ctx, &payload.Object_Vector{
					Id:     ids[i],
					Vector: train[i],
				})
				if err != nil {
					b.Error(err)
				}
			}

			time.Sleep(wait)

			i = 0
			b.Run("Search objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					_, err := client.Search(ctx, &payload.Search_Request{
						Vector: query[i],
						Config: searchConfig,
					})
					if err != nil {
						bb.Error(err)
					}
					i++
				}
			})

			i = 0
			b.Run("Remove objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					_, err := client.Remove(ctx, &payload.Object_ID{
						Id: ids[i],
					})
					if err != nil {
						bb.Error(err)
					}
					i++
				}
			})
		})
	}
}

func BenchmarkValdGatewayStream(rb *testing.B) {
	parseArgs(rb)
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()
	for N, name := range targets {
		if name == "" {
			continue
		}
		rb.Run(name, func(b *testing.B) {
			data := assets.Data(name)(rb)
			if data == nil {
				b.Fatalf("assets %s is nil", name)
			}
			b.Logf("benchmark %s", name)
			ids := data.IDs()
			train := data.Train()
			query := data.Query()

			b.ReportAllocs()
			b.ResetTimer()
			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			client := internal.NewAgentClient(b, ctx, addresses[N])

			sti, err := client.StreamInsert(ctx)
			if err != nil {
				b.Error(err)
			}
			wg := &sync.WaitGroup{}
			i := 0
			b.Run("Insert objects", func(bb *testing.B) {
				if bb.N+i >= len(ids) {
					ids = append(ids, assets.CreateRandomIDs(len(train))...)
				}

				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					err := sti.Send(&payload.Object_Vector{
						Id:     ids[i],
						Vector: train[i%len(train)],
					})
					if err != nil {
						if err == io.EOF {
							log.Error(err)
							return
						}
						bb.Error(err)
					}
					wg.Add(1)
					go func() {
						_, err := sti.Recv()
						if err != nil {
							if err != io.EOF {
								bb.Error(err)
							}
						}
						wg.Done()
					}()
					i++
				}
			})
			for ; i < len(train); i++ {
				err := sti.Send(&payload.Object_Vector{
					Id:     ids[i],
					Vector: train[i%len(train)],
				})
				if err != nil {
					if err == io.EOF {
						log.Error(err)
						return
					}
					b.Error(err)
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, err := sti.Recv()
					if err != nil {
						if err != io.EOF {
							b.Error(err)
						}
					}
				}()
			}
			wg.Wait()
			if err := sti.CloseSend(); err != nil {
				b.Error(err)
			}
			b.Run("CreateIndex", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				_, err := client.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
					PoolSize: 10000,
				})
				if err != nil {
					if err == io.EOF {
						return
					}
					bb.Error(err)
				}
			})
			st, err := client.StreamSearch(ctx)
			if err != nil {
				b.Error(err)
			}

			i = 0
			b.Run("Search objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					err := st.Send(&payload.Search_Request{
						Vector: query[i%len(query)],
						Config: searchConfig,
					})
					if err != nil {
						if err == io.EOF {
							return
						}
						bb.Error(err)
					}
					_, err = st.Recv()
					if err != nil {
						if err != io.EOF {
							bb.Error(i, err)
						}
					}
					i++
				}
			})
			if err := st.CloseSend(); err != nil {
				b.Error(err)
			}

			str, err := client.StreamRemove(ctx)
			if err != nil {
				b.Error(err)
			}

			i = 0
			b.Run("Remove objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					err := str.Send(&payload.Object_ID{
						Id: ids[i%len(ids)],
					})
					if err != nil {
						if err == io.EOF {
							return
						}
						bb.Error(err)
					}
					_, err = str.Recv()
					if err != nil {
						if err != io.EOF {
							bb.Error(err)
						}
					}
					i++
				}
			})
			if err := str.CloseSend(); err != nil {
				b.Error(err)
			}
		})
	}
}
