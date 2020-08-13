package lz4

import "io"

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
	return m.Close()
}

// Header calls HeaderFunc.
func (m *MockWriter) Header() *Header {
	return m.HeaderFunc()
}

// Flush calls FlushFunc.
func (m *MockWriter) Flush() error {
	return m.FlushFunc()
}

// MockBuilder represents mock struct of Builder.
type MockBuilder struct {
	NewWriterFunc func(w io.Writer) Writer
	NewReaderFunc func(r io.Reader) Reader
}

// NewWriter calls NewWriterFunc.
func (m *MockBuilder) NewWriter(w io.Writer) Writer {
	return m.NewWriterFunc(w)
}

// NewReader calls NewReader.
func (m *MockBuilder) NewReader(r io.Reader) Reader {
	return m.NewReaderFunc(r)
}
