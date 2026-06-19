//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	testdata "github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/test/goleak"
)

var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestNewServerConfig(t *testing.T) {
	if err := testdata.Run(t.Context(), t, func(tt *testing.T, opts []Option) (*Config, error) {
		tt.Helper()
		return NewServerConfig(opts...)
	}, []testdata.Case[*Config, []Option]{
		{
			Name: "returns cfg and nil when option is not empty",
			Want: testdata.Result[*Config]{
				Val: func() *Config {
					cfg := new(tls.Config)
					cfg.Certificates = make([]tls.Certificate, 1)
					cfg.Certificates[0], _ = tls.LoadX509KeyPair(testdata.GetTestdataPath("tls/server.crt"),
						testdata.GetTestdataPath("tls/server.key"))
					pool, _ := NewX509CertPool(testdata.GetTestdataPath("tls/ca.pem"))
					cfg.ClientCAs = pool
					cfg.ClientAuth = tls.RequireAndVerifyClientCert
					return cfg
				}(),
				Err: nil,
			},
			Args: []Option{
				WithCert(testdata.GetTestdataPath("tls/server.crt")),
				WithKey(testdata.GetTestdataPath("tls/server.key")),
				WithCa(testdata.GetTestdataPath("tls/ca.pem")),
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[*Config], got testdata.Result[*Config]) error {
				tt.Helper()
				if !errors.Is(got.Err, want.Err) {
					return fmt.Errorf("got_error: \"%wv\",\n\t\t\t\twawt: \"%#v\"", got.Err, want.Err)
				}
				if len(got.Val.Certificates) != 1 && len(got.Val.Certificates) != len(want.Val.Certificates) {
					return errors.New("Certificates length is wrong")
				}
				wb := conv.Btoa(want.Val.Certificates[0].Certificate[0])
				gb := conv.Btoa(got.Val.Certificates[0].Certificate[0])
				if wb != gb {
					return errors.Errorf("Certificates[0] want: %v, but got: %v", wb, gb)
				}
				if ok := got.Val.ClientCAs.Equal(want.Val.ClientCAs); !ok {
					return errors.Errorf("ClientCAs.Equal want: %v, got: %v", wb, gb)
				}
				if want.Val.ClientAuth != got.Val.ClientAuth {
					return errors.Errorf("ClientAuth want: %v, but got: %v", wb, gb)
				}
				return nil
			},
		},
		{
			Name: "returns nil and error when option is empty",
			Want: testdata.Result[*Config]{
				Err: errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			Name: "returns nil and error when cert path is empty",
			Args: []Option{
				WithKey(testdata.GetTestdataPath("tls/server.key")),
			},
			Want: testdata.Result[*Config]{
				Err: errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			Name: "returns nil and error when key path is empty",
			Args: []Option{
				WithCert(testdata.GetTestdataPath("tls/server.crt")),
			},
			Want: testdata.Result[*Config]{
				Err: errors.ErrTLSCertOrKeyNotFound,
			},
		},
		{
			Name: "returns nil and error when contents of cert file is invalid",
			Args: []Option{
				WithCert(testdata.GetTestdataPath("tls/invalid-server.crt")),
				WithKey(testdata.GetTestdataPath("tls/server.key")),
			},
			Want: testdata.Result[*Config]{
				Err: errors.New("tls: failed to find \"CERTIFICATE\" PEM block in certificate input after skipping PEM blocks of the following types: [CERTIFICATE REQUEST]"),
			},
		},
		{
			Name: "returns nil and error when contents of ca file is invalid",
			Args: []Option{
				WithCert(testdata.GetTestdataPath("tls/server.crt")),
				WithKey(testdata.GetTestdataPath("tls/server.key")),
				WithCa(testdata.GetTestdataPath("tls/invalid-ca.pem")),
			},
			Want: testdata.Result[*Config]{
				Err: errors.ErrNoCertsAddedToPool,
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestNewClientConfig(t *testing.T) {
	if err := testdata.Run(t.Context(), t, func(tt *testing.T, opts []Option) (*Config, error) {
		tt.Helper()
		return NewClientConfig(opts...)
	}, []testdata.Case[*Config, []Option]{
		{
			Name: "returns cfg and nil when option is empty",
			CheckFunc: func(tt *testing.T, want testdata.Result[*Config], got testdata.Result[*Config]) error {
				tt.Helper()
				if got.Err != nil {
					return errors.New("")
				}
				if got.Val == nil {
					return errors.New("config is nil")
				}
				return nil
			},
		},
		{
			Name: "returns cfg and nil when cert and key option is not empty",
			Args: []Option{
				WithCert(testdata.GetTestdataPath("tls/server.crt")),
				WithKey(testdata.GetTestdataPath("tls/server.key")),
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[*Config], got testdata.Result[*Config]) error {
				tt.Helper()
				if !errors.Is(got.Err, want.Err) {
					return errors.Errorf("got_error: \"%s\",\n\t\t\t\twant: \"%s\"", got.Err.Error(), want.Err.Error())
				}
				if got.Val == nil {
					return errors.New("config is nil")
				}
				if len(got.Val.Certificates) != 1 {
					return errors.Errorf("invalid certificate was set. %v", got.Val.Certificates)
				}
				return nil
			},
		},
		{
			Name: "returns nil and error when contents of ca file is invalid",
			Args: []Option{
				WithCa(testdata.GetTestdataPath("tls/invalid-ca.pem")),
			},
			Want: testdata.Result[*Config]{
				Err: errors.ErrNoCertsAddedToPool,
			},
		},
		{
			Name: "returns nil and error when contents of cert file is invalid",
			Args: []Option{
				WithCert(testdata.GetTestdataPath("tls/invalid-server.crt")),
				WithKey(testdata.GetTestdataPath("tls/server.key")),
			},
			Want: testdata.Result[*Config]{
				Err: errors.New("tls: failed to find \"CERTIFICATE\" PEM block in certificate input after skipping PEM blocks of the following types: [CERTIFICATE REQUEST]"),
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[*Config], got testdata.Result[*Config]) error {
				tt.Helper()
				if !errors.Is(got.Err, want.Err) {
					return errors.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%v\"", got.Err, want.Err)
				}
				if got.Val != nil {
					return errors.Errorf("config is not nil: %v", got.Val)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestNewX509CertPool(t *testing.T) {
	if err := testdata.Run(t.Context(), t, func(tt *testing.T, path string) (*x509.CertPool, error) {
		tt.Helper()
		return NewX509CertPool(path)
	}, []testdata.Case[*x509.CertPool, string]{
		{
			Name: "returns pool and nil when the pool exists and adds the cert into pool",
			Args: testdata.GetTestdataPath("tls/ca.pem"),
			Want: testdata.Result[*x509.CertPool]{
				Val: func() *x509.CertPool {
					path := testdata.GetTestdataPath("tls/ca.pem")
					pool, err := x509.SystemCertPool()
					if err != nil {
						pool = x509.NewCertPool()
					}
					b, err := file.ReadFile(path)
					if err == nil && b != nil {
						pool.AppendCertsFromPEM(b)
					}
					return pool
				}(),
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[*x509.CertPool], got testdata.Result[*x509.CertPool]) error {
				if got.Err != nil {
					return errors.Errorf("err is not nil. err: %v", got.Err)
				}
				if got.Val == nil {
					return errors.New("got is nil")
				}
				if ok := got.Val.Equal(want.Val); !ok {
					return errors.New("cert pool is not equals")
				}
				return nil
			},
		},
		{
			Name: "returns nil and error when contents of path is invalid",
			Args: testdata.GetTestdataPath("tls/invalid-ca.pem"),
			Want: testdata.Result[*x509.CertPool]{
				Err: errors.ErrNoCertsAddedToPool,
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[*x509.CertPool], got testdata.Result[*x509.CertPool]) error {
				if !errors.Is(got.Err, want.Err) {
					return errors.Errorf("err not equals. want: %v, but got: %v", want.Err, got.Err)
				}
				if got.Val != nil {
					return errors.Errorf("got is not nil: %v", got.Val)
				}

				return nil
			},
		},
		{
			Name: "returns nil and error when path dose not exist",
			Args: "not_exist",
			CheckFunc: func(tt *testing.T, want testdata.Result[*x509.CertPool], got testdata.Result[*x509.CertPool]) error {
				if got.Err == nil {
					return errors.New("err is nil")
				}
				if got.Val != nil {
					return errors.Errorf("got is not nil: %v, want: %v", got, want)
				}
				return nil
			},
		},
		{
			Name: "returns nil and error when path is empty",
			CheckFunc: func(tt *testing.T, want testdata.Result[*x509.CertPool], got testdata.Result[*x509.CertPool]) error {
				if got.Err == nil {
					return errors.New("err is nil")
				}
				if got.Val != nil {
					return errors.Errorf("got is not nil: %v, want: %v", got, want)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_loadCRL(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	type want struct {
// 		want  map[string]struct{}
// 		want1 time.Time
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, map[string]struct{}, time.Time, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]struct{}, got1 time.Time, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		if !reflect.DeepEqual(got1, w.want1) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
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
// 		           path:"",
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
// 			got, got1, err := loadCRL(test.args.path)
// 			if err := checkFunc(test.want, got, got1, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_newCredential(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want *credentials
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *credentials, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *credentials, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           opts:nil,
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
// 		           opts:nil,
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
// 			got, err := newCredential(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_loadKeyPair(t *testing.T) {
// 	type args struct {
// 		role     string
// 		certPath string
// 		keyPath  string
// 	}
// 	type want struct {
// 		want tls.Certificate
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, tls.Certificate, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got tls.Certificate, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           role:"",
// 		           certPath:"",
// 		           keyPath:"",
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
// 		           role:"",
// 		           certPath:"",
// 		           keyPath:"",
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
// 			got, err := loadKeyPair(test.args.role, test.args.certPath, test.args.keyPath)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_credentials_reloadCert(t *testing.T) {
// 	type fields struct {
// 		crlNextUpdate time.Time
// 		certPtr       atomic.Pointer[tls.Certificate]
// 		revoked       map[string]struct{}
// 		cfg           *Config
// 		crl           string
// 		sn            string
// 		ca            string
// 		key           string
// 		cert          string
// 		clientAuth    tls.ClientAuthType
// 		insecure      bool
// 		hotReload     bool
// 	}
// 	type want struct {
// 		want *tls.Certificate
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *tls.Certificate, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *tls.Certificate, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		       fields: fields {
// 		           crlNextUpdate:time.Time{},
// 		           certPtr:nil,
// 		           revoked:nil,
// 		           cfg:Config{},
// 		           crl:"",
// 		           sn:"",
// 		           ca:"",
// 		           key:"",
// 		           cert:"",
// 		           clientAuth:nil,
// 		           insecure:false,
// 		           hotReload:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           crlNextUpdate:time.Time{},
// 		           certPtr:nil,
// 		           revoked:nil,
// 		           cfg:Config{},
// 		           crl:"",
// 		           sn:"",
// 		           ca:"",
// 		           key:"",
// 		           cert:"",
// 		           clientAuth:nil,
// 		           insecure:false,
// 		           hotReload:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			c := &credentials{
// 				crlNextUpdate: test.fields.crlNextUpdate,
// 				certPtr:       test.fields.certPtr,
// 				revoked:       test.fields.revoked,
// 				cfg:           test.fields.cfg,
// 				crl:           test.fields.crl,
// 				sn:            test.fields.sn,
// 				ca:            test.fields.ca,
// 				key:           test.fields.key,
// 				cert:          test.fields.cert,
// 				clientAuth:    test.fields.clientAuth,
// 				insecure:      test.fields.insecure,
// 				hotReload:     test.fields.hotReload,
// 			}
//
// 			got, err := c.reloadCert()
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_credentials_ensureCRL(t *testing.T) {
// 	type fields struct {
// 		crlNextUpdate time.Time
// 		certPtr       atomic.Pointer[tls.Certificate]
// 		revoked       map[string]struct{}
// 		cfg           *Config
// 		crl           string
// 		sn            string
// 		ca            string
// 		key           string
// 		cert          string
// 		clientAuth    tls.ClientAuthType
// 		insecure      bool
// 		hotReload     bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           crlNextUpdate:time.Time{},
// 		           certPtr:nil,
// 		           revoked:nil,
// 		           cfg:Config{},
// 		           crl:"",
// 		           sn:"",
// 		           ca:"",
// 		           key:"",
// 		           cert:"",
// 		           clientAuth:nil,
// 		           insecure:false,
// 		           hotReload:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           crlNextUpdate:time.Time{},
// 		           certPtr:nil,
// 		           revoked:nil,
// 		           cfg:Config{},
// 		           crl:"",
// 		           sn:"",
// 		           ca:"",
// 		           key:"",
// 		           cert:"",
// 		           clientAuth:nil,
// 		           insecure:false,
// 		           hotReload:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			c := &credentials{
// 				crlNextUpdate: test.fields.crlNextUpdate,
// 				certPtr:       test.fields.certPtr,
// 				revoked:       test.fields.revoked,
// 				cfg:           test.fields.cfg,
// 				crl:           test.fields.crl,
// 				sn:            test.fields.sn,
// 				ca:            test.fields.ca,
// 				key:           test.fields.key,
// 				cert:          test.fields.cert,
// 				clientAuth:    test.fields.clientAuth,
// 				insecure:      test.fields.insecure,
// 				hotReload:     test.fields.hotReload,
// 			}
//
// 			c.ensureCRL()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_credentials_attachCRLChainChecker(t *testing.T) {
// 	type args struct {
// 		cfg *tls.Config
// 	}
// 	type fields struct {
// 		crlNextUpdate time.Time
// 		certPtr       atomic.Pointer[tls.Certificate]
// 		revoked       map[string]struct{}
// 		cfg           *Config
// 		crl           string
// 		sn            string
// 		ca            string
// 		key           string
// 		cert          string
// 		clientAuth    tls.ClientAuthType
// 		insecure      bool
// 		hotReload     bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cfg:nil,
// 		       },
// 		       fields: fields {
// 		           crlNextUpdate:time.Time{},
// 		           certPtr:nil,
// 		           revoked:nil,
// 		           cfg:Config{},
// 		           crl:"",
// 		           sn:"",
// 		           ca:"",
// 		           key:"",
// 		           cert:"",
// 		           clientAuth:nil,
// 		           insecure:false,
// 		           hotReload:false,
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
// 		           cfg:nil,
// 		           },
// 		           fields: fields {
// 		           crlNextUpdate:time.Time{},
// 		           certPtr:nil,
// 		           revoked:nil,
// 		           cfg:Config{},
// 		           crl:"",
// 		           sn:"",
// 		           ca:"",
// 		           key:"",
// 		           cert:"",
// 		           clientAuth:nil,
// 		           insecure:false,
// 		           hotReload:false,
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
// 			c := &credentials{
// 				crlNextUpdate: test.fields.crlNextUpdate,
// 				certPtr:       test.fields.certPtr,
// 				revoked:       test.fields.revoked,
// 				cfg:           test.fields.cfg,
// 				crl:           test.fields.crl,
// 				sn:            test.fields.sn,
// 				ca:            test.fields.ca,
// 				key:           test.fields.key,
// 				cert:          test.fields.cert,
// 				clientAuth:    test.fields.clientAuth,
// 				insecure:      test.fields.insecure,
// 				hotReload:     test.fields.hotReload,
// 			}
//
// 			c.attachCRLChainChecker(test.args.cfg)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_systemOrNewPool(t *testing.T) {
// 	type want struct {
// 		want *x509.CertPool
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, *x509.CertPool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *x509.CertPool) error {
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
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := systemOrNewPool()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_processCert(t *testing.T) {
// 	type args struct {
// 		path string
// 		cert *x509.Certificate
// 		pool *x509.CertPool
// 		seen map[string]struct{}
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
// 		           path:"",
// 		           cert:nil,
// 		           pool:nil,
// 		           seen:nil,
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
// 		           path:"",
// 		           cert:nil,
// 		           pool:nil,
// 		           seen:nil,
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
// 			got := processCert(test.args.path, test.args.cert, test.args.pool, test.args.seen)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_verifyCertChain(t *testing.T) {
// 	type args struct {
// 		cert *x509.Certificate
// 		pool *x509.CertPool
// 		now  time.Time
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cert:nil,
// 		           pool:nil,
// 		           now:time.Time{},
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
// 		           cert:nil,
// 		           pool:nil,
// 		           now:time.Time{},
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
// 			err := verifyCertChain(test.args.cert, test.args.pool, test.args.now)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_fingerprint(t *testing.T) {
// 	type args struct {
// 		data []byte
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
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
// 		           data:nil,
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
// 		           data:nil,
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
// 			got := fingerprint(test.args.data)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_parsePEMCertificates(t *testing.T) {
// 	type args struct {
// 		pemBytes []byte
// 	}
// 	type want struct {
// 		want []*x509.Certificate
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []*x509.Certificate, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []*x509.Certificate, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           pemBytes:nil,
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
// 		           pemBytes:nil,
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
// 			got, err := parsePEMCertificates(test.args.pemBytes)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
