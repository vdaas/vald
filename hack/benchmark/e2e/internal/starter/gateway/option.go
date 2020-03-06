package gateway

type Option func(*server)

var (
	defaultOptions = []Option{}
)
