package compress

import "io"

// mockWriteCloser is the mock structure of io.WriteCloser.
type mockWriteCloser struct {
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
}

var _ io.WriteCloser = (*mockWriteCloser)(nil)

func (m *mockWriteCloser) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

func (m *mockWriteCloser) Close() error {
	return m.CloseFunc()
}

// mockReadCloser is the mock structure of io.ReadCloser.
type mockReadCloser struct {
	ReadFunc  func(p []byte) (n int, err error)
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
}

var _ io.WriteCloser = (*mockReadCloser)(nil)

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

func (m *mockReadCloser) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

func (m *mockReadCloser) Close() error {
	return m.CloseFunc()
}
