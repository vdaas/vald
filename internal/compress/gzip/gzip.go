package gzip

import (
	"io"

	"github.com/klauspost/compress/gzip"
)

// These constants are copied from the gzip package
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

// Builder an interface to create Writer and Reader implementation.
type Builder interface {
	NewWriterLevel(w io.Writer, level int) (Writer, error)
	NewReader(r io.Reader) (Reader, error)
}

type builder struct{}

// NewBuilder returns ReaderWriter implementation.
func NewBuilder() Builder {
	return new(builder)
}

// NewWriterLevel returns Writer implementation.
func (*builder) NewWriterLevel(w io.Writer, level int) (Writer, error) {
	gw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return nil, err
	}
	return gw, nil
}

// NewReader returns Reader implementation.
func (*builder) NewReader(r io.Reader) (Reader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return gr, nil
}
