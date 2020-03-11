package gateway

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/starter"
	"github.com/vdaas/vald/pkg/gateway/vald/config"
	"github.com/vdaas/vald/pkg/gateway/vald/usecase"
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

	daemon, err := usecase.New(s.cfg)
	if err != nil {
		tb.Fatal(err)
	}

	go func() {

	}()
	_ = daemon

	return func() {}
}
