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
	"fmt"
	"io"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/e2e/benchmark/internal"
	"github.com/vdaas/vald/internal/log"
)

const (
	assetDir   = "../../assets"
	configDir  = assetDir + "/config/"
	datasetDir = assetDir + "/dataset/"
)

var (
	searchConfig = &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.01,
	}
	dataset = []string{
		"fashion-mnist-784-euclidean",
		"mnist-784-euclidean",
		// "sift-128-euclidean",
		// "nytimes-256-angular",
		// "glove-25-angular",
		// "glove-50-angular",
		// "glove-100-angular",
		// "glove-200-angular",
	}
)

func init() {
	log.Init(log.DefaultGlg())
}

func BenchmarkAgentNGTRESTSequential(b *testing.B) {
}

func BenchmarkAgentNGTgRPCSequential(rb *testing.B) {
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()
	for _, name := range dataset {
		rb.Run(name, func (b *testing.B){
			b.ReportAllocs()
			b.ResetTimer()
			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			internal.StartAgentNGTServer(b, ctx, configDir+name+".yaml")

			ids, train, test := internal.LoadDataAndIDs(b, datasetDir+name+".hdf5")

			client := internal.NewAgentClient(b, ctx, "localhost", 8082)

			b.Run(fmt.Sprintf("Insert %d objects", len(train)), func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for i, vector := range train {
					_, err := client.Insert(ctx, &payload.Object_Vector{
						Id: &payload.Object_ID{
							Id: ids[i],
						},
						Vector: vector,
					})
					if err != nil {
						bb.Error(err)
					}
				}
			})

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
			b.Run(fmt.Sprintf("StreamSearch %d objects", len(test)), func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for _, data := range test {
					_, err := client.Search(ctx, &payload.Search_Request{
						Vector: &payload.Object_Vector{
							Vector: data,
						},
						Config: searchConfig,
					})
					if err != nil {
						bb.Error(err)
					}
				}
			})
			b.Run(fmt.Sprintf("StreamRemove %d objects", len(ids)/2), func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for _, id := range ids[:len(ids)/2] {
					_, err := client.Remove(ctx, &payload.Object_ID{
						Id: id,
					})
					if err != nil {
						bb.Error(err)
					}
				}
			})

		})
	}
}

func BenchmarkAgentNGTgRPCStream(rb *testing.B) {
	rctx, rcancel := context.WithCancel(context.Background())
	defer rcancel()
	rb.ReportAllocs()
	rb.ResetTimer()
	for _, name := range dataset {
		rb.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			ctx, cancel := context.WithCancel(rctx)
			defer cancel()

			internal.StartAgentNGTServer(b, ctx, configDir+name+".yaml")

			ids, train, test := internal.LoadDataAndIDs(b, datasetDir+name+".hdf5")

			client := internal.NewAgentClient(b, ctx, "localhost", 8082)

			sti, err := client.StreamInsert(ctx)
			if err != nil {
				b.Error(err)
			}
			b.Run(fmt.Sprintf("StreamInsert %d objects", len(train)), func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for i, data := range train {
					err := sti.Send(&payload.Object_Vector{
						Id: &payload.Object_ID{
							Id: ids[i],
						},
						Vector: data,
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
				}
			})
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

			b.Run(fmt.Sprintf("StreamSearch %d objects", len(test)), func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
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
						bb.Error(err)
					}
					_, err = st.Recv()
					if err != nil {
						if err != io.EOF {
							bb.Error(err)
						}
					}
				}
			})
			if err := st.CloseSend(); err != nil {
				b.Error(err)
			}

			str, err := client.StreamRemove(ctx)
			if err != nil {
				b.Error(err)
			}

			b.Run(fmt.Sprintf("StreamRemove %d objects", len(ids)/2), func(bb *testing.B) {
				bb.ReportAllocs()
				bb.ResetTimer()
				for _, id := range ids[:len(ids)/2] {
					err := str.Send(&payload.Object_ID{
						Id: id,
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
				}
			})
			if err := str.CloseSend(); err != nil {
				b.Error(err)
			}
		})
	}
}
