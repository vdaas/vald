package log

import (
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/glg"
	"github.com/vdaas/vald/internal/log/level"
	loggertype "github.com/vdaas/vald/internal/log/logger_type"
	"github.com/vdaas/vald/internal/log/retry"
)

type Option func(*option)

var (
	defaultOptions = []Option{
		WithLogger(
			glg.New(
				glg.WithRetry(
					retry.New(
						retry.WithError(Error),
						retry.WithWarn(Warn),
					),
				),
			),
		),
	}
)

type option struct {
	loggerType loggertype.LoggerType
	level      level.Level
	format     format.Format
	logger     Logger
}

func WithLogger(logger Logger) Option {
	return func(o *option) {
		if logger == nil {
			return
		}
		o.logger = logger
	}
}

func WithLoggerType(str string) Option {
	return func(o *option) {
		if str == "" {
			return
		}
		o.loggerType = loggertype.Atot(str)
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
