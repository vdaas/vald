package glg

import (
	"strings"
)

// Option represetns option for GlgLogger.
type Option func(*GlgLogger)

var (
	defaultOpts = []Option{
		WithLevel(DEBUG.String()),
	}
)

func WithEnableJSON() Option {
	return func(g *GlgLogger) {}
}

func WithLevel(lv string) Option {
	return func(g *GlgLogger) {
		g.lv = toLevel(strings.ToUpper(lv))
	}
}
