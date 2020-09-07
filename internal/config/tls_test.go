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

// Package config providers configuration type and load configuration logic
package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/tls"
)

func TestTLS_Bind(t *testing.T) {
	type fields struct {
		Enabled bool
		Cert    string
		Key     string
		CA      string
	}
	type want struct {
		want *TLS
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *TLS) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *TLS) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns TLS when all fields contain no prefix/suffix symbol",
			fields: fields{
				Enabled: true,
				Cert:    "cert",
				Key:     "key",
				CA:      "ca",
			},
			want: want{
				want: &TLS{
					Enabled: true,
					Cert:    "cert",
					Key:     "key",
					CA:      "ca",
				},
			},
		},
		{
			name: "returns TLS with environment variable when it contains `_` prefix and suffix",
			fields: fields{
				Enabled: true,
				Cert:    "_cert_",
				Key:     "_key_",
				CA:      "_ca_",
			},
			beforeFunc: func() {
				_ = os.Setenv("cert", "tls_cert")
				_ = os.Setenv("key", "tls_key")
				_ = os.Setenv("ca", "tls_ca")
			},
			afterFunc: func() {
				_ = os.Unsetenv("cert")
				_ = os.Unsetenv("key")
				_ = os.Unsetenv("ca")
			},
			want: want{
				want: &TLS{
					Enabled: true,
					Cert:    "tls_cert",
					Key:     "tls_key",
					CA:      "tls_ca",
				},
			},
		},
		{
			name: "returns TLS when all fields are empty",
			want: want{
				want: new(TLS),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			t := &TLS{
				Enabled: test.fields.Enabled,
				Cert:    test.fields.Cert,
				Key:     test.fields.Key,
				CA:      test.fields.CA,
			}

			got := t.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestTLS_Opts(t *testing.T) {
	type fields struct {
		Enabled bool
		Cert    string
		Key     string
		CA      string
	}
	type want struct {
		want []tls.Option
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []tls.Option) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []tls.Option) error {
		if len(w.want) != len(got) {
			return errors.Errorf("len(got) = %d, len(want) = %d", len(got), len(w.want))
		}
		for i := range w.want {
			ok := false
			for j := range got {
				if reflect.ValueOf(w.want[i]).Pointer() == reflect.ValueOf(got[j]).Pointer() {
					ok = true
					break
				}
			}
			if !ok {
				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "returns []tls.Option",
			fields: fields{
				Enabled: true,
				Cert:    "cert",
				Key:     "key",
				CA:      "ca",
			},
			want: want{
				want: []tls.Option{
					tls.WithCa("ca"),
					tls.WithCert("cert"),
					tls.WithKey("key"),
				},
			},
		},
		{
			name: "returns []tls.Option",
			want: want{
				want: []tls.Option{
					tls.WithCa(""),
					tls.WithCert(""),
					tls.WithKey(""),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			t := &TLS{
				Enabled: test.fields.Enabled,
				Cert:    test.fields.Cert,
				Key:     test.fields.Key,
				CA:      test.fields.CA,
			}

			got := t.Opts()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
