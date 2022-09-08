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

package watch

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestWithErrGroup(t *testing.T) {
	type T = watch
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
			name: "set nothing when eg is nil",
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

			got := WithErrGroup(test.args.eg)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDirs(t *testing.T) {
	type T = watch
	type args struct {
		dirs []string
	}
	type field struct {
		dirs map[string]struct{}
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		field      field
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
			name: "set success when dirs is `vdaas`, `vald`",
			args: args{
				dirs: []string{
					"vdaas", "vald",
				},
			},
			want: want{
				obj: &T{
					dirs: map[string]struct{}{
						"vdaas": {},
						"vald":  {},
					},
				},
			},
		},

		{
			name: "update success when dir field contains map[`team`]struct{}{} and dirs is `vdaas`, `vald`",
			args: args{
				dirs: []string{
					"vdaas", "vald",
				},
			},
			field: field{
				dirs: map[string]struct{}{
					"team": {},
				},
			},
			want: want{
				obj: &T{
					dirs: map[string]struct{}{
						"team":  {},
						"vdaas": {},
						"vald":  {},
					},
				},
			},
		},

		{
			name: "set nothing when dirs is nil",
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

			got := WithDirs(test.args.dirs...)
			obj := &T{
				dirs: test.field.dirs,
			}
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithOnChange(t *testing.T) {
	type T = watch
	type args struct {
		f func(ctx context.Context, name string) error
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
		if reflect.ValueOf(w.obj.onChange).Pointer() != reflect.ValueOf(obj.onChange).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			f := func(context.Context, string) error {
				return nil
			}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: f,
				},
				want: want{
					obj: &watch{
						onChange: f,
					},
				},
			}
		}(),

		{
			name: "set nothing when f is nil",
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

			got := WithOnChange(test.args.f)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithOnCreate(t *testing.T) {
	type T = watch
	type args struct {
		f func(ctx context.Context, name string) error
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
		if reflect.ValueOf(w.obj.onCreate).Pointer() != reflect.ValueOf(obj.onCreate).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			f := func(context.Context, string) error {
				return nil
			}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: f,
				},
				want: want{
					obj: &watch{
						onCreate: f,
					},
				},
			}
		}(),

		{
			name: "set nothing when f is nil",
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

			got := WithOnCreate(test.args.f)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithOnChmod(t *testing.T) {
	type T = watch
	type args struct {
		f func(ctx context.Context, name string) error
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
		if reflect.ValueOf(w.obj.onChmod).Pointer() != reflect.ValueOf(obj.onChmod).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			f := func(context.Context, string) error {
				return nil
			}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: f,
				},
				want: want{
					obj: &watch{
						onChmod: f,
					},
				},
			}
		}(),

		{
			name: "set nothing when f is nil",
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

			got := WithOnChmod(test.args.f)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithOnRename(t *testing.T) {
	type T = watch
	type args struct {
		f func(ctx context.Context, name string) error
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
		if reflect.ValueOf(w.obj.onRename).Pointer() != reflect.ValueOf(obj.onRename).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			f := func(context.Context, string) error {
				return nil
			}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: f,
				},
				want: want{
					obj: &watch{
						onRename: f,
					},
				},
			}
		}(),

		{
			name: "set nothing when f is nil",
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

			got := WithOnRename(test.args.f)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithOnDelete(t *testing.T) {
	type T = watch
	type args struct {
		f func(ctx context.Context, name string) error
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
		if reflect.ValueOf(w.obj.onDelete).Pointer() != reflect.ValueOf(obj.onDelete).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			f := func(context.Context, string) error {
				return nil
			}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: f,
				},
				want: want{
					obj: &watch{
						onDelete: f,
					},
				},
			}
		}(),

		{
			name: "set nothing when f is nil",
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

			got := WithOnDelete(test.args.f)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithOnWrite(t *testing.T) {
	type T = watch
	type args struct {
		f func(ctx context.Context, name string) error
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
		if reflect.ValueOf(w.obj.onWrite).Pointer() != reflect.ValueOf(obj.onWrite).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			f := func(context.Context, string) error {
				return nil
			}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: f,
				},
				want: want{
					obj: &watch{
						onWrite: f,
					},
				},
			}
		}(),

		{
			name: "set nothing when f is nil",
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

			got := WithOnWrite(test.args.f)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithOnError(t *testing.T) {
	type T = watch
	type args struct {
		f func(ctx context.Context, err error) error
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
		if reflect.ValueOf(w.obj.onError).Pointer() != reflect.ValueOf(obj.onError).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			f := func(context.Context, error) error {
				return nil
			}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: f,
				},
				want: want{
					obj: &watch{
						onError: f,
					},
				},
			}
		}(),

		{
			name: "set nothing when f is nil",
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

			got := WithOnError(test.args.f)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
