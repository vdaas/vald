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
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got = %v, want %v", got1, w.want1)
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
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
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
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
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
	type want struct {
	}
	type test struct {
		name       string
		args       args
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
		           f: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
			if err := test.checkFunc(test.want); err != nil {
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
		egctx            context.Context
		cancel           func()
		wg               sync.WaitGroup
		limitation       chan struct{}
		enableLimitation atomic.Value
		cancelOnce       sync.Once
		mu               sync.RWMutex
		emap             map[string]struct{}
		errs             []error
		err              error
	}
	type want struct {
	}
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
		           limit: 0,
		       },
		       fields: fields {
		           egctx: nil,
		           cancel: nil,
		           wg: sync.WaitGroup{},
		           limitation: nil,
		           enableLimitation: nil,
		           cancelOnce: sync.Once{},
		           mu: sync.RWMutex{},
		           emap: nil,
		           errs: nil,
		           err: nil,
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
		           limit: 0,
		           },
		           fields: fields {
		           egctx: nil,
		           cancel: nil,
		           wg: sync.WaitGroup{},
		           limitation: nil,
		           enableLimitation: nil,
		           cancelOnce: sync.Once{},
		           mu: sync.RWMutex{},
		           emap: nil,
		           errs: nil,
		           err: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				egctx:            test.fields.egctx,
				cancel:           test.fields.cancel,
				wg:               test.fields.wg,
				limitation:       test.fields.limitation,
				enableLimitation: test.fields.enableLimitation,
				cancelOnce:       test.fields.cancelOnce,
				mu:               test.fields.mu,
				emap:             test.fields.emap,
				errs:             test.fields.errs,
				err:              test.fields.err,
			}

			g.Limitation(test.args.limit)
			if err := test.checkFunc(test.want); err != nil {
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
		wg               sync.WaitGroup
		limitation       chan struct{}
		enableLimitation atomic.Value
		cancelOnce       sync.Once
		mu               sync.RWMutex
		emap             map[string]struct{}
		errs             []error
		err              error
	}
	type want struct {
	}
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
		           f: nil,
		       },
		       fields: fields {
		           egctx: nil,
		           cancel: nil,
		           wg: sync.WaitGroup{},
		           limitation: nil,
		           enableLimitation: nil,
		           cancelOnce: sync.Once{},
		           mu: sync.RWMutex{},
		           emap: nil,
		           errs: nil,
		           err: nil,
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
		           f: nil,
		           },
		           fields: fields {
		           egctx: nil,
		           cancel: nil,
		           wg: sync.WaitGroup{},
		           limitation: nil,
		           enableLimitation: nil,
		           cancelOnce: sync.Once{},
		           mu: sync.RWMutex{},
		           emap: nil,
		           errs: nil,
		           err: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				egctx:            test.fields.egctx,
				cancel:           test.fields.cancel,
				wg:               test.fields.wg,
				limitation:       test.fields.limitation,
				enableLimitation: test.fields.enableLimitation,
				cancelOnce:       test.fields.cancelOnce,
				mu:               test.fields.mu,
				emap:             test.fields.emap,
				errs:             test.fields.errs,
				err:              test.fields.err,
			}

			g.Go(test.args.f)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_doCancel(t *testing.T) {
	type fields struct {
		egctx            context.Context
		cancel           func()
		wg               sync.WaitGroup
		limitation       chan struct{}
		enableLimitation atomic.Value
		cancelOnce       sync.Once
		mu               sync.RWMutex
		emap             map[string]struct{}
		errs             []error
		err              error
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
		           egctx: nil,
		           cancel: nil,
		           wg: sync.WaitGroup{},
		           limitation: nil,
		           enableLimitation: nil,
		           cancelOnce: sync.Once{},
		           mu: sync.RWMutex{},
		           emap: nil,
		           errs: nil,
		           err: nil,
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
		           egctx: nil,
		           cancel: nil,
		           wg: sync.WaitGroup{},
		           limitation: nil,
		           enableLimitation: nil,
		           cancelOnce: sync.Once{},
		           mu: sync.RWMutex{},
		           emap: nil,
		           errs: nil,
		           err: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &group{
				egctx:            test.fields.egctx,
				cancel:           test.fields.cancel,
				wg:               test.fields.wg,
				limitation:       test.fields.limitation,
				enableLimitation: test.fields.enableLimitation,
				cancelOnce:       test.fields.cancelOnce,
				mu:               test.fields.mu,
				emap:             test.fields.emap,
				errs:             test.fields.errs,
				err:              test.fields.err,
			}

			g.doCancel()
			if err := test.checkFunc(test.want); err != nil {
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			err := Wait()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_group_Wait(t *testing.T) {
	type fields struct {
		egctx            context.Context
		cancel           func()
		wg               sync.WaitGroup
		limitation       chan struct{}
		enableLimitation atomic.Value
		cancelOnce       sync.Once
		mu               sync.RWMutex
		emap             map[string]struct{}
		errs             []error
		err              error
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           egctx: nil,
		           cancel: nil,
		           wg: sync.WaitGroup{},
		           limitation: nil,
		           enableLimitation: nil,
		           cancelOnce: sync.Once{},
		           mu: sync.RWMutex{},
		           emap: nil,
		           errs: nil,
		           err: nil,
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
		           egctx: nil,
		           cancel: nil,
		           wg: sync.WaitGroup{},
		           limitation: nil,
		           enableLimitation: nil,
		           cancelOnce: sync.Once{},
		           mu: sync.RWMutex{},
		           emap: nil,
		           errs: nil,
		           err: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			g := &group{
				egctx:            test.fields.egctx,
				cancel:           test.fields.cancel,
				wg:               test.fields.wg,
				limitation:       test.fields.limitation,
				enableLimitation: test.fields.enableLimitation,
				cancelOnce:       test.fields.cancelOnce,
				mu:               test.fields.mu,
				emap:             test.fields.emap,
				errs:             test.fields.errs,
				err:              test.fields.err,
			}

			err := g.Wait()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
