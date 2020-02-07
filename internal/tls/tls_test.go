//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	dummyCertPath := "./testdata/dummyServer.crt"
	dummyKeyPath := "./testdata/dummyServer.key"
	dummyCaPath := "./testdata/dummyCa.pem"
	defaultArgs := args{
		opts: []Option{
			WithCert(dummyCertPath),
			WithKey(dummyKeyPath),
			WithCa(dummyCaPath),
		},
	}
	tests := []struct {
		name      string
		args      args
		want      *Config
		checkFunc func(*Config, *Config) error
		wantErr   error
	}{
		{
			name: "return value MinVersion test.",
			args: defaultArgs,
			want: func() *Config {
				conf := defaultConfig
				cert, _ := tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
				conf.Certificates = []tls.Certificate{cert}
				return conf
			}(),
			checkFunc: func(got, want *tls.Config) error {
				if got.MinVersion != want.MinVersion {
					return fmt.Errorf("MinVersion not Matched :\tgot %d\twant %d", got.MinVersion, want.MinVersion)
				}
				return nil
			},
		},
		{
			name: "return value CurvePreferences test.",
			args: defaultArgs,
			want: func() *Config {
				conf := defaultConfig
				cert, _ := tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
				conf.Certificates = []tls.Certificate{cert}
				return conf
			}(),
			checkFunc: func(got, want *tls.Config) error {
				if len(got.CurvePreferences) != len(want.CurvePreferences) {
					return fmt.Errorf("CurvePreferences not Matched length:\tgot %d\twant %d", len(got.CurvePreferences), len(want.CurvePreferences))
				}
				for _, actualValue := range got.CurvePreferences {
					var match bool
					for _, expectedValue := range want.CurvePreferences {
						if actualValue == expectedValue {
							match = true
							break
						}
					}

					if !match {
						return fmt.Errorf("CurvePreferences not Find :\twant %s", string(want.MinVersion))
					}
				}
				return nil
			},
		},
		{
			name: "return value SessionTicketsDisabled test.",
			args: defaultArgs,
			want: func() *Config {
				conf := defaultConfig
				cert, _ := tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
				conf.Certificates = []tls.Certificate{cert}
				return conf
			}(),
			checkFunc: func(got, want *tls.Config) error {
				if got.SessionTicketsDisabled != want.SessionTicketsDisabled {
					return fmt.Errorf("SessionTicketsDisabled not matched :\tgot %v\twant %v", got.SessionTicketsDisabled, want.SessionTicketsDisabled)
				}
				return nil
			},
		},
		{
			name: "return value ClientAuth test.",
			args: defaultArgs,
			want: func() *Config {
				conf := defaultConfig
				cert, _ := tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
				conf.Certificates = []tls.Certificate{cert}
				return conf
			}(),
			checkFunc: func(got, want *tls.Config) error {
				if got.ClientAuth != want.ClientAuth {
					return fmt.Errorf("ClientAuth not Matched :\tgot %d \twant %d", got.ClientAuth, want.ClientAuth)
				}
				return nil
			},
		},
		{
			name: "cert file not found return error",
			args: func() args {
				a := defaultArgs
				a.opts = append(a.opts, WithCert(""))
				return a
			}(),
			wantErr: errors.ErrTLSCertOrKeyNotFound,
		},
		{
			name: "key file not found return error",
			args: func() args {
				a := defaultArgs
				a.opts = append(a.opts, WithCert(""))
				return a
			}(),
			wantErr: errors.ErrTLSCertOrKeyNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got, gotErr := New(tt.args.opts...)
			// if tt.checkFunc != nil {
			// 	if err := tt.checkFunc(got, tt.want); err != nil {
			// 		t.Errorf("NewTLSConfig() error = %v", err)
			// 		return
			// 	}
			// }
			// if gotErr != nil {
			// 	if tt.wantErr == nil {
			// 		t.Errorf("NewTLSConfig() error = %v, wantErr = %v", gotErr, tt.wantErr)
			// 	} else if gotErr.Error() != tt.wantErr.Error() {
			// 		t.Errorf("NewTLSConfig() error = %v, wantErr = %v", gotErr, tt.wantErr)
			// 	}
			// } else if tt.wantErr != nil {
			// 	t.Errorf("NewTLSConfig() error = %v, wantErr = %v", gotErr, tt.wantErr)
			// }
		})
	}
}

