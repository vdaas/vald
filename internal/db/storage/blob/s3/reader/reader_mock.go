package reader

// MockReadCloser represents mock of io.ReadCloser.
type MockReadCloser struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

// Read calls ReadFunc.
func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

// Close calls CloseFunc.
func (m *MockReadCloser) Close() error {
	return m.CloseFunc()
}
