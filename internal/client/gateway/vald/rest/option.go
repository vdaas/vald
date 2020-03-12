// Package rest provides rest client functions
package rest

// Option is gatewayClient configure
type Option func(*gatewayClient)

var (
	defaultOptions = []Option{
		WithAddr("http://127.0.0.1:8080"),
	}
)

// WithAddr returns Option that sets addr
func WithAddr(addr string) Option {
	return func(c *gatewayClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}
