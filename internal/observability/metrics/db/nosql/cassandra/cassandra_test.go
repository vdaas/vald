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

// Package cassandra provides cassandra metrics functions
package cassandra

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type want struct {
		wantO Observer
		err   error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, Observer, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotO Observer, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotO, w.wantO) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotO, w.wantO)
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

			gotO, err := New()
			if err := test.checkFunc(test.want, gotO, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_cassandraMetrics_Measurement(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		queryTotal         metrics.Int64Measure
		queryAttemptsTotal metrics.Int64Measure
		queryLatency       metrics.Float64Measure
		keyspaceKey        metrics.Key
		clusterNameKey     metrics.Key
		dataCenterKey      metrics.Key
		hostIDKey          metrics.Key
		hostPortKey        metrics.Key
		rackKey            metrics.Key
		versionKey         metrics.Key
		mu                 sync.Mutex
		ms                 []metrics.MeasurementWithTags
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
		           queryAttemptsTotal: nil,
		           queryLatency: nil,
		           keyspaceKey: nil,
		           clusterNameKey: nil,
		           dataCenterKey: nil,
		           hostIDKey: nil,
		           hostPortKey: nil,
		           rackKey: nil,
		           versionKey: nil,
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
		           queryTotal: nil,
		           queryAttemptsTotal: nil,
		           queryLatency: nil,
		           keyspaceKey: nil,
		           clusterNameKey: nil,
		           dataCenterKey: nil,
		           hostIDKey: nil,
		           hostPortKey: nil,
		           rackKey: nil,
		           versionKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
			cm := &cassandraMetrics{
				queryTotal:         test.fields.queryTotal,
				queryAttemptsTotal: test.fields.queryAttemptsTotal,
				queryLatency:       test.fields.queryLatency,
				keyspaceKey:        test.fields.keyspaceKey,
				clusterNameKey:     test.fields.clusterNameKey,
				dataCenterKey:      test.fields.dataCenterKey,
				hostIDKey:          test.fields.hostIDKey,
				hostPortKey:        test.fields.hostPortKey,
				rackKey:            test.fields.rackKey,
				versionKey:         test.fields.versionKey,
				mu:                 test.fields.mu,
				ms:                 test.fields.ms,
			}

			got, err := cm.Measurement(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_cassandraMetrics_MeasurementWithTags(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		queryTotal         metrics.Int64Measure
		queryAttemptsTotal metrics.Int64Measure
		queryLatency       metrics.Float64Measure
		keyspaceKey        metrics.Key
		clusterNameKey     metrics.Key
		dataCenterKey      metrics.Key
		hostIDKey          metrics.Key
		hostPortKey        metrics.Key
		rackKey            metrics.Key
		versionKey         metrics.Key
		mu                 sync.Mutex
		ms                 []metrics.MeasurementWithTags
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
		           queryAttemptsTotal: nil,
		           queryLatency: nil,
		           keyspaceKey: nil,
		           clusterNameKey: nil,
		           dataCenterKey: nil,
		           hostIDKey: nil,
		           hostPortKey: nil,
		           rackKey: nil,
		           versionKey: nil,
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
		           queryTotal: nil,
		           queryAttemptsTotal: nil,
		           queryLatency: nil,
		           keyspaceKey: nil,
		           clusterNameKey: nil,
		           dataCenterKey: nil,
		           hostIDKey: nil,
		           hostPortKey: nil,
		           rackKey: nil,
		           versionKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
			cm := &cassandraMetrics{
				queryTotal:         test.fields.queryTotal,
				queryAttemptsTotal: test.fields.queryAttemptsTotal,
				queryLatency:       test.fields.queryLatency,
				keyspaceKey:        test.fields.keyspaceKey,
				clusterNameKey:     test.fields.clusterNameKey,
				dataCenterKey:      test.fields.dataCenterKey,
				hostIDKey:          test.fields.hostIDKey,
				hostPortKey:        test.fields.hostPortKey,
				rackKey:            test.fields.rackKey,
				versionKey:         test.fields.versionKey,
				mu:                 test.fields.mu,
				ms:                 test.fields.ms,
			}

			got, err := cm.MeasurementWithTags(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_cassandraMetrics_View(t *testing.T) {
	t.Parallel()
	type fields struct {
		queryTotal         metrics.Int64Measure
		queryAttemptsTotal metrics.Int64Measure
		queryLatency       metrics.Float64Measure
		keyspaceKey        metrics.Key
		clusterNameKey     metrics.Key
		dataCenterKey      metrics.Key
		hostIDKey          metrics.Key
		hostPortKey        metrics.Key
		rackKey            metrics.Key
		versionKey         metrics.Key
		mu                 sync.Mutex
		ms                 []metrics.MeasurementWithTags
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
		           queryAttemptsTotal: nil,
		           queryLatency: nil,
		           keyspaceKey: nil,
		           clusterNameKey: nil,
		           dataCenterKey: nil,
		           hostIDKey: nil,
		           hostPortKey: nil,
		           rackKey: nil,
		           versionKey: nil,
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
		           queryTotal: nil,
		           queryAttemptsTotal: nil,
		           queryLatency: nil,
		           keyspaceKey: nil,
		           clusterNameKey: nil,
		           dataCenterKey: nil,
		           hostIDKey: nil,
		           hostPortKey: nil,
		           rackKey: nil,
		           versionKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
			cm := &cassandraMetrics{
				queryTotal:         test.fields.queryTotal,
				queryAttemptsTotal: test.fields.queryAttemptsTotal,
				queryLatency:       test.fields.queryLatency,
				keyspaceKey:        test.fields.keyspaceKey,
				clusterNameKey:     test.fields.clusterNameKey,
				dataCenterKey:      test.fields.dataCenterKey,
				hostIDKey:          test.fields.hostIDKey,
				hostPortKey:        test.fields.hostPortKey,
				rackKey:            test.fields.rackKey,
				versionKey:         test.fields.versionKey,
				mu:                 test.fields.mu,
				ms:                 test.fields.ms,
			}

			got := cm.View()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_cassandraMetrics_ObserveQuery(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		q   cassandra.ObservedQuery
	}
	type fields struct {
		queryTotal         metrics.Int64Measure
		queryAttemptsTotal metrics.Int64Measure
		queryLatency       metrics.Float64Measure
		keyspaceKey        metrics.Key
		clusterNameKey     metrics.Key
		dataCenterKey      metrics.Key
		hostIDKey          metrics.Key
		hostPortKey        metrics.Key
		rackKey            metrics.Key
		versionKey         metrics.Key
		mu                 sync.Mutex
		ms                 []metrics.MeasurementWithTags
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
		           q: nil,
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryAttemptsTotal: nil,
		           queryLatency: nil,
		           keyspaceKey: nil,
		           clusterNameKey: nil,
		           dataCenterKey: nil,
		           hostIDKey: nil,
		           hostPortKey: nil,
		           rackKey: nil,
		           versionKey: nil,
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
		           q: nil,
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryAttemptsTotal: nil,
		           queryLatency: nil,
		           keyspaceKey: nil,
		           clusterNameKey: nil,
		           dataCenterKey: nil,
		           hostIDKey: nil,
		           hostPortKey: nil,
		           rackKey: nil,
		           versionKey: nil,
		           mu: sync.Mutex{},
		           ms: nil,
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
			cm := &cassandraMetrics{
				queryTotal:         test.fields.queryTotal,
				queryAttemptsTotal: test.fields.queryAttemptsTotal,
				queryLatency:       test.fields.queryLatency,
				keyspaceKey:        test.fields.keyspaceKey,
				clusterNameKey:     test.fields.clusterNameKey,
				dataCenterKey:      test.fields.dataCenterKey,
				hostIDKey:          test.fields.hostIDKey,
				hostPortKey:        test.fields.hostPortKey,
				rackKey:            test.fields.rackKey,
				versionKey:         test.fields.versionKey,
				mu:                 test.fields.mu,
				ms:                 test.fields.ms,
			}

			cm.ObserveQuery(test.args.ctx, test.args.q)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
