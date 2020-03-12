// Package grpc provides grpc client functions
package grpc

import (
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/net/grpc"
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

func WithGRPCClientOption(opts ...grpc.Option) Option {
	return func(c *ngtdClient) {
		if len(opts) != 0 {
			c.opts = opts
		}
	}
}
