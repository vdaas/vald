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
package zap

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestWithLevel(t *testing.T) {
	type T = logger
	type args struct {
		lv string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "do nothing if lv is empty string",
			args: args{
				lv: "",
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "if lv is debug, DEBUG level will be set",
			args: args{
				lv: "debug",
			},
			want: want{
				obj: &T{
					level: level.DEBUG,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			got := WithLevel(test.args.lv)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithFormat(t *testing.T) {
	type T = logger
	type args struct {
		fmt string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "do nothing if fmt is empty string",
			args: args{
				fmt: "",
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "if fmt is json, JSON format will be set",
			args: args{
				fmt: "json",
			},
			want: want{
				obj: &T{
					format: format.JSON,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			got := WithFormat(test.args.fmt)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithCaller(t *testing.T) {
	type T = logger
	type args struct {
		enable bool
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "given value will be set to enableCaller field",
			args: args{
				enable: true,
			},
			want: want{
				obj: &T{
					enableCaller: true,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			got := WithCaller(test.args.enable)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
