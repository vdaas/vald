//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"io"
	"math"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

var cio = NewCopier(0)

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return cio.Copy(dst, src)
}

type Copier interface {
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
}

type copier struct {
	bufSize int64
	pool    sync.Pool
}

const (
	defaultBufferSize int = 64 * 1024
)

func NewCopier(size int) Copier {
	c := new(copier)
	if size > 0 {
		atomic.StoreInt64(&c.bufSize, int64(size))
	} else {
		atomic.StoreInt64(&c.bufSize, int64(defaultBufferSize))
	}
	c.pool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, int(atomic.LoadInt64(&c.bufSize))))
		},
	}
	return c
}

func (c *copier) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	if dst == nil || src == nil {
		return 0, errors.New("empty source or destination")
	}
	var (
		wt io.WriterTo
		rf io.ReaderFrom
		ok bool
	)
	if wt, ok = src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rf, ok = dst.(io.ReaderFrom); ok {
		return rf.ReadFrom(src)
	}

	var (
		limit int64 = math.MaxInt64
		size  int64 = atomic.LoadInt64(&c.bufSize)
		l     *io.LimitedReader
		buf   *bytes.Buffer
	)
	if l, ok = src.(*io.LimitedReader); ok && l.N >= 1 && size > l.N {
		limit = l.N
		size = limit
	}
	buf, ok = c.pool.Get().(*bytes.Buffer)
	if !ok || buf == nil {
		buf = bytes.NewBuffer(make([]byte, size))
	}
	defer func() {
		if atomic.LoadInt64(&c.bufSize) < size {
			atomic.StoreInt64(&c.bufSize, size)
			buf.Grow(int(size))
		}
		buf.Reset()
		c.pool.Put(buf)
	}()
	if size > int64(buf.Cap()) {
		size = int64(buf.Cap())
	}
	var nr, nw int
	for err != io.EOF {
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
					return written, errors.New("invalid write result")
				}
				nw = 0
			}
			written += int64(nw)
			if err != nil {
				return written, err
			}
			if nr != int(nw) {
				return written, io.ErrShortWrite
			}
		}
		if err != nil && err != io.EOF {
			return written, err
		}
	}
	return written, nil
}
