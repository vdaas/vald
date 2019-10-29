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
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/e2e/benchmark/internal"
	"github.com/vdaas/vald/hack/e2e/benchmark/internal/dataset"
	"github.com/vdaas/vald/internal/log"
)

var (
	searchConfig = &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.01,
	}
	targets []string
	datasetVar string
	once sync.Once
)

func init() {
	log.Init(log.DefaultGlg())

	datasetList := make([]string, 0, len(dataset.Data))
	for key := range dataset.Data {
		datasetList = append(datasetList, "\t"+key)
	}
	sort.Strings(datasetList)
	flag.StringVar(&datasetVar, "dataset", "", "available dataset(choice with comma)\n" + strings.Join(datasetList, "\n"))
}

func parseArgs() {
	once.Do(func() {
		flag.Parse()
		targets = strings.Split(strings.TrimSpace(datasetVar), ",")
	})
}

func BenchmarkAgentNGTRESTSequential(rb *testing.B) {
	parseArgs()
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()

	for _, name := range targets {
		if name == "" {
			continue
		}

		rb.Run(name, func(b *testing.B) {
			data := dataset.Data[name](rb)
			ids := data.IDs()
			train := data.Train()
			query := data.Query()

			b.ReportAllocs()
			b.ResetTimer()
			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			internal.StartAgentNGTServer(b, ctx, data)

			buffers := make([]*bytes.Buffer, len(train))
			for i := 0; i < len(train); i++ {
				buf, err := json.Marshal(&payload.Object_Vector{
					Id: &payload.Object_ID{
						Id: ids[i],
					},
					Vector: train[i],
				})
				if err != nil {
					b.Error(err)
				}
				buffers[i] = bytes.NewBuffer(buf)
			}

			i := 0
			b.Run("Insert objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					resp, err := http.Post("http://localhost:8081/insert", "application/json", buffers[i])
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
				resp, err := http.Post("http://localhost:8081/insert", "application/json", buffers[i])
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

			b.Run("CreateIndex", func(bb *testing.B) {
				buf, err := json.Marshal(&payload.Controll_CreateIndexRequest{
					PoolSize: 10000,
				})
				if err != nil {
					bb.Error(err)
				}
				buffer := bytes.NewBuffer(buf)
				bb.ReportAllocs()
				bb.ResetTimer()
				resp, err := http.Post("http://localhost:8081/index/create", "application/json", buffer)
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
					Vector: &payload.Object_Vector{
						Vector: query[i],
					},
					Config: searchConfig,
				})
				if err != nil {
					b.Error(err)
				}
				buffers[i] = bytes.NewBuffer(buf)
			}

			i = 0
			b.Run("Search objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					resp, err := http.Post("http://localhost:8081/search", "application/json", buffers[i])
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
			b.Run("Remove objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					resp, err := http.Post("http://localhost:8081/remove", "application/json", buffers[i])
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
	parseArgs()
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()
	for _, name := range targets {
		if name == "" {
			continue
		}
		rb.Run(name, func (b *testing.B){
			data := dataset.Data[name](rb)
			if data == nil {
				b.Logf("dataset %s is nil", name)
				return
			}
			ids := data.IDs()
			train := data.Train()
			query := data.Query()

			b.ReportAllocs()
			b.ResetTimer()
			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			internal.StartAgentNGTServer(b, ctx, data)

			client := internal.NewAgentClient(b, ctx, "localhost", 8082)

			i := 0
			b.Run("Insert objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					_, err := client.Insert(ctx, &payload.Object_Vector{
						Id: &payload.Object_ID{
							Id: ids[i],
						},
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
					Id: &payload.Object_ID{
						Id: ids[i],
					},
					Vector: train[i],
				})
				if err != nil {
					b.Error(err)
				}
			}

			b.Run("CreateIndex", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				_, err := client.CreateIndex(ctx, &payload.Controll_CreateIndexRequest{
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
						Vector: &payload.Object_Vector{
							Vector: query[i],
						},
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

func BenchmarkAgentNGTgRPCStream(rb *testing.B) {
	parseArgs()
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()
	for _, name := range targets {
		if name == "" {
			continue
		}
		rb.Run(name, func(b *testing.B) {
			data := dataset.Data[name](rb)
			if data == nil {
				b.Logf("dataset %s is nil", name)
				return
			}
			ids := data.IDs()
			train := data.Train()
			query := data.Query()

			b.ReportAllocs()
			b.ResetTimer()
			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			internal.StartAgentNGTServer(b, ctx, data)

			client := internal.NewAgentClient(b, ctx, "localhost", 8082)

			sti, err := client.StreamInsert(ctx)
			if err != nil {
				b.Error(err)
			}
			i := 0
			b.Run("Insert objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					err := sti.Send(&payload.Object_Vector{
						Id: &payload.Object_ID{
							Id: ids[i],
						},
						Vector: train[i],
					})
					if err != nil {
						if err == io.EOF {
							log.Error(err)
							return
						}
						bb.Error(err)
					}
					_, err = sti.Recv()
					if err != nil {
						if err != io.EOF {
							bb.Error(err)
						}
					}
					i++
				}
			})
			for ; i < len(train); i++ {
				err := sti.Send(&payload.Object_Vector{
					Id: &payload.Object_ID{
						Id: ids[i],
					},
					Vector: train[i],
				})
				if err != nil {
					if err == io.EOF {
						log.Error(err)
						return
					}
					b.Error(err)
				}
				_, err = sti.Recv()
				if err != nil {
					if err != io.EOF {
						b.Error(err)
					}
				}
			}
			if err := sti.CloseSend(); err != nil {
				b.Error(err)
			}
			b.Run("CreateIndex", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				_, err := client.CreateIndex(ctx, &payload.Controll_CreateIndexRequest{
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
						Vector: &payload.Object_Vector{
							Vector: query[i],
						},
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
							bb.Error(err)
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
						Id: ids[i],
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
