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

	"github.com/vdaas/vald/internal/compress/lz4"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
)

type lz4Compressor struct {
	gobc             Compressor
	compressionLevel int
	lz4              lz4.LZ4
}

// NewLZ4 returns Compressor implementation.
func NewLZ4(opts ...LZ4Option) (Compressor, error) {
	c := &lz4Compressor{
		lz4: lz4.New(),
	}
	for _, opt := range append(defaultLZ4Opts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

// CompressVector compresses the data given and returns the compressed data.
// If CompressVector fails, it will return an error.
func (l *lz4Compressor) CompressVector(vector []float32) (b []byte, err error) {
	gob, err := l.gobc.CompressVector(vector)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zw := l.lz4.NewWriterLevel(buf, l.compressionLevel)
	defer func() {
		cerr := zw.Close()
		if cerr != nil {
			b = nil
			err = errors.Join(err, cerr)
		}
	}()

	_, err = zw.Write(gob)
	if err != nil {
		return nil, err
	}

	err = zw.Flush()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecompressVector decompresses the compressed data and returns the data.
// If decompress fails, it will return an error.
func (l *lz4Compressor) DecompressVector(bs []byte) ([]float32, error) {
	buf := new(bytes.Buffer)
	zr := l.lz4.NewReader(bytes.NewReader(bs))
	_, err := io.Copy(buf, zr)
	if err != nil {
		return nil, err
	}

	vec, err := l.gobc.DecompressVector(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return vec, nil
}

// Reader returns io.ReadCloser implementation.
func (l *lz4Compressor) Reader(src io.ReadCloser) (io.ReadCloser, error) {
	return &lz4Reader{
		src: src,
		r:   l.lz4.NewReader(src),
	}, nil
}

// Writer returns io.WriteCloser implementation.
func (l *lz4Compressor) Writer(dst io.WriteCloser) (io.WriteCloser, error) {
	return &lz4Writer{
		dst: dst,
		w:   l.lz4.NewWriter(dst),
	}, nil
}

type lz4Reader struct {
	src io.ReadCloser
	r   io.Reader
}

// Read returns the number of bytes for read p (0 <= n <= len(p)).
// If any errors occurs, it will return an error.
func (l *lz4Reader) Read(p []byte) (n int, err error) {
	return l.r.Read(p)
}

// Close closes the reader.
func (l *lz4Reader) Close() (err error) {
	return l.src.Close()
}

type lz4Writer struct {
	dst io.WriteCloser
	w   io.WriteCloser
}

// Write returns the number of bytes written from p (0 <= n <= len(p)).
// If any errors occurs, it will return an error.
func (l *lz4Writer) Write(p []byte) (n int, err error) {
	return l.w.Write(p)
}

// Close closes the writer.
func (l *lz4Writer) Close() (err error) {
	err = l.w.Close()
	if err != nil {
		return errors.Join(l.dst.Close(), err)
	}

	return l.dst.Close()
}
