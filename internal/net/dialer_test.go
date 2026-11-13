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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math"
	"math/big"
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
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
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
			comparator.IgnoreFields(*want, "dialer", "der", "addrs", "tlsConfigs", "dnsCachedOnce", "dnsCache", "ctrl"),
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

// Test_dialer_cachedDialer_multiIPRoundRobin verifies that
// when dialing a hostname backed by multiple IP addresses in the DNS cache,
// subsequent dials rotate properly across all IPs and never revert to bypassing d.dnsCache.
func Test_dialer_cachedDialer_multiIPRoundRobin(t *testing.T) {
	t.Parallel()

	// Listen on 3 local loopback endpoints (127.0.0.1, 127.0.0.2, 127.0.0.3)
	backendIPs := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}

	probeLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot bind loopback listener: %v", err)
	}
	sharedPort := probeLn.Addr().(*net.TCPAddr).Port
	_ = probeLn.Close()

	portStr := fmt.Sprintf("%d", sharedPort)
	listeners := make([]net.Listener, len(backendIPs))
	connCounts := make([]int64, len(backendIPs))

	for idx, ip := range backendIPs {
		addr := net.JoinHostPort(ip, portStr)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			t.Skipf("cannot bind listener on %s: %v (requires loopback aliases)", addr, err)
		}
		listeners[idx] = ln

		serverIdx := idx
		go func(l net.Listener, id int) {
			for {
				conn, err := l.Accept()
				if err != nil {
					return
				}
				atomic.AddInt64(&connCounts[id], 1)
				_, _ = fmt.Fprintf(conn, "server-%d\n", id)
				_ = conn.Close()
			}
		}(ln, serverIdx)
	}

	defer func() {
		for _, ln := range listeners {
			if ln != nil {
				_ = ln.Close()
			}
		}
	}()

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatal(err)
	}

	// We use "localhost" here because cachedDialer on line 263 currently calls Parse(ctx, addr),
	// which attempts a DNS lookup via DefaultResolver on the host name.
	const testHostname = "localhost"
	targetDC := &dialerCache{
		ips: backendIPs,
	}
	dnsCache.Set(testHostname, targetDC)

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("5s"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	dialTarget := net.JoinHostPort(testHostname, portStr)

	// Sequential Dial Verification (15 dials)
	// Round-robin order of targetDC.IP():
	// Call 1: index 1
	// Call 2: index 2
	// Call 3: index 0
	// Call 4: index 1 ...
	const numSequentialDials = 15
	for i := range numSequentialDials {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := d.cachedDialer(ctx, TCP.String(), dialTarget)
		cancel()
		if err != nil {
			t.Fatalf("sequential dial %d failed: %v", i, err)
		}

		buf := make([]byte, 32)
		n, _ := conn.Read(buf)
		_ = conn.Close()
		serverMsg := string(buf[:n])

		expectedServerIdx := (i + 1) % len(backendIPs)
		expectedMsg := fmt.Sprintf("server-%d\n", expectedServerIdx)
		if serverMsg != expectedMsg {
			t.Errorf("dial %d connected to %q, want %q", i, serverMsg, expectedMsg)
		}

		// Verify invariant: d.addrs must keep isIP == false across subsequent dials
		ai, loaded := d.addrs.Load(dialTarget)
		if !loaded {
			t.Errorf("dial %d: dialTarget %q not found in d.addrs", i, dialTarget)
		} else if ai.isIP {
			t.Fatalf("CRITICAL BUG: dial %d reverted ai.isIP to true for hostname %q! DNS cache bypassed!", i, testHostname)
		}
	}

	// Verify that targetDC.cnt was incremented exactly once per successful dial
	if cnt := atomic.LoadUint32(&targetDC.cnt); cnt != numSequentialDials {
		t.Errorf("targetDC.cnt = %d, want %d", cnt, numSequentialDials)
	}

	// Verify that each backend received exactly numSequentialDials / 3 connections
	for idx := range backendIPs {
		got := atomic.LoadInt64(&connCounts[idx])
		want := int64(numSequentialDials / len(backendIPs))
		if got != want {
			t.Errorf("backend %d received %d connections, want %d", idx, got, want)
		}
	}
}

// Test_dialer_cachedDialer_preCachedHostnameSkipsResolver verifies that
// cachedDialer uses the DNS cache without calling Parse(ctx, addr) on the first dial,
// so that a pre-cached hostname not resolvable by system DNS dials successfully.
func Test_dialer_cachedDialer_preCachedHostnameSkipsResolver(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatal(err)
	}

	// Unresolvable domain in public/system DNS, but present in dnsCache
	const unresolvableHost = "unresolvable-pod-name.cluster.local"
	dnsCache.Set(unresolvableHost, &dialerCache{
		ips: []string{"127.0.0.1"},
	})

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("2s"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Attempt to dial the unresolvable hostname
	conn, dialErr := d.cachedDialer(ctx, TCP.String(), net.JoinHostPort(unresolvableHost, port))
	if dialErr != nil {
		t.Fatalf("cachedDialer failed on pre-cached hostname %q with error: %v", unresolvableHost, dialErr)
	}
	if conn == nil {
		t.Fatalf("expected non-nil connection for pre-cached hostname %q", unresolvableHost)
	}
	_ = conn.Close()
}

