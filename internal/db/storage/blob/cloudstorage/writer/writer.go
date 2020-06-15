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
			return nil, err
		}
	}
	return w, nil
}

func (w *writer) Open(ctx context.Context) (err error) {
	w.Writer, err = w.bucket.NewWriter(ctx, w.key, w.opts)
	if err != nil {
		return err
	}
	return nil
}

func (w *writer) Close() error {
	if w.Writer == nil {
		return errors.ErrStorageWriterNotOpened
	}
	return w.Close()
}

func (w *writer) Write(p []byte) (n int, err error) {
	return 0, nil
}
