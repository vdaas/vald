//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/vdaas/vald/apis/grpc/manager/backup"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/metric"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/config"
	handler "github.com/vdaas/vald/pkg/manager/backup/mysql/handler/grpc"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/handler/rest"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/router"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	mySQL         service.MySQL
	server        starter.Server
	observability observability.Observability
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	m, err := service.New(cfg.MySQL)
	if err != nil {
		return nil, err
	}
	g := handler.New(handler.WithMySQL(m))
	eg := errgroup.Get()

	obs, err := observability.New(cfg.Observability)
	if err != nil {
		return nil, err
	}

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
								rest.WithBackup(g),
							),
						),
					)),
			}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(srv *grpc.Server) {
					backup.RegisterBackupServer(srv, g)
				}),
				server.WithGRPCOption(
					grpc.StatsHandler(metric.NewServerHandler()),
				),
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
		eg:            eg,
		cfg:           cfg,
		mySQL:         m,
		server:        srv,
		observability: obs,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	log.Info("daemon pre-start")
	err := r.mySQL.Connect(ctx)
	if err != nil {
		return err
	}
	return r.observability.PreStart(ctx)
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 2)
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		sech := r.server.ListenAndServe(ctx)
		oech := r.observability.Start(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-sech:
			case err = <-oech:
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

func (r *run) PreStop(ctx context.Context) error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	r.observability.Stop(ctx)
	return r.server.Shutdown(ctx)
}

func (r *run) PostStop(ctx context.Context) error {
	return r.mySQL.Close(ctx)
}