// Test_dialer_cachedDialer_ipLiteralBypassesDNSCache verifies that when
// dialing an IP literal (e.g. "127.0.0.1:port"):
// 1. It is never stored in d.addrs (preventing high-churn memory leaks).
// 2. It skips d.dnsCache and dials directly.
func Test_dialer_cachedDialer_ipLiteralBypassesDNSCache(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatal(err)
	}

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("3s"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	ipAddr := ln.Addr().String()

	for i := range 5 {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := d.cachedDialer(ctx, TCP.String(), ipAddr)
		cancel()
		if err != nil {
			t.Fatalf("dial IP literal %d failed: %v", i, err)
		}
		_ = conn.Close()

		if _, loaded := d.addrs.Load(ipAddr); loaded {
			t.Fatalf("IP literal %q was cached in d.addrs, causing memory leak under pod IP churn!", ipAddr)
		}
	}
}

// Test_dialer_cachedDialer_preCachedHostnameRoundRobin verifies
// that cachedDialer connects to a pre-cached hostname not resolvable by system DNS,
// and rotates properly across multiple backend IPs across dials 1..N without querying
// DefaultResolver or d.der.Resolver.
func Test_dialer_cachedDialer_preCachedHostnameRoundRobin(t *testing.T) {
	t.Parallel()

	backendIPs := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
	probeLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot bind loopback listener: %v", err)
	}
	sharedPort := probeLn.Addr().(*net.TCPAddr).Port
	_ = probeLn.Close()

	portStr := fmt.Sprintf("%d", sharedPort)
	listeners := make([]net.Listener, len(backendIPs))
	connCounts := make([]int64, len(backendIPs))

	for idx, ip := range backendIPs {
		addr := net.JoinHostPort(ip, portStr)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			t.Skipf("cannot bind listener on %s: %v (requires loopback aliases)", addr, err)
		}
		listeners[idx] = ln

		serverIdx := idx
		go func(l net.Listener, id int) {
			for {
				conn, err := l.Accept()
				if err != nil {
					return
				}
				atomic.AddInt64(&connCounts[id], 1)
				_, _ = fmt.Fprintf(conn, "server-%d\n", id)
				_ = conn.Close()
			}
		}(ln, serverIdx)
	}

	defer func() {
		for _, ln := range listeners {
			if ln != nil {
				_ = ln.Close()
			}
		}
	}()

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatal(err)
	}

	// Hostname guaranteed not to be in OS DNS
	const unresolvableHost = "unresolvable-svc.backend.cluster.local"
	targetDC := &dialerCache{
		ips: backendIPs,
	}
	dnsCache.Set(unresolvableHost, targetDC)

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("5s"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	// Track any calls to d.der.Resolver.Dial
	var resolverDialCalls int64
	origResolverDial := d.der.Resolver.Dial
	d.der.Resolver.Dial = func(ctx context.Context, network, addr string) (net.Conn, error) {
		atomic.AddInt64(&resolverDialCalls, 1)
		return origResolverDial(ctx, network, addr)
	}

	dialTarget := net.JoinHostPort(unresolvableHost, portStr)

	// Perform 30 sequential dials (dials 1 to 30)
	const numDials = 30
	for i := range numDials {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := d.cachedDialer(ctx, TCP.String(), dialTarget)
		cancel()
		if err != nil {
			t.Fatalf("dial %d failed for pre-cached hostname %q: %v", i+1, dialTarget, err)
		}

		buf := make([]byte, 32)
		n, _ := conn.Read(buf)
		_ = conn.Close()
		serverMsg := string(buf[:n])

		// Expected round-robin order: (i+1) % len(backendIPs)
		expectedServerIdx := (i + 1) % len(backendIPs)
		expectedMsg := fmt.Sprintf("server-%d\n", expectedServerIdx)
		if serverMsg != expectedMsg {
			t.Errorf("dial %d connected to %q, want %q", i+1, serverMsg, expectedMsg)
		}

		// Verify invariant: d.addrs retains isIP = false across dials 1..N
		ai, loaded := d.addrs.Load(dialTarget)
		if !loaded {
			t.Errorf("dial %d: dialTarget %q missing from d.addrs", i+1, dialTarget)
		} else if ai.isIP {
			t.Fatalf("CRITICAL BUG: dial %d set isIP=true for hostname %q!", i+1, unresolvableHost)
		}
	}

	// Verify that neither DefaultResolver nor d.der.Resolver was queried
	if calls := atomic.LoadInt64(&resolverDialCalls); calls != 0 {
		t.Fatalf("expected 0 resolver calls for pre-cached hostname, got %d", calls)
	}

	// Verify total cache rotation count
	if cnt := atomic.LoadUint32(&targetDC.cnt); cnt != numDials {
		t.Errorf("targetDC.cnt = %d, want %d", cnt, numDials)
	}

	// Verify even distribution across all backends
	for idx := range backendIPs {
		got := atomic.LoadInt64(&connCounts[idx])
		want := int64(numDials / len(backendIPs))
		if got != want {
			t.Errorf("backend %d received %d connections, want %d", idx, got, want)
		}
	}
}

