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
	mode   string
	lv     string
	format string
	logger Logger
}

func WithLogger(logger Logger) Option {
	return func(o *option) {
		if logger == nil {
			return
		}
		o.logger = logger
	}
}

func WithMode(mode string) Option {
	return func(o *option) {
		if mode == "" {
			return
		}
		o.mode = mode
	}
}

func WithLevel(lv string) Option {
	return func(o *option) {
		if lv == "" {
			return
		}
		o.lv = lv
	}
}

func WithFormat(format string) Option {
	return func(o *option) {
		if format == "" {
			return
		}
		o.format = format
	}
}
