package mirror

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/net/grpc"
)

type Client interface {
	vald.Client
	mirror.MirrorClient
	GRPCClient() grpc.Client
	Start(context.Context) (<-chan error, error)
	Stop(context.Context) error
}

type client struct {
	addrs []string
	c     grpc.Client
	vald.Client
}

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (c *client) Start(ctx context.Context) (<-chan error, error) {
	return c.c.StartConnectionMonitor(ctx)
}

func (c *client) Stop(ctx context.Context) error {
	return c.Stop(ctx)
}

func (c *client) GRPCClient() grpc.Client {
	return c.c
}

// // Register is the RPC to register other mirror servers.
// Register(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error)
// // Advertise is the RPC to advertise other mirror servers.
// Advertise(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error)

func (c *client) Register(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error) {
	return nil, nil
}

func (c *client) Advertise(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error) {
	return nil, nil
}
