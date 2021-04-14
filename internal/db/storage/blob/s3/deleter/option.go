package deleter

import (
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3"
	"github.com/vdaas/vald/internal/errors"
)

type Option func(d *deleter) error

var (
	defaultOpts = []Option{}
)

// WithService returns the option to set s for writer.
func WithService(s *s3.S3) Option {
	return func(d *deleter) error {
		if s == nil {
			return errors.NewErrInvalidOption("service", s)
		}
		d.service = s
		return nil
	}
}

// WithBucket returns the option to set bucket for writer.
func WithBucket(bucket string) Option {
	return func(d *deleter) error {
		if len(bucket) == 0 {
			return errors.NewErrInvalidOption("bucket", bucket)
		}
		d.bucket = bucket
		return nil
	}
}
