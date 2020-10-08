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

type readWriter interface {
	NewReader(opts ...reader.Option) reader.Reader
	NewWriter(opts ...writer.Option) writer.Writer
}

type rw struct{}

func newRW() readWriter {
	return new(rw)
}

func (*rw) NewReader(opts ...reader.Option) reader.Reader {
	return reader.New(opts...)
}

func (*rw) NewWriter(opts ...writer.Option) writer.Writer {
	return writer.New(opts...)
}

type client struct {
	eg      errgroup.Group
	session *session.Session
	service *s3.S3
	bucket  string

	maxPartSize  int64
	maxChunkSize int64

	readerBackoffEnabled bool
	readerBackoffOpts    []backoff.Option

	readWriter readWriter
}

// New returns Bucket implementation.
func New(opts ...Option) (blob.Bucket, error) {
	c := &client{
		readWriter: newRW(),
	}
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	if c.session == nil {
		return nil, errors.ErrS3SessionNotFound
	}

	c.service = s3.New(c.session)

	return c, nil
}

// Open do nothing
func (c *client) Open(ctx context.Context) (err error) {
	return nil
}

// Close do nothing
func (c *client) Close() error {
	return nil
}

// Reader creates s3 reader and calls reader.Open.
func (c *client) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	r := c.readWriter.NewReader(
		reader.WithErrGroup(c.eg),
		reader.WithService(c.service),
		reader.WithBucket(c.bucket),
		reader.WithKey(key),
		reader.WithMaxChunkSize(c.maxChunkSize),
		reader.WithBackoff(c.readerBackoffEnabled),
		reader.WithBackoffOpts(c.readerBackoffOpts...),
	)

	return r, r.Open(ctx)
}

// Writer creates s3 writer and calls writer.Open.
func (c *client) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	w := c.readWriter.NewWriter(
		writer.WithErrGroup(c.eg),
		writer.WithService(c.service),
		writer.WithBucket(c.bucket),
		writer.WithKey(key),
		writer.WithMaxPartSize(c.maxPartSize),
	)

	return w, w.Open(ctx)
}
