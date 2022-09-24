package downloader

import (
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for client.
type Option func(c *client) error

var defaultOptions = []Option{
	WithMaxPartSize(12 * 1024 * 1024),
	WithConcurrency(manager.DefaultDownloadConcurrency),
}

// WithBucket returns the option to set bucket for reader.
func WithBucket(bucket string) Option {
	return func(c *client) error {
		if len(bucket) == 0 {
			return errors.NewErrInvalidOption("bucket", bucket)
		}
		c.bucket = bucket
		return nil
	}
}

// WithAPIClient returns the option to set client for reader.
func WithAPIClient(ac manager.DownloadAPIClient) Option {
	return func(c *client) error {
		if ac == nil {
			return errors.NewErrInvalidOption("apiclient", ac)
		}
		c.client = ac
		return nil
	}
}

// WithConcurrency retunrs the option to set the concurrency.
// The minimum number of goroutine is 5.
func WithConcurrency(n int) Option {
	return func(c *client) error {
		if n >= manager.DefaultDownloadConcurrency {
			c.concurrency = n
		}
		return nil
	}
}

// WithMaxPartSize retunrs the option to set the partSize.
// The minimum allowed part size is 5MB.
func WithMaxPartSize(size int64) Option {
	return func(c *client) error {
		if size >= manager.DefaultDownloadPartSize {
			c.partSize = size
		}
		return nil
	}
}
