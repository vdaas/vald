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

// Package errgroup provides server global wait group for graceful kill all goroutine
package errgroup

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"

	"go.uber.org/goleak"
)

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
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Group, got1 context.Context) error {
		if got, want := got.(*group), w.want.(*group); !reflect.DeepEqual(got.emap, want.emap) &&
			!reflect.DeepEqual(got.enableLimitation, want.enableLimitation) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got = %v, want %v", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx := context.Background()
			egctx, cancel := context.WithCancel(ctx)

			var enableLimitation atomic.Value
			enableLimitation.Store(false)

			return test{
				name: "returns (g, ctx)",
				args: args{
					ctx: ctx,
				},
				want: want{
					want: &group{
						egctx:            egctx,
						cancel:           cancel,
						enableLimitation: enableLimitation,
						emap:             make(map[string]struct{}),
					},
					want1: egctx,
				},
			}
		}(),
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

			got, got1 := New(test.args.ctx)
			if err := test.checkFunc(test.want, got, got1); err != nil {
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
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotEgctx context.Context) error {
		if !reflect.DeepEqual(gotEgctx, w.wantEgctx) {
			return errors.Errorf("got = %v, want %v", gotEgctx, w.wantEgctx)
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
			egctx, cancel := context.WithCancel(ctx)

			return test{
				name: "returns egctx when once.Do is called",
				args: args{
					ctx: ctx,
				},
				want: want{
					wantEgctx: egctx,
				},
				afterFunc: func(a args) {
					cancel()
					defaultBeforeFunc(a)
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			test.beforeFunc(test.args)

			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotEgctx := Init(test.args.ctx)
			if err := test.checkFunc(test.want, gotEgctx); err != nil {
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got Group) error {
		if got, want := got.(*group), w.want.(*group); !reflect.DeepEqual(got.egctx, want.egctx) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	defaultBeforeFunc := func() {
		instance, once = nil, sync.Once{}
	}
	defaultAfterFunc := func() {
		defaultBeforeFunc()
	}

	tests := []test{
		func() test {
			ctx := context.Background()
			egctx, cancel := context.WithCancel(ctx)

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
				beforeFunc: func() {
					instance = g
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			test.beforeFunc()

			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}
			defer test.afterFunc()

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := Get()
			if err := test.checkFunc(test.want, got); err != nil {
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
		beforeFunc func(args)
		afterFunc  func(args)
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
				beforeFunc: func(args) {
					g := new(group)
					g.enableLimitation.Store(false)
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

			Go(test.args.f)
			if err := test.checkFunc(instance); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_Limitation(t *testing.T) {
	type args struct {
		limit int
	}
	type fields struct {
		limitation       chan struct{}
		enableLimitation atomic.Value
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
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, g Group) error {
		got, want := g.(*group), w.want.(*group)
		if !reflect.DeepEqual(got.enableLimitation, want.enableLimitation) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		if got.limitation != nil && want.limitation != nil {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "set disable when limit is 0",
			args: args{
				limit: 0,
			},
			fields: fields{
				limitation: make(chan struct{}),
			},
			want: want{
				want: &group{
					enableLimitation: func() atomic.Value {
						var el atomic.Value
						el.Store(false)
						return el
					}(),
				},
			},
		},

		{
			name: "set enable when limit is 1",
			args: args{
				limit: 1,
			},
			fields: fields{
				limitation: make(chan struct{}),
			},
			want: want{
				want: &group{
					enableLimitation: func() atomic.Value {
						var el atomic.Value
						el.Store(true)
						return el
					}(),
				},
			},
		},
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
			g := &group{
				limitation:       test.fields.limitation,
				enableLimitation: test.fields.enableLimitation,
			}

			g.Limitation(test.args.limit)
			if err := test.checkFunc(test.want, g); err != nil {
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
		egctx            context.Context
		cancel           func()
		limitation       chan struct{}
		enableLimitation atomic.Value
		emap             map[string]struct{}
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(*fields) Group
		checkFunc  func(Group) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(g Group) error {
		return nil
	}
	defaultCreateFunc := func(fields *fields) Group {
		return &group{
			egctx:            fields.egctx,
			cancel:           fields.cancel,
			limitation:       fields.limitation,
			enableLimitation: fields.enableLimitation,
			emap:             fields.emap,
		}
	}
	tests := []test{
		func() test {
			var calledCnt int32

			egctx, cancel := context.WithCancel(context.Background())
			cancel()

			return test{
				name: "f is not called when reached limit and cancel g.egctx",
				args: args{
					f: func() error {
						atomic.AddInt32(&calledCnt, 1)
						return nil
					},
				},
				fields: fields{
					egctx:      egctx,
					limitation: make(chan struct{}),
				},
				createFunc: func(fields *fields) Group {
					g := defaultCreateFunc(fields)
					g.(*group).enableLimitation.Store(true)
					return g
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

			egctx, cancel := context.WithCancel(context.Background())

			return test{
				name: "f is called and f returns nil",
				args: args{
					f: func() error {
						atomic.AddInt32(&calledCnt, 1)
						return nil
					},
				},
				fields: fields{
					egctx:      egctx,
					cancel:     cancel,
					limitation: make(chan struct{}, 1),
				},
				createFunc: func(fields *fields) Group {
					g := defaultCreateFunc(fields)
					g.(*group).enableLimitation.Store(true)
					return g
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

		func() test {
			var (
				calledCnt       int32
				cancelCalledCnt int32
			)

			egctx, cancel := context.WithCancel(context.Background())

			return test{
				name: "f is called and f returns error",
				args: args{
					f: func() error {
						atomic.AddInt32(&calledCnt, 1)
						return errors.New("err1")
					},
				},
				fields: fields{
					egctx: egctx,
					cancel: func() {
						cancel()
						atomic.AddInt32(&cancelCalledCnt, 1)
					},
					emap:       make(map[string]struct{}),
					limitation: make(chan struct{}, 1),
				},
				createFunc: func(fields *fields) Group {
					g := defaultCreateFunc(fields)
					g.(*group).enableLimitation.Store(true)
					return g
				},
				checkFunc: func(got Group) error {
					if err := got.Wait(); err == nil {
						return errors.New("Wait returns nil")
					}

					if got, want := int(atomic.LoadInt32(&calledCnt)), 1; got != want {
						return errors.Errorf("calledCnt = %v, want: %v", got, want)
					}

					if got, want := int(atomic.LoadInt32(&cancelCalledCnt)), 1; got != want {
						return errors.Errorf("cancel called = %v, want: %v", got, want)
					}
					return nil
				},
			}
		}(),
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
			if test.createFunc == nil {
				test.createFunc = defaultCreateFunc
			}
			g := test.createFunc(&test.fields)

			g.Go(test.args.f)
			if err := test.checkFunc(g); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_doCancel(t *testing.T) {
	type fields struct {
		cancel func()
	}
	type test struct {
		name       string
		fields     fields
		createFunc func(*fields) *group
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
					cancel: func() {
						called = true
					},
				},
				createFunc: func(fields *fields) *group {
					g := &group{
						cancel: fields.cancel,
					}
					return g
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

			g := test.createFunc(&test.fields)

			g.doCancel()
			if err := test.checkFunc(); err != nil {
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
		createFunc func() Group
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns nil when Wait returns nil",
			afterFunc: func() {
				instance = nil
			},
			createFunc: func() Group {
				return new(group)
			},
		},
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

			instance = test.createFunc()

			err := Wait()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_Wait(t *testing.T) {
	type fields struct {
		limitation chan struct{}
		errs       []error
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		createFunc func(*fields) Group
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	defaultCreateFunc := func(fields *fields) Group {
		return &group{
			limitation: fields.limitation,
			errs:       fields.errs,
		}
	}
	tests := []test{
		{
			name: "returns nil when g.errs is nil",
			fields: fields{
				limitation: make(chan struct{}),
			},
		},

		{
			name: "returns error when g.errs is not nil",
			fields: fields{
				limitation: make(chan struct{}),
				errs: []error{
					errors.New("err1"),
					errors.New("err2"),
				},
			},
			want: want{
				err: errors.Wrap(errors.New("err1"), errors.New("err2").Error()),
			},
		},
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
			if test.createFunc == nil {
				test.createFunc = defaultCreateFunc
			}

			g := test.createFunc(&test.fields)

			err := g.Wait()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
