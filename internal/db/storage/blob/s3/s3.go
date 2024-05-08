//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

package s3

import (
	"context"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/reader"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/writer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type client struct {
	eg      errgroup.Group
	session *session.Session
	service *s3.S3
	bucket  string

	maxPartSize  int64
	maxChunkSize int64

	reader reader.Reader
	writer writer.Writer

	readerBackoffEnabled bool
	readerBackoffOpts    []backoff.Option
}

// New returns blob.Bucket implementation if no error occurs.
func New(opts ...Option) (b blob.Bucket, err error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	if c.session == nil {
		return nil, errors.NewErrInvalidOption("session", c.session)
	}

	c.service = s3.New(c.session)

	if c.writer == nil {
		c.writer = writer.New(
			writer.WithErrGroup(c.eg),
			writer.WithService(c.service),
			writer.WithBucket(c.bucket),
			writer.WithMaxPartSize(c.maxPartSize),
		)
	}

	if c.reader == nil {
		c.reader, err = reader.New(
			reader.WithErrGroup(c.eg),
			reader.WithService(c.service),
			reader.WithBucket(c.bucket),
			reader.WithMaxChunkSize(c.maxChunkSize),
			reader.WithBackoff(c.readerBackoffEnabled),
			reader.WithBackoffOpts(c.readerBackoffOpts...),
		)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Open does nothing. Always returns nil.
func (*client) Open(context.Context) (err error) {
	return nil
}

// Close does nothing. Always returns nil.
func (*client) Close() error {
	return nil
}

// Reader creates reader.Reader implementation and returns it.
// An error will be returned when the reader initialization fails or an error occurs in reader.Open.
func (c *client) Reader(ctx context.Context, key string) (rc io.ReadCloser, err error) {
	err = c.reader.Open(ctx, key)
	if err != nil {
		return nil, err
	}
	return c.reader, nil
}

// Writer creates writer.Writer implementation and returns it.
// An error will be returned when the writer initialization fails or an error occurs in writer.Open.
func (c *client) Writer(ctx context.Context, key string) (wc io.WriteCloser, err error) {
	err = c.writer.Open(ctx, key)
	if err != nil {
		return nil, err
	}
	return c.writer, nil
}
