package log

import "github.com/vdaas/vald/internal/log/glg"

type Option func(*option)

var (
	defaultOptions = []Option{
		WithLogger(glg.Default()),
	}
)

type option struct {
	logger Logger
}

func (o *option) apply(opts ...Option) *option {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithLogger(logger Logger) Option {
	return func(o *option) {
		o.logger = logger
	}
}
