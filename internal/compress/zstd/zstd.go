// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package zstd

import (
	"github.com/klauspost/compress/zstd"
	"github.com/vdaas/vald/internal/io"
)

type (
	// EOption is type alias of zstd.EOption.
	EOption = zstd.EOption

	// DOption is type alias of zstd.DOption.
	DOption = zstd.DOption
)

// Encoder represents an interface for Encoder of zstd.
type Encoder interface {
	io.WriteCloser
	ReadFrom(r io.Reader) (n int64, err error)
}

// Decoder represents an interface for Decoder of zstd.
type Decoder interface {
	io.Reader
	Close()
	WriteTo(w io.Writer) (int64, error)
}

// Zstd is an interface to create Writer and Reader implementation.
type Zstd interface {
	NewWriter(w io.Writer, opts ...EOption) (Encoder, error)
	NewReader(r io.Reader, opts ...DOption) (Decoder, error)
}

type compress struct{}

// New returns Zstd implementation.
func New() Zstd {
	return new(compress)
}

// NewWriter returns Encoder implementation.
func (*compress) NewWriter(w io.Writer, opts ...EOption) (Encoder, error) {
	return zstd.NewWriter(w, opts...)
}

// NewReader returns Decoder implementation.
func (*compress) NewReader(r io.Reader, opts ...DOption) (Decoder, error) {
	return zstd.NewReader(r, opts...)
}
