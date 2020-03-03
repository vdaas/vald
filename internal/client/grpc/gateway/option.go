package gateway

import "github.com/vdaas/vald/internal/config"

type Option func(*gatewayClient)

var (
	defaultOptions = []Option{
		WithAddr("0.0.0.0:8081"),
		WithStreamConcurrency(5),
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

func WithStreamConcurrency(n int) Option {
	return func(c *gatewayClient) {
		if n > 0 {
			c.streamConcurrency = n
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
