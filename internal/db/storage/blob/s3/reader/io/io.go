//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
package io

import (
	"context"
	"io"

	iio "github.com/vdaas/vald/internal/io"
)

// IO represents an interface to create object for io.
type IO interface {
	NewReaderWithContext(ctx context.Context, r io.Reader) (io.Reader, error)
	NewReadCloserWithContext(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error)
}

type ctxio struct{}

// New returns IO implementation.
func New() IO {
	return new(ctxio)
}

// NewReaderWithContext calls io.NewReaderWithContext.
func (*ctxio) NewReaderWithContext(ctx context.Context, r io.Reader) (io.Reader, error) {
	return iio.NewReaderWithContext(ctx, r)
}

// NewReadCloserWithContext calls io.NewReadCloserWithContext.
func (*ctxio) NewReadCloserWithContext(ctx context.Context, r io.ReadCloser) (io.ReadCloser, error) {
	return iio.NewReadCloserWithContext(ctx, r)
}
