package usecase

import (
	"context"

	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/agent/ngt/config"
	"github.com/vdaas/vald/pkg/agent/ngt/handler/grpc"
	"github.com/vdaas/vald/pkg/agent/ngt/handler/rest"
	"github.com/vdaas/vald/pkg/agent/ngt/router"
	"github.com/vdaas/vald/pkg/agent/ngt/service"
)

type Runner runner.Runner

type run struct {
	cfg    *config.Data
	server service.Server
}

func New(cfg *config.Data) (Runner, error) {
	ngt, err := service.NewNGT(cfg.NGT)
	if err != nil {
		return nil, err
	}
	g := grpc.New(grpc.WithNGT(ngt))

	srv, err := service.NewServer(
		service.WithConfig(cfg.Server),
		service.WithREST(
			router.New(
				router.WithHandler(
					rest.New(
						rest.WithAgent(g),
					),
				),
			),
		),
		service.WithGRPC(g),
		// TODO add GraphQL handler
	)

	if err != nil {
		return nil, err
	}

	return &run{
		cfg:    cfg,
		server: srv,
	}, nil
}

func (r *run) PreStart() error {
	return nil
}

func (r *run) Start(ctx context.Context) <-chan error {
	return r.server.ListenAndServe(ctx)
}

func (r *run) PreStop() error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}
