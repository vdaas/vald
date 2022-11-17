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
package compress

import "github.com/vdaas/vald/internal/io"

// MockCompressor represents mock of Compressor.
type MockCompressor struct {
	CompressVectorFunc   func(vector []float32) (bytes []byte, err error)
	DecompressVectorFunc func(bytes []byte) (vector []float32, err error)
	ReaderFunc           func(src io.ReadCloser) (io.ReadCloser, error)
	WriterFunc           func(dst io.WriteCloser) (io.WriteCloser, error)
}

// CompressVector calls CompressVectorFunc.
func (m *MockCompressor) CompressVector(vector []float32) (bytes []byte, err error) {
	return m.CompressVectorFunc(vector)
}

// DecompressVector calls DecompressVectorFunc.
func (m *MockCompressor) DecompressVector(bytes []byte) (vector []float32, err error) {
	return m.DecompressVectorFunc(bytes)
}

// Reader calls ReaderFunc.
func (m *MockCompressor) Reader(src io.ReadCloser) (io.ReadCloser, error) {
	return m.ReaderFunc(src)
}

// Writer calls WriterFunc.
func (m *MockCompressor) Writer(dst io.WriteCloser) (io.WriteCloser, error) {
	return m.WriterFunc(dst)
}

// MockReadCloser represents mock of ReadCloser.
type MockReadCloser struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

// Read calls ReadFunc.
func (m *MockReadCloser) Read(p []byte) (int, error) {
	return m.ReadFunc(p)
}

// Close calls CloseFunc.
func (m *MockReadCloser) Close() error {
	return m.CloseFunc()
}

// MockWriteCloser represents mock of WriteCloser.
type MockWriteCloser struct {
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
}

// Write calls WriterFunc.
func (m *MockWriteCloser) Write(p []byte) (int, error) {
	return m.WriteFunc(p)
}

// Close calls CloseFunc.
func (m *MockWriteCloser) Close() error {
	return m.CloseFunc()
}
