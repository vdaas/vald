package log

import (
	glglogger "github.com/kpango/glg"
	"github.com/vdaas/vald/internal/log/glg"
	"github.com/vdaas/vald/internal/log/retry"
)

type Option func(*option)

var (
	defaultOptions = []Option{
		WithLogger(
			glg.New(
				glglogger.Get(),
				glg.WithRetry(retry.New(
					Warn,
					Error,
				)),
			),
		),
	}
)

type option struct {
	logger Logger
}

func WithLogger(logger Logger) Option {
	return func(o *option) {
		o.logger = logger
	}
}
