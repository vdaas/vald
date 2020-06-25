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

package reader

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	ctxio "github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

type reader struct {
	eg      errgroup.Group
	service *s3.S3
	bucket  string
	key     string

	pr io.ReadCloser
	wg *sync.WaitGroup

	maxChunkSize int64
}

type Reader interface {
	Open(ctx context.Context) error
	io.ReadCloser
}

func New(opts ...Option) Reader {
	r := new(reader)
	for _, opt := range append(defaultOpts, opts...) {
		opt(r)
	}

	return r
}

func (r *reader) Open(ctx context.Context) (err error) {
	var pw io.WriteCloser

	r.pr, pw = io.Pipe()

	r.wg = new(sync.WaitGroup)
	r.wg.Add(1)
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer r.wg.Done()
		defer pw.Close()

		var offset int64

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			log.Debugf("loading %d bytes...", offset)
			body, err := r.getObject(ctx, offset, r.maxChunkSize)
			if err != nil {
				return err
			}

			body, err = ctxio.NewReadCloserWithContext(ctx, body)
			if err != nil {
				return err
			}

			chunk, err := io.Copy(pw, body)
			if err != nil {
				return err
			}

			if chunk < r.maxChunkSize {
				log.Debugf("read %d bytes.", offset+chunk)
				return nil
			}

			offset += chunk

			err = body.Close()
			if err != nil {
				return err
			}
		}
	}))

	return nil
}

func (r *reader) getObject(ctx context.Context, offset, length int64) (io.ReadCloser, error) {
	resp, err := r.service.GetObjectWithContext(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(r.bucket),
			Key:    aws.String(r.key),
			Range: aws.String("bytes=" +
				strconv.FormatInt(offset, 10) +
				"-" +
				strconv.FormatInt(offset+length-1, 10),
			),
		},
	)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, errors.NewErrBlobNoSuchBucket(err, r.bucket)
			case s3.ErrCodeNoSuchKey:
				return nil, errors.NewErrBlobNoSuchKey(err, r.key)
			case "InvalidRange":
				return ioutil.NopCloser(bytes.NewReader(nil)), nil
			}
		}

		return nil, err
	}

	return resp.Body, nil
}

func (r *reader) Close() error {
	if r.pr != nil {
		return r.pr.Close()
	}

	if r.wg != nil {
		r.wg.Wait()
	}

	return nil
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.pr == nil {
		return 0, errors.ErrStorageReaderNotOpened
	}

	return r.pr.Read(p)
}
