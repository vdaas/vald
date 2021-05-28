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
package io

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// A version of bytes.Buffer without ReadFrom and WriteTo
type Buffer struct {
	bytes.Buffer
	io.ReaderFrom // conflicts with and hides bytes.Buffer's ReaderFrom.
	io.WriterTo   // conflicts with and hides bytes.Buffer's WriterTo.
}

func TestCopy(t *testing.T) {
	rb := new(Buffer)
	wb := new(Buffer)
	txt := "hello, world."
	rb.WriteString(txt)
	Copy(wb, rb)
	if wb.String() != txt {
		t.Errorf("Copy did not work properly")
	}
}

func TestCopyNegative(t *testing.T) {
	rb := new(Buffer)
	wb := new(Buffer)
	rb.WriteString("hello")
	Copy(wb, &io.LimitedReader{R: rb, N: -1})
	if wb.String() != "" {
		t.Errorf("Copy on LimitedReader with N<0 copied data")
	}
}

func TestCopyBuffer(t *testing.T) {
	rb := new(Buffer)
	wb := new(Buffer)
	rb.WriteString("hello, world.")
	Copy(wb, rb) // Tiny buffer to keep it honest.
	if wb.String() != "hello, world." {
		t.Errorf("CopyBuffer did not work properly")
	}
}

func TestCopyBufferNil(t *testing.T) {
	rb := new(Buffer)
	wb := new(Buffer)
	rb.WriteString("hello, world.")
	Copy(wb, rb) // Should allocate a buffer.
	if wb.String() != "hello, world." {
		t.Errorf("CopyBuffer did not work properly")
	}
}

func TestCopyReadFrom(t *testing.T) {
	rb := new(Buffer)
	wb := new(bytes.Buffer) // implements ReadFrom.
	rb.WriteString("hello, world.")
	Copy(wb, rb)
	if wb.String() != "hello, world." {
		t.Errorf("Copy did not work properly")
	}
}

func TestCopyWriteTo(t *testing.T) {
	rb := new(bytes.Buffer) // implements WriteTo.
	wb := new(Buffer)
	rb.WriteString("hello, world.")
	Copy(wb, rb)
	if wb.String() != "hello, world." {
		t.Errorf("Copy did not work properly")
	}
}

type writeToChecker struct {
	bytes.Buffer
	writeToCalled bool
}

func (wt *writeToChecker) WriteTo(w io.Writer) (int64, error) {
	wt.writeToCalled = true
	return wt.Buffer.WriteTo(w)
}

func TestCopyPriority(t *testing.T) {
	rb := new(writeToChecker)
	wb := new(bytes.Buffer)
	rb.WriteString("hello, world.")
	Copy(wb, rb)
	if wb.String() != "hello, world." {
		t.Errorf("Copy did not work properly")
	} else if !rb.writeToCalled {
		t.Errorf("WriteTo was not prioritized over ReadFrom")
	}
}

type zeroErrReader struct {
	err error
}

func (r zeroErrReader) Read(p []byte) (int, error) {
	return copy(p, []byte{0}), r.err
}

type errWriter struct {
	err error
}

func (w errWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestCopyReadErrWriteErr(t *testing.T) {
	er, ew := errors.New("readError"), errors.New("writeError")
	r, w := zeroErrReader{err: er}, errWriter{err: ew}
	n, err := Copy(w, r)
	if n != 0 || err != ew {
		t.Errorf("Copy(zeroErrReader, errWriter) = %d, %v; want 0, writeError", n, err)
	}
}

type largeWriter struct {
	err error
}

func (w largeWriter) Write(p []byte) (int, error) {
	return len(p) + 1, w.err
}

func TestCopyLargeWriter(t *testing.T) {
	want := errors.New("invalid write result")
	rb := new(Buffer)
	wb := largeWriter{}
	rb.WriteString("hello, world.")
	_, err := Copy(wb, rb)
	if err.Error() != want.Error() {
		t.Errorf("Copy error: got %v, want %v", err, want)
	}

	want = errors.New("largeWriterError")
	rb = new(Buffer)
	wb = largeWriter{err: want}
	rb.WriteString("hello, world.")
	if _, err := Copy(wb, rb); err != want {
		t.Errorf("Copy error: got %v, want %v", err, want)
	}
}
