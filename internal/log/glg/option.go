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
		WithConfig(&config.Log{
			Level: log.DEBUG.String(),
		}),
	}
)

func WithConfig(cfg *config.Log) Option {
	return func(g *glglogger) {
		if cfg == nil {
			return
		}

		if cfg.Format == "json" {
			// TODO:
		}

		g.lv = log.ToLevel(strings.ToUpper(cfg.Level))
	}
}
