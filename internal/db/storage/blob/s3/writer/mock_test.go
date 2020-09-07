package writer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3manager"
)

// S3Manager represents mock of s3manager.S3Manager.
type MockS3Manager struct {
	NewUploaderWithClientFunc func(s3iface.S3API, ...func(*s3manager.Uploader)) s3manager.UploadClient
}

// NewUploaderWithClient calls NewUNewUploaderWithClientFunc.
func (m *MockS3Manager) NewUploaderWithClient(svc s3iface.S3API, opts ...func(*s3manager.Uploader)) s3manager.UploadClient {
	return m.NewUploaderWithClientFunc(svc, opts...)
}

type MockUploadClient struct {
	UploadWithContextFunc func(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

func (m *MockUploadClient) UploadWithContext(ctx aws.Context, input *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return m.UploadWithContextFunc(ctx, input, opts...)
}

// MockWriteCloser represents mock of io.WriteCloser.
type MockWriteCloser struct {
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (m *MockWriteCloser) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

func (m *MockWriteCloser) Close() error {
	return m.CloseFunc()
}
