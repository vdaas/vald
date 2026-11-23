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
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

// NOTE: This variable is for observability package.
//
//	This will be fixed when refactoring the observability package.
var (
	// metrics counts Do calls per breaker name and resulting state. Do runs on
	// every gRPC call that has a breaker, so the hot-path lookup and counters
	// are lock-free; metricsMu only serializes cold misses and EvictAddr.
	metrics sync.Map[string, *stateCounts] //nolint:gochecknoglobals // process-wide counters shared by every breakerManager and read by the observability package; see the NOTE above
	//nolint:gochecknoglobals // serializes cold breaker/metrics registration against EvictAddr
	metricsMu sync.Mutex
)

// CircuitBreaker is a state machine to prevent doing processes that are likely to fail.
type CircuitBreaker interface {
	Do(ctx context.Context, key string, fn func(ctx context.Context) (any, error)) (val any, err error)
}

// stateCounts holds one counter per State (StateUnknown..StateHalfOpen) padded to a 64B cache line
// to avoid false sharing across per-target breakers.
type stateCounts struct {
	counts [StateHalfOpen + 1]atomic.Int64
	_      [32]byte
}

type breakerManager struct {
	opts []BreakerOption
	m    sync.Map[string, *breaker]
}

// New returns CircuitBreaker object if no error occurs.
func New(opts ...Option) (CircuitBreaker, error) {
	bm := &breakerManager{}
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(bm); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return bm, nil
}

type breakerStateError struct {
	err error
	st  State
}

func (e *breakerStateError) Error() string {
	msg := "circuitbreaker state is " + e.st.String() + ", this error is not caused by circuitbreaker"
	if e.err == nil {
		return msg
	}
	return msg + ": " + e.err.Error()
}

func (e *breakerStateError) Unwrap() error {
	return e.err
}

// Handle provides direct access to a resolved key's breaker state machine
// without repeating the manager's map lookup or boxing results in any.
type Handle struct {
	b *breaker
}

// Resolve returns a cached Handle for key, creating the breaker and metrics
// entry if needed.
func (bm *breakerManager) Resolve(key string) (*Handle, error) {
	br, err := bm.resolveBreaker(key)
	if err != nil {
		return nil, err
	}
	return &br.handle, nil
}

func (bm *breakerManager) resolveBreaker(key string) (*breaker, error) {
	// Pre-loading to prevent a lot of object generation on the hot path.
	if br, ok := bm.m.Load(key); ok {
		return br, nil
	}
	metricsMu.Lock()
	defer metricsMu.Unlock()
	if br, ok := bm.m.Load(key); ok {
		return br, nil
	}
	nb, err := newBreaker(key, bm.opts...)
	if err != nil {
		return nil, err
	}
	c, _ := metrics.LoadOrStore(key, new(stateCounts))
	nb.counts = c
	bm.m.Store(key, nb)
	return nb, nil
}

// IsEvicted reports whether the underlying breaker has been evicted via EvictAddr.
func (h *Handle) IsEvicted() bool {
	return h == nil || h.b == nil || h.b.evicted.Load()
}

func (h *Handle) recordCount(st State) {
	h.b.recordCount(st)
}

func wrapStateError(st State, err error) error {
	switch st {
	case StateClosed:
		return &breakerStateError{err: err, st: st}
	case StateOpen:
		if !errors.Is(err, errors.ErrCircuitBreakerOpenState) {
			return errors.Join(err, errors.ErrCircuitBreakerOpenState)
		}
	case StateHalfOpen:
		if !errors.Is(err, errors.ErrCircuitBreakerHalfOpenFlowLimitation) {
			return errors.Join(err, errors.ErrCircuitBreakerHalfOpenFlowLimitation)
		}
	case StateUnknown:
	}
	return err
}

type breakerErrKind uint8

const (
	breakerErrNone breakerErrKind = iota
	breakerErrMarkWithSuccess
	breakerErrIgnorable
)

// matchMarkWithSuccessAs and matchIgnorableAs are kept out of the hot-path
// frame so taking &serr/&igerr never causes heap escapes on errors that do not
// implement As(any) bool.
//
//go:noinline

