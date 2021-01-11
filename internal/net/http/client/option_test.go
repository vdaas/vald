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

package client

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
)

func TestWithProxy(t *testing.T) {
	type T = transport
	type args struct {
		px func(*http.Request) (*url.URL, error)
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		func() test {
			p := func(*http.Request) (*url.URL, error) {
				return nil, nil
			}
			return test{
				name: "set proxy success",
				args: args{
					px: p,
				},
				want: want{
					obj: &T{
						Transport: &http.Transport{
							Proxy: p,
						},
					},
				},
			}
		}(),
		{
			name: "return error when proxy is nil",
			args: args{
				px: nil,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrInvalidOption("proxy", nil),
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

			got := WithProxy(test.args.px)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDialContext(t *testing.T) {
	type T = transport
	type args struct {
		dx func(ctx context.Context, network, addr string) (net.Conn, error)
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		func() test {
			d := func(ctx context.Context, network, addr string) (net.Conn, error) {
				return nil, nil
			}
			return test{
				name: "set dial context success",
				args: args{
					dx: d,
				},
				want: want{
					obj: &T{
						Transport: &http.Transport{
							DialContext: d,
						},
					},
				},
			}
		}(),
		{
			name: "return error when dial context is nil",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrInvalidOption("dialContext", nil),
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

			got := WithDialContext(test.args.dx)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithTLSHandshakeTimeout(t *testing.T) {
	type T = transport
	type args struct {
		dur string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set timeout success",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						TLSHandshakeTimeout: 5 * time.Second,
					},
				},
			},
		},
		{
			name: "set timeout failed with invalid value",
			args: args{
				dur: "dummy",
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrCriticalOption("TLSHandshakeTimeout", "dummy", errors.New("invalid timeout value: dummy\t:timeout parse error out put failed: time: invalid duration \"dummy\"")),
			},
		},
		{
			name: "set timeout failed with empty value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrInvalidOption("TLSHandshakeTimeout", ""),
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

			got := WithTLSHandshakeTimeout(test.args.dur)
			obj := &T{
				Transport: new(http.Transport),
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithEnableKeepAlives(t *testing.T) {
	type T = transport
	type args struct {
		enable bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set enable success",
			args: args{
				enable: true,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						DisableKeepAlives: false,
					},
				},
			},
		},
		{
			name: "set disable success",
			args: args{
				enable: false,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						DisableKeepAlives: true,
					},
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

			got := WithEnableKeepAlives(test.args.enable)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithEnableCompression(t *testing.T) {
	type T = transport
	type args struct {
		enable bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set enable success",
			args: args{
				enable: true,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						DisableCompression: false,
					},
				},
			},
		},
		{
			name: "set disable success",
			args: args{
				enable: false,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						DisableCompression: true,
					},
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

			got := WithEnableCompression(test.args.enable)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxIdleConns(t *testing.T) {
	type T = transport
	type args struct {
		cn int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set conn success",
			args: args{
				cn: 5,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						MaxIdleConns: 5,
					},
				},
			},
		},
		{
			name: "set conn success with default value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						MaxIdleConns: 0,
					},
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

			got := WithMaxIdleConns(test.args.cn)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxIdleConnsPerHost(t *testing.T) {
	type T = transport
	type args struct {
		cn int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set conn per host success",
			args: args{
				cn: 5,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						MaxIdleConnsPerHost: 5,
					},
				},
			},
		},
		{
			name: "set conn per host success with default value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						MaxIdleConnsPerHost: 0,
					},
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

			got := WithMaxIdleConnsPerHost(test.args.cn)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxConnsPerHost(t *testing.T) {
	type T = transport
	type args struct {
		cn int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set conn per host success",
			args: args{
				cn: 5,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						MaxConnsPerHost: 5,
					},
				},
			},
		},
		{
			name: "set conn per host success with default value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						MaxConnsPerHost: 0,
					},
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

			got := WithMaxConnsPerHost(test.args.cn)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithIdleConnTimeout(t *testing.T) {
	type T = transport
	type args struct {
		dur string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set timeout success",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						IdleConnTimeout: 5 * time.Second,
					},
				},
			},
		},
		{
			name: "set timeout failed with invalid value",
			args: args{
				dur: "dummy",
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrCriticalOption("idleConnTimeout", "dummy", errors.New("invalid timeout value: dummy\t:timeout parse error out put failed: time: invalid duration \"dummy\"")),
			},
		},
		{
			name: "set timeout failed with empty value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrInvalidOption("idleConnTimeout", ""),
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

			got := WithIdleConnTimeout(test.args.dur)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithResponseHeaderTimeout(t *testing.T) {
	type T = transport
	type args struct {
		dur string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set timeout success",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						ResponseHeaderTimeout: 5 * time.Second,
					},
				},
			},
		},
		{
			name: "set timeout failed with invalid value",
			args: args{
				dur: "dummy",
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrCriticalOption("responseHeaderTimeout", "dummy", errors.New("invalid timeout value: dummy\t:timeout parse error out put failed: time: invalid duration \"dummy\"")),
			},
		},
		{
			name: "set timeout failed with empty value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrInvalidOption("responseHeaderTimeout", ""),
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

			got := WithResponseHeaderTimeout(test.args.dur)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithExpectContinueTimeout(t *testing.T) {
	type T = transport
	type args struct {
		dur string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set timeout success",
			args: args{
				dur: "5s",
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						ExpectContinueTimeout: 5 * time.Second,
					},
				},
			},
		},
		{
			name: "set timeout failed with invalid value",
			args: args{
				dur: "dummy",
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrCriticalOption("expectContinueTimeout", "dummy", errors.New("invalid timeout value: dummy\t:timeout parse error out put failed: time: invalid duration \"dummy\"")),
			},
		},
		{
			name: "set timeout failed with empty value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrInvalidOption("expectContinueTimeout", ""),
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

			got := WithExpectContinueTimeout(test.args.dur)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithProxyConnectHeader(t *testing.T) {
	type T = transport
	type args struct {
		header http.Header
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set header success",
			args: args{
				header: http.Header(
					map[string][]string{"dummy": []string{"val"}},
				),
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						ProxyConnectHeader: http.Header(
							map[string][]string{"dummy": []string{"val"}},
						),
					},
				},
			},
		},
		{
			name: "return error when header is nil",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
				},
				err: errors.NewErrInvalidOption("proxyConnectHeader", http.Header(nil)),
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

			got := WithProxyConnectHeader(test.args.header)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxResponseHeaderBytes(t *testing.T) {
	type T = transport
	type args struct {
		bs int64
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set max response header byte",
			args: args{
				bs: 5,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						MaxResponseHeaderBytes: 5,
					},
				},
			},
		},
		{
			name: "set max response header byte with default value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						MaxResponseHeaderBytes: 0,
					},
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

			got := WithMaxResponseHeaderBytes(test.args.bs)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithWriteBufferSize(t *testing.T) {
	type T = transport
	type args struct {
		bs int64
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set write buffer size",
			args: args{
				bs: 5,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						WriteBufferSize: 5,
					},
				},
			},
		},
		{
			name: "set write buffer size with default value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						WriteBufferSize: 0,
					},
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

			got := WithWriteBufferSize(test.args.bs)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithReadBufferSize(t *testing.T) {
	type T = transport
	type args struct {
		bs int64
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set buffer size success",
			args: args{
				bs: 5,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						ReadBufferSize: 5,
					},
				},
			},
		},
		{
			name: "set buffer size success with default value",
			args: args{},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						ReadBufferSize: 0,
					},
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

			got := WithReadBufferSize(test.args.bs)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithForceAttemptHTTP2(t *testing.T) {
	type T = transport
	type args struct {
		force bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set force http2 success",
			args: args{
				force: true,
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{
						ForceAttemptHTTP2: true,
					},
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

			got := WithForceAttemptHTTP2(test.args.force)
			obj := &T{
				Transport: &http.Transport{},
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBackoffOpts(t *testing.T) {
	type T = transport
	type args struct {
		opts []backoff.Option
	}
	type fields struct {
		backoffOpts []backoff.Option
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, transportComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set backoff opts success when origin backoff opts is nil",
			args: args{
				opts: []backoff.Option{backoff.WithEnableErrorLog()},
			},
			want: want{
				obj: &T{
					Transport:   &http.Transport{},
					backoffOpts: []backoff.Option{backoff.WithEnableErrorLog()},
				},
			},
		},
		{
			name: "set backoff opts success when origin backoff opts is not nil",
			args: args{
				opts: []backoff.Option{backoff.WithEnableErrorLog()},
			},
			fields: fields{
				backoffOpts: []backoff.Option{backoff.WithRetryCount(20)},
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
					backoffOpts: []backoff.Option{
						backoff.WithRetryCount(20),
						backoff.WithEnableErrorLog(),
					},
				},
			},
		},
		{
			name: "return error when opt is empty",
			args: args{},
			fields: fields{
				backoffOpts: []backoff.Option{backoff.WithRetryCount(20)},
			},
			want: want{
				obj: &T{
					Transport: &http.Transport{},
					backoffOpts: []backoff.Option{
						backoff.WithRetryCount(20),
					},
				},
				err: func() error {
					var bo []backoff.Option
					return errors.NewErrInvalidOption("backoffOpts", bo)
				}(),
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

			got := WithBackoffOpts(test.args.opts...)
			obj := &T{
				Transport:   &http.Transport{},
				backoffOpts: test.fields.backoffOpts,
			}
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
