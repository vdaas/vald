// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

	iconfig "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/recover"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/operator/vald/config"
	"github.com/vdaas/vald/pkg/operator/vald/service"
)

// runnableErrChanBuf is the buffer size for the error channel returned by Start.
// It equals the number of runnable error sources: observability, operator, server.
const runnableErrChanBuf = 3

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	observability observability.Observability
	server        starter.Server
	operator      service.Operator
}

func New(cfg *config.Data) (_ runner.Interface, err error) {
	eg := errgroup.Get()
	operator, err := service.New(cfg.Operator)
	if err != nil {
		return nil, err
	}

	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithGRPC(func(cfg *iconfig.Server) []server.Option {
			return []server.Option{
				server.WithGRPCOption(
					grpc.ChainUnaryInterceptor(recover.Interceptor()),
					grpc.ChainStreamInterceptor(recover.StreamInterceptor()),
				),
			}
		}),
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

	return &run{
		eg:            eg,
		cfg:           cfg,
		observability: obs,
		server:        srv,
		operator:      operator,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	if r.observability != nil {
		return r.observability.PreStart(ctx)
	}
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, runnableErrChanBuf)
	var oech, dech, sech <-chan error
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		if r.observability != nil {
			oech = r.observability.Start(ctx)
		}
		dech, err = r.operator.Start(ctx)
		if err != nil {
			ech <- err
			return err
		}

		sech = r.server.ListenAndServe(ctx)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-oech:
			case err = <-dech:
			case err = <-sech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ech <- err:
				}
			}
		}
	}))
	return ech, nil
}

func (*run) PreStop(_ context.Context) error {
	return nil
}

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

func (*run) PostStop(_ context.Context) error {
	return nil
}
