//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestCircuitBreaker_Bind(t *testing.T) {
	type fields struct {
		ClosedErrorRate      float32
		HalfOpenErrorRate    float32
		MinSamples           int64
		OpenTimeout          string
		ClosedRefreshTimeout string
	}
	type want struct {
		want *CircuitBreaker
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *CircuitBreaker) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *CircuitBreaker) error {
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
		           ClosedErrorRate: 0,
		           HalfOpenErrorRate: 0,
		           MinSamples: 0,
		           OpenTimeout: "",
		           ClosedRefreshTimeout: "",
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
		           ClosedErrorRate: 0,
		           HalfOpenErrorRate: 0,
		           MinSamples: 0,
		           OpenTimeout: "",
		           ClosedRefreshTimeout: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			cb := &CircuitBreaker{
				ClosedErrorRate:      test.fields.ClosedErrorRate,
				HalfOpenErrorRate:    test.fields.HalfOpenErrorRate,
				MinSamples:           test.fields.MinSamples,
				OpenTimeout:          test.fields.OpenTimeout,
				ClosedRefreshTimeout: test.fields.ClosedRefreshTimeout,
			}

			got := cb.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
