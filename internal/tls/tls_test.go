//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package tls provides implementation of Go API for tls certificate provider
package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	testdata "github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/test/goleak"
)

var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want *Config
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *Config, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *Config, err error) error {
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
			name: "returns cfg and nil when option is not empty",
			args: args{
				opts: []Option{
					WithCert(testdata.GetTestdataPath("tls/dummyServer.crt")),
					WithKey(testdata.GetTestdataPath("tls/dummyServer.key")),
					WithCa(testdata.GetTestdataPath("tls/dummyCa.pem")),
				},
			},
			want: want{
				want: func() *Config {
					cfg := new(tls.Config)

					cfg.Certificates = make([]tls.Certificate, 1)
					cfg.Certificates[0], _ = tls.LoadX509KeyPair(testdata.GetTestdataPath("tls/dummyServer.crt"),
						testdata.GetTestdataPath("tls/dummyServer.key"))

					pool, _ := NewX509CertPool(testdata.GetTestdataPath("tls/dummyCa.pem"))
					cfg.ClientCAs = pool
					cfg.ClientAuth = tls.RequireAndVerifyClientCert

					return cfg
				}(),
			},
			checkFunc: func(w want, c *tls.Config, err error) error {
				if !errors.Is(err, w.err) {
					return fmt.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}

				if len(c.Certificates) != 1 && len(c.Certificates) != len(w.want.Certificates) {
					return errors.New("Certificates length is wrong")
				}

				want := string(w.want.Certificates[0].Certificate[0])
				got := string(c.Certificates[0].Certificate[0])
				if want != got {
					return errors.Errorf("Certificates[0] want: %v, but got: %v", want, got)
				}

				if ok := c.ClientCAs.Equal(w.want.ClientCAs); !ok {
					return errors.Errorf("ClientCAs.Equal want: %v, got: %v", want, got)
				}

				if got, want := c.ClientAuth, w.want.ClientAuth; want != got {
					return errors.Errorf("ClientAuth want: %v, but got: %v", want, got)
				}
				return nil
			},
		},
		{
			name: "returns nil and error when option is empty",
			args: args{},
			want: want{
				err: errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			name: "returns nil and error when cert path is empty",
			args: args{
				opts: []Option{
					WithKey(testdata.GetTestdataPath("tls/dummyServer.key")),
				},
			},
			want: want{
				err: errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			name: "returns nil and error when key path is empty",
			args: args{
				opts: []Option{
					WithCert(testdata.GetTestdataPath("tls/dummyServer.crt")),
				},
			},
			want: want{
				err: errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			name: "returns nil and error when contents of cert file is invalid",
			args: args{
				opts: []Option{
					WithCert(testdata.GetTestdataPath("tls/invalid.crt")),
					WithKey(testdata.GetTestdataPath("tls/dummyServer.key")),
				},
			},
			want: want{
				err: errors.New("tls: failed to find any PEM data in certificate input"),
			},
		},
		{
			name: "returns nil and error when contents of ca file is invalid",
			args: args{
				opts: []Option{
					WithCert(testdata.GetTestdataPath("tls/dummyServer.crt")),
					WithKey(testdata.GetTestdataPath("tls/dummyServer.key")),
					WithCa(testdata.GetTestdataPath("tls/invalid.pem")),
				},
			},
			want: want{
				err: errors.ErrCertificationFailed,
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

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewClientConfig(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want *Config
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *Config, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *Config, err error) error {
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
			name: "returns cfg and nil when option is empty",
			checkFunc: func(w want, c *Config, err error) error {
				if !errors.Is(err, w.err) {
					return fmt.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}
				if c == nil {
					return errors.New("config is nil")
				}
				return nil
			},
		},
		{
			name: "returns cfg and nil when cert and key option is not empty",
			args: args{
				opts: []Option{
					WithCert(testdata.GetTestdataPath("tls/dummyServer.crt")),
					WithKey(testdata.GetTestdataPath("tls/dummyServer.key")),
				},
			},
			checkFunc: func(w want, c *Config, err error) error {
				if !errors.Is(err, w.err) {
					return fmt.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}
				if c == nil {
					return errors.New("config is nil")
				}
				if len(c.Certificates) != 1 {
					return errors.Errorf("invalid certificate was set. %v", c.Certificates)
				}
				return nil
			},
		},
		{
			name: "returns nil and error when contents of ca file is invalid",
			args: args{
				opts: []Option{
					WithCa(testdata.GetTestdataPath("tls/invalid.pem")),
				},
			},
			want: want{
				err: errors.ErrCertificationFailed,
			},
		},
		{
			name: "returns nil and error when contents of cert file is invalid",
			args: args{
				opts: []Option{
					WithCert(testdata.GetTestdataPath("tls/invalid.crt")),
					WithKey(testdata.GetTestdataPath("tls/dummyServer.key")),
				},
			},
			checkFunc: func(w want, c *Config, err error) error {
				wantErr := "tls: failed to find any PEM data in certificate input"
				if err.Error() != wantErr {
					return fmt.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}
				if c != nil {
					return errors.Errorf("config is not nil: %v", c)
				}
				return nil
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

			got, err := NewClientConfig(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestNewX509CertPool(t *testing.T) {
	type args struct {
		path string
	}
	type want struct {
		want *x509.CertPool
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *x509.CertPool, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *x509.CertPool, err error) error {
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
			name: "returns pool and nil when the pool exists and adds the cert into pool",
			args: args{
				path: testdata.GetTestdataPath("tls/dummyServer.crt"),
			},
			want: want{
				want: func() *x509.CertPool {
					path := testdata.GetTestdataPath("tls/dummyServer.crt")
					pool, err := x509.SystemCertPool()
					if err != nil {
						pool = x509.NewCertPool()
					}
					b, err := file.ReadFile(path)
					if err == nil {
						pool.AppendCertsFromPEM(b)
					}
					return pool
				}(),
			},
			checkFunc: func(w want, cp *x509.CertPool, err error) error {
				if err != nil {
					return errors.Errorf("err is not nil. err: %v", err)
				}
				if cp == nil {
					return errors.New("got is nil")
				}
				if ok := cp.Equal(w.want); !ok {
					return errors.Errorf("not equals. want: %#v, got: %#v", w.want, cp)
				}
				return nil
			},
		},
		{
			name: "returns nil and error when contents of path is invalid",
			args: args{
				path: testdata.GetTestdataPath("tls/invalid.pem"),
			},
			want: want{
				err: errors.ErrCertificationFailed,
			},
			checkFunc: func(w want, cp *x509.CertPool, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("err not equals. want: %v, but got: %v", w.err, err)
				}
				if cp == nil {
					return errors.Errorf("got is nil: %v", cp)
				}
				return nil
			},
		},
		{
			name: "returns nil and error when path dose not exist",
			args: args{
				path: "not_exist",
			},
			checkFunc: func(w want, cp *x509.CertPool, err error) error {
				if err == nil {
					return errors.New("err is nil")
				}
				if cp != nil {
					return errors.Errorf("got is not nil: %v, want: %v", cp, w)
				}
				return nil
			},
		},
		{
			name: "returns nil and error when path is empty",
			checkFunc: func(w want, cp *x509.CertPool, err error) error {
				if err == nil {
					return errors.New("err is nil")
				}
				if cp != nil {
					return errors.Errorf("got is not nil: %v, want: %v", cp, w)
				}
				return nil
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

			got, err := NewX509CertPool(test.args.path)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
