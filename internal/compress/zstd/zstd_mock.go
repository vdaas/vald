package zstd

import (
	"io"

	"github.com/klauspost/compress/zstd"
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

// MockDecoder represents
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

// MockTranscoder represents mock struct of Transcoder.
type MockTranscoder struct {
	NewWriterFunc func(w io.Writer, opts ...zstd.EOption) (Encoder, error)
	NewReaderFunc func(r io.Reader, opts ...zstd.DOption) (Decoder, error)
}

// NewWriter calls NewWriterFunc.
func (m *MockTranscoder) NewWriter(w io.Writer, opts ...zstd.EOption) (Encoder, error) {
	return m.NewWriterFunc(w, opts...)
}

// NewReader calls NewReader.
func (m *MockTranscoder) NewReader(r io.Reader, opts ...zstd.DOption) (Decoder, error) {
	return m.NewReaderFunc(r, opts...)
}
