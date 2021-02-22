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

// Package tcp provides tcp option
package net

import (
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/cache/gache"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/control"
	"go.uber.org/goleak"
)

func TestWithCache(t *testing.T) {
	type T = dialer
	type args struct {
		c cache.Cache
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
		func() test {
			c := gache.New()
			return test{
				name: "set cache success",
				args: args{
					c: c,
				},
				want: want{
					obj: &T{
						cache: c,
					},
				},
			}
		}(),
		{
			name: "set cache to nil success",
			args: args{
				c: nil,
			},
			want: want{
				obj: &T{
					cache: nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithCache(test.args.c)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDNSRefreshDuration(t *testing.T) {
	// Change interface type to the type of object you are testing
	type T = dialer
	type args struct {
		dur string
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
			name: "set success when dur is valid",
			args: args{
				dur: "10s",
			},
			want: want{
				obj: &T{
					dnsRefreshDuration:    10 * time.Second,
					dnsRefreshDurationStr: "10s",
				},
			},
		},
		{
			name: "set success when dur is invalid",
			args: args{
				dur: "dummy",
			},
			want: want{
				obj: &T{
					dnsRefreshDuration:    30 * time.Minute,
					dnsRefreshDurationStr: "30m",
				},
			},
		},
		{
			name: "set success when dur is empty",
			want: want{
				obj: &T{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithDNSRefreshDuration(test.args.dur)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDNSCacheExpiration(t *testing.T) {
	type T = dialer
	type args struct {
		dur string
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
			name: "set success when dur is valid",
			args: args{
				dur: "10s",
			},
			want: want{
				obj: &T{
					dnsCacheExpiration:    10 * time.Second,
					dnsCacheExpirationStr: "10s",
					dnsCache:              true,
				},
			},
		},
		{
			name: "set success when dur is invalid",
			args: args{
				dur: "dummy",
			},
			want: want{
				obj: &T{
					dnsCacheExpiration:    1 * time.Hour,
					dnsCacheExpirationStr: "1h",
					dnsCache:              true,
				},
			},
		},
		{
			name: "set success when dur is empty",
			want: want{
				obj: &T{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithDNSCacheExpiration(test.args.dur)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDialerTimeout(t *testing.T) {
	type T = dialer
	type args struct {
		dur string
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
			name: "set success when dur is valid",
			args: args{
				dur: "10s",
			},
			want: want{
				obj: &T{
					dialerTimeout: 10 * time.Second,
				},
			},
		},
		{
			name: "set success when dur is invalid",
			args: args{
				dur: "dummy",
			},
			want: want{
				obj: &T{
					dialerTimeout: 30 * time.Second,
				},
			},
		},
		{
			name: "set success when dur is empty",
			want: want{
				obj: &T{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithDialerTimeout(test.args.dur)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDialerKeepAlive(t *testing.T) {
	type T = dialer
	type args struct {
		dur string
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
			name: "set success when dur is valid",
			args: args{
				dur: "10s",
			},
			want: want{
				obj: &T{
					dialerKeepAlive: 10 * time.Second,
				},
			},
		},
		{
			name: "set success when dur is invalid",
			args: args{
				dur: "dummy",
			},
			want: want{
				obj: &T{
					dialerKeepAlive: 30 * time.Second,
				},
			},
		},
		{
			name: "set success when dur is empty",
			want: want{
				obj: &T{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithDialerKeepAlive(test.args.dur)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithTLS(t *testing.T) {
	type T = dialer
	type args struct {
		cfg *tls.Config
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
		func() test {
			cfg := new(tls.Config)
			return test{
				name: "set success when cfg is not nil",
				args: args{
					cfg: cfg,
				},
				want: want{
					obj: &T{
						tlsConfig: cfg,
					},
				},
			}
		}(),
		{
			name: "set success when cfg is nil",
			want: want{
				obj: &T{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithTLS(test.args.cfg)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithEnableDNSCache(t *testing.T) {
	type T = dialer
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func()
		afterFunc  func()
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "dnsCache enabled",
			want: want{
				obj: &T{
					dnsCache: true,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithEnableDNSCache()
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDisableDNSCache(t *testing.T) {
	type T = dialer
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func()
		afterFunc  func()
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "dnsCache disabled",
			want: want{
				obj: &T{
					dnsCache: false,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithDisableDNSCache()
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithEnableDialerDualStack(t *testing.T) {
	type T = dialer
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func()
		afterFunc  func()
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "DialerDualStack enabled",
			want: want{
				obj: &T{
					dialerDualStack: true,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithEnableDialerDualStack()
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDisableDialerDualStack(t *testing.T) {
	type T = dialer
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func()
		afterFunc  func()
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "DialerDualStack disabled",
			want: want{
				obj: &T{
					dialerDualStack: false,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithDisableDialerDualStack()
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDialerFallbackDelay(t *testing.T) {
	t.Parallel()
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		dur string
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           dur: "",
		       },
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
		           args: args {
		           dur: "",
		           },
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
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }

			   got := WithDialerFallbackDelay(test.args.dur)
			   obj := new(T)
			   if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }
			   got := WithDialerFallbackDelay(test.args.dur)
			   obj := new(T)
			   got(obj)
			   if err := test.checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithErrGroup(t *testing.T) {
	t.Parallel()
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		eg errgroup.Group
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           eg: nil,
		       },
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
		           args: args {
		           eg: nil,
		           },
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
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }

			   got := WithErrGroup(test.args.eg)
			   obj := new(T)
			   if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }
			   got := WithErrGroup(test.args.eg)
			   obj := new(T)
			   got(obj)
			   if err := test.checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}

func TestWithSocketFlag(t *testing.T) {
	t.Parallel()
	// Change interface type to the type of object you are testing
	type T = interface{}
	type args struct {
		flg control.SocketFlag
	}
	type want struct {
		obj *T
		// Uncomment this line if the option returns an error, otherwise delete it
		// err error
	}
	type test struct {
		name string
		args args
		want want
		// Use the first line if the option returns an error. otherwise use the second line
		// checkFunc  func(want, *T, error) error
		// checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	// Uncomment this block if the option returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T, err error) error {
	       if !errors.Is(err, w.err) {
	           return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
	       }
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	// Uncomment this block if the option do not returns an error, otherwise delete it
	/*
	   defaultCheckFunc := func(w want, obj *T) error {
	       if !reflect.DeepEqual(obj, w.obj) {
	           return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
	       }
	       return nil
	   }
	*/

	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           flg: nil,
		       },
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
		           args: args {
		           flg: nil,
		           },
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
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			// Uncomment this block if the option returns an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }

			   got := WithSocketFlag(test.args.flg)
			   obj := new(T)
			   if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/

			// Uncomment this block if the option do not return an error, otherwise delete it
			/*
			   if test.checkFunc == nil {
			       test.checkFunc = defaultCheckFunc
			   }
			   got := WithSocketFlag(test.args.flg)
			   obj := new(T)
			   got(obj)
			   if err := test.checkFunc(test.want, obj); err != nil {
			       tt.Errorf("error = %v", err)
			   }
			*/
		})
	}
}
