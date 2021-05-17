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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		func() test {
			var (
				uuid   = "535f54cc-b6d3-11eb-8529-0242ac130003"
				vector = []float32{1, 2, 3}
			)
			uiim := make(map[string]index)
			for uuid, index := range map[string]index{
				uuid: index{
					vector: vector,
				},
				"535f57ce-b6d3-11eb-8529-0242ac130003": index{
					vector: []float32{
						4, 5, 6,
					},
				},
			} {
				uiim[uuid] = index
			}
			return test{
				name: "return ([]float32{1, 2, 3}, true) when the uuid exists",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
				},
				want: want{
					want:  vector,
					want1: true,
				},
			}
		}(),
		func() test {
			var (
				uuid = "535f54cc-b6d3-11eb-8529-0242ac130003"
			)
			uiim := make(map[string]index)
			for uuid, index := range map[string]index{
				"535f58f0-b6d3-11eb-8529-0242ac130003": index{
					vector: []float32{1, 2, 3},
				},
				"535f57ce-b6d3-11eb-8529-0242ac130003": index{
					vector: []float32{
						4, 5, 6,
					},
				},
			} {
				uiim[uuid] = index
			}
			return test{
				name: "return ([]float32{1, 2, 3}, true) when the uuid exists",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: uiim,
				},
				want: want{
					want1: false,
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args, *vqueue)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "insert success",
				args: args{
					i: index{},
				},
				fields: fields{
					uii:  nil,
					uiim: nil,
				},
				want: want{},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

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
			if test.beforeFunc != nil {
				test.beforeFunc(test.args, v)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
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
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           d: key{},
		           },
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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

			v.addDelete(test.args.d)
			if err := test.checkFunc(test.want); err != nil {
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
		wantUii []index
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []index) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotUii []index) error {
		if !reflect.DeepEqual(gotUii, w.wantUii) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUii, w.wantUii)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
			if err := test.checkFunc(test.want, gotUii); err != nil {
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
		wantUdk []key
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []key) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotUdk []key) error {
		if !reflect.DeepEqual(gotUdk, w.wantUdk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUdk, w.wantUdk)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
		           fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: sync.Mutex{},
		           uiim: uiim{},
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
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
			if err := test.checkFunc(test.want, gotUdk); err != nil {
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
			size := 10
			uii := make([]index, size)
			return test{
				name: "return size of un inserted index",
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
			size := 10
			udk := make([]key, size)
			return test{
				name: "return size of undeleted key",
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
			size := 10
			ich := make(chan index, size)
			for i := 0; i < size; i++ {
				ich <- index{}
			}

			return test{
				name: "return size of insert queue",
				fields: fields{
					ich: ich,
				},
				want: want{
					want: size,
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
			size := 10
			dch := make(chan key, size)
			for i := 0; i < size; i++ {
				dch <- key{}
			}

			return test{
				name: "return size of delete queue",
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
