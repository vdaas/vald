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
package otlp

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want exporter.Exporter
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, exporter.Exporter, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got exporter.Exporter, err error) error {
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
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_exp_initTracer(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		metricsExporter     metric.Exporter
// 		meterProvider       *metric.MeterProvider
// 		traceExporter       *otlptrace.Exporter
// 		traceProvider       *trace.TracerProvider
// 		collectorEndpoint   string
// 		serviceName         string
// 		metricsViews        []metrics.View
// 		attributes          []attribute.KeyValue
// 		tBatchTimeout       time.Duration
// 		tExportTimeout      time.Duration
// 		tMaxExportBatchSize int
// 		tMaxQueueSize       int
// 		mExportInterval     time.Duration
// 		mExportTimeout      time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
// 		       },
// 		       fields: fields {
// 		           metricsExporter:nil,
// 		           meterProvider:nil,
// 		           traceExporter:nil,
// 		           traceProvider:nil,
// 		           collectorEndpoint:"",
// 		           serviceName:"",
// 		           metricsViews:nil,
// 		           attributes:nil,
// 		           tBatchTimeout:nil,
// 		           tExportTimeout:nil,
// 		           tMaxExportBatchSize:0,
// 		           tMaxQueueSize:0,
// 		           mExportInterval:nil,
// 		           mExportTimeout:nil,
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
// 		           },
// 		           fields: fields {
// 		           metricsExporter:nil,
// 		           meterProvider:nil,
// 		           traceExporter:nil,
// 		           traceProvider:nil,
// 		           collectorEndpoint:"",
// 		           serviceName:"",
// 		           metricsViews:nil,
// 		           attributes:nil,
// 		           tBatchTimeout:nil,
// 		           tExportTimeout:nil,
// 		           tMaxExportBatchSize:0,
// 		           tMaxQueueSize:0,
// 		           mExportInterval:nil,
// 		           mExportTimeout:nil,
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
// 			e := &exp{
// 				metricsExporter:     test.fields.metricsExporter,
// 				meterProvider:       test.fields.meterProvider,
// 				traceExporter:       test.fields.traceExporter,
// 				traceProvider:       test.fields.traceProvider,
// 				collectorEndpoint:   test.fields.collectorEndpoint,
// 				serviceName:         test.fields.serviceName,
// 				metricsViews:        test.fields.metricsViews,
// 				attributes:          test.fields.attributes,
// 				tBatchTimeout:       test.fields.tBatchTimeout,
// 				tExportTimeout:      test.fields.tExportTimeout,
// 				tMaxExportBatchSize: test.fields.tMaxExportBatchSize,
// 				tMaxQueueSize:       test.fields.tMaxQueueSize,
// 				mExportInterval:     test.fields.mExportInterval,
// 				mExportTimeout:      test.fields.mExportTimeout,
// 			}
//
// 			err := e.initTracer(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_exp_initMeter(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		metricsExporter     metric.Exporter
// 		meterProvider       *metric.MeterProvider
// 		traceExporter       *otlptrace.Exporter
// 		traceProvider       *trace.TracerProvider
// 		collectorEndpoint   string
// 		serviceName         string
// 		metricsViews        []metrics.View
// 		attributes          []attribute.KeyValue
// 		tBatchTimeout       time.Duration
// 		tExportTimeout      time.Duration
// 		tMaxExportBatchSize int
// 		tMaxQueueSize       int
// 		mExportInterval     time.Duration
// 		mExportTimeout      time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
// 		       },
// 		       fields: fields {
// 		           metricsExporter:nil,
// 		           meterProvider:nil,
// 		           traceExporter:nil,
// 		           traceProvider:nil,
// 		           collectorEndpoint:"",
// 		           serviceName:"",
// 		           metricsViews:nil,
// 		           attributes:nil,
// 		           tBatchTimeout:nil,
// 		           tExportTimeout:nil,
// 		           tMaxExportBatchSize:0,
// 		           tMaxQueueSize:0,
// 		           mExportInterval:nil,
// 		           mExportTimeout:nil,
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
// 		           },
// 		           fields: fields {
// 		           metricsExporter:nil,
// 		           meterProvider:nil,
// 		           traceExporter:nil,
// 		           traceProvider:nil,
// 		           collectorEndpoint:"",
// 		           serviceName:"",
// 		           metricsViews:nil,
// 		           attributes:nil,
// 		           tBatchTimeout:nil,
// 		           tExportTimeout:nil,
// 		           tMaxExportBatchSize:0,
// 		           tMaxQueueSize:0,
// 		           mExportInterval:nil,
// 		           mExportTimeout:nil,
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
// 			e := &exp{
// 				metricsExporter:     test.fields.metricsExporter,
// 				meterProvider:       test.fields.meterProvider,
// 				traceExporter:       test.fields.traceExporter,
// 				traceProvider:       test.fields.traceProvider,
// 				collectorEndpoint:   test.fields.collectorEndpoint,
// 				serviceName:         test.fields.serviceName,
// 				metricsViews:        test.fields.metricsViews,
// 				attributes:          test.fields.attributes,
// 				tBatchTimeout:       test.fields.tBatchTimeout,
// 				tExportTimeout:      test.fields.tExportTimeout,
// 				tMaxExportBatchSize: test.fields.tMaxExportBatchSize,
// 				tMaxQueueSize:       test.fields.tMaxQueueSize,
// 				mExportInterval:     test.fields.mExportInterval,
// 				mExportTimeout:      test.fields.mExportTimeout,
// 			}
//
// 			err := e.initMeter(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_exp_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		metricsExporter     metric.Exporter
// 		meterProvider       *metric.MeterProvider
// 		traceExporter       *otlptrace.Exporter
// 		traceProvider       *trace.TracerProvider
// 		collectorEndpoint   string
// 		serviceName         string
// 		metricsViews        []metrics.View
// 		attributes          []attribute.KeyValue
// 		tBatchTimeout       time.Duration
// 		tExportTimeout      time.Duration
// 		tMaxExportBatchSize int
// 		tMaxQueueSize       int
// 		mExportInterval     time.Duration
// 		mExportTimeout      time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
// 		       },
// 		       fields: fields {
// 		           metricsExporter:nil,
// 		           meterProvider:nil,
// 		           traceExporter:nil,
// 		           traceProvider:nil,
// 		           collectorEndpoint:"",
// 		           serviceName:"",
// 		           metricsViews:nil,
// 		           attributes:nil,
// 		           tBatchTimeout:nil,
// 		           tExportTimeout:nil,
// 		           tMaxExportBatchSize:0,
// 		           tMaxQueueSize:0,
// 		           mExportInterval:nil,
// 		           mExportTimeout:nil,
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
// 		           },
// 		           fields: fields {
// 		           metricsExporter:nil,
// 		           meterProvider:nil,
// 		           traceExporter:nil,
// 		           traceProvider:nil,
// 		           collectorEndpoint:"",
// 		           serviceName:"",
// 		           metricsViews:nil,
// 		           attributes:nil,
// 		           tBatchTimeout:nil,
// 		           tExportTimeout:nil,
// 		           tMaxExportBatchSize:0,
// 		           tMaxQueueSize:0,
// 		           mExportInterval:nil,
// 		           mExportTimeout:nil,
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
// 			e := &exp{
// 				metricsExporter:     test.fields.metricsExporter,
// 				meterProvider:       test.fields.meterProvider,
// 				traceExporter:       test.fields.traceExporter,
// 				traceProvider:       test.fields.traceProvider,
// 				collectorEndpoint:   test.fields.collectorEndpoint,
// 				serviceName:         test.fields.serviceName,
// 				metricsViews:        test.fields.metricsViews,
// 				attributes:          test.fields.attributes,
// 				tBatchTimeout:       test.fields.tBatchTimeout,
// 				tExportTimeout:      test.fields.tExportTimeout,
// 				tMaxExportBatchSize: test.fields.tMaxExportBatchSize,
// 				tMaxQueueSize:       test.fields.tMaxQueueSize,
// 				mExportInterval:     test.fields.mExportInterval,
// 				mExportTimeout:      test.fields.mExportTimeout,
// 			}
//
// 			err := e.Start(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_exp_Stop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		metricsExporter     metric.Exporter
// 		meterProvider       *metric.MeterProvider
// 		traceExporter       *otlptrace.Exporter
// 		traceProvider       *trace.TracerProvider
// 		collectorEndpoint   string
// 		serviceName         string
// 		metricsViews        []metrics.View
// 		attributes          []attribute.KeyValue
// 		tBatchTimeout       time.Duration
// 		tExportTimeout      time.Duration
// 		tMaxExportBatchSize int
// 		tMaxQueueSize       int
// 		mExportInterval     time.Duration
// 		mExportTimeout      time.Duration
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
// 		       },
// 		       fields: fields {
// 		           metricsExporter:nil,
// 		           meterProvider:nil,
// 		           traceExporter:nil,
// 		           traceProvider:nil,
// 		           collectorEndpoint:"",
// 		           serviceName:"",
// 		           metricsViews:nil,
// 		           attributes:nil,
// 		           tBatchTimeout:nil,
// 		           tExportTimeout:nil,
// 		           tMaxExportBatchSize:0,
// 		           tMaxQueueSize:0,
// 		           mExportInterval:nil,
// 		           mExportTimeout:nil,
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
// 		           },
// 		           fields: fields {
// 		           metricsExporter:nil,
// 		           meterProvider:nil,
// 		           traceExporter:nil,
// 		           traceProvider:nil,
// 		           collectorEndpoint:"",
// 		           serviceName:"",
// 		           metricsViews:nil,
// 		           attributes:nil,
// 		           tBatchTimeout:nil,
// 		           tExportTimeout:nil,
// 		           tMaxExportBatchSize:0,
// 		           tMaxQueueSize:0,
// 		           mExportInterval:nil,
// 		           mExportTimeout:nil,
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
// 			e := &exp{
// 				metricsExporter:     test.fields.metricsExporter,
// 				meterProvider:       test.fields.meterProvider,
// 				traceExporter:       test.fields.traceExporter,
// 				traceProvider:       test.fields.traceProvider,
// 				collectorEndpoint:   test.fields.collectorEndpoint,
// 				serviceName:         test.fields.serviceName,
// 				metricsViews:        test.fields.metricsViews,
// 				attributes:          test.fields.attributes,
// 				tBatchTimeout:       test.fields.tBatchTimeout,
// 				tExportTimeout:      test.fields.tExportTimeout,
// 				tMaxExportBatchSize: test.fields.tMaxExportBatchSize,
// 				tMaxQueueSize:       test.fields.tMaxQueueSize,
// 				mExportInterval:     test.fields.mExportInterval,
// 				mExportTimeout:      test.fields.mExportTimeout,
// 			}
//
// 			err := e.Stop(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
