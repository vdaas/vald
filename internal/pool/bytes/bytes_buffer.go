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

// Package bytes provides buffer pool functionality for pooling bytes buffer
package bytes

import (
	"bytes"

	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/pool"
)

type Buffer interface {
	pool.Extender
	pool.Flusher
	Bytes() []byte
	String() string
	Len() int
	Cap() int
	Truncate(n int)
	Reset()
	Grow(n int)
	io.Writer
	io.WriterTo
	io.StringWriter
	io.ByteWriter
	io.RuneReader
	io.ByteReader
	io.RuneScanner
	io.ByteScanner
	io.Reader
	io.ReaderFrom
	WriteRune(r rune) (n int, err error)
	Next(n int) []byte
	ReadBytes(delim byte) (line []byte, err error)
	ReadString(delim byte) (line string, err error)
}

type Pool interface {
	Get() Buffer
	Put(buf Buffer)
}

type bytesBuffer struct {
	*bytes.Buffer
}

type bufferPool struct {
	p pool.Buffer
}

func New(poolSize, bufferSize int) Pool {
	return &bufferPool{
		p: pool.New(
			pool.WithConstructor(func(size uint64) interface{} {
				return &bytesBuffer{
					Buffer: bytes.NewBuffer(make([]byte, 0, size)),
				}
			}),
			pool.WithLimit(uint64(poolSize)),
			pool.WithSize(uint64(bufferSize)),
		),
	}
}

func (b *bufferPool) Get() Buffer {
	buf, ok := b.p.Get().(Buffer)
	if ok {
		return buf
	}
	return nil
}

func (b *bufferPool) Put(buf Buffer) {
	b.p.PutWithResize(buf, uint64(buf.Cap()))
}

func (b *bytesBuffer) Extend(size uint64) (data interface{}) {
	b.Grow(int(size))
	return b
}

func (b *bytesBuffer) Flush() (data interface{}) {
	b.Reset()
	return
}
