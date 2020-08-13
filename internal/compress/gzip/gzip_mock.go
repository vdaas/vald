package gzip

import "io"

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

// MockBuilder represents mock struct of Builder.
type MockBuilder struct {
	NewWriterLevelFunc func(w io.Writer, level int) (Writer, error)
	NewReaderFunc      func(r io.Reader) (Reader, error)
}

// NewWriterLevel calls NewWriterLevelFunc.
func (m *MockBuilder) NewWriterLevel(w io.Writer, level int) (Writer, error) {
	return m.NewWriterLevelFunc(w, level)
}

// NewReader calls NewReaderFunc.
func (m *MockBuilder) NewReader(r io.Reader) (Reader, error) {
	return m.NewReaderFunc(r)
}
