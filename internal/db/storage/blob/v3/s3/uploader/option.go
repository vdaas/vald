package uploader

import (
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for client.
type Option func(c *client) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithContentType("application/octet-stream"),
	WithMaxPartSize(12 * 1024 * 1024),
	WithConcurrency(manager.DefaultUploadConcurrency),
}

// WithErrGroup returns the option to set the eg.
func WithErrGroup(eg errgroup.Group) Option {
	return func(c *client) error {
		if eg == nil {
			return errors.NewErrInvalidOption("errgroup", eg)
		}
		c.eg = eg
		return nil
	}
}

// WithBucket returns the option to set bucket for writer.
func WithBucket(bucket string) Option {
	return func(c *client) error {
		if len(bucket) == 0 {
			return errors.NewErrInvalidOption("bucket", bucket)
		}
		c.bucket = bucket
		return nil
	}
}

// WithAPIClient returns the option to set client for writer.
func WithAPIClient(ac manager.UploadAPIClient) Option {
	return func(c *client) error {
		if ac == nil {
			return errors.NewErrInvalidOption("apiclient", ac)
		}
		c.client = ac
		return nil
	}
}

// WithContentType returns the option to set ct for writer.
func WithContentType(ct string) Option {
	return func(c *client) error {
		if len(ct) == 0 {
			return errors.NewErrInvalidOption("contentType", ct)
		}
		c.contentType = ct
		return nil
	}
}

// WithConcurrency retunrs the option to set the concurrency.
// The minimum number of goroutine is 5.
func WithConcurrency(n int) Option {
	return func(c *client) error {
		if n >= manager.DefaultUploadConcurrency {
			c.concurrency = n
		}
		return nil
	}
}

// WithMaxPartSize retunrs the option to set the partSize.
// The minimum allowed part size is 5MB.
func WithMaxPartSize(size int64) Option {
	return func(c *client) error {
		if size >= manager.DefaultUploadPartSize {
			c.partSize = size
		}
		return nil
	}
}
