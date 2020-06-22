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

// Package tcp provides tcp option
package tcp

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/cache/gache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"

	"go.uber.org/goleak"
)

var (
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
		goleak.IgnoreTopFunction("github.com/kpango/gache.(*gache).StartExpired.func1"),
	}
)

func TestNewDialer(t *testing.T) {
	type args struct {
		opts []DialerOption
	}
	type want struct {
		wantDer Dialer
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Dialer, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotDer Dialer, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %+v, want %+v", err, w.err)
		}

		want := w.wantDer.(*dialer)
		got := gotDer.(*dialer)
		opts := []cmp.Option{cmp.AllowUnexported(*want),
			cmpopts.IgnoreFields(*want, "dialer", "der"),
			cmp.Comparer(func(x, y cache.Cache) bool {
				if x == nil && y == nil {
					return true
				}
				return !(x == nil && y != nil) || !(y == nil && x != nil)
			}),
			cmp.Comparer(func(x, y tls.Config) bool {
				return reflect.DeepEqual(x, y)
			}),
		}

		if diff := cmp.Diff(*want, *got, opts...); diff != "" {
			return errors.Errorf("err: %s", diff)
		}
		if got.dialer == nil {
			return errors.Errorf("dialer is not initialized")
		}
		if got.der == nil {
			return errors.Errorf("der is not initialized")
		}
		return nil
	}
	tests := []test{
		func() test {
			d := &net.Dialer{
				Timeout:   time.Second * 30,
				KeepAlive: time.Second * 30,
				DualStack: true,
				Control:   Control,
			}
			d.Resolver = &net.Resolver{
				PreferGo: false,
				Dial:     d.DialContext,
			}

			return test{
				name: "returns dialer when option is empty",
				want: want{
					wantDer: &dialer{
						dialerKeepAlive: time.Second * 30,
						dialerTimeout:   time.Second * 30,
						dialerDualStack: true,
						der:             d,
						dialer:          d.DialContext,
					},
				},
			}
		}(),
		func() test {
			d := &net.Dialer{
				Timeout:   time.Second * 30,
				KeepAlive: time.Second * 30,
				DualStack: true,
				Control:   Control,
			}
			d.Resolver = &net.Resolver{
				PreferGo: false,
				Dial:     d.DialContext,
			}
			tc := new(tls.Config)
			c := gache.New()

			return test{
				name: "returns dialer when option is not empty",
				args: args{
					opts: []DialerOption{
						WithTLS(tc),
						WithCache(c),
						WithEnableDNSCache(),
						WithDNSRefreshDuration("5s"),
						WithDNSCacheExpiration("10s"),
					},
				},
				want: want{
					wantDer: &dialer{
						dialerKeepAlive:       time.Second * 30,
						dialerTimeout:         time.Second * 30,
						dnsRefreshDuration:    time.Second * 5,
						dnsCacheExpiration:    time.Second * 10,
						dnsRefreshDurationStr: "5s",
						dnsCacheExpirationStr: "10s",
						dnsCache:              true,
						dialerDualStack:       true,
						der:                   d,
						dialer:                d.DialContext,
						cache:                 c,
						tlsConfig:             tc,
					},
				},
			}
		}(),
		func() test {
			tc := new(tls.Config)

			d := &net.Dialer{
				Timeout:   time.Second * 30,
				KeepAlive: time.Second * 30,
				DualStack: true,
				Control:   Control,
			}
			d.Resolver = &net.Resolver{
				PreferGo: false,
				Dial:     d.DialContext,
			}

			return test{
				name: "returns dialer when tls option is not empty and connection confirmation succeeds",
				args: args{
					opts: []DialerOption{WithTLS(tc)},
				},
				want: want{
					wantDer: &dialer{
						dialerKeepAlive: time.Second * 30,
						dialerTimeout:   time.Second * 30,
						dialerDualStack: true,
						der:             d,
						dialer:          d.DialContext,
						tlsConfig:       tc,
					},
				},
				checkFunc: func(w want, gotDer Dialer, err error) error {
					if err := defaultCheckFunc(w, gotDer, err); err != nil {
						return err
					}

					f := gotDer.GetDialer()
					conn, err := f(context.Background(), "tcp", "google.com:80")
					if err != nil {
						return errors.Errorf("err is not nil: %v", err)
					}
					if conn == nil {
						return errors.Errorf("conn is nil")
					}

					return nil
				},
			}
		}(),
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

			gotDer, err := NewDialer(test.args.opts...)
			if err := test.checkFunc(test.want, gotDer, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dialer_GetDialer(t *testing.T) {
	type fields struct {
		dialer func(ctx context.Context, network, addr string) (net.Conn, error)
	}
	type want struct {
		want func(ctx context.Context, network, addr string) (net.Conn, error)
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, func(ctx context.Context, network, addr string) (net.Conn, error)) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got func(ctx context.Context, network, addr string) (net.Conn, error)) error {
		if reflect.ValueOf(w.want).Pointer() != reflect.ValueOf(got).Pointer() {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			fn := func(ctx context.Context, network, addr string) (net.Conn, error) {
				return nil, nil
			}
			return test{
				name: "get success",
				fields: fields{
					dialer: fn,
				},
				want: want{
					want: fn,
				},
			}
		}(),
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
			d := &dialer{
				dialer: test.fields.dialer,
			}

			got := d.GetDialer()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dialer_lookup(t *testing.T) {
	type args struct {
		ctx  context.Context
		addr string
	}
	type fields struct {
		cache                 cache.Cache
		dnsCache              bool
		tlsConfig             *tls.Config
		dnsRefreshDurationStr string
		dnsCacheExpirationStr string
		dnsRefreshDuration    time.Duration
		dnsCacheExpiration    time.Duration
		dialerTimeout         time.Duration
		dialerKeepAlive       time.Duration
		dialerDualStack       bool
		der                   *net.Dialer
		dialer                func(ctx context.Context, network, addr string) (net.Conn, error)
	}
	type want struct {
		wantIps []string
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []string, error, *dialer) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotIps []string, err error, d *dialer) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotIps, w.wantIps) {
			return errors.Errorf("got = %v, want %v", gotIps, w.wantIps)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ips when lookupIpAddr returns ips",
			args: args{
				ctx:  context.Background(),
				addr: "google.com",
			},
			fields: fields{
				cache: gache.New(),
				der: &net.Dialer{
					Resolver: &net.Resolver{
						PreferGo: false,
					},
				},
			},
			checkFunc: func(w want, gotIps []string, err error, d *dialer) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got error = %v, want %v", err, w.err)
				}

				if len(gotIps) == 0 {
					return errors.New("ips is empty")
				}
				return nil
			},
		},
		{
			name: "return ips when lookupIpAddr returns ips and the cache is set",
			args: args{
				ctx:  context.Background(),
				addr: "google.com",
			},
			fields: fields{
				cache: gache.New(),
				der: &net.Dialer{
					Resolver: &net.Resolver{
						PreferGo: false,
					},
				},
			},
			checkFunc: func(w want, gotIps []string, err error, d *dialer) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got error = %v, want %v", err, w.err)
				}

				if len(gotIps) == 0 {
					return errors.New("ips is empty")
				}

				// check the cache is set
				if _, ok := d.cache.Get("google.com"); !ok {
					return errors.New("cache is not set")
				}

				// execute lookup again and check the result is the same
				gotIps1, err1 := d.lookup(context.Background(), "google.com")
				if err1 != nil {
					return err1
				}
				if !reflect.DeepEqual(gotIps, gotIps1) {
					return errors.Errorf("got = %v, got1 %v", gotIps, gotIps)
				}

				// check the cache is set
				if _, ok := d.cache.Get("google.com"); !ok {
					return errors.New("cache is not set")
				}

				return nil
			},
		},
		{
			name: "return cached ips when the cache hits",
			args: args{
				ctx:  context.Background(),
				addr: "addr",
			},
			fields: fields{
				cache: func() cache.Cache {
					g := gache.New()
					g.Set("addr", []string{"999.999.999.999"})
					return g
				}(),
				der: &net.Dialer{
					Resolver: &net.Resolver{
						PreferGo: false,
					},
				},
			},
			want: want{
				wantIps: []string{"999.999.999.999"},
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
			d := &dialer{
				cache:                 test.fields.cache,
				dnsCache:              test.fields.dnsCache,
				tlsConfig:             test.fields.tlsConfig,
				dnsRefreshDurationStr: test.fields.dnsRefreshDurationStr,
				dnsCacheExpirationStr: test.fields.dnsCacheExpirationStr,
				dnsRefreshDuration:    test.fields.dnsRefreshDuration,
				dnsCacheExpiration:    test.fields.dnsCacheExpiration,
				dialerTimeout:         test.fields.dialerTimeout,
				dialerKeepAlive:       test.fields.dialerKeepAlive,
				dialerDualStack:       test.fields.dialerDualStack,
				der:                   test.fields.der,
				dialer:                test.fields.dialer,
			}

			gotIps, err := d.lookup(test.args.ctx, test.args.addr)
			if err := test.checkFunc(test.want, gotIps, err, d); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dialer_StartDialerCache(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		cache                 cache.Cache
		dnsCache              bool
		tlsConfig             *tls.Config
		dnsRefreshDurationStr string
		dnsCacheExpirationStr string
		dnsRefreshDuration    time.Duration
		dnsCacheExpiration    time.Duration
		dialerTimeout         time.Duration
		dialerKeepAlive       time.Duration
		dialerDualStack       bool
		der                   *net.Dialer
		dialer                func(ctx context.Context, network, addr string) (net.Conn, error)
	}
	type want struct {
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(*dialer) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(d *dialer) error {
		return nil
	}
	tests := []test{
		func() test {
			addr := "google.com"
			ips := []string{"0.0.0.0"}
			cache, _ := cache.New(cache.WithExpireDuration("500ms"), cache.WithExpireCheckDuration("500ms"))

			return test{
				name: "hook is called when expired",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					dnsCache:           true,
					cache:              cache,
					dnsRefreshDuration: time.Millisecond * 500,
					dnsCacheExpiration: time.Millisecond * 500,
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						}},
				},
				checkFunc: func(d *dialer) error {
					time.Sleep(time.Second)

					val, _ := d.cache.Get(addr)
					if reflect.DeepEqual(val, ips) {
						return errors.New("cache is not cleared")
					}
					return nil
				},
			}
		}(),
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
			d := &dialer{
				cache:                 test.fields.cache,
				dnsCache:              test.fields.dnsCache,
				tlsConfig:             test.fields.tlsConfig,
				dnsRefreshDurationStr: test.fields.dnsRefreshDurationStr,
				dnsCacheExpirationStr: test.fields.dnsCacheExpirationStr,
				dnsRefreshDuration:    test.fields.dnsRefreshDuration,
				dnsCacheExpiration:    test.fields.dnsCacheExpiration,
				dialerTimeout:         test.fields.dialerTimeout,
				dialerKeepAlive:       test.fields.dialerKeepAlive,
				dialerDualStack:       test.fields.dialerDualStack,
				der:                   test.fields.der,
				dialer:                test.fields.dialer,
			}

			d.StartDialerCache(test.args.ctx)
			if err := test.checkFunc(d); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialer_DialContext(t *testing.T) {
	type args struct {
		ctx     context.Context
		network string
		address string
	}
	type fields struct {
		cache                 cache.Cache
		dnsCache              bool
		tlsConfig             *tls.Config
		dnsRefreshDurationStr string
		dnsCacheExpirationStr string
		dnsRefreshDuration    time.Duration
		dnsCacheExpiration    time.Duration
		dialerTimeout         time.Duration
		dialerKeepAlive       time.Duration
		dialerDualStack       bool
		der                   *net.Dialer
		dialer                func(ctx context.Context, network, addr string) (net.Conn, error)
	}
	type want struct {
		want net.Conn
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, net.Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got net.Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           network: "",
		           address: "",
		       },
		       fields: fields {
		           cache: nil,
		           dnsCache: false,
		           tlsConfig: nil,
		           dnsRefreshDurationStr: "",
		           dnsCacheExpirationStr: "",
		           dnsRefreshDuration: nil,
		           dnsCacheExpiration: nil,
		           dialerTimeout: nil,
		           dialerKeepAlive: nil,
		           dialerDualStack: false,
		           der: nil,
		           dialer: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx: nil,
		           network: "",
		           address: "",
		           },
		           fields: fields {
		           cache: nil,
		           dnsCache: false,
		           tlsConfig: nil,
		           dnsRefreshDurationStr: "",
		           dnsCacheExpirationStr: "",
		           dnsRefreshDuration: nil,
		           dnsCacheExpiration: nil,
		           dialerTimeout: nil,
		           dialerKeepAlive: nil,
		           dialerDualStack: false,
		           der: nil,
		           dialer: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			d := &dialer{
				cache:                 test.fields.cache,
				dnsCache:              test.fields.dnsCache,
				tlsConfig:             test.fields.tlsConfig,
				dnsRefreshDurationStr: test.fields.dnsRefreshDurationStr,
				dnsCacheExpirationStr: test.fields.dnsCacheExpirationStr,
				dnsRefreshDuration:    test.fields.dnsRefreshDuration,
				dnsCacheExpiration:    test.fields.dnsCacheExpiration,
				dialerTimeout:         test.fields.dialerTimeout,
				dialerKeepAlive:       test.fields.dialerKeepAlive,
				dialerDualStack:       test.fields.dialerDualStack,
				der:                   test.fields.der,
				dialer:                test.fields.dialer,
			}

			got, err := d.DialContext(test.args.ctx, test.args.network, test.args.address)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dialer_cachedDialer(t *testing.T) {
	type args struct {
		dctx    context.Context
		network string
		addr    string
	}
	type fields struct {
		cache                 cache.Cache
		dnsCache              bool
		tlsConfig             *tls.Config
		dnsRefreshDurationStr string
		dnsCacheExpirationStr string
		dnsRefreshDuration    time.Duration
		dnsCacheExpiration    time.Duration
		dialerTimeout         time.Duration
		dialerKeepAlive       time.Duration
		dialerDualStack       bool
		der                   *net.Dialer
		dialer                func(ctx context.Context, network, addr string) (net.Conn, error)
	}
	type want struct {
		wantConn net.Conn
		err      error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, net.Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotConn net.Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotConn, w.wantConn) {
			return errors.Errorf("got = %v, want %v", gotConn, w.wantConn)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           dctx: nil,
		           network: "",
		           addr: "",
		       },
		       fields: fields {
		           cache: nil,
		           dnsCache: false,
		           tlsConfig: nil,
		           dnsRefreshDurationStr: "",
		           dnsCacheExpirationStr: "",
		           dnsRefreshDuration: nil,
		           dnsCacheExpiration: nil,
		           dialerTimeout: nil,
		           dialerKeepAlive: nil,
		           dialerDualStack: false,
		           der: nil,
		           dialer: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           dctx: nil,
		           network: "",
		           addr: "",
		           },
		           fields: fields {
		           cache: nil,
		           dnsCache: false,
		           tlsConfig: nil,
		           dnsRefreshDurationStr: "",
		           dnsCacheExpirationStr: "",
		           dnsRefreshDuration: nil,
		           dnsCacheExpiration: nil,
		           dialerTimeout: nil,
		           dialerKeepAlive: nil,
		           dialerDualStack: false,
		           der: nil,
		           dialer: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			d := &dialer{
				cache:                 test.fields.cache,
				dnsCache:              test.fields.dnsCache,
				tlsConfig:             test.fields.tlsConfig,
				dnsRefreshDurationStr: test.fields.dnsRefreshDurationStr,
				dnsCacheExpirationStr: test.fields.dnsCacheExpirationStr,
				dnsRefreshDuration:    test.fields.dnsRefreshDuration,
				dnsCacheExpiration:    test.fields.dnsCacheExpiration,
				dialerTimeout:         test.fields.dialerTimeout,
				dialerKeepAlive:       test.fields.dialerKeepAlive,
				dialerDualStack:       test.fields.dialerDualStack,
				der:                   test.fields.der,
				dialer:                test.fields.dialer,
			}

			gotConn, err := d.cachedDialer(test.args.dctx, test.args.network, test.args.addr)
			if err := test.checkFunc(test.want, gotConn, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
