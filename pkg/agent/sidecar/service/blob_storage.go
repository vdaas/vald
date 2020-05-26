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

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/session"
	"github.com/vdaas/vald/internal/errors"
)

type BlobStorage interface {
	Start(ctx context.Context) (<-chan error, error)
}

type bs struct {
	storageType config.BlobStorageType
	bucketName  string

	endpoint        string
	region          string
	accessKey       string
	secretAccessKey string
	token           string

	bucket blob.Bucket
}

func NewBlobStorage(opts ...BlobStorageOption) (BlobStorage, error) {
	b := new(bs)
	for _, opt := range append(defaultBlobStorageOpts, opts...) {
		if err := opt(b); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return b, nil
}

func (b *bs) initBucket() (err error) {
	switch b.storageType {
	case config.S3:
		s, err := session.New(
			session.WithEndpoint(b.endpoint),
			session.WithRegion(b.region),
			session.WithAccessKey(b.accessKey),
			session.WithSecretAccessKey(b.secretAccessKey),
			session.WithToken(b.token),
		).Session()
		if err != nil {
			return err
		}

		b.bucket = s3.New(
			s3.WithSession(s),
			s3.WithBucket(b.bucketName),
		)

	default:
		return errors.ErrInvalidStorageType
	}

	return nil
}

func (b *bs) Start(ctx context.Context) (<-chan error, error) {
	return nil, nil
}
