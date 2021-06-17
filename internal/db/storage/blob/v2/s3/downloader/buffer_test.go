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
package downloader

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func Test_newBuffer(t *testing.T) {
	type args struct {
		ctx       context.Context
		chunkSize int
	}
	type want struct {
		want ReadWriterAtCloserBuffer
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, ReadWriterAtCloserBuffer) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got ReadWriterAtCloserBuffer) error {
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
		           chunkSize: 0,
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
		           chunkSize: 0,
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

			got := newBuffer(test.args.ctx, test.args.chunkSize)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_buffer_WriteAt(t *testing.T) {
	type args struct {
		p   []byte
		pos int64
	}
	type fields struct {
		mu     sync.RWMutex
		data   *data
		pos    int64
		cur    int64
		nextCh chan int64
		ctx    context.Context
		cancel context.CancelFunc
	}
	type want struct {
		wantN int
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotN int, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           p: nil,
		           pos: 0,
		       },
		       fields: fields {
		           mu: sync.RWMutex{},
		           data: data{},
		           pos: 0,
		           cur: 0,
		           nextCh: nil,
		           ctx: nil,
		           cancel: nil,
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
		           p: nil,
		           pos: 0,
		           },
		           fields: fields {
		           mu: sync.RWMutex{},
		           data: data{},
		           pos: 0,
		           cur: 0,
		           nextCh: nil,
		           ctx: nil,
		           cancel: nil,
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
			b := &buffer{
				mu:     test.fields.mu,
				data:   test.fields.data,
				pos:    test.fields.pos,
				cur:    test.fields.cur,
				nextCh: test.fields.nextCh,
				ctx:    test.fields.ctx,
				cancel: test.fields.cancel,
			}

			gotN, err := b.WriteAt(test.args.p, test.args.pos)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_buffer_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	type fields struct {
		mu     sync.RWMutex
		data   *data
		pos    int64
		cur    int64
		nextCh chan int64
		ctx    context.Context
		cancel context.CancelFunc
	}
	type want struct {
		wantN int
		err   error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, int, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotN int, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotN, w.wantN) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           p: nil,
		       },
		       fields: fields {
		           mu: sync.RWMutex{},
		           data: data{},
		           pos: 0,
		           cur: 0,
		           nextCh: nil,
		           ctx: nil,
		           cancel: nil,
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
		           p: nil,
		           },
		           fields: fields {
		           mu: sync.RWMutex{},
		           data: data{},
		           pos: 0,
		           cur: 0,
		           nextCh: nil,
		           ctx: nil,
		           cancel: nil,
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
			b := &buffer{
				mu:     test.fields.mu,
				data:   test.fields.data,
				pos:    test.fields.pos,
				cur:    test.fields.cur,
				nextCh: test.fields.nextCh,
				ctx:    test.fields.ctx,
				cancel: test.fields.cancel,
			}

			gotN, err := b.Read(test.args.p)
			if err := test.checkFunc(test.want, gotN, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_buffer_Close(t *testing.T) {
	type fields struct {
		mu     sync.RWMutex
		data   *data
		pos    int64
		cur    int64
		nextCh chan int64
		ctx    context.Context
		cancel context.CancelFunc
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
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
		       fields: fields {
		           mu: sync.RWMutex{},
		           data: data{},
		           pos: 0,
		           cur: 0,
		           nextCh: nil,
		           ctx: nil,
		           cancel: nil,
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
		           mu: sync.RWMutex{},
		           data: data{},
		           pos: 0,
		           cur: 0,
		           nextCh: nil,
		           ctx: nil,
		           cancel: nil,
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
			b := &buffer{
				mu:     test.fields.mu,
				data:   test.fields.data,
				pos:    test.fields.pos,
				cur:    test.fields.cur,
				nextCh: test.fields.nextCh,
				ctx:    test.fields.ctx,
				cancel: test.fields.cancel,
			}

			err := b.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_data_add(t *testing.T) {
	type args struct {
		n *data
	}
	type fields struct {
		pos  int64
		size int64
		p    []byte
		next *data
	}
	type want struct {
		want *data
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *data) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *data) error {
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
		           n: data{},
		       },
		       fields: fields {
		           pos: 0,
		           size: 0,
		           p: nil,
		           next: data{},
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
		           n: data{},
		           },
		           fields: fields {
		           pos: 0,
		           size: 0,
		           p: nil,
		           next: data{},
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
			d := &data{
				pos:  test.fields.pos,
				size: test.fields.size,
				p:    test.fields.p,
				next: test.fields.next,
			}

			got := d.add(test.args.n)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
