//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

import (
	"encoding/gob"

	"github.com/vdaas/vald/internal/io"
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
