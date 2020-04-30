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

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"

	"go.uber.org/goleak"
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotDer, w.wantDer) {
			return errors.Errorf("got = %v, want %v", gotDer, w.wantDer)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts: nil,
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
		           opts: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
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
			defer goleak.VerifyNone(t)
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
		wantIps map[int]string
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, map[int]string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotIps map[int]string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotIps, w.wantIps) {
			return errors.Errorf("got = %v, want %v", gotIps, w.wantIps)
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
		           ctx: nil,
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
			defer goleak.VerifyNone(t)
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
			if err := test.checkFunc(test.want, gotIps, err); err != nil {
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
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
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
			defer goleak.VerifyNone(t)
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
			if err := test.checkFunc(test.want); err != nil {
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
			defer goleak.VerifyNone(t)
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
			defer goleak.VerifyNone(t)
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
