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

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/reader"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/writer"
)

type s3client struct {
	session *session.Session
	service *s3.S3
	bucket  string
}

func New(opts ...Option) blob.Bucket {
	s := new(s3client)
	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}

	s.service = s3.New(s.session)

	return s
}

func (s *s3client) Open(ctx context.Context) (err error) {
	return nil
}

func (s *s3client) Close() error {
	return nil
}

func (s *s3client) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	r := reader.New(
		reader.WithService(s.service),
		reader.WithBucket(s.bucket),
		reader.WithKey(key),
	)

	return r, r.Open(ctx)
}

func (s *s3client) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	w := writer.New(
		writer.WithService(s.service),
		writer.WithBucket(s.bucket),
		writer.WithKey(key),
	)

	return w, w.Open(ctx)
}
