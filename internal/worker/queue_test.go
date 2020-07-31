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

	"github.com/google/go-cmp/cmp"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
		goleak.IgnoreTopFunction("github.com/vdaas/vald/internal/worker.(*queue).Start.func1"),
	}
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
		opts := []cmp.Option{
			cmp.Comparer(func(want, got Queue) bool {
				return want != nil && got != nil
			}),
		}
		if diff := cmp.Diff(w.want, got, opts...); diff != "" {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			q := new(queue)
			for _, opt := range defaultQueueOpts {
				_ = opt(q)
			}
			q.qLen.Store(uint64(0))
			q.running.Store(false)
			q.inCh = make(chan JobFunc, q.buffer)
			q.outCh = make(chan JobFunc, 1)
			return test{
				name: "set success",
				want: want{
					want: q,
				},
			}
		}(),
		func() test {
			q := new(queue)
			opts := []QueueOption{
				WithQueueCheckDuration("invalid"),
			}
			var err error
			for _, opt := range opts {
				err = opt(q)
			}
			return test{
				name: "got error when opts is invalid.",
				args: args{
					opts: opts,
				},
				want: want{
					err: err,
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
		opts := []cmp.Option{
			cmp.Comparer(func(want, got <-chan error) bool {
				return (want != nil && got != nil) || (want == nil && got == nil)
			}),
		}
		if diff := cmp.Diff(w.want, got, opts...); diff != "" {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			inCh := make(chan JobFunc, 10)
			for i := 0; i < 10; i++ {
				inCh <- func(context.Context) error {
					return nil
				}
			}
			return test{
				name: "Start success.",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					buffer: 10,
					eg:     errgroup.Get(),
					qcdur:  100 * time.Microsecond,
					inCh:   inCh,
					outCh:  make(chan JobFunc, 1),
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
				},
				want: want{
					want: make(chan error, 1),
				},
			}
		}(),
		func() test {
			return test{
				name: "Start failed when queue is already running.",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					buffer: 0,
					eg:     errgroup.Get(),
					qcdur:  100 * time.Microsecond,
					inCh:   make(chan JobFunc),
					outCh:  make(chan JobFunc, 1),
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
				},
				want: want{
					err: errors.ErrQueueIsAlreadyRunning(),
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			return test{
				name: "Start failed when ctx.Done before inCh send.",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					buffer: 10,
					eg:     errgroup.Get(),
					qcdur:  100 * time.Microsecond,
					inCh:   make(chan JobFunc, 10),
					outCh:  make(chan JobFunc, 1),
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
				},
				want: want{
					want: make(chan error, 1),
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			inCh := make(chan JobFunc, 10)
			inCh <- func(context.Context) error {
				return nil
			}
			cancel()
			return test{
				name: "Start failed when ctx.Done after inCh send.",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					buffer: 10,
					eg:     errgroup.Get(),
					qcdur:  100 * time.Microsecond,
					inCh:   inCh,
					outCh:  make(chan JobFunc),
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
				},
				want: want{
					want: make(chan error, 1),
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
		{
			name: "Get true when queue is running",
			fields: fields{
				running: func() atomic.Value {
					v := new(atomic.Value)
					v.Store(true)
					return *v
				}(),
			},
			want: want{
				want: true,
			},
			checkFunc: defaultCheckFunc,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
		func() test {
			return test{
				name: "return nil when push success.",
				args: args{
					ctx: context.Background(),
					job: func(context.Context) error {
						return nil
					},
				},
				fields: fields{
					buffer: 10,
					eg:     nil,
					qcdur:  100 * time.Microsecond,
					inCh:   make(chan JobFunc, 10),
					outCh:  make(chan JobFunc),
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			return test{
				name: "return error when job is nil.",
				args: args{
					ctx: context.Background(),
					job: nil,
				},
				fields: fields{
					buffer: 0,
					eg:     nil,
					qcdur:  100 * time.Microsecond,
					inCh:   nil,
					outCh:  nil,
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
				},
				want: want{
					err: errors.ErrJobFuncIsNil(),
				},
			}
		}(),
		func() test {
			return test{
				name: "return error when queue is not running.",
				args: args{
					ctx: context.Background(),
					job: func(context.Context) error {
						return nil
					},
				},
				fields: fields{
					buffer: 0,
					eg:     nil,
					qcdur:  100 * time.Microsecond,
					inCh:   nil,
					outCh:  nil,
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
				},
				want: want{
					err: errors.ErrQueueIsNotRunning(),
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			return test{
				name: "return error when ctx.Done.",
				args: args{
					ctx: ctx,
					job: func(context.Context) error {
						return nil
					},
				},
				fields: fields{
					buffer: 1,
					eg:     errgroup.Get(),
					qcdur:  100 * time.Microsecond,
					inCh:   make(chan JobFunc, 1),
					outCh:  make(chan JobFunc, 0),
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
				},
				want: want{
					err: context.Canceled,
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
		if reflect.ValueOf(w.want).Pointer() != reflect.ValueOf(got).Pointer() {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx := context.Background()
			f := JobFunc(func(context.Context) error {
				return nil
			})
			outCh := make(chan JobFunc, 1)
			outCh <- f
			return test{
				name: "return (JobFunc, nil) when pop success.",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					buffer: 10,
					eg:     errgroup.Get(),
					qcdur:  100 * time.Microsecond,
					inCh:   make(chan JobFunc, 10),
					outCh:  outCh,
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(1))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
				},
				want: want{
					want: f,
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
		if reflect.ValueOf(w.want).Pointer() != reflect.ValueOf(got).Pointer() {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return (nil, error) when queue is not running.",
			args: args{
				ctx:   context.Background(),
				retry: 1,
			},
			fields: fields{
				buffer: 0,
				eg:     errgroup.Get(),
				qcdur:  1 * time.Microsecond,
				inCh:   make(chan JobFunc, 1),
				outCh:  make(chan JobFunc, 1),
				qLen: func() atomic.Value {
					v := new(atomic.Value)
					v.Store(uint64(0))
					return *v
				}(),
				running: func() atomic.Value {
					v := new(atomic.Value)
					v.Store(false)
					return *v
				}(),
			},
			want: want{
				want: nil,
				err:  errors.ErrQueueIsNotRunning(),
			},
		},
		func() test {
			ctx := context.Background()
			f := JobFunc(func(context.Context) error {
				return nil
			})
			outCh := make(chan JobFunc, 10)
			outCh <- nil
			for i := 0; i < 9; i++ {
				outCh <- f
			}
			return test{
				name: "return (JobFunc, nil).",
				args: args{
					ctx:   ctx,
					retry: 10,
				},
				fields: fields{
					buffer: 10,
					eg:     errgroup.Get(),
					qcdur:  100 * time.Microsecond,
					inCh:   make(chan JobFunc, 10),
					outCh:  outCh,
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
				},
				want: want{
					want: f,
					err:  nil,
				},
			}
		}(),
		func() test {
			ctx := context.Background()
			f := JobFunc(func(context.Context) error {
				return nil
			})
			outCh := make(chan JobFunc, 10)
			outCh <- nil
			for i := 0; i < 9; i++ {
				outCh <- f
			}
			return test{
				name: "return (nil, error) when retry is 0.",
				args: args{
					ctx:   ctx,
					retry: 0,
				},
				fields: fields{
					buffer: 10,
					eg:     errgroup.Get(),
					qcdur:  100 * time.Microsecond,
					inCh:   make(chan JobFunc, 10),
					outCh:  outCh,
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
				},
				want: want{
					want: nil,
					err:  errors.ErrJobFuncIsNil(),
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			return test{
				name: "return (JobFunc, nil).",
				args: args{
					ctx:   ctx,
					retry: 0,
				},
				fields: fields{
					buffer: 10,
					eg:     errgroup.Get(),
					qcdur:  100 * time.Microsecond,
					inCh:   make(chan JobFunc, 10),
					outCh:  make(chan JobFunc, 0),
					qLen: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(uint64(0))
						return *v
					}(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(true)
						return *v
					}(),
				},
				want: want{
					want: nil,
					err:  context.Canceled,
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
		qLen atomic.Value
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
		{
			name: "Get qLen when qLen is stored.",
			fields: fields{
				qLen: func() atomic.Value {
					v := new(atomic.Value)
					v.Store(uint64(0))
					return *v
				}(),
			},
			want: want{
				want: 0,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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
				qLen: test.fields.qLen,
			}

			got := q.Len()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
