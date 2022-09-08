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

// Package gache provides implementation of cache using gache
package gache

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestDefaultOptions(t *testing.T) {
	type args struct{}
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
		opts := []cmp.Option{
			cmp.AllowUnexported(*w.want),
			cmp.AllowUnexported(*got),
			cmp.Comparer(func(want, got *cache) bool {
				return want.gache != nil && got.gache != nil
			}),
		}
		if diff := cmp.Diff(w.want, got, opts...); diff != "" {
			return errors.Errorf("got = %v, want = %v", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "set succuess",
			want: want{
				want: &cache{
					gache: gache.New(),
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
			g := new(cache)
			for _, opt := range defaultOptions() {
				opt(g)
			}
			if err := checkFunc(test.want, g); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithGache(t *testing.T) {
	type T = cache
	type args struct {
		g gache.Gache
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
			ga := gache.New()
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
	type T = cache
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
			got := WithExpiredHook(test.args.f)
			want := new(T)
			got(want)
			if err := checkFunc(test.want, want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithExpireDuration(t *testing.T) {
	type T = cache
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
			got := WithExpireDuration(test.args.dur)
			want := new(T)
			got(want)
			if err := checkFunc(test.want, want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithExpireCheckDuration(t *testing.T) {
	type T = cache
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
			got := WithExpireCheckDuration(test.args.dur)
			want := new(T)
			got(want)
			if err := checkFunc(test.want, want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_defaultOptions(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = interface{}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func()
		afterFunc  func()
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got error = %v, want %v", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got = %v, want %v", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got = %v, want %v", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       want: want {
		           obj: new(T),
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           want: want {
		               obj: new(T),
		           },
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }

			   got := defaultOptions()
			   obj := new(T)
			   if err := checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }
			   got := defaultOptions()
			   obj := new(T)
			   got(obj)
			   if err := checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}
