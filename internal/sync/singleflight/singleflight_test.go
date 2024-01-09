//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package singleflight represents zero time caching
package singleflight

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	type want struct {
		want Group[any]
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, Group[any]) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got Group[any]) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns Group implementation",
			want: want{
				want: &group[any]{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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

			got := New[any]()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_group_Do(t *testing.T) {
	type args[V any] struct {
		ctx context.Context
		key string
		fn  func(context.Context) (V, error)
	}
	type want[V any] struct {
		wantV      V
		wantShared bool
		err        error
	}
	type test[V any] struct {
		name       string
		args       args[V]
		want       want[V]
		beforeFunc func(args[V])
		execFunc   func(*testing.T, args[V]) (V, bool, error)
		checkFunc  func(want[V], V, bool, error) error
		afterFunc  func(args[V])
	}
	tests := []test[string]{
		func() test[string] {
			// routine1
			key1 := "req_1"
			var cnt1 uint32

			// the unparam lint rule is disabled here because we need to match the interface to singleflight implementation.
			// if this rule is not disabled, if will warns that the error will always return null.
			//nolint:unparam
			fn1 := func(context.Context) (string, error) {
				atomic.AddUint32(&cnt1, 1)
				return "res_1", nil
			}

			// routine 2
			key2 := "req_2"
			var cnt2 uint32

			// the unparam lint rule is disabled here because we need to match the interface to singleflight implementation.
			// if this rule is not disabled, if will warns that the error will always return null.
			//nolint:unparam
			fn2 := func(context.Context) (string, error) {
				atomic.AddUint32(&cnt2, 1)
				return "res_2", nil
			}

			return test[string]{
				name: "returns (v, false, nil) when Do is called with another key",
				args: args[string]{
					key: key1,
					ctx: context.Background(),
					fn:  fn1,
				},
				want: want[string]{
					wantV:      "res_1",
					wantShared: false,
					err:        nil,
				},
				execFunc: func(t *testing.T, a args[string]) (got string, gotShared bool, err error) {
					t.Helper()
					g := New[string]()

					wg := new(sync.WaitGroup)
					wg.Add(1)
					go func() {
						got, gotShared, err = g.Do(a.ctx, a.key, a.fn)
						wg.Done()
					}()

					wg.Add(1)
					go func() {
						_, _, _ = g.Do(a.ctx, key2, fn2)
						wg.Done()
					}()

					wg.Wait()
					return got, gotShared, err
				},
				checkFunc: func(w want[string], gotV string, gotShared bool, err error) error {
					if got, want := int(atomic.LoadUint32(&cnt1)), 1; got != want {
						return errors.Errorf("cnt got = %d, want = %d", got, want)
					}
					if got, want := int(atomic.LoadUint32(&cnt2)), 1; got != want {
						return errors.Errorf("cnt got = %d, want = %d", got, want)
					}
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					if !reflect.DeepEqual(gotV, w.wantV) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotV, w.wantV)
					}
					if !reflect.DeepEqual(gotShared, w.wantShared) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotShared, w.wantShared)
					}
					return nil
				},
			}
		}(),
		func() test[string] {
			// routine1
			var cnt1 uint32

			// the unparam lint rule is disabled here because we need to match the interface to singleflight implementation.
			// if this rule is not disabled, if will warns that the error will always return null.
			//nolint:unparam
			fn1 := func(context.Context) (string, error) {
				atomic.AddUint32(&cnt1, 1)
				time.Sleep(time.Millisecond * 500)
				return "res_1", nil
			}

			// routine 2
			var cnt2 uint32

			// the unparam lint rule is disabled here because we need to match the interface to singleflight implementation.
			// if this rule is not disabled, if will warns that the error will always return null.
			//nolint:unparam
			fn2 := func(context.Context) (string, error) {
				atomic.AddUint32(&cnt2, 1)
				return "res_2", nil
			}

			w := want[string]{
				wantV:      "res_1",
				wantShared: true,
				err:        nil,
			}

			checkFunc := func(w want[string], gotV string, gotShared bool, err error) error {
				c1 := int(atomic.LoadUint32(&cnt1))
				c2 := int(atomic.LoadUint32(&cnt2))
				// since there is a chance that the go routine 2 is executed before routine 1, we need to check if either one is executed
				if !((c1 == 1 && c2 == 0) || (c1 == 0 && c2 == 1)) {
					return errors.Errorf("cnt1 and cnt2 is executed, %d, %d", c1, c2)
				}
				if !errors.Is(err, w.err) {
					return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
				}
				if !reflect.DeepEqual(gotV, w.wantV) {
					return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotV, w.wantV)
				}
				if !reflect.DeepEqual(gotShared, w.wantShared) {
					return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotShared, w.wantShared)
				}
				return nil
			}

			return test[string]{
				name: "returns (v, true, nil) when Do is called with the same key",
				args: args[string]{
					key: "req_1",
					ctx: context.Background(),
					fn:  fn1,
				},
				want: w,
				execFunc: func(t *testing.T, a args[string]) (string, bool, error) {
					t.Helper()

					g := New[string]()
					wg := new(sync.WaitGroup)
					var got, got1 string
					var gotShared, gotShared1 bool
					var err, err1 error

					wg.Add(1)
					go func() {
						got, gotShared, err = g.Do(a.ctx, a.key, fn1)
						wg.Done()
					}()

					// call with the same key but with another function
					wg.Add(1)
					time.Sleep(time.Millisecond * 100)
					go func() {
						got1, gotShared1, err1 = g.Do(a.ctx, a.key, fn2)
						wg.Done()
					}()

					wg.Wait()

					if err := checkFunc(w, got1, gotShared1, err1); err != nil {
						t.Fatal(err)
					}

					return got, gotShared, err
				},
				checkFunc: checkFunc,
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}

			gotV, gotShared, err := test.execFunc(t, test.args)

			if err := test.checkFunc(test.want, gotV, gotShared, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDo(t *testing.T) {
	g := New[string]()
	v, _, err := g.Do(context.Background(), "key", func(context.Context) (string, error) {
		return "bar", nil
	})
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}

func TestDoErr(t *testing.T) {
	g := New[any]()
	someErr := errors.New("Some error")
	v, _, err := g.Do(context.Background(), "key", func(context.Context) (any, error) {
		return nil, someErr
	})
	if err != someErr {
		t.Errorf("Do error = %v; want someErr %v", err, someErr)
	}
	if v != nil {
		t.Errorf("unexpected non-nil value %#v", v)
	}
}

func TestDoDupSuppress(t *testing.T) {
	g := New[string]()
	var wg1, wg2 sync.WaitGroup
	c := make(chan string, 1)
	var calls int32
	fn := func(context.Context) (string, error) {
		if atomic.AddInt32(&calls, 1) == 1 {
			// First invocation.
			wg1.Done()
		}
		v := <-c
		c <- v // pump; make available for any future calls

		time.Sleep(10 * time.Millisecond) // let more goroutines enter Do

		return v, nil
	}

	const n = 10
	wg1.Add(1)
	for i := 0; i < n; i++ {
		wg1.Add(1)
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			wg1.Done()
			s, _, err := g.Do(context.Background(), "key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
				return
			}
			if s != "bar" {
				t.Errorf("Do = %T %v; want %q", s, s, "bar")
			}
		}()
	}
	wg1.Wait()
	// At least one goroutine is in fn now and all of them have at
	// least reached the line before the Do.
	c <- "bar"
	wg2.Wait()
	if got := atomic.LoadInt32(&calls); got <= 0 || got >= n {
		t.Errorf("number of calls = %d; want over 0 and less than %d", got, n)
	}
}

