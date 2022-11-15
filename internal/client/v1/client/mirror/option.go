package mirror

import (
	"context"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/net/grpc"
)

type Option func(c *client) error

var defaultOpts = []Option{}

func WithAddrs(addrs ...string) Option {
	return func(c *client) error {
		if addrs == nil {
			return nil
		}
		if c.addrs != nil {
			c.addrs = append(c.addrs, addrs...)
		} else {
			c.addrs = addrs
		}
		return nil
	}
}

func WithValdClient(vc vald.Client) Option {
	return func(c *client) error {
		if vc != nil {
			c.Client = vc
			if c.c != nil {
				err := c.c.Close(context.Background())
				if err != nil {
					return err
				}
			}
			c.c = c.Client.GRPCClient()
		}
		return nil
	}
}

func WithClient(gc grpc.Client) Option {
	return func(c *client) error {
		if gc != nil {
			c.c = gc
		}
		return nil
	}
}
