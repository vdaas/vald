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

	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/recover"
	"github.com/vdaas/vald/internal/observability"
	backoffmetrics "github.com/vdaas/vald/internal/observability/metrics/backoff"
	infometrics "github.com/vdaas/vald/internal/observability/metrics/info"
	benchmarkmetrics "github.com/vdaas/vald/internal/observability/metrics/tools/benchmark"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/operator/benchmark/config"
	handler "github.com/vdaas/vald/pkg/operator/benchmark/handler/grpc"
	"github.com/vdaas/vald/pkg/operator/benchmark/handler/rest"
	"github.com/vdaas/vald/pkg/operator/benchmark/router"
	"github.com/vdaas/vald/pkg/operator/benchmark/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	operator      service.Operator
	h             handler.Benchmark
	server        starter.Server
	observability observability.Observability
}

var JOB_NAMESPACE = os.Getenv("JOB_NAMESPACE")

func New(cfg *config.Data) (r runner.Interface, err error) {
	log.Info("pkg/operator/benchmark/cmd start")

	eg := errgroup.Get()

	log.Info("pkg/operator/benchmark/cmd success")

	operator, err := service.New(
		service.WithErrGroup(eg),
		service.WithJobNamespace(JOB_NAMESPACE),
		service.WithJobImageRepository(cfg.Job.Image.Repository),
		service.WithJobImageTag(cfg.Job.Image.Tag),
		service.WithJobImagePullPolicy(cfg.Job.Image.PullPolicy),
	)
	if err != nil {
		return nil, err
	}

	h, err := handler.New(handler.WithOperator(operator))
	if err != nil {
		return nil, err
	}

	grpcServerOptions := []server.Option{
		server.WithGRPCRegisterar(func(srv *grpc.Server) {
			// TODO register grpc server handler here
		}),
		server.WithGRPCOption(
			grpc.ChainUnaryInterceptor(recover.Interceptor()),
			grpc.ChainStreamInterceptor(recover.StreamInterceptor()),
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
		obs, err = observability.NewWithConfig(
			cfg.Observability,
			benchmarkmetrics.New(operator),
			infometrics.New(config.BenchmarkOperatorInfo, "Benchmark Operator info", *cfg.Job.Image),
			backoffmetrics.New(),
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
								rest.WithBenchmark(h),
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
	log.Info("pkg/operator/benchmark/cmd end")

	return &run{
		eg:            eg,
		cfg:           cfg,
		operator:      operator,
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
	if r.operator != nil {
		return r.operator.PreStart(ctx)
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

		dech, err = r.operator.Start(ctx)
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
					log.Error(err)
					return errors.Wrap(ctx.Err(), err.Error())
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

func (*run) PostStop(context.Context) error {
	return nil
}
