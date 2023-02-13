//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	mclient "github.com/vdaas/vald/internal/client/v1/client/mirror"
	vclient "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	backoffmetrics "github.com/vdaas/vald/internal/observability/metrics/backoff"
	cbmetrics "github.com/vdaas/vald/internal/observability/metrics/circuitbreaker"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/gateway/mirror/config"
	handler "github.com/vdaas/vald/pkg/gateway/mirror/handler/grpc"
	"github.com/vdaas/vald/pkg/gateway/mirror/handler/rest"
	"github.com/vdaas/vald/pkg/gateway/mirror/router"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	server        starter.Server
	lbc           vclient.Client
	gateway       service.Gateway
	observability observability.Observability
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()

	mcOpts, err := cfg.Mirror.Mirror.Opts()
	if err != nil {
		return nil, err
	}
	mcOpts = append(mcOpts, grpc.WithErrGroup(eg))

	mc, err := mclient.New(
		mclient.WithAddrs(cfg.Mirror.Mirror.Addrs...),
		mclient.WithClient(grpc.New(mcOpts...)),
	)
	if err != nil {
		return nil, err
	}

	selfMcOpts, err := cfg.Mirror.SelfMirror.Opts()
	if err != nil {
		return nil, err
	}
	selfMcOpts = append(selfMcOpts, grpc.WithErrGroup(eg))

	selfMc, err := mclient.New(
		mclient.WithAddrs(cfg.Mirror.SelfMirror.Addrs...),
		mclient.WithClient(grpc.New(selfMcOpts...)),
	)
	if err != nil {
		return nil, err
	}

	lbcOpts, err := cfg.Mirror.LB.Opts()
	if err != nil {
		return nil, err
	}
	lbcOpts = append(lbcOpts, grpc.WithErrGroup(eg))

	lbc, err := vclient.New(
		vclient.WithAddrs(cfg.Mirror.LB.Addrs...),
		vclient.WithClient(grpc.New(lbcOpts...)),
	)
	if err != nil {
		return nil, err
	}

	gateway, err := service.NewGateway(
		service.WithErrGroup(eg),
		service.WithSelfMirror(selfMc),
		service.WithMirror(mc),
		service.WithPodName(cfg.Mirror.PodName),
		service.WithAdvertiseInterval(cfg.Mirror.AdvertiseInterval),
	)
	if err != nil {
		return nil, err
	}

	v := handler.New(
		handler.WithValdClient(lbc),
		handler.WithErrGroup(eg),
		handler.WithGateway(gateway),
		handler.WithStreamConcurrency(cfg.Server.GetGRPCStreamConcurrency()),
	)

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			mirror.RegisterValdServerWithMirror(srv, v)
		}),
		server.WithPreStopFunction(func() error {
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
		lbc:           lbc,
		gateway:       gateway,
		observability: obs,
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
	var gech, lech, sech, oech <-chan error
	var err error

	sech = r.server.ListenAndServe(ctx)
	if r.gateway != nil {
		gech, err = r.gateway.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.lbc != nil {
		lech, err = r.lbc.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.observability != nil {
		oech = r.observability.Start(ctx)
	}

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-oech:
			case err = <-gech:
			case err = <-sech:
			case err = <-lech:
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
