package reader

import (
	"context"
	"io"

	"github.com/vdaas/vald/internal/errors"
	"gocloud.dev/blob"
)

type reader struct {
	key string

	bucket *blob.Bucket
	opts   *blob.ReaderOptions
	*blob.Reader
}

// Reader is the interface that groups the basic Open and Read and Close methods.
type Reader interface {
	Open(ctx context.Context) error
	io.ReadCloser
}

// New returns Reader implementation.
func New(opts ...Option) (Reader, error) {
	r := new(reader)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}
	return r, nil
}

func (r *reader) Open(ctx context.Context) (err error) {
	r.Reader, err = r.bucket.NewReader(ctx, r.key, r.opts)
	if err != nil {
		return errors.Wrap(err, "failed to create reader")
	}
	return nil
}

func (r *reader) Close() error {
	if r.Reader == nil {
		return errors.ErrStorageReaderNotOpened
	}
	return r.Reader.Close()
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.Reader == nil {
		return 0, errors.ErrStorageReaderNotOpened
	}
	return r.Reader.Read(p)
}
