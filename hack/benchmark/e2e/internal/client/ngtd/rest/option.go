package rest

type Option func(*ngtdClient)

var (
	defaultOptions = []Option{
		WithAddr("http://127.0.0.1:8200"),
	}
)

func WithAddr(addr string) Option {
	return func(c *ngtdClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}
