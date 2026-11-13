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
package circuitbreaker

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func Test_breaker_isReady(t *testing.T) {
	type fields struct {
		closedErrShouldTrip   Tripper
		halfOpenErrShouldTrip Tripper
		count                 *count
		key                   string
		minSamples            int64
		openTimeout           time.Duration
		openExp               int64
		closedRefreshTimeout  time.Duration
		closedRefreshExp      int64
		tripped               int32
		closedErrRate         float32
		halfOpenErrRate       float32
	}
	type want struct {
		err    error
		wantSt State
	}
	type test struct {
		want       want
		checkFunc  func(want, State, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		name       string
		fields     fields
	}
	defaultCheckFunc := func(w want, gotSt State, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotSt, w.wantSt) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSt, w.wantSt)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return the StateClose and nil when the current state is Close",
				fields: fields{
					key:              "insertRPC",
					tripped:          0,
					closedRefreshExp: time.Now().Add(100 * time.Second).UnixNano(),
				},
				want: want{
					wantSt: StateClosed,
					err:    nil,
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
		func() test {
			cnt := &count{
				successes: 1,
			}
			return test{
				name: "return the StateHalfOpen and nil when the current state is HalfOpen",
				fields: fields{
					key:     "insertRPC",
					tripped: 1,
					openExp: time.Now().Add(-100 * time.Second).UnixNano(),
					count:   cnt,
				},
				want: want{
					wantSt: StateHalfOpen,
					err:    nil,
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
		func() test {
			cnt := &count{}
			return test{
				name: "return the StateHalfOpen and error when the current state is HalfOpen but the flow is being limited",
				fields: fields{
					key:     "insertRPC",
					tripped: 1,
					openExp: time.Now().Add(-100 * time.Second).UnixNano(),
					count:   cnt,
				},
				want: want{
					wantSt: StateHalfOpen,
					err:    errors.ErrCircuitBreakerHalfOpenFlowLimitation,
				},
				checkFunc: func(w want, s State, err error) error {
					if err := defaultCheckFunc(w, s, err); err != nil {
						return err
					}
					if got := cnt.Fails(); got != 0 {
						return fmt.Errorf("failures is not equals. want: %d, but got: %d", 0, got)
					}
					return nil
				},
			}
		}(),
		func() test {
			return test{
				name: "return the StateOpen and error when the current state is Open",
				fields: fields{
					key:     "insertRPC",
					tripped: 1,
					openExp: time.Now().Add(100 * time.Second).UnixNano(),
				},
				want: want{
					wantSt: StateOpen,
					err:    errors.ErrCircuitBreakerOpenState,
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			b := &breaker{
				key:                   test.fields.key,
				count:                 test.fields.count,
				tripped:               test.fields.tripped,
				closedErrRate:         test.fields.closedErrRate,
				closedErrShouldTrip:   test.fields.closedErrShouldTrip,
				halfOpenErrRate:       test.fields.halfOpenErrRate,
				halfOpenErrShouldTrip: test.fields.halfOpenErrShouldTrip,
				minSamples:            test.fields.minSamples,
				openTimeout:           test.fields.openTimeout,
				openExp:               test.fields.openExp,
				closedRefreshTimeout:  test.fields.closedRefreshTimeout,
				closedRefreshExp:      test.fields.closedRefreshExp,
			}

			gotSt, err := b.isReady()
			if err := checkFunc(test.want, gotSt, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_breaker_success(t *testing.T) {
	type fields struct {
		closedErrShouldTrip   Tripper
		halfOpenErrShouldTrip Tripper
		count                 *count
		key                   string
		minSamples            int64
		openTimeout           time.Duration
		openExp               int64
		closedRefreshTimeout  time.Duration
		closedRefreshExp      int64
		tripped               int32
		closedErrRate         float32
		halfOpenErrRate       float32
	}
	type want struct{}
	type test struct {
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T, *breaker)
		name       string
		fields     fields
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		func() test {
			cnt := &count{
				successes: 10,
				failures:  10,
			}
			halfOpenErrRate := float32(0.5)
			minSamples := int64(10)
			return test{
				name: "the current state change from HalfOpen to Close when the success rate is higher",
				fields: fields{
					key:                   "insertRPC",
					count:                 cnt,
					tripped:               1,
					openExp:               time.Now().Add(-100 * time.Second).UnixNano(),
					halfOpenErrRate:       halfOpenErrRate,
					halfOpenErrShouldTrip: NewRateTripper(halfOpenErrRate, minSamples),
					minSamples:            minSamples,
				},
				checkFunc: defaultCheckFunc,
				afterFunc: func(t *testing.T, b *breaker) {
					t.Helper()
					if b.tripped != 0 {
						t.Errorf("state did not change: %d", b.tripped)
					}
				},
			}
		}(),
		func() test {
			cnt := &count{
				successes: 10,
				failures:  11,
			}
			halfOpenErrRate := float32(0.5)
			minSamples := int64(10)
			return test{
				name: "the current state do not change from HalfOpen to Close when the success rate is less",
				fields: fields{
					key:                   "insertRPC",
					count:                 cnt,
					tripped:               1,
					openExp:               time.Now().Add(-100 * time.Second).UnixNano(),
					halfOpenErrRate:       halfOpenErrRate,
					halfOpenErrShouldTrip: NewRateTripper(halfOpenErrRate, minSamples),
					minSamples:            minSamples,
				},
				checkFunc: defaultCheckFunc,
				afterFunc: func(t *testing.T, b *breaker) {
					t.Helper()
					if b.tripped != 1 {
						t.Errorf("state changed: %d", b.tripped)
					}
				},
			}
		}(),
		func() test {
			// success() counts this call first, so the samples reach exactly
			// minSamples (9 failures + 1 success). The tripper ignores that many
			// samples, so the breaker used to close at a 90% error rate.
			cnt := &count{
				successes: 0,
				failures:  9,
			}
			halfOpenErrRate := float32(0.5)
			minSamples := int64(10)
			return test{
				name: "the current state do not change from HalfOpen to Close at exactly minSamples",
				fields: fields{
					key:                   "insertRPC",
					count:                 cnt,
					tripped:               1,
					openExp:               time.Now().Add(-100 * time.Second).UnixNano(),
					halfOpenErrRate:       halfOpenErrRate,
					halfOpenErrShouldTrip: NewRateTripper(halfOpenErrRate, minSamples),
					minSamples:            minSamples,
				},
				checkFunc: defaultCheckFunc,
				afterFunc: func(t *testing.T, b *breaker) {
					t.Helper()
					if b.tripped != 1 {
						t.Errorf("state changed at exactly minSamples: %d", b.tripped)
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			b := &breaker{
				key:                   test.fields.key,
				count:                 test.fields.count,
				tripped:               test.fields.tripped,
				closedErrRate:         test.fields.closedErrRate,
				closedErrShouldTrip:   test.fields.closedErrShouldTrip,
				halfOpenErrRate:       test.fields.halfOpenErrRate,
				halfOpenErrShouldTrip: test.fields.halfOpenErrShouldTrip,
				minSamples:            test.fields.minSamples,
				openTimeout:           test.fields.openTimeout,
				openExp:               test.fields.openExp,
				closedRefreshTimeout:  test.fields.closedRefreshTimeout,
				closedRefreshExp:      test.fields.closedRefreshExp,
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, b)
			}

			b.success()
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_breaker_fail(t *testing.T) {
	type fields struct {
		closedErrShouldTrip   Tripper
		halfOpenErrShouldTrip Tripper
		count                 *count
		key                   string
		minSamples            int64
		openTimeout           time.Duration
		openExp               int64
		closedRefreshTimeout  time.Duration
		closedRefreshExp      int64
		tripped               int32
		closedErrRate         float32
		halfOpenErrRate       float32
	}
	type want struct{}
	type test struct {
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T, *breaker)
		name       string
		fields     fields
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		func() test {
			cnt := &count{
				successes: 10,
				failures:  11,
			}
			closedErrRate := float32(0.5)
			minSamples := int64(20)
			return test{
				name: "the current state change from Close to Open when the failure rate is higher",
				fields: fields{
					key:                 "insertRPC",
					count:               cnt,
					tripped:             0,
					closedErrRate:       closedErrRate,
					closedRefreshExp:    time.Now().Add(100 * time.Second).UnixNano(),
					closedErrShouldTrip: NewRateTripper(closedErrRate, minSamples),
					minSamples:          minSamples,
				},
				checkFunc: defaultCheckFunc,
				afterFunc: func(t *testing.T, b *breaker) {
					t.Helper()
					if b.tripped == 0 {
						t.Errorf("state did not change: %d", b.tripped)
					}
					if total := cnt.Total(); total != 0 {
						t.Errorf("count did not reset: %d", total)
					}
				},
			}
		}(),
		func() test {
			cnt := &count{
				successes: 10,
				failures:  11,
			}
			halfOpenErrRate := float32(0.5)
			minSamples := int64(20)
			return test{
				name: "the current state change from HalfOpen to Open when the failure rate is higher",
				fields: fields{
					key:                   "insertRPC",
					count:                 cnt,
					tripped:               1,
					openExp:               time.Now().Add(-100 * time.Second).UnixNano(),
					halfOpenErrRate:       halfOpenErrRate,
					halfOpenErrShouldTrip: NewRateTripper(halfOpenErrRate, minSamples),
					minSamples:            minSamples,
				},
				checkFunc: defaultCheckFunc,
				afterFunc: func(t *testing.T, b *breaker) {
					t.Helper()
					if b.tripped == 0 {
						t.Errorf("state changed: %d", b.tripped)
					}
					if total := b.count.Total(); total != 0 {
						t.Errorf("count did not reset: %d", total)
					}
				},
			}
		}(),
		func() test {
			cnt := &count{
				successes: 10,
				failures:  1,
			}
			halfOpenErrRate := float32(0.5)
			minSamples := int64(10)
			return test{
				name: "the current HalfOpen state dot not change when the failure rate does not reached the setting value",
				fields: fields{
					key:                   "insertRPC",
					count:                 cnt,
					tripped:               1,
					openExp:               time.Now().Add(-100 * time.Second).UnixNano(),
					halfOpenErrRate:       halfOpenErrRate,
					halfOpenErrShouldTrip: NewRateTripper(halfOpenErrRate, minSamples),
					minSamples:            minSamples,
				},
				checkFunc: defaultCheckFunc,
				afterFunc: func(t *testing.T, b *breaker) {
					t.Helper()
					if b.tripped == 0 {
						t.Errorf("state changed: %d", b.tripped)
					}
					if total := b.count.Total(); total == 0 {
						t.Errorf("count resetted: %d", total)
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			b := &breaker{
				key:                   test.fields.key,
				count:                 test.fields.count,
				tripped:               test.fields.tripped,
				closedErrRate:         test.fields.closedErrRate,
				closedErrShouldTrip:   test.fields.closedErrShouldTrip,
				halfOpenErrRate:       test.fields.halfOpenErrRate,
				halfOpenErrShouldTrip: test.fields.halfOpenErrShouldTrip,
				minSamples:            test.fields.minSamples,
				openTimeout:           test.fields.openTimeout,
				openExp:               test.fields.openExp,
				closedRefreshTimeout:  test.fields.closedRefreshTimeout,
				closedRefreshExp:      test.fields.closedRefreshExp,
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, b)
			}

			b.fail()
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_newBreaker(t *testing.T) {
// 	type args struct {
// 		key  string
// 		opts []BreakerOption
// 	}
// 	type want struct {
// 		want *breaker
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *breaker, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *breaker, err error) error {
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
// 		           key:"",
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
// 		           key:"",
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
// 			got, err := newBreaker(test.args.key, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_breaker_do(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		fn  func(ctx context.Context) (val any, err error)
// 	}
// 	type fields struct {
// 		closedErrShouldTrip   Tripper
// 		halfOpenErrShouldTrip Tripper
// 		count                 *count
// 		key                   string
// 		minSamples            int64
// 		openTimeout           time.Duration
// 		openExp               int64
// 		closedRefreshTimeout  time.Duration
// 		closedRefreshExp      int64
// 		tripped               int32
// 		closedErrRate         float32
// 		halfOpenErrRate       float32
// 	}
// 	type want struct {
// 		wantVal any
// 		wantSt  State
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, any, State, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVal any, gotSt State, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVal, w.wantVal) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVal, w.wantVal)
// 		}
// 		if !reflect.DeepEqual(gotSt, w.wantSt) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSt, w.wantSt)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 		           ctx:nil,
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 			b := &breaker{
// 				closedErrShouldTrip:   test.fields.closedErrShouldTrip,
// 				halfOpenErrShouldTrip: test.fields.halfOpenErrShouldTrip,
// 				count:                 test.fields.count,
// 				key:                   test.fields.key,
// 				minSamples:            test.fields.minSamples,
// 				openTimeout:           test.fields.openTimeout,
// 				openExp:               test.fields.openExp,
// 				closedRefreshTimeout:  test.fields.closedRefreshTimeout,
// 				closedRefreshExp:      test.fields.closedRefreshExp,
// 				tripped:               test.fields.tripped,
// 				closedErrRate:         test.fields.closedErrRate,
// 				halfOpenErrRate:       test.fields.halfOpenErrRate,
// 			}
//
// 			gotVal, gotSt, err := b.do(test.args.ctx, test.args.fn)
// 			if err := checkFunc(test.want, gotVal, gotSt, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_breaker_currentState(t *testing.T) {
// 	type fields struct {
// 		closedErrShouldTrip   Tripper
// 		halfOpenErrShouldTrip Tripper
// 		count                 *count
// 		key                   string
// 		minSamples            int64
// 		openTimeout           time.Duration
// 		openExp               int64
// 		closedRefreshTimeout  time.Duration
// 		closedRefreshExp      int64
// 		tripped               int32
// 		closedErrRate         float32
// 		halfOpenErrRate       float32
// 	}
// 	type want struct {
// 		want State
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, State) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got State) error {
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
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 			b := &breaker{
// 				closedErrShouldTrip:   test.fields.closedErrShouldTrip,
// 				halfOpenErrShouldTrip: test.fields.halfOpenErrShouldTrip,
// 				count:                 test.fields.count,
// 				key:                   test.fields.key,
// 				minSamples:            test.fields.minSamples,
// 				openTimeout:           test.fields.openTimeout,
// 				openExp:               test.fields.openExp,
// 				closedRefreshTimeout:  test.fields.closedRefreshTimeout,
// 				closedRefreshExp:      test.fields.closedRefreshExp,
// 				tripped:               test.fields.tripped,
// 				closedErrRate:         test.fields.closedErrRate,
// 				halfOpenErrRate:       test.fields.halfOpenErrRate,
// 			}
//
// 			got := b.currentState()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_breaker_reset(t *testing.T) {
// 	type fields struct {
// 		closedErrShouldTrip   Tripper
// 		halfOpenErrShouldTrip Tripper
// 		count                 *count
// 		key                   string
// 		minSamples            int64
// 		openTimeout           time.Duration
// 		openExp               int64
// 		closedRefreshTimeout  time.Duration
// 		closedRefreshExp      int64
// 		tripped               int32
// 		closedErrRate         float32
// 		halfOpenErrRate       float32
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
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 			b := &breaker{
// 				closedErrShouldTrip:   test.fields.closedErrShouldTrip,
// 				halfOpenErrShouldTrip: test.fields.halfOpenErrShouldTrip,
// 				count:                 test.fields.count,
// 				key:                   test.fields.key,
// 				minSamples:            test.fields.minSamples,
// 				openTimeout:           test.fields.openTimeout,
// 				openExp:               test.fields.openExp,
// 				closedRefreshTimeout:  test.fields.closedRefreshTimeout,
// 				closedRefreshExp:      test.fields.closedRefreshExp,
// 				tripped:               test.fields.tripped,
// 				closedErrRate:         test.fields.closedErrRate,
// 				halfOpenErrRate:       test.fields.halfOpenErrRate,
// 			}
//
// 			b.reset()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_breaker_trip(t *testing.T) {
// 	type fields struct {
// 		closedErrShouldTrip   Tripper
// 		halfOpenErrShouldTrip Tripper
// 		count                 *count
// 		key                   string
// 		minSamples            int64
// 		openTimeout           time.Duration
// 		openExp               int64
// 		closedRefreshTimeout  time.Duration
// 		closedRefreshExp      int64
// 		tripped               int32
// 		closedErrRate         float32
// 		halfOpenErrRate       float32
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
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 			b := &breaker{
// 				closedErrShouldTrip:   test.fields.closedErrShouldTrip,
// 				halfOpenErrShouldTrip: test.fields.halfOpenErrShouldTrip,
// 				count:                 test.fields.count,
// 				key:                   test.fields.key,
// 				minSamples:            test.fields.minSamples,
// 				openTimeout:           test.fields.openTimeout,
// 				openExp:               test.fields.openExp,
// 				closedRefreshTimeout:  test.fields.closedRefreshTimeout,
// 				closedRefreshExp:      test.fields.closedRefreshExp,
// 				tripped:               test.fields.tripped,
// 				closedErrRate:         test.fields.closedErrRate,
// 				halfOpenErrRate:       test.fields.halfOpenErrRate,
// 			}
//
// 			b.trip()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_breaker_isTripped(t *testing.T) {
// 	type fields struct {
// 		closedErrShouldTrip   Tripper
// 		halfOpenErrShouldTrip Tripper
// 		count                 *count
// 		key                   string
// 		minSamples            int64
// 		openTimeout           time.Duration
// 		openExp               int64
// 		closedRefreshTimeout  time.Duration
// 		closedRefreshExp      int64
// 		tripped               int32
// 		closedErrRate         float32
// 		halfOpenErrRate       float32
// 	}
// 	type want struct {
// 		wantOk bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotOk bool) error {
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 		           closedErrShouldTrip:nil,
// 		           halfOpenErrShouldTrip:nil,
// 		           count:count{},
// 		           key:"",
// 		           minSamples:0,
// 		           openTimeout:nil,
// 		           openExp:0,
// 		           closedRefreshTimeout:nil,
// 		           closedRefreshExp:0,
// 		           tripped:0,
// 		           closedErrRate:0,
// 		           halfOpenErrRate:0,
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
// 			b := &breaker{
// 				closedErrShouldTrip:   test.fields.closedErrShouldTrip,
// 				halfOpenErrShouldTrip: test.fields.halfOpenErrShouldTrip,
// 				count:                 test.fields.count,
// 				key:                   test.fields.key,
// 				minSamples:            test.fields.minSamples,
// 				openTimeout:           test.fields.openTimeout,
// 				openExp:               test.fields.openExp,
// 				closedRefreshTimeout:  test.fields.closedRefreshTimeout,
// 				closedRefreshExp:      test.fields.closedRefreshExp,
// 				tripped:               test.fields.tripped,
// 				closedErrRate:         test.fields.closedErrRate,
// 				halfOpenErrRate:       test.fields.halfOpenErrRate,
// 			}
//
// 			gotOk := b.isTripped()
// 			if err := checkFunc(test.want, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// Test_breaker_do_concurrentWrappedErrors calls do concurrently with errors that must be unwrapped.
// The errors.As targets used to be package-level variables, which is a data race under -race.
func Test_breaker_do_concurrentWrappedErrors(t *testing.T) {
	t.Parallel()

	const callers = 32
	base := errors.New("base error")
	tests := []struct {
		wrap func(error) error
		name string
	}{
		{
			name: "mark with success error is unwrapped",
			wrap: errors.NewErrCircuitBreakerMarkWithSuccess,
		},
		{
			name: "ignorable error is unwrapped",
			wrap: errors.NewErrCircuitBreakerIgnorable,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			b, err := newBreaker("do")
			if err != nil {
				t.Fatal(err)
			}
			wrapped := test.wrap(base)

			start := make(chan struct{})
			errs := make(chan error, callers)
			for range callers {
				go func() {
					<-start
					val, _, err := b.do(context.Background(), func(context.Context) (any, error) {
						return nil, wrapped
					})
					if val != nil {
						err = errors.Errorf("unexpected value: %v", val)
					}
					errs <- err
				}()
			}
			close(start)

			timeout := time.After(10 * time.Second)
			for range callers {
				select {
				case err := <-errs:
					// do must return the unwrapped base error, not the marker error around it.
					if !errors.Is(err, base) || errors.Is(err, wrapped) {
						t.Errorf("got error: %v, want: %v", err, base)
					}
				case <-timeout:
					t.Fatal("timed out waiting for do callers")
				}
			}
		})
	}
}

// Test_breaker_trip_keepsClosedRefreshExp guards the trip/currentState race:
// a concurrent currentState that observed tripped==0 just before trip() must not
// read an expired (zeroed) closedRefreshExp and reset() the fresh trip away.
func Test_breaker_trip_keepsClosedRefreshExp(t *testing.T) {
	t.Parallel()
	refresh := time.Now().Add(time.Hour).UnixNano()
	b := &breaker{
		key:              "insertRPC",
		count:            &count{successes: 1, failures: 30},
		openTimeout:      time.Minute,
		closedRefreshExp: refresh,
	}
	b.trip()
	if !b.isTripped() {
		t.Fatalf("trip() did not set tripped")
	}
	if got := b.closedRefreshExp; got != refresh {
		t.Errorf("trip() changed closedRefreshExp: got %d, want %d", got, refresh)
	}
	if st := b.currentState(); st != StateOpen {
		t.Errorf("currentState after trip = %s, want %s", st, StateOpen)
	}
	if total := b.count.Total(); total != 0 {
		t.Errorf("trip() did not reset count: %d", total)
	}
}

func Test_breaker_currentState_closedRefresh(t *testing.T) {
	t.Parallel()
	const refreshTimeout = time.Hour
	future := time.Now().Add(time.Minute).UnixNano()
	tests := []struct {
		name        string
		exp         int64
		wantRefresh bool
	}{
		{name: "zero deadline is refreshed", exp: 0, wantRefresh: true},
		{name: "expired deadline is refreshed", exp: time.Now().Add(-time.Minute).UnixNano(), wantRefresh: true},
		{name: "unexpired deadline is kept", exp: future, wantRefresh: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			b := &breaker{
				key:                  "insertRPC",
				count:                &count{successes: 3, failures: 2},
				closedRefreshTimeout: refreshTimeout,
				closedRefreshExp:     tc.exp,
			}
			before := time.Now().UnixNano()
			if st := b.currentState(); st != StateClosed {
				tt.Fatalf("currentState = %s, want %s", st, StateClosed)
			}
			if b.isTripped() || b.openExp != 0 {
				tt.Errorf("closed refresh touched open state: tripped=%d openExp=%d", b.tripped, b.openExp)
			}
			total := b.count.Total()
			if tc.wantRefresh {
				if total != 0 {
					tt.Errorf("count not reset on refresh: %d", total)
				}
				if b.closedRefreshExp < before+int64(refreshTimeout) {
					tt.Errorf("closedRefreshExp not advanced: %d", b.closedRefreshExp)
				}
				return
			}
			if total != 5 {
				tt.Errorf("count reset before deadline: %d", total)
			}
			if b.closedRefreshExp != tc.exp {
				tt.Errorf("closedRefreshExp changed before deadline: %d", b.closedRefreshExp)
			}
		})
	}
}
