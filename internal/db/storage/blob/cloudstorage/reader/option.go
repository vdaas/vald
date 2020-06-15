package reader

import "gocloud.dev/blob"

// Option configures reader.
type Option func(r *reader) error

var (
	defaultOpts = []Option{}
)

// WithBucket returns Option that set r.bucket.
func WithBucket(bucket *blob.Bucket) Option {
	return func(r *reader) error {
		if bucket != nil {
			r.bucket = bucket
		}
		return nil
	}
}

// WithKey returns Option that set r.key.
func WithKey(key string) Option {
	return func(r *reader) error {
		if len(key) != 0 {
			r.key = key
		}
		return nil
	}
}

// WithBeforeRead returns Option that set r.opts.BeforeRead.
func WithBeforeRead(fn func(asFunc func(interface{}) bool) error) Option {
	return func(r *reader) error {
		if fn != nil {
			if r.opts == nil {
				r.opts = new(blob.ReaderOptions)
			}
			r.opts.BeforeRead = fn
		}
		return nil
	}
}
