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
package circuitbreaker

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func Test_newBreaker(t *testing.T) {
	type args struct {
		key  string
		opts []BreakerOption
	}
	type want struct {
		want *breaker
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *breaker, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *breaker, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
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
		       args: args {
		           key: "",
		           opts: nil,
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
		           args: args {
		           key: "",
		           opts: nil,
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, err := newBreaker(test.args.key, test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_breaker_do(t *testing.T) {
	type args struct {
		ctx context.Context
		fn  func(ctx context.Context) (val interface{}, err error)
	}
	type fields struct {
		key                   string
		count                 atomic.Value
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
		wantVal interface{}
		wantSt  State
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, State, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotVal interface{}, gotSt State, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotVal, w.wantVal) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVal, w.wantVal)
		}
		if !reflect.DeepEqual(gotSt, w.wantSt) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSt, w.wantSt)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           fn: nil,
		       },
		       fields: fields {
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
		           args: args {
		           ctx: nil,
		           fn: nil,
		           },
		           fields: fields {
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
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

			gotVal, gotSt, err := b.do(test.args.ctx, test.args.fn)
			if err := checkFunc(test.want, gotVal, gotSt, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_breaker_isReady(t *testing.T) {
	type fields struct {
		key                   string
		count                 atomic.Value
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
		wantOk bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotOk bool) error {
		if !reflect.DeepEqual(gotOk, w.wantOk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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

			gotOk := b.isReady()
			if err := checkFunc(test.want, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_breaker_success(t *testing.T) {
	type fields struct {
		key                   string
		count                 atomic.Value
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
		count                 atomic.Value
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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

			b.fail()
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_breaker_currentState(t *testing.T) {
	type fields struct {
		key                   string
		count                 atomic.Value
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
		want State
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, State) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got State) error {
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
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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

			got := b.currentState()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_breaker_reset(t *testing.T) {
	type fields struct {
		key                   string
		count                 atomic.Value
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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

			b.reset()
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_breaker_trip(t *testing.T) {
	type fields struct {
		key                   string
		count                 atomic.Value
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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

			b.trip()
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_breaker_isTripped(t *testing.T) {
	type fields struct {
		key                   string
		count                 atomic.Value
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
		wantOk bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotOk bool) error {
		if !reflect.DeepEqual(gotOk, w.wantOk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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
		           key: "",
		           count: nil,
		           tripped: 0,
		           closedErrRate: 0,
		           closedErrShouldTrip: nil,
		           halfOpenErrRate: 0,
		           halfOpenErrShouldTrip: nil,
		           minSamples: 0,
		           openTimeout: nil,
		           openExp: 0,
		           cloedRefreshTimeout: nil,
		           closedRefreshExp: 0,
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

			gotOk := b.isTripped()
			if err := checkFunc(test.want, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
