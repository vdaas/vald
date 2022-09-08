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

// Package cache provides implementation of cache
package cache

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/cache/cacher"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/internal/timeutil"
)

func TestWithExpiredHook(t *testing.T) {
	type args struct {
		f func(context.Context, string)
	}
	type want struct {
		want *cache
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *cache) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *cache) error {
		if reflect.ValueOf(w.want.expiredHook).Pointer() != reflect.ValueOf(got.expiredHook).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			fn := func(context.Context, string) {}
			return test{
				name: "set success when f is not nil",
				args: args{
					f: fn,
				},
				want: want{
					want: &cache{
						expiredHook: fn,
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "set success when f is nil",
				args: args{
					f: nil,
				},
				want: want{
					want: &cache{},
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

			got := new(cache)
			opts := WithExpiredHook(test.args.f)
			opts(got)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithType(t *testing.T) {
	type args struct {
		mo string
	}
	type want struct {
		want *cache
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *cache) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *cache) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			val := "gache"
			return test{
				name: "set success when len(mo) is not 0",
				args: args{
					mo: val,
				},
				want: want{
					want: &cache{
						cacher: cacher.ToType(val),
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "set success when len(mo) is 0",
				want: want{
					want: &cache{},
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

			got := new(cache)
			opts := WithType(test.args.mo)
			opts(got)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithExpireDuration(t *testing.T) {
	type args struct {
		dur string
	}
	type want struct {
		want *cache
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *cache) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *cache) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			val := "10s"
			dur, _ := timeutil.Parse(val)
			return test{
				name: "set success when dur is legal parameter",
				args: args{
					dur: val,
				},
				want: want{
					want: &cache{
						expireDur: dur,
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "set success when dur is empty",
				want: want{
					want: &cache{},
				},
			}
		}(),
		func() test {
			val := "invalid"
			return test{
				name: "set success when dur is invalid",
				args: args{
					dur: val,
				},
				want: want{
					want: &cache{},
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

			got := new(cache)
			opts := WithExpireDuration(test.args.dur)
			opts(got)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithExpireCheckDuration(t *testing.T) {
	type args struct {
		dur string
	}
	type want struct {
		want *cache
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *cache) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *cache) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			val := "10s"
			dur, _ := timeutil.Parse(val)
			return test{
				name: "set success when dur is legal parameter",
				args: args{
					dur: val,
				},
				want: want{
					want: &cache{
						expireCheckDur: dur,
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "set success when dur is empty",
				want: want{
					want: &cache{},
				},
			}
		}(),
		func() test {
			val := "invalid"
			return test{
				name: "set success when dur is invalid",
				args: args{
					dur: val,
				},
				want: want{
					want: &cache{},
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

			got := new(cache)
			opts := WithExpireCheckDuration(test.args.dur)
			opts(got)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
