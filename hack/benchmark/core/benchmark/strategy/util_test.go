//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func Test_wrapErrors(t *testing.T) {
	type args struct {
		errs []error
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
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
		           errs: nil,
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
		           errs: nil,
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			err := wrapErrors(test.args.errs)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_insertAndCreateIndex32(t *testing.T) {
	type args struct {
		ctx     context.Context
		c       algorithm.Bit32
		dataset assets.Dataset
	}
	type want struct {
		wantIds []uint
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []uint, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotIds []uint, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotIds, w.wantIds) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIds, w.wantIds)
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
		           c: nil,
		           dataset: nil,
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
		           c: nil,
		           dataset: nil,
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotIds, err := insertAndCreateIndex32(test.args.ctx, test.args.c, test.args.dataset)
			if err := checkFunc(test.want, gotIds, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_insertAndCreateIndex64(t *testing.T) {
	type args struct {
		ctx     context.Context
		c       algorithm.Bit64
		dataset assets.Dataset
	}
	type want struct {
		wantIds []uint
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []uint, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotIds []uint, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotIds, w.wantIds) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotIds, w.wantIds)
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
		           c: nil,
		           dataset: nil,
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
		           c: nil,
		           dataset: nil,
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotIds, err := insertAndCreateIndex64(test.args.ctx, test.args.c, test.args.dataset)
			if err := checkFunc(test.want, gotIds, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_float32To64(t *testing.T) {
	t.Parallel()
	type args struct {
		x []float32
	}
	type want struct {
		wantY []float64
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []float64) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotY []float64) error {
		if !reflect.DeepEqual(gotY, w.wantY) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotY, w.wantY)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           x: nil,
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
		           x: nil,
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotY := float32To64(test.args.x)
			if err := checkFunc(test.want, gotY); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
