// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/errors"
)

func TestCircuitBreaker_breakerOpts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		cb   *CircuitBreaker
		name string
		want int
	}{
		{
			name: "return no option when nothing is set",
			cb:   new(CircuitBreaker),
			want: 0,
		},
		{
			name: "return only the timeout and min samples options when the error rates are unset",
			cb: &CircuitBreaker{
				OpenTimeout:          "1s",
				ClosedRefreshTimeout: "1m",
				MinSamples:           10,
			},
			want: 3,
		},
		{
			name: "return all options when everything is set",
			cb: &CircuitBreaker{
				OpenTimeout:          "1s",
				ClosedRefreshTimeout: "1m",
				MinSamples:           10,
				ClosedErrorRate:      0.7,
				HalfOpenErrorRate:    0.5,
			},
			want: 5,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			if got := len(tc.cb.breakerOpts()); got != tc.want {
				tt.Errorf("len(breakerOpts()) = %d, want %d", got, tc.want)
			}
		})
	}
}

// An unset error rate must keep the breaker default (0.7) instead of becoming 0,
// which would trip on the first failure once min_samples is exceeded.
func TestCircuitBreaker_breakerOpts_unsetErrorRateKeepsDefault(t *testing.T) {
	t.Parallel()
	cfg := &CircuitBreaker{
		OpenTimeout:          "1h",
		ClosedRefreshTimeout: "1h",
		MinSamples:           1,
	}
	cb, err := circuitbreaker.New(circuitbreaker.WithBreakerOpts(cfg.breakerOpts()...))
	if err != nil {
		t.Fatalf("circuitbreaker.New() error = %v", err)
	}
	ctx := t.Context()
	key := "config-test-unset-error-rate"
	errFail := errors.New("fail")
	success := func(context.Context) (any, error) { return struct{}{}, nil }
	fail := func(context.Context) (any, error) { return nil, errFail }
	for i, fn := range []func(context.Context) (any, error){success, success, fail} {
		if _, err := cb.Do(ctx, key, fn); err != nil && !errors.Is(err, errFail) {
			t.Fatalf("Do() #%d error = %v", i, err)
		}
	}
	// 1 failure out of 3 samples is below the default closed error rate, so the breaker must stay closed.
	if _, err := cb.Do(ctx, key, success); err != nil {
		t.Errorf("Do() after one failure error = %v, want nil (breaker must stay closed)", err)
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestCircuitBreaker_Bind(t *testing.T) {
// 	type fields struct {
// 		OpenTimeout          string
// 		ClosedRefreshTimeout string
// 		MinSamples           int64
// 		ClosedErrorRate      float32
// 		HalfOpenErrorRate    float32
// 	}
// 	type want struct {
// 		want *CircuitBreaker
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *CircuitBreaker) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *CircuitBreaker) error {
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
// 		           OpenTimeout:"",
// 		           ClosedRefreshTimeout:"",
// 		           MinSamples:0,
// 		           ClosedErrorRate:0,
// 		           HalfOpenErrorRate:0,
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
// 		           OpenTimeout:"",
// 		           ClosedRefreshTimeout:"",
// 		           MinSamples:0,
// 		           ClosedErrorRate:0,
// 		           HalfOpenErrorRate:0,
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
// 			cb := &CircuitBreaker{
// 				OpenTimeout:          test.fields.OpenTimeout,
// 				ClosedRefreshTimeout: test.fields.ClosedRefreshTimeout,
// 				MinSamples:           test.fields.MinSamples,
// 				ClosedErrorRate:      test.fields.ClosedErrorRate,
// 				HalfOpenErrorRate:    test.fields.HalfOpenErrorRate,
// 			}
//
// 			got := cb.Bind()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
