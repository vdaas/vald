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

// Package runner provides implementation of process runner
package runner

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestWithName(t *testing.T) {
	type T = runner
	type args struct {
		name string
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
			name: "set success when name is `vald`",
			args: args{
				name: "vald",
			},
			want: want{
				obj: &T{
					name: "vald",
				},
			},
		},

		{
			name: "set nothing when name is empty",
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
			got := WithName(test.args.name)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithVersion(t *testing.T) {
	type T = runner
	type args struct {
		ver string
		max string
		min string
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
			name: "set success when ver is `v1.1.3` and max is `v1.1.5` min is `v1.1.2`",
			args: args{
				ver: "v1.1.3",
				max: "v.1.1.5",
				min: "v.1.1.2",
			},
			want: want{
				obj: &T{
					version:    "v1.1.3",
					maxVersion: "v.1.1.5",
					minVersion: "v.1.1.2",
				},
			},
		},

		{
			name: "set success when ver is empty and max is `v1.1.5` min is `v1.1.2`",
			args: args{
				max: "v.1.1.5",
				min: "v.1.1.2",
			},
			want: want{
				obj: &T{
					maxVersion: "v.1.1.5",
					minVersion: "v.1.1.2",
				},
			},
		},

		{
			name: "set success when ver is `v1.1.3` and max is empty min is `v1.1.2`",
			args: args{
				ver: "v1.1.3",
				min: "v.1.1.2",
			},
			want: want{
				obj: &T{
					version:    "v1.1.3",
					minVersion: "v.1.1.2",
				},
			},
		},

		{
			name: "set success when ver is `v1.1.3` and max is `v1.1.5` min is empty",
			args: args{
				ver: "v1.1.3",
				max: "v.1.1.5",
			},
			want: want{
				obj: &T{
					version:    "v1.1.3",
					maxVersion: "v.1.1.5",
				},
			},
		},

		{
			name: "set nothing when ver and max and min are empty",
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
			got := WithVersion(test.args.ver, test.args.max, test.args.min)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithConfigLoader(t *testing.T) {
	type T = runner
	type args struct {
		f func(string) (interface{}, *config.GlobalConfig, error)
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
		if reflect.ValueOf(w.obj.loadConfig).Pointer() != reflect.ValueOf(obj.loadConfig).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			f := func(string) (interface{}, *config.GlobalConfig, error) {
				return nil, nil, nil
			}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: f,
				},
				want: want{
					obj: &T{
						loadConfig: f,
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
			got := WithConfigLoader(test.args.f)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDaemonInitializer(t *testing.T) {
	type T = runner
	type args struct {
		f func(interface{}) (Runner, error)
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
		if reflect.ValueOf(w.obj.loadConfig).Pointer() != reflect.ValueOf(obj.loadConfig).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			f := func(interface{}) (Runner, error) {
				return nil, nil
			}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: f,
				},
				want: want{
					obj: &T{
						initializeDaemon: f,
					},
				},
			}
		}(),

		{
			name: "set success when f is nil",
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
			got := WithDaemonInitializer(test.args.f)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
