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

// Package pool provides pool functionality for pooling buffer or others
package pool

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts []Option
	}
	type want struct {
		want Buffer
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Buffer) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Buffer) error {
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
		           ctx: nil,
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

			got := New(test.args.ctx, test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Get(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		size      uint64
		length    uint64
		limit     uint64
		allocator Allocator
		extender  Extender
		flusher   Flusher
		pool      sync.Pool
	}
	type want struct {
		want interface{}
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got interface{}) error {
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
			p := &pool{
				size:      test.fields.size,
				length:    test.fields.length,
				limit:     test.fields.limit,
				allocator: test.fields.allocator,
				extender:  test.fields.extender,
				flusher:   test.fields.flusher,
				pool:      test.fields.pool,
			}

			got := p.Get(test.args.ctx)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Put(t *testing.T) {
	type args struct {
		ctx  context.Context
		data interface{}
	}
	type fields struct {
		size      uint64
		length    uint64
		limit     uint64
		allocator Allocator
		extender  Extender
		flusher   Flusher
		pool      sync.Pool
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
		           data: nil,
		       },
		       fields: fields {
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
		           data: nil,
		           },
		           fields: fields {
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
			p := &pool{
				size:      test.fields.size,
				length:    test.fields.length,
				limit:     test.fields.limit,
				allocator: test.fields.allocator,
				extender:  test.fields.extender,
				flusher:   test.fields.flusher,
				pool:      test.fields.pool,
			}

			p.Put(test.args.ctx, test.args.data)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_PutWithResize(t *testing.T) {
	type args struct {
		ctx  context.Context
		data interface{}
		size uint64
	}
	type fields struct {
		size      uint64
		length    uint64
		limit     uint64
		allocator Allocator
		extender  Extender
		flusher   Flusher
		pool      sync.Pool
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
		           data: nil,
		           size: 0,
		       },
		       fields: fields {
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
		           data: nil,
		           size: 0,
		           },
		           fields: fields {
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
			p := &pool{
				size:      test.fields.size,
				length:    test.fields.length,
				limit:     test.fields.limit,
				allocator: test.fields.allocator,
				extender:  test.fields.extender,
				flusher:   test.fields.flusher,
				pool:      test.fields.pool,
			}

			p.PutWithResize(test.args.ctx, test.args.data, test.args.size)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Size(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		size      uint64
		length    uint64
		limit     uint64
		allocator Allocator
		extender  Extender
		flusher   Flusher
		pool      sync.Pool
	}
	type want struct {
		wantSize uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotSize uint64) error {
		if !reflect.DeepEqual(gotSize, w.wantSize) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSize, w.wantSize)
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
			p := &pool{
				size:      test.fields.size,
				length:    test.fields.length,
				limit:     test.fields.limit,
				allocator: test.fields.allocator,
				extender:  test.fields.extender,
				flusher:   test.fields.flusher,
				pool:      test.fields.pool,
			}

			gotSize := p.Size(test.args.ctx)
			if err := test.checkFunc(test.want, gotSize); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_extendSize(t *testing.T) {
	type args struct {
		ctx  context.Context
		size uint64
	}
	type fields struct {
		size      uint64
		length    uint64
		limit     uint64
		allocator Allocator
		extender  Extender
		flusher   Flusher
		pool      sync.Pool
	}
	type want struct {
		wantRet uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRet uint64) error {
		if !reflect.DeepEqual(gotRet, w.wantRet) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRet, w.wantRet)
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
		           size: 0,
		       },
		       fields: fields {
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
		           size: 0,
		           },
		           fields: fields {
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
			p := &pool{
				size:      test.fields.size,
				length:    test.fields.length,
				limit:     test.fields.limit,
				allocator: test.fields.allocator,
				extender:  test.fields.extender,
				flusher:   test.fields.flusher,
				pool:      test.fields.pool,
			}

			gotRet := p.extendSize(test.args.ctx, test.args.size)
			if err := test.checkFunc(test.want, gotRet); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_incrementLength(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		size      uint64
		length    uint64
		limit     uint64
		allocator Allocator
		extender  Extender
		flusher   Flusher
		pool      sync.Pool
	}
	type want struct {
		wantSize uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotSize uint64) error {
		if !reflect.DeepEqual(gotSize, w.wantSize) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSize, w.wantSize)
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
			p := &pool{
				size:      test.fields.size,
				length:    test.fields.length,
				limit:     test.fields.limit,
				allocator: test.fields.allocator,
				extender:  test.fields.extender,
				flusher:   test.fields.flusher,
				pool:      test.fields.pool,
			}

			gotSize := p.incrementLength(test.args.ctx)
			if err := test.checkFunc(test.want, gotSize); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_decrementLength(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		size      uint64
		length    uint64
		limit     uint64
		allocator Allocator
		extender  Extender
		flusher   Flusher
		pool      sync.Pool
	}
	type want struct {
		wantSize uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotSize uint64) error {
		if !reflect.DeepEqual(gotSize, w.wantSize) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSize, w.wantSize)
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
			p := &pool{
				size:      test.fields.size,
				length:    test.fields.length,
				limit:     test.fields.limit,
				allocator: test.fields.allocator,
				extender:  test.fields.extender,
				flusher:   test.fields.flusher,
				pool:      test.fields.pool,
			}

			gotSize := p.decrementLength(test.args.ctx)
			if err := test.checkFunc(test.want, gotSize); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Len(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		size      uint64
		length    uint64
		limit     uint64
		allocator Allocator
		extender  Extender
		flusher   Flusher
		pool      sync.Pool
	}
	type want struct {
		wantSize uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotSize uint64) error {
		if !reflect.DeepEqual(gotSize, w.wantSize) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSize, w.wantSize)
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
			p := &pool{
				size:      test.fields.size,
				length:    test.fields.length,
				limit:     test.fields.limit,
				allocator: test.fields.allocator,
				extender:  test.fields.extender,
				flusher:   test.fields.flusher,
				pool:      test.fields.pool,
			}

			gotSize := p.Len(test.args.ctx)
			if err := test.checkFunc(test.want, gotSize); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_pool_Limit(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		size      uint64
		length    uint64
		limit     uint64
		allocator Allocator
		extender  Extender
		flusher   Flusher
		pool      sync.Pool
	}
	type want struct {
		wantSize uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotSize uint64) error {
		if !reflect.DeepEqual(gotSize, w.wantSize) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSize, w.wantSize)
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
		           size: 0,
		           length: 0,
		           limit: 0,
		           allocator: nil,
		           extender: nil,
		           flusher: nil,
		           pool: sync.Pool{},
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
			p := &pool{
				size:      test.fields.size,
				length:    test.fields.length,
				limit:     test.fields.limit,
				allocator: test.fields.allocator,
				extender:  test.fields.extender,
				flusher:   test.fields.flusher,
				pool:      test.fields.pool,
			}

			gotSize := p.Limit(test.args.ctx)
			if err := test.checkFunc(test.want, gotSize); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_allocator_New(t *testing.T) {
	type args struct {
		ctx  context.Context
		size uint64
	}
	type fields struct {
		f func(size uint64) interface{}
	}
	type want struct {
		wantData interface{}
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotData interface{}) error {
		if !reflect.DeepEqual(gotData, w.wantData) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
		           size: 0,
		       },
		       fields: fields {
		           f: nil,
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
		           size: 0,
		           },
		           fields: fields {
		           f: nil,
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
			a := &allocator{
				f: test.fields.f,
			}

			gotData := a.New(test.args.ctx, test.args.size)
			if err := test.checkFunc(test.want, gotData); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_extender_Extend(t *testing.T) {
	type args struct {
		ctx  context.Context
		data interface{}
		size uint64
	}
	type fields struct {
		f func(data interface{}, size uint64) interface{}
	}
	type want struct {
		wantEdata interface{}
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotEdata interface{}) error {
		if !reflect.DeepEqual(gotEdata, w.wantEdata) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotEdata, w.wantEdata)
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
		           data: nil,
		           size: 0,
		       },
		       fields: fields {
		           f: nil,
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
		           data: nil,
		           size: 0,
		           },
		           fields: fields {
		           f: nil,
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
			e := &extender{
				f: test.fields.f,
			}

			gotEdata := e.Extend(test.args.ctx, test.args.data, test.args.size)
			if err := test.checkFunc(test.want, gotEdata); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_flusher_Flush(t *testing.T) {
	type args struct {
		ctx  context.Context
		data interface{}
	}
	type fields struct {
		f func(data interface{}) interface{}
	}
	type want struct {
		wantFdata interface{}
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotFdata interface{}) error {
		if !reflect.DeepEqual(gotFdata, w.wantFdata) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotFdata, w.wantFdata)
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
		           data: nil,
		       },
		       fields: fields {
		           f: nil,
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
		           data: nil,
		           },
		           fields: fields {
		           f: nil,
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
			f := &flusher{
				f: test.fields.f,
			}

			gotFdata := f.Flush(test.args.ctx, test.args.data)
			if err := test.checkFunc(test.want, gotFdata); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
