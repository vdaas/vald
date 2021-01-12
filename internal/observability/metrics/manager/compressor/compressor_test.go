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

// Package compressor provides functions for compressor stats
package compressor

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/manager/compressor/service"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		c service.Compressor
		r service.Registerer
	}
	type want struct {
		want metrics.Metric
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, metrics.Metric) error
		beforeFunc func(args)
		afterFunc  func(args)
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
		       args: args {
		           c: nil,
		           r: nil,
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
		           c: nil,
		           r: nil,
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

			got := New(test.args.c, test.args.r)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_compressorMetrics_Measurement(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		compressor                  service.Compressor
		registerer                  service.Registerer
		compressorBuffer            metrics.Int64Measure
		compressorTotalRequestedJob metrics.Int64Measure
		compressorTotalCompletedJob metrics.Int64Measure
		registererBuffer            metrics.Int64Measure
		registererTotalRequestedJob metrics.Int64Measure
		registererTotalCompletedJob metrics.Int64Measure
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
		           compressor: nil,
		           registerer: nil,
		           compressorBuffer: nil,
		           compressorTotalRequestedJob: nil,
		           compressorTotalCompletedJob: nil,
		           registererBuffer: nil,
		           registererTotalRequestedJob: nil,
		           registererTotalCompletedJob: nil,
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
		           compressor: nil,
		           registerer: nil,
		           compressorBuffer: nil,
		           compressorTotalRequestedJob: nil,
		           compressorTotalCompletedJob: nil,
		           registererBuffer: nil,
		           registererTotalRequestedJob: nil,
		           registererTotalCompletedJob: nil,
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
			c := &compressorMetrics{
				compressor:                  test.fields.compressor,
				registerer:                  test.fields.registerer,
				compressorBuffer:            test.fields.compressorBuffer,
				compressorTotalRequestedJob: test.fields.compressorTotalRequestedJob,
				compressorTotalCompletedJob: test.fields.compressorTotalCompletedJob,
				registererBuffer:            test.fields.registererBuffer,
				registererTotalRequestedJob: test.fields.registererTotalRequestedJob,
				registererTotalCompletedJob: test.fields.registererTotalCompletedJob,
			}

			got, err := c.Measurement(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_compressorMetrics_MeasurementWithTags(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		compressor                  service.Compressor
		registerer                  service.Registerer
		compressorBuffer            metrics.Int64Measure
		compressorTotalRequestedJob metrics.Int64Measure
		compressorTotalCompletedJob metrics.Int64Measure
		registererBuffer            metrics.Int64Measure
		registererTotalRequestedJob metrics.Int64Measure
		registererTotalCompletedJob metrics.Int64Measure
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
		           compressor: nil,
		           registerer: nil,
		           compressorBuffer: nil,
		           compressorTotalRequestedJob: nil,
		           compressorTotalCompletedJob: nil,
		           registererBuffer: nil,
		           registererTotalRequestedJob: nil,
		           registererTotalCompletedJob: nil,
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
		           compressor: nil,
		           registerer: nil,
		           compressorBuffer: nil,
		           compressorTotalRequestedJob: nil,
		           compressorTotalCompletedJob: nil,
		           registererBuffer: nil,
		           registererTotalRequestedJob: nil,
		           registererTotalCompletedJob: nil,
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
			c := &compressorMetrics{
				compressor:                  test.fields.compressor,
				registerer:                  test.fields.registerer,
				compressorBuffer:            test.fields.compressorBuffer,
				compressorTotalRequestedJob: test.fields.compressorTotalRequestedJob,
				compressorTotalCompletedJob: test.fields.compressorTotalCompletedJob,
				registererBuffer:            test.fields.registererBuffer,
				registererTotalRequestedJob: test.fields.registererTotalRequestedJob,
				registererTotalCompletedJob: test.fields.registererTotalCompletedJob,
			}

			got, err := c.MeasurementWithTags(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_compressorMetrics_View(t *testing.T) {
	t.Parallel()
	type fields struct {
		compressor                  service.Compressor
		registerer                  service.Registerer
		compressorBuffer            metrics.Int64Measure
		compressorTotalRequestedJob metrics.Int64Measure
		compressorTotalCompletedJob metrics.Int64Measure
		registererBuffer            metrics.Int64Measure
		registererTotalRequestedJob metrics.Int64Measure
		registererTotalCompletedJob metrics.Int64Measure
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
		           compressor: nil,
		           registerer: nil,
		           compressorBuffer: nil,
		           compressorTotalRequestedJob: nil,
		           compressorTotalCompletedJob: nil,
		           registererBuffer: nil,
		           registererTotalRequestedJob: nil,
		           registererTotalCompletedJob: nil,
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
		           compressor: nil,
		           registerer: nil,
		           compressorBuffer: nil,
		           compressorTotalRequestedJob: nil,
		           compressorTotalCompletedJob: nil,
		           registererBuffer: nil,
		           registererTotalRequestedJob: nil,
		           registererTotalCompletedJob: nil,
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
			c := &compressorMetrics{
				compressor:                  test.fields.compressor,
				registerer:                  test.fields.registerer,
				compressorBuffer:            test.fields.compressorBuffer,
				compressorTotalRequestedJob: test.fields.compressorTotalRequestedJob,
				compressorTotalCompletedJob: test.fields.compressorTotalCompletedJob,
				registererBuffer:            test.fields.registererBuffer,
				registererTotalRequestedJob: test.fields.registererTotalRequestedJob,
				registererTotalCompletedJob: test.fields.registererTotalCompletedJob,
			}

			got := c.View()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
