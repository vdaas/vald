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
package pool

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func Test_newMaxSlicePool(t *testing.T) {
	type args struct {
		sliceSize int64
	}
	type want struct {
		want *maxSlicePool
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *maxSlicePool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *maxSlicePool) error {
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
		           sliceSize: 0,
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
		           sliceSize: 0,
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

			got := newMaxSlicePool(test.args.sliceSize)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_maxSlicePool_Get(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		allocator      sliceAllocator
		slices         chan *[]byte
		allocations    chan struct{}
		capacityChange chan struct{}
		max            int
		sliceSize      int64
		mtx            sync.RWMutex
	}
	type want struct {
		want *[]byte
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *[]byte, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *[]byte, err error) error {
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
			p := &maxSlicePool{
				allocator:      test.fields.allocator,
				slices:         test.fields.slices,
				allocations:    test.fields.allocations,
				capacityChange: test.fields.capacityChange,
				max:            test.fields.max,
				sliceSize:      test.fields.sliceSize,
				mtx:            test.fields.mtx,
			}

			got, err := p.Get(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_maxSlicePool_Put(t *testing.T) {
	type args struct {
		bs *[]byte
	}
	type fields struct {
		allocator      sliceAllocator
		slices         chan *[]byte
		allocations    chan struct{}
		capacityChange chan struct{}
		max            int
		sliceSize      int64
		mtx            sync.RWMutex
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
		           bs: nil,
		       },
		       fields: fields {
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
		           bs: nil,
		           },
		           fields: fields {
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
			p := &maxSlicePool{
				allocator:      test.fields.allocator,
				slices:         test.fields.slices,
				allocations:    test.fields.allocations,
				capacityChange: test.fields.capacityChange,
				max:            test.fields.max,
				sliceSize:      test.fields.sliceSize,
				mtx:            test.fields.mtx,
			}

			p.Put(test.args.bs)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_maxSlicePool_ModifyCapacity(t *testing.T) {
	type args struct {
		delta int
	}
	type fields struct {
		allocator      sliceAllocator
		slices         chan *[]byte
		allocations    chan struct{}
		capacityChange chan struct{}
		max            int
		sliceSize      int64
		mtx            sync.RWMutex
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
		           delta: 0,
		       },
		       fields: fields {
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
		           delta: 0,
		           },
		           fields: fields {
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
			p := &maxSlicePool{
				allocator:      test.fields.allocator,
				slices:         test.fields.slices,
				allocations:    test.fields.allocations,
				capacityChange: test.fields.capacityChange,
				max:            test.fields.max,
				sliceSize:      test.fields.sliceSize,
				mtx:            test.fields.mtx,
			}

			p.ModifyCapacity(test.args.delta)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_maxSlicePool_notifyCapacity(t *testing.T) {
	type fields struct {
		allocator      sliceAllocator
		slices         chan *[]byte
		allocations    chan struct{}
		capacityChange chan struct{}
		max            int
		sliceSize      int64
		mtx            sync.RWMutex
	}
	type want struct{}
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
			p := &maxSlicePool{
				allocator:      test.fields.allocator,
				slices:         test.fields.slices,
				allocations:    test.fields.allocations,
				capacityChange: test.fields.capacityChange,
				max:            test.fields.max,
				sliceSize:      test.fields.sliceSize,
				mtx:            test.fields.mtx,
			}

			p.notifyCapacity()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_maxSlicePool_SliceSize(t *testing.T) {
	type fields struct {
		allocator      sliceAllocator
		slices         chan *[]byte
		allocations    chan struct{}
		capacityChange chan struct{}
		max            int
		sliceSize      int64
		mtx            sync.RWMutex
	}
	type want struct {
		want int64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, int64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got int64) error {
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
			p := &maxSlicePool{
				allocator:      test.fields.allocator,
				slices:         test.fields.slices,
				allocations:    test.fields.allocations,
				capacityChange: test.fields.capacityChange,
				max:            test.fields.max,
				sliceSize:      test.fields.sliceSize,
				mtx:            test.fields.mtx,
			}

			got := p.SliceSize()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_maxSlicePool_Close(t *testing.T) {
	type fields struct {
		allocator      sliceAllocator
		slices         chan *[]byte
		allocations    chan struct{}
		capacityChange chan struct{}
		max            int
		sliceSize      int64
		mtx            sync.RWMutex
	}
	type want struct{}
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
			p := &maxSlicePool{
				allocator:      test.fields.allocator,
				slices:         test.fields.slices,
				allocations:    test.fields.allocations,
				capacityChange: test.fields.capacityChange,
				max:            test.fields.max,
				sliceSize:      test.fields.sliceSize,
				mtx:            test.fields.mtx,
			}

			p.Close()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_maxSlicePool_empty(t *testing.T) {
	type fields struct {
		allocator      sliceAllocator
		slices         chan *[]byte
		allocations    chan struct{}
		capacityChange chan struct{}
		max            int
		sliceSize      int64
		mtx            sync.RWMutex
	}
	type want struct{}
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
			p := &maxSlicePool{
				allocator:      test.fields.allocator,
				slices:         test.fields.slices,
				allocations:    test.fields.allocations,
				capacityChange: test.fields.capacityChange,
				max:            test.fields.max,
				sliceSize:      test.fields.sliceSize,
				mtx:            test.fields.mtx,
			}

			p.empty()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_maxSlicePool_newSlice(t *testing.T) {
	type fields struct {
		allocator      sliceAllocator
		slices         chan *[]byte
		allocations    chan struct{}
		capacityChange chan struct{}
		max            int
		sliceSize      int64
		mtx            sync.RWMutex
	}
	type want struct {
		want *[]byte
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *[]byte) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *[]byte) error {
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
		           allocator: nil,
		           slices: nil,
		           allocations: nil,
		           capacityChange: nil,
		           max: 0,
		           sliceSize: 0,
		           mtx: sync.RWMutex{},
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
			p := &maxSlicePool{
				allocator:      test.fields.allocator,
				slices:         test.fields.slices,
				allocations:    test.fields.allocations,
				capacityChange: test.fields.capacityChange,
				max:            test.fields.max,
				sliceSize:      test.fields.sliceSize,
				mtx:            test.fields.mtx,
			}

			got := p.newSlice()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestReturnCapacityPoolCloser_ModifyCapacity(t *testing.T) {
	type args struct {
		delta int
	}
	type fields struct {
		ByteSlicePool  ByteSlicePool
		ReturnCapacity int
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
		           delta: 0,
		       },
		       fields: fields {
		           ByteSlicePool: nil,
		           ReturnCapacity: 0,
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
		           delta: 0,
		           },
		           fields: fields {
		           ByteSlicePool: nil,
		           ReturnCapacity: 0,
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
			n := &ReturnCapacityPoolCloser{
				ByteSlicePool:  test.fields.ByteSlicePool,
				ReturnCapacity: test.fields.ReturnCapacity,
			}

			n.ModifyCapacity(test.args.delta)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestReturnCapacityPoolCloser_Close(t *testing.T) {
	type fields struct {
		ByteSlicePool  ByteSlicePool
		ReturnCapacity int
	}
	type want struct{}
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
		           ByteSlicePool: nil,
		           ReturnCapacity: 0,
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
		           ByteSlicePool: nil,
		           ReturnCapacity: 0,
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
			n := &ReturnCapacityPoolCloser{
				ByteSlicePool:  test.fields.ByteSlicePool,
				ReturnCapacity: test.fields.ReturnCapacity,
			}

			n.Close()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
