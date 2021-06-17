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
package downloader

import (
	"context"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/arn"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

type Client interface {
	Open(ctx context.Context, input *s3.GetObjectInput) (rc io.ReadCloser, err error)
	Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput) (n int64, err error)
}

type client struct {
	partSize       int64
	concurrency    int
	client         manager.DownloadAPIClient
	clientOptions  []func(*s3.Options)
	bufferProvider manager.WriterReadFromProvider
	copier         io.Copier
	eg             errgroup.Group
	bo             backoff.Backoff
}

type rangeWriter struct {
	w       io.WriterAt
	offset  int64
	size    int64
	current int64
	rng     string
}

func rangeParam(offset, size int64) string {
	return "bytes=" +
		strconv.FormatInt(offset, 10) +
		"-" +
		strconv.FormatInt(offset+size-1, 10)
}

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return c, nil
}

func (c *client) Open(ctx context.Context, input *s3.GetObjectInput) (rc io.ReadCloser, err error) {
	buf := newBuffer(ctx, 10)
	c.eg.Go(safety.RecoverFunc(func() error {
		_, err = c.Download(ctx, buf, input)
		if err != nil {
			log.Error(err)
		}
		return nil
	}))
	return buf, nil
}

func (c *client) Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput) (written int64, err error) {
	if err := arn.ValidateSupportedARNType(aws.ToString(input.Bucket)); err != nil {
		return 0, err
	}
	if rng := aws.ToString(input.Range); len(rng) > 0 {
		written, _, err = c.downloadChunk(ctx, input, &rangeWriter{
			w:       w,
			offset:  0,
			size:    0,
			current: 0,
			rng:     rng,
		})
		return written, err
	}
	pos := c.partSize
	n, total, err := c.downloadChunk(ctx, input, &rangeWriter{
		w:       w,
		offset:  0,
		size:    c.partSize,
		current: 0,
	})
	if err == nil {
		atomic.AddInt64(&written, n)
	}
	if total >= 0 {
		eg, egctx := errgroup.New(ctx)
		eg.Limitation(c.concurrency)
		for pos < total {
			offset := pos
			eg.Go(safety.RecoverFunc(func() error {
				n, _, err := c.downloadChunk(egctx, input, &rangeWriter{
					w:      w,
					offset: offset,
					size:   c.partSize,
				})
				if err != nil && !errors.Is(err, io.EOF) {
					var rerr interface {
						HTTPStatusCode() int
					}
					if errors.As(err, &rerr) {
						if rerr.HTTPStatusCode() == http.StatusRequestedRangeNotSatisfiable {
							err = nil
						}
					}
					return err
				}
				if n > 0 {
					atomic.AddInt64(&written, n)
				}
				return nil
			}))
			pos += c.partSize
		}
		err = eg.Wait()
		return written, err
	}
	for err == nil {
		offset := pos
		pos += c.partSize
		n, _, err = c.downloadChunk(ctx, input, &rangeWriter{
			w:      w,
			offset: offset,
			size:   c.partSize,
		})
		if n > 0 {
			written += n
		}
	}
	if err != nil {
		var rerr interface {
			HTTPStatusCode() int
		}
		if errors.As(err, &rerr) {
			if rerr.HTTPStatusCode() == http.StatusRequestedRangeNotSatisfiable {
				err = nil
			}
		}
	}
	return written, err
}

func (c *client) downloadChunk(ctx context.Context, in *s3.GetObjectInput, w *rangeWriter) (n, total int64, err error) {
	if len(w.rng) == 0 {
		w.rng = rangeParam(w.offset, w.size)
	}
	in.Range = aws.String(w.rng)
	if c.bo == nil {
		return c.tryDownloadChunk(ctx, in, w)
	}
	_, err = c.bo.Do(ctx, func(ctx context.Context) (interface{}, bool, error) {
		n, total, err = c.tryDownloadChunk(ctx, in, w)
		if err == nil {
			return n, false, nil
		}
		berr, ok := err.(*errors.ErrS3ReadingBody)
		if !ok {
			return n, false, err
		}
		err = berr.Unwrap()
		return n, err != nil && !errors.Is(err, io.EOF), errors.Wrapf(err, "error object part body download interrupted %s", aws.ToString(in.Key))
	})
	return n, total, err
}

func (c *client) tryDownloadChunk(ctx context.Context, in *s3.GetObjectInput, w io.Writer) (n, total int64, err error) {
	if c.bufferProvider != nil {
		var cleanup func()
		w, cleanup = c.bufferProvider.GetReadFrom(w)
		defer cleanup()
	}

	total = int64(-1)
	resp, err := c.client.GetObject(ctx, in, c.clientOptions...)
	if err != nil {
		return 0, total, err
	}
	defer func() {
		_, cerr := io.Copy(io.Discard, resp.Body)
		if cerr != nil {
			err = errors.Wrap(err, cerr.Error())
		}
		cerr = resp.Body.Close()
		if cerr != nil {
			err = errors.Wrap(err, cerr.Error())
		}
	}()
	if resp.ContentRange == nil && resp.ContentLength > 0 {
		total = resp.ContentLength
	} else {
		parts := strings.Split(*resp.ContentRange, "/")
		totalStr := parts[len(parts)-1]
		if totalStr != "*" {
			t, err := strconv.ParseInt(totalStr, 10, 64)
			if err == nil {
				total = t
			}
		}
	}
	if c.copier != nil {
		n, err = c.copier.CopyWithContext(ctx, w, resp.Body)
	} else {
		n, err = io.CopyWithContext(ctx, w, resp.Body)
	}
	if err != nil {
		return n, total, errors.NewErrS3ReadingBody(err)
	}
	return n, total, nil
}

func (r *rangeWriter) Write(p []byte) (n int, err error) {
	if r.current >= r.size && len(r.rng) == 0 {
		return 0, io.EOF
	}
	n, err = r.w.WriteAt(p, r.offset+r.current)
	r.current += int64(n)
	return
}
