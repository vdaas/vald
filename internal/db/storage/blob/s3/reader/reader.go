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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/errors"
)

type reader struct {
	service *s3.S3
	bucket  string
	key     string

	resp *s3.GetObjectOutput
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

	r.resp, err = r.service.GetObjectWithContext(ctx, input)

	return err
}

func (r *reader) Close() error {
	if r.resp != nil {
		return r.resp.Body.Close()
	}

	return nil
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.resp == nil {
		return 0, errors.ErrStorageReaderNotOpened
	}

	return r.resp.Body.Read(p)
}
