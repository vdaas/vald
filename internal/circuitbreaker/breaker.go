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
	count   atomic.Value // type: *count
	tripped int32        // breaker is active or not.

	closedErrRate float64
	openTimeout   time.Duration
	openExpire    int64 // Unix time
}

func newBreaker(opts ...BreakerOption) (*breaker, error) {
	b := &breaker{}
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
	return b, nil
}

// do executes the function given argument when the current breaker state is "Closed" or "Half-Open".
// If the current breaker state is "Open", this function returns ErrCircuitBreakerOpenState.
func (b *breaker) do(ctx context.Context, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, err error) {
	if !b.ready() {
		return nil, errors.ErrCircuitBreakerOpenState
	}

	// count up

	val, err = fn(ctx)
	if err != nil {
		b.fail()
	}

	b.success()
	return val, nil
}

// ready determines the breaker is ready or not.
// If the current breaker state is "Closed" or "Half-Open", this function returns true.
func (b *breaker) ready() (ok bool) {
	st := b.currentState()
	return st == stateClosed || st == stateHalfOpen
}

func (b *breaker) success() {
	switch st := b.currentState(); st {
	// In most cases, an "Open" state will not occur, but reset in that case to be safe.
	case stateHalfOpen, stateOpen:
		b.reset()
		return
	}

	b.count.Load().(*counts).onSuccess()
	// WIP:
	// This function focus on the "Half-Open" state.
	// When current state is "Half-Open" state, this function records the number of successful attempts.
	// And the circuit breaker changes to the "Closed" state after a specified number of consecutive operation invocations have been successful.
}

func (b *breaker) fail() {
	b.count.Load().(*counts).onFailure()

	// WIP:
	// This function focus on the "Closed" and "Half-Open" state.
	// When current state is "Closed" state, this function records the number of failuer attempts.
	// And the circuit breaker changes to the "Open" state if a specified number of consecutive operation invocations have been failed.
	// When current state is "Half-Open" state, the circuit breaker changes to the "Open" state.

	// 1. Increment the consecutiveFailures and clear consecutiveSuccesses by using count.onFailuer function
	// 2. Call b.trip function

	// e.g. the following is an example flow.
	/**
		switch b.currentState() {
		case stateClosed:
		    cnt := b.count.Load().(*counts).onFailure()

		    if cnt >= b.halfClosedFailureThreshold {
			    b.trip() // To "Open"
		    }
		case stateHalfOpen:
	        if b.st.CompareAndSwap(stateHalfOpen, stateOpen) {
		        b.reset(time.Now()) // Reset all counter
		    } else {
		        // The state has been changed to "Closed" during this function calling.
		    }
		}
	**/
}

// currentState returns current breaker state.
// If the tripped flag is not active, this function returns "Closed" state.
// On the other hand, if the tripped flag is active and the expiration date is not reached, returns "Open" state, otherwise returns "Half-Open" state.
func (b *breaker) currentState() (st state) {
	if b.isTripped() {
		now := time.Now().UnixNano()
		if expire := atomic.LoadInt64(&b.openExpire); expire > 0 && expire >= now {
			return stateOpen
		}
		return stateHalfOpen
	}
	return stateClosed
}

func (b *breaker) reset() {
	atomic.StoreInt32(&b.tripped, 0)
	atomic.StoreInt64(&b.openExpire, 0)
}

func (b *breaker) isTripped() (ok bool) {
	return atomic.LoadInt32(&b.tripped) == 1
}

func (b *breaker) trip() {
	atomic.StoreInt32(&b.tripped, 1)
	atomic.StoreInt64(&b.openExpire, time.Now().Add(b.openTimeout).UnixNano())
}

func (b *breaker) unTrip() {
	atomic.StoreInt32(&b.tripped, 0)
}
