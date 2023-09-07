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

// Package io provides io functions
package io

import (
	"io"
	"testing"
)

const (
	readerLength = 1024 * 1024
	bufferLength = 32 * 1024
)

type writer struct {
	io.Writer
}

func (*writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type reader struct {
	pos int
	len int
	io.Reader
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.pos == r.len {
		return 0, io.EOF
	}

	read := r.len - r.pos
	if read > len(p) {
		read = len(p)
	}

	r.pos += read
	return read, nil
}

func BenchmarkStandardIOCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := &writer{}
		r := &reader{len: readerLength}

		io.Copy(w, r)
	}
}

func BenchmarkStandardIOCopyBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := make([]byte, bufferLength)
		w := &writer{}
		r := &reader{len: readerLength}
		io.CopyBuffer(w, r, buf)
	}
}

func BenchmarkValdIOCopy(b *testing.B) {
	c := NewCopier(bufferLength)
	for i := 0; i < b.N; i++ {
		w := &writer{}
		r := &reader{len: readerLength}
		c.Copy(w, r)
	}
}

func BenchmarkValdIOCopyBuffer(b *testing.B) {
	c := NewCopier(bufferLength)
	for i := 0; i < b.N; i++ {
		w := &writer{}
		r := &reader{len: readerLength}
		c.CopyBuffer(w, r, nil)
	}
}

func BenchmarkStandardIOCopyParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := &writer{}
			r := &reader{len: readerLength}
			io.Copy(w, r)
		}
	})
}

func BenchmarkStandardIOCopyBufferParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := make([]byte, bufferLength)
			w := &writer{}
			r := &reader{len: readerLength}
			io.CopyBuffer(w, r, buf)
		}
	})
}

func BenchmarkValdIOCopyParallel(b *testing.B) {
	c := NewCopier(bufferLength)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := &writer{}
			r := &reader{len: readerLength}
			c.Copy(w, r)
		}
	})
}

func BenchmarkValdIOCopyBufferParallel(b *testing.B) {
	c := NewCopier(bufferLength)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := &writer{}
			r := &reader{len: readerLength}
			c.CopyBuffer(w, r, nil)
		}
	})
}
