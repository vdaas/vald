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

// Package service
package service

type BlobStorageOption func(b *bs) error

var (
	defaultBlobStorageOpts = []BlobStorageOption{
		WithBlobStorageCompressAlgorithm("gzip"),
		WithBlobStorageCompressionLevel(-1),
		WithBlobStorageFilenameSuffix(".tar.gz"),
	}
)

func WithBlobStorageType(bst string) BlobStorageOption {
	return func(b *bs) error {
		b.storageType = bst
		return nil
	}
}

func WithBlobStorageBucketName(bn string) BlobStorageOption {
	return func(b *bs) error {
		b.bucketName = bn
		return nil
	}
}

func WithBlobStorageFilename(fn string) BlobStorageOption {
	return func(b *bs) error {
		b.filename = fn
		return nil
	}
}

func WithBlobStorageFilenameSuffix(sf string) BlobStorageOption {
	return func(b *bs) error {
		b.suffix = sf
		return nil
	}
}

func WithBlobStorageEndpoint(ep string) BlobStorageOption {
	return func(b *bs) error {
		b.endpoint = ep
		return nil
	}
}

func WithBlobStorageRegion(rg string) BlobStorageOption {
	return func(b *bs) error {
		b.region = rg
		return nil
	}
}

func WithBlobStorageAccessKey(ak string) BlobStorageOption {
	return func(b *bs) error {
		b.accessKey = ak
		return nil
	}
}

func WithBlobStorageSecretAccessKey(sak string) BlobStorageOption {
	return func(b *bs) error {
		b.secretAccessKey = sak
		return nil
	}
}

func WithBlobStorageToken(tk string) BlobStorageOption {
	return func(b *bs) error {
		b.token = tk
		return nil
	}
}

func WithBlobStorageMultipartUpload(enabled bool) BlobStorageOption {
	return func(b *bs) error {
		b.multipartUpload = enabled
		return nil
	}
}

func WithBlobStorageCompressAlgorithm(al string) BlobStorageOption {
	return func(b *bs) error {
		b.compressAlgorithm = al
		return nil
	}
}

func WithBlobStorageCompressionLevel(level int) BlobStorageOption {
	return func(b *bs) error {
		b.compressionLevel = level
		return nil
	}
}
