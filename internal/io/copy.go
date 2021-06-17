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

// Package io provides io functions
package io

import (
	"bytes"
	"context"
	"math"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/pool"
)

var cio = NewCopier(0)

func Copy(dst Writer, src Reader) (written int64, err error) {
	return cio.Copy(dst, src)
}

func CopyWithContext(ctx context.Context, dst Writer, src Reader) (written int64, err error) {
	return cio.CopyWithContext(ctx, dst, src)
}

type Copier interface {
	Copy(dst Writer, src Reader) (written int64, err error)
	CopyWithContext(ctx context.Context, dst Writer, src Reader) (written int64, err error)
}

type copier struct {
	buffer pool.Buffer
}

const (
	defaultBufferSize int = 64 * 1024
)

func NewCopier(size int) Copier {
	if size <= 0 {
		size = defaultBufferSize
	}
	return &copier{
		buffer: pool.New(context.TODO(),
			pool.WithSize(uint64(size)),
		),
	}
}

func (c *copier) CopyWithContext(ctx context.Context, dst Writer, src Reader) (written int64, err error) {
	csrc, err := NewReaderWithContext(ctx, src)
	if err != nil {
		return 0, err
	}
	return c.Copy(dst, csrc)
}

func (c *copier) Copy(dst Writer, src Reader) (written int64, err error) {
	return c.copyContext(context.TODO(), dst, src)
}

func (c *copier) copyContext(ctx context.Context, dst Writer, src Reader) (written int64, err error) {
	var (
		wt WriterTo
		rf ReaderFrom
		ok bool
	)
	if wt, ok = src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rf, ok = dst.(ReaderFrom); ok {
		return rf.ReadFrom(src)
	}

	var (
		limit int64 = math.MaxInt64
		size        = int64(c.buffer.Size(ctx))
		l     *LimitedReader
		buf   *bytes.Buffer
	)
	if l, ok = src.(*LimitedReader); ok && l.N >= 1 && size > l.N {
		limit = l.N
		size = limit
	}
	buf, ok = c.buffer.Get(ctx).(*bytes.Buffer)
	if !ok || buf == nil {
		buf = bytes.NewBuffer(make([]byte, size))
	}
	if size > int64(buf.Cap()) {
		size = int64(buf.Cap())
	}
	defer c.buffer.PutWithResize(ctx, buf, uint64(size))
	var nr, nw int
	for err != EOF {
		nr, err = src.Read(buf.Bytes()[:size])
		if nr > 0 {
			if int64(nr) > size {
				if int64(nr) >= limit {
					size = limit
				} else {
					size = int64(nr)
				}
			}
			nw, err = dst.Write(buf.Bytes()[0:nr])
			if nw < 0 || nr < int(nw) {
				if err == nil {
					return written, errors.ErrInvalidWriteResult
				}
				nw = 0
			}
			written += int64(nw)
			if err != nil {
				return written, err
			}
			if nr != int(nw) {
				return written, errors.ErrShortWrite
			}
		}
		if err != nil && err != EOF {
			return written, err
		}
	}
	return written, nil
}
