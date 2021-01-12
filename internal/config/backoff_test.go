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

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Backoff) error {
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
		           InitialDuration: "",
		           BackoffTimeLimit: "",
		           MaximumDuration: "",
		           JitterLimit: "",
		           BackoffFactor: 0,
		           RetryCount: 0,
		           EnableErrorLog: false,
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
		           InitialDuration: "",
		           BackoffTimeLimit: "",
		           MaximumDuration: "",
		           JitterLimit: "",
		           BackoffFactor: 0,
		           RetryCount: 0,
		           EnableErrorLog: false,
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
			b := &Backoff{
				InitialDuration:  test.fields.InitialDuration,
				BackoffTimeLimit: test.fields.BackoffTimeLimit,
				MaximumDuration:  test.fields.MaximumDuration,
				JitterLimit:      test.fields.JitterLimit,
				BackoffFactor:    test.fields.BackoffFactor,
				RetryCount:       test.fields.RetryCount,
				EnableErrorLog:   test.fields.EnableErrorLog,
			}

			got := b.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
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
		           InitialDuration: "",
		           BackoffTimeLimit: "",
		           MaximumDuration: "",
		           JitterLimit: "",
		           BackoffFactor: 0,
		           RetryCount: 0,
		           EnableErrorLog: false,
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
		           InitialDuration: "",
		           BackoffTimeLimit: "",
		           MaximumDuration: "",
		           JitterLimit: "",
		           BackoffFactor: 0,
		           RetryCount: 0,
		           EnableErrorLog: false,
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
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
