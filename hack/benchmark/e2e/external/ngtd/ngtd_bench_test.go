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
package ngtd

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/log"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd"
	"github.com/yahoojapan/ngtd/kvs"
	"github.com/yahoojapan/ngtd/model"
	proto "github.com/yahoojapan/ngtd/proto"
	"google.golang.org/grpc"
)

const (
	baseDir = "/tmp/ngtd/"
	port    = 8200
)

var (
	targets    []string
	datasetVar string
	once       sync.Once
)

func init() {
	log.Init()
	glg.Get().SetMode(glg.NONE)
	if err := os.RemoveAll(baseDir); err != nil {
		log.Error(err)
	}
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Error(err)
	}

	flag.StringVar(&datasetVar, "assets", "", "list available assets(choice with comma)")
}

func parseArgs() {
	once.Do(func() {
		flag.Parse()
		targets = strings.Split(strings.TrimSpace(datasetVar), ",")
	})
}

func StartNGTD(tb testing.TB, t ngtd.ServerType, dim int) func() {
	tb.Helper()
	gongt.SetDimension(dim)
	db, err := kvs.NewGoLevel(baseDir + "meta")
	if err != nil {
		tb.Error(err)
	}
	n, err := ngtd.NewNGTD(baseDir+"ngt", db, port)
	if err != nil {
		tb.Error(err)
	}

	go func() {
		err := n.ListenAndServe(t)
		if err != nil {
			tb.Errorf("ngtd returned error: %s", err.Error())
		}
	}()

	time.Sleep(5 * time.Second)

	return func() {
		n.Stop()
		if err := os.RemoveAll(baseDir + "meta"); err != nil {
			tb.Error(err)
		}
		if err := os.RemoveAll(baseDir + "ngt"); err != nil {
			tb.Error(err)
		}
	}
}

func BenchmarkNGTDRESTSequential(rb *testing.B) {
	parseArgs()
	rb.ReportAllocs()
	rb.ResetTimer()

	for _, name := range targets {
		if name == "" {
			continue
		}

		rb.Run(name, func(b *testing.B) {
			data := assets.Data(name)(rb)
			ids := data.IDs()
			train := data.Train()
			query := data.Query()

			defer StartNGTD(b, ngtd.HTTP, data.Dimension())()

			buffers := make([]*bytes.Buffer, len(train))
			for i := 0; i < len(train); i++ {
				buf, err := json.Marshal(&model.InsertRequest{
					Vector: train[i],
					ID:     ids[i],
				})
				if err != nil {
					b.Error(err)
				}
				buffers[i] = bytes.NewBuffer(buf)
			}

			i := 0
			url := fmt.Sprintf("http://localhost:%d/insert", port)
			b.ReportAllocs()
			b.ResetTimer()

			b.Run("Insert objects", func(bb *testing.B) {
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

			b.Run("CreateIndex", func(bb *testing.B) {
				url := fmt.Sprintf("http://localhost:%d/index/create/10000", port)
				bb.ReportAllocs()
				bb.ResetTimer()
				resp, err := http.Get(url)
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
				buf, err := json.Marshal(&model.SearchRequest{
					ID:      ids[i],
					Vector:  train[i],
					Size:    10,
					Epsilon: 0.01,
				})
				if err != nil {
					b.Error(err)
				}
				buffers[i] = bytes.NewBuffer(buf)
			}

			i = 0
			b.Run("Search objects", func(bb *testing.B) {
				url := fmt.Sprintf("http://localhost:%d/search", port)
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

			i = 0
			b.Run("Remove objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					resp, err := http.Get(fmt.Sprintf("http://localhost:%d/remove/%s", port, ids[i]))
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

func BenchmarkNGTDgRPCSequential(rb *testing.B) {
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
			data := assets.Data(name)(rb)
			if data == nil {
				b.Logf("assets %s is nil", name)
				return
			}
			ids := data.IDs()
			train := data.Train()
			query := data.Query()

			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			defer StartNGTD(b, ngtd.GRPC, data.Dimension())()

			conn, err := grpc.DialContext(ctx, fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
			if err != nil {
				b.Error(err)
			}

			client := proto.NewNGTDClient(conn)

			b.ReportAllocs()
			b.ResetTimer()
			i := 0
			b.Run("Insert objects", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for n := 0; n < bb.N; n++ {
					_, err := client.Insert(ctx, &proto.InsertRequest{
						Id:     []byte(ids[i]),
						Vector: train[i],
					})
					if err != nil {
						bb.Error(err)
					}
					i++
				}
			})
			for ; i < len(train); i++ {
				_, err := client.Insert(ctx, &proto.InsertRequest{
					Id:     []byte(ids[i]),
					Vector: train[i],
				})
				if err != nil {
					b.Error(err)
				}
			}

			b.Run("CreateIndex", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				_, err := client.CreateIndex(ctx, &proto.CreateIndexRequest{
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
					_, err := client.Search(ctx, &proto.SearchRequest{
						Vector:  query[i],
						Epsilon: 0.01,
						Size_:   10,
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
					_, err := client.Remove(ctx, &proto.RemoveRequest{
						Id: []byte(ids[i]),
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

func BenchmarkNGTDgRPCStream(rb *testing.B) {
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
			data := assets.Data(name)(rb)
			if data == nil {
				b.Logf("assets %s is nil", name)
				return
			}
			ids := data.IDs()
			train := data.Train()
			query := data.Query()

			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			defer StartNGTD(b, ngtd.GRPC, data.Dimension())()

			conn, err := grpc.DialContext(ctx, fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
			if err != nil {
				b.Error(err)
			}

			client := proto.NewNGTDClient(conn)
			b.ReportAllocs()
			b.ResetTimer()

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
					err := sti.Send(&proto.InsertRequest{
						Id:     []byte(ids[i]),
						Vector: train[i%len(train)],
					})
					if err != nil {
						if err == io.EOF {
							break
						}
						bb.Error(err)
					}
					wg.Add(1)
					go func() {
						defer wg.Done()
						_, err := sti.Recv()
						if err != nil {
							if err != io.EOF {
								bb.Error(err)
							}
						}
					}()
					i++
				}
				log.Info(bb.N, i)
			})
			log.Info("wait")
			wg.Wait()
			log.Info("done")
			for ; i < len(train); i++ {
				err := sti.Send(&proto.InsertRequest{
					Id:     []byte(ids[i]),
					Vector: train[i%len(train)],
				})
				if err != nil {
					if err == io.EOF {
						break
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
			if err := sti.CloseSend(); err != nil {
				b.Error(err)
			}
			wg.Wait()
			b.Run("CreateIndex", func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				_, err := client.CreateIndex(ctx, &proto.CreateIndexRequest{
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
					err := st.Send(&proto.SearchRequest{
						Vector:  query[i%len(query)],
						Size_:   10,
						Epsilon: 0.01,
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
					err := str.Send(&proto.RemoveRequest{
						Id: []byte(ids[i%len(ids)]),
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
