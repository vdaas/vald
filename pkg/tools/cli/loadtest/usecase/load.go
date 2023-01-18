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

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/service"
)

type run struct {
	eg     errgroup.Group
	loader service.Loader
	client grpc.Client
}

// New returns Runner instance.
func New(cfg *config.Data) (r runner.Runner, err error) {
	run := &run{
		eg: errgroup.Get(),
	}

	cOpts, err := cfg.Client.Opts()
	if err != nil {
		return nil, err
	}
	clientOpts := append(
		cOpts,
		grpc.WithAddrs(cfg.Addr),
		grpc.WithErrGroup(run.eg),
	)
	run.client = grpc.New(clientOpts...)

	run.loader, err = service.NewLoader(
		service.WithOperation(cfg.Operation),
		service.WithAddr(cfg.Addr),
		service.WithBatchSize(cfg.BatchSize),
		service.WithDataset(cfg.Dataset),
		service.WithClient(run.client),
		service.WithConcurrency(cfg.Concurrency),
		service.WithProgressDuration(cfg.ProgressDuration),
	)
	if err != nil {
		return nil, err
	}

	return run, nil
}

// PreStart initializes load tester and returns error if occurred.
func (r *run) PreStart(ctx context.Context) (err error) {
	return r.loader.Prepare(ctx)
}

// Start runs load test and returns error if occurred.
func (r *run) Start(ctx context.Context) (<-chan error, error) {
	rech, err := r.client.StartConnectionMonitor(ctx)
	if err != nil {
		return nil, err
	}

	lech := r.loader.Do(ctx)

	ech := make(chan error, 1000) // TODO: fix magic number
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		finalize := func() (err error) {
			var errs error
			if r.client != nil {
				err = r.client.Close(ctx)
				if err != nil {
					errs = errors.Wrap(errs, err.Error())
				}
			}
			err = ctx.Err()
			if err != nil && !errors.Is(err, context.Canceled) {
				errs = errors.Wrap(errs, err.Error())
			}
			return errs
		}
		for {
			select {
			case <-ctx.Done():
				return finalize()
			case err = <-rech:
			case err = <-lech:
			}
			if err != nil {
				log.Error(err)
				select {
				case <-ctx.Done():
					return finalize()
				case ech <- err:
				}
			}
		}
	}))
	return ech, nil
}

// PreStop does nothing.
func (*run) PreStop(ctx context.Context) error {
	return nil
}

// Stop does nothing.
func (r *run) Stop(ctx context.Context) error {
	return nil
}

// PostStop does nothing.
func (*run) PostStop(ctx context.Context) error {
	return nil
}
