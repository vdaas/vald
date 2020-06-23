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

	lz4 "github.com/pierrec/lz4/v3"
	"github.com/vdaas/vald/internal/errors"
)

type lz4Compressor struct {
	gobc             Compressor
	compressionLevel int
}

func NewLZ4(opts ...LZ4Option) (Compressor, error) {
	c := new(lz4Compressor)
	for _, opt := range append(defaultLZ4Opts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

func (l *lz4Compressor) CompressVector(vector []float32) ([]byte, error) {
	gob, err := l.gobc.CompressVector(vector)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zw := lz4.NewWriter(buf)
	zw.Header.CompressionLevel = l.compressionLevel
	defer func() {
		cerr := zw.Close()
		if cerr != nil {
			err = errors.Wrap(err, cerr.Error())
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

func (l *lz4Compressor) DecompressVector(bs []byte) ([]float32, error) {
	buf := new(bytes.Buffer)
	zr := lz4.NewReader(bytes.NewReader(bs))
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

func (l *lz4Compressor) Reader(src io.Reader) (io.Reader, error) {
	return lz4.NewReader(src), nil
}

func (l *lz4Compressor) Writer(dst io.WriteCloser) (io.WriteCloser, error) {
	return &lz4Writer{
		dst: dst,
		w:   lz4.NewWriter(dst),
	}, nil
}

type lz4Writer struct {
	dst io.WriteCloser
	w   io.WriteCloser
}

func (l *lz4Writer) Write(p []byte) (n int, err error) {
	return l.w.Write(p)
}

func (l *lz4Writer) Close() (err error) {
	err = l.w.Close()
	if err != nil {
		return errors.Wrap(l.dst.Close(), err.Error())
	}

	return l.dst.Close()
}
