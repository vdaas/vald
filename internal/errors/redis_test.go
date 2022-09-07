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

// Package errors provides error types and function
package errors

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestErrRedisInvalidKVVKPrefic(t *testing.T) {
	type fields struct {
		kv string
		vk string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	const str = "vdaas"
	tests := []test{
		func() test {
			return test{
				name: "return an ErrRedisInvalidKVVKPrefix error when kv and vk are not empty",
				fields: fields{
					kv: str,
					vk: str,
				},
				want: want{
					want: Errorf("kv index and vk prefix must be defferent.\t(kv: %s,\tvk: %s)", str, str),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an ErrRedisInvalidKVVKPrefix error when kv is not empty and vk is empty",
				fields: fields{
					kv: str,
				},
				want: want{
					want: Errorf("kv index and vk prefix must be defferent.\t(kv: %s,\tvk: %s)", str, ""),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an ErrRedisInvalidKVVKPrefix error when kv is not empty and vk is not empty",
				fields: fields{
					vk: str,
				},
				want: want{
					want: Errorf("kv index and vk prefix must be defferent.\t(kv: %s,\tvk: %s)", "", str),
				},
			}
		}(),
		func() test {
			return test{
				name:   "return an ErrRedisInvalidKVVKPrefix error when kv and vk are empty",
				fields: fields{},
				want: want{
					want: Errorf("kv index and vk prefix must be defferent.\t(kv: %s,\tvk: %s)", "", ""),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRedisInvalidKVVKPrefix(test.fields.kv, test.fields.vk)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewErrRedisNotFoundIdentity(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return an NewErrRedisNotFoundIdentity error",
				want: want{
					want: &ErrRedisNotFoundIdentity{
						err: New("error redis entry not found"),
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := NewErrRedisNotFoundIdentity()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRdisNotFound(t *testing.T) {
	type fields struct {
		key string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrRedisNotFound error when key is not empty",
			fields: fields{
				key: "vdaas",
			},
			want: want{
				want: Wrap(NewErrRedisNotFoundIdentity(), "error redis key 'vdaas' not found"),
			},
		},
		{
			name:   "return an ErrRedisNotFound error when key is empty",
			fields: fields{},
			want: want{
				want: Wrap(NewErrRedisNotFoundIdentity(), "error redis key '' not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRedisNotFound(test.fields.key)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRedisInvalidOption(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrRedisInvalidOption error",
			want: want{
				want: New("error redis invalid option"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRedisInvalidOption
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRedisGetOperationFailed(t *testing.T) {
	type fields struct {
		key string
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	const key = "vdaas"
	err := New("redis error")
	tests := []test{
		func() test {
			return test{
				name: "return a wraped error when key is not empty and err is not nil",
				fields: fields{
					key: key,
					err: err,
				},
				want: want{
					want: Wrapf(err, "Failed to fetch key (%s)", key),
				},
			}
		}(),
		func() test {
			return test{
				name: "return a wraped error when key is not empty and err is nil",
				fields: fields{
					key: key,
				},
				want: want{
					want: Wrapf(nil, "Failed to fetch key (%s)", key),
				},
			}
		}(),
		func() test {
			return test{
				name: "return a wraped error when key is empty and err is not nil",
				fields: fields{
					err: err,
				},
				want: want{
					want: Wrap(err, "Failed to fetch key ()"),
				},
			}
		}(),
		func() test {
			return test{
				name:   "return a wraped error when key is empty and err is nil",
				fields: fields{},
				want: want{
					want: Wrap(nil, "Failed to fetch key ()"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRedisGetOperationFailed(test.fields.key, test.fields.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRedisSetOperationFailed(t *testing.T) {
	type fields struct {
		key string
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	const key = "vdaas"
	err := New("redis error")
	tests := []test{
		func() test {
			return test{
				name: "return a wraped error when key is not empty and err is not nil",
				fields: fields{
					key: key,
					err: err,
				},
				want: want{
					want: Wrapf(err, "Failed to set key (%s)", key),
				},
			}
		}(),
		func() test {
			return test{
				name: "return a wraped error when key is not empty and err is nil",
				fields: fields{
					key: key,
				},
				want: want{
					want: Wrapf(nil, "Failed to set key (%s)", key),
				},
			}
		}(),
		func() test {
			return test{
				name: "return a wraped error when key is empty and err is not nil",
				fields: fields{
					err: err,
				},
				want: want{
					want: Wrap(err, "Failed to set key ()"),
				},
			}
		}(),
		func() test {
			return test{
				name:   "return a wraped error when key is empty and err is nil",
				fields: fields{},
				want: want{
					want: Wrap(nil, "Failed to set key ()"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRedisSetOperationFailed(test.fields.key, test.fields.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRedisDeleteOperationFailed(t *testing.T) {
	type fields struct {
		key string
		err error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	const key = "vdaas"
	err := New("redis error")
	tests := []test{
		func() test {
			return test{
				name: "return a wraped error when key is not empty and err is not nil",
				fields: fields{
					key: key,
					err: err,
				},
				want: want{
					want: Wrapf(err, "Failed to delete key (%s)", key),
				},
			}
		}(),
		func() test {
			return test{
				name: "return a wraped error when key is not empty and err is nil",
				fields: fields{
					key: key,
				},
				want: want{
					want: Wrapf(nil, "Failed to delete key (%s)", key),
				},
			}
		}(),
		func() test {
			return test{
				name: "return a wraped error when key is empty and err is not nil",
				fields: fields{
					err: err,
				},
				want: want{
					want: Wrap(err, "Failed to delete key ()"),
				},
			}
		}(),
		func() test {
			return test{
				name:   "return a wraped error when key is empty and err is nil",
				fields: fields{},
				want: want{
					want: Wrap(nil, "Failed to delete key ()"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRedisDeleteOperationFailed(test.fields.key, test.fields.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidConfigVersion(t *testing.T) {
	type fields struct {
		cur string
		con string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	const cur = "1.0.1"
	const con = "1.1.0"
	tests := []test{
		func() test {
			return test{
				name: "return an config invalid error when cur and con are not empty",
				fields: fields{
					cur: cur,
					con: con,
				},
				want: want{
					want: Errorf("invalid config version %s not satisfies version constraints %s", cur, con),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an config invalid error when cur is empty and con is not empty",
				fields: fields{
					con: con,
				},
				want: want{
					want: Errorf("invalid config version  not satisfies version constraints %s", con),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an config invalid error when cur is not empty and con is empty",
				fields: fields{
					cur: cur,
				},
				want: want{
					want: Errorf("invalid config version %s not satisfies version constraints ", cur),
				},
			}
		}(),
		func() test {
			return test{
				name:   "return an config invalid error when cur and con are empty",
				fields: fields{},
				want: want{
					want: New("invalid config version  not satisfies version constraints "),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrInvalidConfigVersion(test.fields.cur, test.fields.con)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRedisAddrsNotFound(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrRedisAddrsNotFound error",
			want: want{
				want: New("error redis addrs not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRedisAddrsNotFound
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRedisConnectionPingFailed(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrRedisConnectionPingFailed error",
			want: want{
				want: New("error redis connection ping failed"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrRedisConnectionPingFailed
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRedisNotFoundIdentity_Error(t *testing.T) {
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
			name: "Success when err is not nil",
			fields: fields{
				err: New("Not found identity"),
			},
			want: want{
				want: "Not found identity",
			},
		},
		{
			name:   "Success when err is nil",
			fields: fields{},
			want: want{
				want: "expected err is nil: ErrRedisNotFoundIdentity",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			e := &ErrRedisNotFoundIdentity{
				err: test.fields.err,
			}

			got := e.Error()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrRedisNotFoundIdentity_Unwrap(t *testing.T) {
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
		{
			name:   "returns nil when err is nil",
			fields: fields{},
			want:   want{},
		},
		{
			name: "returns err when err is not nil",
			fields: fields{
				err: New("err: Redis not found identity"),
			},
			want: want{
				err: New("err: Redis not found identity"),
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
			e := &ErrRedisNotFoundIdentity{
				err: test.fields.err,
			}

			err := e.Unwrap()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestIsErrRedisNotFound(t *testing.T) {
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
			name: "return false when err is not ErrRedisNotFoundIdentity ",
			args: args{
				err: New("err: Redis not found identity"),
			},
			want: want{},
		},
		{
			name: "return false when err does not match ErrRedisNotFoundIdentity",
			args: args{
				err: &ErrRedisNotFoundIdentity{
					err: New("err: Redis not found identity"),
				},
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false when err is nil",
			args: args{},
			want: want{},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := IsErrRedisNotFound(test.args.err)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
