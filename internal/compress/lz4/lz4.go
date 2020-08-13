package lz4

import (
	"io"

	lz4 "github.com/pierrec/lz4/v3"
)

// Header is type alias of lz4.Header.
type Header = lz4.Header

// Reader represents an interface for Reader of lz4.
type Reader interface {
	io.Reader
}

// Writer represents an interface for Writer of lz4.
type Writer interface {
	io.WriteCloser
	Header() *Header
	Flush() error
}

type writer struct {
	*lz4.Writer
}

// Header returns lz4.Writer.Header object.
func (w *writer) Header() *Header {
	return &w.Writer.Header
}

// Transporter an interface to create Writer and Reader implementation.
type Transporter interface {
	NewWriter(w io.Writer) Writer
	NewReader(r io.Reader) Reader
}

type transporter struct{}

// NewTransporter returns Transporter implementation.
func NewTransporter() Transporter {
	return new(transporter)
}

// NewWriter returns Writer implementation.
func (*transporter) NewWriter(w io.Writer) Writer {
	return &writer{
		Writer: lz4.NewWriter(w),
	}
}

// NewReader returns Reader implementation.
func (*transporter) NewReader(r io.Reader) Reader {
	return lz4.NewReader(r)
}
