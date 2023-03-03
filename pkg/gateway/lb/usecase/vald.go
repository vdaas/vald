//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package usecase represents gateways usecase layer
package usecase

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	backoffmetrics "github.com/vdaas/vald/internal/observability/metrics/backoff"
	cbmetrics "github.com/vdaas/vald/internal/observability/metrics/circuitbreaker"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/gateway/lb/config"
	handler "github.com/vdaas/vald/pkg/gateway/lb/handler/grpc"
	"github.com/vdaas/vald/pkg/gateway/lb/handler/rest"
	"github.com/vdaas/vald/pkg/gateway/lb/router"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	server        starter.Server
	observability observability.Observability
	gateway       service.Gateway
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()

	var gateway service.Gateway

	cOpts, err := cfg.Gateway.Discoverer.Client.Opts()
	if err != nil {
		return nil, err
	}
	// skipcq: CRT-D0001
	dopts := append(
		cOpts,
		grpc.WithErrGroup(eg))
	acOpts, err := cfg.Gateway.Discoverer.AgentClientOptions.Opts()
	if err != nil {
		return nil, err
	}
	// skipcq: CRT-D0001
	aopts := append(
		acOpts,
		grpc.WithErrGroup(eg))

	client, err := discoverer.New(
		discoverer.WithAutoConnect(true),
		discoverer.WithName(cfg.Gateway.AgentName),
		discoverer.WithNamespace(cfg.Gateway.AgentNamespace),
		discoverer.WithPort(cfg.Gateway.AgentPort),
		discoverer.WithServiceDNSARecord(cfg.Gateway.AgentDNS),
		discoverer.WithDiscovererClient(grpc.New(dopts...)),
		discoverer.WithDiscoverDuration(cfg.Gateway.Discoverer.Duration),
		discoverer.WithOptions(aopts...),
		discoverer.WithNodeName(cfg.Gateway.NodeName),
	)
	if err != nil {
		return nil, err
	}
	gateway, err = service.NewGateway(
		service.WithErrGroup(eg),
		service.WithDiscoverer(client),
	)
	if err != nil {
		return nil, err
	}

	v := handler.New(
		handler.WithGateway(gateway),
		handler.WithErrGroup(eg),
		handler.WithReplicationCount(cfg.Gateway.IndexReplica),
		handler.WithStreamConcurrency(cfg.Server.GetGRPCStreamConcurrency()),
		handler.WithMultiConcurrency(cfg.Gateway.MultiOperationConcurrency),
	)

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			vald.RegisterValdServer(srv, v)
		}),
		server.WithPreStopFunction(func() error {
			// TODO notify another gateway and scheduler
			return nil
		}),
	}

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(
			cfg.Observability,
			backoffmetrics.New(),
			cbmetrics.New(),
		)
		if err != nil {
			return nil, err
		}
	}

	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithREST(func(sc *config.Server) []server.Option {
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
		starter.WithGRPC(func(sc *config.Server) []server.Option {
			return grpcServerOptions
		}),
	)
	if err != nil {
		return nil, err
	}

	return &run{
		eg:            eg,
		cfg:           cfg,
		server:        srv,
		observability: obs,
		gateway:       gateway,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	if r.observability != nil {
		return r.observability.PreStart(ctx)
	}
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 6)
	var gech, sech, oech <-chan error
	var err error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
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
				return ctx.Err()
			case err = <-oech:
			case err = <-gech:
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

func (*run) PreStop(context.Context) error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	if r.observability != nil {
		r.observability.Stop(ctx)
	}
	return r.server.Shutdown(ctx)
}

func (*run) PostStop(context.Context) error {
	return nil
}
