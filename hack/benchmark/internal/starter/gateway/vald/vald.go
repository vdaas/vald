package vald

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/starter"
	"github.com/vdaas/vald/pkg/gateway/vald/config"
)

type server struct {
	cfg *config.Data
}

func New(opts ...Option) starter.Starter {
	srv := new(server)
	for _, opt := range append(defaultOptions, opts...) {
		opt(srv)
	}
	return srv
}

func (s *server) Run(ctx context.Context, tb testing.TB) func() {
	tb.Helper()

	// TODO (@hlts2): Make when divided gateway.

	return func() {}
}
