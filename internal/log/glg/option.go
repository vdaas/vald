package glg

import (
	"strings"

	"github.com/vdaas/vald/internal/log"
)

// Option represetns option for glglogger.
type Option func(*glglogger)

var (
	defaultOpts = []Option{
		WithLevel(log.DEBUG.String()),
	}
)

func WithEnableJSON() Option {
	return func(g *glglogger) {

	}
}

func WithLevel(lv string) Option {
	return func(g *glglogger) {
		g.lv = log.ToLevel(strings.ToUpper(lv))
	}
}

func WithMode(lv string) Option {
	return func(g *glglogger) {

	}
}
