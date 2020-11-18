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

package s3

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/reader"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/writer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	}
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want blob.Bucket
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, blob.Bucket, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got blob.Bucket, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		opts := []comparator.Option{
			comparator.AllowUnexported(client{}, aws.Config{}),
			comparator.IgnoreTypes(request.Handlers{}),
			comparator.Comparer(func(want, got errgroup.Group) bool {
				return reflect.DeepEqual(want, got)
			}),
			comparator.Comparer(func(want, got aws.Config) bool {
				return reflect.DeepEqual(want, got)
			}),
			comparator.Comparer(func(want, got *session.Session) bool {
				return reflect.DeepEqual(want, got)
			}),
			comparator.Comparer(func(want, got func(string) (reader.Reader, error)) bool {
				return want != nil && got != nil
			}),
			comparator.Comparer(func(want, got func(string) writer.Writer) bool {
				return want != nil && got != nil
			}),
		}

		if diff := comparator.Diff(w.want, got, opts...); len(diff) != 0 {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns error when option is empty and session field is nil",
				want: want{
					want: nil,
					err:  errors.NewErrInvalidOption("session", nil),
				},
			}
		}(),

		func() test {
			err := errors.New("err")
			opt := func(*client) error {
				return err
			}
			return test{
				name: "returns error when the function field to initialize reader is nil",
				args: args{
					opts: []Option{
						WithSession(func() *session.Session {
							sess, _ := session.NewSession()
							return sess
						}()),
						opt,
					},
				},
				want: want{
					want: nil,
					err:  errors.ErrOptionFailed(err, reflect.ValueOf(opt)),
				},
			}
		}(),

		{
			name: "returns error when the function field to initialize reader is nil",
			args: args{
				opts: []Option{
					WithSession(func() *session.Session {
						sess, _ := session.NewSession()
						return sess
					}()),
					func(c *client) error {
						c.readerFunc = nil
						return nil
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.NewErrInvalidOption("readerFunc", nil),
			},
		},

		{
			name: "returns error when the function field to initialize writer is nil",
			args: args{
				opts: []Option{
					WithSession(func() *session.Session {
						sess, _ := session.NewSession()
						return sess
					}()),
					func(c *client) error {
						c.writerFunc = nil
						return nil
					},
				},
			},
			want: want{
				want: nil,
				err:  errors.NewErrInvalidOption("writerFunc", nil),
			},
		},

		func() test {

			sess, _ := session.NewSession()
			return test{
				name: "returns nil when no error occurs internally",
				args: args{
					opts: []Option{
						WithSession(sess),
					},
				},
				want: want{
					want: &client{
						eg:      errgroup.Get(),
						session: sess,
						service: s3.New(sess),
						readerFunc: func(key string) (reader.Reader, error) {
							return nil, nil
						},
						writerFunc: func(key string) writer.Writer {
							return nil
						},
					},
					err: nil,
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_Open(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		eg          errgroup.Group
		session     *session.Session
		service     *s3.S3
		bucket      string
		maxPartSize int64
		readerFunc  func(key string) (reader.Reader, error)
		writerFunc  func(key string) writer.Writer
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns nil",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				eg:          test.fields.eg,
				session:     test.fields.session,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
				readerFunc:  test.fields.readerFunc,
				writerFunc:  test.fields.writerFunc,
			}

			err := c.Open(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Close(t *testing.T) {
	type fields struct {
		eg          errgroup.Group
		session     *session.Session
		service     *s3.S3
		bucket      string
		maxPartSize int64
		readerFunc  func(key string) (reader.Reader, error)
		writerFunc  func(key string) writer.Writer
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "retursn nil",
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				eg:          test.fields.eg,
				session:     test.fields.session,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
				readerFunc:  test.fields.readerFunc,
				writerFunc:  test.fields.writerFunc,
			}

			err := c.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_client_Reader(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		eg          errgroup.Group
		session     *session.Session
		service     *s3.S3
		bucket      string
		maxPartSize int64
		readerFunc  func(key string) (reader.Reader, error)
		writerFunc  func(key string) writer.Writer
	}
	type want struct {
		want io.ReadCloser
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, io.ReadCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got io.ReadCloser, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns error when there is no function field to initialize the reader",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				want: want{
					want: nil,
					err:  errors.ErrNilObject,
				},
			}
		}(),

		func() test {
			err := errors.New("err")

			return test{
				name: "returns error when the reader initialization fails",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					readerFunc: func(key string) (reader.Reader, error) {
						return nil, err
					},
				},
				want: want{
					want: nil,
					err:  err,
				},
			}
		}(),

		func() test {
			err := errors.New("err")

			r := &reader.MockReader{
				OpenFunc: func(ctx context.Context) error {
					return err
				},
			}
			return test{
				name: "returns error when the open method of reader fails",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					readerFunc: func(key string) (reader.Reader, error) {
						return r, nil
					},
				},
				want: want{
					want: r,
					err:  err,
				},
			}
		}(),

		func() test {
			r := &reader.MockReader{
				OpenFunc: func(ctx context.Context) error {
					return nil
				},
			}
			return test{
				name: "returns nil when no error occurs internally",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					readerFunc: func(key string) (reader.Reader, error) {
						return r, nil
					},
				},
				want: want{
					want: r,
					err:  nil,
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				eg:          test.fields.eg,
				session:     test.fields.session,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
				readerFunc:  test.fields.readerFunc,
				writerFunc:  test.fields.writerFunc,
			}

			got, err := c.Reader(test.args.ctx, test.args.key)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_client_Writer(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	type fields struct {
		eg          errgroup.Group
		session     *session.Session
		service     *s3.S3
		bucket      string
		maxPartSize int64
		readerFunc  func(key string) (reader.Reader, error)
		writerFunc  func(key string) writer.Writer
	}
	type want struct {
		want io.WriteCloser
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, io.WriteCloser, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got io.WriteCloser, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns error when there is no function field to initialize the writer",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				want: want{
					want: nil,
					err:  errors.ErrNilObject,
				},
			}
		}(),

		func() test {
			err := errors.New("err")

			w := &writer.MockWriter{
				OpenFunc: func(ctx context.Context) error {
					return err
				},
			}
			return test{
				name: "returns error when the open method of writer fails",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					writerFunc: func(key string) writer.Writer {
						return w
					},
				},
				want: want{
					want: w,
					err:  err,
				},
			}
		}(),

		func() test {
			w := &writer.MockWriter{
				OpenFunc: func(ctx context.Context) error {
					return nil
				},
			}
			return test{
				name: "returns nil when no error occurs internally",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					writerFunc: func(key string) writer.Writer {
						return w
					},
				},
				want: want{
					want: w,
					err:  nil,
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &client{
				eg:          test.fields.eg,
				session:     test.fields.session,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
				readerFunc:  test.fields.readerFunc,
				writerFunc:  test.fields.writerFunc,
			}

			got, err := c.Writer(test.args.ctx, test.args.key)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
