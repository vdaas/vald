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

package watch

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"go.uber.org/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	goleak.IgnoreTopFunction("syscall.Syscall6"),
	goleak.IgnoreTopFunction("syscall.syscall6"),
}

func TestMain(m *testing.M) {
	log.Init()

	code := m.Run()

	os.Exit(code)
}

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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns (Watcher, nil) when option is not nil",
			args: args{
				opts: []Option{
					WithDirs("watch.go"),
				},
			},
			want: want{
				err: nil,
			},
			checkFunc: func(w want, got Watcher, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}
				if got == nil {
					return errors.Errorf("got is nil")
				}
				return nil
			},
		},

		{
			name: "returns (nil, error) when option is nil and w.dirs is empty",
			want: want{
				err: errors.ErrWatchDirNotFound,
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

func Test_watch_init(t *testing.T) {
	type fields struct {
		w    *fsnotify.Watcher
		dirs map[string]struct{}
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
		beforeFunc func(*testing.T, *fields)
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *watch, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns no such file or directory error when file not exists",
			fields: fields{
				dirs: map[string]struct{}{
					"vald.go": {},
				},
			},
			want: want{
				err: syscall.Errno(0x2),
			},
		},

		{
			name: "returns no such file or directory error when directory not exists",
			fields: fields{
				dirs: map[string]struct{}{
					"test": {},
				},
			},
			want: want{
				err: syscall.Errno(0x2),
			},
		},

		{
			name: "returns no such file or directory error when some file not exists",
			fields: fields{
				dirs: map[string]struct{}{
					"watch.go": {},
					"vald.go":  {},
				},
			},
			want: want{
				err: syscall.Errno(0x2),
			},
		},

		{
			name: "returns nil when watcher already created and initialize success",
			fields: fields{
				dirs: map[string]struct{}{
					"../watch":      {},
					"watch.go":      {},
					"watch_test.go": {},
				},
				w: func() *fsnotify.Watcher {
					w, _ := fsnotify.NewWatcher()
					return w
				}(),
			},
			checkFunc: func(w want, got *watch, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}
				if got == nil {
					return errors.New("got is nil")
				}
				if got.w == nil {
					return errors.New("got w is nil")
				}
				return nil
			},
			want: want{
				err: nil,
			},
		},

		{
			name: "returns nil when initialize success",
			fields: fields{
				dirs: map[string]struct{}{
					"../watch":      {},
					"watch.go":      {},
					"watch_test.go": {},
				},
			},
			checkFunc: func(w want, got *watch, err error) error {
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}
				if got == nil {
					return errors.New("got is nil")
				}
				if got.w == nil {
					return errors.New("got w is nil")
				}
				return nil
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(tt, &test.fields)
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &watch{
				w:    test.fields.w,
				dirs: test.fields.dirs,
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
		dirs     map[string]struct{}
		eg       errgroup.Group
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
		fieldsFunc func(*testing.T) fields
		want       want
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		if !errors.Is(<-got, <-w.want) {
			return errors.Errorf("errCh got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}

		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "returns channel with error when watcher closed and initialize fails",
				args: args{
					ctx: ctx,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					w, err := fsnotify.NewWatcher()
					if err != nil {
						t.Fatal(err)
					}
					w.Errors = make(chan error, 1)

					w.Errors <- errors.New("err")
					w.Close()

					return fields{
						w:  w,
						eg: errgroup.Get(),
						dirs: map[string]struct{}{
							"vald": {},
						},
					}
				},
				afterFunc: func(args) {
					cancel()
				},
				want: want{
					want: func() chan error {
						ch := make(chan error, 1)
						ch <- errors.New("err")
						close(ch)
						return ch
					}(),
					err: nil,
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "return channel with error when the write event occurs and onChange and onWrite hook returns error",
				args: args{
					ctx: ctx,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()

					w, err := fsnotify.NewWatcher()
					if err != nil {
						t.Fatal(err)
					}
					w.Events = make(chan fsnotify.Event, 1)

					w.Events <- fsnotify.Event{
						Name: "watch.go",
						Op:   fsnotify.Write,
					}

					return fields{
						w:  w,
						eg: errgroup.Get(),
						onChange: func(ctx context.Context, name string) error {
							if got, want := name, "watch.go"; got != want {
								t.Errorf("onChange name got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, want)
							}
							return errors.New("err1")
						},
						onWrite: func(ctx context.Context, name string) error {
							if got, want := name, "watch.go"; got != want {
								t.Errorf("onWrite name got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, want)
							}
							return errors.New("err2")
						},
					}
				},
				afterFunc: func(args) {
					cancel()
				},
				want: want{
					want: func() chan error {
						ch := make(chan error, 2)
						ch <- errors.New("err1")
						ch <- errors.New("err2")
						close(ch)
						return ch
					}(),
					err: nil,
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "return channel with error when onWrite hook return error and send to returned channel",
				args: args{
					ctx: ctx,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					w, err := fsnotify.NewWatcher()
					if err != nil {
						t.Fatal(err)
					}
					w.Events = make(chan fsnotify.Event, 1)

					w.Events <- fsnotify.Event{
						Name: "watch.go",
						Op:   fsnotify.Write,
					}

					return fields{
						w:  w,
						eg: errgroup.Get(),
						onWrite: func(ctx context.Context, name string) error {
							if got, want := name, "watch.go"; got != want {
								t.Errorf("onWrite name got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, want)
							}
							return errors.New("err")
						},
					}
				},
				afterFunc: func(args) {
					cancel()
				},
				want: want{
					want: func() chan error {
						ch := make(chan error, 1)
						ch <- errors.New("err")
						close(ch)
						return ch
					}(),
					err: nil,
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "return channel with error when onCreate hook return error and send to returned channel",
				args: args{
					ctx: ctx,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					w, err := fsnotify.NewWatcher()
					if err != nil {
						t.Fatal(err)
					}
					w.Events = make(chan fsnotify.Event, 1)

					w.Events <- fsnotify.Event{
						Name: "watch.go",
						Op:   fsnotify.Create,
					}

					return fields{
						w:  w,
						eg: errgroup.Get(),
						onCreate: func(ctx context.Context, name string) error {
							if got, want := name, "watch.go"; got != want {
								t.Errorf("onWrite name got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, want)
							}
							return errors.New("err")
						},
					}
				},
				afterFunc: func(args) {
					cancel()
				},
				want: want{
					want: func() chan error {
						ch := make(chan error, 1)
						ch <- errors.New("err")
						close(ch)
						return ch
					}(),
					err: nil,
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "return channel with error when onDelete hook return error and send to returned channel",
				args: args{
					ctx: ctx,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					w, err := fsnotify.NewWatcher()
					if err != nil {
						t.Fatal(err)
					}
					w.Events = make(chan fsnotify.Event, 1)

					w.Events <- fsnotify.Event{
						Name: "watch.go",
						Op:   fsnotify.Remove,
					}

					return fields{
						w:  w,
						eg: errgroup.Get(),
						onDelete: func(ctx context.Context, name string) error {
							if got, want := name, "watch.go"; got != want {
								t.Errorf("onWrite name got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, want)
							}
							return errors.New("err")
						},
					}
				},
				afterFunc: func(args) {
					cancel()
				},
				want: want{
					want: func() chan error {
						ch := make(chan error, 1)
						ch <- errors.New("err")
						close(ch)
						return ch
					}(),
					err: nil,
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "return channel with error when onChmod hook return error and send to returned channel",
				args: args{
					ctx: ctx,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					w, err := fsnotify.NewWatcher()
					if err != nil {
						t.Fatal(err)
					}
					w.Events = make(chan fsnotify.Event, 1)

					w.Events <- fsnotify.Event{
						Name: "watch.go",
						Op:   fsnotify.Chmod,
					}

					return fields{
						w:  w,
						eg: errgroup.Get(),
						onChmod: func(ctx context.Context, name string) error {
							if got, want := name, "watch.go"; got != want {
								t.Errorf("onWrite name got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, want)
							}
							return errors.New("err")
						},
					}
				},
				afterFunc: func(args) {
					cancel()
				},
				want: want{
					want: func() chan error {
						ch := make(chan error, 1)
						ch <- errors.New("err")
						close(ch)
						return ch
					}(),
					err: nil,
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "return channel with error when onRename hook return error and send to returned channel",
				args: args{
					ctx: ctx,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					w, err := fsnotify.NewWatcher()
					if err != nil {
						t.Fatal(err)
					}
					w.Events = make(chan fsnotify.Event, 1)

					w.Events <- fsnotify.Event{
						Name: "watch.go",
						Op:   fsnotify.Rename,
					}

					return fields{
						w:  w,
						eg: errgroup.Get(),
						onRename: func(ctx context.Context, name string) error {
							if got, want := name, "watch.go"; got != want {
								t.Errorf("onWrite name got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, want)
							}
							return errors.New("err")
						},
					}
				},
				afterFunc: func(args) {
					cancel()
				},
				want: want{
					want: func() chan error {
						ch := make(chan error, 1)
						ch <- errors.New("err")
						close(ch)
						return ch
					}(),
					err: nil,
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "return channel with error when onWrite and onRename hook return error and send to returned channel",
				args: args{
					ctx: ctx,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					w, err := fsnotify.NewWatcher()
					if err != nil {
						t.Fatal(err)
					}
					w.Events = make(chan fsnotify.Event, 2)

					w.Events <- fsnotify.Event{
						Name: "watch.go",
						Op:   fsnotify.Write,
					}

					w.Events <- fsnotify.Event{
						Name: "watch.go",
						Op:   fsnotify.Rename,
					}

					return fields{
						w:  w,
						eg: errgroup.Get(),
						onWrite: func(ctx context.Context, name string) error {
							if got, want := name, "watch.go"; got != want {
								t.Errorf("onWrite name got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, want)
							}
							return errors.New("err1")
						},
						onRename: func(ctx context.Context, name string) error {
							if got, want := name, "watch.go"; got != want {
								t.Errorf("onRename name got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, want)
							}
							return errors.New("err2")
						},
					}
				},
				afterFunc: func(args) {
					cancel()
				},
				want: want{
					want: func() chan error {
						ch := make(chan error, 2)
						ch <- errors.New("err1")
						ch <- errors.New("err2")
						close(ch)
						return ch
					}(),
					err: nil,
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
			fields := test.fieldsFunc(tt)

			w := &watch{
				w:        fields.w,
				eg:       fields.eg,
				dirs:     fields.dirs,
				onChange: fields.onChange,
				onCreate: fields.onCreate,
				onRename: fields.onRename,
				onDelete: fields.onDelete,
				onWrite:  fields.onWrite,
				onChmod:  fields.onChmod,
				onError:  fields.onError,
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
		w    *fsnotify.Watcher
		dirs map[string]struct{}
	}
	type want struct {
		err  error
		want *watch
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *watch, error) error
		beforeFunc func(*testing.T, *fields, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *watch, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		if got, want := len(got.dirs), len(w.want.dirs); got != want {
			return errors.Errorf("dirs length = %d, want %d", got, want)
		}
		for name := range w.want.dirs {
			if _, ok := got.dirs[name]; !ok {
				return errors.Errorf("dirs %s is not exists", name)
			}
		}
		return nil
	}
	defaultBeforeFunc := func(t *testing.T, fields *fields, args args) {
		t.Helper()

		var err error
		fields.w, err = fsnotify.NewWatcher()
		if err != nil {
			t.Fatal(err)
		}
	}
	tests := []test{
		{
			name: "returns nil when add success",
			args: args{
				dirs: []string{
					"./watch.go",
					"./option.go",
				},
			},
			fields: fields{
				dirs: map[string]struct{}{
					"watch_test.go": {},
				},
			},
			want: want{
				err: nil,
				want: &watch{
					dirs: map[string]struct{}{
						"watch_test.go": {},
						"./watch.go":    {},
						"./option.go":   {},
					},
				},
			},
		},

		{
			name: "returns nil when directory add success",
			args: args{
				dirs: []string{
					"../watch",
				},
			},
			fields: fields{
				dirs: make(map[string]struct{}),
			},
			want: want{
				err: nil,
				want: &watch{
					dirs: map[string]struct{}{
						"../watch": {},
					},
				},
			},
		},

		{
			name: "returns no such file or directory error when some file not exists",
			args: args{
				dirs: []string{
					"watch.go",
					"vald.go",
				},
			},
			fields: fields{
				dirs: make(map[string]struct{}),
			},
			want: want{
				err: syscall.Errno(0x2),
				want: &watch{
					dirs: map[string]struct{}{
						"watch.go": {},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			test.beforeFunc(tt, &test.fields, test.args)

			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &watch{
				w:    test.fields.w,
				dirs: test.fields.dirs,
			}

			err := w.Add(test.args.dirs...)
			if err := test.checkFunc(test.want, w, err); err != nil {
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
		w    *fsnotify.Watcher
		dirs map[string]struct{}
	}
	type want struct {
		want *watch
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *watch, error) error
		beforeFunc func(*testing.T, *fields, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got *watch, err error) error {
		if w.err == nil {
			if err != nil {
				return errors.Errorf("got error is not nil: %v", err)
			}
		} else {
			if err == nil {
				return errors.New("got error is nil")
			}

			if !strings.Contains(err.Error(), w.err.Error()) {
				return errors.Errorf("got error  %v, not contains: %v", err, w.err.Error())
			}
		}

		if got, want := len(got.dirs), len(w.want.dirs); got != want {
			return errors.Errorf("dirs length = %d, want %d", got, want)
		}

		for name := range w.want.dirs {
			if _, ok := got.dirs[name]; !ok {
				return errors.Errorf("dirs %s is not exists", name)
			}
		}
		return nil
	}
	defaultBeforeFunc := func(t *testing.T, fields *fields, args args) {
		t.Helper()

		var err error
		fields.w, err = fsnotify.NewWatcher()
		if err != nil {
			t.Fatal(err)
		}

		for name := range fields.dirs {
			if err = fields.w.Add(name); err != nil {
				t.Fatal(err)
			}
		}
	}
	tests := []test{
		{
			name: "returns nil when remove success",
			args: args{
				dirs: []string{
					"watch.go",
					"watch_test.go",
				},
			},
			fields: fields{
				dirs: map[string]struct{}{
					"watch.go":      {},
					"watch_test.go": {},
					"option.go":     {},
				},
			},
			want: want{
				want: &watch{
					dirs: map[string]struct{}{
						"option.go": {},
					},
				},
				err: nil,
			},
		},

		{
			name: "returns nil when directory remove success",
			args: args{
				dirs: []string{
					"../watch",
				},
			},
			fields: fields{
				dirs: map[string]struct{}{
					"../watch": {},
				},
			},
			want: want{
				want: &watch{
					dirs: map[string]struct{}{},
				},
				err: nil,
			},
		},

		{
			name: "returns non-exist error when some file not exists",
			args: args{
				dirs: []string{
					"watch.go",
					"vald.go",
					"watch_test.go",
				},
			},
			fields: fields{
				dirs: map[string]struct{}{
					"watch.go":      {},
					"watch_test.go": {},
				},
			},
			want: want{
				want: &watch{
					dirs: map[string]struct{}{
						"watch_test.go": {},
					},
				},
				err: fmt.Errorf("can't remove non-existent"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			test.beforeFunc(tt, &test.fields, test.args)

			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			w := &watch{
				w:    test.fields.w,
				dirs: test.fields.dirs,
			}

			err := w.Remove(test.args.dirs...)
			if err := test.checkFunc(test.want, w, err); err != nil {
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
		w    *fsnotify.Watcher
		dirs map[string]struct{}
	}
	type want struct {
		err  error
		want *watch
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *watch, error) error
		beforeFunc func(*testing.T, *fields, args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *watch, err error) error {
		if w.err == nil {
			if err != nil {
				return errors.Errorf("got error is not nil: %v", err)
			}
		} else {
			if err == nil {
				return errors.New("got error is nil")
			}

			if !strings.Contains(err.Error(), w.err.Error()) {
				return errors.Errorf("got error  %v, not contains: %v", err, w.err.Error())
			}
		}

		if got, want := len(got.dirs), len(w.want.dirs); got != want {
			return errors.Errorf("dirs length = %d, want %d", got, want)
		}

		for name := range w.want.dirs {
			if _, ok := got.dirs[name]; !ok {
				return errors.Errorf("dirs %s is not exists", name)
			}
		}
		return nil
	}
	defaultBeforeFunc := func(t *testing.T, fields *fields, args args) {
		t.Helper()

		var err error
		fields.w, err = fsnotify.NewWatcher()
		if err != nil {
			t.Fatal(err)
		}
	}
	tests := []test{
		{
			name: "returns nil when stop success",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				dirs: map[string]struct{}{
					"../watch":      {},
					"watch.go":      {},
					"watch_test.go": {},
				},
			},
			beforeFunc: func(t *testing.T, fields *fields, args args) {
				t.Helper()
				defaultBeforeFunc(t, fields, args)

				for name := range fields.dirs {
					err := fields.w.Add(name)
					if err != nil {
						t.Fatal(err)
					}
				}
			},
			want: want{
				want: &watch{
					dirs: make(map[string]struct{}),
				},
				err: nil,
			},
		},

		{
			name: "returns non-exist error when the file not exists",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				dirs: map[string]struct{}{
					"watch.go": {},
				},
			},
			want: want{
				want: &watch{
					dirs: make(map[string]struct{}),
				},
				err: fmt.Errorf("can't remove non-existent"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			test.beforeFunc(tt, &test.fields, test.args)
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			w := &watch{
				w:    test.fields.w,
				dirs: test.fields.dirs,
			}

			err := w.Stop(test.args.ctx)
			if err := test.checkFunc(test.want, w, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
