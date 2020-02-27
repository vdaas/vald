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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/log"
)

func BenchmarkAgentNGTRESTSequential(rb *testing.B) {
	parseArgs(rb)
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()

	for N, name := range targets {
		address := addresses[N]
		if address == "" {
			address = "localhost:8081"
		}
		if name == "" {
			continue
		}

		rb.Run(name, func(b *testing.B) {
			data := assets.Data(name)(rb)
			ids := data.IDs()
			train := data.Train()
			query := data.Query()

			b.ReportAllocs()
			b.ResetTimer()
			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			if strings.Contains(address, "localhost") {
				internal.StartAgentNGTServer(b, ctx, data)
			}

			buffers := make([]*bytes.Buffer, len(train))
			for i := 0; i < len(train); i++ {
				buf, err := json.Marshal(&payload.Object_Vector{
					Id:     ids[i],
					Vector: train[i],
				})
				if err != nil {
					b.Error(err)
				}
				buffers[i] = bytes.NewBuffer(buf)
			}

			i := 0
			url := fmt.Sprintf("http://%s/insert", address)
			b.Run("Insert objects", func(bb *testing.B) {
				for bb.N+i >= len(ids) {
					ids = append(ids, assets.CreateSequentialIDs(len(train))...)
				}

				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					resp, err := http.Post(url, "application/json", buffers[i])
					if err != nil {
						bb.Error(err)
					}
					_, err = io.Copy(ioutil.Discard, resp.Body)
					if err != nil {
						bb.Error(err)
					}
					err = resp.Body.Close()
					if err != nil {
						bb.Error(err)
					}

					i++
				}
			})
			for ; i < len(train); i++ {
				resp, err := http.Post(url, "application/json", buffers[i])
				if err != nil {
					b.Error(err)
				}
				_, err = io.Copy(ioutil.Discard, resp.Body)
				if err != nil {
					b.Error(err)
				}
				err = resp.Body.Close()
				if err != nil {
					b.Error(err)
				}
			}

			url = fmt.Sprintf("http://%s/index/create", address)
			b.Run("CreateIndex", func(bb *testing.B) {
				buf, err := json.Marshal(&payload.Control_CreateIndexRequest{
					PoolSize: 10000,
				})
				if err != nil {
					bb.Error(err)
				}
				buffer := bytes.NewBuffer(buf)
				bb.ReportAllocs()
				bb.ResetTimer()
				resp, err := http.Post(url, "application/json", buffer)
				if err != nil {
					bb.Error(err)
				}
				_, err = io.Copy(ioutil.Discard, resp.Body)
				if err != nil {
					bb.Error(err)
				}
				err = resp.Body.Close()
				if err != nil {
					bb.Error(err)
				}
			})

			buffers = make([]*bytes.Buffer, len(query))
			for i := 0; i < len(query); i++ {
				buf, err := json.Marshal(&payload.Search_Request{
					Vector: query[i],
					Config: searchConfig,
				})
				if err != nil {
					b.Error(err)
				}
				buffers[i] = bytes.NewBuffer(buf)
			}

			i = 0
			url = fmt.Sprintf("http://%s/search", address)
			b.Run("Search objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					resp, err := http.Post(url, "application/json", buffers[i])
					if err != nil {
						bb.Error(err)
					}
					_, err = io.Copy(ioutil.Discard, resp.Body)
					if err != nil {
						bb.Error(err)
					}
					err = resp.Body.Close()
					if err != nil {
						bb.Error(err)
					}

					i++
				}
			})

			buffers = make([]*bytes.Buffer, len(ids))
			for i := 0; i < len(ids); i++ {
				buf, err := json.Marshal(&payload.Object_ID{
					Id: ids[i],
				})
				if err != nil {
					b.Error(err)
				}
				buffers[i] = bytes.NewBuffer(buf)
			}

			i = 0
			url = fmt.Sprintf("http://%s/remove", address)
			b.Run("Remove objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					resp, err := http.Post(url, "application/json", buffers[i])
					if err != nil {
						bb.Error(err)
					}
					_, err = io.Copy(ioutil.Discard, resp.Body)
					if err != nil {
						bb.Error(err)
					}
					err = resp.Body.Close()
					if err != nil {
						bb.Error(err)
					}

					i++
				}
			})
		})
	}
}

func BenchmarkAgentNGTgRPCSequential(rb *testing.B) {
	parseArgs(rb)
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()
	for N, name := range targets {
		address := addresses[N]
		if address == "" {
			address = "localhost:8082"
		}

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

			if strings.Contains(address, "localhost") {
				internal.StartAgentNGTServer(b, ctx, data)
			}

			client := internal.NewAgentClient(b, ctx, address)

			i := 0
			b.Run("Insert objects", func(bb *testing.B) {
				for bb.N+i >= len(ids) {
					ids = append(ids, assets.CreateSequentialIDs(len(train))...)
				}

				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					_, err := client.Insert(ctx, &payload.Object_Vector{
						Id:     ids[i],
						Vector: train[i%len(train)],
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

			i = 0
			b.Run("Search objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					_, err := client.Search(ctx, &payload.Search_Request{
						Vector: query[i%len(query)],
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
						Id: ids[i%len(ids)],
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

func BenchmarkAgentNGTgRPCStream(rb *testing.B) {
	parseArgs(rb)
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()
	for N, name := range targets {
		address := addresses[N]
		if address == "" {
			address = "localhost:8082"
		}

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

			if strings.Contains(address, "localhost") ||
				strings.Contains(address, "127.0.0.1") ||
				strings.Contains(address, "0.0.0.0") {
				internal.StartAgentNGTServer(b, ctx, data)
			}

			client := internal.NewAgentClient(b, ctx, address)

			sti, err := client.StreamInsert(ctx)
			if err != nil {
				b.Error(err)
			}
			wg := &sync.WaitGroup{}
			i := 0
			b.Run("Insert objects", func(bb *testing.B) {
				for bb.N+i >= len(ids) {
					ids = append(ids, assets.CreateSequentialIDs(len(train))...)
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
