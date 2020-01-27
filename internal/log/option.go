package log

import (
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/glg"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/log/mode"
	retry "github.com/vdaas/vald/internal/log/retry"
)

type Option func(*option)

var (
	defaultOptions = []Option{
		WithLogger(
			glg.New(
				glg.WithRetry(
					retry.New(
						Error,
						Warn,
					),
				),
			),
		),
	}
)

type option struct {
	mode   mode.Mode
	level  level.Level
	format format.Format
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

func WithMode(str string) Option {
	return func(o *option) {
		if str == "" {
			return
		}
		o.mode = mode.Atom(str)
	}
}

func WithLevel(str string) Option {
	return func(o *option) {
		if str == "" {
			return
		}
		o.level = level.Atol(str)
	}
}

func WithFormat(str string) Option {
	return func(o *option) {
		if str == "" {
			return
		}
		o.format = format.Atof(str)
	}
}
