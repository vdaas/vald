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

package usecase

import (
	"context"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/recover"
	"github.com/vdaas/vald/internal/net/grpc/metric"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/test/data/hdf5"
	"github.com/vdaas/vald/pkg/benchmark/job/search/config"
	handler "github.com/vdaas/vald/pkg/benchmark/job/search/handler/grpc"
	"github.com/vdaas/vald/pkg/benchmark/job/search/handler/rest"
	"github.com/vdaas/vald/pkg/benchmark/job/search/router"
	search "github.com/vdaas/vald/pkg/benchmark/job/search/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Config
	sj            search.SearchJob
	h             handler.Benchmark
	server        starter.Server
	observability observability.Observability
}

func New(cfg *config.Config) (r runner.Runner, err error) {
	log.Info("pkg/benchmark/job/search/cmd start")
	eg := errgroup.Get()
	copts, err := cfg.SearchJob.GatewayClient.Opts()
	if err != nil {
		return nil, err
	}

	c, err := vald.New(
		vald.WithAddrs(cfg.SearchJob.GatewayClient.Addrs...),
		vald.WithClient(grpc.New(copts...)),
	)
	if err != nil {
		return nil, err
	}

	// TODO: impl bind config
	d, err := hdf5.New()
	if err != nil {
		return nil, err
	}
	log.Info("pkg/benchmark/job/search/cmd success d")

	sj, err := search.New(
		search.WithErrGroup(eg),
		search.WithValdClient(c),
		search.WithDimension(cfg.SearchJob.Dimension),
		search.WithNum(cfg.SearchJob.Num),
		search.WithMinNum(cfg.SearchJob.MinNum),
		search.WithRadius(cfg.SearchJob.Radius),
		search.WithEpsilon(cfg.SearchJob.Epsilon),
		search.WithTimeout(cfg.SearchJob.Timeout),
		search.WithHdf5(d),
	)
	if err != nil {
		return nil, err
	}

	h, err := handler.New()
	if err != nil {
		return nil, err
	}

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			// TODO register grpc server handler here
		}),
		server.WithGRPCOption(
			grpc.ChainUnaryInterceptor(recover.RecoverInterceptor()),
			grpc.ChainStreamInterceptor(recover.RecoverStreamInterceptor()),
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
							// TODO pass grpc handler to REST option
							),
						),
					)),
			}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return grpcServerOptions
		}),
	)
	if err != nil {
		return nil, err
	}
	log.Info("pkg/benchmark/job/search/cmd end")

	return &run{
		eg:            eg,
		cfg:           cfg,
		sj:            sj,
		h:             h,
		server:        srv,
		observability: obs,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	if r.observability != nil {
		if err := r.observability.PreStart(ctx); err != nil {
			return err
		}
	}
	if r.sj != nil {
		return r.sj.PreStart(ctx)
	}
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 3)
	var oech, dech, sech <-chan error
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		if r.observability != nil {
			oech = r.observability.Start(ctx)
		}
		dech, err = r.sj.Start(ctx)

		if err != nil {
			ech <- err
			return err
		}

		r.h.Start(ctx)

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
