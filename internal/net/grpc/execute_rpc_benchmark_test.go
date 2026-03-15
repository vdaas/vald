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
package grpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type echoServer struct {
	vald.UnimplementedValdServer
}

func (s *echoServer) Insert(
	ctx context.Context, req *payload.Insert_Request,
) (*payload.Object_Location, error) {
	return &payload.Object_Location{}, nil
}

func startEchoServer(t testing.TB, eg errgroup.Group) (string, func()) {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	srv := ggrpc.NewServer()
	vald.RegisterValdServer(srv, &echoServer{})

	eg.Go(func() error {
		return srv.Serve(lis)
	})

	return lis.Addr().String(), func() {
		srv.Stop()
	}
}

func BenchmarkExecuteRPC(b *testing.B) {
	eg := errgroup.Get()
	addr, closer := startEchoServer(b, eg)
	defer closer()

	ctx := context.Background()
	client := grpc.New(
		"benchmark-client",
		grpc.WithAddrs(addr),
		grpc.WithDialOptions(
			ggrpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		grpc.WithConnectionPoolSize(1),
	)
	defer client.Close(ctx)

	for range 10 {
		_, err := client.Connect(ctx, addr)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	req := &payload.Insert_Request{}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, err := client.Do(ctx, addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			c := vald.NewValdClient(conn)
			return c.Insert(ctx, req)
		})
		if err != nil {
			b.Fatalf("Do failed: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkExecuteRPCParallel(b *testing.B) {
	eg := errgroup.Get()
	addr, closer := startEchoServer(b, eg)
	defer closer()

	ctx := context.Background()
	client := grpc.New(
		"benchmark-client",
		grpc.WithAddrs(addr),
		grpc.WithDialOptions(
			ggrpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		grpc.WithConnectionPoolSize(1),
	)
	defer client.Close(ctx)

	for range 10 {
		_, err := client.Connect(ctx, addr)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	req := &payload.Insert_Request{}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.Do(ctx, addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
				c := vald.NewValdClient(conn)
				return c.Insert(ctx, req)
			})
			if err != nil {
				b.Errorf("Do failed: %v", err)
			}
		}
	})
	b.StopTimer()
}
