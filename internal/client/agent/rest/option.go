package rest

type Option func(*agentClient)

var (
	defaultOptions = []Option{
		WithAddr("0.0.0.0:8081"),
	}
)

func WithAddr(addr string) Option {
	return func(ac *agentClient) {
		if len(addr) != 0 {
			ac.addr = addr
		}
	}
}
