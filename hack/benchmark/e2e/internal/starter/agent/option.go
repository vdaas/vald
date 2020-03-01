package agent

type Option func(*server)

var (
	defaultOptions = []Option{}
)
