// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package lz4

import "github.com/vdaas/vald/internal/io"

// MockReader represents mock struct of Reader.
type MockReader struct {
	ReadFunc func(p []byte) (n int, err error)
}

// Read calls ReadFunc.
func (m *MockReader) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

// MockWriter represents mock struct of Writer.
type MockWriter struct {
	WriteFunc  func(p []byte) (n int, err error)
	CloseFunc  func() error
	HeaderFunc func() *Header
	FlushFunc  func() error
}

// Write calls WriteFunc.
func (m *MockWriter) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

// Close calls CloseFunc.
func (m *MockWriter) Close() error {
	return m.CloseFunc()
}

// Header calls HeaderFunc.
func (m *MockWriter) Header() *Header {
	return m.HeaderFunc()
}

// Flush calls FlushFunc.
func (m *MockWriter) Flush() error {
	return m.FlushFunc()
}

// MockLZ4 represents mock struct of LZ4.
type MockLZ4 struct {
	NewWriterFunc      func(w io.Writer) Writer
	NewWriterLevelFunc func(w io.Writer, level int) Writer
	NewReaderFunc      func(r io.Reader) Reader
}

// NewWriter calls NewWriterFunc.
func (m *MockLZ4) NewWriter(w io.Writer) Writer {
	return m.NewWriterFunc(w)
}

// NewWriterLevel calls NewWriterLevelFunc.
func (m *MockLZ4) NewWriterLevel(w io.Writer, level int) Writer {
	return m.NewWriterLevelFunc(w, level)
}

// NewReader calls NewReader.
func (m *MockLZ4) NewReader(r io.Reader) Reader {
	return m.NewReaderFunc(r)
}
