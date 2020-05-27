//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"fmt"

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
	cfg    *config.Data
	loader service.Loader
	client grpc.Client
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	run := &run{
		cfg: cfg,
		eg:  errgroup.Get(),
	}

	return run, nil
}

func (r *run) PreStart(ctx context.Context) (err error) {
	r.client = grpc.New(
		grpc.WithAddrs(append([]string{r.cfg.Addr}, r.cfg.Client.Addrs...)...),
		grpc.WithInsecure(r.cfg.Client.DialOption.Insecure),
		grpc.WithErrGroup(r.eg),
	)

	opts := []service.Option{
		service.WithAddr(r.cfg.Addr),
		service.WithDataset(r.cfg.Dataset),
		service.WithClient(r.client),
	}
	switch Atoo(r.cfg.Method) {
	case Insert:
		r.loader, err = service.NewInsert(opts...)
	case Search:
		r.loader, err = service.NewSearch(opts...)
	default:
		return fmt.Errorf("unsupported method")
	}

	return r.loader.Prepare(ctx)
}

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
				err = r.client.Close()
				if err != nil {
					errs = errors.Wrap(errs, err.Error())
				}
			}
			err = ctx.Err()
			if err != nil && err != context.Canceled {
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

func (r *run) PreStop(ctx context.Context) error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	return nil
}

func (r *run) PostStop(ctx context.Context) error {
	return nil
}
