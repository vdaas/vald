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
	"bytes"
	"io"
	"reflect"

	"github.com/klauspost/compress/zstd"
	"github.com/vdaas/vald/internal/errors"
)

type zstdCompressor struct {
	gobc     Compressor
	eoptions []zstd.EOption
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

func (z *zstdCompressor) CompressVector(vector []float32) ([]byte, error) {
	gob, err := z.gobc.CompressVector(vector)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zw, err := zstd.NewWriter(buf, z.eoptions...)
	if err != nil {
		return nil, err
	}

	_, err = zw.ReadFrom(bytes.NewReader(gob))
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (z *zstdCompressor) DecompressVector(bs []byte) ([]float32, error) {
	buf := new(bytes.Buffer)
	zr, err := zstd.NewReader(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	_, err = zr.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	vec, err := z.gobc.DecompressVector(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return vec, nil
}

func (z *zstdCompressor) Reader(src io.Reader) (io.Reader, error) {
	return zstd.NewReader(src)
}

func (z *zstdCompressor) Writer(dst io.Writer) (io.WriteCloser, error) {
	return zstd.NewWriter(dst, z.eoptions...)
}
