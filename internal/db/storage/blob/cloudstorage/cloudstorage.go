package cloudstorage

import (
	"context"
	"io"

	iblob "github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/errors"
	"gocloud.dev/blob"
)

type client struct {
	urlstr string

	bucket *blob.Bucket

	readerOpts *blob.ReaderOptions
	writerOpts *blob.WriterOptions
}

// New returns blob.Bucket implementation.
func New(opts ...Option) (iblob.Bucket, error) {
	c := new(client)

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *client) Open(ctx context.Context) (err error) {
	c.bucket, err = blob.OpenBucket(ctx, c.urlstr)
	if err != nil {
		return err
	}
	return
}
func (c *client) Close() error {
	if c.bucket == nil {
		return errors.ErrBucketNotOpened
	}
	return c.bucket.Close()
}

func (c *client) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	if c.bucket == nil {
		return nil, errors.ErrBucketNotOpened
	}
	return c.bucket.NewReader(ctx, key, c.readerOpts)
}

func (c *client) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	if c.bucket == nil {
		return nil, errors.ErrBucketNotOpened
	}
	return c.bucket.NewWriter(ctx, key, c.writerOpts)
}
