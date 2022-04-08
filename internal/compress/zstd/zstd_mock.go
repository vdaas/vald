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
package zstd

import (
	"github.com/klauspost/compress/zstd"
	"github.com/vdaas/vald/internal/io"
)

// MockEncoder represents mock of Encoder.
type MockEncoder struct {
	WriteFunc    func(p []byte) (n int, err error)
	CloseFunc    func() error
	ReadFromFunc func(r io.Reader) (n int64, err error)
}

// Write calls WriteFunc.
func (m *MockEncoder) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

// Close calls CloseFunc.
func (m *MockEncoder) Close() error {
	return m.CloseFunc()
}

// ReadFrom calls ReadFromFunc.
func (m *MockEncoder) ReadFrom(r io.Reader) (n int64, err error) {
	return m.ReadFromFunc(r)
}

// MockDecoder represents.
type MockDecoder struct {
	CloseFunc   func()
	ReadFunc    func(p []byte) (int, error)
	WriteToFunc func(w io.Writer) (int64, error)
}

// Close calls CloseFunc.
func (m *MockDecoder) Close() {
	m.CloseFunc()
}

// Read calls ReadFunc.
func (m *MockDecoder) Read(p []byte) (int, error) {
	return m.ReadFunc(p)
}

// WriteTo calls WriteToFunc.
func (m *MockDecoder) WriteTo(w io.Writer) (int64, error) {
	return m.WriteToFunc(w)
}

// MockZstd represents mock struct of Zstd.
type MockZstd struct {
	NewWriterFunc func(w io.Writer, opts ...zstd.EOption) (Encoder, error)
	NewReaderFunc func(r io.Reader, opts ...zstd.DOption) (Decoder, error)
}

// NewWriter calls NewWriterFunc.
func (m *MockZstd) NewWriter(w io.Writer, opts ...zstd.EOption) (Encoder, error) {
	return m.NewWriterFunc(w, opts...)
}

// NewReader calls NewReader.
func (m *MockZstd) NewReader(r io.Reader, opts ...zstd.DOption) (Decoder, error) {
	return m.NewReaderFunc(r, opts...)
}
