package cloudstorage

import (
	"context"
	"io"

	iblob "github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/cloudstorage/reader"
	"gocloud.dev/blob"
)

type client struct {
	urlstr string

	bucket *blob.Bucket
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
		// TODO: return err
		return nil
	}
	return c.bucket.Close()
}

func (c *client) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	if c.bucket == nil {
		// TODO: return nil, error
		return nil, nil
	}

	r, err := reader.New(
		reader.WithBucket(c.bucket),
		reader.WithKey(key),
	)
	if err != nil {
		return nil, err
	}

	return r, r.Open(ctx)
}

func (c *client) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	return nil, nil
}
