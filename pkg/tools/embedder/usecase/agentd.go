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

	embedderpb "github.com/vdaas/vald/apis/grpc/v1/embedder"
	vclient "github.com/vdaas/vald/internal/client/v1/client/vald"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	infometrics "github.com/vdaas/vald/internal/observability/metrics/info"
	memmetrics "github.com/vdaas/vald/internal/observability/metrics/mem"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/tools/embedder/config"
	handler "github.com/vdaas/vald/pkg/tools/embedder/handler/grpc"
	"github.com/vdaas/vald/pkg/tools/embedder/handler/rest"
	"github.com/vdaas/vald/pkg/tools/embedder/router"
	"github.com/vdaas/vald/pkg/tools/embedder/service"
)

type run struct {
	eg            errgroup.Group
	valdClient    vclient.Client
	metaClient    service.MetaClient
	server        starter.Server
	observability observability.Observability
}

func New(cfg *config.Data) (runner.Runner, error) {
	eg := errgroup.Get()
	copts, err := cfg.Client.Opts()
	if err != nil {
		return nil, err
	}
	valdClient, err := vclient.New(vclient.WithClient(grpc.New("Embedder Vald Client", append(copts, grpc.WithErrGroup(eg))...)))
	if err != nil {
		return nil, err
	}
	openAI, err := service.NewOpenAI(service.WithToken(cfg.LLM.OpenAI.Token), service.WithOpenAIModel(cfg.LLM.OpenAI.Model))
	if err != nil {
		return nil, err
	}
	serviceOpts := []service.Option{service.WithValdClient(valdClient), service.WithLLM(openAI)}
	var metaClient service.MetaClient
	if cfg.Meta != nil && cfg.Meta.Client != nil && len(cfg.Meta.Client.Addrs) != 0 {
		mopts, err := cfg.Meta.Client.Opts()
		if err != nil {
			return nil, err
		}
		metaClient, err = service.NewMetaClient(grpc.New("Embedder Meta Client", append(mopts, grpc.WithErrGroup(eg))...))
		if err != nil {
			return nil, err
		}
		serviceOpts = append(serviceOpts, service.WithMetaClient(metaClient))
	}
	embedderService, err := service.New(serviceOpts...)
	if err != nil {
		return nil, err
	}
	g, err := handler.New(handler.WithEmbedder(embedderService))
	if err != nil {
		return nil, err
	}
	var obs observability.Observability
	if cfg.Observability != nil && cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(cfg.Observability, infometrics.New("embedder_info", "Embedder info", cfg.GlobalConfig), memmetrics.New())
		if err != nil {
			return nil, err
		}
	}
	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithREST(func(sc *iconf.Server) []server.Option {
			return []server.Option{server.WithHTTPHandler(router.New(router.WithTimeout(sc.HTTP.HandlerTimeout), router.WithErrGroup(eg), router.WithHandler(rest.New(rest.WithEmbedder(g)))))}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return []server.Option{server.WithGRPCRegisterar(func(srv *grpc.Server) { embedderpb.RegisterEmbedderServer(srv, g) })}
		}),
	)
	if err != nil {
		return nil, err
	}
	return &run{eg: eg, valdClient: valdClient, metaClient: metaClient, server: srv, observability: obs}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	if r.observability != nil {
		return r.observability.PreStart(ctx)
	}
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 4)
	var oech, vech, mech, sech <-chan error
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		if r.observability != nil {
			oech = r.observability.Start(ctx)
		}
		vech, err = r.valdClient.Start(ctx)
		if err != nil {
			return err
		}
		if r.metaClient != nil {
			mech, err = r.metaClient.Start(ctx)
			if err != nil {
				return err
			}
		}
		sech = r.server.ListenAndServe(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-oech:
			case err = <-vech:
			case err = <-mech:
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
	if r.metaClient != nil {
		if err := r.metaClient.Stop(ctx); err != nil {
			return err
		}
	}
	if r.valdClient != nil {
		if err := r.valdClient.Stop(ctx); err != nil {
			return err
		}
	}
	return r.server.Shutdown(ctx)
}

func (*run) PostStop(context.Context) error {
	return nil
}
