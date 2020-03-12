// Package grpc provides grpc client functions
package grpc

import "github.com/vdaas/vald/internal/config"

type Option func(*gatewayClient)

var (
	defaultOptions = []Option{
		WithAddr("0.0.0.0:8081"),
		WithGRPCClientConfig(&config.GRPCClient{
			Addrs: []string{
				"0.0.0.0:8081",
			},
		}),
	}
)

func WithAddr(addr string) Option {
	return func(c *gatewayClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}

func WithGRPCClientConfig(cfg *config.GRPCClient) Option {
	return func(c *gatewayClient) {
		if cfg != nil {
			c.cfg = cfg.Bind()
		}
	}
}