func matchMarkWithSuccessAs(asErr interface{ As(any) bool }) (error, bool) {
	var serr *errors.ErrCircuitBreakerMarkWithSuccess
	if asErr.As(&serr) && serr != nil {
		return serr.Unwrap(), true
	}
	return nil, false
}

//go:noinline

func matchIgnorableAs(asErr interface{ As(any) bool }) (error, bool) {
	var igerr *errors.ErrCircuitBreakerIgnorable
	if asErr.As(&igerr) && igerr != nil {
		return igerr.Unwrap(), true
	}
	return nil, false
}

// classifyWrappedBreakerErr walks err without calling errors.As with interface
// pointer targets, preserving ErrCircuitBreakerMarkWithSuccess > ErrCircuitBreakerIgnorable
// precedence with zero heap allocations on standard/wrapped errors.
func classifyWrappedBreakerErr(err error) (kind breakerErrKind, unwrapped error) {
	var (
		ignUnwrapped error
		hasIgn       bool
	)
	for cur := err; cur != nil; {
		//nolint:errorlint // explicit unwrap walk avoids errors.As interface-target heap escapes
		switch e := cur.(type) {
		case *errors.ErrCircuitBreakerMarkWithSuccess:
			if e == nil {
				cur = nil
				continue
			}
			return breakerErrMarkWithSuccess, e.Unwrap()
		case *errors.ErrCircuitBreakerIgnorable:
			if e == nil {
				cur = nil
				continue
			}
			if !hasIgn {
				ignUnwrapped = e.Unwrap()
				hasIgn = true
			}
			cur = e.Unwrap()
			continue
		}
		if asErr, ok := cur.(interface{ As(any) bool }); ok {
			if u, ok := matchMarkWithSuccessAs(asErr); ok {
				return breakerErrMarkWithSuccess, u
			}
			if !hasIgn {
				if u, ok := matchIgnorableAs(asErr); ok {
					ignUnwrapped = u
					hasIgn = true
				}
			}
		}
		//nolint:errorlint // explicit unwrap walk handles Unwrap() error / []error without errors.As recursion
		switch e := cur.(type) {
		case interface{ Unwrap() error }:
			cur = e.Unwrap()
		case interface{ Unwrap() []error }:
			for _, sub := range e.Unwrap() {
				if sub == nil {
					continue
				}
				switch k, u := classifyWrappedBreakerErr(sub); k {
				case breakerErrMarkWithSuccess:
					return breakerErrMarkWithSuccess, u
				case breakerErrIgnorable:
					if !hasIgn {
						ignUnwrapped = u
						hasIgn = true
					}
				case breakerErrNone:
				}
			}
			if hasIgn {
				return breakerErrIgnorable, ignUnwrapped
			}
			return breakerErrNone, nil
		default:
			if hasIgn {
				return breakerErrIgnorable, ignUnwrapped
			}
			return breakerErrNone, nil
		}
	}
	if hasIgn {
		return breakerErrIgnorable, ignUnwrapped
	}
	return breakerErrNone, nil
}

// Allow checks whether a call is permitted under the breaker's current state.
// If rejected, it updates ignore/state counters and returns the wrapped breaker error.
func (h *Handle) Allow() error {
	if h == nil || h.b == nil {
		return nil
	}
	if st, err := h.b.isReady(); err != nil {
		h.b.count.onIgnore()
		h.recordCount(st)
		return wrapStateError(st, err)
	}
	return nil
}

// FinishSuccess records a successful call on the breaker and updates metrics.
func (h *Handle) FinishSuccess() {
	if h == nil || h.b == nil {
		return
	}
	h.recordCount(h.b.success())
}

func canWrapBreakerErr(err error) bool {
	//nolint:errorlint // fast capability check before recursive unwrap walk
	switch err.(type) {
	case interface{ Unwrap() error }, interface{ Unwrap() []error }, interface{ As(any) bool }:
		return true
	default:
		return false
	}
}

