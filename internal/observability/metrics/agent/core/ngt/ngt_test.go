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

// Package ngt provides functions for ngt stats
package ngt

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		n service.NGT
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
		           n: nil,
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
		           n: nil,
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

			got := New(test.args.n)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngtMetrics_Measurement(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		ngt                       service.NGT
		indexCount                metrics.Int64Measure
		uncommittedIndexCount     metrics.Int64Measure
		insertVCacheCount         metrics.Int64Measure
		deleteVCacheCount         metrics.Int64Measure
		completedCreateIndexTotal metrics.Int64Measure
		executedProactiveGCTotal  metrics.Int64Measure
		isIndexing                metrics.Int64Measure
		isSaving                  metrics.Int64Measure
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
		           ngt: nil,
		           indexCount: nil,
		           uncommittedIndexCount: nil,
		           insertVCacheCount: nil,
		           deleteVCacheCount: nil,
		           completedCreateIndexTotal: nil,
		           executedProactiveGCTotal: nil,
		           isIndexing: nil,
		           isSaving: nil,
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
		           ngt: nil,
		           indexCount: nil,
		           uncommittedIndexCount: nil,
		           insertVCacheCount: nil,
		           deleteVCacheCount: nil,
		           completedCreateIndexTotal: nil,
		           executedProactiveGCTotal: nil,
		           isIndexing: nil,
		           isSaving: nil,
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
			n := &ngtMetrics{
				ngt:                       test.fields.ngt,
				indexCount:                test.fields.indexCount,
				uncommittedIndexCount:     test.fields.uncommittedIndexCount,
				insertVCacheCount:         test.fields.insertVCacheCount,
				deleteVCacheCount:         test.fields.deleteVCacheCount,
				completedCreateIndexTotal: test.fields.completedCreateIndexTotal,
				executedProactiveGCTotal:  test.fields.executedProactiveGCTotal,
				isIndexing:                test.fields.isIndexing,
				isSaving:                  test.fields.isSaving,
			}

			got, err := n.Measurement(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngtMetrics_MeasurementWithTags(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		ngt                       service.NGT
		indexCount                metrics.Int64Measure
		uncommittedIndexCount     metrics.Int64Measure
		insertVCacheCount         metrics.Int64Measure
		deleteVCacheCount         metrics.Int64Measure
		completedCreateIndexTotal metrics.Int64Measure
		executedProactiveGCTotal  metrics.Int64Measure
		isIndexing                metrics.Int64Measure
		isSaving                  metrics.Int64Measure
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
		           ngt: nil,
		           indexCount: nil,
		           uncommittedIndexCount: nil,
		           insertVCacheCount: nil,
		           deleteVCacheCount: nil,
		           completedCreateIndexTotal: nil,
		           executedProactiveGCTotal: nil,
		           isIndexing: nil,
		           isSaving: nil,
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
		           ngt: nil,
		           indexCount: nil,
		           uncommittedIndexCount: nil,
		           insertVCacheCount: nil,
		           deleteVCacheCount: nil,
		           completedCreateIndexTotal: nil,
		           executedProactiveGCTotal: nil,
		           isIndexing: nil,
		           isSaving: nil,
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
			n := &ngtMetrics{
				ngt:                       test.fields.ngt,
				indexCount:                test.fields.indexCount,
				uncommittedIndexCount:     test.fields.uncommittedIndexCount,
				insertVCacheCount:         test.fields.insertVCacheCount,
				deleteVCacheCount:         test.fields.deleteVCacheCount,
				completedCreateIndexTotal: test.fields.completedCreateIndexTotal,
				executedProactiveGCTotal:  test.fields.executedProactiveGCTotal,
				isIndexing:                test.fields.isIndexing,
				isSaving:                  test.fields.isSaving,
			}

			got, err := n.MeasurementWithTags(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngtMetrics_View(t *testing.T) {
	t.Parallel()
	type fields struct {
		ngt                       service.NGT
		indexCount                metrics.Int64Measure
		uncommittedIndexCount     metrics.Int64Measure
		insertVCacheCount         metrics.Int64Measure
		deleteVCacheCount         metrics.Int64Measure
		completedCreateIndexTotal metrics.Int64Measure
		executedProactiveGCTotal  metrics.Int64Measure
		isIndexing                metrics.Int64Measure
		isSaving                  metrics.Int64Measure
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
		           ngt: nil,
		           indexCount: nil,
		           uncommittedIndexCount: nil,
		           insertVCacheCount: nil,
		           deleteVCacheCount: nil,
		           completedCreateIndexTotal: nil,
		           executedProactiveGCTotal: nil,
		           isIndexing: nil,
		           isSaving: nil,
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
		           ngt: nil,
		           indexCount: nil,
		           uncommittedIndexCount: nil,
		           insertVCacheCount: nil,
		           deleteVCacheCount: nil,
		           completedCreateIndexTotal: nil,
		           executedProactiveGCTotal: nil,
		           isIndexing: nil,
		           isSaving: nil,
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
			n := &ngtMetrics{
				ngt:                       test.fields.ngt,
				indexCount:                test.fields.indexCount,
				uncommittedIndexCount:     test.fields.uncommittedIndexCount,
				insertVCacheCount:         test.fields.insertVCacheCount,
				deleteVCacheCount:         test.fields.deleteVCacheCount,
				completedCreateIndexTotal: test.fields.completedCreateIndexTotal,
				executedProactiveGCTotal:  test.fields.executedProactiveGCTotal,
				isIndexing:                test.fields.isIndexing,
				isSaving:                  test.fields.isSaving,
			}

			got := n.View()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
