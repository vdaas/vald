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
package gzip

import "github.com/vdaas/vald/internal/io"

// MockReader represents mock of Reader.
type MockReader struct {
	ReadFunc        func(p []byte) (n int, err error)
	CloseFunc       func() error
	ResetFunc       func(r io.Reader) error
	MultistreamFunc func(ok bool)
}

// Read calls ReadFunc.
func (m *MockReader) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

// Close calls CloseFunc.
func (m *MockReader) Close() error {
	return m.CloseFunc()
}

// Reset calls ResetFunc.
func (m *MockReader) Reset(r io.Reader) error {
	return m.ResetFunc(r)
}

// Multistream calls MultistreamFunc.
func (m *MockReader) Multistream(ok bool) {
	m.MultistreamFunc(ok)
}

// MockWriter represents mock of Writer.
type MockWriter struct {
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
	ResetFunc func(w io.Writer)
	FlushFunc func() error
}

// Write calls WriteFunc.
func (m *MockWriter) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

// Close calls CloseFunc.
func (m *MockWriter) Close() error {
	return m.CloseFunc()
}

// Reset calls ResetFunc.
func (m *MockWriter) Reset(w io.Writer) {
	m.ResetFunc(w)
}

// Flush calls FlushFunc.
func (m *MockWriter) Flush() error {
	return m.FlushFunc()
}

// MockGzip represents mock struct of Gzip.
type MockGzip struct {
	NewWriterLevelFunc func(w io.Writer, level int) (Writer, error)
	NewReaderFunc      func(r io.Reader) (Reader, error)
}

// NewWriterLevel calls NewWriterLevelFunc.
func (m *MockGzip) NewWriterLevel(w io.Writer, level int) (Writer, error) {
	return m.NewWriterLevelFunc(w, level)
}

// NewReader calls NewReaderFunc.
func (m *MockGzip) NewReader(r io.Reader) (Reader, error) {
	return m.NewReaderFunc(r)
}
