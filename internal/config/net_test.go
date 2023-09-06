//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/control"
	testdata "github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestDNS_Bind(t *testing.T) {
	type fields struct {
		CacheEnabled    bool
		RefreshDuration string
		CacheExpiration string
	}
	type want struct {
		want *DNS
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *DNS) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *DNS) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			cacheEnabled := true
			refreshDuration := "10m"
			cacheExpiration := "24h"
			return test{
				name: "return DNS when all parameters are not nil or empty",
				fields: fields{
					CacheEnabled:    cacheEnabled,
					RefreshDuration: refreshDuration,
					CacheExpiration: cacheExpiration,
				},
				want: want{
					want: &DNS{
						CacheEnabled:    cacheEnabled,
						RefreshDuration: refreshDuration,
						CacheExpiration: cacheExpiration,
					},
				},
			}
		}(),
		func() test {
			cacheEnabled := true
			envPrefix := "DNS_BIND_"
			p := map[string]string{
				envPrefix + "REFRESH_DURATION": "10m",
				envPrefix + "CACHE_EXPIRATION": "24h",
			}
			return test{
				name: "return DNS when string values are set as environment value",
				fields: fields{
					CacheEnabled:    cacheEnabled,
					RefreshDuration: "_" + envPrefix + "REFRESH_DURATION_",
					CacheExpiration: "_" + envPrefix + "CACHE_EXPIRATION_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range p {
						t.Setenv(k, v)
					}
				},
				want: want{
					want: &DNS{
						CacheEnabled:    cacheEnabled,
						RefreshDuration: "10m",
						CacheExpiration: "24h",
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return DNS when no parameters are set",
				fields: fields{},
				want: want{
					want: &DNS{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			d := &DNS{
				CacheEnabled:    test.fields.CacheEnabled,
				RefreshDuration: test.fields.RefreshDuration,
				CacheExpiration: test.fields.CacheExpiration,
			}

			got := d.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDialer_Bind(t *testing.T) {
	type fields struct {
		Timeout          string
		Keepalive        string
		FallbackDelay    string
		DualStackEnabled bool
	}
	type want struct {
		want *Dialer
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Dialer) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Dialer) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			timeout := "3s"
			keepalive := "5m"
			fallbackDelay := "10m"
			dualStackEnabled := false
			return test{
				name: "return Dialer when fields are not empty",
				fields: fields{
					Timeout:          timeout,
					Keepalive:        keepalive,
					FallbackDelay:    fallbackDelay,
					DualStackEnabled: dualStackEnabled,
				},
				want: want{
					want: &Dialer{
						Timeout:          timeout,
						Keepalive:        keepalive,
						FallbackDelay:    fallbackDelay,
						DualStackEnabled: dualStackEnabled,
					},
				},
			}
		}(),
		func() test {
			envPrefix := "DIALER_BIND_"
			p := map[string]string{
				envPrefix + "TIMEOUT":          "3s",
				envPrefix + "KEEP_ALIVE":       "5m",
				envPrefix + "DUAL_STACK_DELAY": "10m",
			}
			return test{
				name: "return Dialer when fields are set as environment value",
				fields: fields{
					Timeout:       "_" + envPrefix + "TIMEOUT_",
					Keepalive:     "_" + envPrefix + "KEEP_ALIVE_",
					FallbackDelay: "_" + envPrefix + "DUAL_STACK_DELAY_",
				},
				want: want{
					want: &Dialer{
						Timeout:          "3s",
						Keepalive:        "5m",
						FallbackDelay:    "10m",
						DualStackEnabled: false,
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range p {
						t.Setenv(k, v)
					}
				},
			}
		}(),
		func() test {
			return test{
				name:   "return Dialer when all fields are empty",
				fields: fields{},
				want: want{
					want: &Dialer{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			d := &Dialer{
				Timeout:          test.fields.Timeout,
				Keepalive:        test.fields.Keepalive,
				FallbackDelay:    test.fields.FallbackDelay,
				DualStackEnabled: test.fields.DualStackEnabled,
			}

			got := d.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestSocketOption_Bind(t *testing.T) {
	type fields struct {
		ReusePort                bool
		ReuseAddr                bool
		TCPFastOpen              bool
		TCPNoDelay               bool
		TCPCork                  bool
		TCPQuickAck              bool
		TCPDeferAccept           bool
		IPTransparent            bool
		IPRecoverDestinationAddr bool
	}
	type want struct {
		want *SocketOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *SocketOption) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *SocketOption) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			reusePort := true
			reuseAddr := true
			tcpFastOpen := true
			tcpNoDelay := false
			tcpCork := false
			tcpQuickAck := true
			tcpDefferAccept := true
			ipTransparent := true
			ipRecoverDestinationAddr := false
			return test{
				name: "return SocketOption when all parameters are set",
				fields: fields{
					ReusePort:                reusePort,
					ReuseAddr:                reuseAddr,
					TCPFastOpen:              tcpFastOpen,
					TCPNoDelay:               tcpNoDelay,
					TCPCork:                  tcpCork,
					TCPQuickAck:              tcpQuickAck,
					TCPDeferAccept:           tcpDefferAccept,
					IPTransparent:            ipTransparent,
					IPRecoverDestinationAddr: ipRecoverDestinationAddr,
				},
				want: want{
					want: &SocketOption{
						ReusePort:                reusePort,
						ReuseAddr:                reuseAddr,
						TCPFastOpen:              tcpFastOpen,
						TCPNoDelay:               tcpNoDelay,
						TCPCork:                  tcpCork,
						TCPQuickAck:              tcpQuickAck,
						TCPDeferAccept:           tcpDefferAccept,
						IPTransparent:            ipTransparent,
						IPRecoverDestinationAddr: ipRecoverDestinationAddr,
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return SocketOption when all parameters are not set",
				fields: fields{},
				want: want{
					want: &SocketOption{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			s := &SocketOption{
				ReusePort:                test.fields.ReusePort,
				ReuseAddr:                test.fields.ReuseAddr,
				TCPFastOpen:              test.fields.TCPFastOpen,
				TCPNoDelay:               test.fields.TCPNoDelay,
				TCPCork:                  test.fields.TCPCork,
				TCPQuickAck:              test.fields.TCPQuickAck,
				TCPDeferAccept:           test.fields.TCPDeferAccept,
				IPTransparent:            test.fields.IPTransparent,
				IPRecoverDestinationAddr: test.fields.IPRecoverDestinationAddr,
			}

			got := s.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestSocketOption_ToSocketFlag(t *testing.T) {
	type fields struct {
		socketOpts *SocketOption
	}
	type want struct {
		want control.SocketFlag
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, control.SocketFlag) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got control.SocketFlag) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name:   "return flg when socketOpts is nil",
			fields: fields{},
			want: want{
				want: 0x00,
			},
		},
		{
			name: "return flg when socketOpts.ReuseAddr true and others are false",
			fields: fields{
				socketOpts: &SocketOption{
					ReuseAddr: true,
				},
			},
			want: want{
				want: 0x2,
			},
		},
		{
			name: "return flg when socketOpts.ReusePort true and others are false",
			fields: fields{
				socketOpts: &SocketOption{
					ReusePort: true,
				},
			},
			want: want{
				want: 0x1,
			},
		},
		{
			name: "return flg when socketOpts.TCPFastOpen true and others are false",
			fields: fields{
				socketOpts: &SocketOption{
					TCPFastOpen: true,
				},
			},
			want: want{
				want: 0x4,
			},
		},
		{
			name: "return flg when socketOpts.TCPCork true and others are false",
			fields: fields{
				socketOpts: &SocketOption{
					TCPCork: true,
				},
			},
			want: want{
				want: 0x10,
			},
		},
		{
			name: "return flg when socketOpts.TCPNoDelay true and others are false",
			fields: fields{
				socketOpts: &SocketOption{
					TCPNoDelay: true,
				},
			},
			want: want{
				want: 0x8,
			},
		},
		{
			name: "return flg when socketOpts.TCPDeferAccept true and others are false",
			fields: fields{
				socketOpts: &SocketOption{
					TCPDeferAccept: true,
				},
			},
			want: want{
				want: 0x40,
			},
		},
		{
			name: "return flg when socketOpts.TCPQuickAck true and others are false",
			fields: fields{
				socketOpts: &SocketOption{
					TCPQuickAck: true,
				},
			},
			want: want{
				want: 0x20,
			},
		},
		{
			name: "return flg when socketOpts.IPTransparent true and others are false",
			fields: fields{
				socketOpts: &SocketOption{
					IPTransparent: true,
				},
			},
			want: want{
				want: 0x80,
			},
		},
		{
			name: "return flg when socketOpts.IPRecoverDestinationAddr true and others are false",
			fields: fields{
				socketOpts: &SocketOption{
					IPRecoverDestinationAddr: true,
				},
			},
			want: want{
				want: 0x100,
			},
		},
		{
			name: "return flg when all fields of socketOpts are true",
			fields: fields{
				socketOpts: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               true,
					TCPCork:                  true,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            true,
					IPRecoverDestinationAddr: true,
				},
			},
			want: want{
				want: 0x1ff,
			},
		},
		{
			name: "return flg when all fields of socketOpts are false",
			fields: fields{
				socketOpts: &SocketOption{},
			},
			want: want{
				want: 0x0,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			s := test.fields.socketOpts

			got := s.ToSocketFlag()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNet_Bind(t *testing.T) {
	type fields struct {
		DNS          *DNS
		Dialer       *Dialer
		SocketOption *SocketOption
		TLS          *TLS
	}
	type want struct {
		want *Net
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Net) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Net) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			dns := &DNS{
				CacheEnabled:    true,
				RefreshDuration: "10m",
				CacheExpiration: "24h",
			}
			dialer := &Dialer{
				Timeout:          "3s",
				Keepalive:        "5m",
				FallbackDelay:    "10m",
				DualStackEnabled: false,
			}
			socketOption := &SocketOption{
				ReusePort:                true,
				ReuseAddr:                true,
				TCPFastOpen:              true,
				TCPNoDelay:               false,
				TCPCork:                  false,
				TCPQuickAck:              true,
				TCPDeferAccept:           true,
				IPTransparent:            true,
				IPRecoverDestinationAddr: false,
			}
			tls := &TLS{
				Enabled: false,
			}
			return test{
				name: "return Net when all fields are set",
				fields: fields{
					DNS:          dns,
					Dialer:       dialer,
					SocketOption: socketOption,
					TLS:          tls,
				},
				want: want{
					want: &Net{
						DNS:          dns,
						Dialer:       dialer,
						SocketOption: socketOption,
						TLS:          tls,
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return Net when all fields are empty",
				fields: fields{},
				want: want{
					want: &Net{},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			tr := &Net{
				DNS:          test.fields.DNS,
				Dialer:       test.fields.Dialer,
				SocketOption: test.fields.SocketOption,
				TLS:          test.fields.TLS,
			}

			got := tr.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNet_Opts(t *testing.T) {
	type fields struct {
		DNS          *DNS
		Dialer       *Dialer
		SocketOption *SocketOption
		TLS          *TLS
	}
	type want struct {
		want []net.DialerOption
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []net.DialerOption, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotOpts []net.DialerOption, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(len(gotOpts), len(w.want)) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOpts, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return 8 net.DialerOption and nil error when all fields are not empty and TLS is not enabled",
			fields: fields{
				DNS: &DNS{
					CacheEnabled:    true,
					RefreshDuration: "10m",
					CacheExpiration: "24h",
				},
				Dialer: &Dialer{
					Timeout:          "3s",
					Keepalive:        "5m",
					FallbackDelay:    "10m",
					DualStackEnabled: true,
				},
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               false,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            true,
					IPRecoverDestinationAddr: false,
				},
				TLS: &TLS{
					Enabled: false,
				},
			},
			want: want{
				want: make([]net.DialerOption, 8),
			},
		},
		{
			name: "return 6 net.DialerOption and nil error when all fields are not empty and TLS/DNS Cache/Dialer DualStack is not enabled",
			fields: fields{
				DNS: &DNS{
					CacheEnabled:    false,
					RefreshDuration: "10m",
					CacheExpiration: "24h",
				},
				Dialer: &Dialer{
					Timeout:          "3s",
					Keepalive:        "5m",
					FallbackDelay:    "10m",
					DualStackEnabled: false,
				},
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               false,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            true,
					IPRecoverDestinationAddr: false,
				},
				TLS: &TLS{
					Enabled: false,
				},
			},
			want: want{
				want: make([]net.DialerOption, 6),
			},
		},
		{
			name: "return nil net.DialerOption an error when all fields are not empty and tls.New() returns error",
			fields: fields{
				DNS: &DNS{
					CacheEnabled:    true,
					RefreshDuration: "10m",
					CacheExpiration: "24h",
				},
				Dialer: &Dialer{
					Timeout:          "3s",
					Keepalive:        "5m",
					FallbackDelay:    "10m",
					DualStackEnabled: true,
				},
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               false,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            true,
					IPRecoverDestinationAddr: false,
				},
				TLS: &TLS{
					Enabled: true,
				},
			},
			want: want{
				want: make([]net.DialerOption, 0),
				err:  errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			name: "return 8 net.DialerOption and nil error when all fields are not empty and tls.New() returns nil error",
			fields: fields{
				DNS: &DNS{
					CacheEnabled:    true,
					RefreshDuration: "10m",
					CacheExpiration: "24h",
				},
				Dialer: &Dialer{
					Timeout:          "3s",
					Keepalive:        "5m",
					FallbackDelay:    "10m",
					DualStackEnabled: true,
				},
				SocketOption: &SocketOption{
					ReusePort:                true,
					ReuseAddr:                true,
					TCPFastOpen:              true,
					TCPNoDelay:               false,
					TCPCork:                  false,
					TCPQuickAck:              true,
					TCPDeferAccept:           true,
					IPTransparent:            true,
					IPRecoverDestinationAddr: false,
				},
				TLS: &TLS{
					Enabled: true,
					Cert:    testdata.GetTestdataPath("tls/dummyServer.crt"),
					Key:     testdata.GetTestdataPath("tls/dummyServer.key"),
					CA:      testdata.GetTestdataPath("tls/dummyCa.pem"),
				},
			},
			want: want{
				want: make([]net.DialerOption, 9),
			},
		},
		{
			name:   "return 0 net.DialerOption and nil error when all fields are empty",
			fields: fields{},
			want: want{
				want: make([]net.DialerOption, 0),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			tr := &Net{
				DNS:          test.fields.DNS,
				Dialer:       test.fields.Dialer,
				SocketOption: test.fields.SocketOption,
				TLS:          test.fields.TLS,
			}

			got, err := tr.Opts()
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
