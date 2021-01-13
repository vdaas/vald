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
)

func TestRoundTripper_Bind(t *testing.T) {
	type fields struct {
		TLSHandshakeTimeout   string
		MaxIdleConns          int
		MaxIdleConnsPerHost   int
		MaxConnsPerHost       int
		IdleConnTimeout       string
		ResponseHeaderTimeout string
		ExpectContinueTimeout string
		MaxResponseHeaderSize int64
		WriteBufferSize       int64
		ReadBufferSize        int64
		ForceAttemptHTTP2     bool
	}
	type want struct {
		want *RoundTripper
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *RoundTripper) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *RoundTripper) error {
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
		           TLSHandshakeTimeout: "",
		           MaxIdleConns: 0,
		           MaxIdleConnsPerHost: 0,
		           MaxConnsPerHost: 0,
		           IdleConnTimeout: "",
		           ResponseHeaderTimeout: "",
		           ExpectContinueTimeout: "",
		           MaxResponseHeaderSize: 0,
		           WriteBufferSize: 0,
		           ReadBufferSize: 0,
		           ForceAttemptHTTP2: false,
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
		           TLSHandshakeTimeout: "",
		           MaxIdleConns: 0,
		           MaxIdleConnsPerHost: 0,
		           MaxConnsPerHost: 0,
		           IdleConnTimeout: "",
		           ResponseHeaderTimeout: "",
		           ExpectContinueTimeout: "",
		           MaxResponseHeaderSize: 0,
		           WriteBufferSize: 0,
		           ReadBufferSize: 0,
		           ForceAttemptHTTP2: false,
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
			r := &RoundTripper{
				TLSHandshakeTimeout:   test.fields.TLSHandshakeTimeout,
				MaxIdleConns:          test.fields.MaxIdleConns,
				MaxIdleConnsPerHost:   test.fields.MaxIdleConnsPerHost,
				MaxConnsPerHost:       test.fields.MaxConnsPerHost,
				IdleConnTimeout:       test.fields.IdleConnTimeout,
				ResponseHeaderTimeout: test.fields.ResponseHeaderTimeout,
				ExpectContinueTimeout: test.fields.ExpectContinueTimeout,
				MaxResponseHeaderSize: test.fields.MaxResponseHeaderSize,
				WriteBufferSize:       test.fields.WriteBufferSize,
				ReadBufferSize:        test.fields.ReadBufferSize,
				ForceAttemptHTTP2:     test.fields.ForceAttemptHTTP2,
			}

			got := r.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestTransport_Bind(t *testing.T) {
	type fields struct {
		RoundTripper *RoundTripper
		Backoff      *Backoff
	}
	type want struct {
		want *Transport
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Transport) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Transport) error {
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
		           RoundTripper: RoundTripper{},
		           Backoff: Backoff{},
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
		           RoundTripper: RoundTripper{},
		           Backoff: Backoff{},
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
			t := &Transport{
				RoundTripper: test.fields.RoundTripper,
				Backoff:      test.fields.Backoff,
			}

			got := t.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
