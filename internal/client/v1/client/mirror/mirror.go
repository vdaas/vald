package mirror

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	apiName = "vald/internal/client/v1/client/mirror"
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
	if c.c == nil {
		if c.Client != nil {
			c.c = c.Client.GRPCClient()
		} else {
			if len(c.addrs) == 0 {
				return nil, errors.ErrGRPCTargetAddrNotFound
			}
			c.c = grpc.New(grpc.WithAddrs(c.addrs...))
		}
	}
	if c.Client == nil {
		if len(c.addrs) == 0 {
			return nil, errors.ErrGRPCTargetAddrNotFound
		}
		var err error
		c.Client, err = vald.New(
			vald.WithAddrs(c.addrs...),
			vald.WithClient(c.c),
		)
		if err != nil {
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

func (c *client) Register(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (res *payload.Mirror_Targets, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Register")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = c.c.RoundRobin(ctx, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		res, err = mirror.NewMirrorClient(conn).Register(ctx, in, copts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Advertise(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (res *payload.Mirror_Targets, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Advertise")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = c.c.RoundRobin(ctx, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		res, err = mirror.NewMirrorClient(conn).Advertise(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
