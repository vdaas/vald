// Package strategy provides strategy for e2e testing functions
package strategy

import "github.com/vdaas/vald/internal/client"

type SearchOption func(*search)

var (
	defaultSearchOptions = []SearchOption{
		WithSearchParallel(false),
	}
)

func WithSearchParallel(flag bool) SearchOption {
	return func(s *search) {
		s.parallel = flag
	}
}

func WithSearchConfig(cfg *client.SearchConfig) SearchOption {
	return func(s *search) {
		if cfg != nil {
			s.cfg = cfg
		}
	}
}
