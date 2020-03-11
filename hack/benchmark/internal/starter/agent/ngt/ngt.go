package ngt

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/hack/benchmark/internal/starter"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/agent/ngt/config"
	"github.com/vdaas/vald/pkg/agent/ngt/usecase"
)

const name = "agent-ngt"

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

	log.Init()

	ctx, cancel := context.WithCancel(ctx)

	daemon, err := usecase.New(s.cfg)
	if err != nil {
		tb.Fatal(err)
	}

	go func() {
		err := runner.Run(ctx, daemon, name)
		if err != nil {
			tb.Fatalf("agent runner returned error %s", err.Error())
		}
	}()

	time.Sleep(5 * time.Second)

	return func() {
		cancel()
	}
}
