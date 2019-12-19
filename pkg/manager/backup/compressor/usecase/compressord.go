//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/manager/backup/compressor/config"
	handler "github.com/vdaas/vald/pkg/manager/backup/compressor/handler/grpc"
	"github.com/vdaas/vald/pkg/manager/backup/compressor/handler/rest"
	"github.com/vdaas/vald/pkg/manager/backup/compressor/router"
	"github.com/vdaas/vald/pkg/manager/backup/compressor/service"
)

type run struct {
	eg         errgroup.Group
	cfg        *config.Data
	backup     service.Backup
	compressor service.Compressor
	server     starter.Server
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
				)...,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	compressor, err := service.NewCompressor(
		service.WithCompressAlgorithm(cfg.Compressor.CompressAlgorithm),
		service.WithLimitation(cfg.Compressor.ConcurrentLimit),
		service.WithBuffer(cfg.Compressor.Buffer),
		service.WithErrGroup(eg),
	)
	if err != nil {
		return nil, err
	}
	g := handler.New(
		handler.WithCompressor(compressor),
		handler.WithBackup(b),
	)

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
								rest.WithCompress(g),
							),
						),
					)),
			}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(srv *grpc.Server) {
					backup.RegisterCompressServer(srv, g)
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
		eg:         eg,
		cfg:        cfg,
		backup:     b,
		compressor: compressor,
		server:     srv,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	log.Info("daemon pre-start")
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 2)
	var bech, cech, sech <-chan error
	var err error
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
			case err = <-sech:
			case err = <-bech:
			case err = <-cech:
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
	return r.server.Shutdown(ctx)
}

func (r *run) PostStop(ctx context.Context) error {
	return nil
}
