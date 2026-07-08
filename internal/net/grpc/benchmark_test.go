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
	"fmt"
	"net"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/sync/errgroup"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type echoServer struct {
	vald.UnimplementedValdServer
}

func (s *echoServer) Insert(
	ctx context.Context, req *payload.Insert_Request,
) (*payload.Object_Location, error) {
	return &payload.Object_Location{}, nil
}

func startBroadcastServer(t testing.TB, eg errgroup.Group) (net.Listener, func()) {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	srv := ggrpc.NewServer()
	vald.RegisterValdServer(srv, &echoServer{})

	eg.Go(func() error {
		return srv.Serve(lis)
	})

	return lis, func() {
		srv.Stop()
	}
}

func startEchoServer(t testing.TB, eg errgroup.Group) (string, func()) {
	t.Helper()

	// Get available port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	addr := lis.Addr().(*net.TCPAddr)
	port := addr.Port
	lis.Close()

	srv, err := server.New(
		server.WithServerMode(server.GRPC),
		server.WithHost("127.0.0.1"),
		server.WithPort(uint16(port)),
		server.WithNetwork("tcp"),
		server.WithGRPCRegisterar(func(s *ggrpc.Server) {
			vald.RegisterValdServer(s, &echoServer{})
		}),
	)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	errCh := make(chan error, 1)
	eg.Go(func() error {
		return srv.ListenAndServe(t.Context(), errCh)
	})

	return fmt.Sprintf("127.0.0.1:%d", port), func() {
		srv.Shutdown(t.Context())
	}
}

func BenchmarkBroadcast(b *testing.B) {
	// Number of agents to simulate
	const numAgents = 30
	const concurrency = 2

	eg := errgroup.Get()

	// Start servers
	listeners := make([]net.Listener, numAgents)
	closers := make([]func(), numAgents)
	for i := range numAgents {
		lis, closeFunc := startBroadcastServer(b, eg)
		listeners[i] = lis
		closers[i] = closeFunc
	}
	defer func() {
		for _, c := range closers {
			c()
		}
	}()

	// multi-listener dialer
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		// parse addr to get index
		var idx int
		fmt.Sscanf(addr, "bufconn-%d", &idx)
		return listeners[idx].(*bufconn.Listener).Dial()
	}

	// generate dummy addrs
	addrs := make([]string, numAgents)
	for i := range numAgents {
		addrs[i] = fmt.Sprintf("bufconn-%d", i)
	}

	// Initialize Client
	client := grpc.New(
		"benchmark-client",
		grpc.WithAddrs(addrs...),
		grpc.WithDialOptions(
			ggrpc.WithTransportCredentials(insecure.NewCredentials()),
			ggrpc.WithContextDialer(dialer),
		),
		grpc.WithConnectionPoolSize(1),
	)
	defer client.Close(b.Context())

	_, err := client.StartConnectionMonitor(b.Context())
	if err != nil {
		b.Fatal(err)
	}
	req := &payload.Insert_Request{}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		err := client.RangeConcurrent(b.Context(), concurrency, func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			c := vald.NewValdClient(conn)
			_, err := c.Insert(ctx, req)
			return err
		})
		if err != nil {
			b.Fatalf("Broadcast failed: %v", err)
		}
	}
	b.StopTimer()
	// NOT IMPLEMENTED BELOW
}

func BenchmarkExecuteRPC(b *testing.B) {
	eg := errgroup.Get()
	addr, closer := startEchoServer(b, eg)
	defer closer()

	client := grpc.New(
		"benchmark-client",
		grpc.WithAddrs(addr),
		grpc.WithDialOptions(
			ggrpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		grpc.WithConnectionPoolSize(2),
	)
	defer client.Close(b.Context())

	_, err := client.StartConnectionMonitor(b.Context())
	if err != nil {
		b.Fatal(err)
	}
	req := &payload.Insert_Request{}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, err := client.Do(b.Context(), addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
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

	client := grpc.New(
		"benchmark-client",
		grpc.WithAddrs(addr),
		grpc.WithDialOptions(
			ggrpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		grpc.WithConnectionPoolSize(10),
	)
	defer client.Close(b.Context())

	_, err := client.StartConnectionMonitor(b.Context())
	if err != nil {
		b.Fatal(err)
	}

	req := &payload.Insert_Request{}

	b.ResetTimer()
	b.ReportAllocs()

	b.SetParallelism(1000) // Simulate high load

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.Do(b.Context(), addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
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
