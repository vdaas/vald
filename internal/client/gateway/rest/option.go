package rest

type Option func(*gatewayClient)

var (
	defaultOptions = []Option{
		WithAddr("http://127.0.0.1:8080"),
	}
)

func WithAddr(addr string) Option {
	return func(c *gatewayClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}
