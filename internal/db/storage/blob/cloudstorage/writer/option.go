package writer

import "gocloud.dev/blob"

// Option configures *writer.
type Option func(w *writer) error

var (
	defaultOpts = []Option{}
)

// WithKey returns Option that sets w.key.
func WithKey(key string) Option {
	return func(w *writer) error {
		if len(key) != 0 {
			w.key = key
		}
		return nil
	}
}

// WithBucket returns Option that sets w.bucket.
func WithBucket(bucket *blob.Bucket) Option {
	return func(w *writer) error {
		if bucket != nil {
			w.bucket = bucket
		}
		return nil
	}
}

// WithWriterOptions returns Option that sets w.opts.
func WithWriterOptions(opts *blob.WriterOptions) Option {
	return func(w *writer) error {
		if opts != nil {
			w.opts = opts
		}
		return nil
	}
}
