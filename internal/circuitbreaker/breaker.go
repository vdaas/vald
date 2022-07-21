//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package circuitbreaker

import (
	"context"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

type breaker struct {
	key     string       // breaker key for logging
	count   atomic.Value // type: *count
	tripped int32        // tripped flag. when flag value is 1, breaker state is "Open" or "HalfOpen".

	closedErrRate         float32
	closedErrShouldTrip   Tripper
	halfOpenErrRate       float32
	halfOpenErrShouldTrip Tripper
	minSamples            int64
	openTimeout           time.Duration
	openExp               int64 // unix time
	cloedRefreshTimeout   time.Duration
	closedRefreshExp      int64 // unix time
}

func newBreaker(key string, opts ...BreakerOption) (*breaker, error) {
	b := &breaker{
		key: key,
	}
	for _, opt := range append(defaultBreakerOpts, opts...) {
		if err := opt(b); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	b.count.Store(&count{})

	if b.closedErrShouldTrip == nil {
		b.closedErrShouldTrip = NewRateTripper(b.closedErrRate, b.minSamples)
	}
	if b.halfOpenErrShouldTrip == nil {
		b.halfOpenErrShouldTrip = NewRateTripper(b.halfOpenErrRate, b.minSamples)
	}
	return b, nil
}

// do executes the function given argument when the current breaker state is "Closed" or "Half-Open".
// If the current breaker state is "Open", this function returns ErrCircuitBreakerOpenState.
func (b *breaker) do(ctx context.Context, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, st State, err error) {
	if !b.isReady() {
		return nil, StateOpen, errors.ErrCircuitBreakerOpenState
	}
	val, err = fn(ctx)
	if err != nil {
		serr := &errors.ErrCircuitBreakerMarkWithSuccess{}
		if errors.As(err, &serr) {
			b.success()
			return nil, b.currentState(), serr.Unwrap()
		}

		igerr := &errors.ErrCircuitBreakerIgnorable{}
		if errors.As(err, &igerr) {
			return nil, b.currentState(), igerr.Unwrap()
		}

		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return nil, b.currentState(), err
		}

		b.fail()
		return nil, b.currentState(), err
	}
	b.success()
	return val, b.currentState(), nil
}

// isReady determines the breaker is ready or not.
// If the current breaker state is "Closed" or "Half-Open", this function returns true.
func (b *breaker) isReady() (ok bool) {
	st := b.currentState()
	return st == StateClosed || st == StateHalfOpen
}

func (b *breaker) success() {
	b.count.Load().(*count).onSuccess()
	if st := b.currentState(); st == StateHalfOpen {
		log.Infof("the operation succeeded, circuit breaker state for '%s' changed,\tfrom: %s, to: %s", b.key, st.String(), StateClosed.String())
		b.reset()
	}
}

func (b *breaker) fail() {
	cnt := b.count.Load().(*count)
	cnt.onFail()

	var ok bool
	var st State
	switch st = b.currentState(); st {
	case StateHalfOpen:
		ok = b.halfOpenErrShouldTrip.ShouldTrip(cnt)
	case StateClosed:
		ok = b.closedErrShouldTrip.ShouldTrip(cnt)
	default:
		return
	}
	if ok {
		log.Infof("the operation failed, circuit breaker state for '%s' changed,\nfrom: %s, to: %s", b.key, st.String(), StateOpen.String())
		b.trip()
	}
}

// currentState returns current breaker state.
// If the tripped flag is not active, this function returns "Closed" state.
// On the other hand, if the tripped flag is active and the expiration date is not reached, returns "Open" state, otherwise returns "Half-Open" state.
func (b *breaker) currentState() State {
	now := time.Now().UnixNano()
	if b.isTripped() {
		if expire := atomic.LoadInt64(&b.openExp); expire > 0 && expire > now {
			return StateOpen
		}
		return StateHalfOpen
	}
	if expire := atomic.LoadInt64(&b.closedRefreshExp); expire == 0 || now > expire {
		log.Infof("the closed state expired, circuit breaker state for '%s' refleshed,\nto: %s", b.key, StateClosed.String())
		b.reset()
	}
	return StateClosed
}

func (b *breaker) reset() {
	atomic.StoreInt32(&b.tripped, 0)
	atomic.StoreInt64(&b.openExp, 0)
	atomic.StoreInt64(&b.closedRefreshExp, time.Now().Add(b.cloedRefreshTimeout).UnixNano())
	b.count.Load().(*count).reset()
}

func (b *breaker) trip() {
	atomic.StoreInt32(&b.tripped, 1)
	atomic.StoreInt64(&b.openExp, time.Now().Add(b.openTimeout).UnixNano())
	atomic.StoreInt64(&b.closedRefreshExp, 0)
	b.count.Load().(*count).reset()
}

func (b *breaker) isTripped() (ok bool) {
	return atomic.LoadInt32(&b.tripped) == 1
}
