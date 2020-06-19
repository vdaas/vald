package cloudstorage

import "gocloud.dev/blob"

// Option configures client of google cloud storage.
type Option func(*client) error

var (
	defaultOpts = []Option{}
)

// WithBeforeRead returns Option that sets c.readerOpts.BeforeRead.
func WithBeforeRead(fn func(asFunc func(interface{}) bool) error) Option {
	return func(c *client) error {
		if fn != nil {
			if c.readerOpts == nil {
				c.readerOpts = new(blob.ReaderOptions)
			}
			c.readerOpts.BeforeRead = fn
		}
		return nil
	}
}

// WithWriteBufferSize returns Option that sets c.writerOpts.BufferSize.
func WithWriteBufferSize(size int) Option {
	return func(c *client) error {
		if size > 0 {
			if c.writerOpts == nil {
				c.writerOpts = new(blob.WriterOptions)
			}
			c.writerOpts.BufferSize = size
		}
		return nil
	}
}

// WithWriteCacheControl returns Option that sets c.writerOpts.CacheControl.
func WithWriteCacheControl(str string) Option {
	return func(c *client) error {
		if len(str) != 0 {
			if c.writerOpts == nil {
				c.writerOpts = new(blob.WriterOptions)
			}
			c.writerOpts.CacheControl = str
		}
		return nil
	}
}

// WithWriteContentDisposition returns Option that sets c.writerOpts.ContentDisposition.
func WithWriteContentDisposition(str string) Option {
	return func(c *client) error {
		if len(str) != 0 {
			if c.writerOpts == nil {
				c.writerOpts = new(blob.WriterOptions)
			}
			c.writerOpts.ContentDisposition = str
		}
		return nil
	}
}

// WithWriteContentEncoding returns Option that sets c.writerOpts.Encoding.
func WithWriteContentEncoding(str string) Option {
	return func(c *client) error {
		if len(str) != 0 {
			if c.writerOpts == nil {
				c.writerOpts = new(blob.WriterOptions)
			}
			c.writerOpts.ContentEncoding = str
		}
		return nil
	}
}

// WithWriteContentLanguage returns Option that sets c.writerOpts.ContentLanguage.
func WithWriteContentLanguage(str string) Option {
	return func(c *client) error {
		if len(str) != 0 {
			if c.writerOpts == nil {
				c.writerOpts = new(blob.WriterOptions)
			}
			c.writerOpts.ContentLanguage = str
		}
		return nil
	}
}

// WithWriteContentType returns Option that sets c.writerOpts.ContentType.
func WithWriteContentType(str string) Option {
	return func(c *client) error {
		if len(str) != 0 {
			if c.writerOpts == nil {
				c.writerOpts = new(blob.WriterOptions)
			}
			c.writerOpts.ContentType = str
		}
		return nil
	}
}

// WithWriteContentMD5 returns Option that sets c.writerOpts.MD5.
func WithWriteContentMD5(b []byte) Option {
	return func(c *client) error {
		if len(b) != 0 {
			if c.writerOpts == nil {
				c.writerOpts = new(blob.WriterOptions)
			}
			c.writerOpts.ContentMD5 = b
		}
		return nil
	}
}

// WithWriteMetadata returns Option that sets c.writerOpts.Metadata.
func WithWriteMetadata(meta map[string]string) Option {
	return func(c *client) error {
		if meta != nil && len(meta) != 0 {
			if c.writerOpts == nil {
				c.writerOpts = new(blob.WriterOptions)
			}
			c.writerOpts.Metadata = meta
		}
		return nil
	}
}

// WithBeforeWrite returns Option that sets c.writeOpts.BeforeWrite.
func WithBeforeWrite(f func(asFunc func(interface{}) bool) error) Option {
	return func(c *client) error {
		if f != nil {
			if c.writerOpts == nil {
				c.writerOpts = new(blob.WriterOptions)
			}
			c.writerOpts.BeforeWrite = f
		}
		return nil
	}
}
