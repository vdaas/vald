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

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	vald "github.com/vdaas/vald/apis/grpc/v1/vald"
	iconf "github.com/vdaas/vald/internal/config"
<<<<<<< HEAD
=======
	"github.com/vdaas/vald/internal/errgroup"
>>>>>>> feature/agent/qbg
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	faissmetrics "github.com/vdaas/vald/internal/observability/metrics/agent/core/faiss"
	infometrics "github.com/vdaas/vald/internal/observability/metrics/info"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
<<<<<<< HEAD
	"github.com/vdaas/vald/internal/sync/errgroup"
=======
>>>>>>> feature/agent/qbg
	"github.com/vdaas/vald/pkg/agent/core/faiss/config"
	handler "github.com/vdaas/vald/pkg/agent/core/faiss/handler/grpc"
	"github.com/vdaas/vald/pkg/agent/core/faiss/handler/rest"
	"github.com/vdaas/vald/pkg/agent/core/faiss/router"
	"github.com/vdaas/vald/pkg/agent/core/faiss/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	faiss         service.Faiss
	server        starter.Server
	observability observability.Observability
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	faiss, err := service.New(
		cfg.Faiss,
		service.WithErrGroup(errgroup.Get()),
		service.WithEnableInMemoryMode(cfg.Faiss.EnableInMemoryMode),
		service.WithIndexPath(cfg.Faiss.IndexPath),
		service.WithAutoIndexCheckDuration(cfg.Faiss.AutoIndexCheckDuration),
		service.WithAutoSaveIndexDuration(cfg.Faiss.AutoSaveIndexDuration),
		service.WithAutoIndexDurationLimit(cfg.Faiss.AutoIndexDurationLimit),
		service.WithAutoIndexLength(cfg.Faiss.AutoIndexLength),
		service.WithInitialDelayMaxDuration(cfg.Faiss.InitialDelayMaxDuration),
		service.WithMinLoadIndexTimeout(cfg.Faiss.MinLoadIndexTimeout),
		service.WithMaxLoadIndexTimeout(cfg.Faiss.MaxLoadIndexTimeout),
		service.WithLoadIndexTimeoutFactor(cfg.Faiss.LoadIndexTimeoutFactor),
		service.WithProactiveGC(cfg.Faiss.EnableProactiveGC),
		service.WithCopyOnWrite(cfg.Faiss.EnableCopyOnWrite),
	)
	if err != nil {
		return nil, err
	}

	g, err := handler.New(
		handler.WithFaiss(faiss),
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
			faissmetrics.New(faiss),
			infometrics.New("agent_core_faiss_info", "Agent Faiss info", *cfg.Faiss),
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
	)
	if err != nil {
		return nil, err
	}

	return &run{
		eg:            eg,
		faiss:         faiss,
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
		nech = r.faiss.Start(ctx)
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
	r.faiss.Close(ctx)
	return nil
}
