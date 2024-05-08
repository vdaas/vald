// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package lz4

import (
	lz4 "github.com/pierrec/lz4/v3"
	"github.com/vdaas/vald/internal/io"
)

// Header is type alias of lz4.Header.
type Header = lz4.Header

// Reader represents an interface for Reader of lz4.
type Reader interface {
	io.Reader
}

// Writer represents an interface for Writer of lz4.
type Writer interface {
	io.WriteCloser
	Header() *Header
	Flush() error
}

type writer struct {
	*lz4.Writer
}

// Header returns lz4.Writer.Header object.
func (w *writer) Header() *Header {
	return &w.Writer.Header
}

// LZ4 is an interface to create Writer and Reader implementation.
type LZ4 interface {
	NewWriterLevel(w io.Writer, level int) Writer
	NewWriter(w io.Writer) Writer
	NewReader(r io.Reader) Reader
}

type compress struct{}

// New returns LZ4 implementation.
func New() LZ4 {
	return new(compress)
}

// NewWriterLevel returns Writer implementation.
func (*compress) NewWriterLevel(w io.Writer, level int) Writer {
	lw := lz4.NewWriter(w)
	lw.Header.CompressionLevel = level
	return &writer{
		Writer: lw,
	}
}

// NewWriter returns Writer implementation.
func (*compress) NewWriter(w io.Writer) Writer {
	return &writer{
		Writer: lz4.NewWriter(w),
	}
}

// NewReader returns Reader implementation.
func (*compress) NewReader(r io.Reader) Reader {
	return lz4.NewReader(r)
}
