//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	// TODO
	// which is the better library of zstd algorithm?
	// "github.com/valyala/gozstd"
	"reflect"

	"github.com/DataDog/zstd"
	"github.com/vdaas/vald/internal/errors"
)

type zstdCompressor struct {
	gobc             Compressor
	compressionLevel int
}

func NewZstd(opts ...ZstdOption) (Compressor, error) {
	c := new(zstdCompressor)
	for _, opt := range append(defaultZstdOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

func (z *zstdCompressor) CompressVector(vector []float64) ([]byte, error) {
	gob, err := z.gobc.CompressVector(vector)
	if err != nil {
		return nil, err
	}

	return zstd.CompressLevel(nil, gob, z.compressionLevel)
}

func (z *zstdCompressor) DecompressVector(bs []byte) ([]float64, error) {
	bufbytes, err := zstd.Decompress(nil, bs)
	if err != nil {
		return nil, err
	}

	vec, err := z.gobc.DecompressVector(bufbytes)
	if err != nil {
		return nil, err
	}

	return vec, nil
}
