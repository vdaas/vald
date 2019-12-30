//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

// Package compress provides compress functions
package compress

import (
	"bytes"
	"encoding/gob"
)

type gobCompressor struct {
}

func NewGob() Compressor {
	return &gobCompressor{}
}

func (g *gobCompressor) CompressVector(vector []float64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(vector)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *gobCompressor) DecompressVector(bs []byte) ([]float64, error) {
	var vector []float64
	err := gob.NewDecoder(bytes.NewBuffer(bs)).Decode(&vector)
	if err != nil {
		return nil, err
	}

	return vector, nil
}
