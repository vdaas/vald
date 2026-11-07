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

// Package blobstorage builds the agent-sidecar blob storage from the sidecar
// config. It is shared by the initcontainer and sidecar usecases so that the
// storage wiring (and every new BlobStorage config field) lives in one place.
package blobstorage

import (
	"github.com/vdaas/vald/internal/db/storage/blob/cloudstorage"
	"github.com/vdaas/vald/internal/db/storage/blob/cloudstorage/urlopener"
	"github.com/vdaas/vald/internal/db/storage/blob/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/session"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/http/client"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/sidecar/config"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

// New builds the HTTP client and the blob storage from cfg. clientOpts are
// appended after the common HTTP client options; the sidecar passes its
// transport backoff there, the initcontainer passes none.
func New(
	cfg *config.Data, eg errgroup.Group, clientOpts ...client.Option,
) (storage.Storage, error) {
	netOpts, err := cfg.AgentSidecar.Client.Net.Opts()
	if err != nil {
		return nil, err
	}

	dialer, err := net.NewDialer(netOpts...)
	if err != nil {
		return nil, err
	}

	client, err := client.New(append([]client.Option{
		client.WithDialContext(dialer.DialContext),
		client.WithTLSHandshakeTimeout(cfg.AgentSidecar.Client.Transport.RoundTripper.TLSHandshakeTimeout),
		client.WithMaxIdleConns(cfg.AgentSidecar.Client.Transport.RoundTripper.MaxIdleConns),
		client.WithMaxIdleConnsPerHost(cfg.AgentSidecar.Client.Transport.RoundTripper.MaxIdleConnsPerHost),
		client.WithMaxConnsPerHost(cfg.AgentSidecar.Client.Transport.RoundTripper.MaxConnsPerHost),
		client.WithIdleConnTimeout(cfg.AgentSidecar.Client.Transport.RoundTripper.IdleConnTimeout),
		client.WithResponseHeaderTimeout(cfg.AgentSidecar.Client.Transport.RoundTripper.ResponseHeaderTimeout),
		client.WithExpectContinueTimeout(cfg.AgentSidecar.Client.Transport.RoundTripper.ExpectContinueTimeout),
		client.WithMaxResponseHeaderBytes(cfg.AgentSidecar.Client.Transport.RoundTripper.MaxResponseHeaderSize),
		client.WithWriteBufferSize(cfg.AgentSidecar.Client.Transport.RoundTripper.WriteBufferSize),
		client.WithReadBufferSize(cfg.AgentSidecar.Client.Transport.RoundTripper.ReadBufferSize),
		client.WithForceAttemptHTTP2(cfg.AgentSidecar.Client.Transport.RoundTripper.ForceAttemptHTTP2),
	}, clientOpts...)...)
	if err != nil {
		return nil, err
	}

	return storage.New(
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
			session.WithMaxRetries(cfg.AgentSidecar.BlobStorage.S3.MaxRetries),
			session.WithForcePathStyle(cfg.AgentSidecar.BlobStorage.S3.ForcePathStyle),
			session.WithUseAccelerate(cfg.AgentSidecar.BlobStorage.S3.UseAccelerate),
			session.WithUseARNRegion(cfg.AgentSidecar.BlobStorage.S3.UseARNRegion),
			session.WithUseDualStack(cfg.AgentSidecar.BlobStorage.S3.UseDualStack),
			session.WithEnableSSL(cfg.AgentSidecar.BlobStorage.S3.EnableSSL),
			session.WithEnableParamValidation(cfg.AgentSidecar.BlobStorage.S3.EnableParamValidation),
			session.WithEnable100Continue(cfg.AgentSidecar.BlobStorage.S3.Enable100Continue),
			session.WithEnableContentMD5Validation(cfg.AgentSidecar.BlobStorage.S3.EnableContentMD5Validation),
			session.WithEnableEndpointDiscovery(cfg.AgentSidecar.BlobStorage.S3.EnableEndpointDiscovery),
			session.WithEnableEndpointHostPrefix(cfg.AgentSidecar.BlobStorage.S3.EnableEndpointHostPrefix),
			session.WithHTTPClient(client),
		),
		storage.WithS3Opts(
			s3.WithMaxPartSize(cfg.AgentSidecar.BlobStorage.S3.MaxPartSize),
			s3.WithMaxChunkSize(cfg.AgentSidecar.BlobStorage.S3.MaxChunkSize),
			s3.WithReaderBackoff(cfg.AgentSidecar.RestoreBackoffEnabled),
			s3.WithReaderBackoffOpts(cfg.AgentSidecar.RestoreBackoff.Opts()...),
		),
		storage.WithCloudStorageURLOpenerOpts(
			urlopener.WithCredentialsFile(cfg.AgentSidecar.BlobStorage.CloudStorage.Client.CredentialsFilePath),
			urlopener.WithCredentialsJSON(cfg.AgentSidecar.BlobStorage.CloudStorage.Client.CredentialsJSON),
			urlopener.WithHTTPClient(client),
		),
		storage.WithCloudStorageOpts(
			cloudstorage.WithURL(cfg.AgentSidecar.BlobStorage.CloudStorage.URL),
			cloudstorage.WithWriteBufferSize(cfg.AgentSidecar.BlobStorage.CloudStorage.WriteBufferSize),
			cloudstorage.WithWriteCacheControl(cfg.AgentSidecar.BlobStorage.CloudStorage.WriteCacheControl),
			cloudstorage.WithWriteContentDisposition(cfg.AgentSidecar.BlobStorage.CloudStorage.WriteContentDisposition),
			cloudstorage.WithWriteContentEncoding(cfg.AgentSidecar.BlobStorage.CloudStorage.WriteContentEncoding),
			cloudstorage.WithWriteContentLanguage(cfg.AgentSidecar.BlobStorage.CloudStorage.WriteContentLanguage),
			cloudstorage.WithWriteContentType(cfg.AgentSidecar.BlobStorage.CloudStorage.WriteContentType),
		),
		storage.WithCompressAlgorithm(cfg.AgentSidecar.Compress.CompressAlgorithm),
		storage.WithCompressionLevel(cfg.AgentSidecar.Compress.CompressionLevel),
	)
}
