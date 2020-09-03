package s3manager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
)

type (
	// Uploade is type alias of s3manager.Uploader.
	Uploade = s3manager.Uploader
	// UploadInput is type alias of s3manager.UploadeInput.
	UploadInput = s3manager.UploadInput
	// UploadOutput is type alias of s3manager.UploadeOutput.
	UploadOutput = s3manager.UploadOutput
)

// Uploader represents an interface to upload to s3.
type Uploader interface {
	UploadWithContext(ctx aws.Context, input *UploadInput, opts ...func(*Uploade)) (*UploadOutput, error)
}

// S3Manager represents an interface to create object of s3manager package.
type S3Manager interface {
	NewUploaderWithClient(svc s3iface.S3API, options ...func(*Uploade)) Uploader
}

type s3mngr struct{}

// New returns S3Manager implementation.
func New() S3Manager {
	return new(s3mngr)
}

func (*s3mngr) NewUploaderWithClient(svc s3iface.S3API, options ...func(*Uploade)) Uploader {
	return s3manager.NewUploaderWithClient(svc, options...)
}
