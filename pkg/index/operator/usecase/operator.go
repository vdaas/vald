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

	iconfig "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/recover"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/index/operator/config"
	"github.com/vdaas/vald/pkg/index/operator/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	observability observability.Observability
	server        starter.Server
	operator      service.Operator
}

// New returns Runner instance.
func New(cfg *config.Data) (_ runner.Runner, err error) {
	eg := errgroup.Get()
	operator, err := service.New(
		cfg.Operator.Namespace,
		cfg.Operator.AgentName,
		cfg.Operator.RotatorName,
		cfg.Operator.TargetReadReplicaIDAnnotationsKey,
		cfg.Operator.JobTemplates.Rotate,
		service.WithReadReplicaEnabled(cfg.Operator.ReadReplicaEnabled),
		service.WithReadReplicaLabelKey(cfg.Operator.ReadReplicaLabelKey),
		service.WithRotationJobConcurrency(cfg.Operator.RotationJobConcurrency),
	)
	if err != nil {
		return nil, err
	}

	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithGRPC(func(cfg *iconfig.Server) []server.Option {
			return []server.Option{
				server.WithGRPCOption(
					grpc.ChainUnaryInterceptor(recover.RecoverInterceptor()),
					grpc.ChainStreamInterceptor(recover.RecoverStreamInterceptor()),
				),
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(
			cfg.Observability,
		)
		if err != nil {
			return nil, err
		}
	}

	return &run{
		eg:            eg,
		cfg:           cfg,
		observability: obs,
		server:        srv,
		operator:      operator,
	}, nil
}

// PreStart is a method called before execution of Start, and it invokes the PreStart method of observability.
func (r *run) PreStart(ctx context.Context) error {
	if r.observability != nil {
		return r.observability.PreStart(ctx)
	}
	return nil
}

// Start is a method used to initiate an operation in the run, and it returns a channel for receiving errors
// during the operation and an error representing any initialization errors.
func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 3) //nolint:gomnd
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

// PostStop is a method called after execution of Stop.
func (*run) PostStop(_ context.Context) error {
	return nil
}
