// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"slices"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/filter/egress"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
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
	stopClients   context.CancelFunc
}

func New(cfg *config.Data) (r runner.Interface, err error) {
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
		client.WithClient(grpc.New("Gateway Client", copts...)),
	)
	if err != nil {
		return nil, err
	}
	var iaddrs []string
	var iidle string
	if f := cfg.IngressFilters; f != nil {
		var as []string
		if f.Client != nil {
			as = append(as, f.Client.Addrs...)
			if f.Client.DialOption != nil {
				iidle = f.Client.DialOption.IdleTimeout
			}
		}
		if f.Vectorizer != "" {
			as = append(as, f.Vectorizer)
		}
		iaddrs = sortedUnique(slices.Concat(as, f.SearchFilters, f.InsertFilters, f.UpdateFilters, f.UpsertFilters))
	}
	ic, err := ingress.New(
		ingress.WithAddrs(iaddrs...),
		ingress.WithClient(grpc.New("Ingress Filter Client", icopts...)),
		ingress.WithIdleTimeout(iidle),
	)
	if err != nil {
		return nil, err
	}
	var eaddrs []string
	var eidle string
	if f := cfg.EgressFilters; f != nil {
		var as []string
		if f.Client != nil {
			as = append(as, f.Client.Addrs...)
			if f.Client.DialOption != nil {
				eidle = f.Client.DialOption.IdleTimeout
			}
		}
		eaddrs = sortedUnique(slices.Concat(as, f.DistanceFilters, f.ObjectFilters))
	}
	ec, err := egress.New(
		egress.WithAddrs(eaddrs...),
		egress.WithClient(grpc.New("Egress Filter Client", ecopts...)),
		egress.WithIdleTimeout(eidle),
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
		server.WithGRPCRegisterar(func(srv *grpc.Server) {
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
	// ctx is cancelled by the shutdown signal, before the server drains its
	// in-flight requests, and a cancelled client closes its connections. Run
	// the outbound clients on a context that only PostStop cancels.
	cctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	r.stopClients = cancel
	if r.ingress != nil {
		iech, err = r.ingress.Start(cctx)
		if err != nil {
			cancel()
			close(ech)
			return nil, err
		}
	}
	if r.egress != nil {
		eech, err = r.egress.Start(cctx)
		if err != nil {
			cancel()
			close(ech)
			return nil, err
		}
	}
	if r.client != nil {
		cech, err = r.client.Start(cctx)
		if err != nil {
			cancel()
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

func (r *run) Stop(ctx context.Context) (err error) {
	// Shut the server down first so spans and metrics of the requests it drains
	// are still exported before observability is stopped.
	err = r.server.Shutdown(ctx)
	if r.observability != nil {
		err = errors.Join(err, r.observability.Stop(ctx))
	}
	return err
}

func (r *run) PostStop(ctx context.Context) error {
	if r.stopClients != nil {
		defer r.stopClients()
	}
	// All three clients are always stopped, in this order, and their errors joined.
	return errors.Join(r.client.Stop(ctx), r.ingress.Stop(ctx), r.egress.Stop(ctx))
}

// sortedUnique sorts addrs and drops duplicates; it returns nil for no addresses.
func sortedUnique(addrs []string) []string {
	if len(addrs) == 0 {
		return nil
	}
	slices.Sort(addrs)
	return slices.Compact(addrs)
}
