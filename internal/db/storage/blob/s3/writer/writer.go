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

package writer

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

type writer struct {
	service *s3.S3
	bucket  string
	key     string

	maxPartSize int64

	resp *s3.CreateMultipartUploadOutput

	completedParts []*s3.CompletedPart
}

type Writer interface {
	Open(ctx context.Context) error
	io.WriteCloser
}

func New(opts ...Option) Writer {
	w := new(writer)
	for _, opt := range append(defaultOpts, opts...) {
		opt(w)
	}

	return w
}

func (w *writer) Open(ctx context.Context) (err error) {
	input := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(w.bucket),
		Key:    aws.String(w.key),
	}

	w.completedParts = make([]*s3.CompletedPart, 0)
	w.resp, err = w.service.CreateMultipartUploadWithContext(ctx, input)

	return err
}

func (w *writer) Close() error {
	if w.resp != nil {
		return w.complete()
	}

	return nil
}

func (w *writer) Write(p []byte) (n int, err error) {
	if w.resp == nil {
		return 0, errors.ErrStorageWriterNotOpened
	}

	// TODO: do concurrently
	var cur, pl int64
	size := int64(len(p))
	remaining := size
	pn := int64(1)
	for cur = 0; remaining != 0; cur += pl {
		if remaining < w.maxPartSize {
			pl = remaining
		} else {
			pl = w.maxPartSize
		}

		cp, err := w.upload(p[cur:cur+pl], pn)
		if err != nil {
			return int(size - remaining), errors.Wrap(w.abort(), err.Error())
		}

		remaining -= pl
		pn++
		w.completedParts = append(w.completedParts, cp)
	}

	return int(size), nil
}

func (w *writer) upload(p []byte, n int64) (*s3.CompletedPart, error) {
	pn := aws.Int64(n)
	input := &s3.UploadPartInput{
		Body:          bytes.NewReader(p),
		Bucket:        w.resp.Bucket,
		Key:           w.resp.Key,
		PartNumber:    pn,
		UploadId:      w.resp.UploadId,
		ContentLength: aws.Int64(int64(len(p))),
	}

	// TODO: backoff
	res, err := w.service.UploadPart(input)
	if err != nil {
		return nil, err
	}

	return &s3.CompletedPart{
		ETag:       res.ETag,
		PartNumber: pn,
	}, nil
}

func (w *writer) abort() error {
	input := &s3.AbortMultipartUploadInput{
		Bucket:   w.resp.Bucket,
		Key:      w.resp.Key,
		UploadId: w.resp.UploadId,
	}

	res, err := w.service.AbortMultipartUpload(input)

	log.Debugf("s3 upload aborted: %s", res.String())

	return err
}

func (w *writer) complete() error {
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   w.resp.Bucket,
		Key:      w.resp.Key,
		UploadId: w.resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: w.completedParts,
		},
	}

	res, err := w.service.CompleteMultipartUpload(input)
	if err != nil {
		return err
	}

	log.Debugf("s3 upload completed: %s", res.String())

	return nil
}
