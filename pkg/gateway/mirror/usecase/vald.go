// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	bometrics "github.com/vdaas/vald/internal/observability/metrics/backoff"
	cbmetrics "github.com/vdaas/vald/internal/observability/metrics/circuitbreaker"
	mirrmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/mirror"
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
	dialer        net.Dialer
	cfg           *config.Data
	server        starter.Server
	client        mirror.Client
	gateway       service.Gateway
	mirror        service.Mirror
	discover      service.Discovery
	observability observability.Observability
}

// New returns Runner instance.
func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()

	netOpts, err := cfg.Mirror.Net.Opts()
	if err != nil {
		return nil, err
	}
	dialer, err := net.NewDialer(netOpts...)
	if err != nil {
		return nil, err
	}

	cOpts, err := cfg.Mirror.Client.Opts()
	if err != nil {
		return nil, err
	}
	// skipcq: CRT-D0001
	cOpts = append(cOpts, grpc.WithErrGroup(eg))

	client, err := mirror.New(
		mirror.WithAddrs(cfg.Mirror.Client.Addrs...),
		mirror.WithClient(grpc.New(cOpts...)),
	)
	if err != nil {
		return nil, err
	}

	gateway, err := service.NewGateway(
		service.WithErrGroup(eg),
		service.WithMirrorClient(client),
		service.WithPodName(cfg.Mirror.PodName),
	)
	if err != nil {
		return nil, err
	}
	mirror, err := service.NewMirror(
		service.WithErrorGroup(eg),
		service.WithRegisterDuration(cfg.Mirror.RegisterDuration),
		service.WithGatewayAddrs(cfg.Mirror.GatewayAddr),
		service.WithSelfMirrorAddrs(cfg.Mirror.SelfMirrorAddr),
		service.WithGateway(gateway),
	)
	if err != nil {
		return nil, err
	}
	discover, err := service.NewDiscovery(
		service.WithDiscoveryNamespace(cfg.Mirror.Namespace),
		service.WithDiscoveryGroup(cfg.Mirror.Group),
		service.WithDiscoveryDuration(cfg.Mirror.DiscoveryDuration),
		service.WithDiscoverySelfMirrorAddrs(cfg.Mirror.SelfMirrorAddr),
		service.WithDiscoveryColocation(cfg.Mirror.Colocation),
		service.WithDiscoveryDialer(dialer),
		service.WithDiscoveryMirror(mirror),
		service.WithDiscoveryErrGroup(eg),
	)
	if err != nil {
		return nil, err
	}

	v, err := handler.New(
		handler.WithValdAddr(cfg.Mirror.GatewayAddr),
		handler.WithErrGroup(eg),
		handler.WithGateway(gateway),
		handler.WithMirror(mirror),
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
			bometrics.New(),
			cbmetrics.New(),
			mirrmetrics.New(mirror),
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
		dialer:        dialer,
		cfg:           cfg,
		server:        srv,
		client:        client,
		gateway:       gateway,
		mirror:        mirror,
		discover:      discover,
		observability: obs,
	}, nil
}

// PreStart is a method called before execution of Start.
func (r *run) PreStart(ctx context.Context) error {
	if r.dialer != nil {
		r.dialer.StartDialerCache(ctx)
	}
	if r.observability != nil {
		return r.observability.PreStart(ctx)
	}
	return nil
}

// Start is a method used to initiate an operation in the run, and it returns a channel for receiving errors
// during the operation and an error representing any initialization errors.
func (r *run) Start(ctx context.Context) (_ <-chan error, err error) { // skipcq: GO-R1005
	ech := make(chan error, 6)
	var mech, dech, cech, sech, oech <-chan error

	sech = r.server.ListenAndServe(ctx)
	if r.client != nil {
		cech, err = r.client.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.mirror != nil {
		mech = r.mirror.Start(ctx)
	}
	if r.discover != nil {
		dech, err = r.discover.Start(ctx)
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
					return errors.Join(ctx.Err(), err)
				case ech <- err:
				}
			}
		}
	}))
	return ech, nil
}

// PreStop is a method called before execution of Stop.
func (*run) PreStop(_ context.Context) error {
	return nil
}

// Stop is a method used to stop an operation in the run.
func (r *run) Stop(ctx context.Context) (errs error) {
	if r.observability != nil {
		if err := r.observability.Stop(ctx); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if r.server != nil {
		if err := r.server.Shutdown(ctx); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

// PtopStop is a method called after execution of Stop.
func (*run) PostStop(_ context.Context) error {
	return nil
}
