//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	backoffmetrics "github.com/vdaas/vald/internal/observability/metrics/backoff"
	cbmetrics "github.com/vdaas/vald/internal/observability/metrics/circuitbreaker"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/discoverer/k8s/config"
	handler "github.com/vdaas/vald/pkg/discoverer/k8s/handler/grpc"
	"github.com/vdaas/vald/pkg/discoverer/k8s/handler/rest"
	"github.com/vdaas/vald/pkg/discoverer/k8s/router"
	"github.com/vdaas/vald/pkg/discoverer/k8s/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	dsc           service.Discoverer
	h             handler.DiscovererServer
	server        starter.Server
	observability observability.Observability
	der           net.Dialer
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()
	netOpts, err := cfg.Discoverer.Net.Opts()
	if err != nil {
		return nil, err
	}
	der, err := net.NewDialer(netOpts...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	dsc, err := service.New(
		cfg.Discoverer.Selectors,
		service.WithDiscoverDuration(cfg.Discoverer.DiscoveryDuration),
		service.WithErrGroup(eg),
		service.WithName(cfg.Discoverer.Name),
		service.WithNamespace(cfg.Discoverer.Namespace),
		service.WithDialer(der),
	)
	if err != nil {
		return nil, err
	}
	h, err := handler.New(
		handler.WithDiscoverer(dsc),
	)
	if err != nil {
		return nil, err
	}

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			discoverer.RegisterDiscovererServer(srv, h)
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
		starter.WithREST(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithHTTPHandler(
					router.New(
						router.WithTimeout(sc.HTTP.HandlerTimeout),
						router.WithErrGroup(eg),
						router.WithHandler(
							rest.New(
								rest.WithDiscoverer(h),
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
		der:           der,
		eg:            eg,
		cfg:           cfg,
		dsc:           dsc,
		h:             h,
		server:        srv,
		observability: obs,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	if r.der != nil {
		r.der.StartDialerCache(ctx)
	}
	if r.observability != nil {
		return r.observability.PreStart(ctx)
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
		dech, err = r.dsc.Start(ctx)
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
