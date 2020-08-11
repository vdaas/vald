package compress

import "io"

type mockCompressor struct {
	CompressVectorFunc   func(vector []float32) (bytes []byte, err error)
	DecompressVectorFunc func(bytes []byte) (vector []float32, err error)
	ReaderFunc           func(src io.ReadCloser) (io.ReadCloser, error)
	WriterFunc           func(dst io.WriteCloser) (io.WriteCloser, error)
}

func (m *mockCompressor) CompressVector(vector []float32) (bytes []byte, err error) {
	return m.CompressVectorFunc(vector)
}

func (m *mockCompressor) DecompressVector(bytes []byte) (vector []float32, err error) {
	return m.DecompressVectorFunc(bytes)
}

func (m *mockCompressor) Reader(src io.ReadCloser) (io.ReadCloser, error) {
	return m.ReaderFunc(src)
}

func (m *mockCompressor) Writer(dst io.WriteCloser) (io.WriteCloser, error) {
	return m.WriterFunc(dst)
}

type mockReadCloser struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (m *mockReadCloser) Read(p []byte) (int, error) {
	return m.ReadFunc(p)
}

func (m *mockReadCloser) Close() error {
	return m.CloseFunc()
}

type mockWriteCloser struct {
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (m *mockWriteCloser) Write(p []byte) (int, error) {
	return m.WriteFunc(p)
}

func (m *mockWriteCloser) Close() error {
	return m.CloseFunc()
}
