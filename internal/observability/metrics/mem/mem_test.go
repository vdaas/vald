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

// Package mem provides memory metrics functions
package mem

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type want struct {
		want metrics.Metric
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, metrics.Metric) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got metrics.Metric) error {
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt)
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

			got := New()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_memory_Measurement(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		alloc        metrics.Int64Measure
		totalAlloc   metrics.Int64Measure
		sys          metrics.Int64Measure
		mallocs      metrics.Int64Measure
		frees        metrics.Int64Measure
		heapAlloc    metrics.Int64Measure
		heapSys      metrics.Int64Measure
		heapIdle     metrics.Int64Measure
		heapInuse    metrics.Int64Measure
		heapReleased metrics.Int64Measure
		stackInuse   metrics.Int64Measure
		stackSys     metrics.Int64Measure
		pauseTotalMs metrics.Int64Measure
		numGC        metrics.Int64Measure
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
		           alloc: nil,
		           totalAlloc: nil,
		           sys: nil,
		           mallocs: nil,
		           frees: nil,
		           heapAlloc: nil,
		           heapSys: nil,
		           heapIdle: nil,
		           heapInuse: nil,
		           heapReleased: nil,
		           stackInuse: nil,
		           stackSys: nil,
		           pauseTotalMs: nil,
		           numGC: nil,
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
		           alloc: nil,
		           totalAlloc: nil,
		           sys: nil,
		           mallocs: nil,
		           frees: nil,
		           heapAlloc: nil,
		           heapSys: nil,
		           heapIdle: nil,
		           heapInuse: nil,
		           heapReleased: nil,
		           stackInuse: nil,
		           stackSys: nil,
		           pauseTotalMs: nil,
		           numGC: nil,
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
			defer goleak.VerifyNone(tt)
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
			m := &memory{
				alloc:        test.fields.alloc,
				totalAlloc:   test.fields.totalAlloc,
				sys:          test.fields.sys,
				mallocs:      test.fields.mallocs,
				frees:        test.fields.frees,
				heapAlloc:    test.fields.heapAlloc,
				heapSys:      test.fields.heapSys,
				heapIdle:     test.fields.heapIdle,
				heapInuse:    test.fields.heapInuse,
				heapReleased: test.fields.heapReleased,
				stackInuse:   test.fields.stackInuse,
				stackSys:     test.fields.stackSys,
				pauseTotalMs: test.fields.pauseTotalMs,
				numGC:        test.fields.numGC,
			}

			got, err := m.Measurement(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_memory_MeasurementWithTags(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		alloc        metrics.Int64Measure
		totalAlloc   metrics.Int64Measure
		sys          metrics.Int64Measure
		mallocs      metrics.Int64Measure
		frees        metrics.Int64Measure
		heapAlloc    metrics.Int64Measure
		heapSys      metrics.Int64Measure
		heapIdle     metrics.Int64Measure
		heapInuse    metrics.Int64Measure
		heapReleased metrics.Int64Measure
		stackInuse   metrics.Int64Measure
		stackSys     metrics.Int64Measure
		pauseTotalMs metrics.Int64Measure
		numGC        metrics.Int64Measure
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
		           alloc: nil,
		           totalAlloc: nil,
		           sys: nil,
		           mallocs: nil,
		           frees: nil,
		           heapAlloc: nil,
		           heapSys: nil,
		           heapIdle: nil,
		           heapInuse: nil,
		           heapReleased: nil,
		           stackInuse: nil,
		           stackSys: nil,
		           pauseTotalMs: nil,
		           numGC: nil,
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
		           alloc: nil,
		           totalAlloc: nil,
		           sys: nil,
		           mallocs: nil,
		           frees: nil,
		           heapAlloc: nil,
		           heapSys: nil,
		           heapIdle: nil,
		           heapInuse: nil,
		           heapReleased: nil,
		           stackInuse: nil,
		           stackSys: nil,
		           pauseTotalMs: nil,
		           numGC: nil,
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
			defer goleak.VerifyNone(tt)
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
			m := &memory{
				alloc:        test.fields.alloc,
				totalAlloc:   test.fields.totalAlloc,
				sys:          test.fields.sys,
				mallocs:      test.fields.mallocs,
				frees:        test.fields.frees,
				heapAlloc:    test.fields.heapAlloc,
				heapSys:      test.fields.heapSys,
				heapIdle:     test.fields.heapIdle,
				heapInuse:    test.fields.heapInuse,
				heapReleased: test.fields.heapReleased,
				stackInuse:   test.fields.stackInuse,
				stackSys:     test.fields.stackSys,
				pauseTotalMs: test.fields.pauseTotalMs,
				numGC:        test.fields.numGC,
			}

			got, err := m.MeasurementWithTags(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_memory_View(t *testing.T) {
	t.Parallel()
	type fields struct {
		alloc        metrics.Int64Measure
		totalAlloc   metrics.Int64Measure
		sys          metrics.Int64Measure
		mallocs      metrics.Int64Measure
		frees        metrics.Int64Measure
		heapAlloc    metrics.Int64Measure
		heapSys      metrics.Int64Measure
		heapIdle     metrics.Int64Measure
		heapInuse    metrics.Int64Measure
		heapReleased metrics.Int64Measure
		stackInuse   metrics.Int64Measure
		stackSys     metrics.Int64Measure
		pauseTotalMs metrics.Int64Measure
		numGC        metrics.Int64Measure
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
		           alloc: nil,
		           totalAlloc: nil,
		           sys: nil,
		           mallocs: nil,
		           frees: nil,
		           heapAlloc: nil,
		           heapSys: nil,
		           heapIdle: nil,
		           heapInuse: nil,
		           heapReleased: nil,
		           stackInuse: nil,
		           stackSys: nil,
		           pauseTotalMs: nil,
		           numGC: nil,
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
		           alloc: nil,
		           totalAlloc: nil,
		           sys: nil,
		           mallocs: nil,
		           frees: nil,
		           heapAlloc: nil,
		           heapSys: nil,
		           heapIdle: nil,
		           heapInuse: nil,
		           heapReleased: nil,
		           stackInuse: nil,
		           stackSys: nil,
		           pauseTotalMs: nil,
		           numGC: nil,
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
			defer goleak.VerifyNone(tt)
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
			m := &memory{
				alloc:        test.fields.alloc,
				totalAlloc:   test.fields.totalAlloc,
				sys:          test.fields.sys,
				mallocs:      test.fields.mallocs,
				frees:        test.fields.frees,
				heapAlloc:    test.fields.heapAlloc,
				heapSys:      test.fields.heapSys,
				heapIdle:     test.fields.heapIdle,
				heapInuse:    test.fields.heapInuse,
				heapReleased: test.fields.heapReleased,
				stackInuse:   test.fields.stackInuse,
				stackSys:     test.fields.stackSys,
				pauseTotalMs: test.fields.pauseTotalMs,
				numGC:        test.fields.numGC,
			}

			got := m.View()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
