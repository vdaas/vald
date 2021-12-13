// Package downloader provides download functions for s3.
package downloader

import (
	"context"
	"io"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/vdaas/vald/internal/db/storage/blob/v3/s3/file"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

var (
	errBlobNoSuchBucket = new(types.NoSuchBucket)
	errBlobNoSuchKey    = new(types.NoSuchKey)
)

type Client interface {
	Download(ctx context.Context, key string) (rc io.ReadCloser, err error)
}

type client struct {
	client manager.DownloadAPIClient

	bucket      string
	concurrency int
	partSize    int64
}

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return c, nil
}

func (c *client) Download(ctx context.Context, key string) (rc io.ReadCloser, err error) {
	f, err := file.CreateTemp()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if cerr := f.Close(); cerr != nil {
				err = errors.Wrap(err, cerr.Error())
			}
		}
	}()

	err = c.download(ctx, key, f)
	if err != nil {
		if errors.As(err, &errBlobNoSuchBucket) ||
			errors.As(err, &errBlobNoSuchKey) {
			log.Warn(err)
			return f, nil
		}
		return nil, err
	}

	// After the write operation is complete, shift Seek to read from the beginning of the file.
	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c *client) download(ctx context.Context, key string, w io.WriterAt) (err error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	_, err = manager.NewDownloader(c.client).Download(ctx, w, input, func(d *manager.Downloader) {
		d.Concurrency = c.concurrency
		d.PartSize = c.partSize
	})
	if err != nil {
		return err
	}
	return nil
}
