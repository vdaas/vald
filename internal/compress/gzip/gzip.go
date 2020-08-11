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

// ReaderWriter an interface to create Writer and Reader implementation.
type ReaderWriter interface {
	NewWriterLevel(w io.Writer, level int) (Writer, error)
	NewReader(r io.Reader) (Reader, error)
}

type readerWriter struct{}

// NewReaderWriter returns ReaderWriter implementation.
func NewReaderWriter() ReaderWriter {
	return new(readerWriter)
}

// NewWriterLevel returns Writer implementation.
func (rw *readerWriter) NewWriterLevel(w io.Writer, level int) (Writer, error) {
	return NewWriterLevel(w, level)
}

// NewReader returns Reader implementation.
func (rw *readerWriter) NewReader(r io.Reader) (Reader, error) {
	return NewReader(r)
}

// Reader represents an interface for Reader of gzip.
type Reader interface {
	io.ReadCloser
	Reset(r io.Reader) error
	Multistream(ok bool)
}

// NewReader returns Reader implementation.
func NewReader(r io.Reader) (Reader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return gr, nil
}

// Writer represents an interface for Writer of gzip.
type Writer interface {
	io.WriteCloser
	Reset(w io.Writer)
	Flush() error
}

// NewWriterLevel returns Writer implementation.
func NewWriterLevel(w io.Writer, level int) (Writer, error) {
	gw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return nil, err
	}
	return gw, nil
}
