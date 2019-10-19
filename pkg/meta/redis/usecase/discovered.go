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

	"github.com/vdaas/vald/apis/grpc/discoverer"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/meta/redis/config"
	handler "github.com/vdaas/vald/pkg/meta/redis/handler/grpc"
	"github.com/vdaas/vald/pkg/meta/redis/handler/rest"
	"github.com/vdaas/vald/pkg/meta/redis/router"
	"github.com/vdaas/vald/pkg/meta/redis/service"
	"google.golang.org/grpc"
)

type run struct {
	cfg    *config.Data
	dsc    service.Discoverer
	server starter.Server
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	dsc, err := service.New()
	if err != nil {
		return nil, err
	}
	g := handler.New(handler.WithDiscoverer(dsc))

	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithREST(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithHTTPHandler(
					router.New(
						router.WithTimeout(sc.HTTP.HandlerTimeout),
						router.WithErrGroup(errgroup.Get()),
						router.WithHandler(
							rest.New(
								rest.WithDiscoverer(g),
							),
						),
					)),
			}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithGRPCRegistFunc(func(srv *grpc.Server) {
					discoverer.RegisterDiscovererServer(srv, g)
				}),
				server.WithPreStartFunc(func() error {
					// TODO check unbackupped upstream
					return nil
				}),
				server.WithPreStopFunction(func() error {
					// TODO backup all index data here
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
		dsc:    dsc,
		server: srv,
	}, nil
}

func (r *run) PreStart() error {
	log.Info("daemon start")
	return nil
}

func (r *run) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 2)
	go func() {
		log.Info("daemon start")
		defer close(ech)
		dech := r.dsc.Start(ctx)
		sech := r.server.ListenAndServe(ctx)
		for {
			select {
			case ech <- <-dech:
			case ech <- <-sech:
			}
		}
	}()
	return ech
}

func (r *run) PreStop() error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}
