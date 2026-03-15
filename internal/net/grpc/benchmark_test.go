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
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type broadcastServer struct {
	discoverer.UnimplementedDiscovererServer
}

func (s *broadcastServer) Nodes(
	ctx context.Context, req *payload.Discoverer_Request,
) (*payload.Info_Nodes, error) {
	return &payload.Info_Nodes{}, nil
}

func startBroadcastServer(t testing.TB, eg errgroup.Group) (net.Listener, func()) {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	srv := ggrpc.NewServer()
	discoverer.RegisterDiscovererServer(srv, &broadcastServer{})

	eg.Go(func() error {
		return srv.Serve(lis)
	})

	return lis, func() {
		srv.Stop()
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
	ctx := context.Background()
	client := grpc.New(
		"benchmark-client",
		grpc.WithAddrs(addrs...),
		grpc.WithDialOptions(
			ggrpc.WithTransportCredentials(insecure.NewCredentials()),
			ggrpc.WithContextDialer(dialer),
		),
		grpc.WithConnectionPoolSize(1),
	)
	defer client.Close(ctx)

	// Connect to all agents
	var connectedCount int64
	for _, addr := range addrs {
		eg.Go(func() error {
			// Retry connect loop since server might take a moment to start
			for range 10 {
				_, err := client.Connect(ctx, addr)
				if err == nil {
					atomic.AddInt64(&connectedCount, 1)
					return nil
				}
				time.Sleep(10 * time.Millisecond)
			}
			return nil
		})
	}
	eg.Wait()

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		err := client.RangeConcurrent(ctx, concurrency, func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			c := discoverer.NewDiscovererClient(conn)
			_, err := c.Nodes(ctx, &payload.Discoverer_Request{})
			return err
		})
		if err != nil {
			b.Fatalf("Broadcast failed: %v", err)
		}
	}
	b.StopTimer()
	// NOT IMPLEMENTED BELOW
}
