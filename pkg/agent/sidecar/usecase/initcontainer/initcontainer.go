// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package initcontainer

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/agent/sidecar"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/sidecar/config"
	handler "github.com/vdaas/vald/pkg/agent/sidecar/handler/grpc"
	"github.com/vdaas/vald/pkg/agent/sidecar/handler/rest"
	"github.com/vdaas/vald/pkg/agent/sidecar/router"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/restorer"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
	"github.com/vdaas/vald/pkg/agent/sidecar/usecase/internal/blobstorage"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	server        starter.Server
	observability observability.Observability
	rs            restorer.Restorer
}

func New(cfg *config.Data) (r runner.Interface, err error) {
	log.Info("Initialized in initcontainer mode")

	eg := errgroup.Get()

	var (
		rs restorer.Restorer
		bs storage.Storage
	)

	bs, err = blobstorage.New(cfg, eg)
	if err != nil {
		return nil, err
	}

	rs, err = restorer.New(
		restorer.WithErrGroup(eg),
		restorer.WithDir(cfg.AgentSidecar.WatchDir),
		restorer.WithBlobStorage(bs),
		restorer.WithBackoff(cfg.AgentSidecar.RestoreBackoffEnabled),
		restorer.WithBackoffOpts(cfg.AgentSidecar.RestoreBackoff.Opts()...),
	)
	if err != nil {
		return nil, err
	}

	g := handler.New()

	grpcServerOptions := []server.Option{
		server.WithGRPCRegisterar(func(srv *grpc.Server) {
			sidecar.RegisterSidecarServer(srv, g)
		}),
		server.WithPreStopFunction(func() error {
			// TODO notify another gateway and scheduler
			return nil
		}),
	}

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(cfg.Observability)
		if err != nil {
			return nil, err
		}
	}

	srv, err := starter.New(
		starter.WithConfig(cfg.Server),
		starter.WithREST(func(sc *iconf.Server) []server.Option {
			return []server.Option{
				server.WithHTTPHandler(
					router.New(
						router.WithHandler(
							rest.New(
								rest.WithSidecar(g),
							),
						),
					),
				),
			}
		}),
		starter.WithGRPC(func(sc *iconf.Server) []server.Option {
			return grpcServerOptions
		}),
		// TODO add GraphQL handler
	)
	if err != nil {
		return nil, err
	}

	return &run{
		eg:            eg,
		cfg:           cfg,
		server:        srv,
		observability: obs,
		rs:            rs,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	if r.observability != nil {
		return r.observability.PreStart(ctx)
	}
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 5)
	var rsech, sech, oech <-chan error
	var err error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
	}
	if r.rs != nil {
		rsech, err = r.rs.Start(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}
	sech = r.server.ListenAndServe(ctx)
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-rsech:
				if errors.Is(err, io.EOF) {
					err = nil
				}
			case err = <-oech:
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
	return ech, nil
}

func (r *run) PreStop(ctx context.Context) error {
	if r.rs != nil {
		return r.rs.PreStop(ctx)
	}
	return nil
}

func (r *run) Stop(ctx context.Context) (err error) {
	// Shut the server down first so spans and metrics of the requests it drains
	// are still exported before observability is stopped.
	err = r.server.Shutdown(ctx)
	if r.observability != nil {
		err = errors.Join(err, r.observability.Stop(ctx))
	}
	return err
}

func (*run) PostStop(context.Context) error {
	return nil
}
