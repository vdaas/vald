//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package storage provides blob storage service
package storage

import (
	"github.com/vdaas/vald/internal/db/storage/blob/cloudstorage"
	"github.com/vdaas/vald/internal/db/storage/blob/cloudstorage/urlopener"
	"github.com/vdaas/vald/internal/db/storage/blob/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/session"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Option func(b *bs) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithCompressAlgorithm("gzip"),
	WithCompressionLevel(-1),
	WithFilenameSuffix(".tar.gz"),
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(b *bs) error {
		b.eg = eg
		return nil
	}
}

func WithType(bst string) Option {
	return func(b *bs) error {
		b.storageType = bst
		return nil
	}
}

func WithBucketName(bn string) Option {
	return func(b *bs) error {
		b.bucketName = bn
		return nil
	}
}

func WithFilename(fn string) Option {
	return func(b *bs) error {
		b.filename = fn
		return nil
	}
}

func WithFilenameSuffix(sf string) Option {
	return func(b *bs) error {
		b.suffix = sf
		return nil
	}
}

func WithS3Opts(opts ...s3.Option) Option {
	return func(b *bs) error {
		if b.s3Opts == nil {
			b.s3Opts = opts
			return nil
		}

		b.s3Opts = append(b.s3Opts, opts...)

		return nil
	}
}

func WithS3SessionOpts(opts ...session.Option) Option {
	return func(b *bs) error {
		if b.s3SessionOpts == nil {
			b.s3SessionOpts = opts
			return nil
		}

		b.s3SessionOpts = append(b.s3SessionOpts, opts...)

		return nil
	}
}

func WithCloudStorageOpts(opts ...cloudstorage.Option) Option {
	return func(b *bs) error {
		if b.cloudStorageOpts == nil {
			b.cloudStorageOpts = opts
		}

		b.cloudStorageOpts = append(b.cloudStorageOpts, opts...)

		return nil
	}
}

func WithCloudStorageURLOpenerOpts(opts ...urlopener.Option) Option {
	return func(b *bs) error {
		if b.cloudStorageURLOpenerOpts == nil {
			b.cloudStorageURLOpenerOpts = opts
		}

		b.cloudStorageURLOpenerOpts = append(b.cloudStorageURLOpenerOpts, opts...)

		return nil
	}
}

func WithCompressAlgorithm(al string) Option {
	return func(b *bs) error {
		b.compressAlgorithm = al
		return nil
	}
}

func WithCompressionLevel(level int) Option {
	return func(b *bs) error {
		b.compressionLevel = level
		return nil
	}
}
