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

// Package storage provides blob storage service
package storage

import "github.com/vdaas/vald/internal/errgroup"

type Option func(b *bs) error

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
		WithCompressAlgorithm("gzip"),
		WithCompressionLevel(-1),
		WithFilenameSuffix(".tar.gz"),
	}
)

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

func WithEndpoint(ep string) Option {
	return func(b *bs) error {
		b.endpoint = ep
		return nil
	}
}

func WithRegion(rg string) Option {
	return func(b *bs) error {
		b.region = rg
		return nil
	}
}

func WithAccessKey(ak string) Option {
	return func(b *bs) error {
		b.accessKey = ak
		return nil
	}
}

func WithSecretAccessKey(sak string) Option {
	return func(b *bs) error {
		b.secretAccessKey = sak
		return nil
	}
}

func WithToken(tk string) Option {
	return func(b *bs) error {
		b.token = tk
		return nil
	}
}

func WithMaxPartSizeKB(kb int) Option {
	return func(b *bs) error {
		b.maxPartSize = int64(kb) * 1024
		return nil
	}
}

func WithMaxPartSizeMB(mb int) Option {
	return func(b *bs) error {
		b.maxPartSize = int64(mb) * 1024 * 1024
		return nil
	}
}

func WithMaxPartSizeGB(gb int) Option {
	return func(b *bs) error {
		b.maxPartSize = int64(gb) * 1024 * 1024 * 1024
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
