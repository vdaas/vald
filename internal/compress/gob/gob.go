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

// Builder an interface to create Encoder and Decoder implementation.
type Builder interface {
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

type builder struct{}

// NewBuilder returns Builder implementation.
func NewBuilder() Builder {
	return new(builder)
}

// NewEncoder returns Encoder implementation.
func (*builder) NewEncoder(w io.Writer) Encoder {
	return gob.NewEncoder(w)
}

// NewDecoder returns Decoder implementation.
func (*builder) NewDecoder(r io.Reader) Decoder {
	return gob.NewDecoder(r)
}
