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

// Package net provides net functionality for vald's network connection
package net

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/cache/gache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/tls"
)

func Test_dialerCache_IP(t *testing.T) {
	t.Parallel()
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			d := &dialerCache{
				ips: test.fields.ips,
				cnt: test.fields.cnt,
			}

			got := d.IP()
			if err := checkFunc(d, test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialerCache_Len(t *testing.T) {
	t.Parallel()
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			d := &dialerCache{
				ips: test.fields.ips,
				cnt: test.fields.cnt,
			}

			got := d.Len()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewDialer(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []DialerOption
	}
	type want struct {
		dialer Dialer
		err    error
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

		if w.dialer == nil && gotDer == nil {
			return nil
		}
		if w.dialer == nil && gotDer != nil || w.dialer != nil && gotDer == nil {
			return errors.Errorf("got: \"%+v\",\n\t\t\t\twant: \"%+v\"", gotDer, w.dialer)
		}

		want, ok := w.dialer.(*dialer)
		if !ok {
			return errors.Errorf("want: \"%+v\" is not a dialer", w.dialer)
		}
		got, ok := gotDer.(*dialer)
		if !ok {
			return errors.Errorf("got: \"%+v\" is not a dialer", gotDer)
		}
		if diff := cmp.Diff(*want, *got,
			cmpopts.IgnoreFields(*want, "dialer", "der", "addrs", "dnsCachedOnce", "dnsCache", "ctrl", "tmu"),
			cmp.AllowUnexported(*want),
			cmp.Comparer(func(x, y cache.Cache) bool {
				if x == nil && y == nil {
					return true
				}
				return !(x == nil && y != nil) || !(y == nil && x != nil)
			}),
			cmp.Comparer(func(x, y *tls.Config) bool {
				return reflect.DeepEqual(x, y)
			}),
		); diff != "" {
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
		{
			name: "returns dialer when option is empty",
			want: want{
				dialer: func() (der Dialer) {
					der, _ = NewDialer()
					return der
				}(),
			},
		},
		func() test {
			opts := []DialerOption{
				WithTLS(func() *tls.Config {
					c, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
					if err != nil {
						return nil
					}
					return c
				}()),
				WithEnableDNSCache(),
				WithDNSRefreshDuration("5s"),
				WithDNSCacheExpiration("10s"),
			}
			return test{
				name: "returns dialer when option is not empty",
				args: args{
					opts: opts,
				},
				want: want{
					dialer: func() (der Dialer) {
						der, _ = NewDialer(opts...)
						return der
					}(),
				},
			}
		}(),
		func() test {
			opts := []DialerOption{
				WithTLS(func() *tls.Config {
					c, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
					if err != nil {
						return nil
					}
					return c
				}()),
				WithEnableDNSCache(),
				WithDNSRefreshDuration("50s"),
				WithDNSCacheExpiration("10s"),
			}
			return test{
				name: "return error when refresh duration > cache expiration and cache enabled",
				args: args{
					opts: opts,
				},
				want: want{
					err: errors.ErrInvalidDNSConfig(50*time.Second, 10*time.Second),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			gotDer, err := NewDialer(test.args.opts...)
			if err := checkFunc(test.want, gotDer, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialer_GetDialer(t *testing.T) {
	t.Parallel()
	type fields struct {
		dialer func(ctx context.Context, network, addr string) (Conn, error)
	}
	type want struct {
		want func(ctx context.Context, network, addr string) (Conn, error)
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, func(ctx context.Context, network, addr string) (Conn, error)) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got func(ctx context.Context, network, addr string) (Conn, error)) error {
		if reflect.ValueOf(w.want).Pointer() != reflect.ValueOf(got).Pointer() {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			fn := func(ctx context.Context, network, addr string) (Conn, error) {
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			d := &dialer{
				dialer: test.fields.dialer,
			}

			got := d.GetDialer()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialer_lookup(t *testing.T) {
	t.Parallel()
	type args struct {
		addr string
	}
	type want struct {
		want *dialerCache
		err  error
	}
	type test struct {
		name       string
		args       args
		opts       []DialerOption
		want       want
		checkFunc  func(context.Context, want, *dialerCache, error, *dialer) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(ctx context.Context, w want, got *dialerCache, err error, d *dialer) error {
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
				addr: "google.com",
			},
			opts: []DialerOption{},
			checkFunc: func(ctx context.Context, w want, got *dialerCache, err error, d *dialer) error {
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
				addr: "google.com",
			},
			opts: []DialerOption{
				WithDNSCache(gache.New()),
				WithDisableDialerDualStack(),
			},
			checkFunc: func(ctx context.Context, w want, got *dialerCache, err error, d *dialer) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got.Len() == 0 {
					return errors.New("ips is empty")
				}

				// check the cache is set
				if _, ok := d.dnsCache.Get("google.com"); !ok {
					return errors.New("cache is not set")
				}

				// execute lookup again and check the result is the same
				dc1, err1 := d.lookup(ctx, "google.com")
				if err1 != nil {
					return err1
				}
				if !reflect.DeepEqual(got, dc1) {
					return errors.Errorf("previous = %v, now %v", got, dc1)
				}

				// check the cache is set
				if _, ok := d.dnsCache.Get("google.com"); !ok {
					return errors.New("cache is not set")
				}

				return nil
			},
		},
		{
			name: "return cached ips when the cache hits",
			args: args{
				addr: "addr",
			},
			opts: []DialerOption{
				WithDNSCache(func() cache.Cache {
					g := gache.New()
					g.Set("addr", &dialerCache{
						ips: []string{"999.999.999.999"},
					})
					return g
				}()),
				WithDisableDialerDualStack(),
			},
			want: want{
				want: &dialerCache{
					ips: []string{"999.999.999.999"},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

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

			der, err := NewDialer(test.opts...)
			if err != nil {
				tt.Errorf("failed to initialize dialer: %v", err)
			}
			d, ok := der.(*dialer)
			if !ok {
				tt.Errorf("NewDialer return value Dialer is not *dialer: %v", der)
			}
			got, err := d.lookup(ctx, test.args.addr)
			if err := checkFunc(ctx, test.want, got, err, d); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialer_StartDialerCache(t *testing.T) {
	t.Parallel()
	type args struct{}
	type want struct{}
	type test struct {
		name       string
		args       args
		opts       []DialerOption
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
			addr := "localhost"
			ips := []string{}

			return test{
				name: "cache refresh when it is expired",

				opts: []DialerOption{
					WithEnableDNSCache(),
					WithDNSRefreshDuration("100ms"),
					WithDNSCacheExpiration("100ms"),
					WithDisableDialerDualStack(),
					WithDialerTimeout("1m"),
					WithDialerKeepalive("1m"),
				},
				beforeFunc: func(d *dialer) {
					d.dnsCache.Set(addr, &dialerCache{
						ips: ips,
					})
				},
				checkFunc: func(d *dialer) error {
					validateFn := func(ipLen int) error {
						val, ok := d.dnsCache.Get(addr)
						if !ok {
							return errors.New("cache not found")
						}
						if val == nil || len(val.(*dialerCache).ips) != ipLen {
							return errors.Errorf("cache is not correct, gotLen: %v, want: %v", val, ipLen)
						}
						return nil
					}

					// ensure the cache exists
					if err := validateFn(0); err != nil {
						return errors.Errorf("invalid cache err: %e", err)
					}

					// check cache update until timeout
					timeout := time.After(5 * time.Second)
					ticker := time.Tick(20 * time.Millisecond)
					for {
						select {
						case <-timeout:
							if err := validateFn(2); err != nil {
								return errors.Errorf("cache is not updated, err: %v", err)
							}
						case <-ticker:
							if err := validateFn(2); err == nil {
								return nil
							}
						}
					}
				},
			}
		}(),
		func() test {
			addr := "invalid"

			return test{
				name: "cache deleted when it is expired and the address is invalid or not available anymore",

				opts: []DialerOption{
					WithEnableDNSCache(),
					WithDNSRefreshDuration("100ms"),
					WithDNSCacheExpiration("100ms"),
					WithDisableDialerDualStack(),
					WithDialerTimeout("1m"),
					WithDialerKeepalive("1m"),
				},
				checkFunc: func(d *dialer) error {
					// ensure the cache exists
					_, ok := d.dnsCache.Get(addr)
					if ok {
						return errors.New("cache found")
					}
					// sleep and wait the cache removed
					time.Sleep(500 * time.Millisecond)

					// get again and check if the cache deleted
					if _, ok := d.dnsCache.Get(addr); ok {
						return errors.New("cache found")
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			der, err := NewDialer(test.opts...)
			if err != nil {
				tt.Errorf("failed to initialize dialer: %v", err)
			}
			d, ok := der.(*dialer)
			if !ok {
				tt.Errorf("NewDialer return value Dialer is not *dialer: %v", der)
			}
			if test.beforeFunc != nil {
				test.beforeFunc(d)
			}

			d.StartDialerCache(ctx)
			if err := checkFunc(d); err != nil {
				tt.Errorf("error = %v", err)
			}
			if test.afterFunc != nil {
				test.afterFunc(test.args)
			}
		})
	}
}

func Test_dialer_DialContext(t *testing.T) {
	t.Parallel()
	type args struct {
		network string
		address string
	}
	type fields struct {
		dialer func(ctx context.Context, network, addr string) (Conn, error)
	}
	type want struct {
		want Conn
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Conn, err error) error {
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
				network: "dummyNetwork",
				address: "dummyAddress",
			},
			fields: fields{
				dialer: func(ctx context.Context, network, addr string) (Conn, error) {
					if network == "dummyNetwork" && addr == "dummyAddress" {
						return nil, nil
					}
					return nil, errors.New("invalid error")
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

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
			d := &dialer{
				dialer: test.fields.dialer,
			}

			got, err := d.DialContext(ctx, test.args.network, test.args.address)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialer_cachedDialer(t *testing.T) {
	t.Parallel()
	type args struct {
		network string
		addr    string
	}
	type want struct {
		wantConn Conn
		err      error
	}
	type test struct {
		name       string
		args       args
		opts       []DialerOption
		want       want
		checkFunc  func(*dialer, context.Context, want, Conn, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotConn, w.wantConn) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotConn, w.wantConn)
		}
		return nil
	}
	tests := []test{
		func() test {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			return test{
				name: "return conn",
				args: args{
					network: TCP.String(),
					addr:    srv.URL[len("http://"):],
				},
				opts: nil,
				checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil: %v, want: %#v, got: %#v", err, w, gotConn)
					}

					if gotConn == nil {
						return errors.New("conn is nil")
					}
					return nil
				},
				afterFunc: func(*testing.T) {
					srv.Close()
				},
			}
		}(),
		func() test {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			srv.TLS.InsecureSkipVerify = true

			return test{
				name: "return tls conn",
				args: args{
					network: TCP.String(),
					addr:    srv.URL[len("https://"):],
				},
				opts: []DialerOption{
					WithTLS(func() *tls.Config {
						c, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
						if err != nil {
							return nil
						}
						return c
					}()),
					WithDialerTimeout("30s"),
				},
				checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil: %v", err)
					}

					if gotConn == nil {
						return errors.New("conn is nil")
					}
					return nil
				},
				afterFunc: func(t *testing.T) {
					srv.Close()
				},
			}
		}(), {
			name: "returns error when missing port in address",
			args: args{
				network: TCP.String(),
				addr:    "addr",
			},
			opts: nil,
			checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
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

			addr := "invalid_ip"

			c, err := cache.New()
			if err != nil {
				t.Error(err)
			}
			host, port, _ := SplitHostPort(strings.TrimPrefix(strings.TrimPrefix(srv.URL, "https://"), "http://"))

			return test{
				name: "return cached ip connection",
				args: args{
					network: TCP.String(),
					addr:    addr + ":" + strconv.FormatUint(uint64(port), 10),
				},
				opts: []DialerOption{
					WithDNSCache(c),
				},
				beforeFunc: func(t *testing.T) {
					// set the hostname 'invalid_ip' to the host name of the cache with the test server ip address
					c.Set("invalid_ip", &dialerCache{
						ips: []string{
							host,
						},
					})
				},
				checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil, %v", err)
					}
					if gotConn == nil {
						return errors.New("conn is nil")
					}

					if _, ok := c.Get(addr); !ok {
						return errors.New("cache value is deleted")
					}
					return nil
				},
				afterFunc: func(t *testing.T) {
					srv.Close()
				},
			}
		}(),
		func() test {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))

			addr := "invalid_ip"

			c, err := cache.New()
			if err != nil {
				t.Error(err)
			}
			host, port, _ := SplitHostPort(strings.TrimPrefix(strings.TrimPrefix(srv.URL, "https://"), "http://"))

			return test{
				name: "return cached ip tls connection",
				args: args{
					network: TCP.String(),
					addr:    addr + ":" + strconv.FormatUint(uint64(port), 10),
				},
				opts: []DialerOption{
					WithDNSCache(c),
					WithTLS(func() *tls.Config {
						c, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
						if err != nil {
							return nil
						}
						return c
					}()),
				},
				beforeFunc: func(t *testing.T) {
					// set the hostname 'invalid_ip' to the host name of the cache with the test server ip address
					c.Set("invalid_ip", &dialerCache{
						ips: []string{
							host,
						},
					})
				},
				checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
					if err != nil {
						return errors.Wrap(err, "err is not nil")
					}
					if gotConn == nil {
						return errors.New("conn is nil")
					}

					if _, ok := c.Get(addr); !ok {
						return errors.New("cache value is deleted")
					}
					return nil
				},
				afterFunc: func(t *testing.T) {
					srv.Close()
				},
			}
		}(),
		func() test {
			addr := "invalid_ip"

			c, err := cache.New()
			if err != nil {
				t.Error(err)
			}
			return test{
				name: "remove cache when dial failed",
				args: args{
					network: TCP.String(),
					addr:    addr + ":80",
				},
				opts: []DialerOption{
					WithDNSCache(c),
				},
				beforeFunc: func(*testing.T) {
					c.Set(addr, &dialerCache{
						ips: []string{
							addr,
						},
					})
				},
				checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
					if err == nil {
						return errors.New("err is nil")
					}
					if gotConn != nil {
						return errors.New("conn is nil")
					}

					if _, ok := c.Get(addr); ok {
						return errors.New("cache value is not deleted")
					}
					return nil
				},
			}
		}(),
		func() test {
			addr := "google.com"
			c, err := cache.New()
			if err != nil {
				t.Error(err)
			}
			return test{
				name: "retry when cache invalid and cache will be deleted",
				args: args{
					network: TCP.String(),
					addr:    addr + ":80",
				},
				opts: []DialerOption{
					WithDNSCache(c),
				},
				beforeFunc: func(*testing.T) {
					c.Set(addr, &dialerCache{
						ips: []string{
							"invalid_ip",
						},
					})
				},
				checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil: %v", err)
					}
					if gotConn == nil {
						return errors.New("conn is nil")
					}

					if c, ok := c.Get(addr); ok {
						return errors.Errorf("cache value is set: %+v", c)
					}
					return nil
				},
			}
		}(),
		func() test {
			srvNums := 20
			srvs := make([]*httptest.Server, srvNums)
			hosts := make([]string, srvNums)
			ports := make([]uint16, srvNums)
			addrs := make([]string, srvNums)

			for i := 0; i < srvNums; i++ {
				content := fmt.Sprint(i)
				srvs[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					fmt.Fprint(w, content)
				}))
				hosts[i], ports[i], _ = SplitHostPort(strings.TrimPrefix(strings.TrimPrefix(srvs[i].URL, "https://"), "http://"))
				addrs[i] = JoinHostPort(hosts[i], ports[i])
			}

			c, err := cache.New()
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "return cached ip connection in round robin order",
				args: args{
					network: TCP.String(),
					addr:    addrs[0],
				},
				opts: []DialerOption{
					WithDNSCache(c),
					WithEnableDNSCache(),
				},
				beforeFunc: func(t *testing.T) {
					c.Set(addrs[0], &dialerCache{
						ips: hosts,
					})
				},
				checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
					check := func(gotConn Conn, gotErr error, port uint16, srvContent string) error {
						defer func() {
							if gotConn != nil {
								_ = gotConn.Close()
							}
						}()

						if gotErr != nil {
							return errors.Errorf("err is not nil: %v", gotErr)
						}
						if gotConn == nil {
							return errors.New("conn is nil")
						}

						// check the connection made on the same port
						_, p, _ := net.SplitHostPort(gotConn.RemoteAddr().String())
						if p != strconv.Itoa(int(port)) {
							return errors.Errorf("unexcepted port number, except: %d, got: %s", port, p)
						}

						// read the output from the server and check if it is equals to the count
						fmt.Fprintf(gotConn, "GET / HTTP/1.0\r\n\r\n")
						buf, _ := io.ReadAll(gotConn)
						content := strings.Split(string(buf), "\n")[5] // skip HTTP header
						if content != srvContent {
							return errors.Errorf("excepted output from server, got: %v, want: %v", content, srvContent)
						}

						return nil
					}

					// check the return of the returned connection
					if err := check(gotConn, err, ports[0], "0"); err != nil {
						return errors.Errorf("check return connection, err: %v", err)
					}

					// check all remaining connection
					for i := 1; i < srvNums; i++ {
						c, e := d.cachedDialer(ctx, TCP.String(), addrs[i])
						srvContent := fmt.Sprint(i)
						if err := check(c, e, ports[i], srvContent); err != nil {
							return err
						}
					}

					// check all the connections again and it should start with index 0,
					// and the count should not be reset
					for i := 0; i < srvNums; i++ {
						c, e := d.cachedDialer(ctx, TCP.String(), addrs[i])
						srvContent := fmt.Sprint(i)
						if err := check(c, e, ports[i], srvContent); err != nil {
							return err
						}
					}

					return nil
				},
				afterFunc: func(t *testing.T) {
					for _, s := range srvs {
						s.Close()
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

			c, err := cache.New()
			if err != nil {
				t.Error(err)
			}
			c.Set(addr, &dialerCache{
				ips: []string{
					host, host, host,
				},
				cnt: math.MaxUint32,
			})

			return test{
				name: "reset cache count when it is overflow",
				args: args{
					network: TCP.String(),
					addr:    addr + ":" + fmt.Sprint(port),
				},
				opts: []DialerOption{
					WithDNSCache(c),
					WithDialerTimeout("10s"),
				},
				beforeFunc: func(t *testing.T) {
					c.Set(addr, &dialerCache{
						cnt: math.MaxUint32,
						ips: []string{host, host},
					})
				},
				checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil: %v", err)
					}
					if gotConn == nil {
						return errors.New("conn is nil")
					}

					c, _ := d.dnsCache.Get(addr)
					if dc := c.(*dialerCache); dc.cnt == math.MaxUint32 {
						return errors.Errorf("count do not reset, cnt: %v", dc.cnt)
					}

					return nil
				},
				afterFunc: func(t *testing.T) {
					srv.Close()
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt) // even test is failed, ensure afterFunc is executed to avoid goleak error
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			der, err := NewDialer(test.opts...)
			if err != nil {
				tt.Errorf("failed to initialize dialer: %v", err)
			}
			d, ok := der.(*dialer)
			if !ok {
				tt.Errorf("NewDialer return value Dialer is not *dialer: %v", der)
			}

			gotConn, gotErr := d.cachedDialer(ctx, test.args.network, test.args.addr)
			if err := checkFunc(d, ctx, test.want, gotConn, gotErr); err != nil {
				tt.Errorf("error = %v", err)
			}

			// call without defer to ensure the server is closed before checking with goleak
			if test.afterFunc != nil {
				test.afterFunc(tt)
			}
		})
	}
}

func Test_dialer_dial(t *testing.T) {
	t.Parallel()
	type args struct {
		network string
		addr    string
	}
	type fields struct {
		tlsConfig *tls.Config
		der       *net.Dialer
	}
	type want struct {
		want Conn
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Conn, err error) error {
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
				network: TCP.String(),
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
				network: TCP.String(),
				addr:    "google.com:443",
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
				network: TCP.String(),
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
			checkFunc: func(w want, got Conn, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got != nil {
					return errors.New("got is not nil")
				}
				return nil
			},
			want: want{
				err: &net.AddrError{Err: "missing port in address", Addr: "invalid_address"},
			},
		},
		{
			name: "return error if empty address",
			args: args{
				network: TCP.String(),
			},
			fields: fields{
				tlsConfig: nil,
				der: &net.Dialer{
					Timeout:   time.Minute,
					KeepAlive: time.Minute,
					DualStack: true,
				},
			},
			checkFunc: func(w want, got Conn, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got != nil {
					return errors.New("got is not nil")
				}
				return nil
			},
			want: want{
				err: errors.New("missing address"),
			},
		},
		{
			name: "return error if invalid network",
			args: args{
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
			checkFunc: func(w want, got Conn, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got != nil {
					return errors.New("got is not nil")
				}
				return nil
			},
			want: want{
				err: net.UnknownNetworkError("invalid"),
			},
		},
		{
			name: "return error if empty network",
			args: args{},
			fields: fields{
				tlsConfig: nil,
				der: &net.Dialer{
					Timeout:   time.Minute,
					KeepAlive: time.Minute,
					DualStack: true,
				},
			},
			checkFunc: func(w want, got Conn, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if got != nil {
					return errors.New("got is not nil")
				}
				return nil
			},
			want: want{
				err: net.UnknownNetworkError(""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

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
			d := &dialer{
				tlsConfig: test.fields.tlsConfig,
				der:       test.fields.der,
			}

			got, err := d.dial(ctx, test.args.network, test.args.addr)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_dialer_cacheExpireHook(t *testing.T) {
	t.Parallel()
	type args struct {
		addr string
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		opts       []DialerOption
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
					addr: addr,
				},
				opts: []DialerOption{
					WithEnableDNSCache(),
					WithDNSRefreshDuration("100ms"),
					WithDNSCacheExpiration("100ms"),
					WithDisableDialerDualStack(),
					WithDialerTimeout("1m"),
					WithDialerKeepalive("1m"),
				},
				checkFunc: func(d *dialer) error {
					// get again and check if the cache is updated
					if _, ok := d.dnsCache.Get(addr); !ok {
						return errors.New("cache not found")
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			der, err := NewDialer(test.opts...)
			if err != nil {
				tt.Errorf("failed to initialize dialer: %v", err)
			}
			d, ok := der.(*dialer)
			if !ok {
				tt.Errorf("NewDialer return value Dialer is not *dialer: %v", der)
			}
			if test.beforeFunc != nil {
				test.beforeFunc(d)
			}

			d.cacheExpireHook(ctx, test.args.addr)
			if err := checkFunc(d); err != nil {
				tt.Errorf("error = %v", err)
			}
			cancel()
		})
	}
}

func Test_dialer_tlsHandshake(t *testing.T) {
	t.Parallel()
	type args struct {
		conn    net.Conn
		network string
		addr    string
	}
	type want struct {
		want *tls.Conn
		err  error
	}
	type test struct {
		name       string
		args       args
		opts       []DialerOption
		want       want
		checkFunc  func(want, *tls.Conn, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *tls.Conn, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())

			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			srv.TLS.InsecureSkipVerify = true

			addr := srv.URL[len("https://"):]

			d, err := NewDialer()
			if err != nil {
				t.Error(err)
			}
			conn, err := d.DialContext(ctx, TCP.String(), addr)
			if err != nil || conn == nil {
				t.Errorf("failed to dial: %s, err: %v", addr, err)
			}

			return test{
				name: "return tls connection with handshake success with default timeout",
				args: args{
					network: TCP.String(),
					addr:    addr,
					conn:    conn,
				},
				opts: []DialerOption{
					WithDialerTimeout("30s"),
					WithTLS(func() *tls.Config {
						tlsCfg, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
						if err != nil {
							t.Fatal(err)
						}
						return tlsCfg
					}()),
				},
				checkFunc: func(w want, c *tls.Conn, err error) error {
					if c == nil || !c.ConnectionState().HandshakeComplete {
						return errors.Errorf("Handshake not completed, got: %+v\terr: %v", c, err)
					}
					return nil
				},
				afterFunc: func(t1 *testing.T) {
					srv.Close()
					conn.Close()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())

			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			srv.TLS.InsecureSkipVerify = true

			addr := srv.URL[len("https://"):]

			// create a new dialer to create the connection instead of using the original dialer
			c, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
			if err != nil {
				t.Error(err)
			}
			di, err := NewDialer(WithTLS(c))
			if err != nil {
				t.Errorf("failed to initialize dialer: %v", err)
			}

			conn, err := di.DialContext(ctx, TCP.String(), addr)
			if err != nil || conn == nil {
				t.Errorf("failed to dial: %s, err: %v", addr, err)
			}

			return test{
				name: "return error when handshake timeout",
				args: args{
					network: TCP.String(),
					addr:    addr,
					conn:    conn,
				},
				opts: []DialerOption{
					WithDialerTimeout("1us"),
					WithTLS(func() *tls.Config {
						c, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
						if err != nil {
							return nil
						}
						return c
					}()),
				},
				want: want{
					err: context.DeadlineExceeded,
				},
				afterFunc: func(t1 *testing.T) {
					srv.Close()
					conn.Close()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			srv := httptest.NewTLSServer(h)
			srv.TLS.InsecureSkipVerify = true

			d, err := NewDialer()
			if err != nil {
				t.Error(err)
			}

			host, port, _ := SplitHostPort(strings.TrimPrefix(strings.TrimPrefix(srv.URL, "https://"), "http://"))
			addr := host + ":" + strconv.FormatUint(uint64(port), 10)

			conn, err := d.DialContext(ctx, "tcp", addr)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "return error when host not found",
				args: args{
					network: TCP.String(),
					addr:    addr,
					conn:    conn,
				},
				opts: []DialerOption{
					WithTLS(func() *tls.Config {
						c, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
						if err != nil {
							return nil
						}
						return c
					}()),
				},
				beforeFunc: func(t1 *testing.T) {
					// close the server before the test
					srv.Close()
				},
				checkFunc: func(w want, c *tls.Conn, err error) error {
					if err == nil {
						return errors.New("Handshake completed even server has been gone")
					}
					return nil
				},
				afterFunc: func(t1 *testing.T) {
					conn.Close()
					cancel()
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			der, err := NewDialer(test.opts...)
			if err != nil {
				tt.Errorf("failed to initialize dialer: %v", err)
			}
			d, ok := der.(*dialer)
			if !ok || d == nil {
				tt.Errorf("NewDialer return value Dialer is not *dialer: %v", der)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}

			got, err := d.tlsHandshake(ctx, test.args.conn, test.args.network, test.args.addr)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_dialer_lookupIPAddrs(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		host string
	}
	type fields struct {
		dnsCache              cache.Cache
		enableDNSCache        bool
		dnsCachedOnce         sync.Once
		tlsConfig             *tls.Config
		tmu                   sync.RWMutex
		dnsRefreshDurationStr string
		dnsCacheExpirationStr string
		dnsRefreshDuration    time.Duration
		dnsCacheExpiration    time.Duration
		dialerTimeout         time.Duration
		dialerKeepalive       time.Duration
		dialerFallbackDelay   time.Duration
		ctrl                  control.SocketController
		sockFlg               control.SocketFlag
		dialerDualStack       bool
		addrs                 sync.Map
		der                   *net.Dialer
		dialer                func(ctx context.Context, network, addr string) (Conn, error)
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
		checkFunc  func(want, []string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotIps []string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotIps, w.wantIps) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIps, w.wantIps)
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
		           host: "",
		       },
		       fields: fields {
		           dnsCache: nil,
		           enableDNSCache: false,
		           dnsCachedOnce: sync.Once{},
		           tlsConfig: nil,
		           tmu: sync.RWMutex{},
		           dnsRefreshDurationStr: "",
		           dnsCacheExpirationStr: "",
		           dnsRefreshDuration: nil,
		           dnsCacheExpiration: nil,
		           dialerTimeout: nil,
		           dialerKeepalive: nil,
		           dialerFallbackDelay: nil,
		           ctrl: nil,
		           sockFlg: nil,
		           dialerDualStack: false,
		           addrs: sync.Map{},
		           der: net.Dialer{},
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
		           host: "",
		           },
		           fields: fields {
		           dnsCache: nil,
		           enableDNSCache: false,
		           dnsCachedOnce: sync.Once{},
		           tlsConfig: nil,
		           tmu: sync.RWMutex{},
		           dnsRefreshDurationStr: "",
		           dnsCacheExpirationStr: "",
		           dnsRefreshDuration: nil,
		           dnsCacheExpiration: nil,
		           dialerTimeout: nil,
		           dialerKeepalive: nil,
		           dialerFallbackDelay: nil,
		           ctrl: nil,
		           sockFlg: nil,
		           dialerDualStack: false,
		           addrs: sync.Map{},
		           der: net.Dialer{},
		           dialer: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			d := &dialer{
				dnsCache:              test.fields.dnsCache,
				enableDNSCache:        test.fields.enableDNSCache,
				dnsCachedOnce:         test.fields.dnsCachedOnce,
				tlsConfig:             test.fields.tlsConfig,
				tmu:                   test.fields.tmu,
				dnsRefreshDurationStr: test.fields.dnsRefreshDurationStr,
				dnsCacheExpirationStr: test.fields.dnsCacheExpirationStr,
				dnsRefreshDuration:    test.fields.dnsRefreshDuration,
				dnsCacheExpiration:    test.fields.dnsCacheExpiration,
				dialerTimeout:         test.fields.dialerTimeout,
				dialerKeepalive:       test.fields.dialerKeepalive,
				dialerFallbackDelay:   test.fields.dialerFallbackDelay,
				ctrl:                  test.fields.ctrl,
				sockFlg:               test.fields.sockFlg,
				dialerDualStack:       test.fields.dialerDualStack,
				addrs:                 test.fields.addrs,
				der:                   test.fields.der,
				dialer:                test.fields.dialer,
			}

			gotIps, err := d.lookupIPAddrs(test.args.ctx, test.args.host)
			if err := checkFunc(test.want, gotIps, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
