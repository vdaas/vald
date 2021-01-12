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
package tcp

import (
	"context"
	"crypto/tls"
	stderrors "errors"
	"fmt"
	"io/ioutil"
	"math"
	stdnet "net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/cache/gache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"go.uber.org/goleak"
)

var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
	goleak.IgnoreTopFunction("net._C2func_getaddrinfo"),
}

func Test_dialerCache_IP(t *testing.T) {
	type fields struct {
		ips []string
		cnt uint32
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(*dialerCache, want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(d *dialerCache, w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return IP directly if size is 1",
			fields: fields{
				ips: []string{
					"a",
				},
			},
			want: want{
				want: "a",
			},
			checkFunc: func(d *dialerCache, w want, got string) error {
				if err := defaultCheckFunc(d, w, got); err != nil {
					return err
				}

				for i := 1; i < 100; i++ {
					if d.IP() != "a" {
						return errors.New("invalid output")
					}
					if d.cnt != 0 {
						return errors.New("invalid cnt")
					}
				}
				return nil
			},
		},
		{
			name: "return IP in round robin order",
			fields: fields{
				ips: []string{
					"a", "b", "c",
				},
			},
			want: want{
				want: "b",
			},
			checkFunc: func(d *dialerCache, w want, got string) error {
				if err := defaultCheckFunc(d, w, got); err != nil {
					return err
				}

				for i := 1; i < 100; i++ {
					idx := (i + 1) % len(d.ips)
					if s := d.IP(); s != d.ips[idx] {
						return errors.New("invalid output")
					}
					if d.cnt != uint32(i+1) {
						return errors.New("invalid cnt")
					}
				}
				return nil
			},
		},
		{
			name: "cnt reset when it is about to overflow",
			fields: fields{
				ips: []string{
					"a", "b", "c",
				},
				cnt: math.MaxUint32,
			},
			want: want{
				want: "a",
			},
			checkFunc: func(d *dialerCache, w want, got string) error {
				if err := defaultCheckFunc(d, w, got); err != nil {
					return err
				}
				if d.cnt != 0 {
					return errors.New("invalid cnt")
				}
				return nil
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
			d := &dialerCache{
				ips: test.fields.ips,
				cnt: test.fields.cnt,
			}

			got := d.IP()
			if err := test.checkFunc(d, test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialerCache_Len(t *testing.T) {
	type fields struct {
		ips []string
		cnt uint32
	}
	type want struct {
		want uint32
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint32) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint32) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ips length",
			fields: fields{
				ips: []string{"a"},
			},
			want: want{want: uint32(1)},
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
			d := &dialerCache{
				ips: test.fields.ips,
				cnt: test.fields.cnt,
			}

			got := d.Len()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

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
			return errors.Errorf("got_error: \"%+v\",\n\t\t\t\twant: \"%+v\"", err, w.err)
		}

		if w.wantDer == nil && gotDer == nil {
			return nil
		}
		if w.wantDer == nil && gotDer != nil || w.wantDer != nil && gotDer == nil {
			return errors.Errorf("got: \"%+v\",\n\t\t\t\twant: \"%+v\"", gotDer, w.wantDer)
		}

		want := w.wantDer.(*dialer)
		got := gotDer.(*dialer)
		opts := []cmp.Option{
			cmp.AllowUnexported(*want),
			cmpopts.IgnoreFields(*want, "dialer", "der", "addrs"),
			cmp.Comparer(func(x, y cache.Cache) bool {
				if x == nil && y == nil {
					return true
				}
				return !(x == nil && y != nil) || !(y == nil && x != nil)
			}),
			cmp.Comparer(func(x, y *tls.Config) bool {
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

			return test{
				name: "return error when refresh duration > cache expiration and cache enabled",
				args: args{
					opts: []DialerOption{
						WithTLS(tc),
						WithEnableDNSCache(),
						WithDNSRefreshDuration("50s"),
						WithDNSCacheExpiration("10s"),
					},
				},
				want: want{
					err: errors.ErrInvalidDNSConfig(50*time.Second, 10*time.Second),
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
		want *dialerCache
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *dialerCache, error, *dialer) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *dialerCache, err error, d *dialer) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
			checkFunc: func(w want, got *dialerCache, err error, d *dialer) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got.Len() == 0 {
					return errors.New("ips is empty")
				}
				return nil
			},
		},
		{
			name: "return ips when lookupIpAddr returns and the cache is set",
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
			checkFunc: func(w want, got *dialerCache, err error, d *dialer) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got.Len() == 0 {
					return errors.New("ips is empty")
				}

				// check the cache is set
				if _, ok := d.cache.Get("google.com"); !ok {
					return errors.New("cache is not set")
				}

				// execute lookup again and check the result is the same
				dc1, err1 := d.lookup(context.Background(), "google.com")
				if err1 != nil {
					return err1
				}
				if !reflect.DeepEqual(got, dc1) {
					return errors.Errorf("got = %v, got1 %v", got, dc1)
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
					g.Set("addr", &dialerCache{
						ips: []string{"999.999.999.999"},
					})
					return g
				}(),
				der: &net.Dialer{
					Resolver: &net.Resolver{
						PreferGo: false,
					},
				},
			},
			want: want{
				want: &dialerCache{
					ips: []string{"999.999.999.999"},
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
			got, err := d.lookup(test.args.ctx, test.args.addr)
			if err := test.checkFunc(test.want, got, err, d); err != nil {
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
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(*dialer) error
		beforeFunc func(*dialer)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(d *dialer) error {
		return nil
	}
	tests := []test{
		func() test {
			addr := "google.com"
			ips := []string{"0.0.0.0"}
			ctx, cancel := context.WithCancel(context.Background())

			return test{
				name: "cache refresh when it is expired",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					dnsCache:           true,
					dnsRefreshDuration: time.Millisecond * 100,
					dnsCacheExpiration: time.Millisecond * 100,
				},
				beforeFunc: func(d *dialer) {
					d.cache, _ = cache.New(cache.WithExpireDuration("300ms"), cache.WithExpireCheckDuration("100ms"),
						cache.WithExpiredHook(d.cacheExpireHook))
					d.cache.Set(addr, &dialerCache{ips: ips})

					d.der = &net.Dialer{
						Timeout:   time.Minute,
						KeepAlive: time.Minute,
						DualStack: d.dialerDualStack,
						Control:   Control,
						Resolver: &net.Resolver{
							PreferGo: false,
							Dial:     d.dialer,
						},
					}
				},
				checkFunc: func(d *dialer) error {
					// ensure the cache exists
					val, ok := d.cache.Get(addr)
					if !ok {
						return errors.New("cache not found")
					}
					if !reflect.DeepEqual(val.(*dialerCache).ips, ips) {
						return errors.New("cache is not correct")
					}

					// sleep and wait the cache update
					time.Sleep(500 * time.Millisecond)

					// get again and check if the cache is updated
					val, ok = d.cache.Get(addr)
					if !ok {
						return errors.New("cache not found")
					}
					if reflect.DeepEqual(val.(*dialerCache).ips, ips) {
						return errors.New("cache is not updated")
					}
					return nil
				},
				afterFunc: func(args) {
					cancel()
				},
			}
		}(),
		func() test {
			addr := "invalid"
			ips := []string{"0.0.0.0"}
			ctx, cancel := context.WithCancel(context.Background())

			return test{
				name: "cache deleted when it is expired and the address is invalid or not available anymore",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					dnsCache:           true,
					dnsRefreshDuration: time.Millisecond * 100,
					dnsCacheExpiration: time.Millisecond * 100,
				},
				beforeFunc: func(d *dialer) {
					d.cache, _ = cache.New(cache.WithExpireDuration("300ms"), cache.WithExpireCheckDuration("100ms"),
						cache.WithExpiredHook(d.cacheExpireHook))
					d.cache.Set(addr, &dialerCache{ips: ips})

					d.der = &net.Dialer{
						Timeout:   time.Minute,
						KeepAlive: time.Minute,
						DualStack: d.dialerDualStack,
						Control:   Control,
						Resolver: &net.Resolver{
							PreferGo: false,
							Dial:     d.dialer,
						},
					}
				},
				checkFunc: func(d *dialer) error {
					// ensure the cache exists
					val, ok := d.cache.Get(addr)
					if !ok {
						return errors.New("cache not found")
					}
					if !reflect.DeepEqual(val.(*dialerCache).ips, ips) {
						return errors.New("cache is not correct")
					}

					// sleep and wait the cache removed
					time.Sleep(500 * time.Millisecond)

					// get again and check if the cache deleted
					if _, ok := d.cache.Get(addr); ok {
						return errors.New("cache found")
					}
					return nil
				},
				afterFunc: func(args) {
					cancel()
					time.Sleep(500 * time.Millisecond)
				},
			}
		}(),
	}

	log.Init()
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
			if test.beforeFunc != nil {
				test.beforeFunc(d)
			}

			d.StartDialerCache(test.args.ctx)
			if err := test.checkFunc(d); err != nil {
				tt.Errorf("error = %v", err)
			}
			if test.afterFunc != nil {
				test.afterFunc(test.args)
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
		dialer func(ctx context.Context, network, addr string) (net.Conn, error)
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return non nil error",
			args: args{
				ctx:     context.Background(),
				network: "dummyNetwork",
				address: "dummyAddress",
			},
			fields: fields{
				dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
					if network == "dummyNetwork" && addr == "dummyAddress" {
						return nil, nil
					}
					return nil, errors.New("invalid error")
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
			d := &dialer{
				dialer: test.fields.dialer,
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
		checkFunc  func(*dialer, want, net.Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(d *dialer, w want, gotConn net.Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotConn, w.wantConn) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotConn, w.wantConn)
		}
		return nil
	}
	tests := []test{
		{
			name: "return conn",
			args: args{
				dctx:    context.Background(),
				network: "tcp",
				addr:    "google.com:80",
			},
			fields: fields{
				der: &net.Dialer{
					Resolver: &net.Resolver{
						PreferGo: false,
					},
				},
			},
			checkFunc: func(d *dialer, w want, gotConn net.Conn, err error) error {
				if err != nil {
					return errors.Errorf("err is not nil: %v", err)
				}

				if gotConn == nil {
					return errors.New("conn is nil")
				}
				return nil
			},
		},
		{
			name: "return tls conn",
			args: args{
				dctx:    context.Background(),
				network: "tcp",
				addr:    "google.com:80",
			},
			fields: fields{
				der: &net.Dialer{
					Resolver: &net.Resolver{
						PreferGo: false,
					},
				},
				tlsConfig: new(tls.Config),
			},
			checkFunc: func(d *dialer, w want, gotConn net.Conn, err error) error {
				if err != nil {
					return errors.Errorf("err is not nil: %v", err)
				}

				if gotConn == nil {
					return errors.New("conn is nil")
				}
				return nil
			},
		},
		{
			name: "returns error when missing port in address",
			args: args{
				dctx:    context.Background(),
				network: "tcp",
				addr:    "addr",
			},
			fields: fields{
				der: &net.Dialer{
					Resolver: &net.Resolver{
						PreferGo: false,
					},
				},
			},
			checkFunc: func(d *dialer, w want, gotConn net.Conn, err error) error {
				if err == nil {
					return errors.New("err is nil")
				}

				if gotConn != nil {
					return errors.Errorf("conn is not nil: %v", gotConn)
				}

				return nil
			},
		},
		func() test {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			host, port, _ := net.SplitHostPort(srv.URL[len("http://"):])

			addr := "invalid_ip"

			// set the hostname 'invalid_ip' to the host name of the cache with the test server ip address
			cache, _ := cache.New()
			cache.Set(addr, &dialerCache{
				ips: []string{
					host,
				},
			})

			return test{
				name: "return cached ip connection",
				args: args{
					dctx:    context.Background(),
					network: "tcp",
					addr:    addr + ":" + strconv.FormatUint(uint64(port), 10),
				},
				fields: fields{
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						},
					},
					cache:    cache,
					dnsCache: true,
				},
				checkFunc: func(d *dialer, w want, gotConn net.Conn, err error) error {
					if err != nil {
						return errors.New("err is not nil")
					}
					if gotConn == nil {
						return errors.New("conn is nil")
					}

					if _, ok := cache.Get(addr); !ok {
						return errors.New("cache value is deleted")
					}
					return nil
				},
				afterFunc: func(args) {
					srv.Close()
				},
			}
		}(),
		func() test {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))

			host, port, _ := net.SplitHostPort(srv.URL[len("https://"):])

			addr := "invalid_ip"

			// set the hostname 'invalid_ip' to the host name of the cache with the test server ip address
			cache, _ := cache.New()
			cache.Set(addr, &dialerCache{
				ips: []string{
					host,
				},
			})

			return test{
				name: "return cached ip tls connection",
				args: args{
					dctx:    context.Background(),
					network: "tcp",
					addr:    addr + ":" + strconv.FormatUint(uint64(port), 10),
				},
				fields: fields{
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						},
					},
					cache:     cache,
					dnsCache:  true,
					tlsConfig: new(tls.Config),
				},
				checkFunc: func(d *dialer, w want, gotConn net.Conn, err error) error {
					if err != nil {
						return errors.New("err is not nil")
					}
					if gotConn == nil {
						return errors.New("conn is nil")
					}

					if _, ok := cache.Get(addr); !ok {
						return errors.New("cache value is deleted")
					}
					return nil
				},
				afterFunc: func(args) {
					srv.Close()
				},
			}
		}(),
		func() test {
			addr := "invalid_ip"

			cache, _ := cache.New()
			cache.Set(addr, &dialerCache{
				ips: []string{
					addr,
				},
			})

			return test{
				name: "remove cache when dial failed",
				args: args{
					dctx:    context.Background(),
					network: "tcp",
					addr:    addr + ":80",
				},
				fields: fields{
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						},
					},
					cache:    cache,
					dnsCache: true,
				},
				checkFunc: func(d *dialer, w want, gotConn net.Conn, err error) error {
					if err == nil {
						return errors.New("err is nil")
					}
					if gotConn != nil {
						return errors.New("conn is nil")
					}

					if _, ok := cache.Get(addr); ok {
						return errors.New("cache value is not deleted")
					}
					return nil
				},
			}
		}(),
		func() test {
			addr := "google.com"

			cache, _ := cache.New()
			cache.Set(addr, &dialerCache{
				ips: []string{
					"invalid_ip",
				},
			})

			return test{
				name: "retry when cache invalid and cache will be deleted",
				args: args{
					dctx:    context.Background(),
					network: "tcp",
					addr:    addr + ":80",
				},
				fields: fields{
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						},
					},
					cache:    cache,
					dnsCache: true,
				},
				checkFunc: func(d *dialer, w want, gotConn net.Conn, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil: %v", err)
					}
					if gotConn == nil {
						return errors.New("conn is nil")
					}

					if c, ok := cache.Get(addr); ok {
						return errors.Errorf("cache value is set: %+v", c)
					}
					return nil
				},
			}
		}(),
		func() test {
			srvNums := 20
			srvs := make([]*httptest.Server, 0, srvNums)
			hosts := make([]string, 0, srvNums)
			ports := make([]string, 0, srvNums)

			// create servers that will return the server number
			for i := 0; i < srvNums; i++ {
				content := fmt.Sprint(i)
				hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, content)
					w.WriteHeader(200)
				})
				srvs = append(srvs, httptest.NewServer(hf))
				h, p, _ := net.SplitHostPort(srvs[i].URL[len("http://"):])
				hosts = append(hosts, h)
				ports = append(ports, fmt.Sprint(p))
			}

			addr := "address"

			cache, _ := cache.New()
			cache.Set(addr, &dialerCache{
				ips: hosts,
			})

			return test{
				name: "return cached ip connection in round robin order",
				args: args{
					dctx:    context.Background(),
					network: "tcp",
					addr:    addr + ":" + ports[0],
				},
				fields: fields{
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						},
					},
					cache:    cache,
					dnsCache: true,
				},
				checkFunc: func(d *dialer, w want, gotConn net.Conn, err error) error {
					c, _ := d.cache.Get(addr)
					dc := c.(*dialerCache)

					check := func(gotConn net.Conn, gotErr error, cnt int, port string, srvContent string) error {
						defer func() {
							_ = gotConn.Close()
						}()

						if gotErr != nil {
							return errors.Errorf("err is not nil: %v", gotErr)
						}
						if gotConn == nil {
							return errors.New("conn is nil")
						}
						if c := atomic.LoadUint32(&dc.cnt); c != uint32(cnt) {
							return errors.Errorf("cnt not correct, %d, except: %d", c, cnt)
						}

						// check the connection made is the same excepted
						_, p, _ := net.SplitHostPort(gotConn.RemoteAddr().String())
						if fmt.Sprint(p) != port {
							return errors.Errorf("unexcepted port number, except: %v, got: %v", port, p)
						}

						// read the output from the server and check if it is equals to the count
						fmt.Fprintf(gotConn, "GET / HTTP/1.0\r\n\r\n")
						buf, _ := ioutil.ReadAll(gotConn)
						content := strings.Split(string(buf), "\n")[5] // skip HTTP header
						if content != srvContent {
							return errors.Errorf("excepted output from server, got: %v, want: %v", content, fmt.Sprint(cnt))
						}

						return nil
					}

					// check the return of the returned connection
					if err := check(gotConn, err, 1, ports[0], "0"); err != nil {
						return err
					}

					// check all the connection
					for i := 1; i < srvNums; i++ {
						c, e := d.cachedDialer(context.Background(), "tcp", addr+":"+ports[i])
						srvContent := fmt.Sprint(i)
						if err := check(c, e, i+1, ports[i], srvContent); err != nil {
							return err
						}
					}

					// check all the connections again and it should start with index 0,
					// and the count should not be reset
					for i := 0; i < srvNums; i++ {
						c, e := d.cachedDialer(context.Background(), "tcp", addr+":"+ports[i])
						cnt := srvNums + i + 1
						srvContent := fmt.Sprint(i)
						if err := check(c, e, cnt, ports[i], srvContent); err != nil {
							return err
						}
					}

					return nil
				},
				afterFunc: func(args) {
					for _, srv := range srvs {
						srv.Close()
					}
				},
			}
		}(),
		func() test {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			host, port, _ := net.SplitHostPort(srv.URL[len("http://"):])

			addr := "invalid_ip"

			cache, _ := cache.New()
			cache.Set(addr, &dialerCache{
				ips: []string{
					host, host,
				},
				cnt: math.MaxUint32,
			})

			return test{
				name: "reset cache count when it is  overflow",
				args: args{
					dctx:    context.Background(),
					network: "tcp",
					addr:    addr + ":" + fmt.Sprint(port),
				},
				fields: fields{
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						},
					},
					cache:    cache,
					dnsCache: true,
				},
				checkFunc: func(d *dialer, w want, gotConn net.Conn, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil: %v", err)
					}
					if gotConn == nil {
						return errors.New("conn is nil")
					}

					c, _ := d.cache.Get(addr)
					if dc := c.(*dialerCache); dc.cnt != 0 {
						return errors.Errorf("count do not reset, cnt: %v", dc.cnt)
					}

					return nil
				},
				afterFunc: func(args) {
					srv.Close()
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

			gotConn, gotErr := d.cachedDialer(test.args.dctx, test.args.network, test.args.addr)
			if err := test.checkFunc(d, test.want, gotConn, gotErr); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialer_dial(t *testing.T) {
	type args struct {
		ctx     context.Context
		network string
		addr    string
	}
	type fields struct {
		tlsConfig *tls.Config
		der       *net.Dialer
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if got == nil {
			return errors.New("got is nil")
		}
		return nil
	}
	tests := []test{
		{
			name: "return conn",
			args: args{
				ctx:     context.Background(),
				network: "tcp",
				addr:    "google.com:80",
			},
			fields: fields{
				tlsConfig: nil,
				der: &net.Dialer{
					Timeout:   time.Minute,
					KeepAlive: time.Minute,
					DualStack: true,
				},
			},
		},
		{
			name: "return tls conn",
			args: args{
				ctx:     context.Background(),
				network: "tcp",
				addr:    "google.com:80",
			},
			fields: fields{
				tlsConfig: new(tls.Config),
				der: &net.Dialer{
					Timeout:   time.Minute,
					KeepAlive: time.Minute,
					DualStack: true,
				},
			},
		},
		{
			name: "return error if invalid address",
			args: args{
				ctx:     context.Background(),
				network: "tcp",
				addr:    "invalid_address",
			},
			fields: fields{
				tlsConfig: nil,
				der: &net.Dialer{
					Timeout:   time.Minute,
					KeepAlive: time.Minute,
					DualStack: true,
				},
			},
			checkFunc: func(w want, got net.Conn, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got != nil {
					return errors.New("got is not nil")
				}
				return nil
			},
			want: want{
				err: &stdnet.AddrError{Err: "missing port in address", Addr: "invalid_address"},
			},
		},
		{
			name: "return error if empty address",
			args: args{
				ctx:     context.Background(),
				network: "tcp",
			},
			fields: fields{
				tlsConfig: nil,
				der: &net.Dialer{
					Timeout:   time.Minute,
					KeepAlive: time.Minute,
					DualStack: true,
				},
			},
			checkFunc: func(w want, got net.Conn, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got != nil {
					return errors.New("got is not nil")
				}
				return nil
			},
			want: want{
				err: stderrors.New("missing address"),
			},
		},
		{
			name: "return error if invalid network",
			args: args{
				ctx:     context.Background(),
				network: "invalid",
			},
			fields: fields{
				tlsConfig: nil,
				der: &net.Dialer{
					Timeout:   time.Minute,
					KeepAlive: time.Minute,
					DualStack: true,
				},
			},
			checkFunc: func(w want, got net.Conn, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got != nil {
					return errors.New("got is not nil")
				}
				return nil
			},
			want: want{
				err: stdnet.UnknownNetworkError("invalid"),
			},
		},
		{
			name: "return error if empty network",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				tlsConfig: nil,
				der: &net.Dialer{
					Timeout:   time.Minute,
					KeepAlive: time.Minute,
					DualStack: true,
				},
			},
			checkFunc: func(w want, got net.Conn, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got != nil {
					return errors.New("got is not nil")
				}
				return nil
			},
			want: want{
				err: stdnet.UnknownNetworkError(""),
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
				tlsConfig: test.fields.tlsConfig,
				der:       test.fields.der,
			}

			got, err := d.dial(test.args.ctx, test.args.network, test.args.addr)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialer_cacheExpireHook(t *testing.T) {
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
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(*dialer) error
		beforeFunc func(*dialer)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(*dialer) error {
		return nil
	}
	tests := []test{
		func() test {
			addr := "google.com"
			return test{
				name: "cache refresh",
				args: args{
					ctx:  context.Background(),
					addr: addr,
				},
				fields: fields{
					dnsCache:           true,
					dnsRefreshDuration: time.Millisecond * 100,
					dnsCacheExpiration: time.Millisecond * 100,
				},
				beforeFunc: func(d *dialer) {
					d.cache, _ = cache.New()

					d.der = &net.Dialer{
						Timeout:   time.Minute,
						KeepAlive: time.Minute,
						DualStack: d.dialerDualStack,
						Control:   Control,
						Resolver: &net.Resolver{
							PreferGo: false,
							Dial:     d.dialer,
						},
					}
				},
				checkFunc: func(d *dialer) error {
					// get again and check if the cache is updated
					if _, ok := d.cache.Get(addr); !ok {
						return errors.New("cache not found")
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			if test.beforeFunc != nil {
				test.beforeFunc(d)
			}

			d.cacheExpireHook(test.args.ctx, test.args.addr)
			if err := test.checkFunc(d); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
