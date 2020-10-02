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

		if (w.want == nil && got != nil) || (w.want != nil && got == nil) {
			return errors.Errorf("got: %v, want: %v", w.want, got)
		} else if w.want == nil && got == nil {
			return nil
		}

		wantC, gotC := w.want.(*client), got.(*client)

		clientComparator := []comparator.Option{
			comparator.AllowUnexported(*gotC),
			comparator.Comparer(func(x, y *session.Session) bool {
				return x != nil && y != nil
			}),
			comparator.Comparer(func(x, y *s3.S3) bool {
				return x != nil && y != nil
			}),
			comparator.Comparer(func(x, y errgroup.Group) bool {
				return reflect.DeepEqual(x, y)
			}),
		}

		if diff := comparator.Diff(*wantC, *gotC, clientComparator...); diff != "" {
			return errors.New(diff)
		}

		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns error when option is empty and s3 session is nil",
				args: args{
					opts: nil,
				},
				want: want{
					want: nil,
					err:  errors.ErrS3SessionNotFound,
				},
			}
		}(),

		func() test {
			sess, _ := session.NewSession()
			return test{
				name: "returns error when option is not empty",
				args: args{
					opts: []Option{
						WithSession(sess),
					},
				},
				want: want{
					want: &client{
						eg:         errgroup.Get(),
						session:    sess,
						service:    s3.New(sess),
						readWriter: newRW(),
					},
					err: nil,
				},
			}
		}(),

		func() test {
			opts := []Option{
				func(c *client) error {
					return errors.New("err")
				},
			}
			return test{
				name: "returns error when option is not empty",
				args: args{
					opts: opts,
				},
				want: want{
					want: nil,
					err:  errors.ErrOptionFailed(errors.New("err"), reflect.ValueOf(opts[0])),
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
			name: "returns nil",
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
		readWriter  readWriter
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
			r := &MockReader{
				OpenFunc: func(ctx context.Context) error {
					return nil
				},
			}
			return test{
				name: "returns (ReadCloser, nil) when Open successes",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					readWriter: &MockRW{
						NewReaderFunc: func(opts ...reader.Option) reader.Reader {
							return r
						},
					},
				},
				want: want{
					want: r,
					err:  nil,
				},
			}
		}(),

		func() test {
			r := &MockReader{
				OpenFunc: func(ctx context.Context) error {
					return errors.New("err")
				},
			}
			return test{
				name: "returns (ReadCloser, error) when Open fails",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					readWriter: &MockRW{
						NewReaderFunc: func(opts ...reader.Option) reader.Reader {
							return r
						},
					},
				},
				want: want{
					want: r,
					err:  errors.New("err"),
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
				readWriter:  test.fields.readWriter,
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
		readWriter  readWriter
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
			w := &MockWriter{
				OpenFunc: func(ctx context.Context) error {
					return nil
				},
			}
			return test{
				name: "returns (ReadCloser, nil) when Open successes",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					readWriter: &MockRW{
						NewWriterFunc: func(opts ...writer.Option) writer.Writer {
							return w
						},
					},
				},
				want: want{
					want: w,
					err:  nil,
				},
			}
		}(),

		func() test {
			w := &MockWriter{
				OpenFunc: func(ctx context.Context) error {
					return errors.New("err")
				},
			}
			return test{
				name: "returns (WriteCloser, error) when Open fails",
				args: args{
					ctx: context.Background(),
					key: "key",
				},
				fields: fields{
					readWriter: &MockRW{
						NewWriterFunc: func(opts ...writer.Option) writer.Writer {
							return w
						},
					},
				},
				want: want{
					want: w,
					err:  errors.New("err"),
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
				readWriter:  test.fields.readWriter,
			}

			got, err := c.Writer(test.args.ctx, test.args.key)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_newRW(t *testing.T) {
	type want struct {
		want readWriter
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, readWriter) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got readWriter) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns readWriter",
			want: want{
				want: new(rw),
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

			got := newRW()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_rw_NewReader(t *testing.T) {
	type args struct {
		opts []reader.Option
	}
	type want struct {
		want reader.Reader
	}
	type test struct {
		name       string
		args       args
		r          *rw
		want       want
		checkFunc  func(want, reader.Reader) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got reader.Reader) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns reader.Reader",
			args: args{
				opts: nil,
			},
			want: want{
				want: reader.New(),
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
			r := &rw{}

			got := r.NewReader(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_rw_NewWriter(t *testing.T) {
	type args struct {
		opts []writer.Option
	}
	type want struct {
		want writer.Writer
	}
	type test struct {
		name       string
		args       args
		r          *rw
		want       want
		checkFunc  func(want, writer.Writer) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got writer.Writer) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns writer.Writer",
			args: args{
				opts: nil,
			},
			want: want{
				want: writer.New(),
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
			r := &rw{}

			got := r.NewWriter(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
