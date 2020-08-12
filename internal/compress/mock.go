package compress

import "io"

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
