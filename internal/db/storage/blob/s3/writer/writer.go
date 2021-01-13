//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"context"
	"io"
	"reflect"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3manager"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

type writer struct {
	eg        errgroup.Group
	s3manager s3manager.S3Manager
	service   s3iface.S3API
	bucket    string

	contentType string
	maxPartSize int64

	pw io.WriteCloser
	wg *sync.WaitGroup
}

// Writer represents an interface to write to s3.
type Writer interface {
	Open(ctx context.Context, key string) error
	io.WriteCloser
}

// New returns Writer implementation.
func New(opts ...Option) Writer {
	w := &writer{
		s3manager: s3manager.New(),
	}
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(w); err != nil {
			log.Warn(errors.ErrOptionFailed(err, reflect.ValueOf(opt)))
		}
	}

	return w
}

// Open creates io.Pipe. When the write method is called, the written data will be uploaded to s3.
// Open method returns an error to align the interface, but it doesn't actually return an error.
func (w *writer) Open(ctx context.Context, key string) (err error) {
	w.wg = new(sync.WaitGroup)

	var pr io.ReadCloser

	pr, w.pw = io.Pipe()

	w.wg.Add(1)

	w.eg.Go(safety.RecoverFunc(func() (err error) {
		defer w.wg.Done()
		defer pr.Close()

		return w.upload(ctx, key, pr)
	}))

	return err
}

// Close closes the writer.
func (w *writer) Close() error {
	if w.pw != nil {
		return w.pw.Close()
	}

	if w.wg != nil {
		w.wg.Wait()
	}

	return nil
}

// Write writes len(p) bytes from p to the underlying data stream. The written data will be uploaded to s3.
func (w *writer) Write(p []byte) (n int, err error) {
	if w.pw == nil {
		return 0, errors.ErrStorageWriterNotOpened
	}

	return w.pw.Write(p)
}

func (w *writer) upload(ctx context.Context, key string, body io.Reader) (err error) {
	client := w.s3manager.NewUploaderWithClient(
		w.service,
		func(u *s3manager.Uploader) {
			u.PartSize = w.maxPartSize
		},
	)
	input := &s3manager.UploadInput{
		Bucket:      aws.String(w.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(w.contentType),
	}

	res, err := client.UploadWithContext(ctx, input)
	if err != nil {
		log.Error("upload failed with error: ", err)
		return err
	}

	log.Infof("s3 upload completed: %s", res.Location)

	return nil
}
