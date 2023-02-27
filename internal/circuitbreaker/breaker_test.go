// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func Test_breaker_isReady(t *testing.T) {
	type fields struct {
		key                   string
		count                 *count
		tripped               int32
		closedErrRate         float32
		closedErrShouldTrip   Tripper
		halfOpenErrRate       float32
		halfOpenErrShouldTrip Tripper
		minSamples            int64
		openTimeout           time.Duration
		openExp               int64
		cloedRefreshTimeout   time.Duration
		closedRefreshExp      int64
	}
	type want struct {
		wantSt State
		err    error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, State, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
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
				cloedRefreshTimeout:   test.fields.cloedRefreshTimeout,
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
		key                   string
		count                 *count
		tripped               int32
		closedErrRate         float32
		closedErrShouldTrip   Tripper
		halfOpenErrRate       float32
		halfOpenErrShouldTrip Tripper
		minSamples            int64
		openTimeout           time.Duration
		openExp               int64
		cloedRefreshTimeout   time.Duration
		closedRefreshExp      int64
	}
	type want struct{}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T, *breaker)
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
				cloedRefreshTimeout:   test.fields.cloedRefreshTimeout,
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
		key                   string
		count                 *count
		tripped               int32
		closedErrRate         float32
		closedErrShouldTrip   Tripper
		halfOpenErrRate       float32
		halfOpenErrShouldTrip Tripper
		minSamples            int64
		openTimeout           time.Duration
		openExp               int64
		cloedRefreshTimeout   time.Duration
		closedRefreshExp      int64
	}
	type want struct{}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T, *breaker)
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
						t.Errorf("count reseted: %d", total)
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
				cloedRefreshTimeout:   test.fields.cloedRefreshTimeout,
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
