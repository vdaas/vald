//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/test/goleak"
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
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *RoundTripper) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name:   "return RoundTripper when all parameters are nil",
				fields: fields{},
				want: want{
					want: new(RoundTripper),
				},
			}
		}(),
		func() test {
			tlsHandshakeTimeout := "5s"
			maxIdleConns := 20
			maxIdleConnsPerHost := 3
			maxConnsPerHost := 10
			idleConnTimeout := "10s"
			responseHeaderTimeout := "5s"
			expectContinueTimeout := "5s"
			maxResponseHeaderSize := int64(20)
			writeBufferSize := int64(2000)
			readBufferSize := int64(2000)
			forceAttemptHTTP2 := true

			return test{
				name: "return RoundTripper when all parameters are not nil",
				fields: fields{
					TLSHandshakeTimeout:   tlsHandshakeTimeout,
					MaxIdleConns:          maxIdleConns,
					MaxIdleConnsPerHost:   maxIdleConnsPerHost,
					MaxConnsPerHost:       maxConnsPerHost,
					IdleConnTimeout:       idleConnTimeout,
					ResponseHeaderTimeout: responseHeaderTimeout,
					ExpectContinueTimeout: expectContinueTimeout,
					MaxResponseHeaderSize: maxResponseHeaderSize,
					WriteBufferSize:       writeBufferSize,
					ReadBufferSize:        readBufferSize,
					ForceAttemptHTTP2:     forceAttemptHTTP2,
				},
				want: want{
					want: &RoundTripper{
						TLSHandshakeTimeout:   tlsHandshakeTimeout,
						MaxIdleConns:          maxIdleConns,
						MaxIdleConnsPerHost:   maxIdleConnsPerHost,
						MaxConnsPerHost:       maxConnsPerHost,
						IdleConnTimeout:       idleConnTimeout,
						ResponseHeaderTimeout: responseHeaderTimeout,
						ExpectContinueTimeout: expectContinueTimeout,
						MaxResponseHeaderSize: maxResponseHeaderSize,
						WriteBufferSize:       writeBufferSize,
						ReadBufferSize:        readBufferSize,
						ForceAttemptHTTP2:     forceAttemptHTTP2,
					},
				},
			}
		}(),
		func() test {
			tlsHandshakeTimeout := "5s"
			idleConnTimeout := "10s"
			responseHeaderTimeout := "5s"
			expectContinueTimeout := "5s"
			envPrefix := "bind_env_test_rt"
			m := map[string]string{
				envPrefix + "tlsHandshakeTimeout":   tlsHandshakeTimeout,
				envPrefix + "idleConnTimeout":       idleConnTimeout,
				envPrefix + "responseHeaderTimeout": responseHeaderTimeout,
				envPrefix + "expectContinueTimeout": expectContinueTimeout,
			}

			return test{
				name: "return RoundTripper when the data is loaded environment variable",
				fields: fields{
					TLSHandshakeTimeout:   "_" + envPrefix + "tlsHandshakeTimeout_",
					IdleConnTimeout:       "_" + envPrefix + "idleConnTimeout_",
					ResponseHeaderTimeout: "_" + envPrefix + "responseHeaderTimeout_",
					ExpectContinueTimeout: "_" + envPrefix + "expectContinueTimeout_",
				},
				want: want{
					want: &RoundTripper{
						TLSHandshakeTimeout:   tlsHandshakeTimeout,
						IdleConnTimeout:       idleConnTimeout,
						ResponseHeaderTimeout: responseHeaderTimeout,
						ExpectContinueTimeout: expectContinueTimeout,
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						t.Setenv(k, v)
					}
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
			if err := checkFunc(test.want, got); err != nil {
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
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Transport) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name:   "return Transport when all parameters are nil",
				fields: fields{},
				want: want{
					want: &Transport{
						RoundTripper: new(RoundTripper),
						Backoff:      new(Backoff),
					},
				},
			}
		}(),
		func() test {
			roundTripper := &RoundTripper{
				TLSHandshakeTimeout: "1s",
			}
			backoff := &Backoff{
				InitialDuration: "1s",
			}

			return test{
				name: "return Transport when all parameters are not nil",
				fields: fields{
					RoundTripper: roundTripper,
					Backoff:      backoff,
				},
				want: want{
					&Transport{
						RoundTripper: roundTripper,
						Backoff:      backoff,
					},
				},
			}
		}(),
		func() test {
			tlsHandshakeTimeout := "5s"
			initialDuration := "1s"
			envPrefix := "bind_env_test_t"
			m := map[string]string{
				envPrefix + "tlsHandshakeTimeout": tlsHandshakeTimeout,
				envPrefix + "initialDuration":     initialDuration,
			}

			return test{
				name: "return Transport when the data is loaded environment variable",
				fields: fields{
					RoundTripper: &RoundTripper{
						TLSHandshakeTimeout: "_" + envPrefix + "tlsHandshakeTimeout_",
					},
					Backoff: &Backoff{
						InitialDuration: "_" + envPrefix + "initialDuration_",
					},
				},
				want: want{
					&Transport{
						RoundTripper: &RoundTripper{
							TLSHandshakeTimeout: tlsHandshakeTimeout,
						},
						Backoff: &Backoff{
							InitialDuration: initialDuration,
						},
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						t.Setenv(k, v)
					}
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
			t := &Transport{
				RoundTripper: test.fields.RoundTripper,
				Backoff:      test.fields.Backoff,
			}

			got := t.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
