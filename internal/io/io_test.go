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
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNewReaderWithContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		r   io.Reader
	}
	type want struct {
		want io.Reader
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, io.Reader, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got io.Reader, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "success when context.Context and io.Reader are not nil",
			args: args{
				ctx: context.Background(),
				r:   &bytes.Buffer{},
			},
			want: want{
				want: &ctxReader{
					ctx: context.Background(),
					r:   &bytes.Buffer{},
				},
				err: nil,
			},
		},
		{
			name: "fail when io.Reader is nil",
			args: args{
				ctx: context.Background(),
				r:   nil,
			},
			want: want{
				want: nil,
				err:  errors.NewErrReaderNotProvided,
			},
		},
		{
			name: "fail when context.Context is nil",
			args: args{
				ctx: nil,
				r:   &bytes.Buffer{},
			},
			want: want{
				want: nil,
				err:  errors.NewErrContextNotProvided,
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

			got, err := NewReaderWithContext(test.args.ctx, test.args.r)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewReadCloserWithContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		r   io.ReadCloser
	}
	type want struct {
		want io.ReadCloser
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, io.ReadCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got io.ReadCloser, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "success when context.Context and io.ReadCloser are not nil",
			args: args{
				ctx: context.Background(),
				r:   NopCloser(&bytes.Buffer{}),
			},
			want: want{
				want: &ctxReader{
					ctx: context.Background(),
					r:   NopCloser(&bytes.Buffer{}),
				},
				err: nil,
			},
		},
		{
			name: "fail when io.ReadCloser is nil",
			args: args{
				ctx: context.Background(),
				r:   nil,
			},
			want: want{
				want: nil,
				err:  errors.NewErrReaderNotProvided,
			},
		},
		{
			name: "fail when context.Context is nil",
			args: args{
				ctx: nil,
				r:   NopCloser(&bytes.Buffer{}),
			},
			want: want{
				want: nil,
				err:  errors.NewErrContextNotProvided,
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

			got, err := NewReadCloserWithContext(test.args.ctx, test.args.r)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ctxReader_Read(t *testing.T) {
	t.Parallel()
	type args struct {
		p []byte
	}
	type fields struct {
		ctx context.Context
		r   io.Reader
	}
	type want struct {
		wantN int
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotN int, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got = %v, want %v", gotN, w.wantN)
		}
		return nil
	}
	tests := []test{
		func() test {
			txt := "hello, world."
			r := &bytes.Buffer{}
			r.WriteString(txt)
			return test{
				name: "success when doing nothing",
				args: args{
					p: make([]byte, 64),
				},
				fields: fields{
					ctx: context.Background(),
					r:   r,
				},
				want: want{
					wantN: len(txt),
					err:   nil,
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "fail when calling cancel function",
				args: args{
					p: make([]byte, 64),
				},
				fields: fields{
					ctx: ctx,
					r:   &bytes.Buffer{},
				},
				want: want{
					wantN: 0,
					err:   context.Canceled,
				},
				beforeFunc: func(args) {
					cancel()
				},
			}
		}(),
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
			r := &ctxReader{
				ctx: test.fields.ctx,
				r:   test.fields.r,
			}

			gotN, err := r.Read(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ctxReader_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		ctx context.Context
		r   io.Reader
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx := context.Background()
			return test{
				name: "success when doing nothing",
				fields: fields{
					ctx: ctx,
					r:   &bytes.Buffer{},
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "fail when calling cancel function",
				fields: fields{
					ctx: ctx,
					r:   &bytes.Buffer{},
				},
				want: want{
					err: context.Canceled,
				},
				beforeFunc: func() {
					cancel()
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			return test{
				name: "success with Closer",
				fields: fields{
					ctx: ctx,
					r:   NopCloser(&bytes.Buffer{}),
				},
				want: want{
					err: nil,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			r := &ctxReader{
				ctx: test.fields.ctx,
				r:   test.fields.r,
			}

			err := r.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewWriterWithContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		w   io.Writer
	}
	type want struct {
		want io.Writer
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, io.Writer, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got io.Writer, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx := context.Background()
			w := &bytes.Buffer{}
			return test{
				name: "success when context.Context and io.Writer is not nil",
				args: args{
					ctx: ctx,
					w:   w,
				},
				want: want{
					want: &ctxWriter{
						ctx: ctx,
						w:   w,
					},
					err: nil,
				},
			}
		}(),
		{
			name: "fail when io.Writer is nil",
			args: args{
				ctx: context.Background(),
				w:   nil,
			},
			want: want{
				want: nil,
				err:  errors.NewErrWriterNotProvided,
			},
		},
		{
			name: "fail when context.Context is nil",
			args: args{
				ctx: nil,
				w:   &bytes.Buffer{},
			},
			want: want{
				want: nil,
				err:  errors.NewErrContextNotProvided,
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

			got, err := NewWriterWithContext(test.args.ctx, test.args.w)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

type nopWriteCloser struct {
	*bytes.Buffer
}

func (*nopWriteCloser) Close() error {
	return nil
}

func TestNewWriteCloserWithContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		w   io.WriteCloser
	}
	type want struct {
		want io.WriteCloser
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, io.WriteCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got io.WriteCloser, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "success when context.Context and io.WriteClose are not nil",
			args: args{
				ctx: context.Background(),
				w:   &nopWriteCloser{},
			},
			want: want{
				want: &ctxWriter{
					ctx: context.Background(),
					w:   &nopWriteCloser{},
				},
				err: nil,
			},
		},
		{
			name: "fail when io.WriteCloser is nil",
			args: args{
				ctx: context.Background(),
				w:   nil,
			},
			want: want{
				want: nil,
				err:  errors.NewErrWriterNotProvided,
			},
		},
		{
			name: "fail when context.Context is nil",
			args: args{
				ctx: nil,
				w:   &nopWriteCloser{},
			},
			want: want{
				want: nil,
				err:  errors.NewErrContextNotProvided,
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

			got, err := NewWriteCloserWithContext(test.args.ctx, test.args.w)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ctxWriter_Write(t *testing.T) {
	t.Parallel()
	type args struct {
		p []byte
	}
	type fields struct {
		ctx context.Context
		w   io.Writer
	}
	type want struct {
		wantN int
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotN int, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got = %v, want %v", gotN, w.wantN)
		}
		return nil
	}
	tests := []test{
		func() test {
			txt := "hello, world."
			ctx := context.Background()
			w := &bytes.Buffer{}
			return test{
				name: "success when doing nothing",
				args: args{
					p: []byte(txt),
				},
				fields: fields{
					ctx: ctx,
					w:   w,
				},
				want: want{
					wantN: len(txt),
					err:   nil,
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "fail when calling cancel function",
				args: args{
					[]byte{},
				},
				fields: fields{
					ctx: ctx,
					w:   &bytes.Buffer{},
				},
				want: want{
					wantN: 0,
					err:   context.Canceled,
				},
				beforeFunc: func(args) {
					cancel()
				},
			}
		}(),
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
			w := &ctxWriter{
				ctx: test.fields.ctx,
				w:   test.fields.w,
			}

			gotN, err := w.Write(test.args.p)
			if err := checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ctxWriter_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		ctx context.Context
		w   io.Writer
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx := context.Background()
			return test{
				name: "success without Closer",
				fields: fields{
					ctx: ctx,
					w:   &bytes.Buffer{},
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "fail when calling cancel function",
				fields: fields{
					ctx: ctx,
					w:   &bytes.Buffer{},
				},
				want: want{
					err: context.Canceled,
				},
				beforeFunc: func() {
					cancel()
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			return test{
				name: "success with Closer",
				fields: fields{
					ctx: ctx,
					w:   &nopWriteCloser{},
				},
				want: want{
					err: nil,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			w := &ctxWriter{
				ctx: test.fields.ctx,
				w:   test.fields.w,
			}

			err := w.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNewEOFReader(t *testing.T) {
// 	type want struct {
// 		want Reader
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, Reader) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got Reader) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := NewEOFReader()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_eofReader_Read(t *testing.T) {
// 	type args struct {
// 		in0 []byte
// 	}
// 	type want struct {
// 		wantN int
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		e          *eofReader
// 		want       want
// 		checkFunc  func(want, int, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN int, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
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
// 		           in0:nil,
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
// 			e := &eofReader{}
//
// 			gotN, err := e.Read(test.args.in0)
// 			if err := checkFunc(test.want, gotN, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func TestReadAll(t *testing.T) {
// 	type args struct {
// 		r Reader
// 	}
// 	type want struct {
// 		wantB []byte
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []byte, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotB []byte, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotB, w.wantB) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotB, w.wantB)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           r:nil,
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
// 		           r:nil,
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
//
// 			gotB, err := ReadAll(test.args.r)
// 			if err := checkFunc(test.want, gotB, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
