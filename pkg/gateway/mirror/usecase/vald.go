// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package usecase

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	mclient "github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	backoffmetrics "github.com/vdaas/vald/internal/observability/metrics/backoff"
	cbmetrics "github.com/vdaas/vald/internal/observability/metrics/circuitbreaker"
	mirrormetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/mirror"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/mirror/config"
	handler "github.com/vdaas/vald/pkg/gateway/mirror/handler/grpc"
	"github.com/vdaas/vald/pkg/gateway/mirror/handler/rest"
	"github.com/vdaas/vald/pkg/gateway/mirror/router"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

type run struct {
	eg            errgroup.Group
	der           net.Dialer
	cfg           *config.Data
	server        starter.Server
	c             mclient.Client
	gw            service.Gateway
	mirr          service.Mirror
	dsc           service.Discoverer
	observability observability.Observability
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()
	netOpts, err := cfg.Mirror.Net.Opts()
	if err != nil {
		return nil, err
	}
	der, err := net.NewDialer(netOpts...)
	if err != nil {
		return nil, err
	}

	cOpts, err := cfg.Mirror.Client.Opts()
	if err != nil {
		return nil, err
	}
	cOpts = append(cOpts, grpc.WithErrGroup(eg))

	c, err := mclient.New(
		mclient.WithAddrs(cfg.Mirror.Client.Addrs...),
		mclient.WithClient(grpc.New(cOpts...)),
	)
	if err != nil {
		return nil, err
	}

	gw, err := service.NewGateway(
		service.WithErrGroup(eg),
		service.WithMirrorClient(c),
		service.WithPodName(cfg.Mirror.PodName),
	)
	if err != nil {
		return nil, err
	}
	mirr, err := service.NewMirror(
		service.WithErrorGroup(eg),
		service.WithRegisterDuration(cfg.Mirror.RegisterDuration),
		service.WithValdAddrs(cfg.Mirror.GatewayAddr),
		service.WithSelfMirrorAddrs(cfg.Mirror.SelfMirrorAddr),
		service.WithGateway(gw),
	)
	if err != nil {
		return nil, err
	}
	dsc, err := service.NewDiscoverer(
		service.WithDiscovererNamespace(cfg.Mirror.Namespace),
		service.WithDiscovererGroup(cfg.Mirror.Group),
		service.WithDiscovererDuration(cfg.Mirror.DiscoveryDuration),
		service.WithDiscovererSelfMirrorAddrs(cfg.Mirror.SelfMirrorAddr),
		service.WithDiscovererColocation(cfg.Mirror.Colocation),
		service.WithDiscovererDialer(der),
		service.WithDiscovererMirror(mirr),
		service.WithDiscovererErrGroup(eg),
	)
	if err != nil {
		return nil, err
	}

	v, err := handler.New(
		handler.WithValdAddr(cfg.Mirror.GatewayAddr),
		handler.WithErrGroup(eg),
		handler.WithGateway(gw),
		handler.WithMirror(mirr),
		handler.WithStreamConcurrency(cfg.Server.GetGRPCStreamConcurrency()),
	)
	if err != nil {
		return nil, err
	}

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			vald.RegisterValdServerWithMirror(srv, v)
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
			mirrormetrics.New(mirr),
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
		der:           der,
		cfg:           cfg,
		server:        srv,
		c:             c,
		gw:            gw,
		mirr:          mirr,
		dsc:           dsc,
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
	ech := make(chan error, 6)
	var mech, dech, cech, sech, oech <-chan error
	var err error

	sech = r.server.ListenAndServe(ctx)
	if r.c != nil {
		cech, err = r.c.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.mirr != nil {
		mech = r.mirr.Start(ctx)
	}
	if r.dsc != nil {
		dech, err = r.dsc.Start(ctx)
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
			case err = <-mech:
			case err = <-dech:
			case err = <-cech:
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
	if r.observability != nil {
		r.observability.Stop(ctx)
	}
	return r.server.Shutdown(ctx)
}

func (*run) PostStop(_ context.Context) error {
	return nil
}
