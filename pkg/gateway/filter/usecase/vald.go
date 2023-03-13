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

package usecase

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/filter/egress"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	backoffmetrics "github.com/vdaas/vald/internal/observability/metrics/backoff"
	cbmetrics "github.com/vdaas/vald/internal/observability/metrics/circuitbreaker"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/slices"
	"github.com/vdaas/vald/pkg/gateway/filter/config"
	handler "github.com/vdaas/vald/pkg/gateway/filter/handler/grpc"
	"github.com/vdaas/vald/pkg/gateway/filter/handler/rest"
	"github.com/vdaas/vald/pkg/gateway/filter/router"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	server        starter.Server
	observability observability.Observability
	client        client.Client
	ingress       ingress.Client
	egress        egress.Client
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
	icopts, err := cfg.IngressFilters.Client.Opts()
	if err != nil {
		return nil, err
	}
	ecopts, err := cfg.EgressFilters.Client.Opts()
	if err != nil {
		return nil, err
	}

	c, err := client.New(
		client.WithAddrs(cfg.Client.Addrs...),
		client.WithClient(grpc.New(copts...)),
	)
	if err != nil {
		return nil, err
	}
	var iaddrs []string
	if cfg.IngressFilters != nil {
		as := make([]string, 0, 10)
		if cfg.IngressFilters.Client != nil && cfg.IngressFilters.Client.Addrs != nil {
			as = append(as, cfg.IngressFilters.Client.Addrs...)
		}
		if cfg.IngressFilters.Vectorizer != "" {
			as = append(as, cfg.IngressFilters.Vectorizer)
		}
		if cfg.IngressFilters.SearchFilters != nil {
			as = append(as, cfg.IngressFilters.SearchFilters...)
		}
		if cfg.IngressFilters.InsertFilters != nil {
			as = append(as, cfg.IngressFilters.InsertFilters...)
		}
		if cfg.IngressFilters.UpdateFilters != nil {
			as = append(as, cfg.IngressFilters.UpdateFilters...)
		}
		if cfg.IngressFilters.UpsertFilters != nil {
			as = append(as, cfg.IngressFilters.UpsertFilters...)
		}
		if len(as) != 0 {
			slices.Sort(as)
			dup := make(map[string]bool, len(as))
			iaddrs = make([]string, 0, len(as))
			for _, addr := range as {
				if !dup[addr] {
					dup[addr] = true
					iaddrs = append(iaddrs, addr)
				}
			}
		}
	}
	ic, err := ingress.New(
		ingress.WithAddrs(iaddrs...),
		ingress.WithClient(grpc.New(icopts...)),
	)
	if err != nil {
		return nil, err
	}
	var eaddrs []string
	if cfg.EgressFilters != nil {
		as := make([]string, 0, 10)
		if cfg.EgressFilters.Client != nil && cfg.EgressFilters.Client.Addrs != nil {
			as = append(as, cfg.EgressFilters.Client.Addrs...)
		}
		if cfg.EgressFilters.DistanceFilters != nil {
			as = append(as, cfg.EgressFilters.DistanceFilters...)
		}
		if cfg.EgressFilters.ObjectFilters != nil {
			as = append(as, cfg.EgressFilters.ObjectFilters...)
		}
		if len(as) != 0 {
			slices.Sort(as)
			dup := make(map[string]bool, len(as))
			eaddrs = make([]string, 0, len(as))
			for _, addr := range as {
				if !dup[addr] {
					dup[addr] = true
					eaddrs = append(eaddrs, addr)
				}
			}
		}
	}
	ec, err := egress.New(
		egress.WithAddrs(eaddrs...),
		egress.WithClient(grpc.New(ecopts...)),
	)
	if err != nil {
		return nil, err
	}

	v := handler.New(
		handler.WithValdClient(c),
		handler.WithEgressFilterClient(ec),
		handler.WithIngressFilterClient(ic),
		handler.WithErrGroup(eg),
		handler.WithStreamConcurrency(cfg.Server.GetGRPCStreamConcurrency()),
		handler.WithVectorizerTargets(cfg.IngressFilters.Vectorizer),
		handler.WithSearchFilterTargets(cfg.IngressFilters.SearchFilters...),
		handler.WithInsertFilterTargets(cfg.IngressFilters.InsertFilters...),
		handler.WithUpdateFilterTargets(cfg.IngressFilters.UpdateFilters...),
		handler.WithUpsertFilterTargets(cfg.IngressFilters.UpsertFilters...),
		handler.WithDistanceFilterTargets(cfg.EgressFilters.DistanceFilters...),
		handler.WithObjectFilterTargets(cfg.EgressFilters.ObjectFilters...),
	)

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			vald.RegisterValdServerWithFilter(srv, v)
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
		ingress:       ic,
		egress:        ec,
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
	var eech, iech, sech, oech, cech <-chan error
	var err error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
	}
	if r.ingress != nil {
		iech, err = r.ingress.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	if r.egress != nil {
		eech, err = r.egress.Start(ctx)
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
			case err = <-iech:
			case err = <-eech:
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
		if err != nil {
			err = errors.Wrap(r.ingress.Stop(ctx), errors.Wrap(r.egress.Stop(ctx), err.Error()).Error())
			return
		}
		err = r.ingress.Stop(ctx)
		if err != nil {
			err = errors.Wrap(r.egress.Stop(ctx), err.Error())
			return
		}
		err = r.egress.Stop(ctx)
	}()
	return r.client.Stop(ctx)
}
