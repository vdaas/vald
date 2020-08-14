package gzip

import (
	"io"

	"github.com/klauspost/compress/gzip"
)

// These constants are copied from the gzip package.
const (
	NoCompression       = gzip.NoCompression
	BestSpeed           = gzip.BestSpeed
	BestCompression     = gzip.BestCompression
	DefaultCompression  = gzip.DefaultCompression
	ConstantCompression = gzip.ConstantCompression
	HuffmanOnly         = gzip.HuffmanOnly
)

// Reader represents an interface for Reader of gzip.
type Reader interface {
	io.ReadCloser
	Reset(r io.Reader) error
	Multistream(ok bool)
}

// Writer represents an interface for Writer of gzip.
type Writer interface {
	io.WriteCloser
	Reset(w io.Writer)
	Flush() error
}

// Io is an interface to create Writer and Reader implementation.
type Io interface {
	NewReader(r io.Reader) (Reader, error)
	NewWriterLevel(w io.Writer, level int) (Writer, error)
}

type gzipIo struct{}

// New returns Io implementation.
func New() Io {
	return new(gzipIo)
}

// NewWriterLevel returns Writer implementation.
func (*gzipIo) NewWriterLevel(w io.Writer, level int) (Writer, error) {
	return gzip.NewWriterLevel(w, level)
}

// NewReader returns Reader implementation.
func (*gzipIo) NewReader(r io.Reader) (Reader, error) {
	return gzip.NewReader(r)
}
