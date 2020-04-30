//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package worker provides worker processes
package worker

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"

	"go.uber.org/goleak"
)

func TestNewQueue(t *testing.T) {
	type args struct {
		opts []QueueOption
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := NewQueue(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_queue_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		buffer  int
		eg      errgroup.Group
		qcdur   time.Duration
		inCh    chan JobFunc
		outCh   chan JobFunc
		qLen    atomic.Value
		running atomic.Value
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
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
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			q := &queue{
				buffer:  test.fields.buffer,
				eg:      test.fields.eg,
				qcdur:   test.fields.qcdur,
				inCh:    test.fields.inCh,
				outCh:   test.fields.outCh,
				qLen:    test.fields.qLen,
				running: test.fields.running,
			}

			got, err := q.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_queue_isRunning(t *testing.T) {
	type fields struct {
		buffer  int
		eg      errgroup.Group
		qcdur   time.Duration
		inCh    chan JobFunc
		outCh   chan JobFunc
		qLen    atomic.Value
		running atomic.Value
	}
	type want struct {
		want bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			q := &queue{
				buffer:  test.fields.buffer,
				eg:      test.fields.eg,
				qcdur:   test.fields.qcdur,
				inCh:    test.fields.inCh,
				outCh:   test.fields.outCh,
				qLen:    test.fields.qLen,
				running: test.fields.running,
			}

			got := q.isRunning()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_queue_Push(t *testing.T) {
	type args struct {
		ctx context.Context
		job JobFunc
	}
	type fields struct {
		buffer  int
		eg      errgroup.Group
		qcdur   time.Duration
		inCh    chan JobFunc
		outCh   chan JobFunc
		qLen    atomic.Value
		running atomic.Value
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		           job: nil,
		       },
		       fields: fields {
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
		           job: nil,
		           },
		           fields: fields {
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			q := &queue{
				buffer:  test.fields.buffer,
				eg:      test.fields.eg,
				qcdur:   test.fields.qcdur,
				inCh:    test.fields.inCh,
				outCh:   test.fields.outCh,
				qLen:    test.fields.qLen,
				running: test.fields.running,
			}

			err := q.Push(test.args.ctx, test.args.job)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_queue_Pop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		buffer  int
		eg      errgroup.Group
		qcdur   time.Duration
		inCh    chan JobFunc
		outCh   chan JobFunc
		qLen    atomic.Value
		running atomic.Value
	}
	type want struct {
		want JobFunc
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, JobFunc, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got JobFunc, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
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
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			q := &queue{
				buffer:  test.fields.buffer,
				eg:      test.fields.eg,
				qcdur:   test.fields.qcdur,
				inCh:    test.fields.inCh,
				outCh:   test.fields.outCh,
				qLen:    test.fields.qLen,
				running: test.fields.running,
			}

			got, err := q.Pop(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_queue_pop(t *testing.T) {
	type args struct {
		ctx   context.Context
		retry uint64
	}
	type fields struct {
		buffer  int
		eg      errgroup.Group
		qcdur   time.Duration
		inCh    chan JobFunc
		outCh   chan JobFunc
		qLen    atomic.Value
		running atomic.Value
	}
	type want struct {
		want JobFunc
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, JobFunc, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got JobFunc, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
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
		           retry: 0,
		       },
		       fields: fields {
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
		           retry: 0,
		           },
		           fields: fields {
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			q := &queue{
				buffer:  test.fields.buffer,
				eg:      test.fields.eg,
				qcdur:   test.fields.qcdur,
				inCh:    test.fields.inCh,
				outCh:   test.fields.outCh,
				qLen:    test.fields.qLen,
				running: test.fields.running,
			}

			got, err := q.pop(test.args.ctx, test.args.retry)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_queue_Len(t *testing.T) {
	type fields struct {
		buffer  int
		eg      errgroup.Group
		qcdur   time.Duration
		inCh    chan JobFunc
		outCh   chan JobFunc
		qLen    atomic.Value
		running atomic.Value
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
		           buffer: 0,
		           eg: nil,
		           qcdur: nil,
		           inCh: nil,
		           outCh: nil,
		           qLen: nil,
		           running: nil,
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			q := &queue{
				buffer:  test.fields.buffer,
				eg:      test.fields.eg,
				qcdur:   test.fields.qcdur,
				inCh:    test.fields.inCh,
				outCh:   test.fields.outCh,
				qLen:    test.fields.qLen,
				running: test.fields.running,
			}

			got := q.Len()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
