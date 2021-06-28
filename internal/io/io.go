//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package io provides io functions
package io

import (
	"context"
	"io"

	"github.com/vdaas/vald/internal/errors"
)

type (
	ByteReader    = io.ByteReader
	ByteScanner   = io.ByteScanner
	ByteWriter    = io.ByteWriter
	Closer        = io.Closer
	LimitedReader = io.LimitedReader
	ReadCloser    = io.ReadCloser
	ReadSeeker    = io.ReadSeeker
	Reader        = io.Reader
	ReaderAt      = io.ReaderAt
	ReaderFrom    = io.ReaderFrom
	RuneReader    = io.RuneReader
	RuneScanner   = io.RuneScanner
	Seeker        = io.Seeker
	StringWriter  = io.StringWriter
	WriteCloser   = io.WriteCloser
	Writer        = io.Writer
	WriterAt      = io.WriterAt
	WriterTo      = io.WriterTo
)

const (
	SeekStart   = io.SeekStart
	SeekCurrent = io.SeekCurrent
	SeekEnd     = io.SeekEnd
)

var (
	Pipe             = io.Pipe
	EOF              = io.EOF
	Discard          = io.Discard
	NewSectionReader = io.NewSectionReader
)

type ctxReader struct {
	ctx context.Context
	r   Reader
}

func NewReaderWithContext(ctx context.Context, r Reader) (Reader, error) {
	if ctx == nil {
		return nil, errors.NewErrContextNotProvided()
	}

	if r == nil {
		return nil, errors.NewErrReaderNotProvided()
	}

	return &ctxReader{
		ctx: ctx,
		r:   r,
	}, nil
}

func NewReadCloserWithContext(ctx context.Context, r ReadCloser) (ReadCloser, error) {
	if ctx == nil {
		return nil, errors.NewErrContextNotProvided()
	}

	if r == nil {
		return nil, errors.NewErrReaderNotProvided()
	}

	return &ctxReader{
		ctx: ctx,
		r:   r,
	}, nil
}

func (r *ctxReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
	}
	return r.r.Read(p)
}

func (r *ctxReader) Close() error {
	select {
	case <-r.ctx.Done():
		return r.ctx.Err()
	default:
	}

	if c, ok := r.r.(Closer); ok {
		return c.Close()
	}

	return nil
}

type ctxWriter struct {
	ctx context.Context
	w   Writer
}

func NewWriterWithContext(ctx context.Context, w Writer) (Writer, error) {
	if ctx == nil {
		return nil, errors.NewErrContextNotProvided()
	}

	if w == nil {
		return nil, errors.NewErrWriterNotProvided()
	}

	return &ctxWriter{
		ctx: ctx,
		w:   w,
	}, nil
}

func NewWriteCloserWithContext(ctx context.Context, w WriteCloser) (WriteCloser, error) {
	if ctx == nil {
		return nil, errors.NewErrContextNotProvided()
	}

	if w == nil {
		return nil, errors.NewErrWriterNotProvided()
	}

	return &ctxWriter{
		ctx: ctx,
		w:   w,
	}, nil
}

func (w *ctxWriter) Write(p []byte) (n int, err error) {
	select {
	case <-w.ctx.Done():
		return 0, w.ctx.Err()
	default:
	}
	return w.w.Write(p)
}

func (w *ctxWriter) Close() error {
	select {
	case <-w.ctx.Done():
		return w.ctx.Err()
	default:
	}

	if c, ok := w.w.(Closer); ok {
		return c.Close()
	}

	return nil
}
