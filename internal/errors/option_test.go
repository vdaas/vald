// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package errors

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNewErrInvalidOption(t *testing.T) {
	type args struct {
		name string
		val  interface{}
		errs []error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		func() test {
			name := "WithPort"
			val := 9000
			return test{
				name: "return ErrInvalidOpton when name and val have a value and errs is empty.",
				args: args{
					name: name,
					val:  val,
				},
				want: want{
					want: &ErrInvalidOption{
						err: Errorf("invalid option, name: %s, val: %v", name, val),
					},
				},
			}
		}(),
		func() test {
			name := "WithPort"
			val := 9000
			errs := []error{
				New("set option failed"),
			}
			e := errs[0]
			return test{
				name: "return ErrInvalidOpton when all of parameter has value.",
				args: args{
					name: name,
					val:  val,
					errs: errs,
				},
				want: want{
					want: &ErrInvalidOption{
						err:    Wrapf(e, "invalid option, name: %s, val: %v", name, val),
						origin: e,
					},
				},
			}
		}(),
		func() test {
			name := "WithPort"
			val := 9000
			errs := []error{
				nil,
				New("set option failed."),
			}
			e := errs[1]
			return test{
				name: "return ErrInvalidOpton when all of parameter has value and errs has nil as value.",
				args: args{
					name: name,
					val:  val,
					errs: errs,
				},
				want: want{
					want: &ErrInvalidOption{
						err:    Wrapf(e, "invalid option, name: %s, val: %v", name, val),
						origin: e,
					},
				},
			}
		}(),
		func() test {
			val := 9000
			errs := []error{
				New("set option failed."),
				New("name is nil."),
			}
			e := Wrap(errs[1], errs[0].Error())
			return test{
				name: "return ErrInvalidOpton when name is nil and val and errs have values.",
				args: args{
					val:  val,
					errs: errs,
				},
				want: want{
					want: &ErrInvalidOption{
						err:    Wrapf(e, "invalid option, name: , val: %v", val),
						origin: e,
					},
				},
			}
		}(),
		func() test {
			name := "WithPort"
			errs := []error{
				New("set option failed."),
				New("val is nil."),
			}
			e := Wrap(errs[1], errs[0].Error())
			return test{
				name: "return ErrInvalidOpton when val is nil and name and errs have values.",
				args: args{
					name: name,
					errs: errs,
				},
				want: want{
					want: &ErrInvalidOption{
						err:    Wrapf(e, "invalid option, name: %s, val: %v", name, nil),
						origin: e,
					},
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
			got := NewErrInvalidOption(test.args.name, test.args.val, test.args.errs...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidOption_Error(t *testing.T) {
	type T = string
	type fields struct {
		err    error
		origin error
	}
	type want struct {
		want T
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, T) error
		beforeFunc func()
		afterFunc  func()
	}

	defaultCheckFunc := func(w want, obj T) error {
		if !reflect.DeepEqual(obj, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "return string when e.err is not nil.",
			fields: fields{
				err: New("invalid option. name: WithPort, val: 8080"),
			},
			want: want{
				want: "invalid option. name: WithPort, val: 8080",
			},
		},
		{
			name:   "return string when e.err is nil.",
			fields: fields{},
			want: want{
				want: "expected err is nil: ErrInvalidOption",
			},
		},
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
			e := &ErrInvalidOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Error()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidOption_Unwrap(t *testing.T) {
	type T = error
	type fields struct {
		err    error
		origin error
	}
	type want struct {
		want T
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, T) error
		beforeFunc func()
		afterFunc  func()
	}

	defaultCheckFunc := func(w want, got T) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name:   "return nil when origin is nil.",
			fields: fields{},
			want:   want{},
		},
		{
			name: "return nil when origin is not nil.",
			fields: fields{
				origin: Wrap(New("invalid options"), "WithHost"),
			},
			want: want{
				want: New("WithHost: invalid options"),
			},
		},
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
			e := &ErrInvalidOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Unwrap()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewErrCriticalOption(t *testing.T) {
	type T = error
	type args struct {
		name string
		val  interface{}
		errs []error
	}
	type want struct {
		want T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, got T) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		func() test {
			name := "WithPort"
			val := 9000
			return test{
				name: "return ErrCriticalOption when name and val have a value and errs is empty.",
				args: args{
					name: name,
					val:  val,
				},
				want: want{
					want: &ErrCriticalOption{
						err: Errorf("invalid critical option, name: %s, val: %v", name, val),
					},
				},
			}
		}(),
		func() test {
			name := "WithPort"
			val := 9000
			errs := []error{
				New("set option failed"),
			}
			e := errs[0]
			return test{
				name: "return ErrCriticalOption when all of parameter has value.",
				args: args{
					name: name,
					val:  val,
					errs: errs,
				},
				want: want{
					want: &ErrCriticalOption{
						err:    Wrapf(e, "invalid critical option, name: %s, val: %v", name, val),
						origin: e,
					},
				},
			}
		}(),
		func() test {
			name := "WithPort"
			val := 9000
			errs := []error{
				nil,
				New("set option failed."),
			}
			e := errs[1]
			return test{
				name: "return ErrCriticalOption when all of parameter has value and errs has nil as value.",
				args: args{
					name: name,
					val:  val,
					errs: errs,
				},
				want: want{
					want: &ErrCriticalOption{
						err:    Wrapf(e, "invalid critical option, name: %s, val: %v", name, val),
						origin: e,
					},
				},
			}
		}(),
		func() test {
			val := 9000
			errs := []error{
				New("set option failed."),
				New("name is nil."),
			}
			e := Wrap(errs[1], errs[0].Error())
			return test{
				name: "return ErrCriticalOption when name is nil and val and errs have values.",
				args: args{
					val:  val,
					errs: errs,
				},
				want: want{
					want: &ErrCriticalOption{
						err:    Wrapf(e, "invalid critical option, name: , val: %v", val),
						origin: e,
					},
				},
			}
		}(),
		func() test {
			name := "WithPort"
			errs := []error{
				New("set option failed."),
				New("val is nil."),
			}
			e := Wrap(errs[1], errs[0].Error())
			return test{
				name: "return ErrCriticalOption when val is nil and name and errs have values.",
				args: args{
					name: name,
					errs: errs,
				},
				want: want{
					want: &ErrCriticalOption{
						err:    Wrapf(e, "invalid critical option, name: %s, val: %v", name, nil),
						origin: e,
					},
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
			got := NewErrCriticalOption(test.args.name, test.args.val, test.args.errs...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCriticalOption_Error(t *testing.T) {
	type T = string
	type fields struct {
		err    error
		origin error
	}
	type want struct {
		want T
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, T) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got T) error {
		if !reflect.DeepEqual(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "return string when e.err is not nil",
			fields: fields{
				err: New("invalid option. name: WithPort, val: 8080"),
			},
			want: want{
				want: "invalid option. name: WithPort, val: 8080",
			},
		},
		{
			name:   "return string when e.err is nil",
			fields: fields{},
			want: want{
				want: "expected err is nil: ErrCriticalOption",
			},
		},
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
			e := &ErrCriticalOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Error()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCriticalOption_Unwrap(t *testing.T) {
	type T = error
	type fields struct {
		err    error
		origin error
	}
	type want struct {
		want T
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, T) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got T) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name:   "return nil when origin is nil.",
			fields: fields{},
			want:   want{},
		},
		{
			name: "return nil when origin is not nil.",
			fields: fields{
				origin: Wrap(New("invalid options"), "WithHost"),
			},
			want: want{
				want: New("WithHost: invalid options"),
			},
		},
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
			e := &ErrCriticalOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Unwrap()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewErrIgnoredOption(t *testing.T) {
	type args struct {
		name string
		errs []error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		func() test {
			name := "WithPort"
			return test{
				name: "return ErrIgnoredOption when name and val have a value and errs is empty.",
				args: args{
					name: name,
				},
				want: want{
					want: &ErrIgnoredOption{
						err: Errorf("ignored option, name: %s", name),
					},
				},
			}
		}(),
		func() test {
			name := "WithPort"
			errs := []error{
				New("set option failed"),
			}
			e := errs[0]
			return test{
				name: "return ErrIgnoredOption when all of parameter has value.",
				args: args{
					name: name,
					errs: errs,
				},
				want: want{
					want: &ErrIgnoredOption{
						err:    Wrapf(e, "ignored option, name: %s", name),
						origin: e,
					},
				},
			}
		}(),
		func() test {
			name := "WithPort"
			errs := []error{
				nil,
				New("set option failed."),
			}
			e := errs[1]
			return test{
				name: "return ErrIgnoredOption when all of parameter has value and errs has nil as value.",
				args: args{
					name: name,
					errs: errs,
				},
				want: want{
					want: &ErrIgnoredOption{
						err:    Wrapf(e, "ignored option, name: %s", name),
						origin: e,
					},
				},
			}
		}(),
		func() test {
			errs := []error{
				New("set option failed."),
				New("name is nil."),
			}
			e := Wrap(errs[1], errs[0].Error())
			return test{
				name: "return ErrIgnoredOption when name is nil and val and errs have values.",
				args: args{
					errs: errs,
				},
				want: want{
					want: &ErrIgnoredOption{
						err:    Wrapf(e, "ignored option, name: "),
						origin: e,
					},
				},
			}
		}(),
		func() test {
			name := "WithPort"
			errs := []error{
				New("set option failed."),
				New("val is nil."),
			}
			e := Wrap(errs[1], errs[0].Error())
			return test{
				name: "return ErrIgnoredOption when val is nil and name and errs have values.",
				args: args{
					name: name,
					errs: errs,
				},
				want: want{
					want: &ErrIgnoredOption{
						err:    Wrapf(e, "ignored option, name: %s", name),
						origin: e,
					},
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
			got := NewErrIgnoredOption(test.args.name, test.args.errs...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrIgnoredOption_Error(t *testing.T) {
	type T = string
	type fields struct {
		err    error
		origin error
	}
	type want struct {
		want T
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, T) error
		beforeFunc func()
		afterFunc  func()
	}

	defaultCheckFunc := func(w want, obj T) error {
		if !reflect.DeepEqual(obj, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "return string when e.err is not nil.",
			fields: fields{
				err: New("ignored option. name: WithPort"),
			},
			want: want{
				want: "ignored option. name: WithPort",
			},
		},
		{
			name:   "return string when e.err is nil.",
			fields: fields{},
			want: want{
				want: "expected err is nil: ErrIgnoredOption",
			},
		},
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
			e := &ErrIgnoredOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Error()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrIgnoredOption_Unwrap(t *testing.T) {
	type T = error
	type fields struct {
		err    error
		origin error
	}
	type want struct {
		want T
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, T) error
		beforeFunc func()
		afterFunc  func()
	}

	defaultCheckFunc := func(w want, got T) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name:   "return nil when origin is nil.",
			fields: fields{},
			want:   want{},
		},
		{
			name: "return nil when origin is not nil.",
			fields: fields{
				origin: Wrap(New("ignored options"), "WithHost"),
			},
			want: want{
				want: New("WithHost: ignored options"),
			},
		},
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
			e := &ErrIgnoredOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Unwrap()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
