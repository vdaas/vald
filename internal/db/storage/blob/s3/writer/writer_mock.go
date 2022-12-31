// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package writer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3manager"
)

// MockWriter represents Writer.
type MockWriter struct {
	OpenFunc  func(ctx context.Context, key string) error
	WriteFunc func(p []byte) (n int, err error)
	CloseFunc func() error
}

// Open calls OpenFunc.
func (m *MockWriter) Open(ctx context.Context, key string) error {
	return m.OpenFunc(ctx, key)
}

// Write calls WriteFunc.
func (m *MockWriter) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

// Close calls CloseFunc.
func (m *MockWriter) Close() error {
	return m.CloseFunc()
}

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
