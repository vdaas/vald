// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	"context"
	"crypto/x509"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/cache/cacher"
	"github.com/vdaas/vald/internal/cache/gache"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/test/comparator"
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
		checkFunc  func(*dialerCache, want, string) error
		beforeFunc func()
		afterFunc  func()
		name       string
		want       want
		fields     fields
	}
	defaultCheckFunc := func(_ *dialerCache, w want, got string) error {
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
		checkFunc  func(want, uint32) error
		beforeFunc func()
		afterFunc  func()
		name       string
		fields     fields
		want       want
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
		want       want
		checkFunc  func(want, Dialer, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
		// skipcq: VET-V0008
		//nolint:govet,copylocks
		if diff := comparator.Diff(*want, *got,
			// skipcq: VET-V0008
			//nolint:govet,copylocks
			comparator.IgnoreFields(*want, "dialer", "der", "addrs", "dnsCachedOnce", "dnsCache", "ctrl"),
			// skipcq: VET-V0008
			//nolint:govet,copylocks
			comparator.AllowUnexported(*want),
			comparator.Comparer(func(x, y cacher.Cache[*dialerCache]) bool {
				if x == nil && y == nil {
					return true
				}
				return reflect.DeepEqual(x, y)
			}),
			comparator.Comparer(func(x, y *tls.Config) bool {
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
		fields     fields
		want       want
		checkFunc  func(want, func(ctx context.Context, network, addr string) (Conn, error)) error
		beforeFunc func()
		afterFunc  func()
		name       string
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
		want       want
		checkFunc  func(context.Context, want, *dialerCache, error, *dialer) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
		opts       []DialerOption
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
				WithDNSCache(gache.New[*dialerCache]()),
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
				WithDNSCache(func() cacher.Cache[*dialerCache] {
					g := gache.New[*dialerCache]()
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

			ctx := tt.Context()

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
		args       args
		want       want
		checkFunc  func(*dialer) error
		beforeFunc func(*testing.T, *dialer)
		afterFunc  func(args)
		name       string
		opts       []DialerOption
	}
	defaultCheckFunc := func(d *dialer) error {
		return nil
	}
	tests := []test{
		func() test {
			addr := "localhost"
			var ips []string

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
				beforeFunc: func(t *testing.T, d *dialer) {
					t.Helper()
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
						if val == nil || len(val.ips) != ipLen {
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

			ctx := tt.Context()

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
				test.beforeFunc(tt, d)
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
		want       want
		fields     fields
		checkFunc  func(want, Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		args       args
		name       string
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

			ctx := tt.Context()

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
		want       want
		checkFunc  func(*dialer, context.Context, want, Conn, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		args       args
		name       string
		opts       []DialerOption
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
				w.WriteHeader(http.StatusOK)
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
				afterFunc: func(t *testing.T) {
					t.Helper()
					srv.Close()
				},
			}
		}(),
		func() test {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
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
					t.Helper()
					srv.Close()
				},
			}
		}(),
		{
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
				w.WriteHeader(http.StatusOK)
			}))

			c, err := cache.New[*dialerCache]()
			if err != nil {
				t.Error(err)
			}

			addr := srv.URL[len("http://"):]
			host, _, err := SplitHostPort(addr)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "return cached ip connection",
				args: args{
					network: TCP.String(),
					addr:    addr,
				},
				opts: []DialerOption{
					WithDNSCache(c),
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					// set the hostname 'invalid_ip' to the host name of the cache with the test server ip address
					c.Set(addr, &dialerCache{
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
					t.Helper()
					srv.Close()
				},
			}
		}(),
		func() test {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			srv.TLS.InsecureSkipVerify = true

			c, err := cache.New[*dialerCache]()
			if err != nil {
				t.Error(err)
			}
			tls, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
			if err != nil {
				t.Error(err)
			}

			addr := srv.URL[len("https://"):]
			host, _, err := SplitHostPort(addr)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "return cached ip tls connection",
				args: args{
					network: TCP.String(),
					addr:    addr,
				},
				opts: []DialerOption{
					WithDNSCache(c),
					WithTLS(tls),
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					// set the hostname 'invalid_ip' to the host name of the cache with the test server ip address
					c.Set(addr, &dialerCache{
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
					t.Helper()
					srv.Close()
				},
			}
		}(),
		func() test {
			addr := "invalid_ip"

			c, err := cache.New[*dialerCache]()
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
					WithEnableDNSCache(),
					WithDNSCache(c),
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
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
			c, err := cache.New[*dialerCache]()
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
				beforeFunc: func(t *testing.T) {
					t.Helper()
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

			for i := range srvNums {
				content := fmt.Sprint(i)
				srvs[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, content)
				}))
				hosts[i], ports[i], _ = SplitHostPort(strings.TrimPrefix(strings.TrimPrefix(srvs[i].URL, "https://"), "http://"))
				addrs[i] = JoinHostPort(hosts[i], ports[i])
			}

			c, err := cache.New[*dialerCache]()
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
					t.Helper()
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
							return errors.Errorf("unexpected port number, except: %d, got: %s", port, p)
						}

						// read the output from the server and check if it is equals to the count
						fmt.Fprintf(gotConn, "GET / HTTP/1.0\r\n\r\n")
						buf, err := io.ReadAll(gotConn)
						if err != nil || buf == nil {
							return errors.Errorf("error or buffer is nil,\terror: %v, buf: %v", err, buf)
						}
						content := strings.Split(conv.Btoa(buf), "\n")[5] // skip HTTP header
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
					for i := range srvNums {
						c, e := d.cachedDialer(ctx, TCP.String(), addrs[i])
						srvContent := fmt.Sprint(i)
						if err := check(c, e, ports[i], srvContent); err != nil {
							return err
						}
					}

					return nil
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for _, s := range srvs {
						s.Close()
					}
				},
			}
		}(),
		func() test {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			// Use a host name: an IP literal skips the DNS cache entirely, so the
			// counter would never be exercised.
			ip, port, err := SplitHostPort(srv.URL[len("http://"):])
			if err != nil {
				t.Error(err)
			}
			const host = "localhost"
			addr := JoinHostPort(host, port)

			c, err := cache.New[*dialerCache]()
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "reset cache count when it is overflow",
				args: args{
					network: TCP.String(),
					addr:    addr,
				},
				opts: []DialerOption{
					WithDNSCache(c),
					WithDialerTimeout("10s"),
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					c.Set(host, &dialerCache{
						cnt: math.MaxUint32,
						ips: []string{ip, ip},
					})
				},
				checkFunc: func(d *dialer, ctx context.Context, w want, gotConn Conn, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil: %v", err)
					}
					if gotConn == nil {
						return errors.New("conn is nil")
					}

					dc, ok := d.dnsCache.Get(host)
					if !ok {
						return errors.Errorf("dns cache entry for %s not found", host)
					}
					if dc.cnt != 0 {
						return errors.Errorf("count do not reset, cnt: %v", dc.cnt)
					}

					return nil
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					srv.Close()
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx := tt.Context()

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
				tt.Errorf("error = %v, got(error: %v, conn: %v)", err, gotErr, gotConn)
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
		want       want
		fields     fields
		checkFunc  func(want, Conn, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		args       args
		name       string
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
				err: errors.ErrInvalidAddress(TCP.String(), ""),
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
				err: errors.ErrInvalidAddress(TCP.String(), ""),
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
				err: errors.ErrInvalidAddress(TCP.String(), ""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx := tt.Context()

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
		cache *dialerCache
		addr  string
	}
	type want struct{}
	type test struct {
		want       want
		args       args
		checkFunc  func(*dialer) error
		beforeFunc func(*dialer)
		afterFunc  func(args)
		name       string
		opts       []DialerOption
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
		func() test {
			addr := "10.0.0.1"
			return test{
				name: "IP literal is not refreshed",
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
					if _, ok := d.dnsCache.Get(addr); ok {
						return errors.New("expired IP literal was cached again")
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

			ctx, cancel := context.WithCancel(tt.Context())
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

			d.cacheExpireHook(ctx, test.args.addr, test.args.cache)
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
		want       want
		checkFunc  func(want, *tls.Conn, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		args       args
		name       string
		opts       []DialerOption
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
			ctx, cancel := context.WithCancel(t.Context())

			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
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
				afterFunc: func(t *testing.T) {
					t.Helper()
					srv.Close()
					conn.Close()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())

			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
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
				afterFunc: func(t *testing.T) {
					t.Helper()
					srv.Close()
					conn.Close()
					cancel()
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
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
				beforeFunc: func(t *testing.T) {
					t.Helper()
					// close the server before the test
					srv.Close()
				},
				checkFunc: func(w want, c *tls.Conn, err error) error {
					if err == nil {
						return errors.New("Handshake completed even server has been gone")
					}
					return nil
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
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

			ctx := tt.Context()

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

			got, err := d.tlsHandshake(ctx, test.args.conn, test.args.network, test.args.addr, "")
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_dialer_lookupIPAddrs(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		host string
// 	}
// 	type fields struct {
// 		dnsCache              cacher.Cache[*dialerCache]
// 		ctrl                  control.SocketController
// 		dialer                func(ctx context.Context, network, addr string) (Conn, error)
// 		der                   *net.Dialer
// 		tlsConfig             *tls.Config
// 		addrs                 sync.Map[string, *addrInfo]
// 		dnsRefreshDurationStr string
// 		dnsCacheExpirationStr string
// 		dnsCacheExpiration    time.Duration
// 		dialerTimeout         time.Duration
// 		dialerKeepalive       time.Duration
// 		dialerFallbackDelay   time.Duration
// 		dnsRefreshDuration    time.Duration
// 		sockFlg               control.SocketFlag
// 		dialerDualStack       bool
// 		enableDNSCache        bool
// 	}
// 	type want struct {
// 		wantIps []string
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotIps []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotIps, w.wantIps) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIps, w.wantIps)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           host:"",
// 		       },
// 		       fields: fields {
// 		           dnsCache:nil,
// 		           ctrl:nil,
// 		           dialer:nil,
// 		           der:net.Dialer{},
// 		           tlsConfig:nil,
// 		           addrs:nil,
// 		           dnsRefreshDurationStr:"",
// 		           dnsCacheExpirationStr:"",
// 		           dnsCacheExpiration:nil,
// 		           dialerTimeout:nil,
// 		           dialerKeepalive:nil,
// 		           dialerFallbackDelay:nil,
// 		           dnsRefreshDuration:nil,
// 		           sockFlg:nil,
// 		           dialerDualStack:false,
// 		           enableDNSCache:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           ctx:nil,
// 		           host:"",
// 		           },
// 		           fields: fields {
// 		           dnsCache:nil,
// 		           ctrl:nil,
// 		           dialer:nil,
// 		           der:net.Dialer{},
// 		           tlsConfig:nil,
// 		           addrs:nil,
// 		           dnsRefreshDurationStr:"",
// 		           dnsCacheExpirationStr:"",
// 		           dnsCacheExpiration:nil,
// 		           dialerTimeout:nil,
// 		           dialerKeepalive:nil,
// 		           dialerFallbackDelay:nil,
// 		           dnsRefreshDuration:nil,
// 		           sockFlg:nil,
// 		           dialerDualStack:false,
// 		           enableDNSCache:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			d := &dialer{
// 				dnsCache:              test.fields.dnsCache,
// 				ctrl:                  test.fields.ctrl,
// 				dialer:                test.fields.dialer,
// 				der:                   test.fields.der,
// 				tlsConfig:             test.fields.tlsConfig,
// 				addrs:                 test.fields.addrs,
// 				dnsRefreshDurationStr: test.fields.dnsRefreshDurationStr,
// 				dnsCacheExpirationStr: test.fields.dnsCacheExpirationStr,
// 				dnsCacheExpiration:    test.fields.dnsCacheExpiration,
// 				dialerTimeout:         test.fields.dialerTimeout,
// 				dialerKeepalive:       test.fields.dialerKeepalive,
// 				dialerFallbackDelay:   test.fields.dialerFallbackDelay,
// 				dnsRefreshDuration:    test.fields.dnsRefreshDuration,
// 				sockFlg:               test.fields.sockFlg,
// 				dialerDualStack:       test.fields.dialerDualStack,
// 				enableDNSCache:        test.fields.enableDNSCache,
// 			}
//
// 			gotIps, err := d.lookupIPAddrs(test.args.ctx, test.args.host)
// 			if err := checkFunc(test.want, gotIps, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_isQUICDial(t *testing.T) {
// 	type args struct {
// 		network string
// 		addr    string
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		       args: args {
// 		           network:"",
// 		           addr:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           network:"",
// 		           addr:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := isQUICDial(test.args.network, test.args.addr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// Test_dialer_tlsHandshake_sharedConfigNotMutated verifies that concurrent TLS dials never
// write the shared tls.Config (the ServerName used to be assigned in place, which races with
// the handshakes reading the same config); the ServerName must be set on a per-dial clone.
func Test_dialer_tlsHandshake_sharedConfigNotMutated(t *testing.T) {
	t.Parallel()
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	pool := x509.NewCertPool()
	pool.AddCert(srv.Certificate())
	cfg := &tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
	}
	d := &dialer{
		tlsConfig: cfg,
		der:       &net.Dialer{Timeout: 10 * time.Second},
	}
	// the httptest certificate is valid for 127.0.0.1, so verification only succeeds
	// when the per-dial ServerName is derived from the address.
	addr := srv.Listener.Addr().String()

	const n = 16
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	errs := make(chan error, n)
	for range n {
		go func() {
			conn, err := d.dial(ctx, TCP.String(), addr)
			if err != nil {
				errs <- err
				return
			}
			if conn == nil {
				errs <- errors.New("dial returned nil conn without error")
				return
			}
			errs <- conn.Close()
		}()
	}
	for range n {
		select {
		case err := <-errs:
			if err != nil {
				t.Errorf("failed to dial tls: %v", err)
			}
		case <-ctx.Done():
			t.Fatalf("timed out waiting for tls dials: %v", ctx.Err())
		}
	}

	if cfg.ServerName != "" {
		t.Errorf("shared tls.Config was mutated, ServerName = %q", cfg.ServerName)
	}
}

// Test_dialer_cachedDialer_contextCancelled verifies that the DNS retry loop stops as soon as
// the context is cancelled instead of spinning until dialerTimeout.
func Test_dialer_cachedDialer_contextCancelled(t *testing.T) {
	t.Parallel()
	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSRefreshDuration("5s"),
		WithDNSCacheExpiration("10s"),
		WithDisableDialerDualStack(),
		WithDialerTimeout("2s"),
	)
	if err != nil {
		t.Fatalf("failed to initialize dialer: %v", err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("NewDialer return value Dialer is not *dialer: %v", der)
	}
	const (
		host = "cancelled.invalid"
		addr = host + ":80"
	)
	// pre-populate the parsed address so that the call goes straight to the retry loop
	// (a lookup failure for a non-IP host).
	d.addrs.Store(addr, &addrInfo{addr: addr, host: host, port: "80"})

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	type result struct {
		err     error
		elapsed time.Duration
	}
	done := make(chan result, 1)
	go func() {
		start := time.Now()
		conn, err := d.cachedDialer(ctx, TCP.String(), addr)
		if conn != nil {
			if cerr := conn.Close(); cerr != nil {
				err = errors.Join(err, cerr)
			}
		}
		done <- result{err: err, elapsed: time.Since(start)}
	}()

	select {
	case r := <-done:
		if !errors.Is(r.err, context.Canceled) {
			t.Errorf("got error: %v, want: %v", r.err, context.Canceled)
		}
		// dialerTimeout is 2s; returning well before it proves the loop did not spin on.
		if r.elapsed >= time.Second {
			t.Errorf("cachedDialer kept retrying after ctx cancellation for %v", r.elapsed)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("timed out waiting for cachedDialer to return")
	}
}

// Test_dialer_dial_tlsHandshakeFailureReturnsNilConn verifies that a failed TLS dial returns
// a nil Conn interface. dial used to return tlsHandshake's typed nil *tls.Conn through the
// Conn interface, which made the conn non-nil and cachedDialer's conn.Close() panic.
func Test_dialer_dial_tlsHandshakeFailureReturnsNilConn(t *testing.T) {
	t.Parallel()
	ln, err := net.Listen(TCP.String(), "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	// accept every connection and close it at once, so every TLS handshake fails.
	acceptDone := make(chan struct{})
	defer func() {
		if err := ln.Close(); err != nil {
			t.Errorf("failed to close listener: %v", err)
		}
		// wait for the accept loop so it never logs after the test has returned.
		select {
		case <-acceptDone:
		case <-time.After(10 * time.Second):
			t.Error("timed out waiting for the accept loop to stop")
		}
	}()
	go func() {
		defer close(acceptDone)
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			if err := c.Close(); err != nil {
				t.Logf("failed to close accepted conn: %v", err)
			}
		}
	}()

	d := &dialer{
		tlsConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		der:       &net.Dialer{Timeout: 5 * time.Second},
	}
	type result struct {
		conn Conn
		err  error
	}
	done := make(chan result, 1)
	go func() {
		conn, err := d.dial(t.Context(), TCP.String(), ln.Addr().String())
		done <- result{conn: conn, err: err}
	}()
	select {
	case r := <-done:
		if r.err == nil {
			t.Fatal("dial() error = nil, want a tls handshake error")
		}
		if r.conn != nil {
			t.Errorf("dial() conn = %#v (non-nil interface), want nil", r.conn)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("timed out waiting for dial to return")
	}
}

// Test_dialer_dial_udpSkipsTLS verifies that a udp dial (the resolver dials DNS servers with
// "udp" to port 53) is not TLS-handshaked when TLS is enabled. dial rewrites network to tcp
// for these dials, and the TLS guard used to read the rewritten network.
// Only port 53 takes this path (other udp dials are QUIC), so the test needs to bind a
// privileged port and is skipped where that is not allowed.
func Test_dialer_dial_udpSkipsTLS(t *testing.T) {
	t.Parallel()
	ln, err := net.Listen(TCP.String(), "127.0.0.1:53")
	if err != nil {
		t.Skipf("cannot listen on 127.0.0.1:53 (needs CAP_NET_BIND_SERVICE or a free port 53): %v", err)
	}
	// accept every connection and close it at once, so any TLS handshake fails.
	acceptDone := make(chan struct{})
	defer func() {
		if err := ln.Close(); err != nil {
			t.Errorf("failed to close listener: %v", err)
		}
		select {
		case <-acceptDone:
		case <-time.After(10 * time.Second):
			t.Error("timed out waiting for the accept loop to stop")
		}
	}()
	go func() {
		defer close(acceptDone)
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			if err := c.Close(); err != nil {
				t.Logf("failed to close accepted conn: %v", err)
			}
		}
	}()

	d := &dialer{
		tlsConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		der:       &net.Dialer{Timeout: 5 * time.Second},
	}
	type result struct {
		conn Conn
		err  error
	}
	done := make(chan result, 1)
	go func() {
		conn, err := d.dial(t.Context(), UDP.String(), ln.Addr().String())
		done <- result{conn: conn, err: err}
	}()
	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("dial() error = %v, want nil (udp dials must not be TLS-handshaked)", r.err)
		}
		if r.conn == nil {
			t.Fatal("dial() conn = nil, want a connection")
		}
		defer func() {
			if err := r.conn.Close(); err != nil {
				t.Logf("failed to close conn: %v", err)
			}
		}()
		if _, ok := r.conn.(*tls.Conn); ok {
			t.Error("dial() returned a *tls.Conn for a udp dial, want a plain connection")
		}
	case <-time.After(30 * time.Second):
		t.Fatal("timed out waiting for dial to return")
	}
}

// NewDialer used to pass the keepalive duration to control.New as nanoseconds, while the
// socket controller sets TCP_KEEPIDLE/TCP_KEEPINTVL in seconds.
func TestNewDialer_keepaliveInSeconds(t *testing.T) {
	t.Parallel()
	der, err := NewDialer(WithDialerKeepalive("30s"))
	if err != nil {
		t.Fatalf("failed to initialize dialer: %v", err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("NewDialer return value Dialer is not *dialer: %v", der)
	}
	if want := control.New(d.sockFlg, 30); !reflect.DeepEqual(d.ctrl, want) {
		t.Errorf("ctrl = %#v, want %#v", d.ctrl, want)
	}
}

// Test_dialer_resolverDialSkipsTLS verifies that the resolver's dials never get TLS, while
// normal tcp dials keep it. The Go resolver falls back to tcp (to the DNS server, port 53)
// when a udp answer is truncated, and the resolver used to dial with d.dialer, so that tcp
// fallback was TLS-handshaked against the DNS server and failed.
func Test_dialer_resolverDialSkipsTLS(t *testing.T) {
	t.Parallel()
	ln, err := net.Listen(TCP.String(), "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	// accept every connection and close it at once, so any TLS handshake fails.
	acceptDone := make(chan struct{})
	defer func() {
		if cerr := ln.Close(); cerr != nil {
			t.Errorf("failed to close listener: %v", cerr)
		}
		select {
		case <-acceptDone:
		case <-time.After(10 * time.Second):
			t.Error("timed out waiting for the accept loop to stop")
		}
	}()
	go func() {
		defer close(acceptDone)
		for {
			c, aerr := ln.Accept()
			if aerr != nil {
				return
			}
			if cerr := c.Close(); cerr != nil {
				t.Logf("failed to close accepted conn: %v", cerr)
			}
		}
	}()

	der, err := NewDialer(
		WithTLS(&tls.Config{MinVersion: tls.VersionTLS12}),
		WithDialerTimeout("5s"),
	)
	if err != nil {
		t.Fatalf("failed to initialize dialer: %v", err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("NewDialer return value Dialer is not *dialer: %v", der)
	}
	if d.der == nil || d.der.Resolver == nil || d.der.Resolver.Dial == nil {
		t.Fatal("NewDialer did not set a resolver dial func")
	}

	type result struct {
		conn Conn
		err  error
	}
	tests := []struct {
		dial    func(ctx context.Context, network, addr string) (Conn, error)
		name    string
		wantTLS bool
	}{
		{
			dial:    d.der.Resolver.Dial,
			name:    "resolver tcp dial is plain",
			wantTLS: false,
		},
		{
			dial:    d.DialContext,
			name:    "normal tcp dial keeps TLS",
			wantTLS: true,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			done := make(chan result, 1)
			go func() {
				conn, err := test.dial(tt.Context(), TCP.String(), ln.Addr().String())
				done <- result{conn: conn, err: err}
			}()
			select {
			case r := <-done:
				if r.conn != nil {
					defer func() {
						if err := r.conn.Close(); err != nil {
							tt.Logf("failed to close conn: %v", err)
						}
					}()
				}
				if test.wantTLS {
					// the server closes without a handshake, so a TLS dial must fail.
					if r.err == nil {
						tt.Fatal("dial() error = nil, want a tls handshake error")
					}
					return
				}
				if r.err != nil {
					tt.Fatalf("dial() error = %v, want nil (resolver dials must not be TLS-handshaked)", r.err)
				}
				if r.conn == nil {
					tt.Fatal("dial() conn = nil, want a connection")
				}
				if _, ok := r.conn.(*tls.Conn); ok {
					tt.Error("dial() returned a *tls.Conn for a resolver dial, want a plain connection")
				}
			case <-time.After(30 * time.Second):
				tt.Fatal("timed out waiting for dial to return")
			}
		})
	}
}

// Test_dialer_EvictAddr verifies that EvictAddr drops the cached parse result
// of the address and the DNS cache entry of its host, and nothing else.
func Test_dialer_EvictAddr(t *testing.T) {
	t.Parallel()
	const (
		addr  = "10.0.0.1:8081"
		host  = "10.0.0.1"
		other = "10.0.0.2:8081"
		oHost = "10.0.0.2"
	)
	tests := []struct {
		wantAddrs map[string]bool
		wantDNS   map[string]bool
		name      string
		evict     string
	}{
		{
			name:      "known address",
			evict:     addr,
			wantAddrs: map[string]bool{addr: false, other: true},
			wantDNS:   map[string]bool{host: false, oHost: true},
		},
		{
			name:      "unknown address",
			evict:     "10.0.0.3:8081",
			wantAddrs: map[string]bool{addr: true, other: true},
			wantDNS:   map[string]bool{host: true, oHost: true},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			der, err := NewDialer(
				WithEnableDNSCache(),
				WithDNSRefreshDuration("1m"),
				WithDNSCacheExpiration("1m"),
			)
			if err != nil {
				tt.Fatal(err)
			}
			d, ok := der.(*dialer)
			if !ok {
				tt.Fatalf("NewDialer returned %T, want *dialer", der)
			}
			for a, h := range map[string]string{addr: host, other: oHost} {
				d.addrs.Store(a, &addrInfo{addr: a, host: h, port: "8081", isIP: true})
				d.dnsCache.Set(h, &dialerCache{ips: []string{h}})
			}

			d.EvictAddr(test.evict)

			for a, want := range test.wantAddrs {
				if _, ok := d.addrs.Load(a); ok != want {
					tt.Errorf("addrs[%q] present = %v, want %v", a, ok, want)
				}
			}
			for h, want := range test.wantDNS {
				if _, ok := d.dnsCache.Get(h); ok != want {
					tt.Errorf("dnsCache[%q] present = %v, want %v", h, ok, want)
				}
			}
		})
	}
}

// Test_dialer_cachedDialer_ipLiteralNotCached verifies that dialing an IP literal
// leaves no per-address state, while a host name keeps its cached parse result
// (Disconnect evicts it). gRPC hands the dialer every resolved backend IP:port,
// and those entries used to stay in addrs and the DNS cache forever.
func Test_dialer_cachedDialer_ipLiteralNotCached(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	// subtests run in parallel after this function returns, so close the
	// server only when they finish.
	t.Cleanup(srv.Close)
	ip, port, err := SplitHostPort(srv.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name       string
		host       string
		wantCached bool
	}{
		{
			name:       "ip literal is not cached",
			host:       ip,
			wantCached: false,
		},
		{
			name:       "host name stays cached",
			host:       "localhost",
			wantCached: true,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			der, err := NewDialer(
				WithEnableDNSCache(),
				WithDNSRefreshDuration("1m"),
				WithDNSCacheExpiration("1m"),
				WithDialerTimeout("10s"),
			)
			if err != nil {
				tt.Fatal(err)
			}
			d, ok := der.(*dialer)
			if !ok {
				tt.Fatalf("NewDialer returned %T, want *dialer", der)
			}
			addr := JoinHostPort(test.host, port)
			for range 2 {
				conn, err := d.cachedDialer(tt.Context(), TCP.String(), addr)
				if err != nil {
					tt.Fatalf("cachedDialer(%s) error = %v", addr, err)
				}
				if cerr := conn.Close(); cerr != nil {
					tt.Error(cerr)
				}
			}
			if _, ok := d.addrs.Load(addr); ok != test.wantCached {
				tt.Errorf("addrs[%q] present = %v, want %v", addr, ok, test.wantCached)
			}
			if !test.wantCached {
				if _, ok := d.dnsCache.Get(test.host); ok {
					tt.Errorf("dnsCache[%q] is present for an IP literal", test.host)
				}
			}
		})
	}
}

// Test_dialer_cachedDialer_tlsServerNameIsHost verifies that a TLS dial through
// the DNS cache sends the host name, not the resolved IP it connects to, as the
// TLS server name. It used to send the IP, so certificates issued for the host
// name failed verification on the cache path.
func Test_dialer_cachedDialer_tlsServerNameIsHost(t *testing.T) {
	t.Parallel()
	// The host never resolves: the cached parse result and DNS entry set below
	// route it to the test server, so the dial takes the DNS-cache path.
	const host = "vald-sni-test.invalid"
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil && r.TLS.ServerName == host {
			fmt.Fprint(w, "sni=host")
		}
	}))
	defer srv.Close()
	ip, port, err := SplitHostPort(srv.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := tls.NewClientConfig(tls.WithInsecureSkipVerify(true))
	if err != nil {
		t.Fatal(err)
	}
	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSRefreshDuration("1m"),
		WithDNSCacheExpiration("1m"),
		WithDialerTimeout("10s"),
		WithTLS(cfg),
	)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("NewDialer returned %T, want *dialer", der)
	}
	addr := JoinHostPort(host, port)
	d.addrs.Store(addr, &addrInfo{addr: addr, host: host, port: strconv.Itoa(int(port))})
	d.dnsCache.Set(host, &dialerCache{ips: []string{ip}})

	conn, err := d.cachedDialer(t.Context(), TCP.String(), addr)
	if err != nil {
		t.Fatalf("cachedDialer(%s) error = %v", addr, err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			t.Error(cerr)
		}
	}()
	if err = conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		t.Fatal(err)
	}
	if _, err = fmt.Fprint(conn, "GET / HTTP/1.0\r\n\r\n"); err != nil {
		t.Fatal(err)
	}
	buf, err := io.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}
	if got := conv.Btoa(buf); !strings.HasSuffix(got, "sni=host") {
		t.Errorf("server saw response %q, want TLS server name %q", got, host)
	}
}
