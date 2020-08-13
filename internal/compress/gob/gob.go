package gob

import (
	"encoding/gob"
	"io"
)

// Encoder represents an interface for Encoder of gob.
type Encoder interface {
	Encode(e interface{}) error
}

// Decoder represents an interface for Decoder of gob.
type Decoder interface {
	Decode(e interface{}) error
}

// Transporter an interface to create Encoder and Decoder implementation.
type Transporter interface {
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

type transporter struct{}

// NewTransporter returns Transporter implementation.
func NewTransporter() Transporter {
	return new(transporter)
}

// NewEncoder returns Encoder implementation.
func (*transporter) NewEncoder(w io.Writer) Encoder {
	return gob.NewEncoder(w)
}

// NewDecoder returns Decoder implementation.
func (*transporter) NewDecoder(r io.Reader) Decoder {
	return gob.NewDecoder(r)
}
