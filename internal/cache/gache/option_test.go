//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package gache provides implementation of cache using gache

package gache

import (
	"context"
	"reflect"
	"testing"
	"time"

	gache "github.com/kpango/gache/v2"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestDefaultOptions(t *testing.T) {
	type args struct{}
	type want struct {
		want *cache[any]
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *cache[any]) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, got *cache[any]) error {
		opts := []comparator.Option{
			comparator.AllowUnexported(*w.want),
			comparator.AllowUnexported(*got),
			comparator.Comparer(func(want, got *cache[any]) bool {
				return want.gache != nil && got.gache != nil
			}),
		}
		if diff := comparator.Diff(w.want, got, opts...); diff != "" {
			return errors.Errorf("got = %v, want = %v", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "set succuess",
			want: want{
				want: &cache[any]{
					gache: gache.New[any](),
				},
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
			g := new(cache[any])
			for _, opt := range defaultOptions[any]() {
				opt(g)
			}
			if err := checkFunc(test.want, g); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithGache(t *testing.T) {
	type T = cache[any]
	type args struct {
		g gache.Gache[any]
	}
	type want struct {
		want *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, want *T) error {
		if !reflect.DeepEqual(want, w.want) {
			return errors.Errorf("got = %v, want  = %v", want, w.want)
		}
		return nil
	}

	tests := []test{
		func() test {
			ga := gache.New[any]()
			return test{
				name: "set succuess when g is not nil",
				args: args{
					g: ga,
				},
				want: want{
					want: &T{
						gache: ga,
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "set succuess when g is nil",
				want: want{
					want: new(T),
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
			got := WithGache(test.args.g)
			want := new(T)
			got(want)
			if err := checkFunc(test.want, want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithExpiredHook(t *testing.T) {
	type T = cache[any]
	type args struct {
		f func(context.Context, string)
	}
	type want struct {
		want *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, want *T) error {
		if !reflect.DeepEqual(want, w.want) {
			return errors.Errorf("got = %v, want = %v", want, w.want)
		}
		return nil
	}

	tests := []test{
		func() test {
			fn := func(context.Context, string) {}
			return test{
				name: "set succuess when f is not nil",
				args: args{
					f: fn,
				},
				want: want{
					want: &T{
						expiredHook: fn,
					},
				},
				checkFunc: func(w want, g *T) error {
					if reflect.ValueOf(w.want.expiredHook).Pointer() != reflect.ValueOf(g.expiredHook).Pointer() {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", g, w)
					}
					return nil
				},
			}
		}(),
		func() test {
			return test{
				name: "set succuess when fn is nil",
				want: want{
					want: new(T),
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
			got := WithExpiredHook[any](test.args.f)
			want := new(T)
			got(want)
			if err := checkFunc(test.want, want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithExpireDuration(t *testing.T) {
	type T = cache[any]
	type args struct {
		dur time.Duration
	}
	type want struct {
		want *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, want *T) error {
		if !reflect.DeepEqual(want, w.want) {
			return errors.Errorf("got = %v, want = %v", want, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "set succuess when dur is 0",
			args: args{
				dur: 0,
			},
			want: want{
				want: new(T),
			},
		},
		{
			name: "set succuess when dur is not 0",
			args: args{
				dur: 10,
			},
			want: want{
				want: &T{
					expireDur: 10,
				},
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
			got := WithExpireDuration[any](test.args.dur)
			want := new(T)
			got(want)
			if err := checkFunc(test.want, want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithExpireCheckDuration(t *testing.T) {
	type T = cache[any]
	type args struct {
		dur time.Duration
	}
	type want struct {
		want *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, want *T) error {
		if !reflect.DeepEqual(want, w.want) {
			return errors.Errorf("got = %v, want = %v", want, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "set succuess when dur is 0",
			args: args{
				dur: 0,
			},
			want: want{
				want: new(T),
			},
		},
		{
			name: "set succuess when dur is not 0",
			args: args{
				dur: 10,
			},
			want: want{
				want: &T{
					expireCheckDur: 10,
				},
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
			got := WithExpireCheckDuration[any](test.args.dur)
			want := new(T)
			got(want)
			if err := checkFunc(test.want, want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_defaultOptions(t *testing.T) {
// 	type want struct {
// 		want []Option[V]
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, []Option[V]) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got []Option[V]) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := defaultOptions()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
