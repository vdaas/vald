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

package watch

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Watcher
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Watcher, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Watcher, err error) error {
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

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_watch_init(t *testing.T) {
	type fields struct {
		w        *fsnotify.Watcher
		eg       errgroup.Group
		dirs     map[string]struct{}
		mu       sync.RWMutex
		onChange func(ctx context.Context, name string) error
		onCreate func(ctx context.Context, name string) error
		onRename func(ctx context.Context, name string) error
		onDelete func(ctx context.Context, name string) error
		onWrite  func(ctx context.Context, name string) error
		onChmod  func(ctx context.Context, name string) error
		onError  func(ctx context.Context, err error) error
	}
	type want struct {
		want *watch
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *watch, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *watch, err error) error {
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
		       fields: fields {
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
			w := &watch{
				w:        test.fields.w,
				eg:       test.fields.eg,
				dirs:     test.fields.dirs,
				mu:       test.fields.mu,
				onChange: test.fields.onChange,
				onCreate: test.fields.onCreate,
				onRename: test.fields.onRename,
				onDelete: test.fields.onDelete,
				onWrite:  test.fields.onWrite,
				onChmod:  test.fields.onChmod,
				onError:  test.fields.onError,
			}

			got, err := w.init()
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_watch_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		w        *fsnotify.Watcher
		eg       errgroup.Group
		dirs     map[string]struct{}
		mu       sync.RWMutex
		onChange func(ctx context.Context, name string) error
		onCreate func(ctx context.Context, name string) error
		onRename func(ctx context.Context, name string) error
		onDelete func(ctx context.Context, name string) error
		onWrite  func(ctx context.Context, name string) error
		onChmod  func(ctx context.Context, name string) error
		onError  func(ctx context.Context, err error) error
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
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
			w := &watch{
				w:        test.fields.w,
				eg:       test.fields.eg,
				dirs:     test.fields.dirs,
				mu:       test.fields.mu,
				onChange: test.fields.onChange,
				onCreate: test.fields.onCreate,
				onRename: test.fields.onRename,
				onDelete: test.fields.onDelete,
				onWrite:  test.fields.onWrite,
				onChmod:  test.fields.onChmod,
				onError:  test.fields.onError,
			}

			got, err := w.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_watch_Add(t *testing.T) {
	type args struct {
		dirs []string
	}
	type fields struct {
		w        *fsnotify.Watcher
		eg       errgroup.Group
		dirs     map[string]struct{}
		mu       sync.RWMutex
		onChange func(ctx context.Context, name string) error
		onCreate func(ctx context.Context, name string) error
		onRename func(ctx context.Context, name string) error
		onDelete func(ctx context.Context, name string) error
		onWrite  func(ctx context.Context, name string) error
		onChmod  func(ctx context.Context, name string) error
		onError  func(ctx context.Context, err error) error
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
		           dirs: nil,
		       },
		       fields: fields {
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
		           dirs: nil,
		           },
		           fields: fields {
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
			w := &watch{
				w:        test.fields.w,
				eg:       test.fields.eg,
				dirs:     test.fields.dirs,
				mu:       test.fields.mu,
				onChange: test.fields.onChange,
				onCreate: test.fields.onCreate,
				onRename: test.fields.onRename,
				onDelete: test.fields.onDelete,
				onWrite:  test.fields.onWrite,
				onChmod:  test.fields.onChmod,
				onError:  test.fields.onError,
			}

			err := w.Add(test.args.dirs...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_watch_Remove(t *testing.T) {
	type args struct {
		dirs []string
	}
	type fields struct {
		w        *fsnotify.Watcher
		eg       errgroup.Group
		dirs     map[string]struct{}
		mu       sync.RWMutex
		onChange func(ctx context.Context, name string) error
		onCreate func(ctx context.Context, name string) error
		onRename func(ctx context.Context, name string) error
		onDelete func(ctx context.Context, name string) error
		onWrite  func(ctx context.Context, name string) error
		onChmod  func(ctx context.Context, name string) error
		onError  func(ctx context.Context, err error) error
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
		           dirs: nil,
		       },
		       fields: fields {
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
		           dirs: nil,
		           },
		           fields: fields {
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
			w := &watch{
				w:        test.fields.w,
				eg:       test.fields.eg,
				dirs:     test.fields.dirs,
				mu:       test.fields.mu,
				onChange: test.fields.onChange,
				onCreate: test.fields.onCreate,
				onRename: test.fields.onRename,
				onDelete: test.fields.onDelete,
				onWrite:  test.fields.onWrite,
				onChmod:  test.fields.onChmod,
				onError:  test.fields.onError,
			}

			err := w.Remove(test.args.dirs...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_watch_Stop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		w        *fsnotify.Watcher
		eg       errgroup.Group
		dirs     map[string]struct{}
		mu       sync.RWMutex
		onChange func(ctx context.Context, name string) error
		onCreate func(ctx context.Context, name string) error
		onRename func(ctx context.Context, name string) error
		onDelete func(ctx context.Context, name string) error
		onWrite  func(ctx context.Context, name string) error
		onChmod  func(ctx context.Context, name string) error
		onError  func(ctx context.Context, err error) error
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
		       },
		       fields: fields {
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
		           w: nil,
		           eg: nil,
		           dirs: nil,
		           mu: sync.RWMutex{},
		           onChange: nil,
		           onCreate: nil,
		           onRename: nil,
		           onDelete: nil,
		           onWrite: nil,
		           onChmod: nil,
		           onError: nil,
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
			w := &watch{
				w:        test.fields.w,
				eg:       test.fields.eg,
				dirs:     test.fields.dirs,
				mu:       test.fields.mu,
				onChange: test.fields.onChange,
				onCreate: test.fields.onCreate,
				onRename: test.fields.onRename,
				onDelete: test.fields.onDelete,
				onWrite:  test.fields.onWrite,
				onChmod:  test.fields.onChmod,
				onError:  test.fields.onError,
			}

			err := w.Stop(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
