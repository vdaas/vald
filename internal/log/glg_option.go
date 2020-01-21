package log

import "strings"

// GlgOption represetns option for glglogger.
type GlgOption func(*glglogger)

var (
	defaultGlgOpts = []GlgOption{
		WithLogLevel(DEBUG.String()),
	}
)

func WithLogLevel(lv string) GlgOption {
	return func(g *glglogger) {
		g.lv = level(strings.ToUpper(lv))
	}
}
