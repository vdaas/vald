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

package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestWithErrGroup(t *testing.T) {
	type T = ngt
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
			ctx, cancel := context.WithCancel(context.Background())
			eg, _ := errgroup.New(ctx)

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
				afterFunc: func(a args) {
					cancel()
				},
			}
		}(),
		{
			name: "return nil when eg is nil",
			args: args{
				eg: nil,
			},
			want: want{
				obj: new(T),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithErrGroup(test.args.eg)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithEnableInMemoryMode(t *testing.T) {
	type T = ngt
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
			name: "success to set in memory mode to true",
			args: args{
				enabled: true,
			},
			want: want{
				obj: &T{
					inMem: true,
				},
			},
		},
		{
			name: "success to set in memory mode to false",
			args: args{
				enabled: false,
			},
			want: want{
				obj: &T{
					inMem: false,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithEnableInMemoryMode(test.args.enabled)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithIndexPath(t *testing.T) {
	type T = ngt
	type args struct {
		path string
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
			name: "set success when path is not empty",
			args: args{
				path: "/var/vald",
			},
			want: want{
				obj: &T{
					path: "/var/vald",
				},
			},
		},
		{
			name: "set success when path is empty",
			args: args{
				path: "",
			},
			want: want{
				obj: &T{
					path: "",
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithIndexPath(test.args.path)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithAutoIndexCheckDuration(t *testing.T) {
	type T = ngt
	type args struct {
		dur string
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
			name: "set success when duration is empty string",
			args: args{
				dur: "",
			},
			want: want{
				obj: &T{
					dur: 0,
				},
			},
		},
		{
			name: "set success when duration is a valid duration string",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					dur: 5 * time.Second,
				},
			},
		},
		{
			name: "return error when duration is not a valid duration string",
			args: args{
				dur: "5ss",
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid timeout value: 5ss\t:timeout parse error out put failed: time: unknown unit \"ss\" in duration \"5ss\""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithAutoIndexCheckDuration(test.args.dur)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithAutoIndexDurationLimit(t *testing.T) {
	type T = ngt
	type args struct {
		dur string
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
			name: "set success when duration is empty string",
			args: args{
				dur: "",
			},
			want: want{
				obj: &T{
					lim: 0,
				},
			},
		},
		{
			name: "set success when duration is a valid duration string",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					lim: 5 * time.Second,
				},
			},
		},
		{
			name: "return error when duration is not a valid duration string",
			args: args{
				dur: "5ss",
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid timeout value: 5ss\t:timeout parse error out put failed: time: unknown unit \"ss\" in duration \"5ss\""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithAutoIndexDurationLimit(test.args.dur)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithAutoSaveIndexDuration(t *testing.T) {
	type T = ngt
	type args struct {
		dur string
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
			name: "set success when duration is empty string",
			args: args{
				dur: "",
			},
			want: want{
				obj: &T{
					sdur: 0,
				},
			},
		},
		{
			name: "set success when duration is a valid duration string",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					sdur: 5 * time.Second,
				},
			},
		},
		{
			name: "return error when duration is not a valid duration string",
			args: args{
				dur: "5ss",
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid timeout value: 5ss\t:timeout parse error out put failed: time: unknown unit \"ss\" in duration \"5ss\""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithAutoSaveIndexDuration(test.args.dur)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithAutoIndexLength(t *testing.T) {
	type T = ngt
	type args struct {
		l int
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
			name: "set index length success",
			args: args{
				l: 10,
			},
			want: want{
				obj: &T{
					alen: 10,
				},
			},
		},
		{
			name: "set index length success when length is 0",
			args: args{
				l: 0,
			},
			want: want{
				obj: &T{
					alen: 0,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithAutoIndexLength(test.args.l)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithInitialDelayMaxDuration(t *testing.T) {
	type T = ngt
	type args struct {
		dur string
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
			name: "set success when duration is empty string",
			args: args{
				dur: "",
			},
			want: want{
				obj: &T{
					idelay: 0,
				},
			},
		},
		{
			name: "set success when duration is a valid duration string",
			args: args{
				dur: "5s",
			},
			checkFunc: func(w want, t *T, e error) error {
				if t.idelay == 0 {
					return errors.New("delay value is 0")
				}
				return nil
			},
		},
		{
			name: "return error when duration is not a valid duration string",
			args: args{
				dur: "5ss",
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid timeout value: 5ss\t:timeout parse error out put failed: time: unknown unit \"ss\" in duration \"5ss\""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithInitialDelayMaxDuration(test.args.dur)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMinLoadIndexTimeout(t *testing.T) {
	type T = ngt
	type args struct {
		dur string
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
			name: "set success when duration is empty string",
			args: args{
				dur: "",
			},
			want: want{
				obj: &T{
					minLit: 0,
				},
			},
		},
		{
			name: "set success when duration is a valid duration string",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					minLit: 5 * time.Second,
				},
			},
		},
		{
			name: "return error when duration is not a valid duration string",
			args: args{
				dur: "5ss",
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid timeout value: 5ss\t:timeout parse error out put failed: time: unknown unit \"ss\" in duration \"5ss\""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithMinLoadIndexTimeout(test.args.dur)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxLoadIndexTimeout(t *testing.T) {
	type T = ngt
	type args struct {
		dur string
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
			name: "set success when duration is empty string",
			args: args{
				dur: "",
			},
			want: want{
				obj: &T{
					maxLit: 0,
				},
			},
		},
		{
			name: "set success when duration is a valid duration string",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					maxLit: 5 * time.Second,
				},
			},
		},
		{
			name: "return error when duration is not a valid duration string",
			args: args{
				dur: "5ss",
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid timeout value: 5ss\t:timeout parse error out put failed: time: unknown unit \"ss\" in duration \"5ss\""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithMaxLoadIndexTimeout(test.args.dur)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithLoadIndexTimeoutFactor(t *testing.T) {
	type T = ngt
	type args struct {
		dur string
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
			name: "set success when duration is empty string",
			args: args{
				dur: "",
			},
			want: want{
				obj: &T{
					litFactor: 0,
				},
			},
		},
		{
			name: "set success when duration is a valid duration string",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					litFactor: 5 * time.Second,
				},
			},
		},
		{
			name: "return error when duration is not a valid duration string",
			args: args{
				dur: "5ss",
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid timeout value: 5ss\t:timeout parse error out put failed: time: unknown unit \"ss\" in duration \"5ss\""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithLoadIndexTimeoutFactor(test.args.dur)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDefaultPoolSize(t *testing.T) {
	type T = ngt
	type args struct {
		ps uint32
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
			name: "set pool size success",
			args: args{
				ps: 50,
			},
			want: want{
				obj: &T{
					poolSize: 50,
				},
			},
		},
		{
			name: "set success when pool size is 0",
			args: args{
				ps: 0,
			},
			want: want{
				obj: &T{
					poolSize: 0,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithDefaultPoolSize(test.args.ps)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDefaultRadius(t *testing.T) {
	type T = ngt
	type args struct {
		rad float32
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
			name: "set default radius success",
			args: args{
				rad: 0.4,
			},
			want: want{
				obj: &T{
					radius: 0.4,
				},
			},
		},
		{
			name: "set success when radius is 0",
			args: args{
				rad: 0,
			},
			want: want{
				obj: &T{
					radius: 0,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithDefaultRadius(test.args.rad)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDefaultEpsilon(t *testing.T) {
	type T = ngt
	type args struct {
		epsilon float32
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
			name: "set default epsilon success",
			args: args{
				epsilon: 50,
			},
			want: want{
				obj: &T{
					epsilon: 50,
				},
			},
		},
		{
			name: "set success when epsilon is 0",
			args: args{
				epsilon: 0,
			},
			want: want{
				obj: &T{
					epsilon: 0,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithDefaultEpsilon(test.args.epsilon)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithProactiveGC(t *testing.T) {
	type T = ngt
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
			name: "set proactive GC success",
			args: args{
				enabled: true,
			},
			want: want{
				obj: &T{
					enableProactiveGC: true,
				},
			},
		},
		{
			name: "set proactive GC success when it is false",
			args: args{
				enabled: false,
			},
			want: want{
				obj: &T{
					enableProactiveGC: false,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithProactiveGC(test.args.enabled)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithCopyOnWrite(t *testing.T) {
	type T = ngt
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
			name: "set CoW success",
			args: args{
				enabled: true,
			},
			want: want{
				obj: &T{
					enableCopyOnWrite: true,
				},
			},
		},
		{
			name: "set CoW when it is false",
			args: args{
				enabled: false,
			},
			want: want{
				obj: &T{
					enableCopyOnWrite: false,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := defaultCheckFunc
			if test.checkFunc != nil {
				checkFunc = test.checkFunc
			}

			got := WithCopyOnWrite(test.args.enabled)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