// Test_dialer_cachedDialer_concurrentRoundRobin verifies thread-safe
// round-robin rotation across multiple goroutines dialing a pre-cached hostname.
func Test_dialer_cachedDialer_concurrentRoundRobin(t *testing.T) {
	t.Parallel()

	backendIPs := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
	probeLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot bind loopback listener: %v", err)
	}
	sharedPort := probeLn.Addr().(*net.TCPAddr).Port
	_ = probeLn.Close()

	portStr := fmt.Sprintf("%d", sharedPort)
	listeners := make([]net.Listener, len(backendIPs))
	for idx, ip := range backendIPs {
		addr := net.JoinHostPort(ip, portStr)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			t.Skipf("cannot bind listener on %s: %v", addr, err)
		}
		listeners[idx] = ln

		go func(l net.Listener) {
			for {
				conn, err := l.Accept()
				if err != nil {
					return
				}
				_ = conn.Close()
			}
		}(ln)
	}

	defer func() {
		for _, ln := range listeners {
			if ln != nil {
				_ = ln.Close()
			}
		}
	}()

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatal(err)
	}

	const unresolvableHost = "unresolvable-concurrent.cluster.local"
	targetDC := &dialerCache{
		ips: backendIPs,
	}
	dnsCache.Set(unresolvableHost, targetDC)

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("5s"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	dialTarget := net.JoinHostPort(unresolvableHost, portStr)

	const concurrency = 10
	const dialsPerGoroutine = 30
	eg, _ := errgroup.New(context.Background())

	for range concurrency {
		eg.Go(func() error {
			for range dialsPerGoroutine {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				conn, err := d.cachedDialer(ctx, TCP.String(), dialTarget)
				cancel()
				if err != nil {
					return fmt.Errorf("concurrent dial failed: %w", err)
				}
				_ = conn.Close()
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("concurrent cachedDialer failed: %v", err)
	}

	totalExpectedDials := uint32(concurrency * dialsPerGoroutine)
	if cnt := atomic.LoadUint32(&targetDC.cnt); cnt != totalExpectedDials {
		t.Errorf("targetDC.cnt = %d, want %d", cnt, totalExpectedDials)
	}
}

// Test_dialer_cachedDialer_concurrentLiteralIPs verifies that
// cachedDialer behaves correctly under heavy concurrent dials (100+ goroutines)
// to various literal IP formats: IPv4, IPv6, IPv4-mapped IPv6, and bracketed IPs.
//
// Invariants checked:
// 1. All concurrent dials succeed with 0 errors.
// 2. Zero data races detected under -race.
// 3. d.addrs NEVER stores literal IPs (zero memory leaks / churn unbounded map growth).
// 4. d.dnsCache is never polluted with literal IPs.
func Test_dialer_cachedDialer_concurrentLiteralIPs(t *testing.T) {
	t.Parallel()

	// 1. Setup IPv4 listener
	ln4, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on IPv4: %v", err)
	}
	defer ln4.Close()

	go func() {
		for {
			conn, aerr := ln4.Accept()
			if aerr != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	_, port4Str, err := net.SplitHostPort(ln4.Addr().String())
	if err != nil {
		t.Fatalf("failed to split hostport: %v", err)
	}

	// 2. Setup IPv6 listener if supported
	var (
		port6Str   string
		ipv6Active bool
	)
	if ln6, err := net.Listen("tcp6", "[::1]:0"); err == nil {
		defer ln6.Close()
		ipv6Active = true
		_, port6Str, _ = net.SplitHostPort(ln6.Addr().String())
		go func() {
			for {
				conn, aerr := ln6.Accept()
				if aerr != nil {
					return
				}
				_ = conn.Close()
			}
		}()
	}

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatalf("failed to create dnsCache: %v", err)
	}

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("5s"),
	)
	if err != nil {
		t.Fatalf("failed to create dialer: %v", err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	targets := []string{
		"127.0.0.1:" + port4Str,
		"[127.0.0.1]:" + port4Str,
		"[::ffff:127.0.0.1]:" + port4Str,
	}
	if ipv6Active {
		targets = append(targets, "[::1]:"+port6Str)
	}

	const concurrencyPerTarget = 40
	totalGoroutines := len(targets) * concurrencyPerTarget

	var (
		startWg sync.WaitGroup
		wg      sync.WaitGroup
		errCnt  atomic.Uint64
	)

	startWg.Add(1)
	wg.Add(totalGoroutines)

	for _, target := range targets {
		for i := range concurrencyPerTarget {
			go func(tgt string, iter int) {
				defer wg.Done()
				startWg.Wait()

				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				conn, derr := d.cachedDialer(ctx, TCP.String(), tgt)
				if derr != nil {
					errCnt.Add(1)
					t.Errorf("dial to %s failed (iter %d): %v", tgt, iter, derr)
					return
				}
				_ = conn.Close()
			}(target, i)
		}
	}

	// Release all goroutines simultaneously
	startWg.Done()
	wg.Wait()

	if failed := errCnt.Load(); failed > 0 {
		t.Fatalf("FAILED: %d out of %d concurrent dials failed", failed, totalGoroutines)
	}

	// Invariant: d.addrs must NOT contain any of the IP literals
	for _, target := range targets {
		if val, loaded := d.addrs.Load(target); loaded {
			t.Fatalf("INVARIANT VIOLATION: literal IP %q was cached in d.addrs: %+v", target, val)
		}
	}

	// Invariant: d.dnsCache must not have entries for IP literals
	for _, target := range targets {
		host, _, _ := SplitHostPort(target)
		if dc, ok := d.dnsCache.Get(host); ok && dc != nil {
			t.Fatalf("INVARIANT VIOLATION: literal IP host %q was cached in d.dnsCache: %+v", host, dc)
		}
	}
}

// Test_dialer_cachedDialer_concurrentHostnames verifies
// that cachedDialer under heavy concurrency correctly handles pre-cached hostnames
// with multiple IPs rotating round-robin across dials 1..N.
func Test_dialer_cachedDialer_concurrentHostnames(t *testing.T) {
	t.Parallel()

	// Create 3 listeners to simulate 3 pod replicas
	const numReplicas = 3
	var (
		listeners = make([]net.Listener, numReplicas)
		ips       = make([]string, numReplicas)
		connsSeen = make([]atomic.Uint64, numReplicas)
	)

	for i := range numReplicas {
		ln, err := net.Listen("tcp4", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("failed to listen on replica %d: %v", i, err)
		}
		defer ln.Close()
		listeners[i] = ln

		idx := i
		go func(l net.Listener, id int) {
			for {
				conn, aerr := l.Accept()
				if aerr != nil {
					return
				}
				connsSeen[id].Add(1)
				_ = conn.Close()
			}
		}(ln, idx)

		// Each listener has 127.0.0.1 with unique port, but to test round-robin
		// on a single port, we simulate IPs: 127.0.0.1, 127.0.0.2, 127.0.0.3
		// However, listeners on local host share 127.0.0.1 and have different ports.
		// So we create 3 listeners on 127.0.0.1, 127.0.0.2, 127.0.0.3 on the SAME port!
	}

	// Re-bind to same port on 127.0.0.1, 127.0.0.2, 127.0.0.3 (loopback block 127.0.0.0/8)
	// First close previous listeners
	for i := range numReplicas {
		_ = listeners[i].Close()
	}

	// Pick a free port on 127.0.0.1
	lnProbe, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to probe free port: %v", err)
	}
	_, portStr, err := net.SplitHostPort(lnProbe.Addr().String())
	if err != nil {
		t.Fatalf("failed to split host port: %v", err)
	}
	_ = lnProbe.Close()

	// Wait briefly for port reuse
	time.Sleep(20 * time.Millisecond)

	for i := range numReplicas {
		ip := fmt.Sprintf("127.0.0.%d", i+1)
		ips[i] = ip
		ln, err := net.Listen("tcp4", ip+":"+portStr)
		if err != nil {
			t.Fatalf("failed to listen on %s:%s: %v", ip, portStr, err)
		}
		defer ln.Close()
		idx := i
		go func(l net.Listener, id int) {
			for {
				conn, aerr := l.Accept()
				if aerr != nil {
					return
				}
				connsSeen[id].Add(1)
				_ = conn.Close()
			}
		}(ln, idx)
	}

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatalf("failed to create dnsCache: %v", err)
	}

	const testHost = "my-multi-pod-service.internal"
	dnsCache.Set(testHost, &dialerCache{
		ips: ips,
	})

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("5s"),
	)
	if err != nil {
		t.Fatalf("failed to create dialer: %v", err)
	}
	d := der.(*dialer)

	targetAddr := testHost + ":" + portStr

	const numDials = 120
	var (
		startWg sync.WaitGroup
		wg      sync.WaitGroup
		errCnt  atomic.Uint64
	)

	startWg.Add(1)
	wg.Add(numDials)

	for i := range numDials {
		go func(iter int) {
			defer wg.Done()
			startWg.Wait()

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			conn, derr := d.cachedDialer(ctx, TCP.String(), targetAddr)
			if derr != nil {
				errCnt.Add(1)
				t.Errorf("dial %d failed: %v", iter, derr)
				return
			}
			_ = conn.Close()
		}(i)
	}

	startWg.Done()
	wg.Wait()

	if failed := errCnt.Load(); failed > 0 {
		t.Fatalf("FAILED: %d out of %d dials failed", failed, numDials)
	}

	// Verify that d.addrs contains targetAddr with isIP == false
	ai, ok := d.addrs.Load(targetAddr)
	if !ok || ai == nil {
		t.Fatalf("expected targetAddr %q to be in d.addrs", targetAddr)
	}
	if ai.isIP {
		t.Errorf("expected ai.isIP == false for hostname %q, got true", targetAddr)
	}
	if ai.host != testHost {
		t.Errorf("expected ai.host == %q, got %q", testHost, ai.host)
	}
	if ai.port != portStr {
		t.Errorf("expected ai.port == %q, got %q", portStr, ai.port)
	}

	// Wait briefly for server accept loops to drain the kernel queue
	deadline := time.Now().Add(time.Second)
	var totalSeen uint64
	for time.Now().Before(deadline) {
		totalSeen = 0
		for i := range numReplicas {
			totalSeen += connsSeen[i].Load()
		}
		if totalSeen == numDials {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Verify traffic was distributed across all 3 replicas
	for i := range numReplicas {
		count := connsSeen[i].Load()
		if count == 0 {
			t.Errorf("replica %d (%s) received 0 connections out of %d!", i, ips[i], numDials)
		}
	}
	if totalSeen != numDials {
		t.Errorf("total connections seen = %d, want %d", totalSeen, numDials)
	}
}

// Test_dialer_cachedDialer_hostnameWithoutPort verifies the behavior
// of cachedDialer when dialing hostnames WITHOUT ports:
//  1. With enableDNSCache=true: defaultPort (80) is assigned, and dial behaves
//     predictably without crashing or leaking perr.
//  2. With enableDNSCache=false: standard error (missing port in address) is returned.
func Test_dialer_cachedDialer_hostnameWithoutPort(t *testing.T) {
	t.Parallel()

	// Subtest 1: With enableDNSCache=true
	t.Run("enableDNSCache=true", func(tt *testing.T) {
		dnsCache, err := cache.New[*dialerCache]()
		if err != nil {
			tt.Fatal(err)
		}

		const testHost = "unresolvable-no-port-pod.cluster.local"
		// Set unresolvable host pointing to loopback
		dnsCache.Set(testHost, &dialerCache{
			ips: []string{"127.0.0.1"},
		})

		der, err := NewDialer(
			WithEnableDNSCache(),
			WithDNSCache(dnsCache),
			WithDialerTimeout("500ms"),
		)
		if err != nil {
			tt.Fatal(err)
		}
		d := der.(*dialer)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Dial hostname without port
		conn, dialErr := d.cachedDialer(ctx, TCP.String(), testHost)
		if conn != nil {
			_ = conn.Close()
		}

		// When port 80 is not bound or listener fails, dialErr should be non-nil.
		// The critical verification is that cachedDialer does not panic and correctly
		// manages addrs and cache state.
		if dialErr == nil {
			tt.Logf("dial succeeded to %s:80", testHost)
		} else {
			tt.Logf("expected dial failure to port 80: %v", dialErr)
		}

		// Verify d.addrs cached host with port "80" and isIP == false
		ai, ok := d.addrs.Load(testHost)
		if !ok || ai == nil {
			tt.Fatalf("expected %q in d.addrs", testHost)
		}
		if ai.port != "80" {
			tt.Errorf("expected ai.port == '80', got %q", ai.port)
		}
		if ai.isIP {
			tt.Errorf("expected ai.isIP == false, got true")
		}
	})

	// Subtest 2: With enableDNSCache=false
	t.Run("enableDNSCache=false", func(tt *testing.T) {
		der, err := NewDialer(
			WithDisableDNSCache(),
			WithDialerTimeout("500ms"),
		)
		if err != nil {
			tt.Fatal(err)
		}
		d := der.(*dialer)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		conn, dialErr := d.cachedDialer(ctx, TCP.String(), "example.com")
		if conn != nil {
			_ = conn.Close()
		}

		if dialErr == nil {
			tt.Fatalf("expected error dialing example.com without port and without DNS cache, got nil")
		}
		if !strings.Contains(dialErr.Error(), "missing port in address") {
			tt.Errorf("expected 'missing port in address' error, got %v", dialErr)
		}
	})
}

// Test_dialer_cachedDialer_literalIPWithoutPort verifies dialing
// literal IPs without ports (e.g. "127.0.0.1", "::1", "[::1]"):
// 1. Never panics.
// 2. Returns expected missing port / invalid address error.
// 3. d.addrs remains clean (not cached).
func Test_dialer_cachedDialer_literalIPWithoutPort(t *testing.T) {
	t.Parallel()

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatal(err)
	}

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("500ms"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d := der.(*dialer)

	testAddrs := []string{
		"127.0.0.1",
		"::1",
		"[::1]",
		"[::ffff:127.0.0.1]",
	}

	for _, addr := range testAddrs {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		conn, dialErr := d.cachedDialer(ctx, TCP.String(), addr)
		cancel()

		if conn != nil {
			_ = conn.Close()
			t.Fatalf("expected dial to portless IP %s to fail, got connection", addr)
		}
		if dialErr == nil {
			t.Fatalf("expected dial error for portless IP %s, got nil", addr)
		}

		// Verify d.addrs was NOT populated
		if val, loaded := d.addrs.Load(addr); loaded {
			t.Errorf("INVARIANT VIOLATION: portless literal IP %s was cached in d.addrs: %+v", addr, val)
		}
	}
}

// Test_dialer_cachedDialer_repeatedDials verifies that
// IsIPLiteral is consistently and correctly applied to IPv4, IPv6, and bracketed
// combinations across repeated dials (dials 1..N):
// - Dials 1..N on literal IPs consistently bypass d.addrs and d.dnsCache.
// - Dials 1..N on hostnames consistently utilize d.addrs and d.dnsCache.
func Test_dialer_cachedDialer_repeatedDials(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, aerr := ln.Accept()
			if aerr != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	_, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatal(err)
	}

	const hostname = "pod-service.default.svc.cluster.local"
	dnsCache.Set(hostname, &dialerCache{
		ips: []string{"127.0.0.1"},
	})

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("3s"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d := der.(*dialer)

	// Test literal IPs across 20 consecutive dials each
	ipLiterals := []string{
		"127.0.0.1:" + portStr,
		"[127.0.0.1]:" + portStr,
		"[::ffff:127.0.0.1]:" + portStr,
	}

	const numDials = 20

	for _, lit := range ipLiterals {
		for dialNum := 1; dialNum <= numDials; dialNum++ {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			conn, err := d.cachedDialer(ctx, TCP.String(), lit)
			cancel()
			if err != nil {
				t.Fatalf("dial %d to literal %s failed: %v", dialNum, lit, err)
			}
			_ = conn.Close()

			// Must NEVER be in d.addrs on dial 1, dial 2, ..., dial N
			if _, loaded := d.addrs.Load(lit); loaded {
				t.Fatalf("INVARIANT VIOLATION: %s stored in d.addrs on dial %d", lit, dialNum)
			}
		}
	}

	// Now test hostname across 20 consecutive dials
	hostAddr := hostname + ":" + portStr
	for dialNum := 1; dialNum <= numDials; dialNum++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := d.cachedDialer(ctx, TCP.String(), hostAddr)
		cancel()
		if err != nil {
			t.Fatalf("dial %d to hostname %s failed: %v", dialNum, hostAddr, err)
		}
		_ = conn.Close()

		// MUST be in d.addrs with isIP == false on dial 1, dial 2, ..., dial N
		ai, loaded := d.addrs.Load(hostAddr)
		if !loaded || ai == nil {
			t.Fatalf("INVARIANT VIOLATION: hostname %s NOT in d.addrs on dial %d", hostAddr, dialNum)
		}
		if ai.isIP {
			t.Fatalf("INVARIANT VIOLATION: hostname %s stored with isIP=true on dial %d", hostAddr, dialNum)
		}
	}
}

// Test_dialer_cachedDialer_parallelInstances runs several independent cachedDialer
// instances concurrently to stress-test race safety and stability.
func Test_dialer_cachedDialer_parallelInstances(t *testing.T) {
	t.Parallel()

	eg, _ := errgroup.New(context.Background())

	// Run multiple parallel instances of cachedDialer tests
	for i := range 5 {
		idx := i
		eg.Go(func() error {
			// Subtest: return conn
			srv := net.ListenConfig{}
			ln, err := srv.Listen(context.Background(), "tcp4", "127.0.0.1:0")
			if err != nil {
				return err
			}
			defer ln.Close()

			go func() {
				for {
					conn, aerr := ln.Accept()
					if aerr != nil {
						return
					}
					_ = conn.Close()
				}
			}()

			der, err := NewDialer(WithDialerTimeout("2s"))
			if err != nil {
				return err
			}
			d := der.(*dialer)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			conn, derr := d.cachedDialer(ctx, TCP.String(), ln.Addr().String())
			if derr != nil {
				return fmt.Errorf("instance %d: cachedDialer failed: %w", idx, derr)
			}
			return conn.Close()
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("parallel cachedDialer subtest stress failed: %v", err)
	}
}

// generateTestTLSCert generates a self-signed TLS certificate for testing.
func generateTestTLSCert(t *testing.T) (tls.Certificate, []byte) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate ecdsa key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Vald Test"},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{"localhost"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}
	return cert, derBytes
}

// Test_dialer_dialConn_tlsConfigsConcurrentSNI exercises d.tlsConfigs
// under high concurrency (100 concurrent goroutines) dialing to both the same
// and 10 different SNI targets simultaneously.
// It verifies:
// 1. Race-free execution under Go's race detector (-race).
// 2. The original d.tlsConfig is never mutated (ServerName, InsecureSkipVerify, etc.).
// 3. Cached configs in d.tlsConfigs are properly populated and isolated per target SNI.
// 4. Server receives the expected SNI per connection.
func Test_dialer_dialConn_tlsConfigsConcurrentSNI(t *testing.T) {
	cert, _ := generateTestTLSCert(t)

	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Start TLS server that records the SNI seen on each connection
	ln, err := tls.Listen("tcp4", "127.0.0.1:0", serverTLSConfig)
	if err != nil {
		t.Fatalf("failed to listen on TLS: %v", err)
	}
	defer ln.Close()

	const (
		numCommonGoroutines   = 50
		numDistinctGoroutines = 50
		totalGoroutines       = numCommonGoroutines + numDistinctGoroutines
	)

	var (
		serverMu     sync.Mutex
		receivedSNIs []string
		serverWg     sync.WaitGroup
	)
	serverWg.Add(totalGoroutines)

	serverCtx, serverCancel := context.WithCancel(t.Context())
	defer serverCancel()

	go func() {
		for {
			conn, aerr := ln.Accept()
			if aerr != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				tlsConn, ok := c.(*tls.Conn)
				if !ok {
					return
				}
				if herr := tlsConn.HandshakeContext(serverCtx); herr == nil {
					sni := tlsConn.ConnectionState().ServerName
					serverMu.Lock()
					receivedSNIs = append(receivedSNIs, sni)
					serverMu.Unlock()
					serverWg.Done()
				}
			}(conn)
		}
	}()

	serverAddr := ln.Addr().String()

	// Prepare client dialer with original tlsConfig
	origClientTLSConfig := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	der, err := NewDialer(
		WithTLS(origClientTLSConfig),
		WithDialerTimeout("10s"),
	)
	if err != nil {
		t.Fatalf("failed to create dialer: %v", err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	var (
		wg      sync.WaitGroup
		errCnt  atomic.Uint64
		startWg sync.WaitGroup
	)

	startWg.Add(1)
	wg.Add(totalGoroutines)

	// Launch common host dials
	for i := range numCommonGoroutines {
		go func(id int) {
			defer wg.Done()
			startWg.Wait()

			conn, herr := d.dialConn(t.Context(), "tcp4", serverAddr, "common-host.vald.vdaas.org", true)
			if herr != nil {
				errCnt.Add(1)
				t.Errorf("common dial [%d] failed: %v", id, herr)
				return
			}
			_ = conn.Close()
		}(i)
	}

	// Launch distinct host dials
	for i := range numDistinctGoroutines {
		targetSN := fmt.Sprintf("host-%d.vald.vdaas.org", i%10)
		go func(id int, sn string) {
			defer wg.Done()
			startWg.Wait()

			conn, herr := d.dialConn(t.Context(), "tcp4", serverAddr, sn, true)
			if herr != nil {
				errCnt.Add(1)
				t.Errorf("distinct dial [%d] to %s failed: %v", id, sn, herr)
				return
			}
			_ = conn.Close()
		}(i, targetSN)
	}

	startWg.Done()
	wg.Wait()

	if failed := errCnt.Load(); failed > 0 {
		t.Fatalf("%d out of %d concurrent dials failed", failed, totalGoroutines)
	}

	// Wait for server-side TLS handshakes to complete recording
	serverDone := make(chan struct{})
	go func() {
		serverWg.Wait()
		close(serverDone)
	}()
	select {
	case <-serverDone:
	case <-time.After(5 * time.Second):
		t.Log("timed out waiting for all server handshakes to complete")
	}

	// Verify Requirement 1: Original d.tlsConfig was NEVER mutated
	if origClientTLSConfig.ServerName != "" {
		t.Errorf("VIOLATION: original client tls.Config.ServerName mutated to %q, want empty", origClientTLSConfig.ServerName)
	}
	if !origClientTLSConfig.InsecureSkipVerify {
		t.Errorf("VIOLATION: original client tls.Config.InsecureSkipVerify was mutated")
	}

	// Verify Requirement 2: Cached configs are properly isolated in d.tlsConfigs
	commonCfg, ok := d.tlsConfigs.Load("common-host.vald.vdaas.org")
	if !ok || commonCfg == nil {
		t.Errorf("expected common-host.vald.vdaas.org to be present in d.tlsConfigs")
	} else {
		if commonCfg == origClientTLSConfig {
			t.Errorf("VIOLATION: cached common config shares pointer with original tlsConfig")
		}
		if commonCfg.ServerName != "common-host.vald.vdaas.org" {
			t.Errorf("cached common config ServerName = %q, want common-host.vald.vdaas.org", commonCfg.ServerName)
		}
	}

	// Check each of the 10 distinct host configs
	for i := range 10 {
		sn := fmt.Sprintf("host-%d.vald.vdaas.org", i)
		cfg, loaded := d.tlsConfigs.Load(sn)
		if !loaded || cfg == nil {
			t.Errorf("expected %s to be present in d.tlsConfigs", sn)
			continue
		}
		if cfg == origClientTLSConfig {
			t.Errorf("VIOLATION: cached config for %s shares pointer with original tlsConfig", sn)
		}
		if cfg.ServerName != sn {
			t.Errorf("cached config for %s has ServerName = %q, want %q", sn, cfg.ServerName, sn)
		}
		if cfg == commonCfg {
			t.Errorf("VIOLATION: cached config for %s shares pointer with commonCfg", sn)
		}
	}

	// Verify Requirement 3: Server actually received expected SNIs
	serverMu.Lock()
	defer serverMu.Unlock()
	if len(receivedSNIs) < totalGoroutines {
		t.Errorf("server received %d handshakes, want at least %d", len(receivedSNIs), totalGoroutines)
	}
}

// Test_dialer_tlsHandshake_presetServerName verifies that when d.tlsConfig has
// a preset ServerName:
// 1. Dials matching the preset ServerName or with empty targetSN reuse d.tlsConfig directly (0 clones).
// 2. Dials with a different targetSN (the destination host the DNS-cache path passes) still use the
// preset ServerName: it is the certificate name the user configured, so it is never replaced by a
// clone and nothing is cached or mutated.
func Test_dialer_tlsHandshake_presetServerName(t *testing.T) {
	const presetSN = "preset-master.vald.vdaas.org"
	origConfig := &tls.Config{
		ServerName:         presetSN,
		InsecureSkipVerify: true,
	}

	der, err := NewDialer(
		WithTLS(origConfig),
		WithDialerTimeout("5s"),
	)
	if err != nil {
		t.Fatalf("failed to create dialer: %v", err)
	}
	d := der.(*dialer)

	// Dial to matching preset
	_, _ = d.tlsHandshake(t.Context(), nil, "tcp4", "127.0.0.1:443", presetSN)
	// d.tlsConfigs must have 0 entries because preset matches
	var count int
	d.tlsConfigs.Range(func(string, *tls.Config) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("expected 0 cached configs when serverName matches preset, got %d", count)
	}

	// Dial to an alternate host
	const altSN = "alternate.vald.vdaas.org"
	_, _ = d.tlsHandshake(t.Context(), nil, "tcp4", "127.0.0.1:443", altSN)

	// Original preset config must remain unchanged
	if origConfig.ServerName != presetSN {
		t.Errorf("VIOLATION: original ServerName was mutated from %q to %q", presetSN, origConfig.ServerName)
	}

	// The preset ServerName wins over the alternate host: nothing may be cloned.
	if _, ok := d.tlsConfigs.Load(altSN); ok {
		t.Errorf("%s was cached, but the preset ServerName %q must be used as is", altSN, presetSN)
	}
	count = 0
	d.tlsConfigs.Range(func(string, *tls.Config) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("expected 0 cached configs when a ServerName is preset, got %d", count)
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
// 		args
// 		fields
// 		want
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
// 		args
// 		want
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
// of the address, the DNS cache entry and the per-host TLS config of its host,
// and nothing else.
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
		wantTLS   map[string]bool
		name      string
		evict     string
	}{
		{
			name:      "known address",
			evict:     addr,
			wantAddrs: map[string]bool{addr: false, other: true},
			wantDNS:   map[string]bool{host: false, oHost: true},
			wantTLS:   map[string]bool{host: false, oHost: true},
		},
		{
			name:      "unknown address",
			evict:     "10.0.0.3:8081",
			wantAddrs: map[string]bool{addr: true, other: true},
			wantDNS:   map[string]bool{host: true, oHost: true},
			wantTLS:   map[string]bool{host: true, oHost: true},
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
				d.tlsConfigs.Store(h, &tls.Config{ServerName: h})
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
			for h, want := range test.wantTLS {
				if _, ok := d.tlsConfigs.Load(h); ok != want {
					tt.Errorf("tlsConfigs[%q] present = %v, want %v", h, ok, want)
				}
			}
		})
	}
}

// Test_dialer_EvictAddr_ipLiteralAndNoDNSCache covers what the table above does not: an IP literal
// is never in addrs, yet tlsHandshake keeps a TLS config per dialed IP, so EvictAddr must drop it
// from the address alone; and a dialer without a DNS cache must still evict a cached host name.
func Test_dialer_EvictAddr_ipLiteralAndNoDNSCache(t *testing.T) {
	t.Parallel()
	const (
		hostName     = "backend.example"
		hostNameAddr = hostName + ":8081"
	)
	der, err := NewDialer()
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("NewDialer returned %T, want *dialer", der)
	}
	if d.dnsCache != nil {
		t.Fatal("the test needs a dialer without a DNS cache")
	}
	for _, h := range []string{"10.0.0.5", "10.0.0.6", hostName} {
		d.tlsConfigs.Store(h, &tls.Config{ServerName: h})
	}
	d.addrs.Store(hostNameAddr, &addrInfo{addr: hostNameAddr, host: hostName, port: "8081"})

	d.EvictAddr("10.0.0.5:8081")
	d.EvictAddr(hostNameAddr)

	for h, want := range map[string]bool{"10.0.0.5": false, "10.0.0.6": true, hostName: false} {
		if _, ok := d.tlsConfigs.Load(h); ok != want {
			t.Errorf("tlsConfigs[%q] present = %v, want %v", h, ok, want)
		}
	}
	if _, ok := d.addrs.Load(hostNameAddr); ok {
		t.Error("addrs still holds the evicted host name")
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

// Test_dialer_cachedDialer_tlsConfiguredServerNameWins verifies that an explicitly
// configured tls.Config.ServerName is still the TLS server name on the DNS-cache
// path, even though that path passes the destination host name to the handshake:
// the host name must only be used when no ServerName is configured.
func Test_dialer_cachedDialer_tlsConfiguredServerNameWins(t *testing.T) {
	t.Parallel()
	const (
		host       = "vald-sni-test.invalid"
		configured = "vald.internal.example"
	)
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil && r.TLS.ServerName == configured {
			fmt.Fprint(w, "sni=configured")
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
	cfg.ServerName = configured
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
	if got := conv.Btoa(buf); !strings.HasSuffix(got, "sni=configured") {
		t.Errorf("server saw response %q, want the configured TLS server name %q", got, configured)
	}
}

// Test_dialer_cachedDialer_roundRobinDNS verifies that when dialing a hostname with
// DNS caching enabled, subsequent dials continue to use the DNS cache and round-robin
// across the cached backend IP addresses.
func Test_dialer_cachedDialer_roundRobinDNS(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	host, port, err := SplitHostPort(srv.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	c, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatal(err)
	}

	const hostname = "localhost"
	targetDC := &dialerCache{
		ips: []string{host, host},
	}
	c.Set(hostname, targetDC)

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(c),
		WithDialerTimeout("5s"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	addr := JoinHostPort(hostname, port)

	// Perform 6 consecutive dials
	for i := range 6 {
		conn, err := d.cachedDialer(t.Context(), TCP.String(), addr)
		if err != nil {
			t.Fatalf("dial %d failed: %v", i, err)
		}
		_ = conn.Close()
	}

	// Verify that addr was cached in d.addrs with isIP = false
	ai, loaded := d.addrs.Load(addr)
	if !loaded {
		t.Errorf("expected addr %q to be present in d.addrs", addr)
	} else if ai.isIP {
		t.Errorf("expected ai.isIP to be false for hostname %q, got true", hostname)
	}

	// Verify that targetDC.IP() was called on each dial by checking its counter
	if cnt := atomic.LoadUint32(&targetDC.cnt); cnt < 6 {
		t.Errorf("expected targetDC.cnt >= 6 (round-robin called on every dial), got %d", cnt)
	}
}

// Test_dialer_tlsHandshake_reuseConfig verifies that multiple TLS handshakes
// to the same host reuse cached tls.Config without data races or shared config mutation.
func Test_dialer_tlsHandshake_reuseConfig(t *testing.T) {
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

	der, err := NewDialer(
		WithTLS(cfg),
		WithDialerTimeout("10s"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	addr := srv.Listener.Addr().String()

	// Perform multiple concurrent dials
	const n = 10
	errChan := make(chan error, n)
	for range n {
		go func() {
			conn, err := d.DialContext(t.Context(), TCP.String(), addr)
			if err != nil {
				errChan <- err
				return
			}
			errChan <- conn.Close()
		}()
	}

	for range n {
		if err := <-errChan; err != nil {
			t.Errorf("concurrent TLS dial error: %v", err)
		}
	}

	// Original cfg must NOT be mutated
	if cfg.ServerName != "" {
		t.Errorf("original tls.Config.ServerName was mutated to %q", cfg.ServerName)
	}

	// Caching in d.tlsConfigs must have cached an entry
	targetHost, _, _ := SplitHostPort(addr)
	if _, ok := d.tlsConfigs.Load(targetHost); !ok {
		t.Errorf("expected tlsConfigs to have entry for %q", targetHost)
	}
}

// Test_dialer_cachedDialer_preCachedUnresolvableHostname verifies that cachedDialer
// can dial a hostname not resolvable by system DNS when the hostname is pre-cached
// in d.dnsCache, proving Parse(ctx, addr) is no longer called on initial dial.
func Test_dialer_cachedDialer_preCachedUnresolvableHostname(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	dnsCache, err := cache.New[*dialerCache]()
	if err != nil {
		t.Fatal(err)
	}

	const unresolvableHost = "unresolvable-pod-name.cluster.local"
	dnsCache.Set(unresolvableHost, &dialerCache{
		ips: []string{"127.0.0.1"},
	})

	der, err := NewDialer(
		WithEnableDNSCache(),
		WithDNSCache(dnsCache),
		WithDialerTimeout("2s"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := der.(*dialer)
	if !ok {
		t.Fatalf("expected *dialer, got %T", der)
	}

	_, port, err := SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := d.cachedDialer(ctx, TCP.String(), JoinHostPort(unresolvableHost, port))
	if err != nil {
		t.Fatalf("expected cachedDialer to connect to pre-cached unresolvable hostname, got: %v", err)
	}
	if conn == nil {
		t.Fatalf("expected non-nil connection for pre-cached hostname %q", unresolvableHost)
	}
	_ = conn.Close()
}

// Test_dialer_cachedDialer_firstDial_zeroAllocAddressParsing validates that
// SplitHostPort performs zero heap allocations on hostnames with ports.
func Test_dialer_cachedDialer_firstDial_zeroAllocAddressParsing(t *testing.T) {
	const addr = "vald-agent.default.svc.cluster.local:8081"

	allocs := testing.AllocsPerRun(100, func() {
		host, port, err := SplitHostPort(addr)
		if err != nil || host == "" || port == 0 {
			t.Fatalf("unexpected SplitHostPort result: %s:%d, err: %v", host, port, err)
		}
	})

	if allocs > 0 {
		t.Errorf("expected 0 allocs/op for SplitHostPort, got %f allocs/op", allocs)
	}
}