// Test that singleflight behaves correctly after Forget called.
// See https://github.com/golang/go/issues/31420
func TestForget(t *testing.T) {
	g := New[int]()

	var (
		firstStarted  = make(chan struct{})
		unblockFirst  = make(chan struct{})
		firstFinished = make(chan struct{})
	)

	go func() {
		g.Do(context.Background(), "key", func(ctx context.Context) (i int, e error) {
			close(firstStarted)
			<-unblockFirst
			close(firstFinished)
			return
		})
	}()
	<-firstStarted
	g.Forget("key")

	unblockSecond := make(chan struct{})
	secondResult := g.DoChan(context.Background(), "key", func(ctx context.Context) (i int, e error) {
		t.Log(2, "key")
		<-unblockSecond
		return 2, nil
	})

	close(unblockFirst)
	<-firstFinished

	thirdResult := g.DoChan(context.Background(), "key", func(ctx context.Context) (i int, e error) {
		t.Log(3, "key")
		return 3, nil
	})

	close(unblockSecond)
	<-secondResult
	r := <-thirdResult
	if r.Val != 2 {
		t.Errorf("We should receive result produced by second call, expected: 2, got %d", r.Val)
	}
}

func TestDoChan(t *testing.T) {
	g := New[string]()
	ch := g.DoChan(context.Background(), "key", func(ctx context.Context) (string, error) {
		return "bar", nil
	})

	res := <-ch
	v := res.Val
	err := res.Err
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}

// Test singleflight behaves correctly after Do panic.
// See https://github.com/golang/go/issues/41133
func TestPanicDo(t *testing.T) {
	g := New[any]()
	fn := func(ctx context.Context) (any, error) {
		panic("invalid memory address or nil pointer dereference")
	}

	const n = 5
	waited := int32(n)
	panicCount := int32(0)
	done := make(chan struct{})
	for i := 0; i < n; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					t.Logf("Got panic: %v\n%s", err, info.Get().String())
					atomic.AddInt32(&panicCount, 1)
				}

				if atomic.AddInt32(&waited, -1) == 0 {
					close(done)
				}
			}()

			g.Do(context.Background(), "key", fn)
		}()
	}

	select {
	case <-done:
		if panicCount != n {
			t.Errorf("Expect %d panic, but got %d", n, panicCount)
		}
	case <-time.After(time.Second):
		t.Fatalf("Do hangs")
	}
}

