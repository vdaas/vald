//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		egComparator := func(x, y errgroup.Group) bool {
			return reflect.DeepEqual(x, y)
		}
		atomicValueComparator := func(x, y atomic.Value) bool {
			return reflect.DeepEqual(x.Load(), y.Load())
		}
		want := w.want.(*worker)

		queueOpts := []comparator.Option{
			comparator.AllowUnexported(*(want.queue.(*queue))),
			comparator.Comparer(func(x, y chan JobFunc) bool {
				return len(x) == len(y)
			}),
			comparator.Comparer(egComparator),
			comparator.Comparer(atomicValueComparator),
		}
		opts := []comparator.Option{
			comparator.AllowUnexported(*want),
			comparator.Comparer(func(x, y Queue) bool {
				return comparator.Equal(x, y, queueOpts...)
			}),
			comparator.Comparer(egComparator),
			comparator.Comparer(atomicValueComparator),
		}
		if diff := comparator.Diff(want, got, opts...); diff != "" {
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
					running: func() (v atomic.Value) {
						v.Store(false)
						return v
					}(),
					queue: &queue{
						buffer: 10,
						eg:     errgroup.Get(),
						qcdur:  200 * time.Millisecond,
						qLen: func() (v atomic.Value) {
							v.Store(uint64(0))
							return v
						}(),
						running: func() (v atomic.Value) {
							v.Store(false)
							return v
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
					running: func() (v atomic.Value) {
						v.Store(false)
						return v
					}(),
					eg: errgroup.Get(),
					queue: &queue{
						buffer: 10,
						eg:     errgroup.Get(),
						qcdur:  200 * time.Millisecond,
						qLen: func() (v atomic.Value) {
							v.Store(uint64(0))
							return v
						}(),
						running: func() (v atomic.Value) {
							v.Store(false)
							return v
						}(),
						inCh:  make(chan JobFunc, 10),
						outCh: make(chan JobFunc, 1),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		if w.want == nil && got == nil {
			return nil
		}
		if w.want == nil || got == nil || len(w.want) != len(got) {
			return errors.New("want is not equal to got")
		}

		for e := range w.want {
			if e1 := <-got; !errors.Is(e, e1) {
				return errors.New("want is not equal to got")
			}
		}

		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "Start without error",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					name:       "worker",
					limitation: 10,
					eg:         errgroup.Get(),
					running: func() (v atomic.Value) {
						v.Store(false)
						return v
					}(),
					queue: NewQueueMock(),
				},
				want: want{
					want: func() <-chan error {
						ch := make(chan error, 2)
						close(ch)
						return ch
					}(),
				},
				checkFunc: func(w want, got <-chan error, err error) error {
					cancel()
					return defaultCheckFunc(w, got, err)
				},
			}
		}(),
		{
			name: "return error if it is running",
			args: args{},
			fields: fields{
				name: "test",
				running: func() (v atomic.Value) {
					v.Store(true)
					return v
				}(),
			},
			want: want{
				err: errors.ErrWorkerIsAlreadyRunning("test"),
			},
		},
		{
			name: "return queue start error",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				name:       "worker",
				limitation: 10,
				eg:         errgroup.Get(),
				running: func() (v atomic.Value) {
					v.Store(false)
					return v
				}(),
				queue: &QueueMock{
					StartFunc: func(context.Context) (<-chan error, error) {
						return nil, errors.New("error")
					},
				},
			},
			want: want{
				err: errors.New("error"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, got, err); err != nil {
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
		if w.want == nil && got == nil {
			return nil
		}
		if w.want == nil || got == nil || len(w.want) != len(got) {
			return errors.New("want is not equal to got")
		}

		for e := range w.want {
			if e1 := <-got; !errors.Is(e, e1) {
				return errors.New("want is not equal to got")
			}
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "start job loop with empty job list",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					name:           "test",
					limitation:     1,
					running:        atomic.Value{},
					eg:             errgroup.Get(),
					queue:          NewQueueMock(),
					requestedCount: 0,
					completedCount: 0,
				},
				want: want{
					want: func() <-chan error {
						ch := make(chan error, 1)
						close(ch)
						return ch
					}(),
				},
				checkFunc: func(w want, got <-chan error) error {
					time.Sleep(time.Millisecond * 200)
					cancel()
					time.Sleep(time.Millisecond * 200)

					return defaultCheckFunc(w, got)
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			err := errors.New("error")
			return test{
				name: "start job loop with queue pop error",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					name:       "test",
					limitation: 1,
					running:    atomic.Value{},
					eg:         errgroup.Get(),
					queue: &QueueMock{
						StartFunc: DefaultStartFunc,
						PopFunc: func(context.Context) (JobFunc, error) {
							return nil, err
						},
					},
					requestedCount: 0,
					completedCount: 0,
				},
				checkFunc: func(w want, got <-chan error) error {
					time.Sleep(time.Millisecond * 200)
					cancel()
					time.Sleep(time.Millisecond * 200)

					if len(got) == 0 {
						return errors.New("got chan len 0")
					}
					for e := range got {
						if e != err {
							return errors.New("invalid error")
						}
					}
					return nil
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "start job loop with a job",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					name:       "test",
					limitation: 1,
					running:    atomic.Value{},
					eg:         errgroup.Get(),
					queue: &QueueMock{
						StartFunc: DefaultStartFunc,
						PopFunc: func(context.Context) (JobFunc, error) {
							f := JobFunc(func(context.Context) error {
								return nil
							})
							return f, nil
						},
					},
					requestedCount: 0,
					completedCount: 0,
				},
				checkFunc: func(w want, got <-chan error) error {
					time.Sleep(time.Millisecond * 200)
					cancel()
					time.Sleep(time.Millisecond * 200)

					if len(got) != 0 {
						return errors.New("error returned")
					}
					return nil
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			err := errors.New("error")
			return test{
				name: "start job loop with a job which return error",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					name:       "test",
					limitation: 1,
					running:    atomic.Value{},
					eg:         errgroup.Get(),
					queue: &QueueMock{
						StartFunc: DefaultStartFunc,
						PopFunc: func(context.Context) (JobFunc, error) {
							f := JobFunc(func(context.Context) error {
								return err
							})
							return f, nil
						},
					},
					requestedCount: 0,
					completedCount: 0,
				},
				checkFunc: func(w want, got <-chan error) error {
					time.Sleep(time.Millisecond * 200)
					cancel()
					time.Sleep(time.Millisecond * 200)

					if len(got) == 0 {
						return errors.New("got chan len 0")
					}
					for e := range got {
						if e != err {
							return errors.New("invalid error")
						}
					}
					return nil
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "start job loop with queue pop a nil job without error",
				args: args{
					ctx: ctx,
				},
				fields: fields{
					name:       "test",
					limitation: 1,
					running:    atomic.Value{},
					eg:         errgroup.Get(),
					queue: &QueueMock{
						StartFunc: DefaultStartFunc,
						PopFunc: func(context.Context) (JobFunc, error) {
							return nil, nil
						},
					},
					requestedCount: 0,
					completedCount: 0,
				},
				checkFunc: func(w want, got <-chan error) error {
					time.Sleep(time.Millisecond * 200)
					cancel()
					time.Sleep(time.Millisecond * 200)

					if len(got) != 0 {
						return errors.New("got error")
					}
					return nil
				},
			}
		}(),
	}

	log.Init(log.WithLoggerType(logger.NOP.String()))
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, got); err != nil {
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
	type want struct{}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *worker) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, worker *worker) error {
		return nil
	}
	tests := []test{
		{
			name: "Pause success",
			fields: fields{
				running: func() (v atomic.Value) {
					v.Store(true)
					return v
				}(),
			},
			checkFunc: func(w want, worker *worker) error {
				if worker.running.Load().(bool) != false {
					return errors.New("running is not false")
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, w); err != nil {
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
	type want struct{}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *worker) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, worker *worker) error {
		return nil
	}
	tests := []test{
		{
			name: "Resume success",
			fields: fields{
				running: func() (v atomic.Value) {
					v.Store(false)
					return v
				}(),
			},
			checkFunc: func(w want, worker *worker) error {
				if worker.running.Load().(bool) != true {
					return errors.New("running is not false")
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, w); err != nil {
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return true if it is running",
			fields: fields{
				running: func() (v atomic.Value) {
					v.Store(true)
					return v
				}(),
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false if it is not running",
			fields: fields{
				running: func() (v atomic.Value) {
					v.Store(false)
					return v
				}(),
			},
			want: want{
				want: false,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, got); err != nil {
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return name",
			fields: fields{
				name: "testname",
			},
			want: want{
				want: "testname",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, got); err != nil {
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return queue length",
			fields: fields{
				queue: &QueueMock{
					LenFunc: func() uint64 {
						return uint64(100)
					},
				},
			},
			want: want{
				want: uint64(100),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, got); err != nil {
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return total count",
			fields: fields{
				requestedCount: 1000,
			},
			want: want{
				want: 1000,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, got); err != nil {
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return total count",
			fields: fields{
				completedCount: 1000,
			},
			want: want{
				want: 1000,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, got); err != nil {
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
		checkFunc  func(*worker, want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(worker *worker, w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "return error if the worker is not started yet",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				name: "test",
				running: func() (v atomic.Value) {
					v.Store(false)
					return v
				}(),
			},
			want: want{
				err: errors.ErrWorkerIsNotRunning("test"),
			},
		},
		{
			name: "return error if the job is failed to push to worker queue",
			args: args{
				ctx: context.Background(),
				f: JobFunc(func(context.Context) error {
					return nil
				}),
			},
			fields: fields{
				name: "test",
				running: func() (v atomic.Value) {
					v.Store(true)
					return v
				}(),
				queue: &QueueMock{
					PushFunc: func(context.Context, JobFunc) error {
						return errors.New("queue push error")
					},
				},
			},
			want: want{
				err: errors.New("queue push error"),
			},
		},
		{
			name: "return nil if the job is nil",
			args: args{
				ctx: context.Background(),
				f:   nil,
			},
			fields: fields{
				name: "test",
				running: func() (v atomic.Value) {
					v.Store(true)
					return v
				}(),
				queue: &QueueMock{},
			},
			want: want{},
		},
		{
			name: "request count is incremented if the job is pushed",
			args: args{
				ctx: context.Background(),
				f: JobFunc(func(context.Context) error {
					return nil
				}),
			},
			fields: fields{
				name: "test",
				running: func() (v atomic.Value) {
					v.Store(true)
					return v
				}(),
				queue: &QueueMock{
					PushFunc: func(context.Context, JobFunc) error {
						return nil
					},
				},
				requestedCount: uint64(999),
			},
			want: want{},
			checkFunc: func(worker *worker, w want, err error) error {
				if worker.requestedCount != uint64(1000) {
					return errors.New("requestedCount is not incremented")
				}
				return defaultCheckFunc(worker, w, err)
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(w, test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
