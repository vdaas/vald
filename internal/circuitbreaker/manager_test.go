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
package circuitbreaker

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

func TestBreakerStateError_Parity(t *testing.T) {
	t.Parallel()
	base := errors.New("not found")
	lazy := error(&breakerStateError{err: base, st: StateClosed})
	eager := errors.Wrapf(base, "circuitbreaker state is %s, this error is not caused by circuitbreaker", StateClosed.String())
	if lazy.Error() != eager.Error() {
		t.Errorf("Error() = %q, want %q", lazy.Error(), eager.Error())
	}
	if !errors.Is(lazy, base) {
		t.Errorf("errors.Is(lazy, base) = false")
	}
}

// Do counts every call per breaker name and state; the counters are lock-free
// per key, so concurrent calls must still be counted exactly.
func TestBreakerManager_Do_countsMetrics(t *testing.T) {
	t.Parallel()
	cb, err := New()
	if err != nil {
		t.Fatal(err)
	}
	const (
		key   = "TestBreakerManager_Do_countsMetrics"
		calls = 200
	)
	// The metrics are process-wide and keep their counts, so a repeated run (-count, -shuffle)
	// starts from the totals of the earlier ones; count what this run adds.
	before := Metrics(t.Context())[key][StateClosed]
	var wg sync.WaitGroup
	for range calls {
		wg.Go(func() {
			if _, err := cb.Do(t.Context(), key, func(context.Context) (any, error) {
				return struct{}{}, nil
			}); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if got := Metrics(t.Context())[key][StateClosed] - before; got != calls {
		t.Errorf("Metrics()[%q][closed] grew by %d, want %d", key, got, calls)
	}
}

// EvictAddr removes the breakers and counters of every key that targets the
// address ("<addr>" or "<method>/<addr>") and keeps every other key, including
// keys that only end with the same characters.
func TestBreakerManager_EvictAddr(t *testing.T) {
	t.Parallel()
	const (
		method = "TestBreakerManager_EvictAddr"
		addr   = "10.0.0.1:8081"
	)
	evicted := []string{addr, method + "/" + addr, method + "/192.0.2.1:8081/" + addr}
	kept := []string{
		method + "/1" + addr,
		method + "/" + addr + "/192.0.2.1:8081",
		method + "/192.0.2.1:8081",
	}
	cb, err := New()
	if err != nil {
		t.Fatal(err)
	}
	bm, ok := cb.(*breakerManager)
	if !ok {
		t.Fatalf("New returned %T, want *breakerManager", cb)
	}
	for _, key := range append(slices.Clone(evicted), kept...) {
		if _, err := cb.Do(t.Context(), key, func(context.Context) (any, error) {
			return true, nil
		}); err != nil {
			t.Fatal(err)
		}
	}

	bm.EvictAddr(addr)

	ms := Metrics(t.Context())
	for _, key := range evicted {
		if _, ok := bm.m.Load(key); ok {
			t.Errorf("breaker %q remains after EvictAddr", key)
		}
		if _, ok := ms[key]; ok {
			t.Errorf("metrics for %q remain after EvictAddr", key)
		}
	}
	for _, key := range kept {
		if _, ok := bm.m.Load(key); !ok {
			t.Errorf("breaker %q was evicted", key)
		}
		if _, ok := ms[key]; !ok {
			t.Errorf("metrics for %q were evicted", key)
		}
	}
}

// A call still in flight when EvictAddr runs must not re-create the evicted
// key's counters when it finishes; otherwise the state of removed addresses
// keeps piling up under constant traffic.
func TestBreakerManager_EvictAddr_inFlightCall(t *testing.T) {
	t.Parallel()
	const addr = "10.0.0.2:8081"
	key := "TestBreakerManager_EvictAddr_inFlightCall/" + addr
	cb, err := New()
	if err != nil {
		t.Fatal(err)
	}
	bm, ok := cb.(*breakerManager)
	if !ok {
		t.Fatalf("New returned %T, want *breakerManager", cb)
	}
	started, release := make(chan struct{}), make(chan struct{})
	done := make(chan error, 1)
	go func() {
		_, err := cb.Do(t.Context(), key, func(context.Context) (any, error) {
			close(started)
			<-release
			return true, nil
		})
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(10 * time.Second):
		t.Fatal("call did not start")
	}
	bm.EvictAddr(addr)
	close(release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("call did not finish")
	}
	if _, ok := Metrics(t.Context())[key]; ok {
		t.Errorf("metrics for %q were re-created by an in-flight call after EvictAddr", key)
	}
	if _, ok := bm.m.Load(key); ok {
		t.Errorf("breaker %q remains after EvictAddr", key)
	}
}

func TestBreakerManager_Handle(t *testing.T) {
	t.Parallel()

	t.Run("Resolve Allow FinishSuccess FinishError and IsEvicted parity", func(t *testing.T) {
		t.Parallel()
		const (
			addr = "10.0.0.11:8081"
			key  = "TestBreakerManager_Handle/" + addr
		)
		cb, err := New(
			WithBreakerOpts(
				WithClosedErrorRate(0.5),
				WithMinSamples(1),
				WithOpenTimeout("100ms"),
			),
		)
		if err != nil {
			t.Fatal(err)
		}
		bm := cb.(*breakerManager)
		h, err := bm.Resolve(key)
		if err != nil {
			t.Fatalf("Resolve() error = %v", err)
		}
		if h.IsEvicted() {
			t.Fatal("newly resolved handle should not be evicted")
		}
		beforeClosed := Metrics(t.Context())[key][StateClosed]

		if err := h.Allow(); err != nil {
			t.Fatalf("Allow() error = %v", err)
		}
		h.FinishSuccess()

		// Non-retryable error must not trip the breaker even with minSamples=1.
		nonRetryErr := errors.New("not found")
		wrapped := h.FinishError(nonRetryErr, false)
		if !errors.Is(wrapped, nonRetryErr) {
			t.Fatalf("FinishError(false) = %v, want wrapping %v", wrapped, nonRetryErr)
		}
		if err := h.Allow(); err != nil {
			t.Fatalf("Allow() after non-retryable error = %v", err)
		}

		// MarkWithSuccess wrapper must count as success and unwrap.
		markErr := errors.New("ok-mark")
		if got := h.FinishError(errors.NewErrCircuitBreakerMarkWithSuccess(markErr), true); !errors.Is(got, markErr) {
			t.Fatalf("FinishError(MarkWithSuccess) = %v, want wrapping %v", got, markErr)
		}

		// Ignorable wrapper must not trip even when retryable=true.
		igErr := errors.New("ignorable")
		if got := h.FinishError(errors.NewErrCircuitBreakerIgnorable(igErr), true); !errors.Is(got, igErr) {
			t.Fatalf("FinishError(Ignorable) = %v, want wrapping %v", got, igErr)
		}

		// Context sentinels must not trip even when retryable=true.
		if got := h.FinishError(context.Canceled, true); !errors.Is(got, context.Canceled) {
			t.Fatalf("FinishError(Canceled) = %v, want Canceled", got)
		}

		// Retryable failures should trip the breaker once minSamples and error rate exceed threshold.
		retryErr := errors.New("unavailable")
		for range 10 {
			if err := h.Allow(); err != nil {
				break
			}
			if ferr := h.FinishError(retryErr, true); !errors.Is(ferr, retryErr) {
				t.Fatalf("FinishError() = %v, want %v", ferr, retryErr)
			}
		}
		if err := h.Allow(); !errors.Is(err, errors.ErrCircuitBreakerOpenState) {
			t.Fatalf("Allow() after retryable failures = %v, want ErrCircuitBreakerOpenState", err)
		}
		if got := Metrics(t.Context())[key][StateClosed] - beforeClosed; got == 0 {
			t.Errorf("expected closed metrics to increase")
		}

		bm.EvictAddr(addr)
		if !h.IsEvicted() {
			t.Error("handle should report IsEvicted() == true after EvictAddr")
		}
	})

	t.Run("nil handle is safe", func(t *testing.T) {
		t.Parallel()
		var h *Handle
		if !h.IsEvicted() {
			t.Error("nil handle IsEvicted() = false, want true")
		}
		if err := h.Allow(); err != nil {
			t.Errorf("nil handle Allow() = %v, want nil", err)
		}
		h.FinishSuccess()
		base := errors.New("boom")
		if got := h.FinishError(base, true); !errors.Is(got, base) {
			t.Errorf("nil handle FinishError() = %v, want %v", got, base)
		}
		if got := h.FinishError(nil, true); got != nil {
			t.Errorf("nil handle FinishError(nil) = %v, want nil", got)
		}
	})
}

type customAsError struct {
	target any
	inner  error
	msg    string
}

func (e *customAsError) Error() string { return e.msg }

func (e *customAsError) Unwrap() error { return e.inner }

func (e *customAsError) As(dst any) bool {
	switch p := dst.(type) {
	case **errors.ErrCircuitBreakerMarkWithSuccess:
		if s, ok := e.target.(*errors.ErrCircuitBreakerMarkWithSuccess); ok {
			*p = s
			return true
		}
	case **errors.ErrCircuitBreakerIgnorable:
		if ig, ok := e.target.(*errors.ErrCircuitBreakerIgnorable); ok {
			*p = ig
			return true
		}
	}
	return false
}

func TestBreakerManager_FinishError_WrappedAndZeroAlloc(t *testing.T) {
	t.Parallel()

	base := errors.New("not found")
	markBase := errors.New("mark-success")
	igBase := errors.New("ignorable")

	tests := []struct {
		err        error
		wantUnwrap error
		name       string
		retryable  bool
		wantTrip   bool
		wantMark   bool
	}{
		{
			name:       "fmt wrapped MarkWithSuccess",
			err:        fmt.Errorf("outer: %w", errors.NewErrCircuitBreakerMarkWithSuccess(markBase)),
			retryable:  true,
			wantUnwrap: markBase,
			wantMark:   true,
		},
		{
			name:       "fmt wrapped Ignorable",
			err:        fmt.Errorf("outer: %w", errors.NewErrCircuitBreakerIgnorable(igBase)),
			retryable:  true,
			wantUnwrap: igBase,
		},
		{
			name:       "errors.Join MarkWithSuccess overrides Ignorable",
			err:        errors.Join(errors.NewErrCircuitBreakerIgnorable(igBase), errors.NewErrCircuitBreakerMarkWithSuccess(markBase)),
			retryable:  true,
			wantUnwrap: markBase,
			wantMark:   true,
		},
		{
			name:       "direct Ignorable wrapping MarkWithSuccess honors MarkWithSuccess precedence",
			err:        errors.NewErrCircuitBreakerIgnorable(errors.NewErrCircuitBreakerMarkWithSuccess(markBase)),
			retryable:  true,
			wantUnwrap: markBase,
			wantMark:   true,
		},
		{
			name:       "Ignorable wrapping MarkWithSuccess honors MarkWithSuccess precedence",
			err:        fmt.Errorf("outer: %w", errors.NewErrCircuitBreakerIgnorable(errors.NewErrCircuitBreakerMarkWithSuccess(markBase))),
			retryable:  true,
			wantUnwrap: markBase,
			wantMark:   true,
		},
		{
			name:       "custom As MarkWithSuccess",
			err:        &customAsError{msg: "custom-mark", target: errors.NewErrCircuitBreakerMarkWithSuccess(markBase)},
			retryable:  true,
			wantUnwrap: markBase,
			wantMark:   true,
		},
		{
			name:       "custom As Ignorable with inner MarkWithSuccess",
			err:        &customAsError{msg: "custom-ig", target: errors.NewErrCircuitBreakerIgnorable(igBase), inner: errors.NewErrCircuitBreakerMarkWithSuccess(markBase)},
			retryable:  true,
			wantUnwrap: markBase,
			wantMark:   true,
		},
		{
			name:       "custom As Ignorable only",
			err:        &customAsError{msg: "custom-ig-only", target: errors.NewErrCircuitBreakerIgnorable(igBase)},
			retryable:  true,
			wantUnwrap: igBase,
		},
		{
			name:       "wrapped DeadlineExceeded does not trip",
			err:        fmt.Errorf("rpc timeout: %w", context.DeadlineExceeded),
			retryable:  true,
			wantUnwrap: context.DeadlineExceeded,
		},
		{
			name:       "wrapped non-retryable error does not trip",
			err:        fmt.Errorf("wrapped: %w", base),
			retryable:  false,
			wantUnwrap: base,
		},
		{
			name:       "retryable normal error trips breaker at minSamples=1",
			err:        fmt.Errorf("wrapped: %w", base),
			retryable:  true,
			wantUnwrap: base,
			wantTrip:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cb, err := New(
				WithBreakerOpts(
					WithClosedErrorRate(0.5),
					WithMinSamples(1),
					WithOpenTimeout("100ms"),
				),
			)
			if err != nil {
				t.Fatal(err)
			}
			bm := cb.(*breakerManager)
			h, err := bm.Resolve("TestBreakerManager_FinishError_Wrapped/" + tc.name)
			if err != nil {
				t.Fatal(err)
			}
			// Prime one sample so minSamples=1 threshold (> minSamples) trips on failure.
			if err := h.Allow(); err != nil {
				t.Fatal(err)
			}
			h.FinishSuccess()
			succBefore := h.b.count.Successes()

			if err := h.Allow(); err != nil {
				t.Fatal(err)
			}
			got := h.FinishError(tc.err, tc.retryable)
			if !errors.Is(got, tc.wantUnwrap) {
				t.Fatalf("FinishError() = %v, want wrapping %v", got, tc.wantUnwrap)
			}
			if tc.wantMark && h.b.count.Successes() != succBefore+1 {
				t.Fatalf("Successes() = %d, want %d", h.b.count.Successes(), succBefore+1)
			}
			if tc.wantTrip {
				if err := h.Allow(); !errors.Is(err, errors.ErrCircuitBreakerOpenState) {
					t.Fatalf("Allow() after retryable failure = %v, want ErrCircuitBreakerOpenState", err)
				}
			} else if err := h.Allow(); err != nil {
				t.Fatalf("Allow() unexpectedly errored: %v", err)
			}
		})
	}
}

func TestClassifyWrappedBreakerErr_ZeroAlloc(t *testing.T) {
	leaf := errors.New("rpc error: code = NotFound desc = not found")
	wrapped := fmt.Errorf("rpc wrapped: %w", leaf)
	igDirect := errors.NewErrCircuitBreakerIgnorable(leaf)
	igWrapped := fmt.Errorf("ig wrapped: %w", errors.NewErrCircuitBreakerIgnorable(leaf))

	for _, tc := range []struct {
		err  error
		name string
		kind breakerErrKind
	}{
		{name: "leaf", err: leaf, kind: breakerErrNone},
		{name: "wrapped", err: wrapped, kind: breakerErrNone},
		{name: "ignorable_direct", err: igDirect, kind: breakerErrIgnorable},
		{name: "ignorable_wrapped", err: igWrapped, kind: breakerErrIgnorable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var gotKind breakerErrKind
			allocs := testing.AllocsPerRun(100, func() {
				gotKind, _ = classifyWrappedBreakerErr(tc.err)
			})
			if gotKind != tc.kind {
				t.Fatalf("kind = %v, want %v", gotKind, tc.kind)
			}
			if allocs != 0 {
				t.Fatalf("allocs = %v, want 0", allocs)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want CircuitBreaker
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		want
// 		checkFunc  func(want, CircuitBreaker, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got CircuitBreaker, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           opts:nil,
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
// 		           opts:nil,
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
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_breakerManager_Do(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		key string
// 		fn  func(ctx context.Context) (any, error)
// 	}
// 	type fields struct {
// 		m    sync.Map[string, *breaker]
// 		opts []BreakerOption
// 	}
// 	type want struct {
// 		wantVal any
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		fields
// 		want
// 		checkFunc  func(want, any, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVal any, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVal, w.wantVal) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVal, w.wantVal)
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
// 		           key:"",
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           m:nil,
// 		           opts:nil,
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
// 		           key:"",
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           m:nil,
// 		           opts:nil,
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
// 			bm := &breakerManager{
// 				m:    test.fields.m,
// 				opts: test.fields.opts,
// 			}
//
// 			gotVal, err := bm.Do(test.args.ctx, test.args.key, test.args.fn)
// 			if err := checkFunc(test.want, gotVal, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestMetrics(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type want struct {
// 		wantMs map[string]map[State]int64
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		want
// 		checkFunc  func(want, map[string]map[State]int64) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotMs map[string]map[State]int64) error {
// 		if !reflect.DeepEqual(gotMs, w.wantMs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotMs, w.wantMs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
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
// 		           in0:nil,
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
// 			gotMs := Metrics(test.args.in0)
// 			if err := checkFunc(test.want, gotMs); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
