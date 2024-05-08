//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package reader provides the reader functions for handling with s3.
// This package is wrapping package of "https://github.com/aws/aws-sdk-go".
package reader

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/io"
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

// MockIO represents mock for io.IO.
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

// MockReader represents mock for Reader.
type MockReader struct {
	OpenFunc  func(ctx context.Context, key string) error
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

// Open calls OpenFunc.
func (m *MockReader) Open(ctx context.Context, key string) error {
	return m.OpenFunc(ctx, key)
}

// Read calls ReadFunc.
func (m *MockReader) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

// Close calls CloseFunc.
func (m *MockReader) Close() error {
	return m.CloseFunc()
}
