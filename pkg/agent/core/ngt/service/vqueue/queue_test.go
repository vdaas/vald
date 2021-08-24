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

// Package vqueue manages the vector cache layer for reducing FFI overhead for fast Agent processing.
package vqueue

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Queue
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Queue, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Queue, err error) error {
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_PushInsert(t *testing.T) {
	type args struct {
		uuid   string
		vector []float32
		date   int64
	}
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
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
		           uuid: "",
		           vector: nil,
		           date: 0,
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           uuid: "",
		           vector: nil,
		           date: 0,
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			err := v.PushInsert(test.args.uuid, test.args.vector, test.args.date)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_PushDelete(t *testing.T) {
	type args struct {
		uuid string
		date int64
	}
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
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
		           uuid: "",
		           date: 0,
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           uuid: "",
		           date: 0,
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			err := v.PushDelete(test.args.uuid, test.args.date)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_RangePopInsert(t *testing.T) {
	type args struct {
		ctx context.Context
		now int64
		f   func(uuid string, vector []float32) bool
	}
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
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
		           f: nil,
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           f: nil,
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			v.RangePopInsert(test.args.ctx, test.args.now, test.args.f)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_RangePopDelete(t *testing.T) {
	type args struct {
		ctx context.Context
		now int64
		f   func(uuid string) bool
	}
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
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
		           f: nil,
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           f: nil,
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			v.RangePopDelete(test.args.ctx, test.args.now, test.args.f)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_GetVector(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
	}
	type want struct {
		want  []float32
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float32, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []float32, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           uuid: "",
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           uuid: "",
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			got, got1 := v.GetVector(test.args.uuid)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_IVExists(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
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
		           uuid: "",
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           uuid: "",
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			got := v.IVExists(test.args.uuid)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_DVExists(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
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
		           uuid: "",
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           uuid: "",
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			got := v.DVExists(test.args.uuid)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_addInsert(t *testing.T) {
	type args struct {
		i index
	}
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
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
		           i: index{},
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           i: index{},
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			v.addInsert(test.args.i)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_addDelete(t *testing.T) {
	type args struct {
		d key
	}
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
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
		           d: key{},
		       },
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
		           d: key{},
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           ichSize: 0,
		           dchSize: 0,
		           iBufSize: 0,
		           dBufSize: 0,
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			v.addDelete(test.args.d)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_IVQLen(t *testing.T) {
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
	}
	type want struct {
		wantL int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotL int) error {
		if !reflect.DeepEqual(gotL, w.wantL) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotL, w.wantL)
		}
		return nil
	}
	tests := []test{
		func() test {
			size := 0
			uii := make([]index, size)

			return test{
				name: "return 0 when the capacity and length is 0",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			size := 10
			uii := make([]index, size)

			return test{
				name: "return 10 when the capacity and length is 10",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			c := 10
			l := 5
			uii := make([]index, l, c)

			return test{
				name: "return 5 when the capacity is 10 and the length is 5",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: l,
				},
			}
		}(),
		func() test {
			iniLen := 5
			isrtSize := 2
			uii := make([]index, iniLen, 10)
			for i := 0; i < isrtSize; i++ {
				uii = append(uii, index{})
			}

			return test{
				name: "return 7 when the capacity is 10 and the initial length is 5 but the inserted size is 2",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: iniLen + isrtSize,
				},
			}
		}(),
		func() test {
			size := 10
			uii := make([]index, 0, size)
			for i := 0; i < size; i++ {
				uii = append(uii, index{})
			}

			return test{
				name: "return 10 when the inserted size is 10",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			insertSize := 5
			size := 10
			uii := make([]index, 0, size)
			for i := 0; i < insertSize; i++ {
				uii = append(uii, index{})
			}

			return test{
				name: "return 5 when the capacity is 10 and the inserted size is 5",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: insertSize,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			gotL := v.IVQLen()
			if err := test.checkFunc(test.want, gotL); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_DVQLen(t *testing.T) {
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
	}
	type want struct {
		wantL int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotL int) error {
		if !reflect.DeepEqual(gotL, w.wantL) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotL, w.wantL)
		}
		return nil
	}
	tests := []test{
		func() test {
			size := 0
			udk := make([]key, size)

			return test{
				name: "return 0 when the capacity and length is 0",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			size := 10
			udk := make([]key, size)

			return test{
				name: "return 10 when the capacity and length is 10",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			c := 10
			l := 5
			udk := make([]key, l, c)

			return test{
				name: "return 5 when the capacity is 10 and the length is 5",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: l,
				},
			}
		}(),
		func() test {
			iniLen := 5
			isrtSize := 2
			udk := make([]key, iniLen, 10)
			for i := 0; i < isrtSize; i++ {
				udk = append(udk, key{})
			}

			return test{
				name: "return 7 when the capacity is 10 and the initial length is 5 but the inserted size is 2",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: iniLen + isrtSize,
				},
			}
		}(),
		func() test {
			size := 10
			udk := make([]key, 0, size)
			for i := 0; i < size; i++ {
				udk = append(udk, key{})
			}

			return test{
				name: "return 10 when the inserted size is 10",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: size,
				},
			}
		}(),
		func() test {
			insertSize := 5
			size := 10
			udk := make([]key, 0, size)
			for i := 0; i < insertSize; i++ {
				udk = append(udk, key{})
			}

			return test{
				name: "return 5 when the capacity is 10 and the inserted size is 5",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: insertSize,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			gotL := v.DVQLen()
			if err := test.checkFunc(test.want, gotL); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_IVCLen(t *testing.T) {
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
	}
	type want struct {
		want int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got int) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			size := 0
			ich := make(chan index, size)

			return test{
				name: "return 0 when the buffer size is 0",
				fields: fields{
					ich: ich,
				},
				afterFunc: func() {
					close(ich)
				},
				want: want{
					want: size,
				},
			}
		}(),
		func() test {
			size := 10
			ich := make(chan index, size)

			return test{
				name: "return 0 when the buffer size is 10 but there is no inserted index",
				fields: fields{
					ich: ich,
				},
				afterFunc: func() {
					close(ich)
				},
				want: want{
					want: 0,
				},
			}
		}(),
		func() test {
			size := 10
			ich := make(chan index, size)
			for i := 0; i < size; i++ {
				ich <- index{}
			}

			return test{
				name: "return 10 when the buffer size is 10",
				fields: fields{
					ich: ich,
				},
				afterFunc: func() {
					close(ich)
				},
				want: want{
					want: size,
				},
			}
		}(),
		func() test {
			isrtSize := 5
			size := 10
			ich := make(chan index, size)
			for i := 0; i < isrtSize; i++ {
				ich <- index{}
			}

			return test{
				name: "return 5 when the buffer size is 10 but the inserted size is 5",
				fields: fields{
					ich: ich,
				},
				afterFunc: func() {
					close(ich)
				},
				want: want{
					want: isrtSize,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			got := v.IVCLen()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_DVCLen(t *testing.T) {
	type fields struct {
		ich      chan index
		uii      []index
		imu      sync.Mutex
		uiim     uiim
		dch      chan key
		udk      []key
		dmu      sync.Mutex
		udim     udim
		eg       errgroup.Group
		ichSize  int
		dchSize  int
		iBufSize int
		dBufSize int
	}
	type want struct {
		want int
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got int) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			size := 0
			dch := make(chan key, size)

			return test{
				name: "return 0 when the buffer size is 0",
				fields: fields{
					dch: dch,
				},
				afterFunc: func() {
					close(dch)
				},
				want: want{
					want: size,
				},
			}
		}(),
		func() test {
			size := 10
			dch := make(chan key, size)

			return test{
				name: "return 0 when the buffer size is 10 but there is no inserted key",
				fields: fields{
					dch: dch,
				},
				afterFunc: func() {
					close(dch)
				},
				want: want{
					want: 0,
				},
			}
		}(),
		func() test {
			size := 10
			dch := make(chan key, size)
			for i := 0; i < size; i++ {
				dch <- key{}
			}

			return test{
				name: "return 10 when the buffer size is 10",
				fields: fields{
					dch: dch,
				},
				afterFunc: func() {
					close(dch)
				},
				want: want{
					want: size,
				},
			}
		}(),
		func() test {
			insertSize := 5
			size := 10
			dch := make(chan key, size)
			for i := 0; i < insertSize; i++ {
				dch <- key{}
			}

			return test{
				name: "return 5 when the buffer size is 10 but the inserted size is 5",
				fields: fields{
					dch: dch,
				},
				afterFunc: func() {
					close(dch)
				},
				want: want{
					want: insertSize,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:      test.fields.ich,
				uii:      test.fields.uii,
				imu:      test.fields.imu,
				uiim:     test.fields.uiim,
				dch:      test.fields.dch,
				udk:      test.fields.udk,
				dmu:      test.fields.dmu,
				udim:     test.fields.udim,
				eg:       test.fields.eg,
				ichSize:  test.fields.ichSize,
				dchSize:  test.fields.dchSize,
				iBufSize: test.fields.iBufSize,
				dBufSize: test.fields.dBufSize,
			}

			got := v.DVCLen()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
