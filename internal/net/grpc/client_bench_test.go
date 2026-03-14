package grpc

import (
	"context"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DefaultServerAddr = "localhost:5001"
	DefaultPoolSize   = 4
)

type server struct {
	discoverer.DiscovererServer
}

func (*server) Pods(context.Context, *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	return &payload.Info_Pods{
		Pods: []*payload.Info_Pod{
			{
				Name: "vald is high scalable distributed high-speed approximate nearest neighbor search engine",
			},
		},
	}, nil
}

func (*server) Nodes(context.Context, *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return new(payload.Info_Nodes), nil
}

func ListenAndServe(b *testing.B, addr string) func() {
	b.Helper()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		b.Error(err)
	}

	// skipcq: GO-S0902
	s := grpc.NewServer()
	discoverer.RegisterDiscovererServer(s, new(server))

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		wg.Done()
		if err := s.Serve(lis); err != nil {
			b.Error(err)
		}
	}()

	wg.Wait()
	return func() {
		s.Stop()
	}
}

func BenchmarkClient_ExecuteRPC(b *testing.B) {
	defer ListenAndServe(b, DefaultServerAddr)()

	ctx := b.Context()
	client := New("test",
		WithAddrs(DefaultServerAddr),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	defer client.Close(ctx)

	conn, err := client.Connect(ctx, DefaultServerAddr)
	if err != nil {
		b.Fatal(err)
	}

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		client.(*gRPCClient).executeRPC(ctx, conn, DefaultServerAddr, func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error) {
			return discoverer.NewDiscovererClient(conn).Nodes(ctx, new(payload.Discoverer_Request))
		})
	}
}
