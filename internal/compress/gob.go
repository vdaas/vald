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

	"github.com/vdaas/vald/internal/compress/gob"
	"github.com/vdaas/vald/internal/errors"
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

func (g *gobCompressor) CompressVector(vector []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.transcoder.NewEncoder(buf).Encode(vector)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *gobCompressor) DecompressVector(bs []byte) ([]float32, error) {
	var vector []float32
	err := g.transcoder.NewDecoder(bytes.NewBuffer(bs)).Decode(&vector)
	if err != nil {
		return nil, err
	}

	return vector, nil
}

func (g *gobCompressor) Reader(src io.ReadCloser) (io.ReadCloser, error) {
	return &gobReader{
		src:     src,
		decoder: g.transcoder.NewDecoder(src),
	}, nil
}

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

func (gr *gobReader) Read(p []byte) (n int, err error) {
	err = gr.decoder.Decode(&p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (gr *gobReader) Close() error {
	return gr.src.Close()
}

type gobWriter struct {
	dst     io.WriteCloser
	encoder gob.Encoder
}

func (gw *gobWriter) Write(p []byte) (n int, err error) {
	err = gw.encoder.Encode(&p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (gw *gobWriter) Close() error {
	return gw.dst.Close()
}
