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
package cloudstorage

import (
	"context"
	"net/url"
	"reflect"

	iblob "github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcerrors"
)

type client struct {
	url *url.URL

	urlOpener *gcsblob.URLOpener
	bucket    *blob.Bucket

	readerOpts *blob.ReaderOptions
	writerOpts *blob.WriterOptions
}

// New returns blob.Bucket implementation.
func New(opts ...Option) (iblob.Bucket, error) {
	c := new(client)

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return c, nil
}

func (c *client) Open(ctx context.Context) (err error) {
	if c.url == nil {
		return errors.NewErrInvalidOption("url", c.url)
	}
	if c.urlOpener == nil {
		return errors.NewErrInvalidOption("urlOpener", c.urlOpener)
	}
	c.bucket, err = c.urlOpener.OpenBucketURL(ctx, c.url)
	return err
}

func (c *client) Close() error {
	if c.bucket == nil {
		return errors.ErrBucketNotOpened
	}
	return c.bucket.Close()
}

func (c *client) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	if c.bucket == nil {
		return nil, errors.ErrBucketNotOpened
	}
	rc, err := c.bucket.NewReader(ctx, key, c.readerOpts)
	if err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			log.Warn(err)
			return io.NopCloser(io.NewEOFReader()), nil
		}
		// rc is a nil *blob.Reader here; returning it would yield a non-nil
		// io.ReadCloser interface wrapping a nil pointer.
		return nil, err
	}
	return rc, nil
}

func (c *client) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	if c.bucket == nil {
		return nil, errors.ErrBucketNotOpened
	}
	return c.bucket.NewWriter(ctx, key, c.writerOpts)
}
