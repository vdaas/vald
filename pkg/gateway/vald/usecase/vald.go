//
// Copyright (C) 2019 kpango (Yusuke Kato)
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

	"github.com/vdaas/vald/apis/grpc/vald"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/gateway/vald/config"
	handler "github.com/vdaas/vald/pkg/gateway/vald/handler/grpc"
	"github.com/vdaas/vald/pkg/gateway/vald/handler/rest"
	"github.com/vdaas/vald/pkg/gateway/vald/router"
	"github.com/vdaas/vald/pkg/gateway/vald/service"
	"google.golang.org/grpc"
)

type run struct {
	cfg    *config.Data
	server starter.Server
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	service.New(service.WithDiscoverDuration(cfg.Gateway.Discoverer.Host))
	v, err := service.New(cfg.ValdProxy)
	if err != nil {
		return nil, err
	}
	g := handler.New(handler.WithProxy(v))

	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithREST(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithHTTPHandler(
					router.New(
						router.WithHandler(
							rest.New(
							// rest.WithAgent(g),
							),
						),
					)),
			}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(srv *grpc.Server) {
					// vald.RegisterValdServer(srv, g)
					vald.RegisterValdServer(srv, nil)
				}),
				server.WithPreStopFunction(func() error {
					// TODO notify another gateway and scheduler
					return nil
				}),
			}
		}),
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
