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

// Package io provides io functions
package io

import (
	"bytes"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

const (
	testString = "hello, world."
)

func TestCopy(t *testing.T) {
	// A version of bytes.buffer without ReadFrom and WriteTo
	type buffer struct {
		bytes.Buffer
		io.ReaderFrom // conflicts with and hides bytes.Buffer's ReaderFrom.
		io.WriterTo   // conflicts with and hides bytes.Buffer's WriterTo.
	}
	type args struct {
		dst io.Writer
		src io.Reader
	}
	type want struct {
		wantWritten int64
		wantDst     string
		err         error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, int64, io.Writer, error) error
		beforeFunc func(args)
		afterFunc  func(args)
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
		name       string
		args       args
		want       want
		checkFunc  func(want, Copier) error
		beforeFunc func(args)
		afterFunc  func(args)
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
		wantWritten int64
		wantDst     string
		err         error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int64, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
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
				New: func() interface{} {
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
