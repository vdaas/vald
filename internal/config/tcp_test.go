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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/tcp"
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *DNS) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           CacheEnabled: false,
		           RefreshDuration: "",
		           CacheExpiration: "",
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
		           CacheEnabled: false,
		           RefreshDuration: "",
		           CacheExpiration: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			d := &DNS{
				CacheEnabled:    test.fields.CacheEnabled,
				RefreshDuration: test.fields.RefreshDuration,
				CacheExpiration: test.fields.CacheExpiration,
			}

			got := d.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDialer_Bind(t *testing.T) {
	type fields struct {
		Timeout          string
		KeepAlive        string
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Dialer) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           Timeout: "",
		           KeepAlive: "",
		           DualStackEnabled: false,
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
		           Timeout: "",
		           KeepAlive: "",
		           DualStackEnabled: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			d := &Dialer{
				Timeout:          test.fields.Timeout,
				KeepAlive:        test.fields.KeepAlive,
				DualStackEnabled: test.fields.DualStackEnabled,
			}

			got := d.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestTCP_Bind(t *testing.T) {
	type fields struct {
		DNS    *DNS
		Dialer *Dialer
		TLS    *TLS
	}
	type want struct {
		want *TCP
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *TCP) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *TCP) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           DNS: DNS{},
		           Dialer: Dialer{},
		           TLS: TLS{},
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
		           DNS: DNS{},
		           Dialer: Dialer{},
		           TLS: TLS{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			t := &TCP{
				DNS:    test.fields.DNS,
				Dialer: test.fields.Dialer,
				TLS:    test.fields.TLS,
			}

			got := t.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestTCP_Opts(t *testing.T) {
	type fields struct {
		DNS    *DNS
		Dialer *Dialer
		TLS    *TLS
	}
	type want struct {
		want []tcp.DialerOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []tcp.DialerOption) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []tcp.DialerOption) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           DNS: DNS{},
		           Dialer: Dialer{},
		           TLS: TLS{},
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
		           DNS: DNS{},
		           Dialer: Dialer{},
		           TLS: TLS{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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
			t := &TCP{
				DNS:    test.fields.DNS,
				Dialer: test.fields.Dialer,
				TLS:    test.fields.TLS,
			}

			got := t.Opts()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
