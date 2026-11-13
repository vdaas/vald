// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package worker

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Worker
		err  error
	}
	type test struct {
		want       want
		checkFunc  func(want, Worker, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
			comparator.AllowUnexported(*want.queue.(*queue)),
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
				opts: []Option{
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		want <-chan error
		err  error
	}
	type test struct {
		want       want
		args       args
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		fields     fields
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
			ctx, cancel := context.WithCancel(t.Context())
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
				ctx: t.Context(),
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		want <-chan error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, <-chan error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		fields     fields
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
			ctx, cancel := context.WithCancel(t.Context())
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
			ctx, cancel := context.WithCancel(t.Context())
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
						if !errors.Is(e, err) {
							return errors.New("invalid error")
						}
					}
					return nil
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
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
			ctx, cancel := context.WithCancel(t.Context())
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
						if !errors.Is(e, err) {
							return errors.New("invalid error")
						}
					}
					return nil
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct{}
	type test struct {
		want       want
		checkFunc  func(want, *worker) error
		beforeFunc func()
		afterFunc  func()
		name       string
		fields     fields
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct{}
	type test struct {
		want       want
		checkFunc  func(want, *worker) error
		beforeFunc func()
		afterFunc  func()
		name       string
		fields     fields
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		want bool
	}
	type test struct {
		checkFunc  func(want, bool) error
		beforeFunc func()
		afterFunc  func()
		name       string
		fields     fields
		want       want
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		want string
	}
	type test struct {
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
		name       string
		want       want
		fields     fields
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		want uint64
	}
	type test struct {
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
		name       string
		fields     fields
		want       want
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		want uint64
	}
	type test struct {
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
		name       string
		fields     fields
		want       want
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		want uint64
	}
	type test struct {
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
		name       string
		fields     fields
		want       want
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
		running        atomic.Value
		eg             errgroup.Group
		queue          Queue
		name           string
		qopts          []QueueOption
		limitation     int
		requestedCount uint64
		completedCount uint64
	}
	type want struct {
		err error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(*worker, want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		fields     fields
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
				ctx: t.Context(),
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
				ctx: t.Context(),
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
				ctx: t.Context(),
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
				ctx: t.Context(),
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

// TestWorker_Start_notRunningAfterStop checks that the running flag is false once the worker has stopped.
// The flag must be set before the aggregator goroutine starts, otherwise its deferred Store(false) may run first.
func TestWorker_Start_notRunningAfterStop(t *testing.T) {
	t.Parallel()
	for i := range 100 {
		w, err := New()
		if err != nil {
			t.Fatal(err)
		}
		// the ctx is cancelled before Start so that the aggregator goroutine may exit before Start returns.
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		ech, err := w.Start(ctx)
		if err != nil {
			t.Fatalf("iter %d: Start returned error: %v", i, err)
		}
		if !waitClosed(ech) {
			t.Fatalf("iter %d: worker did not stop in time", i)
		}
		if w.IsRunning() {
			t.Fatalf("iter %d: worker is reported as running after it stopped", i)
		}
	}
}

// TestWorker_startJobLoop_noLeakWhenErrorChNotRead checks that the job loop terminates on ctx cancel
// even if nobody reads the error channel, i.e. errors are not sent with a bare blocking send.
func TestWorker_startJobLoop_noLeakWhenErrorChNotRead(t *testing.T) {
	t.Parallel()
	for i := range 50 {
		eg, _ := errgroup.New(t.Context())
		w, err := New(WithLimitation(1), WithErrGroup(eg))
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(t.Context())
		if _, err := w.(*worker).queue.Start(ctx); err != nil {
			cancel()
			t.Fatalf("iter %d: queue Start returned error: %v", i, err)
		}
		// the returned channel is intentionally never read until the check has finished.
		wech := w.(*worker).startJobLoop(ctx)
		time.Sleep(time.Millisecond)
		cancel()

		var werr error
		done := make(chan struct{})
		go func() {
			werr = eg.Wait()
			close(done)
		}()
		select {
		case <-done:
			if werr != nil && !errors.Is(werr, context.Canceled) {
				t.Errorf("iter %d: unexpected error from errgroup: %v", i, werr)
			}
		case <-time.After(5 * time.Second):
			// release the blocked sender so that the goroutine does not outlive the test.
			go func() {
				for range wech {
				}
			}()
			t.Fatalf("iter %d: job loop did not stop after ctx cancel while the error channel was not read", i)
		}
	}
}

// newJobLoopWorker returns a worker whose queue is replaced by a mock that pops with pop,
// and a fresh errgroup that tracks every goroutine the worker starts.
func newJobLoopWorker(
	t *testing.T, limit int, pop func(context.Context) (JobFunc, error),
) (*worker, errgroup.Group) {
	t.Helper()
	eg, _ := errgroup.New(t.Context())
	w, err := New(WithLimitation(limit), WithErrGroup(eg))
	if err != nil {
		t.Fatal(err)
	}
	wk, ok := w.(*worker)
	if !ok {
		t.Fatalf("unexpected worker type %T", w)
	}
	wk.queue = &QueueMock{PopFunc: pop}
	return wk, eg
}

// waitGroupDone returns a channel that is closed when eg.Wait returns.
// The Wait error is reported through the returned func after the channel is closed.
func waitGroupDone(eg errgroup.Group) (<-chan struct{}, func() error) {
	var werr error
	done := make(chan struct{})
	go func() {
		werr = eg.Wait()
		close(done)
	}()
	return done, func() error { return werr }
}

// TestWorker_startJobLoop_noLeakWhenJobErrorNotRead checks that a failing job does not block forever on the
// error channel after the ctx passed to Start is cancelled and nobody reads the channel.
func TestWorker_startJobLoop_noLeakWhenJobErrorNotRead(t *testing.T) {
	t.Parallel()
	errJob := errors.New("job failed")
	for i := range 50 {
		var calls atomic.Int64
		reached := make(chan struct{})
		w, eg := newJobLoopWorker(t, 1, func(context.Context) (JobFunc, error) {
			if calls.Add(1) == 2 {
				close(reached)
			}
			return func(context.Context) error { return errJob }, nil
		})
		ctx, cancel := context.WithCancel(t.Context())
		wech := w.startJobLoop(ctx)
		done, werr := waitGroupDone(eg)

		// the first failing job fills the unread error channel (capacity 1); a job error does not stop the loop,
		// so it pops a second job whose error send has no reader and blocks while ctx is alive.
		// Wait for that (or for the loop to stop), then cancel ctx: the loop must stop even if the second send is blocked.
		select {
		case <-reached:
		case <-done:
		case <-time.After(5 * time.Second):
			cancel()
			t.Fatalf("iter %d: job loop neither progressed nor stopped", i)
		}
		cancel()

		select {
		case <-done:
			if err := werr(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, errJob) {
				t.Errorf("iter %d: unexpected error from errgroup: %v", i, err)
			}
		case <-time.After(5 * time.Second):
			// release the blocked sender so that the goroutine does not outlive the test.
			go func() {
				for range wech {
				}
			}()
			t.Fatalf("iter %d: job loop did not stop after ctx cancel while the error channel was not read", i)
		}
	}
}

// TestWorker_startJobLoop_deliversAllJobErrors checks that every job error reaches the error channel
// while the ctx passed to Start is alive.
func TestWorker_startJobLoop_deliversAllJobErrors(t *testing.T) {
	t.Parallel()
	const jobs = 3
	jobErrs := make([]error, jobs)
	for i := range jobErrs {
		jobErrs[i] = errors.Errorf("job error %d", i)
	}
	for i := range 50 {
		var popped, ready atomic.Int64
		gate := make(chan struct{})
		w, eg := newJobLoopWorker(t, jobs, func(ctx context.Context) (JobFunc, error) {
			n := int(popped.Add(1)) - 1
			if n >= jobs {
				<-ctx.Done()
				return nil, ctx.Err()
			}
			return func(jctx context.Context) error {
				// all jobs fail at the same time, after every job is running.
				if ready.Add(1) == jobs {
					close(gate)
				}
				select {
				case <-gate:
				case <-jctx.Done():
				}
				return jobErrs[n]
			}, nil
		})
		ctx, cancel := context.WithCancel(t.Context())
		wech := w.startJobLoop(ctx)
		done, _ := waitGroupDone(eg)

		got := make(map[int]struct{}, jobs)
		timeout := time.After(5 * time.Second)
	collect:
		for len(got) < jobs {
			select {
			case err, ok := <-wech:
				if !ok {
					break collect
				}
				for j, want := range jobErrs {
					if errors.Is(err, want) {
						got[j] = struct{}{}
					}
				}
			case <-timeout:
				break collect
			}
		}
		cancel()
		if len(got) != jobs {
			t.Errorf("iter %d: got %d distinct job errors, want %d", i, len(got), jobs)
		}
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			go func() {
				for range wech {
				}
			}()
			t.Fatalf("iter %d: job loop did not stop after ctx cancel", i)
		}
	}
}

// NOT IMPLEMENTED BELOW

// TestWorker_Start_keepsRunningAfterJobError checks that a failing job is reported on the error channel
// but does not stop the worker: the next job still runs with a live ctx, the error channel stays open and
// the worker is reported as running until the ctx passed to Start is cancelled.
// Before the fix the first job error cancelled the job loop's errgroup ctx, so the loop stopped (closing the
// error channel) or ran the next job with an already cancelled ctx.
func TestWorker_Start_keepsRunningAfterJobError(t *testing.T) {
	t.Parallel()
	errJob := errors.New("job failed")
	for i := range 50 {
		var popped atomic.Int64
		ran := make(chan error, 1)
		w, eg := newJobLoopWorker(t, 1, func(ctx context.Context) (JobFunc, error) {
			switch popped.Add(1) {
			case 1:
				return func(context.Context) error { return errJob }, nil
			case 2:
				return func(jctx context.Context) error {
					ran <- jctx.Err()
					return nil
				}, nil
			}
			<-ctx.Done()
			return nil, ctx.Err()
		})
		qech := make(chan error)
		w.queue.(*QueueMock).StartFunc = func(context.Context) (<-chan error, error) {
			return qech, nil
		}
		ctx, cancel := context.WithCancel(t.Context())
		ech, err := w.Start(ctx)
		if err != nil {
			cancel()
			t.Fatalf("iter %d: Start returned error: %v", i, err)
		}
		done, _ := waitGroupDone(eg)
		fail := func(format string, args ...any) {
			cancel()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
			}
			t.Fatalf(format, args...)
		}

		select {
		case err, ok := <-ech:
			if !ok {
				fail("iter %d: error channel closed before the job error was reported", i)
			}
			if !errors.Is(err, errJob) {
				fail("iter %d: got error %v, want %v", i, err, errJob)
			}
		case <-time.After(5 * time.Second):
			fail("iter %d: job error was not reported", i)
		}

		select {
		case jerr := <-ran:
			if jerr != nil {
				fail("iter %d: next job ran with a done ctx: %v", i, jerr)
			}
		case err, ok := <-ech:
			if !ok {
				fail("iter %d: error channel closed after a job error; the worker stopped", i)
			}
			fail("iter %d: unexpected error on the error channel: %v", i, err)
		case <-time.After(5 * time.Second):
			fail("iter %d: next job did not run after a job error", i)
		}
		if !w.IsRunning() {
			fail("iter %d: worker is not reported as running after a job error", i)
		}

		cancel()
		if !waitClosed(ech) {
			t.Fatalf("iter %d: error channel was not closed after ctx cancel", i)
		}
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatalf("iter %d: worker goroutines did not stop after ctx cancel", i)
		}
		if w.IsRunning() {
			t.Fatalf("iter %d: worker is reported as running after ctx cancel", i)
		}
	}
}
