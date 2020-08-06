package compress

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
