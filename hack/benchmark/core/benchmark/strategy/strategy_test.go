//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package strategy provides benchmark strategy
package strategy

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
	"github.com/vdaas/vald/hack/benchmark/internal/errors"

	"go.uber.org/goleak"
)

func Test_newStrategy(t *testing.T) {
	type args struct {
		opts []StrategyOption
	}
	type want struct {
		want benchmark.Strategy
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, benchmark.Strategy) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got benchmark.Strategy) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts: nil,
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
		           opts: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := newStrategy(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_strategy_Init(t *testing.T) {
	type args struct {
		ctx     context.Context
		b       *testing.B
		dataset assets.Dataset
	}
	type fields struct {
		core32     core.Core32
		core64     core.Core64
		initCore32 func(context.Context, *testing.B, assets.Dataset) (core.Core32, core.Closer, error)
		initCore64 func(context.Context, *testing.B, assets.Dataset) (core.Core64, core.Closer, error)
		closer     core.Closer
		propName   string
		preProp32  func(context.Context, *testing.B, core.Core32, assets.Dataset) ([]uint, error)
		preProp64  func(context.Context, *testing.B, core.Core64, assets.Dataset) ([]uint, error)
		mode       core.Mode
		prop32     func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error)
		prop64     func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error)
		parallel   bool
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		           b: testing.B{},
		           dataset: nil,
		       },
		       fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
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
		           b: testing.B{},
		           dataset: nil,
		           },
		           fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &strategy{
				core32:     test.fields.core32,
				core64:     test.fields.core64,
				initCore32: test.fields.initCore32,
				initCore64: test.fields.initCore64,
				closer:     test.fields.closer,
				propName:   test.fields.propName,
				preProp32:  test.fields.preProp32,
				preProp64:  test.fields.preProp64,
				mode:       test.fields.mode,
				prop32:     test.fields.prop32,
				prop64:     test.fields.prop64,
				parallel:   test.fields.parallel,
			}

			err := s.Init(test.args.ctx, test.args.b, test.args.dataset)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_strategy_PreProp(t *testing.T) {
	type args struct {
		ctx     context.Context
		b       *testing.B
		dataset assets.Dataset
	}
	type fields struct {
		core32     core.Core32
		core64     core.Core64
		initCore32 func(context.Context, *testing.B, assets.Dataset) (core.Core32, core.Closer, error)
		initCore64 func(context.Context, *testing.B, assets.Dataset) (core.Core64, core.Closer, error)
		closer     core.Closer
		propName   string
		preProp32  func(context.Context, *testing.B, core.Core32, assets.Dataset) ([]uint, error)
		preProp64  func(context.Context, *testing.B, core.Core64, assets.Dataset) ([]uint, error)
		mode       core.Mode
		prop32     func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error)
		prop64     func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error)
		parallel   bool
	}
	type want struct {
		want []uint
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []uint, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []uint, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
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
		           b: testing.B{},
		           dataset: nil,
		       },
		       fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
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
		           b: testing.B{},
		           dataset: nil,
		           },
		           fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &strategy{
				core32:     test.fields.core32,
				core64:     test.fields.core64,
				initCore32: test.fields.initCore32,
				initCore64: test.fields.initCore64,
				closer:     test.fields.closer,
				propName:   test.fields.propName,
				preProp32:  test.fields.preProp32,
				preProp64:  test.fields.preProp64,
				mode:       test.fields.mode,
				prop32:     test.fields.prop32,
				prop64:     test.fields.prop64,
				parallel:   test.fields.parallel,
			}

			got, err := s.PreProp(test.args.ctx, test.args.b, test.args.dataset)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_strategy_Run(t *testing.T) {
	type args struct {
		ctx     context.Context
		b       *testing.B
		dataset assets.Dataset
		ids     []uint
	}
	type fields struct {
		core32     core.Core32
		core64     core.Core64
		initCore32 func(context.Context, *testing.B, assets.Dataset) (core.Core32, core.Closer, error)
		initCore64 func(context.Context, *testing.B, assets.Dataset) (core.Core64, core.Closer, error)
		closer     core.Closer
		propName   string
		preProp32  func(context.Context, *testing.B, core.Core32, assets.Dataset) ([]uint, error)
		preProp64  func(context.Context, *testing.B, core.Core64, assets.Dataset) ([]uint, error)
		mode       core.Mode
		prop32     func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error)
		prop64     func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error)
		parallel   bool
	}
	type want struct {
	}
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
		           b: testing.B{},
		           dataset: nil,
		           ids: nil,
		       },
		       fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
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
		           b: testing.B{},
		           dataset: nil,
		           ids: nil,
		           },
		           fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &strategy{
				core32:     test.fields.core32,
				core64:     test.fields.core64,
				initCore32: test.fields.initCore32,
				initCore64: test.fields.initCore64,
				closer:     test.fields.closer,
				propName:   test.fields.propName,
				preProp32:  test.fields.preProp32,
				preProp64:  test.fields.preProp64,
				mode:       test.fields.mode,
				prop32:     test.fields.prop32,
				prop64:     test.fields.prop64,
				parallel:   test.fields.parallel,
			}

			s.Run(test.args.ctx, test.args.b, test.args.dataset, test.args.ids)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_strategy_Close(t *testing.T) {
	type fields struct {
		core32     core.Core32
		core64     core.Core64
		initCore32 func(context.Context, *testing.B, assets.Dataset) (core.Core32, core.Closer, error)
		initCore64 func(context.Context, *testing.B, assets.Dataset) (core.Core64, core.Closer, error)
		closer     core.Closer
		propName   string
		preProp32  func(context.Context, *testing.B, core.Core32, assets.Dataset) ([]uint, error)
		preProp64  func(context.Context, *testing.B, core.Core64, assets.Dataset) ([]uint, error)
		mode       core.Mode
		prop32     func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error)
		prop64     func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error)
		parallel   bool
	}
	type want struct {
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
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
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &strategy{
				core32:     test.fields.core32,
				core64:     test.fields.core64,
				initCore32: test.fields.initCore32,
				initCore64: test.fields.initCore64,
				closer:     test.fields.closer,
				propName:   test.fields.propName,
				preProp32:  test.fields.preProp32,
				preProp64:  test.fields.preProp64,
				mode:       test.fields.mode,
				prop32:     test.fields.prop32,
				prop64:     test.fields.prop64,
				parallel:   test.fields.parallel,
			}

			s.Close()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_strategy_float32(t *testing.T) {
	type args struct {
		ctx     context.Context
		b       *testing.B
		dataset assets.Dataset
		ids     []uint
		cnt     *uint64
	}
	type fields struct {
		core32     core.Core32
		core64     core.Core64
		initCore32 func(context.Context, *testing.B, assets.Dataset) (core.Core32, core.Closer, error)
		initCore64 func(context.Context, *testing.B, assets.Dataset) (core.Core64, core.Closer, error)
		closer     core.Closer
		propName   string
		preProp32  func(context.Context, *testing.B, core.Core32, assets.Dataset) ([]uint, error)
		preProp64  func(context.Context, *testing.B, core.Core64, assets.Dataset) ([]uint, error)
		mode       core.Mode
		prop32     func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error)
		prop64     func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error)
		parallel   bool
	}
	type want struct {
	}
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
		           b: testing.B{},
		           dataset: nil,
		           ids: nil,
		           cnt: nil,
		       },
		       fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
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
		           b: testing.B{},
		           dataset: nil,
		           ids: nil,
		           cnt: nil,
		           },
		           fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &strategy{
				core32:     test.fields.core32,
				core64:     test.fields.core64,
				initCore32: test.fields.initCore32,
				initCore64: test.fields.initCore64,
				closer:     test.fields.closer,
				propName:   test.fields.propName,
				preProp32:  test.fields.preProp32,
				preProp64:  test.fields.preProp64,
				mode:       test.fields.mode,
				prop32:     test.fields.prop32,
				prop64:     test.fields.prop64,
				parallel:   test.fields.parallel,
			}

			s.float32(test.args.ctx, test.args.b, test.args.dataset, test.args.ids, test.args.cnt)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_strategy_float64(t *testing.T) {
	type args struct {
		ctx     context.Context
		b       *testing.B
		dataset assets.Dataset
		ids     []uint
		cnt     *uint64
	}
	type fields struct {
		core32     core.Core32
		core64     core.Core64
		initCore32 func(context.Context, *testing.B, assets.Dataset) (core.Core32, core.Closer, error)
		initCore64 func(context.Context, *testing.B, assets.Dataset) (core.Core64, core.Closer, error)
		closer     core.Closer
		propName   string
		preProp32  func(context.Context, *testing.B, core.Core32, assets.Dataset) ([]uint, error)
		preProp64  func(context.Context, *testing.B, core.Core64, assets.Dataset) ([]uint, error)
		mode       core.Mode
		prop32     func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error)
		prop64     func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error)
		parallel   bool
	}
	type want struct {
	}
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
		           b: testing.B{},
		           dataset: nil,
		           ids: nil,
		           cnt: nil,
		       },
		       fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
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
		           b: testing.B{},
		           dataset: nil,
		           ids: nil,
		           cnt: nil,
		           },
		           fields: fields {
		           core32: nil,
		           core64: nil,
		           initCore32: nil,
		           initCore64: nil,
		           closer: nil,
		           propName: "",
		           preProp32: nil,
		           preProp64: nil,
		           mode: nil,
		           prop32: nil,
		           prop64: nil,
		           parallel: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &strategy{
				core32:     test.fields.core32,
				core64:     test.fields.core64,
				initCore32: test.fields.initCore32,
				initCore64: test.fields.initCore64,
				closer:     test.fields.closer,
				propName:   test.fields.propName,
				preProp32:  test.fields.preProp32,
				preProp64:  test.fields.preProp64,
				mode:       test.fields.mode,
				prop32:     test.fields.prop32,
				prop64:     test.fields.prop64,
				parallel:   test.fields.parallel,
			}

			s.float64(test.args.ctx, test.args.b, test.args.dataset, test.args.ids, test.args.cnt)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
