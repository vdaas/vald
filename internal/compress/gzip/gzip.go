//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package gzip

import (
	"github.com/klauspost/compress/gzip"
	"github.com/vdaas/vald/internal/io"
)

// These constants are copied from the gzip package.
const (
	NoCompression       = gzip.NoCompression
	BestSpeed           = gzip.BestSpeed
	BestCompression     = gzip.BestCompression
	DefaultCompression  = gzip.DefaultCompression
	ConstantCompression = gzip.ConstantCompression
	HuffmanOnly         = gzip.HuffmanOnly
)

// Reader represents an interface for Reader of gzip.
type Reader interface {
	io.ReadCloser
	Reset(r io.Reader) error
	Multistream(ok bool)
}

// Writer represents an interface for Writer of gzip.
type Writer interface {
	io.WriteCloser
	Reset(w io.Writer)
	Flush() error
}

// Gzip is an interface to create Writer and Reader implementation.
type Gzip interface {
	NewReader(r io.Reader) (Reader, error)
	NewWriterLevel(w io.Writer, level int) (Writer, error)
}

type compress struct{}

// New returns Gzip implementation.
func New() Gzip {
	return new(compress)
}

// NewWriterLevel returns Writer implementation.
func (*compress) NewWriterLevel(w io.Writer, level int) (Writer, error) {
	return gzip.NewWriterLevel(w, level)
}

// NewReader returns Reader implementation.
func (*compress) NewReader(r io.Reader) (Reader, error) {
	return gzip.NewReader(r)
}
