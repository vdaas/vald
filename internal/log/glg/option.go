package glg

import (
	"strings"

	"github.com/vdaas/vald/internal/log"
)

// Option represetns option for glglogger.
type Option func(*glglogger)

var (
	defaultGlgOpts = []Option{
		WithLogLevel(log.DEBUG.String()),
	}
)

func WithLogLevel(lv string) Option {
	return func(g *glglogger) {
		g.lv = log.ToLevel(strings.ToUpper(lv))
	}
}
