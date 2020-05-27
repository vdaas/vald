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

	"github.com/klauspost/compress/gzip"
	"github.com/vdaas/vald/internal/errors"
)

type gzipCompressor struct {
	gobc             Compressor
	compressionLevel int
}

func NewGzip(opts ...GzipOption) (Compressor, error) {
	c := new(gzipCompressor)
	for _, opt := range append(defaultGzipOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

func (g *gzipCompressor) CompressVector(vector []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	gw, err := gzip.NewWriterLevel(buf, g.compressionLevel)
	if err != nil {
		return nil, err
	}

	gob, err := g.gobc.CompressVector(vector)
	if err != nil {
		return nil, err
	}

	_, err = gw.Write(gob)
	if err != nil {
		return nil, err
	}

	err = gw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *gzipCompressor) DecompressVector(bs []byte) ([]float32, error) {
	buf := new(bytes.Buffer)
	gr, err := gzip.NewReader(bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(buf, gr)
	if err != nil {
		return nil, err
	}

	vec, err := g.gobc.DecompressVector(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return vec, nil
}

func (g *gzipCompressor) Reader(src io.Reader) (io.Reader, error) {
	return gzip.NewReader(src)
}

func (g *gzipCompressor) Writer(dst io.Writer) (io.WriteCloser, error) {
	return gzip.NewWriterLevel(dst, g.compressionLevel)
}
