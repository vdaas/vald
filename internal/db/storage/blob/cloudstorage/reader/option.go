package reader

import "gocloud.dev/blob"

// Option configures reader.
type Option func(r *reader) error

var (
	defaultOpts = []Option{}
)

// WithBucket returns Option that sets r.bucket.
func WithBucket(bucket *blob.Bucket) Option {
	return func(r *reader) error {
		if bucket != nil {
			r.bucket = bucket
		}
		return nil
	}
}

// WithKey returns Option that sets r.key.
func WithKey(key string) Option {
	return func(r *reader) error {
		if len(key) != 0 {
			r.key = key
		}
		return nil
	}
}

// WithReaderOptions returns Option that sets r.opts.
func WithReaderOptions(opts *blob.ReaderOptions) Option {
	return func(r *reader) error {
		if opts != nil {
			r.opts = opts
		}
		return nil
	}
}
