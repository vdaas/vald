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
	"sync/atomic"
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

func Test_vqueue_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
	}
	type want struct {
		want <-chan error
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
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
		           ich: nil,
		           uii: nil,
		           imu: nil,
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: nil,
		           udim: udim{},
		           eg: nil,
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
			}

			got, err := v.Start(test.args.ctx)
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
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
		f   func(uuid string, vector []float32) bool
	}
	type fields struct {
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
			}

			v.RangePopInsert(test.args.ctx, test.args.f)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_RangePopDelete(t *testing.T) {
	type args struct {
		ctx context.Context
		f   func(uuid string) bool
	}
	type fields struct {
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
			}

			v.RangePopDelete(test.args.ctx, test.args.f)
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
		           finalizingInsert: nil,
		           finalizingDelete: nil,
		           closed: nil,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
	}
	type want struct {
		uii []index
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *vqueue) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, v *vqueue) error {
		if !reflect.DeepEqual(v.uii, w.uii) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", v.uii, w.uii)
		}
		return nil
	}
	tests := []test{
		func() test {
			idx := index{
				uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
				date: 1000000000,
			}
			wantUii := []index{
				idx,
			}
			return test{
				name: "add insert successes",
				args: args{
					i: idx,
				},
				fields: fields{
					uii: make([]index, 0),
				},
				want: want{
					uii: wantUii,
				},
			}
		}(),
		func() test {
			preIdx := index{
				uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
				date: 1000000000,
			}
			idx := index{
				uuid: "209c583a-bc35-11eb-8529-0242ac130003",
				date: 2000000000,
			}
			wantUii := []index{
				preIdx, idx,
			}
			return test{
				name: "add insert successes when there is already data",
				args: args{
					i: idx,
				},
				fields: fields{
					uii: []index{
						preIdx,
					},
				},
				want: want{
					uii: wantUii,
				},
			}
		}(),
		func() test {
			preIdx := index{
				uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
				date: 1000000000,
			}
			idx := index{}
			wantUii := []index{
				preIdx, idx,
			}
			return test{
				name: "add insert successes when i is empty",
				args: args{},
				fields: fields{
					uii: []index{
						preIdx,
					},
				},
				want: want{
					uii: wantUii,
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
			}

			v.addInsert(test.args.i)
			if err := test.checkFunc(test.want, v); err != nil {
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
	}
	type want struct {
		udk []key
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *vqueue) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, v *vqueue) error {
		if !reflect.DeepEqual(v.udk, w.udk) {
			return errors.Errorf("udk got: \"%#v\",\n\t\t\t\tudk want: \"%#v\"", v.udk, w.udk)
		}
		return nil
	}
	tests := []test{
		func() test {
			k := key{
				uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
				date: 1000000000,
			}
			wantUdk := []key{
				k,
			}
			return test{
				name: "add delete successes",
				args: args{
					d: k,
				},
				fields: fields{
					udk: make([]key, 0),
				},
				want: want{
					udk: wantUdk,
				},
			}
		}(),
		func() test {
			preK := key{
				uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
				date: 1000000000,
			}
			k := key{
				uuid: "209c583a-bc35-11eb-8529-0242ac130003",
				date: 2000000000,
			}
			wantUdk := []key{
				preK, k,
			}
			return test{
				name: "add delete successes when there is already data",
				args: args{
					d: k,
				},
				fields: fields{
					udk: []key{
						preK,
					},
				},
				want: want{
					udk: wantUdk,
				},
			}
		}(),
		func() test {
			preK := key{
				uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
				date: 1000000000,
			}
			k := key{}
			wantUdk := []key{
				preK, k,
			}
			return test{
				name: "add delete successes when d is empty",
				args: args{},
				fields: fields{
					udk: []key{
						preK,
					},
				},
				want: want{
					udk: wantUdk,
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
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			v := &vqueue{
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
			}

			v.addDelete(test.args.d)
			if err := test.checkFunc(test.want, v); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_flushAndLoadInsert(t *testing.T) {
	type fields struct {
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
	}
	type want struct {
		wantUii  []index
		wantUiim map[string]index
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []index, *vqueue) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotUii []index, v *vqueue) error {
		if !reflect.DeepEqual(gotUii, w.wantUii) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUii, w.wantUii)
		}

		gotUiim := make(map[string]index)
		v.uiim.Range(func(key string, value index) bool {
			gotUiim[key] = value
			return true
		})
		if !reflect.DeepEqual(gotUiim, w.wantUiim) {
			return errors.Errorf("uiim got: \"%#v\",\n\t\t\t\tuiim want: \"%#v\"", gotUiim, w.wantUiim)
		}
		return nil
	}
	tests := []test{
		func() test {
			uii := []index{
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				{
					uuid: "209c583a-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "309c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},
				{
					uuid: "409c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
				{
					uuid: "509c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
			}

			var m uiim
			for _, idx := range uii {
				m.Store(idx.uuid, idx)
			}

			var (
				wantUii  = uii
				wantUiim = make(map[string]index)
			)

			return test{
				name: "return keys when there is no duplicate data in uii",
				fields: fields{
					uii:  uii,
					uiim: m,
				},
				want: want{
					wantUii:  wantUii,
					wantUiim: wantUiim,
				},
			}
		}(),
		func() test {
			uii := []index{
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				{
					uuid: "209c583a-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "309c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},
				{
					uuid: "409c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
				{
					uuid: "509c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
				// The following data are duplicate data.
				{
					uuid: "209c583a-bc35-11eb-8529-0242ac130003",
					date: 2500000000,
				},
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 1500000000,
				},
			}

			var m uiim
			for _, idx := range uii {
				m.Store(idx.uuid, idx)
			}

			var (
				wantUii = []index{
					{
						uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
						date: 1500000000,
					},
					{
						uuid: "209c583a-bc35-11eb-8529-0242ac130003",
						date: 2500000000,
					},
					{
						uuid: "309c5d9e-bc35-11eb-8529-0242ac130003",
						date: 3000000000,
					},
					{
						uuid: "409c5e66-bc35-11eb-8529-0242ac130003",
						date: 4000000000,
					},
					{
						uuid: "509c5f24-bc35-11eb-8529-0242ac130003",
						date: 5000000000,
					},
				}
				wantUiim = make(map[string]index)
			)

			return test{
				name: "return keys when there is duplicate data in uii",
				fields: fields{
					uii:  uii,
					uiim: m,
				},
				want: want{
					wantUii:  wantUii,
					wantUiim: wantUiim,
				},
			}
		}(),
		func() test {
			var (
				uii = make([]index, 0)
				m   uiim
			)

			var (
				wantUii  = make([]index, len(uii))
				wantUiim = make(map[string]index)
			)

			return test{
				name: "return keys when uii is empty",
				fields: fields{
					uii:  uii,
					uiim: m,
				},
				want: want{
					wantUii:  wantUii,
					wantUiim: wantUiim,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
			}

			gotUii := v.flushAndLoadInsert()
			if err := test.checkFunc(test.want, gotUii, v); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_vqueue_flushAndLoadDelete(t *testing.T) {
	type fields struct {
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
	}
	type want struct {
		wantUdk  []key
		wantUdim map[string]int64
		wantUii  []index
		wantUiim map[string]index
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []key, *vqueue) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotUdk []key, v *vqueue) error {
		if !reflect.DeepEqual(gotUdk, w.wantUdk) {
			return errors.Errorf("udk got: \"%#v\",\n\t\t\t\tudk want: \"%#v\"", gotUdk, w.wantUdk)
		}
		if !reflect.DeepEqual(v.uii, w.wantUii) {
			return errors.Errorf("uii got: \"%#v\",\n\t\t\t\tuii want: \"%#v\"", v.uii, w.wantUii)
		}

		gotUdim := make(map[string]int64)
		v.udim.Range(func(key string, value int64) bool {
			gotUdim[key] = value
			return true
		})
		if !reflect.DeepEqual(gotUdim, w.wantUdim) {
			return errors.Errorf("udim got: \"%#v\",\n\t\t\t\tudim want: \"%#v\"", gotUdim, w.wantUdim)
		}

		gotUiim := make(map[string]index)
		v.uiim.Range(func(key string, value index) bool {
			gotUiim[key] = value
			return true
		})
		if !reflect.DeepEqual(gotUiim, w.wantUiim) {
			return errors.Errorf("uiim got: \"%#v\",\n\t\t\t\tuiim want: \"%#v\"", gotUiim, w.wantUiim)
		}
		return nil
	}
	tests := []test{
		func() test {
			udk := []key{
				{
					uuid: "509c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
				{
					uuid: "409c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
				{
					uuid: "309c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},
				{
					uuid: "209c583a-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			var (
				uii  = make([]index, 0)
				uiim uiim
			)

			return test{
				name: "return keys when there is no duplicate data in udk and uii is empty",
				fields: fields{
					udk:  udk,
					udim: udim,
					uii:  uii,
					uiim: uiim,
				},
				want: want{
					wantUdk: []key{
						{
							uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
							date: 1000000000,
						},
						{
							uuid: "209c583a-bc35-11eb-8529-0242ac130003",
							date: 2000000000,
						},
						{
							uuid: "309c5d9e-bc35-11eb-8529-0242ac130003",
							date: 3000000000,
						},
						{
							uuid: "409c5e66-bc35-11eb-8529-0242ac130003",
							date: 4000000000,
						},

						{
							uuid: "509c5f24-bc35-11eb-8529-0242ac130003",
							date: 5000000000,
						},
					},
					wantUdim: make(map[string]int64),
					wantUii:  make([]index, 0),
					wantUiim: make(map[string]index),
				},
			}
		}(),
		func() test {
			udk := []key{
				{
					uuid: "509c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
				{
					uuid: "409c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
				{
					uuid: "309c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},
				{
					uuid: "209c583a-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				// The following data are duplicate data.
				{
					uuid: "209c583a-bc35-11eb-8529-0242ac130003",
					date: 2500000000,
				},
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 1500000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			var (
				uii  = make([]index, 0)
				uiim uiim
			)

			return test{
				name: "return keys when there is duplicate data in udk and uii is empty",
				fields: fields{
					udk:  udk,
					udim: udim,
					uii:  uii,
					uiim: uiim,
				},
				want: want{
					wantUdk: []key{
						{
							uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
							date: 1500000000,
						},
						{
							uuid: "209c583a-bc35-11eb-8529-0242ac130003",
							date: 2500000000,
						},
						{
							uuid: "309c5d9e-bc35-11eb-8529-0242ac130003",
							date: 3000000000,
						},
						{
							uuid: "409c5e66-bc35-11eb-8529-0242ac130003",
							date: 4000000000,
						},

						{
							uuid: "509c5f24-bc35-11eb-8529-0242ac130003",
							date: 5000000000,
						},
					},
					wantUdim: make(map[string]int64),
					wantUii:  make([]index, 0),
					wantUiim: make(map[string]index),
				},
			}
		}(),
		func() test {
			udk := []key{
				{
					uuid: "509c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
				{
					uuid: "409c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
				{
					uuid: "309c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},
				{
					uuid: "209c583a-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				// The following data are duplicate data.
				{
					uuid: "209c583a-bc35-11eb-8529-0242ac130003",
					date: 2500000000,
				},
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 1500000000,
				},
			}
			var udim udim
			for _, key := range udk {
				udim.Store(key.uuid, key.date)
			}

			uii := []index{
				{
					uuid: udk[5].uuid,
					date: 2200000000,
				},
				{
					uuid: udk[6].uuid,
					date: 1600000000,
				},
				{
					uuid: "746bbe1a-bc48-11eb-8529-0242ac130003",
					date: 1500000000,
				},
			}
			var uiim uiim
			for _, idx := range uii {
				uiim.Store(idx.uuid, idx)
			}

			var (
				wantUdim = make(map[string]int64)
				wantUii  = []index{
					uii[1], uii[2],
				}
				wantUiim = make(map[string]index)
			)

			for _, idx := range wantUii {
				wantUiim[idx.uuid] = idx
			}

			return test{
				name: "return keys when there is duplicate data in udk and uii",
				fields: fields{
					udk:  udk,
					udim: udim,
					uii:  uii,
					uiim: uiim,
				},
				want: want{
					wantUdk: []key{
						{
							uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
							date: 1500000000,
						},
						{
							uuid: "209c583a-bc35-11eb-8529-0242ac130003",
							date: 2500000000,
						},
						{
							uuid: "309c5d9e-bc35-11eb-8529-0242ac130003",
							date: 3000000000,
						},
						{
							uuid: "409c5e66-bc35-11eb-8529-0242ac130003",
							date: 4000000000,
						},
						{
							uuid: "509c5f24-bc35-11eb-8529-0242ac130003",
							date: 5000000000,
						},
					},
					wantUdim: wantUdim,
					wantUii:  wantUii,
					wantUiim: wantUiim,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
			}

			gotUdk := v.flushAndLoadDelete()
			if err := test.checkFunc(test.want, gotUdk, v); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_vqueue_IVQLen(t *testing.T) {
	type fields struct {
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
				name: "return 0 when the size of uii is 0",
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
				name: "return 10 when the size of uii is 10",
				fields: fields{
					uii: uii,
				},
				want: want{
					wantL: size,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
				name: "return 0 when the size of udk is 0",
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
				name: "return 10 when the size of udk is 10",
				fields: fields{
					udk: udk,
				},
				want: want{
					wantL: size,
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
			size := 0
			ich := make(chan index, size)

			return test{
				name: "return 10 when the buffer size is 0",
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
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
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             uiim
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             udim
		eg               errgroup.Group
		finalizingInsert atomic.Value
		finalizingDelete atomic.Value
		closed           atomic.Value
		ichSize          int
		dchSize          int
		iBufSize         int
		dBufSize         int
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
			size := 0
			dch := make(chan key, size)

			return test{
				name: "return 10 when the buffer size is 0",
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
				ich:              test.fields.ich,
				uii:              test.fields.uii,
				imu:              test.fields.imu,
				uiim:             test.fields.uiim,
				dch:              test.fields.dch,
				udk:              test.fields.udk,
				dmu:              test.fields.dmu,
				udim:             test.fields.udim,
				eg:               test.fields.eg,
				finalizingInsert: test.fields.finalizingInsert,
				finalizingDelete: test.fields.finalizingDelete,
				closed:           test.fields.closed,
				ichSize:          test.fields.ichSize,
				dchSize:          test.fields.dchSize,
				iBufSize:         test.fields.iBufSize,
				dBufSize:         test.fields.dBufSize,
			}

			got := v.DVCLen()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
