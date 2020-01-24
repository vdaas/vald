package glg

import (
	"strings"

	"github.com/vdaas/vald/internal/log/retry"
)

// Option represetns option for Logger.
type Option func(*Logger)

var (
	defaultOpts = []Option{
		WithLevel(DEBUG.String()),
		WithRetry(retry.NewNop()),
	}
)

func WithEnableJSON() Option {
	return func(g *Logger) {}
}

func WithLevel(lv string) Option {
	return func(g *Logger) {
		g.lv = toLevel(strings.ToUpper(lv))
	}
}

func WithRetry(rt retry.Retry) Option {
	return func(g *Logger) {
		if rt == nil {
			return
		}
		g.rt = rt
	}
}
