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

// Transcoder is an interface to create Encoder and Decoder implementation.
type Transcoder interface {
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

type transcoder struct{}

// New returns Transcoder implementation.
func New() Transcoder {
	return new(transcoder)
}

// NewEncoder returns Encoder implementation.
func (*transcoder) NewEncoder(w io.Writer) Encoder {
	return gob.NewEncoder(w)
}

// NewDecoder returns Decoder implementation.
func (*transcoder) NewDecoder(r io.Reader) Decoder {
	return gob.NewDecoder(r)
}
