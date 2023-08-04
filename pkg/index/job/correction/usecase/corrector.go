package usecase

import (
	"context"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/runner"
)

type run struct {
	eg            errgroup.Group
	// cfg           *config.Data
	// server        starter.Server
	// observability observability.Observability
	// indexer       service.Indexer
}

// FIXME: add config
func New() (r runner.Runner, err error) {
	eg := errgroup.Get()
	return &run{
		eg: eg,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 5) // FIXME: magic number 5
	return ech, nil
}

func (*run) PreStop(context.Context) error {
	return nil
}

func (*run) Stop(context.Context) error {
	return nil
}


func (r *run) PostStop(ctx context.Context) error {
	return nil
}
