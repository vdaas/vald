// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package usecase

import (
	"context"
	"os"
	"syscall"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"

	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/recover"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/index/job/exportation/config"
	"github.com/vdaas/vald/pkg/index/job/exportation/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	observability observability.Observability
	server        starter.Server
	exporter      service.Exporter
}

// New returns Runner instance.
func New(cfg *config.Data) (_ runner.Runner, err error) {
	eg := errgroup.Get()

	gOpts, err := cfg.Exporter.Gateway.Opts()
	if err != nil {
		return nil, err
	}
	// skipcq: CRT-D0001
	gOpts = append(gOpts, grpc.WithErrGroup(eg))

	gateway, err := vald.New(vald.WithClient(grpc.New(gOpts...)))
	if err != nil {
		return nil, err
	}

	exporter, err := service.New(
		service.WithStreamListConcurrency(cfg.Exporter.Concurrency),
		service.WithIndexPath(cfg.Exporter.IndexPath),
		service.WithGateway(gateway),
	)
	if err != nil {
		return nil, err
	}

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(
			cfg.Observability,
		)
		if err != nil {
			return nil, err
		}
	}

	// For health check and metrics
	srv, err := starter.New(starter.WithConfig(cfg.Server))
	if err != nil {
		return nil, err
	}

	return &run{
		eg:            eg,
		cfg:           cfg,
		observability: obs,
		server:        srv,
		exporter:      exporter,
	}, nil
}

// PreStart is a method called before execution of Start, and it invokes the PreStart method of observability.
func (r *run) PreStart(ctx context.Context) error {
	if r.observability != nil {
		return r.observability.PreStart(ctx)
	}
	return nil
}

// Start is a method used to initiate an operation in the run, and it returns a channel for receiving errors
// during the operation and an error representing any initialization errors.
func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 3)
	var sech, oech <-chan error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
	}
	sech = r.server.ListenAndServe(ctx)
	cech, err := r.exporter.StartClient(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer func() {
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				// using Fatal to avoid this process to be zombie
				// skipcq: RVV-A0003
				log.Fatalf("failed to find my pid to kill %v", err)
				return
			}
			log.Info("sending SIGTERM to myself to stop this job")
			if err := p.Signal(syscall.SIGTERM); err != nil {
				log.Error(err)
			}
		}()
		return r.exporter.Start(ctx)
	}))

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-oech:
			case err = <-sech:
			case err = <-cech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return errors.Join(ctx.Err(), err)
				case ech <- err:
				}
			}
		}
	}))
	return ech, nil
}

// PreStop is a method called before execution of Stop.
func (r *run) PreStop(ctx context.Context) error {
	return r.exporter.PreStop(ctx)
}

// Stop is a method used to stop an operation in the run.
func (r *run) Stop(ctx context.Context) (errs error) {
	if r.observability != nil {
		if err := r.observability.Stop(ctx); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if r.server != nil {
		if err := r.server.Shutdown(ctx); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

// PostStop is a method called after execution of Stop.
func (*run) PostStop(_ context.Context) error {
	return nil
}
