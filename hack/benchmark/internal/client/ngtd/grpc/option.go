package grpc

import (
	"github.com/vdaas/vald/internal/config"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
)

type Option func(*ngtdClient)

var (
	defaultOptions = []Option{
		WithAddr("127.0.0.1:8200"),
		WithGRPCClientOption(
			(&config.GRPCClient{
				Addrs: []string{
					"127.0.0.1:8200",
				},
				CallOption: &config.CallOption{
					MaxRecvMsgSize: 100000000000,
				},
				DialOption: &config.DialOption{
					Insecure: true,
				},
			}).Bind().Opts()...,
		),
	}
)

func WithAddr(addr string) Option {
	return func(c *ngtdClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}

func WithStreamConcurrency(n int) Option {
	return func(c *ngtdClient) {
		if n > 0 {
			c.streamConcurrency = n
		}
	}
}

func WithGRPCClientOption(opts ...igrpc.Option) Option {
	return func(c *ngtdClient) {
		if len(opts) != 0 {
			c.opts = opts
		}
	}
}
