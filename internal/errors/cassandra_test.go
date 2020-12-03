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

// Package errors provides error types and function
package errors

import (
	"reflect"
	"testing"

	"go.uber.org/goleak"
)

func TestErrCassandraNotFoundIdentity_Error(t *testing.T) {
	type fields struct {
		err error
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns string when internal error is cassandra not found identity error",
			fields: fields{
				err: New("cassandra not found identity"),
			},
			want: want{
				want: "cassandra not found identity",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ErrCassandraNotFoundIdentity{
				err: test.fields.err,
			}

			got := e.Error()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraNotFoundIdentity_Unwrap(t *testing.T) {
	type fields struct {
		err error
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
		if !Is(err, w.err) {
			return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns nil when internal error is nil",
				fields: fields{
					err: nil,
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			err := New("cassandra not found identity")
			return test{
				name: "returns internal error when internal error is cassandra not found identity error",
				fields: fields{
					err: err,
				},
				want: want{
					err: err,
				},
				checkFunc: func(w want, err error) error {
					if !reflect.DeepEqual(w.err, err) {
						return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ErrCassandraNotFoundIdentity{
				err: test.fields.err,
			}

			err := e.Unwrap()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIsErrCassandraNotFound(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return false when error is not cassandra not found identity",
			args: args{
				err: New("database not found"),
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return true when error is cassandra not found identity",
			args: args{
				err: new(ErrCassandraNotFoundIdentity),
			},
			want: want{
				want: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := IsErrCassandraNotFound(test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraUnavailableIdentity_Error(t *testing.T) {
	type fields struct {
		err error
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns string when internal error is cassandra unavailable identity error",
			fields: fields{
				err: New("cassandra unavailable identity"),
			},
			want: want{
				want: "cassandra unavailable identity",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ErrCassandraUnavailableIdentity{
				err: test.fields.err,
			}

			got := e.Error()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrCassandraUnavailableIdentity_Unwrap(t *testing.T) {
	type fields struct {
		err error
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
		if !Is(err, w.err) {
			return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns nil when internal error is nil",
				fields: fields{
					err: nil,
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			err := New("cassandra unavailable identity")
			return test{
				name: "returns internal error when internal error is cassandra unavailable identity error",
				fields: fields{
					err: err,
				},
				want: want{
					err: err,
				},
				checkFunc: func(w want, err error) error {
					if !reflect.DeepEqual(w.err, err) {
						return Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ErrCassandraUnavailableIdentity{
				err: test.fields.err,
			}

			err := e.Unwrap()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIsErrCassandraUnavailable(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return false when error is not cassandra unavailable identity",
			args: args{
				err: New("database not found"),
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return true when error is cassandra unavailable identity",
			args: args{
				err: new(ErrCassandraUnavailableIdentity),
			},
			want: want{
				want: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := IsErrCassandraUnavailable(test.args.err)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
