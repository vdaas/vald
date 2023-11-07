// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"time"

	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/recover"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/observability/metrics/index/job/correction"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/index/job/correction/config"
	"github.com/vdaas/vald/pkg/index/job/correction/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	observability observability.Observability
	server        starter.Server
	corrector     service.Corrector
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	if cfg.Corrector.IndexReplica == 1 {
		return nil, errors.ErrIndexReplicaOne
	}

	eg := errgroup.Get()

	cOpts, err := cfg.Corrector.Discoverer.Client.Opts()
	if err != nil {
		return nil, err
	}
	// skipcq: CRT-D0001
	dopts := append(
		cOpts,
		grpc.WithErrGroup(eg))

	acOpts, err := cfg.Corrector.Discoverer.AgentClientOptions.Opts()
	if err != nil {
		return nil, err
	}
	// skipcq: CRT-D0001
	aopts := append(
		acOpts,
		grpc.WithErrGroup(eg))

	// Construct discoverer
	discoverer, err := discoverer.New(
		discoverer.WithAutoConnect(true),
		discoverer.WithName(cfg.Corrector.AgentName),
		discoverer.WithNamespace(cfg.Corrector.AgentNamespace),
		discoverer.WithPort(cfg.Corrector.AgentPort),
		discoverer.WithServiceDNSARecord(cfg.Corrector.AgentDNS),
		discoverer.WithDiscovererClient(grpc.New(dopts...)),
		discoverer.WithDiscoverDuration(cfg.Corrector.Discoverer.Duration),
		discoverer.WithOptions(aopts...),
		discoverer.WithNodeName(cfg.Corrector.NodeName),
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

	grpcServerOptions := []server.Option{
		server.WithGRPCOption(
			grpc.ChainUnaryInterceptor(recover.RecoverInterceptor()),
			grpc.ChainStreamInterceptor(recover.RecoverStreamInterceptor()),
		),
	}

	// For health check and metrics
	srv, err := starter.New(starter.WithConfig(cfg.Server),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return grpcServerOptions
		}),
	)
	if err != nil {
		return nil, err
	}

	corrector, err := service.New(cfg, discoverer)
	if err != nil {
		return nil, err
	}

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(
			cfg.Observability,
			correction.New(corrector),
		)
		if err != nil {
			log.Error("failed to initialize observability")
			return nil, err
		}
	}

	return &run{
		eg:            eg,
		cfg:           cfg,
		observability: obs,
		server:        srv,
		corrector:     corrector,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	if r.observability != nil {
		return r.observability.PreStart(ctx)
	}
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	log.Info("starting servers")
	ech := make(chan error, 3) //nolint:gomnd
	var oech <-chan error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
	}
	sech := r.server.ListenAndServe(ctx)
	nech, err := r.corrector.PreStart(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
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

	// main groutine to run the job
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer func() {
			log.Info("fiding my pid to kill myself")
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

		start := time.Now()
		err = r.corrector.Start(ctx)
		if err != nil {
			log.Errorf("index correction process failed: %v", err)
			return err
		}
		end := time.Since(start)
		log.Infof("correction finished in %v", end)
		return nil
	}))

	return ech, nil
}

func (r *run) PreStop(ctx context.Context) error {
	r.corrector.PreStop(ctx)
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	if r.observability != nil {
		r.observability.Stop(ctx)
	}
	if r.server != nil {
		r.server.Shutdown(ctx)
	}
	return nil
}

func (*run) PostStop(_ context.Context) error {
	return nil
}
