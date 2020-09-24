package reader

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
)

// MockS3API represents mock for s3iface.MMockS3API.
type MockS3API struct {
	s3iface.S3API
	GetObjectWithContextFunc func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error)
}

// GetObjectWithContext calls GetObjectWithContextFunc.
func (m *MockS3API) GetObjectWithContext(ctx aws.Context, in *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
	return m.GetObjectWithContextFunc(ctx, in, opts...)
}

// MockIO represents mock for io.IO
type MockIO struct {
	NewReaderWithContextFunc     func(ctx context.Context, r io.Reader) (io.Reader, error)
	NewReadCloserWithContextFunc func(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error)
}

// NewReaderWithContext calls NewReaderWithContextFunc.
func (m *MockIO) NewReaderWithContext(ctx context.Context, r io.Reader) (io.Reader, error) {
	return m.NewReaderWithContextFunc(ctx, r)
}

// NewReadCloserWithContext calls NewReadCloserWithContextFunc.
func (m *MockIO) NewReadCloserWithContext(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
	return m.NewReadCloserWithContextFunc(ctx, r)
}

// MockReadCloser represents mock for io.ReadCloser.
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
