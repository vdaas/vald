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

package writer

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3manager"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type writer struct {
	eg          errgroup.Group
	s3manager   s3manager.S3Manager
	service     s3iface.S3API
	pw          io.WriteCloser
	pr          *io.PipeReader
	wg          *sync.WaitGroup
	err         error // the result of the upload, valid after wg.Wait
	bucket      string
	contentType string
	maxPartSize int64
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
// The result of the upload is returned from Close.
// If the previous stream was not closed, it is aborted and waited for so that its goroutine does not leak.
func (w *writer) Open(ctx context.Context, key string) (err error) {
	if w.pr != nil {
		if cerr := w.pr.CloseWithError(io.ErrClosedPipe); cerr != nil {
			log.Warn(cerr)
		}
		if w.wg != nil {
			w.wg.Wait()
		}
	}

	pr, pw := io.Pipe()
	w.pr, w.pw, w.err = pr, pw, nil

	wg := new(sync.WaitGroup)
	w.wg = wg
	wg.Add(1)

	w.eg.Go(func() (err error) {
		defer wg.Done()

		err = safety.RecoverFunc(func() error {
			return w.upload(ctx, key, pr)
		})()
		// The reader must be closed with the error so that a blocked Write returns it
		// instead of waiting for a consumer that is gone.
		if cerr := pr.CloseWithError(err); cerr != nil {
			log.Warn(cerr)
		}
		w.err = err // read by Close after wg.Wait

		// Close returns the upload error to the caller. It is not returned to the
		// errgroup as well: that may be the process-wide group, where every failed
		// upload would be recorded until shutdown and cancel the daemon context.
		return nil
	})

	return nil
}

// Close closes the writer, waits for the upload to complete and returns the result of the upload.
func (w *writer) Close() (err error) {
	if w.pw != nil {
		err = w.pw.Close()
	}

	if w.wg != nil {
		w.wg.Wait()
		if w.err != nil {
			return w.err
		}
	}

	return err
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
		Bucket:      new(w.bucket),
		Key:         new(key),
		Body:        body,
		ContentType: new(w.contentType),
	}

	res, err := client.UploadWithContext(ctx, input)
	if err != nil {
		log.Errorf("upload failed with error: %s", err)
		return err
	}

	log.Infof("s3 upload completed: %s", res.Location)

	return nil
}