func TestNewClientConfig(t *testing.T) {
	type args struct {
		opts []Option
	}
	dummyCertPath := "./testdata/dummyServer.crt"
	dummyKeyPath := "./testdata/dummyServer.key"
	dummyCaPath := "./testdata/dummyCa.pem"
	defaultArgs := args{
		opts: []Option{
			WithCert(dummyCertPath),
			WithKey(dummyKeyPath),
			WithCa(dummyCaPath),
		},
	}
	tests := []struct {
		name      string
		args      args
		checkFunc func(*Config, *Config) error
		want      *Config
		wantErr   error
	}{
		{
			name: "return value MinVersion test.",
			args: defaultArgs,
			want: func() *Config {
				conf := defaultConfig
				cert, _ := tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
				conf.Certificates = []tls.Certificate{cert}
				return conf
			}(),
			checkFunc: func(got, want *tls.Config) error {
				if got.MinVersion != want.MinVersion {
					return fmt.Errorf("MinVersion not Matched :\tgot %d\twant %d", got.MinVersion, want.MinVersion)
				}
				return nil
			},
		},

		{
			name: "return value CurvePreferences test.",
			args: defaultArgs,
			want: func() *Config {
				conf := defaultConfig
				cert, _ := tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
				conf.Certificates = []tls.Certificate{cert}
				return conf
			}(),
			checkFunc: func(got, want *tls.Config) error {
				if len(got.CurvePreferences) != len(want.CurvePreferences) {
					return fmt.Errorf("CurvePreferences not Matched length:\tgot %d\twant %d", len(got.CurvePreferences), len(want.CurvePreferences))
				}
				for _, actualValue := range got.CurvePreferences {
					var match bool
					for _, expectedValue := range want.CurvePreferences {
						if actualValue == expectedValue {
							match = true
							break
						}
					}

					if !match {
						return fmt.Errorf("CurvePreferences not Find :\twant %s", string(want.MinVersion))
					}
				}
				return nil
			},
		},
		{
			name: "return value SessionTicketsDisabled test.",
			args: defaultArgs,
			want: func() *Config {
				conf := defaultConfig
				cert, _ := tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
				conf.Certificates = []tls.Certificate{cert}
				return conf
			}(),
			checkFunc: func(got, want *tls.Config) error {
				if got.SessionTicketsDisabled != want.SessionTicketsDisabled {
					return fmt.Errorf("SessionTicketsDisabled not matched :\tgot %v\twant %v", got.SessionTicketsDisabled, want.SessionTicketsDisabled)
				}
				return nil
			},
		},
		{
			name: "return value ClientAuth test.",
			args: defaultArgs,
			want: func() *Config {
				conf := defaultConfig
				cert, _ := tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
				conf.Certificates = []tls.Certificate{cert}
				return conf
			}(),
			checkFunc: func(got, want *tls.Config) error {
				if got.ClientAuth != want.ClientAuth {
					return fmt.Errorf("ClientAuth not Matched :\tgot %d \twant %d", got.ClientAuth, want.ClientAuth)
				}
				return nil
			},
		},
		{
			name: "cert file not found return error",
			args: args{
				[]Option{
					WithCert(""),
					WithKey(dummyKeyPath),
					WithCa(dummyCaPath),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewClientConfig(tt.args.opts...)
			if tt.checkFunc != nil {
				if err := tt.checkFunc(got, tt.want); err != nil {
					t.Errorf("NewTLSConfig() error = %v", err)
					return
				}
			}
			if gotErr != nil {
				if tt.wantErr == nil {
					t.Errorf("NewTLSConfig() error = %v, wantErr = %v", gotErr, tt.wantErr)
				} else if gotErr.Error() != tt.wantErr.Error() {
					t.Errorf("NewTLSConfig() error = %v, wantErr = %v", gotErr, tt.wantErr)
				}
			} else if tt.wantErr != nil {
				t.Errorf("NewTLSConfig() error = %v, wantErr = %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestNewX509CertPool(t *testing.T) {
	type test struct {
		name      string
		path      string
		checkFunc func(*x509.CertPool, error) error
	}

	tests := []test{
		func() test {
			path := "./testdata/dummyServer.crt"

			wantFn := func() (got *x509.CertPool, err error) {
				pool := x509.NewCertPool()
				b, err := ioutil.ReadFile("./testdata/dummyServer.crt")
				if err != nil {
					return nil, err
				}

				if !pool.AppendCertsFromPEM(b) {
					return nil, errors.New("faild to add cert")
				}

				return pool, nil
			}

			return test{
				name: "returns pool and nil when the pool exists and adds the cert file into pool",
				path: path,
				checkFunc: func(got *x509.CertPool, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil. err: %v", err)
					}

					if got == nil {
						return errors.New("got is nil")
					}

					want, err := wantFn()
					if err != nil {
						return errors.Errorf("faild to create want object. err:", err)
					}

					if len(got.Subjects()) == 0 {
						return errors.New("cert files are empty")
					}

					if got, want := got.Subjects()[len(got.Subjects())-1], want.Subjects()[0]; !reflect.DeepEqual(got, want) {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		{
			name: "returns nil and error when contents of path is invalid",
			path: "./testdata/invalid.crt",
			checkFunc: func(got *x509.CertPool, err error) error {
				if err == nil {
					return errors.New("err is ")
				} else if !errors.Is(err, errors.ErrCertificationFailed) {
					return errors.Errorf("err not equals. want: %v, but got: %v", errors.ErrCertificationFailed, err)
				}

				if got == nil {
					return errors.Errorf("got is not nil: %v", got)
				}

				return nil
			},
		},

		{
			name: "returns nil and error when path dose not exist",
			path: "not_exist",
			checkFunc: func(got *x509.CertPool, err error) error {
				if err == nil {
					return errors.New("err is nil")
				}

				if got != nil {
					return errors.Errorf("got is not nil: %v", got)
				}
				return nil
			},
		},

		{
			name: "returns nil and error when path is empty",
			checkFunc: func(got *x509.CertPool, err error) error {
				if err == nil {
					return errors.New("err is nil")
				}

				if got != nil {
					return errors.Errorf("got is not nil: %v", got)
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewX509CertPool(tt.path)
			if err = tt.checkFunc(got, err); err != nil {
				t.Error(err)
			}
		})
	}
}
