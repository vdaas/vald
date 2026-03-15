package grpc

import (
	"context"
	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var StartMockServer func() func()

const (
	DefaultServerAddr = "localhost:5001"
	DefaultPoolSize   = 4
)

func BenchmarkClient_ExecuteRPC(b *testing.B) {
	defer StartMockServer()()

	ctx := context.Background()
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
