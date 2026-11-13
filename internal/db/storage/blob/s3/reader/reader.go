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

package reader

import (
	"bytes"
	"context"
	"strconv"

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
	eg             errgroup.Group
	service        s3iface.S3API
	ctxio          ctxio.IO
	bucket         string
	backoffOpts    []backoff.Option
	maxChunkSize   int64
	backoffEnabled bool
}

// stream is the state of a single Open. Every Open returns its own stream, so Read and Close
// of an earlier stream never touch the pipe, wait group or backoff of a later one.
type stream struct {
	pr io.ReadCloser
	bo backoff.Backoff
	wg *sync.WaitGroup
	// cancel stops the goroutine started by Open, including its backoff sleeps, which
	// closing the pipe does not interrupt. It is nil for a stream built without Open.
	cancel context.CancelFunc
}

// isBlobNotRetryable reports whether err is one of the blob errors that must not be retried.
// errors.As writes to its target, so the targets are local to stay race-free across concurrent callers.
func isBlobNotRetryable(err error) bool {
	var (
		noSuchBucket      *errors.ErrBlobNoSuchBucket
		noSuchKey         *errors.ErrBlobNoSuchKey
		invalidChunkRange *errors.ErrBlobInvalidChunkRange
	)
	return errors.As(err, &noSuchBucket) ||
		errors.As(err, &noSuchKey) ||
		errors.As(err, &invalidChunkRange)
}

// Reader is an interface that opens an independent stream for each key.
type Reader interface {
	Open(ctx context.Context, key string) (io.ReadCloser, error)
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

// Open creates io.Pipe and returns a stream that makes the data read from s3 available with its Read method.
// Open method returns an error to align the interface, but it doesn't actually return an error.
// An error that occurs while reading from s3 is returned from Read instead of io.EOF.
// The returned io.ReadCloser is independent of the streams returned by other Open calls.
func (r *reader) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(ctx)
	var bo backoff.Backoff
	if r.backoffEnabled {
		bo = backoff.New(r.backoffOpts...)
	}
	wg.Add(1)
	r.eg.Go(func() error {
		defer wg.Done()
		defer cancel()

		err := safety.RecoverFunc(func() (err error) {
			var offset int64

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				body, err := func() (io.Reader, error) {
					if r.backoffEnabled {
						return r.getObjectWithBackoff(ctx, bo, key, offset, r.maxChunkSize)
					}
					return r.getObject(ctx, key, offset, r.maxChunkSize)
				}()
				if err != nil {
					if isBlobNotRetryable(err) {
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
					// The consumer closed the stream before reading everything;
					// stop without reporting an error to the shared errgroup.
					if errors.Is(err, io.ErrClosedPipe) {
						return nil
					}
					return err
				}

				if chunk < r.maxChunkSize {
					return nil
				}

				offset += chunk
			}
		})()
		// CloseWithError(nil) is io.EOF for the reader, so only a real failure is propagated as an error.
		if cerr := pw.CloseWithError(err); cerr != nil {
			log.Warn(cerr)
		}
		// Read returns the error to the consumer. It is not returned to the errgroup as well:
		// that may be the process-wide group, where every failed read would be recorded until
		// shutdown and cancel the group's context.
		return nil
	})

	return &stream{pr: pr, wg: wg, bo: bo, cancel: cancel}, nil
}

func (r *reader) getObjectWithBackoff(
	ctx context.Context, bo backoff.Backoff, key string, offset, length int64,
) (res io.Reader, err error) {
	if !r.backoffEnabled || bo == nil {
		return r.getObject(ctx, key, offset, length)
	}
	_, err = bo.Do(ctx, func(ctx context.Context) (any, bool, error) {
		res, err = r.getObject(ctx, key, offset, length)
		if err != nil {
			if isBlobNotRetryable(err) {
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

func (r *reader) getObject(
	ctx context.Context, key string, offset, length int64,
) (io.Reader, error) {
	rng := new("bytes=" + strconv.FormatInt(offset, 10) + "-" + strconv.FormatInt(offset+length-1, 10))
	log.Debugf("reading %s", *rng)
	resp, err := r.service.GetObjectWithContext(
		ctx,
		&s3.GetObjectInput{
			Bucket: new(r.bucket),
			Key:    new(key),
			Range:  rng,
		},
	)
	if err != nil {
		var aerr awserr.Error
		if errors.As(err, &aerr) {
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
		if resp.Body != nil {
			if e := resp.Body.Close(); e != nil {
				log.Warn(e)
			}
		}
		return nil, err
	}

	// Presize to the chunk length so a large chunk is not grown by repeated doubling. The
	// Content-Length header is the server's claim, so it is capped at the length requested.
	buf := new(bytes.Buffer)
	if resp.ContentLength != nil && *resp.ContentLength > 0 && length > 0 {
		buf.Grow(int(min(*resp.ContentLength, length)))
	}

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

// Close closes the stream and waits for the goroutine started by its Open to finish.
func (s *stream) Close() (err error) {
	if s.pr != nil {
		err = s.pr.Close()
	}

	// Closing the pipe does not wake a goroutine that sleeps in a backoff retry, so without
	// this Wait would block until the retries are exhausted.
	if s.cancel != nil {
		s.cancel()
	}

	// the goroutine uses bo, so wait for it before closing the backoff.
	if s.wg != nil {
		s.wg.Wait()
	}

	if s.bo != nil {
		s.bo.Close()
	}

	return err
}

// Read reads up to len(p) bytes and returns the number of bytes read.
func (s *stream) Read(p []byte) (n int, err error) {
	if s.pr == nil {
		return 0, errors.ErrStorageReaderNotOpened
	}

	return s.pr.Read(p)
}
