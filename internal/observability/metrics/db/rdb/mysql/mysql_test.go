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

// Package mysql provides mysql metrics functions
package mysql

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/db/rdb/mysql"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type want struct {
		wantE EventReceiver
		err   error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, EventReceiver, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotE EventReceiver, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotE, w.wantE) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotE, w.wantE)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotE, err := New()
			if err := test.checkFunc(test.want, gotE, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_mysqlMetrics_Measurement(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		queryTotal        metrics.Int64Measure
		queryLatency      metrics.Float64Measure
		mu                sync.Mutex
		ms                []metrics.Measurement
		NullEventReceiver mysql.NullEventReceiver
	}
	type want struct {
		want []metrics.Measurement
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []metrics.Measurement, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []metrics.Measurement, err error) error {
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
		           ctx: nil,
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
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
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			mm := &mysqlMetrics{
				queryTotal:        test.fields.queryTotal,
				queryLatency:      test.fields.queryLatency,
				mu:                test.fields.mu,
				ms:                test.fields.ms,
				NullEventReceiver: test.fields.NullEventReceiver,
			}

			got, err := mm.Measurement(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_mysqlMetrics_MeasurementWithTags(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		queryTotal        metrics.Int64Measure
		queryLatency      metrics.Float64Measure
		mu                sync.Mutex
		ms                []metrics.Measurement
		NullEventReceiver mysql.NullEventReceiver
	}
	type want struct {
		want []metrics.MeasurementWithTags
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []metrics.MeasurementWithTags, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []metrics.MeasurementWithTags, err error) error {
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
		           ctx: nil,
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
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
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			mm := &mysqlMetrics{
				queryTotal:        test.fields.queryTotal,
				queryLatency:      test.fields.queryLatency,
				mu:                test.fields.mu,
				ms:                test.fields.ms,
				NullEventReceiver: test.fields.NullEventReceiver,
			}

			got, err := mm.MeasurementWithTags(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_mysqlMetrics_View(t *testing.T) {
	t.Parallel()
	type fields struct {
		queryTotal        metrics.Int64Measure
		queryLatency      metrics.Float64Measure
		mu                sync.Mutex
		ms                []metrics.Measurement
		NullEventReceiver mysql.NullEventReceiver
	}
	type want struct {
		want []*metrics.View
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []*metrics.View) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []*metrics.View) error {
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
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
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
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			mm := &mysqlMetrics{
				queryTotal:        test.fields.queryTotal,
				queryLatency:      test.fields.queryLatency,
				mu:                test.fields.mu,
				ms:                test.fields.ms,
				NullEventReceiver: test.fields.NullEventReceiver,
			}

			got := mm.View()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_mysqlMetrics_SpanStart(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx       context.Context
		eventName string
		query     string
	}
	type fields struct {
		queryTotal        metrics.Int64Measure
		queryLatency      metrics.Float64Measure
		mu                sync.Mutex
		ms                []metrics.Measurement
		NullEventReceiver mysql.NullEventReceiver
	}
	type want struct {
		want context.Context
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, context.Context) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got context.Context) error {
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
		           ctx: nil,
		           eventName: "",
		           query: "",
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
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
		           eventName: "",
		           query: "",
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			mm := &mysqlMetrics{
				queryTotal:        test.fields.queryTotal,
				queryLatency:      test.fields.queryLatency,
				mu:                test.fields.mu,
				ms:                test.fields.ms,
				NullEventReceiver: test.fields.NullEventReceiver,
			}

			got := mm.SpanStart(test.args.ctx, test.args.eventName, test.args.query)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_mysqlMetrics_SpanError(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		err error
	}
	type fields struct {
		queryTotal        metrics.Int64Measure
		queryLatency      metrics.Float64Measure
		mu                sync.Mutex
		ms                []metrics.Measurement
		NullEventReceiver mysql.NullEventReceiver
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           err: nil,
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
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
		           err: nil,
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			mm := &mysqlMetrics{
				queryTotal:        test.fields.queryTotal,
				queryLatency:      test.fields.queryLatency,
				mu:                test.fields.mu,
				ms:                test.fields.ms,
				NullEventReceiver: test.fields.NullEventReceiver,
			}

			mm.SpanError(test.args.ctx, test.args.err)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_mysqlMetrics_SpanFinish(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		queryTotal        metrics.Int64Measure
		queryLatency      metrics.Float64Measure
		mu                sync.Mutex
		ms                []metrics.Measurement
		NullEventReceiver mysql.NullEventReceiver
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
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
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           mu: sync.Mutex{},
		           ms: nil,
		           NullEventReceiver: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			mm := &mysqlMetrics{
				queryTotal:        test.fields.queryTotal,
				queryLatency:      test.fields.queryLatency,
				mu:                test.fields.mu,
				ms:                test.fields.ms,
				NullEventReceiver: test.fields.NullEventReceiver,
			}

			mm.SpanFinish(test.args.ctx)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
