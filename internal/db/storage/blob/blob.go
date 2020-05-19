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

package blob

import (
	"context"
	"io"
	"net/url"
	"reflect"

	"github.com/vdaas/vald/internal/errors"

	"gocloud.dev/blob"
)

type BucketURLOpener = blob.BucketURLOpener

type bucket struct {
	opener BucketURLOpener
	url    string
	bucket *blob.Bucket
}

type Bucket interface {
	Open(ctx context.Context) error
	Close() error
	Reader(ctx context.Context, key string) (io.ReadCloser, error)
	Writer(ctx context.Context, key string) (io.WriteCloser, error)
}

func NewBucket(opts ...Option) (Bucket, error) {
	b := new(bucket)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(b); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return b, nil
}

func (b *bucket) Open(ctx context.Context) (err error) {
	url, err := url.Parse(b.url)
	if err != nil {
		return err
	}
	b.bucket, err = b.opener.OpenBucketURL(ctx, url)
	return err
}

func (b *bucket) Close() error {
	if b.bucket != nil {
		return b.bucket.Close()
	}
	return nil
}

func (b *bucket) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	return b.bucket.NewReader(ctx, key, nil)
}

func (b *bucket) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	return b.bucket.NewWriter(ctx, key, nil)
}
