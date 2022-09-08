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

// Package sidecar provides functions for sidecar stats
package sidecar

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/sidecar/service/observer"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type want struct {
		want MetricsHook
		err  error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, MetricsHook, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got MetricsHook, err error) error {
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

			got, err := New()
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_sidecarMetrics_Measurement(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		uploadTotal    metrics.Int64Measure
		uploadBytes    metrics.Int64Measure
		uploadLatency  metrics.Float64Measure
		storageTypeKey metrics.Key
		bucketNameKey  metrics.Key
		filenameKey    metrics.Key
		mu             sync.Mutex
		ms             []metrics.MeasurementWithTags
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
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
			sm := &sidecarMetrics{
				uploadTotal:    test.fields.uploadTotal,
				uploadBytes:    test.fields.uploadBytes,
				uploadLatency:  test.fields.uploadLatency,
				storageTypeKey: test.fields.storageTypeKey,
				bucketNameKey:  test.fields.bucketNameKey,
				filenameKey:    test.fields.filenameKey,
				mu:             test.fields.mu,
				ms:             test.fields.ms,
			}

			got, err := sm.Measurement(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_sidecarMetrics_MeasurementWithTags(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		uploadTotal    metrics.Int64Measure
		uploadBytes    metrics.Int64Measure
		uploadLatency  metrics.Float64Measure
		storageTypeKey metrics.Key
		bucketNameKey  metrics.Key
		filenameKey    metrics.Key
		mu             sync.Mutex
		ms             []metrics.MeasurementWithTags
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
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
			sm := &sidecarMetrics{
				uploadTotal:    test.fields.uploadTotal,
				uploadBytes:    test.fields.uploadBytes,
				uploadLatency:  test.fields.uploadLatency,
				storageTypeKey: test.fields.storageTypeKey,
				bucketNameKey:  test.fields.bucketNameKey,
				filenameKey:    test.fields.filenameKey,
				mu:             test.fields.mu,
				ms:             test.fields.ms,
			}

			got, err := sm.MeasurementWithTags(test.args.ctx)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_sidecarMetrics_View(t *testing.T) {
	t.Parallel()
	type fields struct {
		uploadTotal    metrics.Int64Measure
		uploadBytes    metrics.Int64Measure
		uploadLatency  metrics.Float64Measure
		storageTypeKey metrics.Key
		bucketNameKey  metrics.Key
		filenameKey    metrics.Key
		mu             sync.Mutex
		ms             []metrics.MeasurementWithTags
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
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
			sm := &sidecarMetrics{
				uploadTotal:    test.fields.uploadTotal,
				uploadBytes:    test.fields.uploadBytes,
				uploadLatency:  test.fields.uploadLatency,
				storageTypeKey: test.fields.storageTypeKey,
				bucketNameKey:  test.fields.bucketNameKey,
				filenameKey:    test.fields.filenameKey,
				mu:             test.fields.mu,
				ms:             test.fields.ms,
			}

			got := sm.View()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_sidecarMetrics_BeforeProcess(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		info *observer.BackupInfo
	}
	type fields struct {
		uploadTotal    metrics.Int64Measure
		uploadBytes    metrics.Int64Measure
		uploadLatency  metrics.Float64Measure
		storageTypeKey metrics.Key
		bucketNameKey  metrics.Key
		filenameKey    metrics.Key
		mu             sync.Mutex
		ms             []metrics.MeasurementWithTags
	}
	type want struct {
		want context.Context
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, context.Context, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got context.Context, err error) error {
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
		           info: nil,
		       },
		       fields: fields {
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
		           info: nil,
		           },
		           fields: fields {
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
			sm := &sidecarMetrics{
				uploadTotal:    test.fields.uploadTotal,
				uploadBytes:    test.fields.uploadBytes,
				uploadLatency:  test.fields.uploadLatency,
				storageTypeKey: test.fields.storageTypeKey,
				bucketNameKey:  test.fields.bucketNameKey,
				filenameKey:    test.fields.filenameKey,
				mu:             test.fields.mu,
				ms:             test.fields.ms,
			}

			got, err := sm.BeforeProcess(test.args.ctx, test.args.info)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_sidecarMetrics_AfterProcess(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		info *observer.BackupInfo
	}
	type fields struct {
		uploadTotal    metrics.Int64Measure
		uploadBytes    metrics.Int64Measure
		uploadLatency  metrics.Float64Measure
		storageTypeKey metrics.Key
		bucketNameKey  metrics.Key
		filenameKey    metrics.Key
		mu             sync.Mutex
		ms             []metrics.MeasurementWithTags
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
		beforeFunc func(args)
		afterFunc  func(args)
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
		           ctx: nil,
		           info: nil,
		       },
		       fields: fields {
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
		           info: nil,
		           },
		           fields: fields {
		           uploadTotal: nil,
		           uploadBytes: nil,
		           uploadLatency: nil,
		           storageTypeKey: nil,
		           bucketNameKey: nil,
		           filenameKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
			sm := &sidecarMetrics{
				uploadTotal:    test.fields.uploadTotal,
				uploadBytes:    test.fields.uploadBytes,
				uploadLatency:  test.fields.uploadLatency,
				storageTypeKey: test.fields.storageTypeKey,
				bucketNameKey:  test.fields.bucketNameKey,
				filenameKey:    test.fields.filenameKey,
				mu:             test.fields.mu,
				ms:             test.fields.ms,
			}

			err := sm.AfterProcess(test.args.ctx, test.args.info)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
