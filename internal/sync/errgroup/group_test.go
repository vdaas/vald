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

package errgroup

import (
	"context"
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/log/nop"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/semaphore"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	// recLog discards like the NOP logger but records Errorf messages so
	// tests can assert that global group task errors are logged.
	log.Init(log.WithLogger(recLog))
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
		want       want
		args       args
		checkFunc  func(want, Group, context.Context) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
		name       string
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
			ctx := t.Context()
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
		args       args
		want       want
		checkFunc  func(want, context.Context) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
		name       string
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
			ctx := t.Context()
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
		func() test {
			ctx, cancel := context.WithCancel(t.Context())

			return test{
				name: "returns ctx unchanged and binds the already initialized instance to ctx",
				args: args{
					ctx: ctx,
				},
				want: want{
					wantEgctx: ctx,
				},
				beforeFunc: func(a args) {
					defaultBeforeFunc(a)
					// the instance is created first with context.Background(), as package init does in production.
					Get()
				},
				checkFunc: func(w want, gotEgctx context.Context) error {
					if err := defaultCheckFunc(w, gotEgctx); err != nil {
						return err
					}
					g, ok := instance.(*group)
					if !ok {
						return errors.Errorf("instance type = %T, want *group", instance)
					}
					if err := g.egctx.Err(); err != nil {
						return errors.Errorf("global group ctx is canceled before ctx is canceled: %v", err)
					}
					cancel()
					select {
					case <-g.egctx.Done():
						return nil
					case <-time.After(5 * time.Second):
						return errors.New("global group ctx is not canceled after ctx is canceled")
					}
				},
				afterFunc: func(t *testing.T, a args) {
					t.Helper()
					cancel()
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
		want       want
		checkFunc  func(want, Group) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		name       string
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
						emap:   make(map[string]struct{}),
					},
				},
			}
		}(),

		func() test {
			g := &group{
				egctx: t.Context(),
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
		args       args
		checkFunc  func(Group) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
	}
	defaultCheckFunc := func(Group) error {
		return nil
	}
	tests := []test{
		func() test {
			var calledCnt atomic.Int32

			return test{
				name: "instance.Go is called when instance is not nil",
				args: args{
					f: func() error {
						calledCnt.Add(1)
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

					if got, want := int(calledCnt.Load()), 1; got != want {
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
		want       want
		fields     fields
		checkFunc  func(want, Group) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
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
		{
			name: "set disable when limit is -1 from 1",
			args: args{
				limit: -1,
			},
			fields: fields{
				sem: semaphore.NewWeighted(1),
			},
			want: want{
				want: &group{
					sem: nil,
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
			g := &group{
				sem: test.fields.sem,
			}

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
		fields     fields
		args       args
		checkFunc  func(Group) error
		beforeFunc func(*testing.T, args, Group)
		afterFunc  func(*testing.T, args)
		name       string
	}
	defaultCheckFunc := func(g Group) error {
		return nil
	}
	tests := []test{
		func() test {
			var calledCnt atomic.Int32

			egctx, cancel := context.WithCancelCause(t.Context())

			limit := 3

			return test{
				name: "f is not called when reached limit and cancel g.egctx",
				args: args{
					f: func() error {
						calledCnt.Add(1)
						return nil
					},
				},
				fields: fields{
					egctx: egctx,
					sem:   semaphore.NewWeighted(int64(limit)),
				},
				beforeFunc: func(t *testing.T, _ args, g Group) {
					t.Helper()
					for range limit {
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

					if got, want := int(calledCnt.Load()), 0; got != want {
						return errors.Errorf("calledCnt = %v, want: %v", got, want)
					}
					return nil
				},
			}
		}(),

		func() test {
			var calledCnt atomic.Int32

			egctx, cancel := context.WithCancelCause(t.Context())

			return test{
				name: "f is called but f returns error and previous process also returns error",
				args: args{
					f: func() error {
						calledCnt.Add(1)
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

					if got, want := int(calledCnt.Load()), 1; got != want {
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
		fields     fields
		checkFunc  func() error
		beforeFunc func()
		afterFunc  func()
		name       string
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
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
		name       string
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
				instance, _ = New(t.Context())
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
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, Group)
		afterFunc  func()
		name       string
		fields     fields
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
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

// Wait blocks until every task has returned and then reports no error. The test used to count
// runtime.NumGoroutine() around a one second Sleep, which depends on the goroutines of every other
// test and of earlier -count runs. In a synctest bubble, synctest.Wait returns once every goroutine
// is durably blocked, so "Wait has not returned while the tasks are held" and "Wait returned after
// they were released" are observed, not timed.
func Test_group_Wait_waitsForAllTasks(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const tasks = 3
		g := &group{}
		release := make(chan struct{})
		var finished atomic.Int32
		for range tasks {
			g.Go(func() error {
				<-release
				finished.Add(1)
				return nil
			})
		}

		var (
			waitErr  error
			returned atomic.Bool
		)
		go func() {
			waitErr = g.Wait()
			returned.Store(true)
		}()

		synctest.Wait() // the tasks wait for release and Wait waits for the tasks
		if returned.Load() {
			t.Fatal("Wait returned before the tasks did")
		}
		close(release)
		synctest.Wait()
		if !returned.Load() {
			t.Fatal("Wait did not return after the tasks returned")
		}
		if got := finished.Load(); got != tasks {
			t.Errorf("finished tasks = %d, want %d", got, tasks)
		}
		if waitErr != nil {
			t.Errorf("Wait() error = %v, want nil", waitErr)
		}
	})
}

// fanOutProbe records how many calls to hold overlap. A call stays in flight
// until more than limit calls overlap (then every call is released at once, so
// an unbounded fan-out is seen without waiting) or a short grace period ends
// (so a bounded fan-out still makes progress).
type fanOutProbe struct {
	over   chan struct{}
	mu     sync.Mutex
	limit  int
	cur    int
	peak   int
	calls  int
	closed bool
}

func (p *fanOutProbe) hold() {
	p.mu.Lock()
	p.cur++
	p.calls++
	p.peak = max(p.peak, p.cur)
	if p.cur > p.limit && !p.closed {
		p.closed = true
		close(p.over)
	}
	p.mu.Unlock()
	timer := time.NewTimer(20 * time.Millisecond)
	select {
	case <-p.over:
	case <-timer.C:
	}
	timer.Stop()
	p.mu.Lock()
	p.cur--
	p.mu.Unlock()
}

func (p *fanOutProbe) stats() (peak, calls int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.peak, p.calls
}

// TestNewBoundedFanOut pins the group the Multi* handlers fan out on: it never
// runs more than limit items at once, limit <= 0 means unbounded (not a
// deadlock), and neither an already-cancelled request nor a failing item makes
// it silently skip the items that had to wait for a slot; the item error is
// still returned by Wait.
func TestNewBoundedFanOut(t *testing.T) {
	t.Parallel()
	const n = 6
	errItem := errors.New("item failed")
	tests := []struct {
		name      string
		limit     int
		cancelled bool
		fails     bool
	}{
		{name: "bounded", limit: 2},
		{name: "bounded with a cancelled request ctx", limit: 2, cancelled: true},
		{name: "bounded with a failing item", limit: 2, fails: true},
		{name: "zero limit is unbounded", limit: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if test.cancelled {
				cancel()
			}
			p := &fanOutProbe{over: make(chan struct{}), limit: 2}
			done := make(chan error, 1)
			go func() {
				eg := NewBoundedFanOut(ctx, test.limit)
				for i := range n {
					eg.Go(func() error {
						p.hold()
						if test.fails && i == 0 {
							return errItem
						}
						return nil
					})
				}
				done <- eg.Wait()
			}()
			select {
			case err := <-done:
				if test.fails {
					if !errors.Is(err, errItem) {
						t.Errorf("Wait() = %v, want %v", err, errItem)
					}
				} else if err != nil {
					t.Fatalf("Wait() = %v, want nil", err)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timed out waiting for the group")
			}
			peak, calls := p.stats()
			if calls != n {
				t.Errorf("ran %d items, want all %d", calls, n)
			}
			if test.limit > 0 && peak > test.limit {
				t.Errorf("peak concurrency = %d, want <= %d", peak, test.limit)
			}
		})
	}
}

// recordingLogger discards everything except Errorf messages, which it records
// so a test can count them by a unique error message.
type recordingLogger struct {
	logger.Interface
	errors []string
	mu     sync.Mutex
}

//nolint:gochecknoglobals // TestMain installs this as the process-wide log.Init logger, so tests must share the one instance to read what the group logged
var recLog = &recordingLogger{Interface: nop.New()}

func (r *recordingLogger) Errorf(format string, vals ...any) {
	msg := fmt.Sprintf(format, vals...)
	r.mu.Lock()
	r.errors = append(r.errors, msg)
	r.mu.Unlock()
}

// count returns how many recorded Errorf messages contain s.
func (r *recordingLogger) count(s string) (n int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.errors {
		if strings.Contains(m, s) {
			n++
		}
	}
	return n
}

// A task error on the process-wide group is logged when it occurs (its
// cancellation is watched by nobody and Wait only runs at shutdown); a local
// group keeps leaving reporting to its caller, and cancellation is never logged.
// Not parallel: it resets the package-level instance.
func TestGlobalGroupLogsTaskErrors(t *testing.T) {
	tests := []struct {
		group   func(t *testing.T) Group
		taskErr error
		name    string
		want    int
	}{
		{
			name: "global group logs a failed task",
			group: func(t *testing.T) Group {
				t.Helper()
				instance, once = nil, sync.Once{}
				t.Cleanup(func() { instance, once = nil, sync.Once{} })
				Init(t.Context())
				return Get()
			},
			taskErr: errors.New("TestGlobalGroupLogsTaskErrors/global failure"),
			want:    1,
		},
		{
			name: "local group does not log a failed task",
			group: func(t *testing.T) Group {
				t.Helper()
				g, _ := New(t.Context())
				return g
			},
			taskErr: errors.New("TestGlobalGroupLogsTaskErrors/local failure"),
			want:    0,
		},
		{
			name: "global group does not log a canceled task",
			group: func(t *testing.T) Group {
				t.Helper()
				instance, once = nil, sync.Once{}
				t.Cleanup(func() { instance, once = nil, sync.Once{} })
				Init(t.Context())
				return Get()
			},
			taskErr: errors.Errorf("TestGlobalGroupLogsTaskErrors/canceled: %w", context.Canceled),
			want:    0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.group(t)
			// recLog is process-wide and keeps its records, so a repeated run (-count, -shuffle)
			// sees the messages of the earlier ones; count what this run adds.
			before := recLog.count(tc.taskErr.Error())
			g.Go(func() error { return tc.taskErr })
			// Wait returns once the task has run; its result is not under test.
			if err := g.Wait(); err != nil {
				t.Logf("Wait() error = %v", err)
			}
			if got := recLog.count(tc.taskErr.Error()) - before; got != tc.want {
				t.Errorf("logged %d times, want %d", got, tc.want)
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
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 		limit  atomic.Int64
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 				limit:  test.fields.limit,
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
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 		limit  atomic.Int64
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 				limit:  test.fields.limit,
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
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 		limit  atomic.Int64
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 				limit:  test.fields.limit,
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
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 		limit  atomic.Int64
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 				limit:  test.fields.limit,
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
// 		sem    *semaphore.Weighted
// 		emap   map[string]struct{}
// 		errs   []error
// 		limit  atomic.Int64
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 		           sem:nil,
// 		           emap:nil,
// 		           errs:nil,
// 		           limit:nil,
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
// 				sem:    test.fields.sem,
// 				emap:   test.fields.emap,
// 				errs:   test.fields.errs,
// 				limit:  test.fields.limit,
// 			}
//
// 			g.done()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// appendErr used to check and insert under separate locks, so concurrent goroutines
// returning the same error could all record it.
func Test_group_appendErr_concurrentDuplicates(t *testing.T) {
	t.Parallel()
	g := new(group)
	const workers = 64
	var wg sync.WaitGroup
	start := make(chan struct{})
	for range workers {
		wg.Go(func() {
			<-start
			g.appendErr(errors.New("same"))
		})
	}
	close(start)
	wg.Wait()
	if got := len(g.errs); got != 1 {
		t.Errorf("recorded %d errors, want 1", got)
	}
}

// TestNewContinueOnError pins that a task error does not cancel the context of
// a NewContinueOnError group, so tasks that wait for a slot after the failed
// one still run, while Wait still returns the error.
func TestNewContinueOnError(t *testing.T) {
	t.Parallel()
	const n = 4
	errTask := errors.New("task failed")
	g, egctx := NewContinueOnError(t.Context())
	// limit 1 runs the failing task inline, so its effect on egctx is settled
	// when Go returns.
	g.SetLimit(1)
	g.Go(func() error { return errTask })
	if err := egctx.Err(); err != nil {
		t.Fatalf("egctx.Err() = %v after a task error, want nil", err)
	}
	g.SetLimit(2)
	var ran atomic.Int64
	for range n {
		g.Go(func() error {
			ran.Add(1)
			return nil
		})
	}
	done := make(chan error, 1)
	go func() { done <- g.Wait() }()
	select {
	case err := <-done:
		if !errors.Is(err, errTask) {
			t.Errorf("Wait() = %v, want %v", err, errTask)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for the group")
	}
	if got := ran.Load(); got != n {
		t.Errorf("ran %d tasks, want all %d", got, n)
	}
}
