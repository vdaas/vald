//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
package zstd

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

// EOption is type alias of zstd.EOption.
type EOption = zstd.EOption

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

// Transcoder is an interface to create Writer and Reader implementation.
type Transcoder interface {
	NewWriter(w io.Writer, opts ...zstd.EOption) (Encoder, error)
	NewReader(r io.Reader, opts ...zstd.DOption) (Decoder, error)
}

type transcoder struct{}

// New returns Transcoder implementation.
func New() Transcoder {
	return new(transcoder)
}

// NewWriter returns Encoder implementation.
func (*transcoder) NewWriter(w io.Writer, opts ...zstd.EOption) (Encoder, error) {
	return zstd.NewWriter(w, opts...)
}

// NewReader returns Decoder implementation.
func (*transcoder) NewReader(r io.Reader, opts ...zstd.DOption) (Decoder, error) {
	return zstd.NewReader(r, opts...)
}
