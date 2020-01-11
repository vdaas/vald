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

	"github.com/vdaas/vald/apis/grpc/vald"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/manager/index/config"
	handler "github.com/vdaas/vald/pkg/manager/index/handler/grpc"
	"github.com/vdaas/vald/pkg/manager/index/handler/rest"
	"github.com/vdaas/vald/pkg/manager/index/router"
	"github.com/vdaas/vald/pkg/manager/index/service"
)

type run struct {
	eg       errgroup.Group
	cfg      *config.Data
	server   starter.Server
	filter   service.Filter
	gateway  service.Gateway
	metadata service.Meta
	backup   service.Backup
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()

	var (
		filter   service.Filter
		gateway  service.Gateway
		metadata service.Meta
		backup   service.Backup
	)

	if addrs := cfg.Gateway.BackupManager.Client.Addrs; len(addrs) == 0 {
		return nil, errors.ErrInvalidBackupConfig
	}

	backup, err = service.NewBackup(
		service.WithBackupAddr(cfg.Gateway.BackupManager.Client.Addrs[0]),
		service.WithBackupClient(
			grpc.New(
				append(cfg.Gateway.BackupManager.Client.Opts(),
					grpc.WithErrGroup(eg),
				)...,
			),
		),
	)
	if err != nil {
		return nil, err
	}
	dscClient := grpc.New(
		append(cfg.Gateway.Discoverer.DiscoverClient.Opts(),
			grpc.WithErrGroup(eg),
		)...,
	)
	agentOpts := cfg.Gateway.Discoverer.AgentClient.Opts()
	gateway, err = service.NewGateway(
		service.WithErrGroup(eg),
		service.WithAgentName(cfg.Gateway.AgentName),
		service.WithAgentPort(cfg.Gateway.AgentPort),
		service.WithAgentServiceDNSARecord(cfg.Gateway.AgentDNS),
		service.WithDiscovererClient(dscClient),
		service.WithDiscovererHostPort(
			cfg.Gateway.Discoverer.Host,
			cfg.Gateway.Discoverer.Port,
		),
		service.WithDiscoverDuration(cfg.Gateway.Discoverer.Duration),
		service.WithAgentOptions(agentOpts...),
	)
	if err != nil {
		return nil, err
	}

	if addrs := cfg.Gateway.Meta.Client.Addrs; len(addrs) == 0 {
		return nil, errors.ErrInvalidMetaDataConfig
	}
	metadata, err = service.NewMeta(
		service.WithMetaAddr(cfg.Gateway.Meta.Client.Addrs[0]),
		service.WithMetaClient(
			grpc.New(
				append(cfg.Gateway.Meta.Client.Opts(),
					grpc.WithErrGroup(eg),
				)...,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	ef := cfg.Gateway.EgressFilter
	if ef != nil &&
		ef.Client != nil &&
		ef.Client.Addrs != nil &&
		len(ef.Client.Addrs) != 0 {
		filter, err = service.NewFilter(
			service.WithFilterClient(
				grpc.New(
					append(ef.Client.Opts(),
						grpc.WithErrGroup(eg),
					)...,
				),
			),
		)
	}

	v := handler.New(
		handler.WithGateway(gateway),
		handler.WithBackup(backup),
		handler.WithMeta(metadata),
		handler.WithFilters(filter),
		handler.WithErrGroup(eg),
		handler.WithReplicationCount(cfg.Gateway.IndexReplica),
	)

	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithREST(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithHTTPHandler(
					router.New(
						router.WithHandler(
							rest.New(
								rest.WithVald(v),
							),
						),
					),
				),
			}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(srv *grpc.Server) {
					vald.RegisterValdServer(srv, v)
				}),
				server.WithPreStopFunction(func() error {
					// TODO notify another gateway and scheduler
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
		eg:       eg,
		cfg:      cfg,
		server:   srv,
		filter:   filter,
		gateway:  gateway,
		metadata: metadata,
		backup:   backup,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 5)
	var bech, fech, mech, gech, sech <-chan error
	var err error
	if r.backup != nil {
		bech, err = r.backup.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.filter != nil {
		fech, err = r.filter.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.metadata != nil {
		mech, err = r.metadata.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.gateway != nil {
		gech, err = r.gateway.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	sech = r.server.ListenAndServe(ctx)
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return nil
			case err = <-bech:
			case err = <-fech:
			case err = <-gech:
			case err = <-mech:
			case err = <-sech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return nil
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
