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

// Package redis provides redis metrics functions
package redis

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/db/kvs/redis"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type want struct {
		wantO MetricsHook
		err   error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, MetricsHook, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotO MetricsHook, err error) error {
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

func Test_redisMetrics_Measurement(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		queryTotal      metrics.Int64Measure
		queryLatency    metrics.Float64Measure
		pipelineTotal   metrics.Int64Measure
		pipelineLatency metrics.Float64Measure
		cmdNameKey      metrics.Key
		numCmdKey       metrics.Key
		mu              sync.Mutex
		ms              []metrics.MeasurementWithTags
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
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
			rm := &redisMetrics{
				queryTotal:      test.fields.queryTotal,
				queryLatency:    test.fields.queryLatency,
				pipelineTotal:   test.fields.pipelineTotal,
				pipelineLatency: test.fields.pipelineLatency,
				cmdNameKey:      test.fields.cmdNameKey,
				numCmdKey:       test.fields.numCmdKey,
				mu:              test.fields.mu,
				ms:              test.fields.ms,
			}

			got, err := rm.Measurement(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisMetrics_MeasurementWithTags(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		queryTotal      metrics.Int64Measure
		queryLatency    metrics.Float64Measure
		pipelineTotal   metrics.Int64Measure
		pipelineLatency metrics.Float64Measure
		cmdNameKey      metrics.Key
		numCmdKey       metrics.Key
		mu              sync.Mutex
		ms              []metrics.MeasurementWithTags
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
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
			rm := &redisMetrics{
				queryTotal:      test.fields.queryTotal,
				queryLatency:    test.fields.queryLatency,
				pipelineTotal:   test.fields.pipelineTotal,
				pipelineLatency: test.fields.pipelineLatency,
				cmdNameKey:      test.fields.cmdNameKey,
				numCmdKey:       test.fields.numCmdKey,
				mu:              test.fields.mu,
				ms:              test.fields.ms,
			}

			got, err := rm.MeasurementWithTags(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisMetrics_View(t *testing.T) {
	t.Parallel()
	type fields struct {
		queryTotal      metrics.Int64Measure
		queryLatency    metrics.Float64Measure
		pipelineTotal   metrics.Int64Measure
		pipelineLatency metrics.Float64Measure
		cmdNameKey      metrics.Key
		numCmdKey       metrics.Key
		mu              sync.Mutex
		ms              []metrics.MeasurementWithTags
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
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
			rm := &redisMetrics{
				queryTotal:      test.fields.queryTotal,
				queryLatency:    test.fields.queryLatency,
				pipelineTotal:   test.fields.pipelineTotal,
				pipelineLatency: test.fields.pipelineLatency,
				cmdNameKey:      test.fields.cmdNameKey,
				numCmdKey:       test.fields.numCmdKey,
				mu:              test.fields.mu,
				ms:              test.fields.ms,
			}

			got := rm.View()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisMetrics_BeforeProcess(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		cmd redis.Cmder
	}
	type fields struct {
		queryTotal      metrics.Int64Measure
		queryLatency    metrics.Float64Measure
		pipelineTotal   metrics.Int64Measure
		pipelineLatency metrics.Float64Measure
		cmdNameKey      metrics.Key
		numCmdKey       metrics.Key
		mu              sync.Mutex
		ms              []metrics.MeasurementWithTags
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
		           cmd: nil,
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
		           cmd: nil,
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
			rm := &redisMetrics{
				queryTotal:      test.fields.queryTotal,
				queryLatency:    test.fields.queryLatency,
				pipelineTotal:   test.fields.pipelineTotal,
				pipelineLatency: test.fields.pipelineLatency,
				cmdNameKey:      test.fields.cmdNameKey,
				numCmdKey:       test.fields.numCmdKey,
				mu:              test.fields.mu,
				ms:              test.fields.ms,
			}

			got, err := rm.BeforeProcess(test.args.ctx, test.args.cmd)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisMetrics_AfterProcess(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		cmd redis.Cmder
	}
	type fields struct {
		queryTotal      metrics.Int64Measure
		queryLatency    metrics.Float64Measure
		pipelineTotal   metrics.Int64Measure
		pipelineLatency metrics.Float64Measure
		cmdNameKey      metrics.Key
		numCmdKey       metrics.Key
		mu              sync.Mutex
		ms              []metrics.MeasurementWithTags
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
		           cmd: nil,
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
		           cmd: nil,
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
			rm := &redisMetrics{
				queryTotal:      test.fields.queryTotal,
				queryLatency:    test.fields.queryLatency,
				pipelineTotal:   test.fields.pipelineTotal,
				pipelineLatency: test.fields.pipelineLatency,
				cmdNameKey:      test.fields.cmdNameKey,
				numCmdKey:       test.fields.numCmdKey,
				mu:              test.fields.mu,
				ms:              test.fields.ms,
			}

			err := rm.AfterProcess(test.args.ctx, test.args.cmd)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisMetrics_BeforeProcessPipeline(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		cmds []redis.Cmder
	}
	type fields struct {
		queryTotal      metrics.Int64Measure
		queryLatency    metrics.Float64Measure
		pipelineTotal   metrics.Int64Measure
		pipelineLatency metrics.Float64Measure
		cmdNameKey      metrics.Key
		numCmdKey       metrics.Key
		mu              sync.Mutex
		ms              []metrics.MeasurementWithTags
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
		           cmds: nil,
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
		           cmds: nil,
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
			rm := &redisMetrics{
				queryTotal:      test.fields.queryTotal,
				queryLatency:    test.fields.queryLatency,
				pipelineTotal:   test.fields.pipelineTotal,
				pipelineLatency: test.fields.pipelineLatency,
				cmdNameKey:      test.fields.cmdNameKey,
				numCmdKey:       test.fields.numCmdKey,
				mu:              test.fields.mu,
				ms:              test.fields.ms,
			}

			got, err := rm.BeforeProcessPipeline(test.args.ctx, test.args.cmds)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_redisMetrics_AfterProcessPipeline(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		cmds []redis.Cmder
	}
	type fields struct {
		queryTotal      metrics.Int64Measure
		queryLatency    metrics.Float64Measure
		pipelineTotal   metrics.Int64Measure
		pipelineLatency metrics.Float64Measure
		cmdNameKey      metrics.Key
		numCmdKey       metrics.Key
		mu              sync.Mutex
		ms              []metrics.MeasurementWithTags
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
		           cmds: nil,
		       },
		       fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
		           cmds: nil,
		           },
		           fields: fields {
		           queryTotal: nil,
		           queryLatency: nil,
		           pipelineTotal: nil,
		           pipelineLatency: nil,
		           cmdNameKey: nil,
		           numCmdKey: nil,
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
			rm := &redisMetrics{
				queryTotal:      test.fields.queryTotal,
				queryLatency:    test.fields.queryLatency,
				pipelineTotal:   test.fields.pipelineTotal,
				pipelineLatency: test.fields.pipelineLatency,
				cmdNameKey:      test.fields.cmdNameKey,
				numCmdKey:       test.fields.numCmdKey,
				mu:              test.fields.mu,
				ms:              test.fields.ms,
			}

			err := rm.AfterProcessPipeline(test.args.ctx, test.args.cmds)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
