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

package reader

import (
	"bytes"
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/vdaas/vald/internal/backoff"
	ctxio "github.com/vdaas/vald/internal/db/storage/blob/s3/reader/io"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type reader struct {
	eg      errgroup.Group
	service s3iface.S3API
	bucket  string

	pr io.ReadCloser
	wg *sync.WaitGroup

	ctxio ctxio.IO

	backoffEnabled bool
	backoffOpts    []backoff.Option
	bo             backoff.Backoff
	maxChunkSize   int64
}

var (
	errBlobNoSuchBucket      = new(errors.ErrBlobNoSuchBucket)
	errBlobNoSuchKey         = new(errors.ErrBlobNoSuchKey)
	errBlobInvalidChunkRange = new(errors.ErrBlobInvalidChunkRange)
)

// Reader is an interface that groups the basic Read and Close and Open methods.
type Reader interface {
	Open(ctx context.Context, key string) error
	io.ReadCloser
}

// New returns Reader implementation.
func New(opts ...Option) (Reader, error) {
	r := new(reader)
	for _, opt := range append(defaultOptions, opts...) {
		opt(r)
	}

	if r.ctxio == nil {
		return nil, errors.NewErrInvalidOption("ctxio", r.ctxio)
	}

	return r, nil
}

// Open creates io.Pipe. After reading the data from s3, make it available with Read method.
// Open method returns an error to align the interface, but it doesn't actually return an error.
func (r *reader) Open(ctx context.Context, key string) (err error) {
	var pw io.WriteCloser

	r.pr, pw = io.Pipe()

	r.wg = new(sync.WaitGroup)
	r.wg.Add(1)
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer r.wg.Done()
		defer pw.Close()

		var offset int64

		if r.backoffEnabled {
			r.bo = backoff.New(r.backoffOpts...)
		}

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			body, err := func() (io.Reader, error) {
				if r.backoffEnabled {
					return r.getObjectWithBackoff(ctx, key, offset, r.maxChunkSize)
				}
				return r.getObject(ctx, key, offset, r.maxChunkSize)
			}()
			if err != nil {
				if errors.As(err, &errBlobNoSuchBucket) ||
					errors.As(err, &errBlobNoSuchKey) ||
					errors.As(err, &errBlobInvalidChunkRange) {
					log.Warn(err)
					return nil
				}
				return err
			}

			body, err = r.ctxio.NewReaderWithContext(ctx, body)
			if err != nil {
				return err
			}

			chunk, err := io.Copy(pw, body)
			if err != nil && !errors.Is(err, io.EOF) {
				return err
			}

			if chunk < r.maxChunkSize {
				return nil
			}

			offset += chunk
		}
	}))

	return nil
}

func (r *reader) getObjectWithBackoff(ctx context.Context, key string, offset, length int64) (res io.Reader, err error) {
	if !r.backoffEnabled || r.bo == nil {
		return r.getObject(ctx, key, offset, length)
	}
	_, err = r.bo.Do(ctx, func(ctx context.Context) (interface{}, bool, error) {
		res, err = r.getObject(ctx, key, offset, length)
		if err != nil {
			if errors.As(err, &errBlobNoSuchBucket) ||
				errors.As(err, &errBlobNoSuchKey) ||
				errors.As(err, &errBlobInvalidChunkRange) {
				return res, false, err
			}
			return res, true, err
		}
		return res, false, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *reader) getObject(ctx context.Context, key string, offset, length int64) (io.Reader, error) {
	rng := aws.String("bytes=" + strconv.FormatInt(offset, 10) + "-" + strconv.FormatInt(offset+length-1, 10))
	log.Debugf("reading %s", *rng)
	resp, err := r.service.GetObjectWithContext(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(r.bucket),
			Key:    aws.String(key),
			Range:  rng,
		},
	)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, errors.NewErrBlobNoSuchBucket(err, r.bucket)
			case s3.ErrCodeNoSuchKey:
				return nil, errors.NewErrBlobNoSuchKey(err, key)
			case "InvalidRange":
				return nil, errors.NewErrBlobInvalidChunkRange(err, *rng)
			}
		}
		return nil, err
	}

	res, err := r.ctxio.NewReadCloserWithContext(ctx, resp.Body)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	defer func() {
		e := res.Close()
		if e != nil {
			log.Warn(e)
		}
	}()

	_, err = io.Copy(buf, res)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Close closes the reader.
func (r *reader) Close() error {
	if r.bo != nil {
		defer r.bo.Close()
	}
	if r.pr != nil {
		return r.pr.Close()
	}

	if r.wg != nil {
		r.wg.Wait()
	}

	return nil
}

// Read reads up to len(p) bytes and returns the number of bytes read.
func (r *reader) Read(p []byte) (n int, err error) {
	if r.pr == nil {
		return 0, errors.ErrStorageReaderNotOpened
	}

	return r.pr.Read(p)
}
