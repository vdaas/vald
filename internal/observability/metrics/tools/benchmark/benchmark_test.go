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
package benchmark

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/operator/benchmark/service"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestNew(t *testing.T) {
	type args struct {
		om service.Operator
	}
	type want struct {
		want metrics.Metric
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, metrics.Metric) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
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
		           om:nil,
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
		           om:nil,
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

			got := New(test.args.om)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_operatorMetrics_View(t *testing.T) {
	type fields struct {
		op service.Operator
	}
	type want struct {
		err  error
		want []metrics.View
	}
	type test struct {
		fields     fields
		checkFunc  func(want, []metrics.View, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, got []metrics.View, err error) error {
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
		       fields: fields {
		           op:nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		       beforeFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		       afterFunc: func(t *testing.T,) {
		           t.Helper()
		       },
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           op:nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		           beforeFunc: func(t *testing.T,) {
		               t.Helper()
		           },
		           afterFunc: func(t *testing.T,) {
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
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			om := &operatorMetrics{
				op: test.fields.op,
			}

			got, err := om.View()
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_operatorMetrics_Register(t *testing.T) {
	type args struct {
		m metrics.Meter
	}
	type fields struct {
		op service.Operator
	}
	type want struct {
		err error
	}
	type test struct {
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
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
		           m:nil,
		       },
		       fields: fields {
		           op:nil,
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
		           m:nil,
		           },
		           fields: fields {
		           op:nil,
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
			om := &operatorMetrics{
				op: test.fields.op,
			}

			err := om.Register(test.args.m)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

type fakeOperator struct {
	service.Operator
	scenario map[v1.ValdBenchmarkScenarioStatus]int64
	job      map[v1.BenchmarkJobStatus]int64
}

func (f *fakeOperator) GetScenarioStatus() map[v1.ValdBenchmarkScenarioStatus]int64 {
	return f.scenario
}

func (f *fakeOperator) GetBenchmarkJobStatus() map[v1.BenchmarkJobStatus]int64 {
	return f.job
}

// Test_operatorMetrics_Register_instruments registers the metrics against a real SDK meter
// with the views returned by View() and asserts the name, description, unit and value of
// every instrument. A view only matches an instrument when name and description are equal,
// so a Gauge result also proves View() and Register() agree on the description.
func Test_operatorMetrics_Register_instruments(t *testing.T) {
	t.Parallel()

	type want struct {
		desc  string
		value int64
	}
	wants := map[string]want{
		"benchmark_operator_applied_scenario":       {desc: "Benchmark Operator applied scenario count", value: 5},
		"benchmark_operator_running_scenario":       {desc: "Benchmark Operator running scenario count", value: 2},
		"benchmark_operator_complete_scenario":      {desc: "Benchmark Operator complete scenario count", value: 3},
		"benchmark_operator_applied_benchmark_job":  {desc: "Benchmark Operator applied benchmark job count", value: 7},
		"benchmark_operator_running_benchmark_job":  {desc: "Benchmark Operator running benchmark job count", value: 4},
		"benchmark_operator_complete_benchmark_job": {desc: "Benchmark Operator complete benchmark job count", value: 3},
	}

	om := &operatorMetrics{
		op: &fakeOperator{
			scenario: map[v1.ValdBenchmarkScenarioStatus]int64{
				v1.BenchmarkScenarioHealthy:   2,
				v1.BenchmarkScenarioCompleted: 3,
			},
			job: map[v1.BenchmarkJobStatus]int64{
				v1.BenchmarkJobHealthy:   1,
				v1.BenchmarkJobNotReady:  3,
				v1.BenchmarkJobCompleted: 3,
			},
		},
	}
	views, err := om.View()
	if err != nil {
		t.Fatalf("View() error = %v", err)
	}
	opts := make([]sdkmetric.Option, 0, len(views)+1)
	reader := sdkmetric.NewManualReader()
	opts = append(opts, sdkmetric.WithReader(reader))
	for _, vw := range views {
		opts = append(opts, sdkmetric.WithView(vw))
	}
	mp := sdkmetric.NewMeterProvider(opts...)
	t.Cleanup(func() {
		if err := mp.Shutdown(context.WithoutCancel(t.Context())); err != nil {
			t.Errorf("Shutdown() error = %v", err)
		}
	})

	if err := om.Register(mp.Meter("test")); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(t.Context(), &rm); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	got := make(map[string]metricdata.Metrics)
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			got[m.Name] = m
		}
	}
	if len(got) != len(wants) {
		t.Errorf("registered %d instruments, want %d: %v", len(got), len(wants), got)
	}
	for name, w := range wants {
		m, ok := got[name]
		if !ok {
			t.Errorf("instrument %q is not registered", name)
			continue
		}
		if m.Description != w.desc {
			t.Errorf("%s: description = %q, want %q", name, m.Description, w.desc)
		}
		if m.Unit != metrics.Dimensionless {
			t.Errorf("%s: unit = %q, want %q", name, m.Unit, metrics.Dimensionless)
		}
		g, ok := m.Data.(metricdata.Gauge[int64])
		if !ok {
			t.Errorf("%s: data = %T, want metricdata.Gauge[int64] (view did not match)", name, m.Data)
			continue
		}
		if len(g.DataPoints) != 1 || g.DataPoints[0].Value != w.value {
			t.Errorf("%s: data points = %v, want one point with value %d", name, g.DataPoints, w.value)
		}
	}
}

// NOT IMPLEMENTED BELOW
