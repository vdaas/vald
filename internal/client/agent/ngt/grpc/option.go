// Package grpc provides gRPC client functions
package grpc

import (
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/net/grpc"
)

// Option is agentClient configure.
type Option func(*agentClient)

var (
	defaultOptions = []Option{
		WithAddr("127.0.0.1:8082"),
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
			}).Bind().Opts()...),
	}
)

// WithAddr returns Option that sets addr.
func WithAddr(addr string) Option {
	return func(c *agentClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}

// WithGRPCClientOption returns Option that sets options for gRPC.
func WithGRPCClientOption(opts ...grpc.Option) Option {
	return func(c *agentClient) {
		if len(opts) != 0 {
			c.opts = append(c.opts, opts...)
		}
	}
}
