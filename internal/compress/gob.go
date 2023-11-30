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

	"github.com/vdaas/vald/internal/compress/gob"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
)

type gobCompressor struct {
	transcoder gob.Transcoder
}

// NewGob returns a Compressor implemented using gob.
func NewGob(opts ...GobOption) (Compressor, error) {
	c := &gobCompressor{
		transcoder: gob.New(),
	}
	for _, opt := range append(defaultGobOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

// CompressVector compresses the data given and returns the compressed data.
// If CompressVector fails, it will return an error.
func (g *gobCompressor) CompressVector(vector []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.transcoder.NewEncoder(buf).Encode(vector)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecompressVector decompresses the compressed data and returns the data.
// If decompress fails, it will return an error.
func (g *gobCompressor) DecompressVector(bs []byte) ([]float32, error) {
	var vector []float32
	err := g.transcoder.NewDecoder(bytes.NewBuffer(bs)).Decode(&vector)
	if err != nil {
		return nil, err
	}

	return vector, nil
}

// Reader returns io.ReadCloser implementation.
func (g *gobCompressor) Reader(src io.ReadCloser) (io.ReadCloser, error) {
	return &gobReader{
		src:     src,
		decoder: g.transcoder.NewDecoder(src),
	}, nil
}

// Writer returns io.WriteCloser implementation.
func (g *gobCompressor) Writer(dst io.WriteCloser) (io.WriteCloser, error) {
	return &gobWriter{
		dst:     dst,
		encoder: g.transcoder.NewEncoder(dst),
	}, nil
}

type gobReader struct {
	src     io.ReadCloser
	decoder gob.Decoder
}

// Read returns the number of bytes for read p (0 <= n <= len(p)).
// If any errors occurs, it will return an error.
func (gr *gobReader) Read(p []byte) (n int, err error) {
	if err = gr.decoder.Decode(&p); err != nil {
		return 0, err
	}

	return len(p), nil
}

// Close closes the reader.
func (gr *gobReader) Close() error {
	return gr.src.Close()
}

type gobWriter struct {
	dst     io.WriteCloser
	encoder gob.Encoder
}

// write returns the number of bytes written from p (0 <= n <= len(p)).
// if any errors occurs, it will return an error.
func (gw *gobWriter) Write(p []byte) (n int, err error) {
	if err = gw.encoder.Encode(&p); err != nil {
		return 0, err
	}

	return len(p), nil
}

// Close closes the writer.
func (gw *gobWriter) Close() error {
	return gw.dst.Close()
}
