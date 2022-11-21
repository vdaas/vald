//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package storage provides blob storage service
package storage

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/compress"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/cloudstorage"
	"github.com/vdaas/vald/internal/db/storage/blob/cloudstorage/urlopener"
	"github.com/vdaas/vald/internal/db/storage/blob/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/session"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
)

type Storage interface {
	Start(ctx context.Context) (<-chan error, error)
	Stop(ctx context.Context) error
	Reader(ctx context.Context) (io.ReadCloser, error)
	Writer(ctx context.Context) (io.WriteCloser, error)
	StorageInfo() *StorageInfo
}

type bs struct {
	eg          errgroup.Group
	storageType string
	bucketName  string
	filename    string
	suffix      string

	s3Opts        []s3.Option
	s3SessionOpts []session.Option

	cloudStorageOpts          []cloudstorage.Option
	cloudStorageURLOpenerOpts []urlopener.Option

	compressAlgorithm string
	compressionLevel  int

	bucket     blob.Bucket
	compressor compress.Compressor
}

func New(opts ...Option) (Storage, error) {
	b := new(bs)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(b); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	err := b.initCompressor()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *bs) initCompressor() (err error) {
	// Without compress
	if b.compressAlgorithm == "" {
		return nil
	}

	switch config.AToCompressAlgorithm(b.compressAlgorithm) {
	case config.GOB:
		b.compressor, err = compress.NewGob()
	case config.GZIP:
		b.compressor, err = compress.NewGzip(
			compress.WithGzipCompressionLevel(b.compressionLevel),
		)
	case config.LZ4:
		b.compressor, err = compress.NewLZ4(
			compress.WithLZ4CompressionLevel(b.compressionLevel),
		)
	case config.ZSTD:
		b.compressor, err = compress.NewZstd(
			compress.WithZstdCompressionLevel(b.compressionLevel),
		)
	default:
		return errors.ErrCompressorNameNotFound(b.compressAlgorithm)
	}

	return err
}

func (b *bs) initBucket(ctx context.Context) (err error) {
	switch config.AtoBST(b.storageType) {
	case config.S3:
		s, err := session.New(b.s3SessionOpts...).Session()
		if err != nil {
			return err
		}

		b.bucket, err = s3.New(
			append(
				b.s3Opts,
				s3.WithErrGroup(b.eg),
				s3.WithSession(s),
				s3.WithBucket(b.bucketName),
			)...,
		)
		if err != nil {
			return err
		}
	case config.CloudStorage:
		uoi, err := urlopener.New(b.cloudStorageURLOpenerOpts...)
		if err != nil {
			return err
		}

		uo, err := uoi.URLOpener(ctx)
		if err != nil {
			return err
		}

		b.bucket, err = cloudstorage.New(
			append(
				b.cloudStorageOpts,
				cloudstorage.WithURLOpener(uo),
			)...,
		)
		if err != nil {
			return err
		}
	default:
		return errors.ErrInvalidStorageType
	}

	return nil
}

func (b *bs) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 1)

	err := b.initBucket(ctx)
	if err != nil {
		return nil, err
	}

	err = b.bucket.Open(ctx)
	if err != nil {
		return nil, err
	}

	return ech, nil
}

func (b *bs) Stop(_ context.Context) error {
	if b.bucket != nil {
		return b.bucket.Close()
	}
	return nil
}

func (b *bs) Reader(ctx context.Context) (r io.ReadCloser, err error) {
	r, err = b.bucket.Reader(ctx, b.filename+b.suffix)
	if err != nil {
		return nil, err
	}

	if b.compressor != nil {
		r, err = b.compressor.Reader(r)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (b *bs) Writer(ctx context.Context) (w io.WriteCloser, err error) {
	w, err = b.bucket.Writer(ctx, b.filename+b.suffix)
	if err != nil {
		return nil, err
	}

	if b.compressor != nil {
		w, err = b.compressor.Writer(w)
		if err != nil {
			return nil, err
		}
	}

	return w, nil
}

func (b *bs) StorageInfo() *StorageInfo {
	return &StorageInfo{
		Type:       config.AtoBST(b.storageType).String(),
		BucketName: b.bucketName,
		Filename:   b.filename + b.suffix,
	}
}
