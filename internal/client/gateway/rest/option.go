package rest

type Option func(*gatewayClient)

var (
	defaultOptions = []Option{
		WithAddr("0.0.0.0:8081"),
	}
)

func WithAddr(addr string) Option {
	return func(c *gatewayClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}
