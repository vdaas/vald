// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package grpc

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/atomic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type benchStreamServer struct {
	vald.UnimplementedValdServer
}

func (s *benchStreamServer) StreamInsert(stream vald.Insert_StreamInsertServer) error {
	return BidirectionalStream(stream.Context(), stream, 10,
		func(ctx context.Context, req *payload.Insert_Request) (*payload.Object_StreamLocation, error) {
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: &payload.Object_Location{
						Uuid: req.GetVector().GetId(),
					},
				},
			}, nil
		},
	)
}

func BenchmarkBidirectionalStream(b *testing.B) {
	const bufSize = 4 * 1024 * 1024
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	vald.RegisterValdServer(s, &benchStreamServer{})
	go func() {
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			b.Error(err)
		}
	}()
	b.Cleanup(func() {
		s.Stop()
		if err := lis.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			b.Error(err)
		}
	})

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		if err := conn.Close(); err != nil {
			b.Error(err)
		}
	})

	client := vald.NewValdClient(conn)

	concurrencies := []int{1, 10}
	dims := []int{128, 512}

	for _, concurrency := range concurrencies {
		for _, dim := range dims {
			vec := make([]float32, dim)
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     "bench-vec",
					Vector: vec,
				},
			}

			b.Run(fmt.Sprintf("Concurrency_%d_Dim_%d", concurrency, dim), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					ctx := b.Context()
					stream, err := client.StreamInsert(ctx)
					if err != nil {
						b.Fatal(err)
					}

					var count atomic.Int64
					const batchSize = 100
					var sent int
					provider := func() (*payload.Insert_Request, bool) {
						if sent >= batchSize {
							return nil, false
						}
						sent++
						return req, true
					}
					callBack := func(res *payload.Object_StreamLocation, err error) bool {
						if err != nil {
							return false
						}
						count.Add(1)
						return true
					}

					err = BidirectionalStreamClient(stream, concurrency, provider, callBack)
					if err != nil && !errors.Is(err, io.EOF) {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
