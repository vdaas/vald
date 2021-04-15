package copier

import (
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3"
	"github.com/vdaas/vald/internal/errors"
)

type Option func(c *copier) error

var (
	defaultOpts = []Option{}
)

// WithService returns the option to set s for copier.
func WithService(s *s3.S3) Option {
	return func(c *copier) error {
		if s == nil {
			return errors.NewErrInvalidOption("service", s)
		}
		c.service = s
		return nil
	}
}

// WithBucket returns the option to set bucket for copier.
func WithBucket(bucket string) Option {
	return func(c *copier) error {
		if len(bucket) == 0 {
			return errors.NewErrInvalidOption("bucket", bucket)
		}
		c.bucket = bucket
		return nil
	}
}
