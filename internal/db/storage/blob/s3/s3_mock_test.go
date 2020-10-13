package s3

import (
	"context"

	"github.com/vdaas/vald/internal/db/storage/blob/s3/reader"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/writer"
)

// MockS3IO represents mock of S3IO.
type MockS3IO struct {
	NewReaderFunc func(opts ...reader.Option) reader.Reader
	NewWriterFunc func(opts ...writer.Option) writer.Writer
}

// NewReader calls NewReaderFunc.
func (m *MockS3IO) NewReader(opts ...reader.Option) reader.Reader {
	return m.NewReaderFunc(opts...)
}

// NewWriter calls NewWriterFunc.
func (m *MockS3IO) NewWriter(opts ...writer.Option) writer.Writer {
	return m.NewWriterFunc(opts...)
}

// MockReader represents mock for s3.Reader.
type MockReader struct {
	OpenFunc  func(ctx context.Context) error
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

// Open calls OpenFunc.
func (m *MockReader) Open(ctx context.Context) error {
	return m.OpenFunc(ctx)
}

// Read calls ReadFunc.
func (m *MockReader) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

// Close calls CloseFunc.
func (m *MockReader) Close() error {
	return m.CloseFunc()
}

// MockWriter represents mock of s3.Writer.
type MockWriter struct {
	OpenFunc  func(ctx context.Context) error
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (m *MockWriter) Open(ctx context.Context) error {
	return m.OpenFunc(ctx)
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

func (m *MockWriter) Close() error {
	return m.CloseFunc()
}
