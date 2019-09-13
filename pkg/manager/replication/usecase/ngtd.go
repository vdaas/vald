//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/agent/ngt/config"
	"github.com/vdaas/vald/pkg/agent/ngt/handler/grpc"
	"github.com/vdaas/vald/pkg/agent/ngt/handler/rest"
	"github.com/vdaas/vald/pkg/agent/ngt/router"
	"github.com/vdaas/vald/pkg/agent/ngt/service"
)

type Runner runner.Runner

type run struct {
	cfg    *config.Data
	server service.Server
}

func New(cfg *config.Data) (Runner, error) {
	ngt, err := service.NewNGT(cfg.NGT)
	if err != nil {
		return nil, err
	}
	g := grpc.New(grpc.WithNGT(ngt))

	srv, err := service.NewServer(
		service.WithConfig(cfg.Server),
		service.WithREST(
			router.New(
				router.WithHandler(
					rest.New(
						rest.WithAgent(g),
					),
				),
			),
		),
		service.WithGRPC(g),
		// TODO add GraphQL handler
	)

	if err != nil {
		return nil, err
	}

	return &run{
		cfg:    cfg,
		server: srv,
	}, nil
}

func (r *run) PreStart() error {
	return nil
}

func (r *run) Start(ctx context.Context) <-chan error {
	return r.server.ListenAndServe(ctx)
}

func (r *run) PreStop() error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}
