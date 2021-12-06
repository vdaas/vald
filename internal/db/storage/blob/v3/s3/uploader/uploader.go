package uploader

import (
	"context"
	"io"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vdaas/vald/internal/db/storage/blob/v3/s3/file"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	iio "github.com/vdaas/vald/internal/io"
)

type Client interface {
	Open(ctx context.Context, key string) (err error)
	io.WriteCloser
}

type client struct {
	eg     errgroup.Group
	client manager.UploadAPIClient

	bucket      string
	contentType string
	concurrency int
	partSize    int64

	pr *io.PipeReader
	pw *io.PipeWriter
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

func (c *client) Open(ctx context.Context, key string) (err error) {
	f, err := file.CreateTemp()
	if err != nil {
		return err
	}
	c.pr, c.pw = io.Pipe()

	c.eg.Go(func() (err error) {
		defer f.Close()
		defer c.pr.Close()

		_, err = iio.Copy(f, c.pr)
		if err != nil {
			return err
		}

		// After the write operation is complete, shift Seek to read from the beginning of the file.
		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}

		err = c.upload(ctx, key, f)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (c *client) Close() (err error) {
	if c.pw != nil {
		// pw.Close only return EOF and nil and never fail, so ignore the error.
		_ = c.pw.Close()
	}

	if werr := c.eg.Wait(); werr != nil {
		return err
	}
	return nil
}

func (c *client) Write(data []byte) (n int, err error) {
	if c.pw == nil {
		return 0, errors.ErrStorageWriterNotOpened
	}
	return c.pw.Write(data)
}

func (c *client) upload(ctx context.Context, key string, body io.Reader) (err error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(c.contentType),
	}

	_, err = manager.NewUploader(c.client).Upload(ctx, input, func(u *manager.Uploader) {
		u.Concurrency = c.concurrency
		u.PartSize = c.partSize
	})
	if err != nil {
		return err
	}
	return nil
}
