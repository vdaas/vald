package grpc

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type broadcastServer struct {
	discoverer.UnimplementedDiscovererServer
}

func (s *broadcastServer) Nodes(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return &payload.Info_Nodes{}, nil
}

func startBroadcastServer(t testing.TB) (string, func()) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	discoverer.RegisterDiscovererServer(s, &broadcastServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
		}
	}()
	return lis.Addr().String(), func() {
		s.Stop()
	}
}

func BenchmarkBroadcast(b *testing.B) {
	// Number of agents to simulate
	const numAgents = 2
	const concurrency = 2

	// Start servers
	addrs := make([]string, numAgents)
	closers := make([]func(), numAgents)
	for i := 0; i < numAgents; i++ {
		addr, closeFunc := startBroadcastServer(b)
		addrs[i] = addr
		closers[i] = closeFunc
	}
	defer func() {
		for _, c := range closers {
			c()
		}
	}()

	// Initialize Client
	ctx := context.Background()
	client := New(
		"benchmark-client",
		WithAddrs(addrs...),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
		WithConnectionPoolSize(1),
	)
	defer client.Close(ctx)

	// Connect to all agents
	var connectedCount int64
	var wg sync.WaitGroup
	wg.Add(numAgents)
	for _, addr := range addrs {
		go func(a string) {
			defer wg.Done()
			_, err := client.Connect(ctx, a)
			if err == nil {
				atomic.AddInt64(&connectedCount, 1)
			}
		}(addr)
	}
	wg.Wait()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := client.RangeConcurrent(ctx, concurrency, func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error {
			c := discoverer.NewDiscovererClient(conn)
			_, err := c.Nodes(ctx, &payload.Discoverer_Request{})
			return err
		})
		if err != nil {
			// b.Fatalf("Broadcast failed: %v", err)
		}
	}
	b.StopTimer()
}
