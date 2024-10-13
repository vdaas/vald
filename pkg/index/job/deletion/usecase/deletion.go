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
	"os"
	"syscall"

	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	iconfig "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/recover"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/index/job/deletion/config"
	"github.com/vdaas/vald/pkg/index/job/deletion/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	observability observability.Observability
	server        starter.Server
	indexer       service.Deleter
}

// New returns Runner instance.
func New(cfg *config.Data) (_ runner.Runner, err error) {
	eg := errgroup.Get()

	dOpts, err := cfg.Deletion.Discoverer.Client.Opts()
	if err != nil {
		return nil, err
	}
	// skipcq: CRT-D0001
	dOpts = append(dOpts, grpc.WithErrGroup(eg))

	acOpts, err := cfg.Deletion.Discoverer.AgentClientOptions.Opts()
	if err != nil {
		return nil, err
	}
	// skipcq: CRT-D0001
	acOpts = append(acOpts, grpc.WithErrGroup(eg))

	discoverer, err := discoverer.New(
		discoverer.WithAutoConnect(true),
		discoverer.WithName(cfg.Deletion.AgentName),
		discoverer.WithNamespace(cfg.Deletion.AgentNamespace),
		discoverer.WithPort(cfg.Deletion.AgentPort),
		discoverer.WithServiceDNSARecord(cfg.Deletion.AgentDNS),
		discoverer.WithDiscovererClient(grpc.New(dOpts...)),
		discoverer.WithDiscoverDuration(cfg.Deletion.Discoverer.Duration),
		discoverer.WithOptions(acOpts...),
		discoverer.WithNodeName(cfg.Deletion.NodeName),
		discoverer.WithOnDiscoverFunc(func(ctx context.Context, c discoverer.Client, addrs []string) error {
			last := len(addrs) - 1
			for i := 0; i < len(addrs)/2; i++ {
				addrs[i], addrs[last-i] = addrs[last-i], addrs[i]
			}
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	indexer, err := service.New(
		service.WithDiscoverer(discoverer),
		service.WithIndexingConcurrency(cfg.Deletion.Concurrency),
		service.WithTargetAddrs(cfg.Deletion.TargetAddrs...),
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
		indexer:       indexer,
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
	ech := make(chan error, 4)
	var sech, oech <-chan error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
	}
	sech = r.server.ListenAndServe(ctx)
	cech, err := r.indexer.StartClient(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer func() {
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				// using Fatal to avoid this process to be zombie
				// skipcq: RVV-A0003
				log.Fatalf("failed to find my pid to kill %v", err)
				return
			}
			log.Info("sending SIGTERM to myself to stop this job")
			if err := p.Signal(syscall.SIGTERM); err != nil {
				log.Error(err)
			}
		}()
		return r.indexer.Start(ctx)
	}))

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-oech:
			case err = <-sech:
			case err = <-cech:
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
