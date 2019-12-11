//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package usecase

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/agent"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/agent/ngt/config"
	handler "github.com/vdaas/vald/pkg/agent/ngt/handler/grpc"
	"github.com/vdaas/vald/pkg/agent/ngt/handler/rest"
	"github.com/vdaas/vald/pkg/agent/ngt/router"
	"github.com/vdaas/vald/pkg/agent/ngt/service"
)

type run struct {
	eg     errgroup.Group
	cfg    *config.Data
	ngt    service.NGT
	server starter.Server
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	ngt, err := service.New(cfg.NGT)
	if err != nil {
		return nil, err
	}
	g := handler.New(handler.WithNGT(ngt))
	eg := errgroup.Get()

	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithREST(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithHTTPHandler(
					router.New(
						router.WithTimeout(sc.HTTP.HandlerTimeout),
						router.WithErrGroup(eg),
						router.WithHandler(
							rest.New(
								rest.WithAgent(g),
							),
						),
					),
				),
			}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(srv *grpc.Server) {
					agent.RegisterAgentServer(srv, g)
				}),
				server.WithPreStartFunc(func() error {
					// TODO check unbackupped upstream
					return nil
				}),
				server.WithPreStopFunction(func() error {
					// TODO backup all index data here
					return nil
				}),
			}
		}),
		// TODO add GraphQL handler
	)

	if err != nil {
		return nil, err
	}

	return &run{
		eg:     eg,
		ngt:    ngt,
		cfg:    cfg,
		server: srv,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	return nil
}

func (r *run) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 2)
	nech := r.ngt.Start(ctx)
	sech := r.server.ListenAndServe(ctx)
	r.eg.Go(safety.RecoverFunc(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-nech:
				if err != nil {
					ech <- err
				}
			case err := <-sech:
				if err != nil {
					ech <- err
				}
				return nil
			}
		}
	}))
	return ech
}

func (r *run) PreStop(ctx context.Context) error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}

func (r *run) PostStop(ctx context.Context) error {
	return nil
}
