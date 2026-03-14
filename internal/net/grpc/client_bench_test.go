package grpc_test

import (
	"context"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	grpc_client "github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/servers/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DefaultServerAddr = "localhost:5001"
	DefaultPoolSize   = 4
)

type testServer struct {
	discoverer.DiscovererServer
}

func (*testServer) Pods(context.Context, *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	return &payload.Info_Pods{
		Pods: []*payload.Info_Pod{
			{
				Name: "vald is high scalable distributed high-speed approximate nearest neighbor search engine",
			},
		},
	}, nil
}

func (*testServer) Nodes(context.Context, *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return new(payload.Info_Nodes), nil
}

func listenAndServe(b *testing.B, addr string) func() {
	b.Helper()
	srv, err := server.New(
		server.WithServerMode(server.GRPC),
		server.WithHost("127.0.0.1"),
		server.WithPort(5001),
		server.WithGRPCRegisterar(func(srv *grpc.Server) {
			discoverer.RegisterDiscovererServer(srv, new(testServer))
		}),
	)
	if err != nil {
		b.Error(err)
	}

	errCh := make(chan error, 1)
	go srv.ListenAndServe(context.Background(), errCh)
	return func() {
		srv.Shutdown(context.Background())
	}
}

func BenchmarkClient_ExecuteRPC(b *testing.B) {
	defer listenAndServe(b, DefaultServerAddr)()

	ctx := b.Context()
	client := grpc_client.New("test",
		grpc_client.WithAddrs(DefaultServerAddr),
		grpc_client.WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
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
		grpc_client.ExecuteRPC(client, ctx, conn, DefaultServerAddr, func(ctx context.Context, conn *grpc_client.ClientConn, copts ...grpc_client.CallOption) (any, error) {
			return discoverer.NewDiscovererClient(conn).Nodes(ctx, new(payload.Discoverer_Request))
		})
	}
}
