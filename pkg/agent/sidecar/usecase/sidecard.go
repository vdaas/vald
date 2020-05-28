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

	"github.com/vdaas/vald/apis/grpc/agent/sidecar"
	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
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
	"github.com/vdaas/vald/pkg/agent/sidecar/service"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	server        starter.Server
	observability observability.Observability
	so            service.StorageObserver
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()

	var (
		so service.StorageObserver
		bs service.BlobStorage
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
	bs, err = service.NewBlobStorage(
		service.WithBlobStorageType(cfg.AgentSidecar.BlobStorage.StorageType),
		service.WithBlobStorageBucketName(cfg.AgentSidecar.BlobStorage.Bucket),
		service.WithBlobStorageFilename(cfg.AgentSidecar.Filename),
		service.WithBlobStorageFilenameSuffix(cfg.AgentSidecar.FilenameSuffix),
		service.WithBlobStorageEndpoint(cfg.AgentSidecar.BlobStorage.S3.Endpoint),
		service.WithBlobStorageRegion(cfg.AgentSidecar.BlobStorage.S3.Region),
		service.WithBlobStorageAccessKey(cfg.AgentSidecar.BlobStorage.S3.AccessKey),
		service.WithBlobStorageSecretAccessKey(cfg.AgentSidecar.BlobStorage.S3.SecretAccessKey),
		service.WithBlobStorageToken(cfg.AgentSidecar.BlobStorage.S3.Token),
		service.WithBlobStorageMultipartUpload(cfg.AgentSidecar.BlobStorage.S3.MultipartUpload),
		service.WithBlobStorageCompressAlgorithm(cfg.AgentSidecar.Compress.CompressAlgorithm),
		service.WithBlobStorageCompressionLevel(cfg.AgentSidecar.Compress.CompressionLevel),
	)
	if err != nil {
		return nil, err
	}

	so, err = service.New(
		service.WithErrGroup(eg),
		service.WithBackupDuration(cfg.AgentSidecar.AutoBackupDuration),
		service.WithBackupDurationLimit(cfg.AgentSidecar.AutoBackupDurationLimit),
		service.WithDir(cfg.AgentSidecar.WatchDir),
		service.WithBlobStorage(bs),
	)
	if err != nil {
		return nil, err
	}

	g := handler.New(handler.WithStorageObserver(so))

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			sidecar.RegisterSidecarServer(srv, g)
		}),
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
		so:            so,
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
	var soech, sech, oech <-chan error
	var err error
	if r.observability != nil {
		oech = r.observability.Start(ctx)
	}
	if r.so != nil {
		soech, err = r.so.Start(ctx)
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
			case err = <-oech:
			case err = <-soech:
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
