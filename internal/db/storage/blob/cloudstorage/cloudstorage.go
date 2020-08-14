package cloudstorage

import (
	"context"
	"io"
	"net/url"
	"reflect"

	iblob "github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/errors"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
)

type client struct {
	url *url.URL

	urlOpener *gcsblob.URLOpener
	bucket    *blob.Bucket

	readerOpts *blob.ReaderOptions
	writerOpts *blob.WriterOptions
}

// New returns blob.Bucket implementation.
func New(opts ...Option) (iblob.Bucket, error) {
	c := new(client)

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return c, nil
}

func (c *client) Open(ctx context.Context) (err error) {
	c.bucket, err = c.urlOpener.OpenBucketURL(ctx, c.url)
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
