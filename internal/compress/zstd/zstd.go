package zstd

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

// EOption is type alias of zstd.EOption.
type EOption = zstd.EOption

// Encoder represents an interface for Encoder of zstd.
type Encoder interface {
	io.WriteCloser
	ReadFrom(r io.Reader) (n int64, err error)
}

// Decoder represents an interface for Decoder of zstd.
type Decoder interface {
	io.Reader
	Close()
	WriteTo(w io.Writer) (int64, error)
}

// Transporter is an interface to create Writer and Reader implementation.
type Transporter interface {
	NewWriter(w io.Writer, opts ...zstd.EOption) (Encoder, error)
	NewReader(r io.Reader, opts ...zstd.DOption) (Decoder, error)
}

type transporter struct{}

// NewTransporter returns Transporter implementation.
func NewTransporter() Transporter {
	return new(transporter)
}

// NewWriter returns Encoder implementation.
func (*transporter) NewWriter(w io.Writer, opts ...zstd.EOption) (Encoder, error) {
	return zstd.NewWriter(w, opts...)
}

// NewReader returns Decoder implementation.
func (*transporter) NewReader(r io.Reader, opts ...zstd.DOption) (Decoder, error) {
	return zstd.NewReader(r, opts...)
}
