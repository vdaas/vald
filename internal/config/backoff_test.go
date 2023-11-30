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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestBackoff_Bind(t *testing.T) {
	type fields struct {
		InitialDuration  string
		BackoffTimeLimit string
		MaximumDuration  string
		JitterLimit      string
		BackoffFactor    float64
		RetryCount       int
		EnableErrorLog   bool
	}
	type want struct {
		want *Backoff
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Backoff) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Backoff) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return Backoff struct when values are not empty",
				fields: fields{
					InitialDuration:  "5m",
					BackoffTimeLimit: "10m",
					MaximumDuration:  "15m",
					JitterLimit:      "3m",
				},
				want: want{
					want: &Backoff{
						InitialDuration:  "5m",
						BackoffTimeLimit: "10m",
						MaximumDuration:  "15m",
						JitterLimit:      "3m",
					},
				},
			}
		}(),
		func() test {
			key := "BACKOFF_BIND_INITIAL_DURATION"
			val := "5m"
			return test{
				name: "return Backoff struct when initialDuration is set via the environment value",
				fields: fields{
					InitialDuration: "_" + key + "_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					t.Setenv(key, val)
				},
				want: want{
					want: &Backoff{
						InitialDuration: val,
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return Backoff struct when values are empty",
				fields: fields{},
				want: want{
					want: &Backoff{},
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
			b := &Backoff{
				InitialDuration:  test.fields.InitialDuration,
				BackoffTimeLimit: test.fields.BackoffTimeLimit,
				MaximumDuration:  test.fields.MaximumDuration,
				JitterLimit:      test.fields.JitterLimit,
			}

			got := b.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestBackoff_Opts(t *testing.T) {
	type fields struct {
		InitialDuration  string
		BackoffTimeLimit string
		MaximumDuration  string
		JitterLimit      string
		BackoffFactor    float64
		RetryCount       int
		EnableErrorLog   bool
	}
	type want struct {
		want []backoff.Option
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []backoff.Option) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []backoff.Option) error {
		if !reflect.DeepEqual(len(w.want), len(got)) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return 6 backoff.Option when EnableErrorLog false and other values are set",
			fields: fields{
				InitialDuration:  "5m",
				BackoffTimeLimit: "10m",
				MaximumDuration:  "15m",
				JitterLimit:      "3m",
				BackoffFactor:    3,
				RetryCount:       100,
				EnableErrorLog:   false,
			},
			want: want{
				want: make([]backoff.Option, 6, 7),
			},
		},
		{
			name: "return 7 backoff.Option when EnableErrorLog false and other values are set",
			fields: fields{
				InitialDuration:  "5m",
				BackoffTimeLimit: "10m",
				MaximumDuration:  "15m",
				JitterLimit:      "3m",
				BackoffFactor:    3,
				RetryCount:       100,
				EnableErrorLog:   true,
			},
			want: want{
				want: make([]backoff.Option, 7),
			},
		},
		{
			name:   "return 6 backoff.Option when all values are set as default value",
			fields: fields{},
			want: want{
				want: make([]backoff.Option, 6, 7),
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
			b := &Backoff{
				InitialDuration:  test.fields.InitialDuration,
				BackoffTimeLimit: test.fields.BackoffTimeLimit,
				MaximumDuration:  test.fields.MaximumDuration,
				JitterLimit:      test.fields.JitterLimit,
				BackoffFactor:    test.fields.BackoffFactor,
				RetryCount:       test.fields.RetryCount,
				EnableErrorLog:   test.fields.EnableErrorLog,
			}

			got := b.Opts()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
