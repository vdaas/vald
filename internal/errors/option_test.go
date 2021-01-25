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
package errors

import (
	"reflect"
	"testing"

	"go.uber.org/goleak"
)

func TestNewErrInvalidOption(t *testing.T) {
	t.Parallel()
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
			var e error
			for _, err := range errs {
				if err == nil {
					continue
				}
				if e != nil {
					e = Wrap(err, e.Error())
				} else {
					e = err
				}
			}
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
			var e error
			for _, err := range errs {
				if err == nil {
					continue
				}
				if e != nil {
					e = Wrap(err, e.Error())
				} else {
					e = err
				}
			}
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
			var e error
			for _, err := range errs {
				if err == nil {
					continue
				}
				if e != nil {
					e = Wrap(err, e.Error())
				} else {
					e = err
				}
			}
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
			var e error
			for _, err := range errs {
				if err == nil {
					continue
				}
				if e != nil {
					e = Wrap(err, e.Error())
				} else {
					e = err
				}
			}
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

	for _, test := range tests {
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
			got := NewErrInvalidOption(test.args.name, test.args.val, test.args.errs...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidOption_Error(t *testing.T) {
	t.Parallel()
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
			name:   "return empty string when e.err is nil.",
			fields: fields{},
			want: want{
				want: "",
			},
		},
		{
			name: "return empty string when e.err is not nil.",
			fields: fields{
				err: Errorf("invalid option. name: WithPort, val: 8080"),
			},
			want: want{
				want: "invalid option. name: WithPort, val: 8080",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ErrInvalidOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Error()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidOption_Unwrap(t *testing.T) {
	t.Parallel()
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ErrInvalidOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Unwrap()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewErrCriticalOption(t *testing.T) {
	t.Parallel()
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
			var e error
			for _, err := range errs {
				if err == nil {
					continue
				}
				if e != nil {
					e = Wrap(err, e.Error())
				} else {
					e = err
				}
			}
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
			var e error
			for _, err := range errs {
				if err == nil {
					continue
				}
				if e != nil {
					e = Wrap(err, e.Error())
				} else {
					e = err
				}
			}
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
			var e error
			for _, err := range errs {
				if err == nil {
					continue
				}
				if e != nil {
					e = Wrap(err, e.Error())
				} else {
					e = err
				}
			}
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
			var e error
			for _, err := range errs {
				if err == nil {
					continue
				}
				if e != nil {
					e = Wrap(err, e.Error())
				} else {
					e = err
				}
			}
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

	for _, test := range tests {
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
			got := NewErrCriticalOption(test.args.name, test.args.val, test.args.errs...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCriticalOption_Error(t *testing.T) {
	t.Parallel()
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
			name:   "return empty string when e.err is nil",
			fields: fields{},
			want: want{
				want: "",
			},
		},
		{
			name: "return empty string when e.err is not nil",
			fields: fields{
				err: Errorf("invalid option. name: WithPort, val: 8080"),
			},
			want: want{
				want: "invalid option. name: WithPort, val: 8080",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ErrCriticalOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Error()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCriticalOption_Unwrap(t *testing.T) {
	t.Parallel()
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ErrCriticalOption{
				err:    test.fields.err,
				origin: test.fields.origin,
			}
			got := e.Unwrap()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
