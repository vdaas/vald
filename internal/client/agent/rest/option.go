package rest

type Option func(*agentClient)

var (
	defaultOptions = []Option{
		WithAddr("http://127.0.0.1:8081"),
	}
)

func WithAddr(addr string) Option {
	return func(ac *agentClient) {
		if len(addr) != 0 {
			ac.addr = addr
		}
	}
}
