//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"

	"github.com/vdaas/vald/internal/compress/zstd"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
)

type zstdCompressor struct {
	gobc     Compressor
	eoptions []zstd.EOption
	zstd     zstd.Zstd
}

// NewZstd returns the zstd compressor object or any initialization error.
func NewZstd(opts ...ZstdOption) (Compressor, error) {
	c := &zstdCompressor{
		zstd: zstd.New(),
	}
	for _, opt := range append(defaultZstdOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

// CompressVector compresses the data given and returns the compressed data.
// If CompressVector fails, it will return an error.
func (z *zstdCompressor) CompressVector(vector []float32) ([]byte, error) {
	gob, err := z.gobc.CompressVector(vector)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zw, err := z.zstd.NewWriter(buf, z.eoptions...)
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

// DecompressVector decompresses the compressed data and returns the data.
// If decompress fails, it will return an error.
func (z *zstdCompressor) DecompressVector(bs []byte) ([]float32, error) {
	buf := new(bytes.Buffer)
	zr, err := z.zstd.NewReader(bytes.NewReader(bs))
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

// Reader returns io.ReadCloser implementation.
func (z *zstdCompressor) Reader(src io.ReadCloser) (io.ReadCloser, error) {
	r, err := z.zstd.NewReader(src)
	if err != nil {
		return nil, err
	}

	return &zstdReader{
		src: src,
		r:   r,
	}, nil
}

// Writer returns io.WriteCloser implementation.
func (z *zstdCompressor) Writer(dst io.WriteCloser) (io.WriteCloser, error) {
	w, err := z.zstd.NewWriter(dst, z.eoptions...)
	if err != nil {
		return nil, err
	}

	return &zstdWriter{
		dst: dst,
		w:   w,
	}, nil
}

type zstdReader struct {
	src io.ReadCloser
	r   io.Reader
}

// Read returns the number of bytes for read p (0 <= n <= len(p)).
// If any errors occurs, it will return an error.
func (z *zstdReader) Read(p []byte) (n int, err error) {
	return z.r.Read(p)
}

// Close closes the reader.
func (z *zstdReader) Close() error {
	return z.src.Close()
}

type zstdWriter struct {
	dst io.WriteCloser
	w   io.WriteCloser
}

// Write returns the number of bytes written from p (0 <= n <= len(p)).
// If any errors occurs, it will return an error.
func (z *zstdWriter) Write(p []byte) (n int, err error) {
	return z.w.Write(p)
}

// Close closes the writer.
func (z *zstdWriter) Close() (err error) {
	err = z.w.Close()
	if err != nil {
		return errors.Join(z.dst.Close(), err)
	}

	return z.dst.Close()
}
