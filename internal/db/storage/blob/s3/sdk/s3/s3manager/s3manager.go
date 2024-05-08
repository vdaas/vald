// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package s3manager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
)

type (
	// Uploader is type alias of s3manager.Uploader.
	Uploader = s3manager.Uploader
	// UploadInput is type alias of s3manager.UploadInput.
	UploadInput = s3manager.UploadInput
	// UploadOutput is type alias of s3manager.UploadOutput.
	UploadOutput = s3manager.UploadOutput
)

// UploadClient represents an interface to upload to s3.
type UploadClient interface {
	UploadWithContext(ctx aws.Context, input *UploadInput, opts ...func(*Uploader)) (*UploadOutput, error)
}

// S3Manager represents an interface to create object of s3manager package.
type S3Manager interface {
	NewUploaderWithClient(svc s3iface.S3API, options ...func(*Uploader)) UploadClient
}

type s3mngr struct{}

// New returns S3Manager implementation.
func New() S3Manager {
	return new(s3mngr)
}

// NewUploaderWithClient returns UploadClient implementation.
func (*s3mngr) NewUploaderWithClient(svc s3iface.S3API, options ...func(*Uploader)) UploadClient {
	return s3manager.NewUploaderWithClient(svc, options...)
}
