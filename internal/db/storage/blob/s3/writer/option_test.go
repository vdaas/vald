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

package writer

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

var (
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	}
)

func TestWithErrGroup(t *testing.T) {
	type T = writer
	type args struct {
		eg errgroup.Group
	}
	type fields struct {
		eg errgroup.Group
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when eg is not nil",
			args: args{
				eg: errgroup.Get(),
			},
			want: want{
				obj: &T{
					eg: errgroup.Get(),
				},
			},
		},

		{
			name: "returns error when eg is nil",
			args: args{
				eg: nil,
			},
			fields: fields{
				eg: errgroup.Get(),
			},
			want: want{
				obj: &T{
					eg: errgroup.Get(),
				},
				err: errors.ErrInvalidOption("errgroup", nil),
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
			got := WithErrGroup(test.args.eg)
			obj := &T{
				eg: test.fields.eg,
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithService(t *testing.T) {
	type T = writer
	type args struct {
		s *s3.S3
	}
	type fields struct {
		service *s3.S3
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			s := new(s3.S3)
			return test{
				name: "set success when s is not nil",
				args: args{
					s: s,
				},
				want: want{
					obj: &T{
						service: s,
					},
				},
			}
		}(),

		func() test {
			s := new(s3.S3)
			var ss *s3.S3
			return test{
				name: "returns error when s is nil",
				args: args{
					s: ss,
				},
				fields: fields{
					service: s,
				},
				want: want{
					obj: &T{
						service: s,
					},
					err: errors.ErrInvalidOption("service", ss),
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
			got := WithService(test.args.s)
			obj := &T{
				service: test.fields.service,
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBucket(t *testing.T) {
	type T = writer
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when bucket is not empty",
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
			name: "returns error when bucket is empty",
			args: args{
				bucket: "",
			},
			want: want{
				obj: new(T),
				err: errors.ErrInvalidOption("bucket", ""),
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
			got := WithBucket(test.args.bucket)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithKey(t *testing.T) {
	type T = writer
	type args struct {
		key string
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when key is not empty",
			args: args{
				key: "key",
			},
			want: want{
				obj: &T{
					key: "key",
				},
			},
		},

		{
			name: "returns error when key is empty",
			args: args{
				key: "",
			},
			want: want{
				obj: new(T),
				err: errors.ErrInvalidOption("key", ""),
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
			got := WithKey(test.args.key)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxPartSize(t *testing.T) {
	type T = writer
	type args struct {
		max int64
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when max is larger then MinUploadPartSize",
			args: args{
				max: s3manager.MinUploadPartSize + 1000,
			},
			want: want{
				obj: &T{
					maxPartSize: s3manager.MinUploadPartSize + 1000,
				},
			},
		},

		{
			name: "returns error when max is smaller then MinUploadPartSize",
			args: args{
				max: 10,
			},
			want: want{
				obj: new(T),
				err: errors.ErrInvalidOption("maxPartSize", 10),
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
			got := WithMaxPartSize(test.args.max)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithContentType(t *testing.T) {
	type T = writer
	type args struct {
		ct string
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when ct is not empty",
			args: args{
				ct: "utf8",
			},
			want: want{
				obj: &T{
					contentType: "utf8",
				},
			},
		},

		{
			name: "returns error when ct is empty",
			args: args{
				ct: "",
			},
			want: want{
				obj: new(T),
				err: errors.ErrInvalidOption("contentType", ""),
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

			got := WithContentType(test.args.ct)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
