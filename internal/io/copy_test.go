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
	"bytes"
	"io"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
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
	tests := []test{
		func() test {
			dst := new(buffer)
			src := new(buffer)
			txt := "hello, world."
			src.WriteString(txt)
			return test{
				name: "copy string",
				args: args{
					dst: dst,
					src: src,
				},
				want: want{
					wantWritten: int64(len(txt)),
					wantDst:     txt,
					err:         nil,
				},
			}
		}(),
		func() test {
			dst := new(buffer)
			src := new(buffer)
			src.WriteString("hello")
			return test{
				name: "copy with LimitedReader",
				args: args{
					dst: dst,
					src: &io.LimitedReader{R: src, N: -1},
				},
				want: want{
					wantWritten: 0,
					wantDst:     "",
					err:         nil,
				},
			}
		}(),
		func() test {
			dst := new(buffer)
			src := new(buffer)
			txt := "hello"
			src.WriteString(txt)
			bufferSize := 32 * 1024
			return test{
				name: "copy with LimitedReader smaller buffer than defaultBufferSize",
				args: args{
					dst: dst,
					src: &io.LimitedReader{R: src, N: int64(bufferSize)},
				},
				want: want{
					wantWritten: int64(len(txt)),
					wantDst:     txt,
					err:         nil,
				},
			}
		}(),
		func() test {
			dst := new(buffer)
			src := new(bytes.Buffer)
			txt := "hello, world."
			src.WriteString(txt)
			return test{
				name: "copy with ReadFrom",
				args: args{
					dst: dst,
					src: src,
				},
				want: want{
					wantWritten: int64(len(txt)),
					wantDst:     txt,
					err:         nil,
				},
			}
		}(),
		func() test {
			dst := new(bytes.Buffer)
			src := new(buffer)
			txt := "hello, world."
			src.WriteString(txt)
			return test{
				name: "copy with WriteTo",
				args: args{
					dst: dst,
					src: src,
				},
				want: want{
					wantWritten: int64(len(txt)),
					wantDst:     txt,
					err:         nil,
				},
				checkFunc: func(w want, gotWritten int64, dst io.Writer, err error) error {
					v := dst.(*bytes.Buffer)
					return checkFunc(w, gotWritten, v.String(), err)
				},
			}
		}(),
		{
			name: "dst is nil",
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
			name: "src is nil",
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := NewCopier(test.args.size)
			if err := test.checkFunc(test.want, got); err != nil {
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
		pool    sync.Pool
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &copier{
				bufSize: test.fields.bufSize,
				pool:    test.fields.pool,
			}
			dst := &bytes.Buffer{}

			gotWritten, err := c.Copy(dst, test.args.src)
			if err := test.checkFunc(test.want, gotWritten, dst.String(), err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
