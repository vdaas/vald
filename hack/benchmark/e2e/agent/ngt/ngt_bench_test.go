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
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/benchmark/e2e/internal"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e/strategy"
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
				buf, err := json.Marshal(&payload.Controll_CreateIndexRequest{
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

func BenchmarkAgentNGT_gRPC_Sequential(b *testing.B) {
	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(nil),
			e2e.WithStrategy(
				strategy.NewInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(nil),
					strategy.WithCreateIndexPoolSize(10000),
				),
				strategy.NewSearch(
					strategy.WithSearchConfig(searchConfig),
				),
				strategy.NewRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}

func BenchmarkAgentNGT_gRPC_Stream(b *testing.B) {
	for _, name := range targets {
		bench := e2e.New(
			b,
			e2e.WithName(name),
			e2e.WithClient(nil),
			e2e.WithStrategy(
				strategy.NewStreamInsert(),
				strategy.NewCreateIndex(
					strategy.WithCreateIndexClient(nil),
					strategy.WithCreateIndexPoolSize(10000),
				),
				strategy.NewStreamSearch(
					strategy.WithStreamSearchConfig(searchConfig),
				),
				strategy.NewStreamRemove(),
			),
		)
		bench.Run(context.Background(), b)
	}
}