func TestGoexitDo(t *testing.T) {
	g := New[any]()
	fn := func(context.Context) (interface{}, error) {
		runtime.Goexit()
		return nil, nil
	}

	const n = 5
	waited := int32(n)
	done := make(chan struct{})
	for i := 0; i < n; i++ {
		go func() {
			var err error
			defer func() {
				if err != nil {
					t.Errorf("Error should be nil, but got: %v", err)
				}
				if atomic.AddInt32(&waited, -1) == 0 {
					close(done)
				}
			}()
			_, _, err = g.Do(context.Background(), "key", fn)
		}()
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("Do hangs")
	}
}

func executable(tb testing.TB) string {
	tb.Helper()
	exe, err := os.Executable()
	if err != nil {
		tb.Skipf("skipping: test executable not found")
	}

	// Control case: check whether exec.Command works at all.
	// (For example, it might fail with a permission error on iOS.)
	cmd := exec.Command(exe, "-test.list=^$")
	cmd.Env = []string{}
	if err := cmd.Run(); err != nil {
		tb.Skipf("skipping: exec appears not to work on %s: %v", runtime.GOOS, err)
	}

	return exe
}

func TestPanicDoChan(t *testing.T) {
	if os.Getenv("TEST_PANIC_DOCHAN") != "" {
		defer func() {
			recover()
		}()

		g := New[any]()
		ch := g.DoChan(context.Background(), "", func(context.Context) (interface{}, error) {
			panic("Panicking in DoChan")
		})
		<-ch
		t.Fatalf("DoChan unexpectedly returned")
	}

	t.Parallel()

	cmd := exec.Command(executable(t), "-test.run="+t.Name(), "-test.v")
	cmd.Env = append(os.Environ(), "TEST_PANIC_DOCHAN=1")
	out := new(bytes.Buffer)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	err := cmd.Wait()
	t.Logf("%s:\n%s", strings.Join(cmd.Args, " "), out)
	if err == nil {
		t.Errorf("Test subprocess passed; want a crash due to panic in DoChan")
	}
	if bytes.Contains(out.Bytes(), []byte("DoChan unexpectedly")) {
		t.Errorf("Test subprocess failed with an unexpected failure mode.")
	}
	if !bytes.Contains(out.Bytes(), []byte("Panicking in DoChan")) {
		t.Errorf("Test subprocess failed, but the crash isn't caused by panicking in DoChan")
	}
}

func TestPanicDoSharedByDoChan(t *testing.T) {
	if os.Getenv("TEST_PANIC_DOCHAN") != "" {
		blocked := make(chan struct{})
		unblock := make(chan struct{})

		g := New[any]()
		go func() {
			defer func() {
				recover()
			}()
			g.Do(context.Background(), "", func(context.Context) (interface{}, error) {
				close(blocked)
				<-unblock
				panic("Panicking in Do")
			})
		}()

		<-blocked
		ch := g.DoChan(context.Background(), "", func(context.Context) (interface{}, error) {
			panic("DoChan unexpectedly executed callback")
		})
		close(unblock)
		<-ch
		t.Fatalf("DoChan unexpectedly returned")
	}

	t.Parallel()

	cmd := exec.Command(executable(t), "-test.run="+t.Name(), "-test.v")
	cmd.Env = append(os.Environ(), "TEST_PANIC_DOCHAN=1")
	out := new(bytes.Buffer)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	err := cmd.Wait()
	t.Logf("%s:\n%s", strings.Join(cmd.Args, " "), out)
	if err == nil {
		t.Errorf("Test subprocess passed; want a crash due to panic in Do shared by DoChan")
	}
	if bytes.Contains(out.Bytes(), []byte("DoChan unexpectedly")) {
		t.Errorf("Test subprocess failed with an unexpected failure mode.")
	}
	if !bytes.Contains(out.Bytes(), []byte("Panicking in Do")) {
		t.Errorf("Test subprocess failed, but the crash isn't caused by panicking in Do")
	}
}

func ExampleGroup() {
	g := New[string]()

	block := make(chan struct{})
	res1c := g.DoChan(context.Background(), "key", func(context.Context) (string, error) {
		<-block
		return "func 1", nil
	})
	res2c := g.DoChan(context.Background(), "key", func(context.Context) (string, error) {
		<-block
		return "func 2", nil
	})
	close(block)

	res1 := <-res1c
	res2 := <-res2c

	// Results are shared by functions executed with duplicate keys.
	fmt.Println("Shared:", res2.Shared)
	// Only the first function is executed: it is registered and started with "key",
	// and doesn't complete before the second function is registered with a duplicate key.
	fmt.Println("Equal results:", res1.Val == res2.Val)
	fmt.Println("Result:", res1.Val)

	// Output:
	// Shared: true
	// Equal results: true
	// Result: func 1
}

// NOT IMPLEMENTED BELOW
