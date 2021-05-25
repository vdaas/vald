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
	"github.com/vdaas/vald/internal/test/comparator"
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
		checkFunc  func(want, error, *vqueue) error
		beforeFunc func(args, *vqueue)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error, _ *vqueue) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return error when the finalizingInsert is true",
				args: args{
					uuid:   "5047f42c-bcfc-11eb-8529-0242ac130003",
					vector: []float32{1, 2, 3},
					date:   1000000000,
				},
				fields: fields{
					finalizingInsert: func() (v atomic.Value) {
						v.Store(true)
						return
					}(),
					ich: make(chan index),
				},
				want: want{
					err: errors.ErrVQueueFinalizing,
				},
			}
		}(),
		func() test {
			return test{
				name: "return error when the closed is true",
				args: args{
					uuid:   "5047f42c-bcfc-11eb-8529-0242ac130003",
					vector: []float32{1, 2, 3},
					date:   1000000000,
				},
				fields: fields{
					finalizingInsert: func() (v atomic.Value) {
						v.Store(false)
						return
					}(),
					closed: func() (v atomic.Value) {
						v.Store(true)
						return
					}(),
					ich: make(chan index),
				},
				want: want{
					err: errors.ErrVQueueFinalizing,
				},
			}
		}(),
		func() test {
			idx := index{
				uuid:   "5047f42c-bcfc-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
				date:   1000000000,
			}
			return test{
				name: "return nil when the push insert successes",
				args: args{
					uuid:   idx.uuid,
					vector: idx.vector,
					date:   idx.date,
				},
				fields: fields{
					finalizingInsert: func() (v atomic.Value) {
						v.Store(false)
						return
					}(),
					closed: func() (v atomic.Value) {
						v.Store(false)
						return
					}(),
					ich: make(chan index, 1),
				},
				checkFunc: func(w want, e error, v *vqueue) error {
					if err := defaultCheckFunc(w, e, v); err != nil {
						return err
					}
					got := <-v.ich
					if !reflect.DeepEqual(got, idx) {
						return errors.Errorf("got_index: \"%#v\",\n\t\t\t\tgot_index: \"%#v\"", got, idx)
					}
					return nil
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			idx := index{}
			return test{
				name: "return nil when the push insert successes and the arguments are empty",
				args: args{},
				fields: fields{
					finalizingInsert: func() (v atomic.Value) {
						v.Store(false)
						return
					}(),
					closed: func() (v atomic.Value) {
						v.Store(false)
						return
					}(),
					ich: make(chan index, 1),
				},
				checkFunc: func(w want, e error, v *vqueue) error {
					if err := defaultCheckFunc(w, e, v); err != nil {
						return err
					}
					got := <-v.ich

					opts := []comparator.Option{
						comparator.AllowUnexported(idx),
						comparator.IgnoreFields(idx, "date"),
					}
					if diff := comparator.Diff(idx, got, opts...); diff != "" {
						return errors.Errorf("got_index diff: %s", diff)
					}
					return nil
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			preIdx := index{
				uuid:   "5047f738-bcfc-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
				date:   900000000,
			}
			idx := index{
				uuid:   "5047f42c-bcfc-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
				date:   1000000000,
			}
			return test{
				name: "return nil when the pre-insert successes and push insert successes",
				args: args{
					uuid:   idx.uuid,
					vector: idx.vector,
					date:   idx.date,
				},
				fields: fields{
					finalizingInsert: func() (v atomic.Value) {
						v.Store(false)
						return
					}(),
					closed: func() (v atomic.Value) {
						v.Store(false)
						return
					}(),
					ich: make(chan index, 2),
				},
				beforeFunc: func(a args, v *vqueue) {
					v.ich <- preIdx
				},
				checkFunc: func(w want, e error, v *vqueue) error {
					if err := defaultCheckFunc(w, e, v); err != nil {
						return err
					}
					got := <-v.ich
					if !reflect.DeepEqual(got, preIdx) {
						return errors.Errorf("got_pre-index: \"%#v\",\n\t\t\t\tgot_pre-index: \"%#v\"", got, idx)
					}

					got = <-v.ich
					if !reflect.DeepEqual(got, idx) {
						return errors.Errorf("got_index: \"%#v\",\n\t\t\t\tgot_index: \"%#v\"", got, idx)
					}
					return nil
				},
				want: want{
					err: nil,
				},
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

			err := v.PushInsert(test.args.uuid, test.args.vector, test.args.date)
			if err := test.checkFunc(test.want, err, v); err != nil {
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
	type want struct {
		wantUii  []index
		wantUiim map[string]index
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []index, map[string]index) error
		beforeFunc func(args, *vqueue)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotUii []index, gotUiim map[string]index) error {
		if !reflect.DeepEqual(gotUii, w.wantUii) {
			return errors.Errorf("uii got: \"%#v\",\n\t\t\t\tuii want: \"%#v\"", gotUii, w.wantUii)
		}
		if !reflect.DeepEqual(gotUiim, w.wantUiim) {
			return errors.Errorf("uiim got: \"%#v\",\n\t\t\t\tuiim want: \"%#v\"", gotUiim, w.wantUiim)
		}
		return nil
	}
	tests := []test{
		func() test {
			idx := index{
				uuid:   "9d2091e2-b864-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
			}
			return test{
				name: "insert success",
				args: args{
					i: idx,
				},
				fields: fields{
					uii:  make([]index, 0),
					uiim: make(map[string]index),
				},
				want: want{
					wantUii: []index{idx},
					wantUiim: map[string]index{
						idx.uuid: idx,
					},
				},
			}
		}(),
		func() test {
			idx := index{
				uuid:   "9d2091e2-b864-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
			}
			preInsertedIdx := index{
				uuid:   "1d2091e2-b864-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
			}
			uii := []index{
				preInsertedIdx,
			}
			uiim := map[string]index{
				preInsertedIdx.uuid: preInsertedIdx,
			}
			return test{
				name: "insert success when there is already data",
				args: args{
					i: idx,
				},
				fields: fields{
					uii:  uii,
					uiim: uiim,
				},
				want: want{
					wantUii: []index{
						preInsertedIdx,
						idx,
					},
					wantUiim: map[string]index{
						preInsertedIdx.uuid: preInsertedIdx,
						idx.uuid:            idx,
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "insert success when the idx is empty",
				args: args{},
				fields: fields{
					uii:  make([]index, 0),
					uiim: make(map[string]index),
				},
				want: want{
					wantUii: []index{
						index{},
					},
					wantUiim: map[string]index{
						"": index{},
					},
				},
			}
		}(),
		func() test {
			idx := index{
				uuid:   "9d2091e2-b864-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
			}
			preInsertedIdx := index{
				uuid:   "9d2091e2-b864-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
			}
			uii := []index{
				preInsertedIdx,
			}
			uiim := map[string]index{
				preInsertedIdx.uuid: preInsertedIdx,
			}
			return test{
				name: "insert success when the same uuid exits",
				args: args{
					i: idx,
				},
				fields: fields{
					uii:  uii,
					uiim: uiim,
				},
				want: want{
					wantUii: []index{
						preInsertedIdx,
						idx,
					},
					wantUiim: map[string]index{
						idx.uuid: idx,
					},
				},
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
			if err := test.checkFunc(test.want, v.uii, v.uiim); err != nil {
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
		udim             map[string]int64
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
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []key, map[string]int64) error
		beforeFunc func(args, *vqueue)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotUdk []key, gotUdim map[string]int64) error {
		if !reflect.DeepEqual(gotUdk, w.wantUdk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUdk, w.wantUdk)
		}
		if !reflect.DeepEqual(gotUdim, w.wantUdim) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUdim, w.wantUdim)
		}
		return nil
	}
	tests := []test{
		func() test {
			k := key{
				uuid:   "9d2091e2-b864-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
				date:   1419933529,
			}
			return test{
				name: "delete success",
				args: args{
					d: k,
				},
				fields: fields{
					udk:  make([]key, 0),
					udim: make(map[string]int64),
				},
				want: want{
					wantUdk: []key{k},
					wantUdim: map[string]int64{
						k.uuid: k.date,
					},
				},
			}
		}(),
		func() test {
			k := key{
				uuid:   "9d2091e2-b864-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
				date:   1419933529,
			}
			preK := key{
				uuid:   "fabc74ae-b915-11eb-8529-0242ac130003",
				vector: []float32{4, 5, 6},
				date:   1019933529,
			}
			udk := []key{
				preK,
			}
			udim := map[string]int64{
				preK.uuid: preK.date,
			}
			return test{
				name: "delete success when there is already data",
				args: args{
					d: k,
				},
				fields: fields{
					udk:  udk,
					udim: udim,
				},
				want: want{
					wantUdk: []key{
						preK,
						k,
					},
					wantUdim: map[string]int64{
						k.uuid:    k.date,
						preK.uuid: preK.date,
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "delete success when the key is empty",
				args: args{},
				fields: fields{
					udk:  make([]key, 0),
					udim: make(map[string]int64),
				},
				want: want{
					wantUdk: []key{
						{},
					},
					wantUdim: map[string]int64{
						"": 0,
					},
				},
			}
		}(),
		func() test {
			k := key{
				uuid:   "9d2091e2-b864-11eb-8529-0242ac130003",
				vector: []float32{1, 2, 3},
				date:   1419933529,
			}
			preK := key{
				uuid:   "9d2091e2-b864-11eb-8529-0242ac130003",
				vector: []float32{4, 5, 6},
				date:   1019933529,
			}
			udk := []key{
				preK,
			}
			udim := map[string]int64{
				preK.uuid: preK.date,
			}
			return test{
				name: "delete success when the same uuid exits",
				args: args{
					d: k,
				},
				fields: fields{
					udk:  udk,
					udim: udim,
				},
				want: want{
					wantUdk: []key{
						preK,
						k,
					},
					wantUdim: map[string]int64{
						k.uuid: k.date,
					},
				},
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
				udim:             test.fields.udim,
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

			v.addDelete(test.args.d)
			if err := test.checkFunc(test.want, v.udk, v.udim); err != nil {
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
			return errors.Errorf("uii got: \"%#v\",\n\t\t\t\tuii want: \"%#v\"", gotUii, w.wantUii)
		}
		if !reflect.DeepEqual(v.uiim, w.wantUiim) {
			return errors.Errorf("uiim got: \"%#v\",\n\t\t\t\tuiim want: \"%#v\"", v.uiim, w.wantUiim)
		}
		return nil
	}
	tests := []test{
		func() test {
			uii := []index{
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "109c583a-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				{
					uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},
				{
					uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
				{
					uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			uiim := make(map[string]index)
			for _, idx := range uii {
				uiim[idx.uuid] = idx
			}
			return test{
				name: "return indexes when there is no duplicate data in uii",
				fields: fields{
					uii:  uii,
					uiim: uiim,
				},
				want: want{
					wantUii: []index{
						{
							uuid: "109c583a-bc35-11eb-8529-0242ac130003",
							date: 1000000000,
						},
						{
							uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
							date: 2000000000,
						},
						{
							uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
							date: 3000000000,
						},
						{
							uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
							date: 4000000000,
						},
						{
							uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
							date: 5000000000,
						},
					},
					wantUiim: make(map[string]index),
				},
			}
		}(),
		func() test {
			uii := []index{
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "109c583a-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				{
					uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},
				{
					uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
				// The following data are duplicate data with uii.
				{
					uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
				{
					uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
					date: 6000000000,
				},
			}
			uiim := make(map[string]index)
			for _, idx := range uii {
				uiim[idx.uuid] = idx
			}
			return test{
				name: "return indexes when there is duplicate data in uii",
				fields: fields{
					uii:  uii,
					uiim: uiim,
				},
				want: want{
					wantUii: []index{
						{
							uuid: "109c583a-bc35-11eb-8529-0242ac130003",
							date: 1000000000,
						},
						{
							uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
							date: 2000000000,
						},
						{
							uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
							date: 3000000000,
						},
						{
							uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
							date: 5000000000,
						},
						{
							uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
							date: 6000000000,
						},
					},
					wantUiim: make(map[string]index),
				},
			}
		}(),
		func() test {
			return test{
				name: "return indexes when the all of data are empty",
				fields: fields{
					uii:  make([]index, 0),
					uiim: make(map[string]index),
				},
				want: want{
					wantUii:  make([]index, 0),
					wantUiim: make(map[string]index),
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
		udim             map[string]int64
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
		if !reflect.DeepEqual(v.udim, w.wantUdim) {
			return errors.Errorf("uiim got: \"%#v\",\n\t\t\t\tudim want: \"%#v\"", v.udim, w.wantUdim)
		}
		if !reflect.DeepEqual(v.uii, w.wantUii) {
			return errors.Errorf("uii got: \"%#v\",\n\t\t\t\tuii want: \"%#v\"", v.uii, w.wantUii)
		}
		if !reflect.DeepEqual(v.uiim, w.wantUiim) {
			return errors.Errorf("uiim got: \"%#v\",\n\t\t\t\tuiim want: \"%#v\"", v.uiim, w.wantUiim)
		}
		return nil
	}
	tests := []test{
		func() test {
			udk := []key{
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "109c583a-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				{
					uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},
				{
					uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
				{
					uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			udim := make(map[string]int64)
			for _, key := range udk {
				udim[key.uuid] = key.date
			}
			return test{
				name: "return keys when there is no duplicate data in udk",
				fields: fields{
					udk:  udk,
					udim: udim,
				},
				want: want{
					wantUdk: []key{
						{
							uuid: "109c583a-bc35-11eb-8529-0242ac130003",
							date: 1000000000,
						},
						{
							uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
							date: 2000000000,
						},
						{
							uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
							date: 3000000000,
						},
						{
							uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
							date: 4000000000,
						},
						{
							uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
							date: 5000000000,
						},
					},
					wantUdim: make(map[string]int64),
				},
			}
		}(),
		func() test {
			udk := []key{
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "109c583a-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				{
					uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},
				{
					uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
				// The following data are duplicate data with udk.
				{
					uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
				{
					uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
					date: 6000000000,
				},
			}
			udim := make(map[string]int64)
			for _, key := range udk {
				udim[key.uuid] = key.date
			}
			return test{
				name: "return keys when there is duplicate data in udk",
				fields: fields{
					udk:  udk,
					udim: udim,
				},
				want: want{
					wantUdk: []key{
						{
							uuid: "109c583a-bc35-11eb-8529-0242ac130003",
							date: 1000000000,
						},
						{
							uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
							date: 2000000000,
						},
						{
							uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
							date: 3000000000,
						},
						{
							uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
							date: 5000000000,
						},
						{
							uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
							date: 6000000000,
						},
					},
					wantUdim: make(map[string]int64),
				},
			}
		}(),
		func() test {
			udk := []key{
				{
					uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
					date: 2000000000,
				},
				{
					uuid: "109c583a-bc35-11eb-8529-0242ac130003",
					date: 1000000000,
				},
				{
					uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3000000000,
				},

				// The following data is duplicate data with udk.
				{
					uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
					date: 3500000000,
				},
				// The following data are duplicate data with uii.
				{
					uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
					date: 5000000000,
				},
				{
					uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
					date: 4000000000,
				},
			}
			udim := make(map[string]int64)
			for _, key := range udk {
				udim[key.uuid] = key.date
			}

			uii := []index{
				{
					uuid: udk[4].uuid,
					date: 4500000000,
				},
				{
					uuid: udk[5].uuid,
					date: 3500000000,
				},
				{
					uuid: "746bbe1a-bc48-11eb-8529-0242ac130003",
					date: 1500000000,
				},
			}
			uiim := make(map[string]index)
			for _, idx := range uii {
				uiim[idx.uuid] = idx
			}
			return test{
				name: "return keys when there is duplicate data in uii and delete it",
				fields: fields{
					udk:  udk,
					udim: udim,
					uii:  uii,
					uiim: uiim,
				},
				want: want{
					wantUdk: []key{
						{
							uuid: "109c583a-bc35-11eb-8529-0242ac130003",
							date: 1000000000,
						},
						{
							uuid: "109c5c86-bc35-11eb-8529-0242ac130003",
							date: 2000000000,
						},
						{
							uuid: "109c5d9e-bc35-11eb-8529-0242ac130003",
							date: 3500000000,
						},
						{
							uuid: "109c5e66-bc35-11eb-8529-0242ac130003",
							date: 4000000000,
						},
						{
							uuid: "109c5f24-bc35-11eb-8529-0242ac130003",
							date: 5000000000,
						},
					},
					wantUdim: make(map[string]int64),
					wantUii: []index{
						{
							uuid: "746bbe1a-bc48-11eb-8529-0242ac130003",
							date: 1500000000,
						},
					},
					wantUiim: map[string]index{
						"746bbe1a-bc48-11eb-8529-0242ac130003": index{
							uuid: "746bbe1a-bc48-11eb-8529-0242ac130003",
							date: 1500000000,
						},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return keys when the all of data are empty",
				fields: fields{
					udk:  make([]key, 0),
					udim: make(map[string]int64),
				},
				want: want{
					wantUdk:  make([]key, 0),
					wantUdim: make(map[string]int64),
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
				udim:             test.fields.udim,
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

func Test_vqueue_IVExists(t *testing.T) {
	type args struct {
		uuid string
	}
	type fields struct {
		ich              chan index
		uii              []index
		imu              sync.Mutex
		uiim             map[string]index
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             map[string]int64
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
		beforeFunc func(args, *vqueue)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			m := []struct {
				i index
				k key
			}{
				{
					i: index{
						uuid: "5047ecac-bcfc-11eb-8529-0242ac130003",
					},
					k: key{
						uuid: "5047ecac-bcfc-11eb-8529-0242ac130003",
					},
				},
			}
			return test{
				name: "return false when the uuid not exits in uiim",
				args: args{
					uuid: "5047f026-bcfc-11eb-8529-0242ac130003",
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				beforeFunc: func(a args, v *vqueue) {
					for _, val := range m {
						v.addInsert(val.i)
						v.addDelete(val.k)
					}
				},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			return test{
				name: "return false when the uuid not exits in uiim and the uiim and udim are empty",
				args: args{
					uuid: "5047f026-bcfc-11eb-8529-0242ac130003",
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			uuid := "5047ecac-bcfc-11eb-8529-0242ac130003"
			idxs := []index{
				{
					uuid: uuid,
				},
			}
			return test{
				name: "return true when the uuid exits in uiim but not in udim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				beforeFunc: func(a args, v *vqueue) {
					for _, idx := range idxs {
						v.addInsert(idx)
					}
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "5047ecac-bcfc-11eb-8529-0242ac130003"
			vals := []struct {
				idx index
				k   key
			}{
				{
					idx: index{
						uuid: uuid,
						date: 1500000000,
					},
					k: key{
						uuid: uuid,
						date: 1000000000,
					},
				},
				{
					idx: index{},
					k:   key{},
				},
			}
			return test{
				name: "return true when the inserted index is the latest than the deleted index",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				beforeFunc: func(a args, v *vqueue) {
					for _, val := range vals {
						v.addInsert(val.idx)
						v.addDelete(val.k)
					}
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "5047ecac-bcfc-11eb-8529-0242ac130003"
			vals := []struct {
				idx index
				k   key
			}{
				{
					idx: index{
						uuid: uuid,
						date: 1000000000,
					},
					k: key{
						uuid: uuid,
						date: 1500000000,
					},
				},
				{
					idx: index{},
					k:   key{},
				},
			}
			return test{
				name: "return true when the inserted index is not the latest than the deleted index",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				beforeFunc: func(a args, v *vqueue) {
					for _, val := range vals {
						v.addInsert(val.idx)
						v.addDelete(val.k)
					}
				},
				want: want{
					want: false,
				},
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
		uiim             map[string]index
		dch              chan key
		udk              []key
		dmu              sync.Mutex
		udim             map[string]int64
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
		beforeFunc func(args, *vqueue)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			m := []struct {
				i index
				k key
			}{
				{
					i: index{
						uuid: "5047ecac-bcfc-11eb-8529-0242ac130003",
					},
					k: key{
						uuid: "5047ecac-bcfc-11eb-8529-0242ac130003",
					},
				},
			}
			return test{
				name: "return false when the uuid not exits in udim",
				args: args{
					uuid: "5047f026-bcfc-11eb-8529-0242ac130003",
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				beforeFunc: func(a args, v *vqueue) {
					for _, val := range m {
						v.addInsert(val.i)
						v.addDelete(val.k)
					}
				},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			return test{
				name: "return false when the uuid not exits in udim and the uiim and udim are empty",
				args: args{
					uuid: "5047f026-bcfc-11eb-8529-0242ac130003",
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				want: want{
					want: false,
				},
			}
		}(),
		func() test {
			uuid := "5047ecac-bcfc-11eb-8529-0242ac130003"
			keys := []key{
				{
					uuid: uuid,
				},
			}
			return test{
				name: "return true when the uuid exits in udim but not in uiim",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				beforeFunc: func(a args, v *vqueue) {
					for _, k := range keys {
						v.addDelete(k)
					}
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "5047ecac-bcfc-11eb-8529-0242ac130003"
			vals := []struct {
				idx index
				k   key
			}{
				{
					idx: index{
						uuid: uuid,
						date: 1000000000,
					},
					k: key{
						uuid: uuid,
						date: 1500000000,
					},
				},
				{
					idx: index{},
					k:   key{},
				},
			}
			return test{
				name: "return true when the deleted index is the latest than the inserted index",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				beforeFunc: func(a args, v *vqueue) {
					for _, val := range vals {
						v.addInsert(val.idx)
						v.addDelete(val.k)
					}
				},
				want: want{
					want: true,
				},
			}
		}(),
		func() test {
			uuid := "5047ecac-bcfc-11eb-8529-0242ac130003"
			vals := []struct {
				idx index
				k   key
			}{
				{
					idx: index{
						uuid: uuid,
						date: 1500000000,
					},
					k: key{
						uuid: uuid,
						date: 1000000000,
					},
				},
				{
					idx: index{},
					k:   key{},
				},
			}
			return test{
				name: "return true when the deleted index is not the latest than the inserted index",
				args: args{
					uuid: uuid,
				},
				fields: fields{
					uiim: make(map[string]index),
					udim: make(map[string]int64),
				},
				beforeFunc: func(a args, v *vqueue) {
					for _, val := range vals {
						v.addInsert(val.idx)
						v.addDelete(val.k)
					}
				},
				want: want{
					want: false,
				},
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

			got := v.DVExists(test.args.uuid)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
