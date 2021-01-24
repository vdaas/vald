//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/metric"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/gateway/meta/config"
	handler "github.com/vdaas/vald/pkg/gateway/meta/handler/grpc"
	"github.com/vdaas/vald/pkg/gateway/meta/handler/rest"
	"github.com/vdaas/vald/pkg/gateway/meta/router"
	"github.com/vdaas/vald/pkg/gateway/meta/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	server        starter.Server
	observability observability.Observability
	metadata      service.Meta
	client        client.Client
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()

	var metadata service.Meta

	mopts := cfg.Meta.Client.Opts()
	copts := cfg.Client.Opts()

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(cfg.Observability)
		if err != nil {
			return nil, err
		}
		copts = append(
			copts,
			grpc.WithDialOptions(
				grpc.WithStatsHandler(metric.NewClientHandler()),
			),
		)
		mopts = append(
			mopts,
			grpc.WithDialOptions(
				grpc.WithStatsHandler(metric.NewClientHandler()),
			),
		)
	}

	if len(cfg.Meta.Client.Addrs) == 0 {
		return nil, errors.ErrInvalidMetaDataConfig
	}
	metadata, err = service.New(
		service.WithMetaClient(
			grpc.New(mopts...),
		),
		service.WithMetaCacheEnabled(cfg.Meta.EnableCache),
		service.WithMetaCacheExpireDuration(cfg.Meta.CacheExpiration),
		service.WithMetaCacheExpiredCheckDuration(cfg.Meta.ExpiredCacheCheckDuration),
	)
	if err != nil {
		return nil, err
	}

	if len(cfg.Client.Addrs) == 0 {
		return nil, errors.ErrGRPCTargetAddrNotFound
	}

	c, err := client.New(
		client.WithAddrs(cfg.Client.Addrs...),
		client.WithClient(grpc.New(copts...)),
	)
	if err != nil {
		return nil, err
	}

	v := handler.New(
		handler.WithValdClient(c),
		handler.WithMeta(metadata),
		handler.WithErrGroup(eg),
		handler.WithStreamConcurrency(cfg.Server.GetGRPCStreamConcurrency()),
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

	if cfg.Observability.Enabled {
		grpcServerOptions = append(
			grpcServerOptions,
			server.WithGRPCOption(
				grpc.StatsHandler(metric.NewServerHandler()),
			),
		)
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
		metadata:      metadata,
		client:        c,
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
	var mech, sech, oech, cech <-chan error
	var err error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
	}
	if r.metadata != nil {
		mech, err = r.metadata.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.client != nil {
		cech, err = r.client.Start(ctx)
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
			case err = <-mech:
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
	return r.client.Stop(ctx)
}
