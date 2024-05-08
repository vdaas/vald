//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// NOT IMPLEMENTED BELOW
//
// func Test_newStrategy(t *testing.T) {
// 	type args struct {
// 		opts []StrategyOption
// 	}
// 	type want struct {
// 		want benchmark.Strategy
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, benchmark.Strategy) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got benchmark.Strategy) error {
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
// 			got := newStrategy(test.args.opts...)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_strategy_Init(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		b       *testing.B
// 		dataset assets.Dataset
// 	}
// 	type fields struct {
// 		core32    algorithm.Bit32
// 		core64    algorithm.Bit64
// 		initBit32 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit32, algorithm.Closer, error)
// 		initBit64 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit64, algorithm.Closer, error)
// 		closer    algorithm.Closer
// 		propName  string
// 		preProp32 func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset) ([]uint, error)
// 		preProp64 func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset) ([]uint, error)
// 		mode      algorithm.Mode
// 		prop32    func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		prop64    func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		parallel  bool
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
// 		           b:testing.B{},
// 		           dataset:nil,
// 		       },
// 		       fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 		           b:testing.B{},
// 		           dataset:nil,
// 		           },
// 		           fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 			s := &strategy{
// 				core32:    test.fields.core32,
// 				core64:    test.fields.core64,
// 				initBit32: test.fields.initBit32,
// 				initBit64: test.fields.initBit64,
// 				closer:    test.fields.closer,
// 				propName:  test.fields.propName,
// 				preProp32: test.fields.preProp32,
// 				preProp64: test.fields.preProp64,
// 				mode:      test.fields.mode,
// 				prop32:    test.fields.prop32,
// 				prop64:    test.fields.prop64,
// 				parallel:  test.fields.parallel,
// 			}
//
// 			err := s.Init(test.args.ctx, test.args.b, test.args.dataset)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_strategy_PreProp(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		b       *testing.B
// 		dataset assets.Dataset
// 	}
// 	type fields struct {
// 		core32    algorithm.Bit32
// 		core64    algorithm.Bit64
// 		initBit32 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit32, algorithm.Closer, error)
// 		initBit64 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit64, algorithm.Closer, error)
// 		closer    algorithm.Closer
// 		propName  string
// 		preProp32 func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset) ([]uint, error)
// 		preProp64 func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset) ([]uint, error)
// 		mode      algorithm.Mode
// 		prop32    func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		prop64    func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		parallel  bool
// 	}
// 	type want struct {
// 		want []uint
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []uint, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []uint, err error) error {
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
// 		           ctx:nil,
// 		           b:testing.B{},
// 		           dataset:nil,
// 		       },
// 		       fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 		           b:testing.B{},
// 		           dataset:nil,
// 		           },
// 		           fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 			s := &strategy{
// 				core32:    test.fields.core32,
// 				core64:    test.fields.core64,
// 				initBit32: test.fields.initBit32,
// 				initBit64: test.fields.initBit64,
// 				closer:    test.fields.closer,
// 				propName:  test.fields.propName,
// 				preProp32: test.fields.preProp32,
// 				preProp64: test.fields.preProp64,
// 				mode:      test.fields.mode,
// 				prop32:    test.fields.prop32,
// 				prop64:    test.fields.prop64,
// 				parallel:  test.fields.parallel,
// 			}
//
// 			got, err := s.PreProp(test.args.ctx, test.args.b, test.args.dataset)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_strategy_Run(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		b       *testing.B
// 		dataset assets.Dataset
// 		ids     []uint
// 	}
// 	type fields struct {
// 		core32    algorithm.Bit32
// 		core64    algorithm.Bit64
// 		initBit32 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit32, algorithm.Closer, error)
// 		initBit64 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit64, algorithm.Closer, error)
// 		closer    algorithm.Closer
// 		propName  string
// 		preProp32 func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset) ([]uint, error)
// 		preProp64 func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset) ([]uint, error)
// 		mode      algorithm.Mode
// 		prop32    func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		prop64    func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		parallel  bool
// 	}
// 	type want struct {
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           b:testing.B{},
// 		           dataset:nil,
// 		           ids:nil,
// 		       },
// 		       fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 		           b:testing.B{},
// 		           dataset:nil,
// 		           ids:nil,
// 		           },
// 		           fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 			s := &strategy{
// 				core32:    test.fields.core32,
// 				core64:    test.fields.core64,
// 				initBit32: test.fields.initBit32,
// 				initBit64: test.fields.initBit64,
// 				closer:    test.fields.closer,
// 				propName:  test.fields.propName,
// 				preProp32: test.fields.preProp32,
// 				preProp64: test.fields.preProp64,
// 				mode:      test.fields.mode,
// 				prop32:    test.fields.prop32,
// 				prop64:    test.fields.prop64,
// 				parallel:  test.fields.parallel,
// 			}
//
// 			s.Run(test.args.ctx, test.args.b, test.args.dataset, test.args.ids)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_strategy_Close(t *testing.T) {
// 	type fields struct {
// 		core32    algorithm.Bit32
// 		core64    algorithm.Bit64
// 		initBit32 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit32, algorithm.Closer, error)
// 		initBit64 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit64, algorithm.Closer, error)
// 		closer    algorithm.Closer
// 		propName  string
// 		preProp32 func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset) ([]uint, error)
// 		preProp64 func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset) ([]uint, error)
// 		mode      algorithm.Mode
// 		prop32    func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		prop64    func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		parallel  bool
// 	}
// 	type want struct {
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &strategy{
// 				core32:    test.fields.core32,
// 				core64:    test.fields.core64,
// 				initBit32: test.fields.initBit32,
// 				initBit64: test.fields.initBit64,
// 				closer:    test.fields.closer,
// 				propName:  test.fields.propName,
// 				preProp32: test.fields.preProp32,
// 				preProp64: test.fields.preProp64,
// 				mode:      test.fields.mode,
// 				prop32:    test.fields.prop32,
// 				prop64:    test.fields.prop64,
// 				parallel:  test.fields.parallel,
// 			}
//
// 			s.Close()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_strategy_float32(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		b       *testing.B
// 		dataset assets.Dataset
// 		ids     []uint
// 		cnt     *uint64
// 	}
// 	type fields struct {
// 		core32    algorithm.Bit32
// 		core64    algorithm.Bit64
// 		initBit32 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit32, algorithm.Closer, error)
// 		initBit64 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit64, algorithm.Closer, error)
// 		closer    algorithm.Closer
// 		propName  string
// 		preProp32 func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset) ([]uint, error)
// 		preProp64 func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset) ([]uint, error)
// 		mode      algorithm.Mode
// 		prop32    func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		prop64    func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		parallel  bool
// 	}
// 	type want struct {
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           b:testing.B{},
// 		           dataset:nil,
// 		           ids:nil,
// 		           cnt:nil,
// 		       },
// 		       fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 		           b:testing.B{},
// 		           dataset:nil,
// 		           ids:nil,
// 		           cnt:nil,
// 		           },
// 		           fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 			s := &strategy{
// 				core32:    test.fields.core32,
// 				core64:    test.fields.core64,
// 				initBit32: test.fields.initBit32,
// 				initBit64: test.fields.initBit64,
// 				closer:    test.fields.closer,
// 				propName:  test.fields.propName,
// 				preProp32: test.fields.preProp32,
// 				preProp64: test.fields.preProp64,
// 				mode:      test.fields.mode,
// 				prop32:    test.fields.prop32,
// 				prop64:    test.fields.prop64,
// 				parallel:  test.fields.parallel,
// 			}
//
// 			s.float32(test.args.ctx, test.args.b, test.args.dataset, test.args.ids, test.args.cnt)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_strategy_float64(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		b       *testing.B
// 		dataset assets.Dataset
// 		ids     []uint
// 		cnt     *uint64
// 	}
// 	type fields struct {
// 		core32    algorithm.Bit32
// 		core64    algorithm.Bit64
// 		initBit32 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit32, algorithm.Closer, error)
// 		initBit64 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit64, algorithm.Closer, error)
// 		closer    algorithm.Closer
// 		propName  string
// 		preProp32 func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset) ([]uint, error)
// 		preProp64 func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset) ([]uint, error)
// 		mode      algorithm.Mode
// 		prop32    func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		prop64    func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset, []uint, *uint64) (interface{}, error)
// 		parallel  bool
// 	}
// 	type want struct {
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           b:testing.B{},
// 		           dataset:nil,
// 		           ids:nil,
// 		           cnt:nil,
// 		       },
// 		       fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 		           b:testing.B{},
// 		           dataset:nil,
// 		           ids:nil,
// 		           cnt:nil,
// 		           },
// 		           fields: fields {
// 		           core32:nil,
// 		           core64:nil,
// 		           initBit32:nil,
// 		           initBit64:nil,
// 		           closer:nil,
// 		           propName:"",
// 		           preProp32:nil,
// 		           preProp64:nil,
// 		           mode:nil,
// 		           prop32:nil,
// 		           prop64:nil,
// 		           parallel:false,
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
// 			s := &strategy{
// 				core32:    test.fields.core32,
// 				core64:    test.fields.core64,
// 				initBit32: test.fields.initBit32,
// 				initBit64: test.fields.initBit64,
// 				closer:    test.fields.closer,
// 				propName:  test.fields.propName,
// 				preProp32: test.fields.preProp32,
// 				preProp64: test.fields.preProp64,
// 				mode:      test.fields.mode,
// 				prop32:    test.fields.prop32,
// 				prop64:    test.fields.prop64,
// 				parallel:  test.fields.parallel,
// 			}
//
// 			s.float64(test.args.ctx, test.args.b, test.args.dataset, test.args.ids, test.args.cnt)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
