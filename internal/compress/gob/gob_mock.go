//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
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
	return m.NewDecoderFunc(r)
}
