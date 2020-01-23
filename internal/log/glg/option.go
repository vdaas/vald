package glg

import (
	"strings"

	"github.com/vdaas/vald/internal/log/retry"
)

// Option represetns option for GlgLogger.
type Option func(*GlgLogger)

var (
	defaultOpts = []Option{
		WithLevel(DEBUG.String()),
		WithRetryOut(func(fn func(vals ...interface{}) error, vals ...interface{}) {
			fn(vals...)
		}),
		WithRetryOutf(func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{}) {
			fn(format, vals...)
		}),
	}
)

func WithEnableJSON() Option {
	return func(g *GlgLogger) {}
}

func WithLevel(lv string) Option {
	return func(g *GlgLogger) {
		if lv == "" {
			return
		}
		g.lv = toLevel(strings.ToUpper(lv))
	}
}

func WithRetryOut(fn retry.Out) Option {
	return func(g *GlgLogger) {
		if fn == nil {
			return
		}
		g.rout = fn
	}
}

func WithRetryOutf(fn retry.Outf) Option {
	return func(g *GlgLogger) {
		if fn == nil {
			return
		}
		g.routf = fn
	}
}
