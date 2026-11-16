// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package io

import (
	"bytes"
	"io"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/goleak"
)

const (
	testString = "hello, world."
)

func TestCopy(t *testing.T) {
	// A version of bytes.buffer without ReadFrom and WriteTo
	type buffer struct {
		io.ReaderFrom
		io.WriterTo
		bytes.Buffer
	}
	type args struct {
		dst io.Writer
		src io.Reader
	}
	type want struct {
		err         error
		wantDst     string
		wantWritten int64
	}
	type test struct {
		args       args
		checkFunc  func(want, int64, io.Writer, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		want       want
	}
	checkFunc := func(w want, gotWritten int64, got string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotWritten, w.wantWritten) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotWritten, w.wantWritten)
		}
		if !reflect.DeepEqual(got, w.wantDst) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.wantDst)
		}
		return nil
	}
	defaultCheckFunc := func(w want, gotWritten int64, dst io.Writer, err error) error {
		v := dst.(*buffer)
		return checkFunc(w, gotWritten, v.String(), err)
	}
	bufferWithStringFunc := func() *buffer {
		buf := new(buffer)
		buf.WriteString(testString)
		return buf
	}
	tests := []test{
		{
			name: "copy success with not nil string",
			args: args{
				dst: new(buffer),
				src: bufferWithStringFunc(),
			},
			want: want{
				wantWritten: int64(len(testString)),
				wantDst:     testString,
				err:         nil,
			},
		},
		{
			name: "copy success with LimitedReader",
			args: args{
				dst: new(buffer),
				src: &io.LimitedReader{R: bufferWithStringFunc(), N: -1},
			},
			want: want{
				wantWritten: 0,
				wantDst:     "",
				err:         nil,
			},
		},
		{
			name: "copy success with LimitedReader smaller buffer size than defaultBufferSize",
			args: args{
				dst: new(buffer),
				src: &io.LimitedReader{R: bufferWithStringFunc(), N: 32 * 1024},
			},
			want: want{
				wantWritten: int64(len(testString)),
				wantDst:     testString,
				err:         nil,
			},
		},
		{
			name: "copy success with using ReadFrom function",
			args: func() args {
				src := new(bytes.Buffer)
				src.WriteString(testString)
				return args{
					dst: new(buffer),
					src: src,
				}
			}(),
			want: want{
				wantWritten: int64(len(testString)),
				wantDst:     testString,
				err:         nil,
			},
		},
		{
			name: "copy success with using WriteTo function",
			args: args{
				dst: new(bytes.Buffer),
				src: bufferWithStringFunc(),
			},
			want: want{
				wantWritten: int64(len(testString)),
				wantDst:     testString,
				err:         nil,
			},
			checkFunc: func(w want, gotWritten int64, dst io.Writer, err error) error {
				v := dst.(*bytes.Buffer)
				return checkFunc(w, gotWritten, v.String(), err)
			},
		},
		{
			name: "copy fail when dst is nil",
			args: args{
				dst: nil,
				src: new(buffer),
			},
			want: want{
				wantWritten: 0,
				wantDst:     "",
				err:         errors.New("empty source or destination"),
			},
			checkFunc: func(w want, gotWritten int64, dst io.Writer, err error) error {
				return checkFunc(w, gotWritten, "", err)
			},
		},
		{
			name: "copy fail when src is nil",
			args: args{
				dst: new(buffer),
				src: nil,
			},
			want: want{
				wantWritten: 0,
				wantDst:     "",
				err:         errors.New("empty source or destination"),
			},
			checkFunc: func(w want, gotWritten int64, dst io.Writer, err error) error {
				return checkFunc(w, gotWritten, "", err)
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotWritten, err := Copy(test.args.dst, test.args.src)
			if err := test.checkFunc(test.want, gotWritten, test.args.dst, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewCopier(t *testing.T) {
	type args struct {
		size int
	}
	type want struct {
		want Copier
	}
	type test struct {
		want       want
		checkFunc  func(want, Copier) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, got Copier) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return default buffer size Copier",
			args: args{
				size: 0,
			},
			checkFunc: func(w want, got Copier) error {
				if got == nil {
					return errors.New("got is nil")
				}
				return nil
			},
		},
		{
			name: "return user set buffer size Copier",
			args: args{
				size: 128 * 1024,
			},
			checkFunc: func(w want, got Copier) error {
				if got == nil {
					return errors.New("got is nil")
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := NewCopier(test.args.size)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_copier_Copy(t *testing.T) {
	type args struct {
		src io.Reader
	}
	type fields struct {
		bufSize int64
	}
	type want struct {
		err         error
		wantDst     string
		wantWritten int64
	}
	type test struct {
		args       args
		checkFunc  func(want, int64, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		want       want
		fields     fields
	}
	defaultCheckFunc := func(w want, gotWritten int64, gotDst string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotWritten, w.wantWritten) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotWritten, w.wantWritten)
		}
		if !reflect.DeepEqual(gotDst, w.wantDst) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDst, w.wantDst)
		}
		return nil
	}
	tests := []test{}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &copier{
				bufSize: test.fields.bufSize,
			}
			c.pool = sync.Pool{
				New: func() any {
					return bytes.NewBuffer(make([]byte, int(atomic.LoadInt64(&c.bufSize))))
				},
			}
			dst := &bytes.Buffer{}

			gotWritten, err := c.Copy(dst, test.args.src)
			if err := checkFunc(test.want, gotWritten, dst.String(), err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestCopyBuffer(t *testing.T) {
// 	type args struct {
// 		src io.Reader
// 		buf []byte
// 	}
// 	type want struct {
// 		wantWritten int64
// 		wantDst     string
// 		err         error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, int64, string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotWritten int64, gotDst string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotWritten, w.wantWritten) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotWritten, w.wantWritten)
// 		}
// 		if !reflect.DeepEqual(gotDst, w.wantDst) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDst, w.wantDst)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           src:nil,
// 		           buf:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           src:nil,
// 		           buf:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			dst := &bytes.Buffer{}
//
// 			gotWritten, err := CopyBuffer(dst, test.args.src, test.args.buf)
// 			if err := checkFunc(test.want, gotWritten, err, dst.String()); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_copier_CopyBuffer(t *testing.T) {
// 	type args struct {
// 		src io.Reader
// 		buf []byte
// 	}
// 	type fields struct {
// 		bufSize int64
// 	}
// 	type want struct {
// 		wantWritten int64
// 		wantDst     string
// 		err         error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int64, string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotWritten int64, gotDst string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotWritten, w.wantWritten) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotWritten, w.wantWritten)
// 		}
// 		if !reflect.DeepEqual(gotDst, w.wantDst) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDst, w.wantDst)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           src:nil,
// 		           buf:nil,
// 		       },
// 		       fields: fields {
// 		           bufSize:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           src:nil,
// 		           buf:nil,
// 		           },
// 		           fields: fields {
// 		           bufSize:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			c := &copier{
// 				bufSize: test.fields.bufSize,
// 			}
// 			dst := &bytes.Buffer{}
//
// 			gotWritten, err := c.CopyBuffer(dst, test.args.src, test.args.buf)
// 			if err := checkFunc(test.want, gotWritten, err, dst.String()); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_copier_copyBuffer(t *testing.T) {
// 	type args struct {
// 		src io.Reader
// 		buf *bytes.Buffer
// 	}
// 	type fields struct {
// 		bufSize int64
// 	}
// 	type want struct {
// 		wantWritten int64
// 		wantDst     string
// 		err         error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int64, string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotWritten int64, gotDst string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotWritten, w.wantWritten) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotWritten, w.wantWritten)
// 		}
// 		if !reflect.DeepEqual(gotDst, w.wantDst) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDst, w.wantDst)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           src:nil,
// 		           buf:bytes.Buffer{},
// 		       },
// 		       fields: fields {
// 		           bufSize:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           src:nil,
// 		           buf:bytes.Buffer{},
// 		           },
// 		           fields: fields {
// 		           bufSize:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			c := &copier{
// 				bufSize: test.fields.bufSize,
// 			}
// 			dst := &bytes.Buffer{}
//
// 			gotWritten, err := c.copyBuffer(dst, test.args.src, test.args.buf)
// 			if err := checkFunc(test.want, gotWritten, err, dst.String()); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// scriptedReader replays steps and counts Read calls. It implements neither
// io.WriterTo nor io.ReaderFrom so that the copy loop is exercised.
type scriptedReader struct {
	steps []scriptedStep
	reads int
}

type scriptedStep struct {
	err  error
	data string
}

var errReadAfterEnd = errors.New("read called after the script ended")

func (r *scriptedReader) Read(p []byte) (int, error) {
	r.reads++
	if r.reads > len(r.steps) {
		return 0, errReadAfterEnd
	}
	s := r.steps[r.reads-1]
	return copy(p, s.data), s.err
}

// writerOnly hides bytes.Buffer's io.ReaderFrom.
type writerOnly struct {
	buf bytes.Buffer
}

func (w *writerOnly) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

const (
	chunkHello = "hello"
	chunkABC   = "abc"
)

func Test_copier_copyBuffer_readError(t *testing.T) {
	errRead := errors.New("read failed")
	type want struct {
		err         error
		wantDst     string
		wantReads   int
		wantWritten int64
	}
	tests := []struct {
		name  string
		steps []scriptedStep
		want  want
	}{
		{
			name:  "read error returned with data stops reading and is returned",
			steps: []scriptedStep{{data: chunkHello, err: errRead}},
			want:  want{wantWritten: 5, wantDst: chunkHello, wantReads: 1, err: errRead},
		},
		{
			name:  "read error returned with data after a good chunk",
			steps: []scriptedStep{{data: chunkABC}, {data: "def", err: errRead}},
			want:  want{wantWritten: 6, wantDst: "abcdef", wantReads: 2, err: errRead},
		},
		{
			name:  "read error without data is returned",
			steps: []scriptedStep{{data: chunkABC}, {err: errRead}},
			want:  want{wantWritten: 3, wantDst: chunkABC, wantReads: 2, err: errRead},
		},
		{
			name:  "EOF returned with data completes without reading again",
			steps: []scriptedStep{{data: chunkHello, err: io.EOF}},
			want:  want{wantWritten: 5, wantDst: chunkHello, wantReads: 1},
		},
		{
			name:  "EOF returned with data after a good chunk completes without reading again",
			steps: []scriptedStep{{data: chunkABC}, {data: "def", err: io.EOF}},
			want:  want{wantWritten: 6, wantDst: "abcdef", wantReads: 2},
		},
		{
			name:  "EOF without data completes",
			steps: []scriptedStep{{data: chunkABC}, {err: io.EOF}},
			want:  want{wantWritten: 3, wantDst: chunkABC, wantReads: 2},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			run := map[string]func(Copier, *writerOnly, *scriptedReader) (int64, error){
				"Copy": func(c Copier, dst *writerOnly, src *scriptedReader) (int64, error) {
					return c.Copy(dst, src)
				},
				"CopyBuffer": func(c Copier, dst *writerOnly, src *scriptedReader) (int64, error) {
					return c.CopyBuffer(dst, src, make([]byte, 16))
				},
			}
			for name, fn := range run {
				dst := new(writerOnly)
				src := &scriptedReader{steps: test.steps}
				gotWritten, err := fn(NewCopier(0), dst, src)
				if !errors.Is(err, test.want.err) {
					tt.Errorf("%s: got error: %v, want: %v", name, err, test.want.err)
				}
				if gotWritten != test.want.wantWritten {
					tt.Errorf("%s: got written: %d, want: %d", name, gotWritten, test.want.wantWritten)
				}
				if got := dst.buf.String(); got != test.want.wantDst {
					tt.Errorf("%s: got dst: %q, want: %q", name, got, test.want.wantDst)
				}
				if src.reads != test.want.wantReads {
					tt.Errorf("%s: got Read calls: %d, want: %d", name, src.reads, test.want.wantReads)
				}
			}
		})
	}
}

func Test_copier_CopyBuffer_emptyBuffer(t *testing.T) {
	tests := []struct {
		name string
		buf  []byte
	}{
		{name: "non-nil empty buffer", buf: []byte{}},
		{name: "non-nil zero length buffer with capacity", buf: make([]byte, 0, 8)},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			dst := new(writerOnly)
			src := &scriptedReader{steps: []scriptedStep{{data: chunkABC, err: io.EOF}}}
			gotWritten, err := NewCopier(0).CopyBuffer(dst, src, test.buf)
			if !errors.Is(err, errors.New("empty buffer in CopyBuffer")) {
				tt.Errorf("got error: %v, want: empty buffer in CopyBuffer", err)
			}
			if gotWritten != 0 || src.reads != 0 || dst.buf.Len() != 0 {
				tt.Errorf("nothing should be copied: written=%d reads=%d dst=%q", gotWritten, src.reads, dst.buf.String())
			}
		})
	}
}
