package gob

import "io"

// MockEncoder represents mock struct of Encoder.
type MockEncoder struct {
	EncodeFunc func(e interface{}) error
}

// Encode calls EncodeFunc.
func (m *MockEncoder) Encode(e interface{}) error {
	return m.EncodeFunc(e)
}

// MockDecoder represents mock struct of Decoder.
type MockDecoder struct {
	DecodeFunc func(e interface{}) error
}

// Decode calls DecodeFunc.
func (m *MockDecoder) Decode(e interface{}) error {
	return m.DecodeFunc(e)
}

// MockTranscoder represents mock struct of Transcoder.
type MockTranscoder struct {
	NewEncoderFunc func(w io.Writer) Encoder
	NewDecoderFunc func(r io.Reader) Decoder
}

// NewEncoder calls NewEncoderFunc.
func (m *MockTranscoder) NewEncoder(w io.Writer) Encoder {
	return m.NewEncoderFunc(w)
}

// NewDecoder calls NewEncoderFunc.
func (m *MockTranscoder) NewDecoder(r io.Reader) Decoder {
	return m.NewDecoder(r)
}
