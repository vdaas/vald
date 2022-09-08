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
	stderrs "errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/go-cmp/cmp"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/reader"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/writer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestWithErrGroup(t *testing.T) {
	type T = client
	type args struct {
		eg errgroup.Group
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			eg := errgroup.Get()
			return test{
				name: "set success when eg is not nil",
				args: args{
					eg: eg,
				},
				want: want{
					obj: &T{
						eg: eg,
					},
				},
			}
		}(),

		func() test {
			return test{
				name: "set nothing when eg is nil",
				args: args{
					eg: nil,
				},
				want: want{
					obj: &T{
						eg: nil,
					},
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

			got := WithErrGroup(test.args.eg)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithSession(t *testing.T) {
	type T = client
	type args struct {
		sess *session.Session
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			sess := new(session.Session)
			return test{
				name: "set success when sess is not nil",
				args: args{
					sess: sess,
				},
				want: want{
					obj: &T{
						session: sess,
					},
				},
			}
		}(),

		func() test {
			return test{
				name: "set nothing when sess is not nil",
				args: args{
					sess: nil,
				},
				want: want{
					obj: &T{
						session: nil,
					},
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

			got := WithSession(test.args.sess)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBucket(t *testing.T) {
	type T = client
	type args struct {
		bucket string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when bucket is `bucket`",
			args: args{
				bucket: "bucket",
			},
			want: want{
				obj: &T{
					bucket: "bucket",
				},
			},
		},

		{
			name: "set success when bucket is empty",
			args: args{
				bucket: "",
			},
			want: want{
				obj: new(T),
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

			got := WithBucket(test.args.bucket)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxPartSize(t *testing.T) {
	type T = client
	type args struct {
		size string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set nothing when size is `100G`",
			args: args{
				size: "100G",
			},
			want: want{
				obj: &T{
					maxPartSize: 107374182400,
				},
			},
		},

		{
			name: "set nothing when size is `1M`",
			args: args{
				size: "1M",
			},
			want: want{
				obj: &T{
					maxPartSize: 0,
				},
			},
		},

		{
			name: "set default and returns error when size given with invalid string",
			args: args{
				size: "a",
			},
			want: want{
				obj: &T{
					maxPartSize: 0,
				},
				err: func() (err error) {
					err = stderrs.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")
					err = errors.Wrap(err, errors.ErrParseUnitFailed("a").Error())
					return
				}(),
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

			got := WithMaxPartSize(test.args.size)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxChunkSize(t *testing.T) {
	type T = client
	type args struct {
		size string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when size is `100G`",
			args: args{
				size: "100G",
			},
			want: want{
				obj: &T{
					maxChunkSize: 107374182400,
				},
				err: nil,
			},
		},

		{
			name: "set nothing when size is `1M`",
			args: args{
				size: "1M",
			},
			want: want{
				obj: &T{
					maxChunkSize: 0,
				},
				err: nil,
			},
		},

		{
			name: "returns error when size is `a`",
			args: args{
				size: "a",
			},
			want: want{
				obj: &T{
					maxChunkSize: 0,
				},
				err: func() (err error) {
					err = stderrs.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")
					err = errors.Wrap(err, errors.ErrParseUnitFailed("a").Error())
					return
				}(),
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

			got := WithMaxChunkSize(test.args.size)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithReaderBackoff(t *testing.T) {
	type T = client
	type args struct {
		enabled bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when enabled is `true`",
			args: args{
				enabled: true,
			},
			want: want{
				obj: &T{
					readerBackoffEnabled: true,
				},
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

			got := WithReaderBackoff(test.args.enabled)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithReaderBackoffOpts(t *testing.T) {
	type T = client
	type args struct {
		opts []backoff.Option
	}
	type fields struct {
		readerBackoffOpts []backoff.Option
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		opts := []cmp.Option{
			cmp.AllowUnexported(*obj),
			cmp.AllowUnexported(*w.obj),
			cmp.Comparer(func(want, got []backoff.Option) bool {
				return len(got) == len(want)
			}),
			cmp.Comparer(func(want, got backoff.Option) bool {
				return reflect.ValueOf(got).Pointer() == reflect.ValueOf(want).Pointer()
			}),
		}
		if diff := cmp.Diff(w.obj, obj, opts...); diff != "" {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}

		return nil
	}

	tests := []test{
		{
			name: "set success when opts is not nil and c.readerBackoffOpts is not nil",
			args: args{
				opts: []backoff.Option{
					backoff.WithRetryCount(1),
				},
			},
			fields: fields{
				readerBackoffOpts: []backoff.Option{
					backoff.WithRetryCount(1),
				},
			},
			want: want{
				obj: &T{
					readerBackoffOpts: []backoff.Option{
						backoff.WithRetryCount(1),
						backoff.WithRetryCount(1),
					},
				},
				err: nil,
			},
		},

		{
			name: "set success when opts is not nil and r.readerBackoffOpts is nil",
			args: args{
				opts: []backoff.Option{
					backoff.WithRetryCount(1),
				},
			},
			want: want{
				obj: &T{
					readerBackoffOpts: []backoff.Option{
						backoff.WithRetryCount(1),
					},
				},
				err: nil,
			},
		},

		{
			name: "set nothing when opts is nil",
			args: args{
				opts: nil,
			},
			want: want{
				obj: new(T),
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

			got := WithReaderBackoffOpts(test.args.opts...)
			obj := &T{
				readerBackoffOpts: test.fields.readerBackoffOpts,
			}
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithReader(t *testing.T) {
	type T = client
	type args struct {
		r reader.Reader
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			r := new(reader.MockReader)
			return test{
				name: "set reader success",
				args: args{
					r: r,
				},
				want: want{
					obj: &T{
						reader: r,
					},
				},
			}
		}(),

		func() test {
			return test{
				name: "set nothing when reader is the default value",
				args: args{
					r: nil,
				},
				want: want{
					obj: new(T),
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

			got := WithReader(test.args.r)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithWriter(t *testing.T) {
	type T = client
	type args struct {
		w writer.Writer
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			w := new(writer.MockWriter)
			return test{
				name: "set writer success",
				args: args{
					w: w,
				},
				want: want{
					obj: &T{
						writer: w,
					},
				},
			}
		}(),

		func() test {
			return test{
				name: "set nothing when writer is the default value",
				args: args{
					w: nil,
				},
				want: want{
					obj: new(T),
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

			got := WithWriter(test.args.w)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
