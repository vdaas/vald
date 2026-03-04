//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/meta/tikv"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	backoffmetrics "github.com/vdaas/vald/internal/observability/metrics/backoff"
	cbmetrics "github.com/vdaas/vald/internal/observability/metrics/circuitbreaker"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/meta/config"
	handler "github.com/vdaas/vald/pkg/gateway/meta/handler/grpc"
	"github.com/vdaas/vald/pkg/gateway/meta/handler/rest"
	"github.com/vdaas/vald/pkg/gateway/meta/router"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	server        starter.Server
	observability observability.Observability
	client        client.Client
	metaClient    tikv.Client
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	if addrs := cfg.Client.Addrs; len(addrs) == 0 {
		return nil, errors.ErrGRPCTargetAddrNotFound
	}
	eg := errgroup.Get()
	copts, err := cfg.Client.Opts()
	if err != nil {
		return nil, err
	}

	c, err := client.New(
		client.WithAddrs(cfg.Client.Addrs...),
		client.WithClient(grpc.New("Gateway Client", copts...)),
	)
	if err != nil {
		return nil, err
	}

	mc, err := tikv.New(
		tikv.WithPDAddrs(cfg.MetadataStore.Addrs...),
	)
	if err != nil {
		return nil, err
	}

	v := handler.New(
		handler.WithValdClient(c),
		handler.WithMetadataClient(mc),
		handler.WithErrGroup(eg),
		handler.WithStreamConcurrency(cfg.Server.GetGRPCStreamConcurrency()),
	)

	grpcServerOptions := []server.Option{
		server.WithGRPCRegisterar(func(srv *grpc.Server) {
			vald.RegisterValdServerWithMetadata(srv, v)
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
		// TODO add GraphQL handler
	)
	if err != nil {
		return nil, err
	}

	return &run{
		eg:            eg,
		cfg:           cfg,
		server:        srv,
		observability: obs,
		client:        c,
		metaClient:    mc,
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
	var sech, oech, cech <-chan error
	var err error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
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

func (*run) PreStop(context.Context) error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	if r.observability != nil {
		r.observability.Stop(ctx)
	}
	return r.server.Shutdown(ctx)
}

func (r *run) PostStop(ctx context.Context) (err error) {
	defer func() {
		if r.metaClient != nil {
			cerr := r.metaClient.Close()
			if cerr != nil {
				if err != nil {
					err = errors.Join(err, cerr)
				} else {
					err = cerr
				}
			}
		}
	}()
	if r.client != nil {
		return r.client.Stop(ctx)
	}
	return nil
}
