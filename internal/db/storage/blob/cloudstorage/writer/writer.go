package writer

import (
	"context"
	"io"

	"github.com/vdaas/vald/internal/errors"
	"gocloud.dev/blob"
)

type writer struct {
	bucket *blob.Bucket
	key    string
	opts   *blob.WriterOptions
	*blob.Writer
}

// Writer is the interface that groups the basic Open and Write and Close methods.
type Writer interface {
	Open(ctx context.Context) error
	io.WriteCloser
}

// New returns Writer implementation.
func New(opts ...Option) (Writer, error) {
	w := new(writer)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(w); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}
	return w, nil
}

func (w *writer) Open(ctx context.Context) (err error) {
	w.Writer, err = w.bucket.NewWriter(ctx, w.key, w.opts)
	if err != nil {
		return errors.Wrap(err, "failed to create writer")
	}
	return nil
}

func (w *writer) Close() error {
	if w.Writer == nil {
		return errors.ErrStorageWriterNotOpened
	}

	err := w.Writer.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close")
	}
	return nil
}

func (w *writer) Write(p []byte) (n int, err error) {
	if w.Writer == nil {
		return 0, errors.ErrStorageWriterNotOpened
	}

	n, err = w.Writer.Write(p)
	if err != nil {
		return 0, errors.Wrap(err, "failed to write")
	}
	return n, nil
}
