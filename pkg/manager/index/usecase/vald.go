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

	"github.com/vdaas/vald/apis/grpc/manager/index"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
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
	eg      errgroup.Group
	cfg     *config.Data
	server  starter.Server
	indexer service.Indexer
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()

	var (
		indexer service.Indexer
	)

	dscClient := grpc.New(
		append(cfg.Indexer.Discoverer.Client.Opts(),
			grpc.WithErrGroup(eg),
		)...,
	)
	agentOpts := cfg.Indexer.Discoverer.AgentClient.Opts()
	indexer, err = service.New(
		service.WithErrGroup(eg),
		service.WithIndexingConcurrency(cfg.Indexer.Concurrency),
		service.WithIndexingDuration(cfg.Indexer.AutoIndexCheckDuration),
		service.WithIndexingDurationLimit(cfg.Indexer.AutoIndexDurationLimit),
		service.WithMinUncommitted(cfg.Indexer.AutoIndexLength),
		service.WithAgentName(cfg.Indexer.AgentName),
		service.WithAgentNamespace(cfg.Indexer.AgentNamespace),
		service.WithAgentPort(cfg.Indexer.AgentPort),
		service.WithAgentServiceDNSARecord(cfg.Indexer.AgentDNS),
		service.WithNodeName(cfg.Indexer.NodeName),
		service.WithDiscovererClient(dscClient),
		service.WithDiscovererHostPort(
			cfg.Indexer.Discoverer.Host,
			cfg.Indexer.Discoverer.Port,
		),
		service.WithDiscoverDuration(cfg.Indexer.Discoverer.Duration),
		service.WithAgentOptions(agentOpts...),
	)
	if err != nil {
		return nil, err
	}
	idx := handler.New(handler.WithIndexer(indexer))

	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithREST(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithHTTPHandler(
					router.New(
						router.WithHandler(
							rest.New(
								rest.WithIndexer(idx),
							),
						),
					),
				),
			}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(srv *grpc.Server) {
					index.RegisterIndexServer(srv, idx)
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
		eg:      eg,
		cfg:     cfg,
		server:  srv,
		indexer: indexer,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 5)
	var iech, sech <-chan error
	var err error
	if r.indexer != nil {
		iech, err = r.indexer.Start(ctx)
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
				return ctx.Err()
			case err = <-iech:
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
	return r.server.Shutdown(ctx)
}

func (r *run) PostStop(ctx context.Context) error {
	return nil
}
