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
	"strings"

	"github.com/vdaas/vald/internal/client/gateway/vald/grpc"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/service"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/service/insert"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/service/search"
)

type run struct {
	eg  errgroup.Group
	cfg *config.Data
	l   service.Load
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	run := new(run)
	run.cfg = cfg
	run.eg = errgroup.Get()

	ctx := context.Background()
	c, err := grpc.New(ctx, grpc.WithAddr(cfg.Address)) // TODO setup vald grpc client
	if err != nil {
		return nil, fmt.Errorf("grpc connection error")
	}
	switch strings.ToLower(cfg.Method) {
	case "insert":
		run.l, err = insert.New(insert.WithDataset(cfg.Dataset), insert.WithWriter(c))
	case "search":
		run.l, err = search.New(search.WithDataset(cfg.Dataset), search.WithReader(c))
	default:
		return nil, fmt.Errorf("unsupported method")
	}

	return run, nil
}

func (r *run) PreStart(ctx context.Context) (err error) {
	return r.l.Prepare(ctx)
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	return r.l.Do(ctx), nil
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
