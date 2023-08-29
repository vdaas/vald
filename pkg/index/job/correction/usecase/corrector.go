// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/index/job/correction/config"
	"github.com/vdaas/vald/pkg/index/job/correction/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	observability observability.Observability
	corrector     service.Corrector
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	if cfg.Gateway.IndexReplica == 1 {
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

	corrector, err := service.New(cfg, discoverer)
	if err != nil {
		return nil, err
	}

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(cfg.Observability)
		if err != nil {
			log.Error("failed to initialize observability")
			return nil, err
		}
	}

	return &run{
		eg:            eg,
		cfg:           cfg,
		observability: obs,
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
	// TODO: Set timeout?
	// ctx, cancel := context.WithTimeout(ctx, time.Microsecond*10)
	// defer cancel()

	defer func() {
		log.Info("fiding my pid to kill myself")
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			// using Fatal to avoid this process to be zombie
			log.Fatalf("failed to find my pid to kill %v", err)
			return
		}

		log.Info("sending SIGTERM to myself to stop this job")
		if err := p.Signal(syscall.SIGTERM); err != nil {
			log.Error(err)
		}
	}()

	log.Info("starting index correction job")
	if r.observability != nil {
		_ = r.observability.Start(ctx) // TODO: listen this returned err channel
	}

	start := time.Now()
	_, err := r.corrector.Start(ctx)
	if err != nil {
		log.Errorf("index correction process failed: %v", err)
		return nil, err
	}
	end := time.Since(start)
	log.Infof("correction finished in %v", end)

	// this ech is just a placeholder to return. this is not a daemon but a job.
	// so after returning, this process will be SIGTERMed by myself immediately.
	ech := make(chan error)
	return ech, nil
}

func (*run) PreStop(context.Context) error {
	return nil
}

func (*run) Stop(context.Context) error {
	return nil
}

func (*run) PostStop(ctx context.Context) error {
	return nil
}
