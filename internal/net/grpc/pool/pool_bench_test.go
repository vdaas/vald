// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package pool_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/servers"
	"github.com/vdaas/vald/internal/servers/server"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type discovererServer struct {
	discoverer.UnimplementedDiscovererServer
}

func (s *discovererServer) Nodes(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return &payload.Info_Nodes{
		Nodes: []*payload.Info_Node{
			{
				Name:         "node-1",
				InternalAddr: "127.0.0.1",
			},
		},
	}, nil
}

func startTestServer(t testing.TB, host string, port uint16) (servers.Listener, string) {
	t.Helper()
	addr := fmt.Sprintf("%s:%d", host, port)
	s, err := server.New(
		server.WithName("discoverer-grpc"),
		server.WithServerMode(server.GRPC),
		server.WithGRPCRegisterar(func(srv *ggrpc.Server) {
			discoverer.RegisterDiscovererServer(srv, &discovererServer{})
		}),
		server.WithHost(host),
		server.WithPort(port),
	)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	l := servers.New(
		servers.WithServer(s),
	)

	go l.ListenAndServe(context.Background())
	return l, addr
}

func Benchmark_ConnPool(b *testing.B) {
	host := "127.0.0.1"
	port := uint16(9093)
	l, addr := startTestServer(b, host, port)
	defer l.Shutdown(context.Background())
	time.Sleep(100 * time.Millisecond)

	ctx := context.Background()
	p, err := pool.New(ctx,
		pool.WithAddr(addr),
		pool.WithSize(4),
		pool.WithDialOptions(ggrpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		b.Fatal(err)
	}
	if _, err := p.Connect(ctx); err != nil {
		b.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		conn, ok := p.Get(ctx)
		if ok {
			client := discoverer.NewDiscovererClient(conn)
			_, err := client.Nodes(ctx, &payload.Discoverer_Request{})
			if err != nil {
				b.Errorf("RPC failed: %v", err)
			}
		} else {
			b.Error("failed to get connection from pool")
		}
	}
}

func BenchmarkParallel_ConnPool(b *testing.B) {
	host := "127.0.0.1"
	port := uint16(9094)
	l, addr := startTestServer(b, host, port)
	defer l.Shutdown(context.Background())
	time.Sleep(100 * time.Millisecond)

	ctx := context.Background()
	p, err := pool.New(ctx,
		pool.WithAddr(addr),
		pool.WithSize(uint64(runtime.GOMAXPROCS(0)*2)),
		pool.WithDialOptions(ggrpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		b.Fatal(err)
	}
	if _, err := p.Connect(ctx); err != nil {
		b.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, ok := p.Get(ctx)
			if ok {
				client := discoverer.NewDiscovererClient(conn)
				_, err := client.Nodes(ctx, &payload.Discoverer_Request{})
				if err != nil {
					b.Errorf("RPC failed: %v", err)
				}
			} else {
				b.Error("failed to get connection from pool")
			}
		}
	})
}

func BenchmarkPool_HighContention(b *testing.B) {
	host := "127.0.0.1"
	port := uint16(9095)
	l, addr := startTestServer(b, host, port)
	defer l.Shutdown(context.Background())
	time.Sleep(100 * time.Millisecond)

	ctx := context.Background()
	// Small pool size to force contention
	p, err := pool.New(ctx,
		pool.WithAddr(addr),
		pool.WithSize(2),
		pool.WithDialOptions(ggrpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		b.Fatal(err)
	}
	if _, err := p.Connect(ctx); err != nil {
		b.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, ok := p.Get(ctx)
			if ok {
				client := discoverer.NewDiscovererClient(conn)
				_, err := client.Nodes(ctx, &payload.Discoverer_Request{})
				if err != nil {
					b.Errorf("RPC failed: %v", err)
				}
			} else {
				b.Error("failed to get connection from pool")
			}
		}
	})
}

func BenchmarkPool_HighContention(b *testing.B) {
	defer ListenAndServe(b, "localhost:5002")()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create pool with small size to force contention
	poolSize := uint64(2)
	p, err := New(ctx,
		WithAddr("localhost:5002"),
		WithSize(poolSize),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		b.Fatal(err)
	}
	p, err = p.Connect(ctx)
	if err != nil {
		b.Fatal(err)
	}

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate high frequency requests
			_ = p.Do(ctx, func(c *ClientConn) error {
				// Simulate short work
				return nil
			})
		}
	})
	b.StopTimer()
}

// NOT IMPLEMENTED BELOW
