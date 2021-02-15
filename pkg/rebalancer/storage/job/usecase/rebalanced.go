//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/storage/blob/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/session"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/recover"
	"github.com/vdaas/vald/internal/net/grpc/metric"
	"github.com/vdaas/vald/internal/observability"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/servers/starter"
	"github.com/vdaas/vald/pkg/rebalancer/storage/job/config"
	handler "github.com/vdaas/vald/pkg/rebalancer/storage/job/handler/grpc"
	"github.com/vdaas/vald/pkg/rebalancer/storage/job/handler/rest"
	"github.com/vdaas/vald/pkg/rebalancer/storage/job/router"
	"github.com/vdaas/vald/pkg/rebalancer/storage/job/service/job"
	"github.com/vdaas/vald/pkg/rebalancer/storage/job/service/storage"
)

type run struct {
	eg            errgroup.Group
	cfg           *config.Data
	rb            job.Rebalancer
	h             handler.Rebalancer
	server        starter.Server
	observability observability.Observability
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	eg := errgroup.Get()
	st, err := storage.New(
		storage.WithErrGroup(eg),
		storage.WithType(cfg.Rebalancer.BlobStorage.StorageType),
		storage.WithBucketName(cfg.Rebalancer.BlobStorage.Bucket),
		// TODO: define filename
		// storage.WithFilename(cfg.Rebalancer.Filename),
		storage.WithFilenameSuffix(cfg.Rebalancer.FilenameSuffix),
		storage.WithS3SessionOpts(
			session.WithEndpoint(cfg.Rebalancer.BlobStorage.S3.Endpoint),
			session.WithRegion(cfg.Rebalancer.BlobStorage.S3.Region),
			session.WithAccessKey(cfg.Rebalancer.BlobStorage.S3.AccessKey),
			session.WithSecretAccessKey(cfg.Rebalancer.BlobStorage.S3.SecretAccessKey),
			session.WithToken(cfg.Rebalancer.BlobStorage.S3.Token),
			session.WithMaxRetries(cfg.Rebalancer.BlobStorage.S3.MaxRetries),
			session.WithForcePathStyle(cfg.Rebalancer.BlobStorage.S3.ForcePathStyle),
			session.WithUseAccelerate(cfg.Rebalancer.BlobStorage.S3.UseAccelerate),
			session.WithUseARNRegion(cfg.Rebalancer.BlobStorage.S3.UseARNRegion),
			session.WithUseDualStack(cfg.Rebalancer.BlobStorage.S3.UseDualStack),
			session.WithEnableSSL(cfg.Rebalancer.BlobStorage.S3.EnableSSL),
			session.WithEnableParamValidation(cfg.Rebalancer.BlobStorage.S3.EnableParamValidation),
			session.WithEnable100Continue(cfg.Rebalancer.BlobStorage.S3.Enable100Continue),
			session.WithEnableContentMD5Validation(cfg.Rebalancer.BlobStorage.S3.EnableContentMD5Validation),
			session.WithEnableEndpointDiscovery(cfg.Rebalancer.BlobStorage.S3.EnableEndpointDiscovery),
			session.WithEnableEndpointHostPrefix(cfg.Rebalancer.BlobStorage.S3.EnableEndpointHostPrefix),
			// TODO: set client
			// session.WithHTTPClient(client),
		),
		storage.WithS3Opts(
			s3.WithMaxPartSize(cfg.Rebalancer.BlobStorage.S3.MaxPartSize),
			s3.WithMaxChunkSize(cfg.Rebalancer.BlobStorage.S3.MaxChunkSize),
			// TODO: backoff opts
			// s3.WithReaderBackoff(cfg.AgentSidecar.RestoreBackoffEnabled),
			// s3.WithReaderBackoffOpts(cfg.AgentSidecar.RestoreBackoff.Opts()...),
		),
		storage.WithCompressAlgorithm(cfg.Rebalancer.Compress.CompressAlgorithm),
		storage.WithCompressionLevel(cfg.Rebalancer.Compress.CompressionLevel),
	)
	if err != nil {
		return nil, err
	}
	rb, err := job.New(
		job.WithStorage(st),
	)
	if err != nil {
		return nil, err
	}
	h, err := handler.New(
		handler.WithDiscoverer(rb),
	)
	if err != nil {
		return nil, err
	}

	grpcServerOptions := []server.Option{
		server.WithGRPCRegistFunc(func(srv *grpc.Server) {
			// TODO register grpc server handler here
		}),
		server.WithGRPCOption(
			grpc.ChainUnaryInterceptor(recover.RecoverInterceptor()),
			grpc.ChainStreamInterceptor(recover.RecoverStreamInterceptor()),
		),
		server.WithPreStartFunc(func() error {
			// TODO check unbackupped upstream
			return nil
		}),
		server.WithPreStopFunction(func() error {
			// TODO backup all index data here
			return nil
		}),
	}

	var obs observability.Observability
	if cfg.Observability.Enabled {
		obs, err = observability.NewWithConfig(cfg.Observability)
		if err != nil {
			return nil, err
		}
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
						router.WithTimeout(sc.HTTP.HandlerTimeout),
						router.WithErrGroup(eg),
						router.WithHandler(
							rest.New(
							// TODO pass grpc handler to REST option
							),
						),
					)),
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
		rb:            rb,
		h:             h,
		server:        srv,
		observability: obs,
	}, nil
}

func (r *run) PreStart(ctx context.Context) error {
	if r.observability != nil {
		return r.observability.PreStart(ctx)
	}
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 3)
	var oech, dech, sech <-chan error
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		if r.observability != nil {
			oech = r.observability.Start(ctx)
		}
		dech, err = r.rb.Start(ctx)

		if err != nil {
			ech <- err
			return err
		}

		r.h.Start(ctx)

		sech = r.server.ListenAndServe(ctx)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-oech:
			case err = <-dech:
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
