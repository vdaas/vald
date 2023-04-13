package mirror

import (
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

func WithClient(gc grpc.Client) Option {
	return func(c *client) error {
		if gc != nil {
			c.c = gc
		}
		return nil
	}
}
