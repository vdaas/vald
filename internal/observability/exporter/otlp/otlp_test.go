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
package otlp

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/exporter"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/test/goleak"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

// NOT IMPLEMENTED BELOW

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want exporter.Exporter
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, exporter.Exporter, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got exporter.Exporter, err error) error {
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
		           opts:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           opts:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_exp_initTracer(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		serviceName         string
		collectorEndpoint   string
		traceExporter       *otlptrace.Exporter
		traceProvider       *trace.TracerProvider
		tBatchTimeout       time.Duration
		tExportTimeout      time.Duration
		tMaxExportBatchSize int
		tMaxQueueSize       int
		metricsExporter     metric.Exporter
		meterProvider       *metric.MeterProvider
		metricsViews        []metrics.View
		mExportInterval     time.Duration
		mExportTimeout      time.Duration
		attributes          []attribute.KeyValue
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx:nil,
		       },
		       fields: fields {
		           serviceName:"",
		           collectorEndpoint:"",
		           traceExporter:nil,
		           traceProvider:nil,
		           tBatchTimeout:nil,
		           tExportTimeout:nil,
		           tMaxExportBatchSize:0,
		           tMaxQueueSize:0,
		           metricsExporter:nil,
		           meterProvider:nil,
		           metricsViews:nil,
		           mExportInterval:nil,
		           mExportTimeout:nil,
		           attributes:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           serviceName:"",
		           collectorEndpoint:"",
		           traceExporter:nil,
		           traceProvider:nil,
		           tBatchTimeout:nil,
		           tExportTimeout:nil,
		           tMaxExportBatchSize:0,
		           tMaxQueueSize:0,
		           metricsExporter:nil,
		           meterProvider:nil,
		           metricsViews:nil,
		           mExportInterval:nil,
		           mExportTimeout:nil,
		           attributes:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			e := &exp{
				serviceName:         test.fields.serviceName,
				collectorEndpoint:   test.fields.collectorEndpoint,
				traceExporter:       test.fields.traceExporter,
				traceProvider:       test.fields.traceProvider,
				tBatchTimeout:       test.fields.tBatchTimeout,
				tExportTimeout:      test.fields.tExportTimeout,
				tMaxExportBatchSize: test.fields.tMaxExportBatchSize,
				tMaxQueueSize:       test.fields.tMaxQueueSize,
				metricsExporter:     test.fields.metricsExporter,
				meterProvider:       test.fields.meterProvider,
				metricsViews:        test.fields.metricsViews,
				mExportInterval:     test.fields.mExportInterval,
				mExportTimeout:      test.fields.mExportTimeout,
				attributes:          test.fields.attributes,
			}

			err := e.initTracer(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_exp_initMeter(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		serviceName         string
		collectorEndpoint   string
		traceExporter       *otlptrace.Exporter
		traceProvider       *trace.TracerProvider
		tBatchTimeout       time.Duration
		tExportTimeout      time.Duration
		tMaxExportBatchSize int
		tMaxQueueSize       int
		metricsExporter     metric.Exporter
		meterProvider       *metric.MeterProvider
		metricsViews        []metrics.View
		mExportInterval     time.Duration
		mExportTimeout      time.Duration
		attributes          []attribute.KeyValue
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx:nil,
		       },
		       fields: fields {
		           serviceName:"",
		           collectorEndpoint:"",
		           traceExporter:nil,
		           traceProvider:nil,
		           tBatchTimeout:nil,
		           tExportTimeout:nil,
		           tMaxExportBatchSize:0,
		           tMaxQueueSize:0,
		           metricsExporter:nil,
		           meterProvider:nil,
		           metricsViews:nil,
		           mExportInterval:nil,
		           mExportTimeout:nil,
		           attributes:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           serviceName:"",
		           collectorEndpoint:"",
		           traceExporter:nil,
		           traceProvider:nil,
		           tBatchTimeout:nil,
		           tExportTimeout:nil,
		           tMaxExportBatchSize:0,
		           tMaxQueueSize:0,
		           metricsExporter:nil,
		           meterProvider:nil,
		           metricsViews:nil,
		           mExportInterval:nil,
		           mExportTimeout:nil,
		           attributes:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			e := &exp{
				serviceName:         test.fields.serviceName,
				collectorEndpoint:   test.fields.collectorEndpoint,
				traceExporter:       test.fields.traceExporter,
				traceProvider:       test.fields.traceProvider,
				tBatchTimeout:       test.fields.tBatchTimeout,
				tExportTimeout:      test.fields.tExportTimeout,
				tMaxExportBatchSize: test.fields.tMaxExportBatchSize,
				tMaxQueueSize:       test.fields.tMaxQueueSize,
				metricsExporter:     test.fields.metricsExporter,
				meterProvider:       test.fields.meterProvider,
				metricsViews:        test.fields.metricsViews,
				mExportInterval:     test.fields.mExportInterval,
				mExportTimeout:      test.fields.mExportTimeout,
				attributes:          test.fields.attributes,
			}

			err := e.initMeter(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_exp_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		serviceName         string
		collectorEndpoint   string
		traceExporter       *otlptrace.Exporter
		traceProvider       *trace.TracerProvider
		tBatchTimeout       time.Duration
		tExportTimeout      time.Duration
		tMaxExportBatchSize int
		tMaxQueueSize       int
		metricsExporter     metric.Exporter
		meterProvider       *metric.MeterProvider
		metricsViews        []metrics.View
		mExportInterval     time.Duration
		mExportTimeout      time.Duration
		attributes          []attribute.KeyValue
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx:nil,
		       },
		       fields: fields {
		           serviceName:"",
		           collectorEndpoint:"",
		           traceExporter:nil,
		           traceProvider:nil,
		           tBatchTimeout:nil,
		           tExportTimeout:nil,
		           tMaxExportBatchSize:0,
		           tMaxQueueSize:0,
		           metricsExporter:nil,
		           meterProvider:nil,
		           metricsViews:nil,
		           mExportInterval:nil,
		           mExportTimeout:nil,
		           attributes:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           serviceName:"",
		           collectorEndpoint:"",
		           traceExporter:nil,
		           traceProvider:nil,
		           tBatchTimeout:nil,
		           tExportTimeout:nil,
		           tMaxExportBatchSize:0,
		           tMaxQueueSize:0,
		           metricsExporter:nil,
		           meterProvider:nil,
		           metricsViews:nil,
		           mExportInterval:nil,
		           mExportTimeout:nil,
		           attributes:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			e := &exp{
				serviceName:         test.fields.serviceName,
				collectorEndpoint:   test.fields.collectorEndpoint,
				traceExporter:       test.fields.traceExporter,
				traceProvider:       test.fields.traceProvider,
				tBatchTimeout:       test.fields.tBatchTimeout,
				tExportTimeout:      test.fields.tExportTimeout,
				tMaxExportBatchSize: test.fields.tMaxExportBatchSize,
				tMaxQueueSize:       test.fields.tMaxQueueSize,
				metricsExporter:     test.fields.metricsExporter,
				meterProvider:       test.fields.meterProvider,
				metricsViews:        test.fields.metricsViews,
				mExportInterval:     test.fields.mExportInterval,
				mExportTimeout:      test.fields.mExportTimeout,
				attributes:          test.fields.attributes,
			}

			err := e.Start(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_exp_Stop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		serviceName         string
		collectorEndpoint   string
		traceExporter       *otlptrace.Exporter
		traceProvider       *trace.TracerProvider
		tBatchTimeout       time.Duration
		tExportTimeout      time.Duration
		tMaxExportBatchSize int
		tMaxQueueSize       int
		metricsExporter     metric.Exporter
		meterProvider       *metric.MeterProvider
		metricsViews        []metrics.View
		mExportInterval     time.Duration
		mExportTimeout      time.Duration
		attributes          []attribute.KeyValue
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx:nil,
		       },
		       fields: fields {
		           serviceName:"",
		           collectorEndpoint:"",
		           traceExporter:nil,
		           traceProvider:nil,
		           tBatchTimeout:nil,
		           tExportTimeout:nil,
		           tMaxExportBatchSize:0,
		           tMaxQueueSize:0,
		           metricsExporter:nil,
		           meterProvider:nil,
		           metricsViews:nil,
		           mExportInterval:nil,
		           mExportTimeout:nil,
		           attributes:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T, args args) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx:nil,
		           },
		           fields: fields {
		           serviceName:"",
		           collectorEndpoint:"",
		           traceExporter:nil,
		           traceProvider:nil,
		           tBatchTimeout:nil,
		           tExportTimeout:nil,
		           tMaxExportBatchSize:0,
		           tMaxQueueSize:0,
		           metricsExporter:nil,
		           meterProvider:nil,
		           metricsViews:nil,
		           mExportInterval:nil,
		           mExportTimeout:nil,
		           attributes:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T, args args) {
		               t.Helper()
		           },
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			e := &exp{
				serviceName:         test.fields.serviceName,
				collectorEndpoint:   test.fields.collectorEndpoint,
				traceExporter:       test.fields.traceExporter,
				traceProvider:       test.fields.traceProvider,
				tBatchTimeout:       test.fields.tBatchTimeout,
				tExportTimeout:      test.fields.tExportTimeout,
				tMaxExportBatchSize: test.fields.tMaxExportBatchSize,
				tMaxQueueSize:       test.fields.tMaxQueueSize,
				metricsExporter:     test.fields.metricsExporter,
				meterProvider:       test.fields.meterProvider,
				metricsViews:        test.fields.metricsViews,
				mExportInterval:     test.fields.mExportInterval,
				mExportTimeout:      test.fields.mExportTimeout,
				attributes:          test.fields.attributes,
			}

			err := e.Stop(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
