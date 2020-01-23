package log

import (
	glglogger "github.com/kpango/glg"
	"github.com/vdaas/vald/internal/log/glg"
)

type Option func(*option)

var (
	defaultOptions = []Option{
		WithLogger(
			glg.New(
				glglogger.Get(),
				glg.WithRetryOut(retryOut),
				glg.WithRetryOutf(retryOutf),
			),
		),
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
		if logger == nil {
			return
		}
		o.logger = logger
	}
}
