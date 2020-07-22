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
	}
)

func TestNew(t *testing.T) {
	type args struct {
		opts []WorkerOption
	}
	type want struct {
		want Worker
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Worker, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Worker, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}

		want := w.want.(*worker)

		avComparer := func(x, y atomic.Value) bool {
			return reflect.DeepEqual(x.Load(), y.Load())
		}
		egComparer := func(x, y errgroup.Group) bool {
			return reflect.DeepEqual(x, y)
		}

		queueOpts := []cmp.Option{
			cmp.AllowUnexported(*(want.queue.(*queue))),
			cmp.Comparer(func(x, y chan JobFunc) bool {
				return len(x) == len(y)
			}),
			cmp.Comparer(egComparer),
			cmp.Comparer(avComparer),
		}

		opts := []cmp.Option{
			cmp.AllowUnexported(*want),
			cmp.Comparer(func(x, y queue) bool {
				return cmp.Equal(x, y, queueOpts...)
			}),
			cmp.Comparer(avComparer),
			cmp.Comparer(egComparer),
		}
		if diff := cmp.Diff(want, got, opts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "return worker without option",
			want: want{
				want: &worker{
					name:       "worker",
					limitation: 10,
					eg:         errgroup.Get(),
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
					queue: &queue{
						buffer: 10,
						eg:     errgroup.Get(),
						qcdur:  200 * time.Millisecond,
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
						inCh:  make(chan JobFunc, 10),
						outCh: make(chan JobFunc, 1),
					},
				},
			},
		},
		{
			name: "return worker with option",
			args: args{
				opts: []WorkerOption{
					WithName("test1"),
				},
			},
			want: want{
				want: &worker{
					name:       "test1",
					limitation: 10,
					running: func() atomic.Value {
						v := new(atomic.Value)
						v.Store(false)
						return *v
					}(),
					eg: errgroup.Get(),
					queue: &queue{
						buffer: 10,
						eg:     errgroup.Get(),
						qcdur:  200 * time.Millisecond,
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
						inCh:  make(chan JobFunc, 10),
						outCh: make(chan JobFunc, 1),
					},
				},
			},
		},
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

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_worker_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			got, err := w.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_worker_startJobLoop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		want <-chan error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error) error {
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			got := w.startJobLoop(test.args.ctx)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_worker_Pause(t *testing.T) {
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
	}
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			w.Pause()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_worker_Resume(t *testing.T) {
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
	}
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			w.Resume()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_worker_IsRunning(t *testing.T) {
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			got := w.IsRunning()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_worker_Name(t *testing.T) {
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			got := w.Name()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_worker_Len(t *testing.T) {
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			got := w.Len()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_worker_TotalRequested(t *testing.T) {
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			got := w.TotalRequested()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_worker_TotalCompleted(t *testing.T) {
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			got := w.TotalCompleted()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_worker_Dispatch(t *testing.T) {
	type args struct {
		ctx context.Context
		f   JobFunc
	}
	type fields struct {
		name           string
		limitation     int
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		qopts          []QueueOption
		requestedCount uint64
		completedCount uint64
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
		           f: nil,
		       },
		       fields: fields {
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
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
		           name: "",
		           limitation: 0,
		           running: nil,
		           eg: nil,
		           queue: nil,
		           qopts: nil,
		           requestedCount: 0,
		           completedCount: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &worker{
				name:           test.fields.name,
				limitation:     test.fields.limitation,
				running:        test.fields.running,
				eg:             test.fields.eg,
				queue:          test.fields.queue,
				qopts:          test.fields.qopts,
				requestedCount: test.fields.requestedCount,
				completedCount: test.fields.completedCount,
			}

			err := w.Dispatch(test.args.ctx, test.args.f)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
