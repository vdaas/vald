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

// Transporter is an interface to create Writer and Reader implementation.
type Transporter interface {
	NewReader(r io.Reader) (Reader, error)
	NewWriterLevel(w io.Writer, level int) (Writer, error)
}

type transporter struct{}

// NewTransporter returns Transporter implementation.
func NewTransporter() Transporter {
	return new(transporter)
}

// NewWriterLevel returns Writer implementation.
func (*transporter) NewWriterLevel(w io.Writer, level int) (Writer, error) {
	return gzip.NewWriterLevel(w, level)
}

// NewReader returns Reader implementation.
func (*transporter) NewReader(r io.Reader) (Reader, error) {
	return gzip.NewReader(r)
}
