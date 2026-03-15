package grpc_test

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/servers/server"
	googlegrpc "google.golang.org/grpc"
)

type serverMock struct {
	discoverer.DiscovererServer
}

func (*serverMock) Pods(context.Context, *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	return &payload.Info_Pods{
		Pods: []*payload.Info_Pod{
			{
				Name: "vald is high scalable distributed high-speed approximate nearest neighbor search engine",
			},
		},
	}, nil
}

func (*serverMock) Nodes(context.Context, *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return new(payload.Info_Nodes), nil
}

func init() {
	grpc.StartMockServer = func() func() {
		srv, err := server.New(
			server.WithHost("127.0.0.1"),
			server.WithPort(5001),
			server.WithServerMode(server.GRPC),
			server.WithGRPCRegisterar(func(srv *googlegrpc.Server) {
				discoverer.RegisterDiscovererServer(srv, new(serverMock))
			}),
		)
		if err != nil {
			panic(err)
		}

		errCh := make(chan error, 1)
		go srv.ListenAndServe(context.Background(), errCh)
		return func() {
			srv.Shutdown(context.Background())
		}
	}
}
