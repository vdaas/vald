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

package initcontainer

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/agent/sidecar"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/storage/blob/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/session"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/metric"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/agent/sidecar/config"
	handler "github.com/vdaas/vald/pkg/agent/sidecar/handler/grpc"
	"github.com/vdaas/vald/pkg/agent/sidecar/handler/rest"
	"github.com/vdaas/vald/pkg/agent/sidecar/router"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/restorer"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	server        starter.Server
	observability observability.Observability
	rs            restorer.Restorer
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	log.Info("Initialized in initcontainer mode")

	eg := errgroup.Get()

	var (
		rs restorer.Restorer
		bs storage.Storage
	)

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(cfg.Observability)
		if err != nil {
			return nil, err
		}
		// TODO observe something
		_ = obs
	}
	bs, err = storage.New(
		storage.WithErrGroup(eg),
		storage.WithType(cfg.AgentSidecar.BlobStorage.StorageType),
		storage.WithBucketName(cfg.AgentSidecar.BlobStorage.Bucket),
		storage.WithFilename(cfg.AgentSidecar.Filename),
		storage.WithFilenameSuffix(cfg.AgentSidecar.FilenameSuffix),
		storage.WithS3SessionOpts(
			session.WithEndpoint(cfg.AgentSidecar.BlobStorage.S3.Endpoint),
			session.WithRegion(cfg.AgentSidecar.BlobStorage.S3.Region),
			session.WithAccessKey(cfg.AgentSidecar.BlobStorage.S3.AccessKey),
			session.WithSecretAccessKey(cfg.AgentSidecar.BlobStorage.S3.SecretAccessKey),
			session.WithToken(cfg.AgentSidecar.BlobStorage.S3.Token),
			session.WithEnableSSL(cfg.AgentSidecar.BlobStorage.S3.EnableSSL),
		),
		storage.WithS3Opts(
			s3.WithMaxPartSize(cfg.AgentSidecar.BlobStorage.S3.MaxPartSize),
		),
		storage.WithCompressAlgorithm(cfg.AgentSidecar.Compress.CompressAlgorithm),
		storage.WithCompressionLevel(cfg.AgentSidecar.Compress.CompressionLevel),
	)
	if err != nil {
		return nil, err
	}

	rs, err = restorer.New(
		restorer.WithErrGroup(eg),
		restorer.WithDir(cfg.AgentSidecar.WatchDir),
		restorer.WithBlobStorage(bs),
	)
	if err != nil {
		return nil, err
	}

	g := handler.New()

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			sidecar.RegisterSidecarServer(srv, g)
		}),
		server.WithGRPCOption(
			grpc.ChainUnaryInterceptor(grpc.RecoverInterceptor()),
			grpc.ChainStreamInterceptor(grpc.RecoverStreamInterceptor()),
		),
		server.WithPreStopFunction(func() error {
			// TODO notify another gateway and scheduler
			return nil
		}),
	}

	if cfg.Observability.Enabled {
		grpcServerOptions = append(
			grpcServerOptions,
			server.WithGRPCOption(
				grpc.StatsHandler(metric.NewServerHandler()),
			),
		)
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
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	if r.observability != nil {
		r.observability.Stop(ctx)
	}
	return r.server.Shutdown(ctx)
}

func (r *run) PostStop(ctx context.Context) error {
	return nil
}
