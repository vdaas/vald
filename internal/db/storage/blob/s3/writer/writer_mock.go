package writer

import "context"

// MockWriter represents Writer.
type MockWriter struct {
	OpenFunc  func(ctx context.Context) error
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
}

// Open calls OpenFunc.
func (m *MockWriter) Open(ctx context.Context) error {
	return m.OpenFunc(ctx)
}

// Write calls WriteFunc.
func (m *MockWriter) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

// Close calls CloseFunc.
func (m *MockWriter) Close() error {
	return m.CloseFunc()
}
