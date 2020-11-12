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

package s3

import (
	"context"
	"io"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/reader"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/writer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

type client struct {
	eg      errgroup.Group
	session *session.Session
	service *s3.S3
	bucket  string

	maxPartSize  int64
	maxChunkSize int64

	readerBackoffEnabled bool
	readerBackoffOpts    []backoff.Option

	readerFunc func(key string) (reader.Reader, error)
	writerFunc func(key string) writer.Writer
}

func New(opts ...Option) (blob.Bucket, error) {
	c := new(client)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	if c.session == nil {
		return nil, errors.NewErrInvalidOption("session", c.session)
	}

	if c.readerFunc == nil {
		return nil, errors.NewErrInvalidOption("readerFunc", c.readerFunc)
	}

	if c.writerFunc == nil {
		return nil, errors.NewErrInvalidOption("writerFunc", c.writerFunc)
	}

	c.service = s3.New(c.session)

	return c, nil
}

func (c *client) Open(ctx context.Context) (err error) {
	return nil
}

func (c *client) Close() error {
	return nil
}

func (c *client) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	if c.readerFunc == nil {
		return nil, errors.ErrNilObject
	}
	r, err := c.readerFunc(key)
	if err != nil {
		return nil, err
	}

	return r, r.Open(ctx)
}

func (c *client) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	if c.writerFunc == nil {
		return nil, errors.ErrNilObject
	}
	w := c.writerFunc(key)
	return w, w.Open(ctx)
}