// finishError classifies an error for both (*breaker).do and (*Handle).FinishError,
// honoring wrapper/context sentinels and retryability while avoiding heap escapes.
func (b *breaker) finishError(err error, retryable bool) (State, error) {
	// Fast-path concrete wrapper types and direct sentinel errors before
	// unwrap tree / errors.Is checks.
	//nolint:errorlint // concrete fast path before unwrap walk below
	if serr, ok := err.(*errors.ErrCircuitBreakerMarkWithSuccess); ok && serr != nil {
		return b.success(), serr.Unwrap()
	}
	//nolint:errorlint // concrete fast path before unwrap walk below
	if igerr, ok := err.(*errors.ErrCircuitBreakerIgnorable); ok && igerr != nil {
		u := igerr.Unwrap()
		//nolint:errorlint // concrete fast path before unwrap walk below
		if serr, ok := u.(*errors.ErrCircuitBreakerMarkWithSuccess); ok && serr != nil {
			return b.success(), serr.Unwrap()
		}
		if canWrapBreakerErr(u) {
			if kind, mu := classifyWrappedBreakerErr(u); kind == breakerErrMarkWithSuccess {
				return b.success(), mu
			}
		}
		return b.fastState(), u
	}
	//nolint:errorlint // direct sentinel fast path before errors.Is fallback below
	if err == context.Canceled || err == context.DeadlineExceeded {
		return b.fastState(), err
	}
	switch kind, u := classifyWrappedBreakerErr(err); kind {
	case breakerErrMarkWithSuccess:
		return b.success(), u
	case breakerErrIgnorable:
		return b.fastState(), u
	case breakerErrNone:
	}
	if !retryable || errors.IsAny(err, context.Canceled, context.DeadlineExceeded) {
		return b.fastState(), err
	}
	b.fail()
	return b.currentState(), err
}

// RecordError records a failed call on the breaker and updates metrics without
// wrapping err in a *breakerStateError, for callers that discard the error.
func (h *Handle) RecordError(err error, retryable bool) {
	if err == nil {
		h.FinishSuccess()
		return
	}
	if h == nil || h.b == nil {
		return
	}
	st, _ := h.b.finishError(err, retryable)
	h.recordCount(st)
}

// FinishError records a failed call on the breaker, honoring wrapper/context
// sentinels and retryability, updates metrics, and returns the state-wrapped error.
func (h *Handle) FinishError(err error, retryable bool) error {
	if err == nil {
		h.FinishSuccess()
		return nil
	}
	if h == nil || h.b == nil {
		return err
	}
	st, uerr := h.b.finishError(err, retryable)
	h.recordCount(st)
	return wrapStateError(st, uerr)
}

// Do invokes the breaker matching the given key.
func (bm *breakerManager) Do(
	ctx context.Context, key string, fn func(ctx context.Context) (any, error),
) (val any, err error) {
	br, err := bm.resolveBreaker(key)
	if err != nil {
		return nil, err
	}

	var st State
	val, st, err = br.do(ctx, fn)
	br.recordCount(st)
	if err != nil {
		return val, wrapStateError(st, err)
	}
	return val, nil
}

// EvictAddr removes the breakers and metrics counters of every key that targets
// addr, so the state of removed addresses does not pile up. The gRPC client
// calls it when it disconnects addr. Its keys are "<addr>" or "<method>/<addr>"
// (grpc.WrapGRPCMethod). A later Do for an evicted key starts a new, closed
// breaker and new counters.
func (bm *breakerManager) EvictAddr(addr string) {
	if addr == "" {
		return
	}
	suffix := "/" + addr
	targets := func(key string) bool {
		return key == addr || strings.HasSuffix(key, suffix)
	}
	metricsMu.Lock()
	defer metricsMu.Unlock()
	bm.m.Range(func(key string, br *breaker) bool {
		if targets(key) {
			if br != nil {
				br.evicted.Store(true)
			}
			bm.m.Delete(key)
		}
		return true
	})
	metrics.Range(func(key string, _ *stateCounts) bool {
		if targets(key) {
			metrics.Delete(key)
		}
		return true
	})
}

func Metrics(context.Context) (ms map[string]map[State]int64) {
	metrics.Range(func(name string, c *stateCounts) bool {
		// Report only the states that occurred, like the previous map did.
		var sts map[State]int64
		for st := range c.counts {
			if n := c.counts[st].Load(); n > 0 {
				if sts == nil {
					sts = make(map[State]int64, len(c.counts))
				}
				sts[State(st)] = n
			}
		}
		if sts != nil {
			if ms == nil {
				ms = make(map[string]map[State]int64)
			}
			ms[name] = sts
		}
		return true
	})
	return ms
}
