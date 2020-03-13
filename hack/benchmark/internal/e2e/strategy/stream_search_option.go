// Package strategy provides strategy for e2e testing functions
package strategy

import "github.com/vdaas/vald/internal/client"

type StreamSearchOption func(*streamSearch)

var (
	defaultStreamSearchOptions = []StreamSearchOption{
		WithStreamSearchConfig(&client.SearchConfig{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.01,
		}),
	}
)

func WithStreamSearchConfig(cfg *client.SearchConfig) StreamSearchOption {
	return func(ss *streamSearch) {
		if cfg != nil {
			ss.cfg = cfg
		}
	}
}
