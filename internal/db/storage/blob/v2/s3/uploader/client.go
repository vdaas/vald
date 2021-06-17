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

// Package uploader uploads operations for AWS S3 objects
package uploader

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/arn"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/pool"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
)

type Uploader interface {
	Open(ctx context.Context, input *s3.PutObjectInput) (wc io.WriteCloser, err error)
	Upload(ctx context.Context, r io.Reader, input *s3.PutObjectInput) (n int64, err error)
}

type client struct {
	partSize       int64
	concurrency    int
	maxUploadParts int32
	client         manager.UploadAPIClient
	clientOptions  []func(*s3.Options)
	bufferProvider manager.ReadSeekerWriteToProvider
	partPool       pool.ByteSlicePool
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

func (c *client) Open(ctx context.Context, input *s3.PutObjectInput) (rc io.WriteCloser, err error) {
	// buf := newBuffer(ctx, 10)
	// c.eg.Go(safety.RecoverFunc(func() error {
	// 	_, err = c.Download(ctx, buf, input)
	// 	if err != nil {
	// 		log.Error(err)
	// 	}
	// 	return nil
	// }))
	// return buf, nil
	return nil, nil
}

func (c *client) Upload(ctx context.Context, r io.Reader, input *s3.PutObjectInput) (written int64, err error) {
	if err := arn.ValidateSupportedARNType(aws.ToString(input.Bucket)); err != nil {
		return 0, err
	}

	var objectSize int64 = -1
	switch r := input.Body.(type) {
	case io.Seeker:
		cur, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize upload: %w", err)
		}
		end, err := r.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize upload: %w", err)
		}

		_, err = r.Seek(cur, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize upload: %w", err)
		}
		objectSize = end - cur
		if objectSize/c.partSize >= int64(c.maxUploadParts) {
			c.partSize = (objectSize / int64(c.maxUploadParts)) + 1
		}
	}
	if c.partPool == nil || c.partPool.SliceSize() != c.partSize {
		c.partPool = pool.NewByteSlicePool(c.partSize)
	} else {
		c.partPool = &pool.ReturnCapacityPoolCloser{ByteSlicePool: c.partPool}
	}
	c.partPool.ModifyCapacity(objectSize + 1)
	defer c.partPool.Close()
	if c.partSize < MinUploadPartSize {
		return 0, fmt.Errorf("part size must be at least %d bytes", MinUploadPartSize)
	}

	// Do one read to determine if we have more than one part
	reader, _, cleanup, err := c.nextReader()
	if err == io.EOF { // single part
		defer cleanup()
		params := input
		params.Body = reader
		var locationRecorder recordLocationClient
		out, err := c.client.PutObject(ctx, params, append(c.clientOptions, locationRecorder.WrapClient())...)
		if err != nil {
			return nil, err
		}

		return &UploadOutput{
			Location:  locationRecorder.location,
			VersionID: out.VersionId,
		}, nil
	} else if err != nil {
		cleanup()
		return nil, fmt.Errorf("read upload data failed: %w", err)
	}
	mu := multiuploader{uploader: u}
	return written, err
}

func (c *client) nextReader() (io.ReadSeeker, int, func(), error) {
	switch r := u.in.Body.(type) {
	case readerAtSeeker:
		var err error

		n := c.partSize
		if u.totalSize >= 0 {
			bytesLeft := u.totalSize - u.readerPos

			if bytesLeft <= c.partSize {
				err = io.EOF
				n = bytesLeft
			}
		}

		var (
			reader  io.ReadSeeker
			cleanup func()
		)

		reader = io.NewSectionReader(r, u.readerPos, n)
		if u.cfg.BufferProvider != nil {
			reader, cleanup = u.cfg.BufferProvider.GetWriteTo(reader)
		} else {
			cleanup = func() {}
		}

		u.readerPos += n

		return reader, int(n), cleanup, err

	default:
		part, err := u.cfg.partPool.Get(u.ctx)
		if err != nil {
			return nil, 0, func() {}, err
		}
		var offset int
		var n int
		b := *part
		for offset < len(b) && err == nil {
			var n int
			n, err = r.Read(b[offset:])
			offset += n
		}
		u.readerPos += int64(n)

		cleanup := func() {
			u.cfg.partPool.Put(part)
		}

		return bytes.NewReader((*part)[0:n]), n, cleanup, err
	}
}
