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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

type writer struct {
	eg      errgroup.Group
	service *s3.S3
	bucket  string
	key     string

	maxPartSize int64
	multipart   bool

	// for single uploads
	pw io.WriteCloser

	// for multipart uploads
	ctx  context.Context
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
	w.ctx = ctx

	if !w.multipart {
		var pr io.ReadCloser

		pr, w.pw = io.Pipe()

		w.eg.Go(safety.RecoverFunc(func() (err error) {
			defer pr.Close()

			return w.upload(pr)
		}))
	}

	return err
}

func (w *writer) Close() error {
	if !w.multipart && w.pw != nil {
		return w.pw.Close()
	}

	if w.multipart && w.resp != nil {
		return w.completeMultipartUpload()
	}

	return nil
}

func (w *writer) Write(p []byte) (n int, err error) {
	if w.ctx == nil {
		return 0, errors.ErrStorageWriterNotOpened
	}

	if w.multipart {
		return w.multipartUpload(p)
	}

	if w.pw == nil {
		return 0, errors.ErrStorageWriterNotOpened
	}

	return w.pw.Write(p)
}

func (w *writer) upload(body io.Reader) (err error) {
	uploader := s3manager.NewUploaderWithClient(w.service)
	input := &s3manager.UploadInput{
		Bucket: aws.String(w.bucket),
		Key:    aws.String(w.key),
		Body:   body,
	}

	res, err := uploader.UploadWithContext(w.ctx, input)
	if err != nil {
		return err
	}

	log.Debugf("s3 upload completed: %s", res.Location)

	return nil
}

func (w *writer) multipartUpload(p []byte) (n int, err error) {
	if w.resp == nil {
		input := &s3.CreateMultipartUploadInput{
			Bucket: aws.String(w.bucket),
			Key:    aws.String(w.key),
		}

		w.completedParts = make([]*s3.CompletedPart, 0)
		w.resp, err = w.service.CreateMultipartUploadWithContext(w.ctx, input)
	}

	// TODO: do concurrently
	var cur, pl int64
	size := int64(len(p))
	remaining := size
	pn := int64(1)
	for cur = 0; remaining != 0; cur += pl {
		select {
		case <-w.ctx.Done():
			return int(size - remaining), errors.Wrap(
				w.abortMultipartUpload(),
				w.ctx.Err().Error(),
			)
		default:
		}

		if remaining < w.maxPartSize {
			pl = remaining
		} else {
			pl = w.maxPartSize
		}

		cp, err := w.uploadPart(p[cur:cur+pl], pn)
		if err != nil {
			return int(size - remaining), errors.Wrap(
				w.abortMultipartUpload(),
				err.Error(),
			)
		}

		remaining -= pl
		pn++
		w.completedParts = append(w.completedParts, cp)
	}

	return int(size), nil
}

func (w *writer) uploadPart(p []byte, n int64) (*s3.CompletedPart, error) {
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
	res, err := w.service.UploadPartWithContext(w.ctx, input)
	if err != nil {
		return nil, err
	}

	return &s3.CompletedPart{
		ETag:       res.ETag,
		PartNumber: pn,
	}, nil
}

func (w *writer) abortMultipartUpload() error {
	input := &s3.AbortMultipartUploadInput{
		Bucket:   w.resp.Bucket,
		Key:      w.resp.Key,
		UploadId: w.resp.UploadId,
	}

	// AbortMultipartUploadWithContext could be used here,
	// however abort should be called even when the context cancelled.
	// this is why abort called without contexts.
	res, err := w.service.AbortMultipartUpload(input)

	log.Debugf("s3 multipart upload aborted: %s", res.String())

	return err
}

func (w *writer) completeMultipartUpload() error {
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   w.resp.Bucket,
		Key:      w.resp.Key,
		UploadId: w.resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: w.completedParts,
		},
	}

	res, err := w.service.CompleteMultipartUploadWithContext(w.ctx, input)
	if err != nil {
		return err
	}

	log.Debugf("s3 multipart upload completed: %s", res.String())

	return nil
}
