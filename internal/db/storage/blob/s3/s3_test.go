//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/reader"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3manager"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/writer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

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
			comparator.Comparer(func(want, got reader.Reader) bool {
				return want != nil && got != nil
			}),
			comparator.Comparer(func(want, got writer.Writer) bool {
				return want != nil && got != nil
			}),
			comparator.Comparer(func(want, got []backoff.Option) bool {
				return reflect.DeepEqual(want, got)
			}),
			comparator.Comparer(func(want, got s3manager.S3Manager) bool {
				return reflect.DeepEqual(want, got)
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
			opt := func(c *client) error {
				return errors.New("err")
			}
			return test{
				name: "returns error when option apply fails",
				args: args{
					opts: []Option{
						opt,
					},
				},
				want: want{
					want: nil,
					err:  errors.ErrOptionFailed(errors.New("err"), reflect.ValueOf(opt)),
				},
			}
		}(),
		func() test {
			sess, _ := session.NewSession()
			r := new(reader.MockReader)
			w := new(writer.MockWriter)
			return test{
				name: "returns bucket and nil when the option apply success and no error occurs internally",
				args: args{
					opts: []Option{
						WithSession(sess),
						WithReader(r),
						WithWriter(w),
					},
				},
				want: want{
					want: &client{
						eg:      errgroup.Get(),
						session: sess,
						service: s3.New(sess),
						reader:  r,
						writer:  w,
					},
					err: nil,
				},
			}
		}(),
		func() test {
			sess, _ := session.NewSession()
			service := s3.New(sess)
			eg := errgroup.Get()
			writer := new(writer.MockWriter)
			return test{
				name: "returns bucket and nil when reader is created and no error occurs internally",
				args: args{
					opts: []Option{
						WithSession(sess),
						WithErrGroup(eg),
						WithBucket("bucket"),
						WithMaxPartSize("100G"),
						func(c *client) error {
							c.writer = writer
							return nil
						},
					},
				},
				want: want{
					want: &client{
						eg:          eg,
						session:     sess,
						service:     service,
						bucket:      "bucket",
						maxPartSize: 107374182400,
						reader: func() (r reader.Reader) {
							r, _ = reader.New(
								reader.WithErrGroup(eg),
								reader.WithService(service),
								reader.WithBucket("bucket"),
								reader.WithMaxChunkSize(107374182400),
							)
							return
						}(),
						writer: writer,
					},
					err: nil,
				},
			}
		}(),
		func() test {
			sess, _ := session.NewSession()
			service := s3.New(sess)
			eg := errgroup.Get()
			reader := new(reader.MockReader)
			return test{
				name: "returns bucket and nil when writer is created and no error occurs internally",
				args: args{
					opts: []Option{
						WithSession(sess),
						WithErrGroup(eg),
						WithBucket("bucket"),
						WithMaxPartSize("100G"),
						func(c *client) error {
							c.reader = reader
							return nil
						},
					},
				},
				want: want{
					want: &client{
						eg:          eg,
						session:     sess,
						service:     service,
						bucket:      "bucket",
						maxPartSize: 107374182400,
						reader:      reader,
						writer: writer.New(
							writer.WithErrGroup(eg),
							writer.WithService(service),
							writer.WithBucket("bucket"),
							writer.WithMaxPartSize(107374182400),
						),
					},
					err: nil,
				},
			}
		}(),
		func() test {
			sess, _ := session.NewSession()
			service := s3.New(sess)
			eg := errgroup.Get()
			return test{
				name: "returns bucket and nil when reader and writer are created and no error occurs internally",
				args: args{
					opts: []Option{
						WithSession(sess),
						WithErrGroup(eg),
						WithBucket("bucket"),
						WithMaxPartSize("100G"),
					},
				},
				want: want{
					want: &client{
						eg:          eg,
						session:     sess,
						service:     service,
						bucket:      "bucket",
						maxPartSize: 107374182400,
						reader: func() (r reader.Reader) {
							r, _ = reader.New(
								reader.WithErrGroup(eg),
								reader.WithService(service),
								reader.WithBucket("bucket"),
								reader.WithMaxChunkSize(107374182400),
							)
							return
						}(),
						writer: writer.New(
							writer.WithErrGroup(eg),
							writer.WithService(service),
							writer.WithBucket("bucket"),
							writer.WithMaxPartSize(107374182400),
						),
					},
					err: nil,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
			c := &client{
				eg:          test.fields.eg,
				session:     test.fields.session,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
			}

			err := c.Open(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
			c := &client{
				eg:          test.fields.eg,
				session:     test.fields.session,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
			}

			err := c.Close()
			if err := checkFunc(test.want, err); err != nil {
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
		reader      reader.Reader
		writer      writer.Writer
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
			opened := false

			r := &reader.MockReader{
				OpenFunc: func(ctx context.Context, key string) error {
					opened = true
					return nil
				},
			}

			return test{
				name: "returns opened reader and nil when open method of reader success",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					reader: r,
				},
				want: want{
					want: r,
					err:  nil,
				},
				checkFunc: func(w want, g io.ReadCloser, gerr error) error {
					err := defaultCheckFunc(w, g, gerr)
					if err != nil {
						return err
					}

					if !opened {
						return errors.New("reader is not opened")
					}

					return nil
				},
			}
		}(),
		func() test {
			err := errors.New("err")
			opened := false
			r := &reader.MockReader{
				OpenFunc: func(ctx context.Context, key string) error {
					opened = true
					return err
				},
			}
			return test{
				name: "returns opened reader and error from the open method of reader",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					reader: r,
				},
				want: want{
					want: nil,
					err:  err,
				},
				checkFunc: func(w want, g io.ReadCloser, gerr error) error {
					err := defaultCheckFunc(w, g, gerr)
					if err != nil {
						return err
					}

					if !opened {
						return errors.New("reader is not opened")
					}

					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
			c := &client{
				eg:          test.fields.eg,
				session:     test.fields.session,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
				reader:      test.fields.reader,
				writer:      test.fields.writer,
			}

			got, err := c.Reader(test.args.ctx, test.args.key)
			if err := checkFunc(test.want, got, err); err != nil {
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
		reader      reader.Reader
		writer      writer.Writer
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
			opened := false
			w := &writer.MockWriter{
				OpenFunc: func(ctx context.Context, key string) error {
					opened = true
					return nil
				},
			}
			return test{
				name: "returns opened writer and nil when open method of writer success",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					writer: w,
				},
				want: want{
					want: w,
					err:  nil,
				},
				checkFunc: func(w want, g io.WriteCloser, gerr error) error {
					err := defaultCheckFunc(w, g, gerr)
					if err != nil {
						return err
					}

					if !opened {
						return errors.New("writer is not opened")
					}

					return nil
				},
			}
		}(),
		func() test {
			err := errors.New("err")
			opened := false
			w := &writer.MockWriter{
				OpenFunc: func(ctx context.Context, key string) error {
					opened = true
					return err
				},
			}
			return test{
				name: "returns opened writer and error from the open method of writer",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					writer: w,
				},
				want: want{
					want: nil,
					err:  err,
				},
				checkFunc: func(w want, g io.WriteCloser, gerr error) error {
					err := defaultCheckFunc(w, g, gerr)
					if err != nil {
						return err
					}

					if !opened {
						return errors.New("writer is not opened")
					}

					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
			c := &client{
				eg:          test.fields.eg,
				session:     test.fields.session,
				service:     test.fields.service,
				bucket:      test.fields.bucket,
				maxPartSize: test.fields.maxPartSize,
				reader:      test.fields.reader,
				writer:      test.fields.writer,
			}

			got, err := c.Writer(test.args.ctx, test.args.key)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
