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

	"github.com/vdaas/vald/apis/grpc/manager/compressor"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/metric"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/manager/compressor/config"
	handler "github.com/vdaas/vald/pkg/manager/compressor/handler/grpc"
	"github.com/vdaas/vald/pkg/manager/compressor/handler/rest"
	"github.com/vdaas/vald/pkg/manager/compressor/router"
	"github.com/vdaas/vald/pkg/manager/compressor/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	backup        service.Backup
	compressor    service.Compressor
	server        starter.Server
	observability observability.Observability
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()

	var (
		b service.Backup
	)

	if addrs := cfg.BackupManager.Client.Addrs; len(addrs) == 0 {
		return nil, errors.ErrInvalidBackupConfig
	}

	b, err = service.NewBackup(
		service.WithBackupAddr(cfg.BackupManager.Client.Addrs[0]),
		service.WithBackupClient(
			grpc.New(
				append(cfg.BackupManager.Client.Opts(),
					grpc.WithErrGroup(eg),
					grpc.WithDialOptions(
						grpc.WithStatsHandler(metric.NewClientHandler()),
					),
				)...,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	c, err := service.NewCompressor(
		service.WithCompressAlgorithm(cfg.Compressor.CompressAlgorithm),
		service.WithCompressionLevel(cfg.Compressor.CompressionLevel),
		service.WithLimitation(cfg.Compressor.ConcurrentLimit),
		service.WithBuffer(cfg.Compressor.Buffer),
		service.WithErrGroup(eg),
	)
	if err != nil {
		return nil, err
	}
	g := handler.New(
		handler.WithCompressor(c),
		handler.WithBackup(b),
	)

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			compressor.RegisterBackupServer(srv, g)
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

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(cfg.Observability)
		if err != nil {
			return nil, err
		}
		grpcServerOptions = append(
			grpcServerOptions,
			server.WithGRPCOption(
				grpc.StatsHandler(metric.NewServerHandler()),
			),
		)
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
			return grpcServerOptions
		}),
		// TODO add GraphQL handler
	)

	if err != nil {
		return nil, err
	}

	return &run{
		eg:            eg,
		cfg:           cfg,
		backup:        b,
		compressor:    c,
		server:        srv,
		observability: obs,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	log.Info("daemon pre-start")
	err := r.compressor.PreStart(ctx)
	if err != nil {
		return err
	}
	return r.observability.PreStart(ctx)
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 4)
	var bech, cech, sech, oech <-chan error
	var err error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
	}
	if r.backup != nil {
		bech, err = r.backup.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.compressor != nil {
		cech = r.compressor.Start(ctx)
	}
	sech = r.server.ListenAndServe(ctx)
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		log.Info("daemon start")
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-oech:
			case err = <-bech:
			case err = <-cech:
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

func (r *run) PreStop(ctx context.Context) error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	if r.observability != nil {
		r.observability.Stop(ctx)
	}
	return r.server.Shutdown(ctx)
}

func (r *run) PostStop(ctx context.Context) error {
	return nil
}
