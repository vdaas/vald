//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package errgroup

import (
	"context"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/semaphore"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	goleak.VerifyTestMain(m)
}

func TestNew(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		want  Group
		want1 context.Context
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Group, context.Context) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got Group, got1 context.Context) error {
		if got, want := got.(*group), w.want.(*group); !reflect.DeepEqual(got.emap, want.emap) && got.cancel != nil {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx := context.Background()
			egctx, cancel := context.WithCancelCause(ctx)

			return test{
				name: "returns (g, ctx)",
				args: args{
					ctx: ctx,
				},
				want: want{
					want: &group{
						egctx:  egctx,
						cancel: cancel,
						emap:   make(map[string]struct{}),
					},
					want1: egctx,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, got1 := New(test.args.ctx)
			if err := checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		wantEgctx context.Context
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, context.Context) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, gotEgctx context.Context) error {
		if !reflect.DeepEqual(gotEgctx, w.wantEgctx) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotEgctx, w.wantEgctx)
		}
		if instance == nil {
			return errors.New("instance is nil")
		}
		return nil
	}
	defaultBeforeFunc := func(args) {
		instance, once = nil, sync.Once{}
	}
	tests := []test{
		func() test {
			ctx := context.Background()
			egctx, cancel := context.WithCancelCause(ctx)

			return test{
				name: "returns egctx when once.Do is called",
				args: args{
					ctx: ctx,
				},
				want: want{
					wantEgctx: egctx,
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					cancel(context.Canceled)
					defaultBeforeFunc(a)
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			test.beforeFunc(test.args)

			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotEgctx := Init(test.args.ctx)
			if err := checkFunc(test.want, gotEgctx); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type want struct {
		want Group
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, Group) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got Group) error {
		if got, want := got.(*group), w.want.(*group); !reflect.DeepEqual(got.egctx, want.egctx) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	initFunc := func() {
		instance, once = nil, sync.Once{}
	}
	defaultBeforeFunc := func(t *testing.T) {
		t.Helper()
		initFunc()
	}
	defaultAfterFunc := func(t *testing.T) {
		t.Helper()
		initFunc()
	}

	tests := []test{
		func() test {
			ctx := context.Background()
			egctx, cancel := context.WithCancelCause(ctx)

			return test{
				name: "returns instance when instance is nil",
				want: want{
					want: &group{
						egctx:  egctx,
						cancel: cancel,
					},
				},
			}
		}(),

		func() test {
			g := &group{
				egctx: context.Background(),
			}
			return test{
				name: "returns instance when instance is not nil",
				want: want{
					want: g,
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					instance = g
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			test.beforeFunc(tt)

			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}
			defer test.afterFunc(tt)

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Get()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGo(t *testing.T) {
	type args struct {
		f func() error
	}
	type test struct {
		name       string
		args       args
		checkFunc  func(Group) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(Group) error {
		return nil
	}
	tests := []test{
		func() test {
			var calledCnt int32

			return test{
				name: "instance.Go is called when instance is not nil",
				args: args{
					f: func() error {
						atomic.AddInt32(&calledCnt, 1)
						return nil
					},
				},
				beforeFunc: func(t *testing.T, _ args) {
					t.Helper()
					g := new(group)
					instance = g
				},
				checkFunc: func(got Group) error {
					if err := got.Wait(); err != nil {
						return err
					}

					if got, want := int(atomic.LoadInt32(&calledCnt)), 1; got != want {
						return errors.Errorf("calledCnt = %v, want: %v", got, want)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			Go(test.args.f)
			if err := checkFunc(instance); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_SetLimit(t *testing.T) {
	type args struct {
		limit int
	}
	type fields struct {
		sem *semaphore.Weighted
	}
	type want struct {
		want Group
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, Group) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, g Group) error {
		got, want := g.(*group), w.want.(*group)
		if (got.sem == nil) != (want.sem == nil) {
			return errors.Errorf("got: \"%#v\",\n\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "set disable when limit is -1",
			args: args{
				limit: -1,
			},
			want: want{
				want: &group{},
			},
		},
		{
			name: "set enable but block all acquireings when limit is 0",
			args: args{
				limit: 0,
			},
			want: want{
				want: &group{
					sem: semaphore.NewWeighted(0),
				},
			},
		},

		{
			name: "set enable when limit is 1",
			args: args{
				limit: 1,
			},
			fields: fields{},
			want: want{
				want: &group{
					sem: semaphore.NewWeighted(1),
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			g := &group{}

			g.SetLimit(test.args.limit)
			if err := checkFunc(test.want, g); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_Go(t *testing.T) {
	type args struct {
		f func() error
	}
	type fields struct {
		egctx  context.Context
		cancel context.CancelCauseFunc
		sem    *semaphore.Weighted
		emap   map[string]struct{}
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func(Group) error
		beforeFunc func(*testing.T, args, Group)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(g Group) error {
		return nil
	}
	tests := []test{
		func() test {
			var calledCnt int32

			egctx, cancel := context.WithCancelCause(context.Background())

			limit := 3

			return test{
				name: "f is not called when reached limit and cancel g.egctx",
				args: args{
					f: func() error {
						atomic.AddInt32(&calledCnt, 1)
						return nil
					},
				},
				fields: fields{
					egctx: egctx,
					sem:   semaphore.NewWeighted(int64(limit)),
				},
				beforeFunc: func(t *testing.T, _ args, g Group) {
					t.Helper()
					for i := 0; i < limit; i++ {
						g.Go(func() error {
							time.Sleep(3 * time.Second)
							return nil
						})
					}
					time.Sleep(time.Second)
					cancel(context.Canceled)
				},
				checkFunc: func(got Group) error {
					if err := got.Wait(); err != nil {
						return err
					}

					if got, want := int(atomic.LoadInt32(&calledCnt)), 0; got != want {
						return errors.Errorf("calledCnt = %v, want: %v", got, want)
					}
					return nil
				},
			}
		}(),

		func() test {
			var calledCnt int32

			egctx, cancel := context.WithCancelCause(context.Background())

			return test{
				name: "f is called but f returns error and previous process also returns error",
				args: args{
					f: func() error {
						atomic.AddInt32(&calledCnt, 1)
						return errors.New("err")
					},
				},
				fields: fields{
					egctx:  egctx,
					cancel: cancel,
					emap:   make(map[string]struct{}),
				},
				beforeFunc: func(t *testing.T, a args, g Group) {
					t.Helper()
					g.Go(func() error {
						return errors.New("err-1")
					})
				},
				checkFunc: func(got Group) error {
					if err := got.Wait(); err == nil {
						return errors.New("err is nil")
					}

					keys := []string{
						"err", "err-1",
					}

					for _, k := range keys {
						if _, ok := got.(*group).emap[k]; !ok {
							return errors.Errorf("emap key: %s not exist", k)
						}
					}

					if got, want := int(atomic.LoadInt32(&calledCnt)), 1; got != want {
						return errors.Errorf("calledCnt = %v, want: %v", got, want)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			g := &group{
				egctx:  test.fields.egctx,
				cancel: test.fields.cancel,
				sem:    test.fields.sem,
				emap:   test.fields.emap,
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, g)
			}

			g.Go(test.args.f)
			if err := checkFunc(g); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_doCancel(t *testing.T) {
	type fields struct {
		cancel context.CancelCauseFunc
	}
	type test struct {
		name       string
		fields     fields
		checkFunc  func() error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		func() test {
			var called bool

			return test{
				name: "g.cancel is called when g.cancel is not nil",
				fields: fields{
					cancel: func(err error) {
						called = true
					},
				},
				checkFunc: func() error {
					if !called {
						return errors.Errorf("got called = %v, want: %v", called, true)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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
			g := &group{
				cancel: test.fields.cancel,
			}

			g.doCancel(context.Canceled)
			if err := checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWait(t *testing.T) {
	type want struct {
		err error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns nil when Wait returns nil",
			beforeFunc: func(t *testing.T) {
				t.Helper()
				instance, _ = New(context.Background())
			},
			afterFunc: func(t *testing.T) {
				t.Helper()
				instance = nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			err := Wait()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_Wait(t *testing.T) {
	type fields struct {
		sem  *semaphore.Weighted
		errs []error
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, Group)
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			var num int32
			return test{
				name:   "returns nil after all goroutne returns",
				fields: fields{},
				beforeFunc: func(t *testing.T, g Group) {
					t.Helper()
					g.Go(func() error {
						atomic.StoreInt32(&num, int32(runtime.NumGoroutine()))
						time.Sleep(time.Second)
						return nil
					})
				},
				checkFunc: func(w want, err error) error {
					if err := defaultCheckFunc(w, err); err != nil {
						return err
					}

					if got, want := int(atomic.LoadInt32(&num)), runtime.NumGoroutine(); got <= want {
						return errors.New("all goroutine not returns")
					}
					return nil
				},
			}
		}(),

		{
			name: "returns error when g.errs is not nil",
			fields: fields{
				errs: []error{
					errors.New("err1"),
					errors.New("err2"),
				},
			},
			want: want{
				err: errors.Join(errors.New("err1"), errors.New("err2")),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			g := &group{
				sem:  test.fields.sem,
				errs: test.fields.errs,
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, g)
			}

			err := g.Wait()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestWithContext(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type want struct {
// 		want  Group
// 		want1 context.Context
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Group, context.Context) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Group, got1 context.Context) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		if !reflect.DeepEqual(got1, w.want1) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, got1 := WithContext(test.args.ctx)
// 			if err := checkFunc(test.want, got, got1); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestTryGo(t *testing.T) {
// 	type args struct {
// 		f func() error
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           f:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           f:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := TryGo(test.args.f)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_group_exec(t *testing.T) {
// 	type args struct {
// 		f func() error
// 	}
// 	type fields struct {
// 		egctx  context.Context
// 		cancel context.CancelCauseFunc
// 		limit  atomic.Int64
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			g := &group{
// 				egctx:  test.fields.egctx,
// 				cancel: test.fields.cancel,
// 				limit:  test.fields.limit,
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 			}
//
// 			g.exec(test.args.f)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_group_run(t *testing.T) {
// 	type args struct {
// 		f func() error
// 	}
// 	type fields struct {
// 		egctx  context.Context
// 		cancel context.CancelCauseFunc
// 		limit  atomic.Int64
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			g := &group{
// 				egctx:  test.fields.egctx,
// 				cancel: test.fields.cancel,
// 				limit:  test.fields.limit,
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 			}
//
// 			g.run(test.args.f)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_group_TryGo(t *testing.T) {
// 	type args struct {
// 		f func() error
// 	}
// 	type fields struct {
// 		egctx  context.Context
// 		cancel context.CancelCauseFunc
// 		limit  atomic.Int64
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			g := &group{
// 				egctx:  test.fields.egctx,
// 				cancel: test.fields.cancel,
// 				limit:  test.fields.limit,
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 			}
//
// 			got := g.TryGo(test.args.f)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_group_appendErr(t *testing.T) {
// 	type args struct {
// 		err error
// 	}
// 	type fields struct {
// 		egctx  context.Context
// 		cancel context.CancelCauseFunc
// 		limit  atomic.Int64
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           err:nil,
// 		       },
// 		       fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           err:nil,
// 		           },
// 		           fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			g := &group{
// 				egctx:  test.fields.egctx,
// 				cancel: test.fields.cancel,
// 				limit:  test.fields.limit,
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 			}
//
// 			g.appendErr(test.args.err)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_group_done(t *testing.T) {
// 	type fields struct {
// 		egctx  context.Context
// 		cancel context.CancelCauseFunc
// 		limit  atomic.Int64
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           egctx:nil,
// 		           cancel:nil,
// 		           limit:nil,
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			g := &group{
// 				egctx:  test.fields.egctx,
// 				cancel: test.fields.cancel,
// 				limit:  test.fields.limit,
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 			}
//
// 			g.done()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
