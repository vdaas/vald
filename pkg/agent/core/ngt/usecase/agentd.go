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

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	vald "github.com/vdaas/vald/apis/grpc/v1/vald"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	ngtmetrics "github.com/vdaas/vald/internal/observability/metrics/agent/core/ngt"
	infometrics "github.com/vdaas/vald/internal/observability/metrics/info"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/agent/core/ngt/config"
	handler "github.com/vdaas/vald/pkg/agent/core/ngt/handler/grpc"
	"github.com/vdaas/vald/pkg/agent/core/ngt/handler/rest"
	"github.com/vdaas/vald/pkg/agent/core/ngt/router"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	ngt           service.NGT
	server        starter.Server
	observability observability.Observability
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	log.Info("Test")
	ngt, err := service.New(
		cfg.NGT,
		service.WithErrGroup(errgroup.Get()),
		service.WithEnableInMemoryMode(cfg.NGT.EnableInMemoryMode),
		service.WithIndexPath(cfg.NGT.IndexPath),
		service.WithAutoIndexCheckDuration(cfg.NGT.AutoIndexCheckDuration),
		service.WithAutoIndexDurationLimit(cfg.NGT.AutoIndexDurationLimit),
		service.WithAutoSaveIndexDuration(cfg.NGT.AutoSaveIndexDuration),
		service.WithAutoIndexLength(cfg.NGT.AutoIndexLength),
		service.WithInitialDelayMaxDuration(cfg.NGT.InitialDelayMaxDuration),
		service.WithMinLoadIndexTimeout(cfg.NGT.MinLoadIndexTimeout),
		service.WithMaxLoadIndexTimeout(cfg.NGT.MaxLoadIndexTimeout),
		service.WithLoadIndexTimeoutFactor(cfg.NGT.LoadIndexTimeoutFactor),
		service.WithDefaultPoolSize(cfg.NGT.DefaultPoolSize),
		service.WithDefaultRadius(cfg.NGT.DefaultRadius),
		service.WithDefaultEpsilon(cfg.NGT.DefaultEpsilon),
		service.WithProactiveGC(cfg.NGT.EnableProactiveGC),
		service.WithCopyOnWrite(cfg.NGT.EnableCopyOnWrite),
	)
	if err != nil {
		return nil, err
	}
	g, err := handler.New(
		handler.WithNGT(ngt),
		handler.WithStreamConcurrency(cfg.Server.GetGRPCStreamConcurrency()),
	)
	if err != nil {
		return nil, err
	}
	eg := errgroup.Get()

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			agent.RegisterAgentServer(srv, g)
			vald.RegisterValdServer(srv, g)
		}),
		server.WithPreStartFunc(func() error {
			return nil
		}),
		server.WithPreStopFunction(func() error {
			return nil
		}),
	}

	var obs observability.Observability
	if cfg.Observability != nil && cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(
			cfg.Observability,
			ngtmetrics.New(ngt),
			infometrics.New("agent_core_ngt_info", "Agent NGT info", *cfg.NGT),
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
								rest.WithAgent(g),
							),
						),
					),
				),
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
		eg:            eg,
		ngt:           ngt,
		cfg:           cfg,
		server:        srv,
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
	ech := make(chan error, 3)
	var oech, nech, sech <-chan error
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		if r.observability != nil {
			oech = r.observability.Start(ctx)
		}
		nech = r.ngt.Start(ctx)
		sech = r.server.ListenAndServe(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-oech:
			case err = <-nech:
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

func (*run) PreStop(ctx context.Context) error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	if r.observability != nil {
		r.observability.Stop(ctx)
	}
	return r.server.Shutdown(ctx)
}

func (r *run) PostStop(ctx context.Context) error {
	r.ngt.Close(ctx)
	return nil
}
