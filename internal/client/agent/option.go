package agent

import "github.com/vdaas/vald/internal/config"

type Option func(*agentClient)

var (
	defaultOptions = []Option{
		WithAddr("0.0.0.0"),
		WithStreamConcurrency(5),
	}
)

func WithAddr(addr string) Option {
	return func(c *agentClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}

func WithStreamConcurrency(n int) Option {
	return func(c *agentClient) {
		if n > 0 {
			c.streamConcurrency = n
		}
	}
}

func WithGRPCClientConfig(cfg *config.GRPCClient) Option {
	return func(c *agentClient) {
		if cfg != nil {
			c.cfg = cfg
		}
	}
}
