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
