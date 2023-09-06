//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"reflect"

	"github.com/vdaas/vald/internal/compress/gzip"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
)

type gzipCompressor struct {
	gobc             Compressor
	compressionLevel int
	gzip             gzip.Gzip
}

// NewGzip returns Compressor implementation.
func NewGzip(opts ...GzipOption) (Compressor, error) {
	c := &gzipCompressor{
		gzip: gzip.New(),
	}
	for _, opt := range append(defaultGzipOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

// CompressVector Compress the data and returns an error if compression fails.
func (g *gzipCompressor) CompressVector(vector []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	gw, err := g.gzip.NewWriterLevel(buf, g.compressionLevel)
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

// CompressVector Decompress the compressed data and returns an error if decompression fails.
func (g *gzipCompressor) DecompressVector(bs []byte) ([]float32, error) {
	buf := new(bytes.Buffer)
	gr, err := g.gzip.NewReader(bytes.NewBuffer(bs))
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

// Reader returns io.ReadCloser implementation.
func (g *gzipCompressor) Reader(src io.ReadCloser) (io.ReadCloser, error) {
	r, err := g.gzip.NewReader(src)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		src: src,
		r:   r,
	}, nil
}

// Writer returns io.WriteCloser implementation.
func (g *gzipCompressor) Writer(dst io.WriteCloser) (io.WriteCloser, error) {
	w, err := g.gzip.NewWriterLevel(dst, g.compressionLevel)
	if err != nil {
		return nil, err
	}

	return &gzipWriter{
		dst: dst,
		w:   w,
	}, nil
}

type gzipReader struct {
	src io.ReadCloser
	r   io.ReadCloser
}

// Read reads up to len(p) bytes into p.
func (g *gzipReader) Read(p []byte) (n int, err error) {
	return g.r.Read(p)
}

// Close closes src and r.
func (g *gzipReader) Close() (err error) {
	err = g.r.Close()
	if err != nil {
		return errors.Join(g.src.Close(), err)
	}

	return g.src.Close()
}

type gzipWriter struct {
	dst io.WriteCloser
	w   io.WriteCloser
}

// Write writes len(p) bytes from p.
func (g *gzipWriter) Write(p []byte) (n int, err error) {
	return g.w.Write(p)
}

// Close closes dst and w.
func (g *gzipWriter) Close() (err error) {
	err = g.w.Close()
	if err != nil {
		return errors.Join(g.dst.Close(), err)
	}

	return g.dst.Close()
}
