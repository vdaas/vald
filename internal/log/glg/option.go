package glg

import (
	"strings"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/log"
)

// Option represetns option for glglogger.
type Option func(*glglogger)

var (
	defaultOpts = []Option{
		WithLevel(log.DEBUG.String()),
	}
)

func WithLevel(lv string) Option {
	return func(g *glglogger) {
		g.lv = log.ToLevel(strings.ToUpper(lv))
	}
}

func WithConfig(cfg *config.Log) Option {
	return func(g *glglogger) {
		if cfg == nil {
			return
		}

		if cfg.Format == "json" {
			WithEnableJSON()(g)
		}

		g.lv = log.ToLevel(strings.ToUpper(cfg.Level))
	}
}

func WithEnableJSON() Option {
	return func(g *glglogger) {
		// TODO: Enable JSON
	}
}
