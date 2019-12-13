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
	"strings"
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
		/*
			{
				name: "return value Certificates test.",
				args: defaultArgs,
				want: func() *Config {
					conf := defaultConfig
					cert, _ := tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
					conf.Certificates = []tls.Certificate{cert}
					return conf
				}(),
				checkFunc: func(got, want *tls.Config) error {
					for _, wantVal := range want.Certificates {
						var notExist = false
						for _, gotVal := range got.Certificates {
							if gotVal.PrivateKey == wantVal.PrivateKey {
								notExist = true
								break
							}
						}
						if notExist {
							return fmt.Errorf("Certificates PrivateKey not Matched :\twant %s", wantVal.PrivateKey)
						}
					}
					return nil
				},
			},
		*/
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
			got, gotErr := New(tt.args.opts...)
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
	type args struct {
		path string
	}

	tests := []struct {
		name      string
		args      args
		want      *x509.CertPool
		checkFunc func(*x509.CertPool, *x509.CertPool) error
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "Check err if file is not exists",
			args: args{
				path: "",
			},
			want: &x509.CertPool{},
			checkFunc: func(*x509.CertPool, *x509.CertPool) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "Check Append CA is correct",
			args: args{
				path: "./testdata/dummyCa.pem",
			},
			want: func() *x509.CertPool {
				wantPool := x509.NewCertPool()
				c, err := ioutil.ReadFile("./testdata/dummyCa.pem")
				if err != nil {
					panic(err)
				}
				if !wantPool.AppendCertsFromPEM(c) {
					panic(errors.New("Error appending certs from PEM"))
				}
				return wantPool
			}(),
			checkFunc: func(want *x509.CertPool, got *x509.CertPool) error {
				for _, wantCert := range want.Subjects() {
					exists := false
					for _, gotCert := range got.Subjects() {
						if strings.EqualFold(string(wantCert), string(gotCert)) {
							exists = true
						}
					}
					if !exists {
						return fmt.Errorf("Error\twant\t%s\t not found", string(wantCert))
					}
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewX509CertPool(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewX509CertPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc != nil {
				err = tt.checkFunc(tt.want, got)
				if err != nil {
					t.Errorf("TestNewX509CertPool error = %s", err)
				}
			}
		})
	}
}
