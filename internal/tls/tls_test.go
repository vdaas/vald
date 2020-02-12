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
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

var (
	dummyCaPath          = "./testdata/dummyCa.pem"
	dummyCertPath        = "./testdata/dummyServer.crt"
	dummyKeyPath         = "./testdata/dummyServer.key"
	dummyInvalidCaPath   = "./testdata/invalid.pem"
	dummyInvalidCertPath = "./testdata/invalid.crt"
)

func TestNew(t *testing.T) {
	type test struct {
		name      string
		opts      []Option
		checkFunc func(*Config, error) error
	}

	tests := []test{
		func() test {
			wantFn := func() (*Config, error) {
				var err error
				c := &credentials{
					cfg: new(tls.Config),
				}

				c.cfg.Certificates = make([]tls.Certificate, 1)
				c.cfg.Certificates[0], err = tls.LoadX509KeyPair(dummyCertPath, dummyKeyPath)
				if err != nil {
					return nil, err
				}

				b, err := ioutil.ReadFile(dummyCertPath)
				if err != nil {
					return nil, err
				}

				pool := x509.NewCertPool()
				if !pool.AppendCertsFromPEM(b) {
					return nil, errors.New("faild to add cert")
				}

				c.cfg.ClientCAs = pool
				c.cfg.ClientAuth = tls.RequireAndVerifyClientCert

				c.cfg.BuildNameToCertificate()

				return c.cfg, nil
			}

			return test{
				name: "returns cfg and nil when option is not empty",
				opts: []Option{
					WithCert(dummyCertPath),
					WithKey(dummyKeyPath),
					WithCa(dummyCertPath),
				},
				checkFunc: func(cfg *tls.Config, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil: %v", err)
					}

					if cfg == nil {
						return errors.New("cfg is nil")
					}

					wantCfg, wantErr := wantFn()
					if wantErr != nil {
						return errors.Errorf("wantErr is not nil: %v", wantErr)
					}

					if len(cfg.Certificates) != 1 && len(cfg.Certificates) != len(wantCfg.Certificates) {
						return errors.New("Certificates length is wrong")
					}

					if got, want := string(wantCfg.Certificates[0].Certificate[0]), string(cfg.Certificates[0].Certificate[0]); want != got {
						return errors.Errorf("Certificates[0] want: %v, but got: %v", want, got)
					}

					if len(cfg.ClientCAs.Subjects()) == 0 {
						return errors.New("subjects are empty")
					}

					l := len(cfg.ClientCAs.Subjects()) - 1
					if got, want := cfg.ClientCAs.Subjects()[l], wantCfg.ClientCAs.Subjects()[0]; !reflect.DeepEqual(got, want) {
						return errors.Errorf("ClientCAs.Subjects want: %v, got: %v", want, got)
					}

					if got, want := cfg.ClientAuth, wantCfg.ClientAuth; want != got {
						return errors.Errorf("ClientAuth want: %v, but got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		{
			name: "returns nil and error when option is empty",
			checkFunc: func(cfg *tls.Config, err error) error {
				if err == nil {
					return errors.New("err is nil")
				} else if !errors.Is(err, errors.ErrTLSCertOrKeyNotFound) {
					return errors.Errorf("want err: %v, got: %v", errors.ErrTLSCertOrKeyNotFound, err)
				}

				if cfg != nil {
					return errors.Errorf("cfg is not nil: %v", cfg)
				}
				return nil
			},
		},

		{
			name: "returns nil and error when cert path is empty",
			opts: []Option{
				WithKey(dummyKeyPath),
			},
			checkFunc: func(cfg *tls.Config, err error) error {
				if err == nil {
					return errors.New("err is nil")
				} else if !errors.Is(err, errors.ErrTLSCertOrKeyNotFound) {
					return errors.Errorf("want err: %v, got: %v", errors.ErrTLSCertOrKeyNotFound, err)
				}

				if cfg != nil {
					return errors.Errorf("cfg is not nil: %v", cfg)
				}
				return nil
			},
		},

		{
			name: "returns nil and error when key path is empty",
			opts: []Option{
				WithCert(dummyCertPath),
			},
			checkFunc: func(cfg *tls.Config, err error) error {
				if err == nil {
					return errors.New("err is nil")
				} else if !errors.Is(err, errors.ErrTLSCertOrKeyNotFound) {
					return errors.Errorf("want err: %v, got: %v", errors.ErrTLSCertOrKeyNotFound, err)
				}

				if cfg != nil {
					return errors.Errorf("cfg is not nil: %v", cfg)
				}
				return nil
			},
		},

		{
			name: "returns nil and error when contents of cert file is invalid",
			opts: []Option{
				WithCert(dummyInvalidCertPath),
				WithKey(dummyKeyPath),
			},
			checkFunc: func(cfg *tls.Config, err error) error {
				if err == nil {
					return errors.New("err is nil")
				} else if want, got := errors.New("tls: failed to find any PEM data in certificate input"), err; want.Error() != got.Error() {
					return errors.Errorf("want err: %v, but got: %v", want, got)
				}

				if cfg != nil {
					return errors.Errorf("cfg is not nil: %v", cfg)
				}
				return nil
			},
		},

		{
			name: "returns nil and error when contents of ca file is invalid",
			opts: []Option{
				WithCert(dummyCertPath),
				WithKey(dummyKeyPath),
				WithCa(dummyInvalidCaPath),
			},
			checkFunc: func(cfg *tls.Config, err error) error {
				if err == nil {
					return errors.New("err is nil")
				} else if !errors.Is(err, errors.ErrCertificationFailed) {
					return errors.Errorf("want err: %v, but got: %v", errors.ErrCertificationFailed, err)
				}

				if cfg != nil {
					return errors.Errorf("cfg is not nil: %v", cfg)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := New(tt.opts...)
			if err := tt.checkFunc(cfg, err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestNewClientConfig(t *testing.T) {
	tests := []struct {
		name      string
		opts      []Option
		checkFunc func(*Config, error) error
	}{
		{
			name: "returns cfg and nil when option is empty",
			checkFunc: func(cfg *Config, err error) error {
				if err != nil {
					return errors.Errorf("err is not nil. err: %v", err)
				}

				if cfg == nil {
					return errors.New("cfg is nil")
				}
				return nil
			},
		},

		{
			name: "returns cfg and nil when cert and key option is not empty",
			opts: []Option{
				WithCert(dummyCertPath),
				WithKey(dummyKeyPath),
			},
			checkFunc: func(cfg *Config, err error) error {
				if err != nil {
					return errors.Errorf("err is not nil. err: %v", err)
				}

				if cfg == nil {
					return errors.New("cfg is nil")
				}

				if len(cfg.Certificates) != 1 {
					return errors.Errorf("invalid certificate was set. %v", cfg.Certificates)
				}
				return nil
			},
		},

		{
			name: "returns cfg and nil when ca option is not empty",
			opts: []Option{
				WithCa(dummyCaPath),
			},
			checkFunc: func(cfg *Config, err error) error {
				if err != nil {
					return errors.Errorf("err is not nil. err: %v", err)
				}

				if cfg == nil {
					return errors.New("cfg is nil")
				}

				if cfg.RootCAs == nil {
					return errors.New("rootca is nil")
				}

				// TODO: added test case
				return nil
			},
		},

		{
			name: "returns nil and error when contents of ca file is invalid",
			opts: []Option{
				WithCa(dummyInvalidCaPath),
			},
			checkFunc: func(cfg *Config, err error) error {
				if err == nil {
					return errors.New("err is nil")
				} else if !errors.Is(err, errors.ErrCertificationFailed) {
					return errors.Errorf("want err: %v, but got: %v", errors.ErrCertificationFailed, err)
				}

				if cfg != nil {
					return errors.Errorf("cfg is not nil: %v", cfg)
				}

				return nil
			},
		},

		{
			name: "returns nil and error when contents of cert file is invalid",
			opts: []Option{
				WithCert(dummyInvalidCertPath),
				WithKey(dummyKeyPath),
			},
			checkFunc: func(cfg *Config, err error) error {
				if err == nil {
					return errors.New("err is nil")
				} else if want, got := errors.New("tls: failed to find any PEM data in certificate input"), err; want.Error() != got.Error() {
					return errors.Errorf("want err: %v, but got: %v", want, got)
				}

				if cfg != nil {
					return errors.Errorf("cfg is not nil: %v", cfg)
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewClientConfig(tt.opts...)
			if err := tt.checkFunc(cfg, err); err != nil {
				t.Error(err)
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
			wantFn := func() (*x509.CertPool, error) {
				pool := x509.NewCertPool()
				b, err := ioutil.ReadFile(dummyCertPath)
				if err != nil {
					return nil, err
				}

				if !pool.AppendCertsFromPEM(b) {
					return nil, errors.New("faild to add cert")
				}

				return pool, nil
			}

			return test{
				name: "returns pool and nil when the pool exists and adds the cert into pool",
				path: dummyCertPath,
				checkFunc: func(pool *x509.CertPool, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil. err: %v", err)
					}

					if pool == nil {
						return errors.New("got is nil")
					}

					want, err := wantFn()
					if err != nil {
						return errors.Errorf("faild to create want object. err:", err)
					}

					if len(pool.Subjects()) == 0 {
						return errors.New("cert files are empty")
					}

					l := len(pool.Subjects()) - 1
					if got, want := pool.Subjects()[l], want.Subjects()[0]; !reflect.DeepEqual(got, want) {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		{
			name: "returns nil and error when contents of path is invalid",
			path: dummyInvalidCertPath,
			checkFunc: func(pool *x509.CertPool, err error) error {
				if err == nil {
					return errors.New("err is ")
				} else if !errors.Is(err, errors.ErrCertificationFailed) {
					return errors.Errorf("err not equals. want: %v, but got: %v", errors.ErrCertificationFailed, err)
				}

				if pool == nil {
					return errors.Errorf("got is not nil: %v", pool)
				}

				return nil
			},
		},

		{
			name: "returns nil and error when path dose not exist",
			path: "not_exist",
			checkFunc: func(pool *x509.CertPool, err error) error {
				if err == nil {
					return errors.New("err is nil")
				}

				if pool != nil {
					return errors.Errorf("got is not nil: %v", pool)
				}
				return nil
			},
		},

		{
			name: "returns nil and error when path is empty",
			checkFunc: func(pool *x509.CertPool, err error) error {
				if err == nil {
					return errors.New("err is nil")
				}

				if pool != nil {
					return errors.Errorf("got is not nil: %v", pool)
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
