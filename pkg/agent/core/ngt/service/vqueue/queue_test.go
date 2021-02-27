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
	t.Parallel()
	type args struct {
		eg errgroup.Group
	}
	type want struct {
		want Queue
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Queue) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Queue) error {
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
		           eg: nil,
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
		           eg: nil,
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

			got := New(test.args.eg)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_Start(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			got, err := v.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_PushInsert(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid   string
		vector []float32
		date   int64
	}
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			err := v.PushInsert(test.args.uuid, test.args.vector, test.args.date)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_PushDelete(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
		date int64
	}
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			err := v.PushDelete(test.args.uuid, test.args.date)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_PopInsert(t *testing.T) {
	t.Parallel()
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
	}
	type want struct {
		wantUuid   string
		wantVector []float32
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string, []float32) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotUuid string, gotVector []float32) error {
		if !reflect.DeepEqual(gotUuid, w.wantUuid) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUuid, w.wantUuid)
		}
		if !reflect.DeepEqual(gotVector, w.wantVector) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVector, w.wantVector)
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			gotUuid, gotVector := v.PopInsert()
			if err := test.checkFunc(test.want, gotUuid, gotVector); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_PopDelete(t *testing.T) {
	t.Parallel()
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
	}
	type want struct {
		wantUuid string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotUuid string) error {
		if !reflect.DeepEqual(gotUuid, w.wantUuid) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUuid, w.wantUuid)
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			gotUuid := v.PopDelete()
			if err := test.checkFunc(test.want, gotUuid); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_RangePopInsert(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		f   func(uuid string, vector []float32) bool
	}
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			v.RangePopInsert(test.args.ctx, test.args.f)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_RangePopDelete(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		f   func(uuid string) bool
	}
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			v.RangePopDelete(test.args.ctx, test.args.f)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_GetVector(t *testing.T) {
	t.Parallel()
	type args struct {
		uuid string
	}
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           imu: sync.Mutex{},
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			got, got1 := v.GetVector(test.args.uuid)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_addInsert(t *testing.T) {
	t.Parallel()
	type args struct {
		i index
	}
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           imu: sync.Mutex{},
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           imu: sync.Mutex{},
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			v.addInsert(test.args.i)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_addDelete(t *testing.T) {
	t.Parallel()
	type args struct {
		d key
	}
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			v.addDelete(test.args.d)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_popInsert(t *testing.T) {
	t.Parallel()
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
	}
	type want struct {
		wantI index
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, index) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotI index) error {
		if !reflect.DeepEqual(gotI, w.wantI) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotI, w.wantI)
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			gotI := v.popInsert()
			if err := test.checkFunc(test.want, gotI); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_popDelete(t *testing.T) {
	t.Parallel()
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
	}
	type want struct {
		wantD key
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, key) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, gotD key) error {
		if !reflect.DeepEqual(gotD, w.wantD) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotD, w.wantD)
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			gotD := v.popDelete()
			if err := test.checkFunc(test.want, gotD); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_flushAndLoadInsert(t *testing.T) {
	t.Parallel()
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			gotUii := v.flushAndLoadInsert()
			if err := test.checkFunc(test.want, gotUii); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_flushAndLoadDelete(t *testing.T) {
	t.Parallel()
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			gotUdk := v.flushAndLoadDelete()
			if err := test.checkFunc(test.want, gotUdk); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_IVQLen(t *testing.T) {
	t.Parallel()
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: sync.Mutex{},
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			gotL := v.IVQLen()
			if err := test.checkFunc(test.want, gotL); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_vqueue_DVQLen(t *testing.T) {
	t.Parallel()
	type fields struct {
		ich        chan index
		uii        []index
		imu        sync.Mutex
		uiil       map[string][]float32
		dch        chan key
		udk        []key
		dmu        sync.Mutex
		eg         errgroup.Group
		finalizing atomic.Value
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           ich: nil,
		           uii: nil,
		           imu: sync.Mutex{},
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
		           uiil: nil,
		           dch: nil,
		           udk: nil,
		           dmu: sync.Mutex{},
		           eg: nil,
		           finalizing: nil,
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
				ich:        test.fields.ich,
				uii:        test.fields.uii,
				imu:        test.fields.imu,
				uiil:       test.fields.uiil,
				dch:        test.fields.dch,
				udk:        test.fields.udk,
				dmu:        test.fields.dmu,
				eg:         test.fields.eg,
				finalizing: test.fields.finalizing,
			}

			gotL := v.DVQLen()
			if err := test.checkFunc(test.want, gotL); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
