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
	"context"
	"io"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
)

type reader struct {
	eg      errgroup.Group
	service *s3.S3
	bucket  string
	key     string

	pr io.ReadCloser
	wg *sync.WaitGroup
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
	input := &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(r.key),
	}

	var pw io.WriteCloser

	r.pr, pw = io.Pipe()

	resp, err := r.service.GetObjectWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return errors.NewErrBlobNoSuchBucket(err, r.bucket)
			case s3.ErrCodeNoSuchKey:
				return errors.NewErrBlobNoSuchKey(err, r.key)
			}
		}

		return err
	}

	r.wg = new(sync.WaitGroup)
	r.wg.Add(1)

	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer r.wg.Done()
		defer resp.Body.Close()
		defer pw.Close()

		_, err = io.Copy(pw, resp.Body)
		return err
	}))

	return nil
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
