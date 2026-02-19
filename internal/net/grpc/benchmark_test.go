package grpc_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/sync/errgroup"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type broadcastServer struct {
	discoverer.UnimplementedDiscovererServer
}

func (s *broadcastServer) Nodes(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return &payload.Info_Nodes{}, nil
}

func startBroadcastServer(t testing.TB, eg errgroup.Group) (string, func()) {
	// Find a free port
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()

	srv, err := server.New(
		server.WithServerMode(server.GRPC),
		server.WithName(fmt.Sprintf("benchmark-server-%d", port)),
		server.WithHost("localhost"),
		server.WithPort(uint16(port)),
		server.WithErrorGroup(eg),
		server.WithGRPCRegisterar(func(s *ggrpc.Server) {
			discoverer.RegisterDiscovererServer(s, &broadcastServer{})
		}),
	)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	ech := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(context.Background(), ech); err != nil {
			// t.Logf("server stopped: %v", err)
		}
	}()

	addr := fmt.Sprintf("localhost:%d", port)
	return addr, func() {
		srv.Shutdown(context.Background())
	}
}

func BenchmarkBroadcast(b *testing.B) {
	// Number of agents to simulate
	const numAgents = 2
	const concurrency = 2

	eg := errgroup.Get()

	// Start servers
	addrs := make([]string, numAgents)
	closers := make([]func(), numAgents)
	for i := 0; i < numAgents; i++ {
		addr, closeFunc := startBroadcastServer(b, eg)
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
	client := grpc.New(
		"benchmark-client",
		grpc.WithAddrs(addrs...),
		grpc.WithDialOptions(ggrpc.WithTransportCredentials(insecure.NewCredentials())),
		grpc.WithConnectionPoolSize(1),
	)
	defer client.Close(ctx)

	// Connect to all agents
	var connectedCount int64
	var wg sync.WaitGroup
	wg.Add(numAgents)
	for _, addr := range addrs {
		go func(a string) {
			defer wg.Done()
			// Retry connect loop since server might take a moment to start
			for i := 0; i < 10; i++ {
				_, err := client.Connect(ctx, a)
				if err == nil {
					atomic.AddInt64(&connectedCount, 1)
					return
				}
			}
		}(addr)
	}
	wg.Wait()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := client.RangeConcurrent(ctx, concurrency, func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
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
