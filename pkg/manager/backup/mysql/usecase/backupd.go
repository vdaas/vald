//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
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
	"github.com/vdaas/vald/pkg/manager/backup/mysql/config"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/handler/grpc"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/handler/rest"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/router"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/service"
)

type Runner runner.Runner

type run struct {
	cfg    *config.Data
	server service.Server
	mySQL  service.MySQL
}

func New(cfg *config.Data) (Runner, error) {
	mySQL, err := service.NewMySQL(cfg.MySQL)
	if err != nil {
		return nil, err
	}
	g := grpc.New(grpc.WithMySQL(mySQL))

	srv, err := service.NewServer(
		service.WithConfig(cfg.Server),
		service.WithREST(
			router.New(
				router.WithHandler(
					rest.New(
						rest.WithBackup(g),
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
		mySQL:  mySQL,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	return r.mySQL.Connect(ctx)
}

func (r *run) Start(ctx context.Context) <-chan error {
	return r.server.ListenAndServe(ctx)
}

func (r *run) PreStop(ctx context.Context) error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}

func (r *run) PostStop(ctx context.Context) error {
	return r.mySQL.Close(ctx)
}
