package circuitbreaker

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

type breaker struct {
	count   atomic.Value // type: *count
	tripped int32

	halfOpenTimeout time.Duration
	openTimeout     time.Duration
	expire          int64
}

func newBreaker(opts ...BreakerOption) (*breaker, error) {
	b := &breaker{}
	for _, opt := range append(defaultBreakerOpts, opts...) {
		if err := opt(b); err != nil {
			return nil, err
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
	// This function focus on the "Half-Open" state.
	// When current state is "Half-Open" state, this function records the number of successful attempts.
	// And the circuit breaker changes to the "Closed" state after a specified number of consecutive operation invocations have been successful.

	// 1. Increment the consecutiveSuccesses and clear consecutiveFailures by using count.onSuccess function
	// 2. Call b.unTrip function

	// e.g. the following is an example flow.
	/**
		switch b.currentState() {
	    case stateHalfOpen:
		    cnt := b.count.Load().(*counts).onSuccess()

		    if cnt >= b.halfOpenSuccessThreshold {
		        b.unTrip() // To "Closed"
		    }
		}
	**/
}

func (b *breaker) fail() {
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
		if expire := atomic.LoadInt64(&b.expire); expire > 0 && now >= expire {
			return stateOpen
		}
		return stateHalfOpen
	}
	return stateClosed
}

func (b *breaker) reset(now time.Time) {
	// This function resets all counter and expire time.
	/**
		b.count.(*count).reset()
	    b.expire.Store(now)
	**/
}

func (b *breaker) isTripped() (ok bool) {
	return atomic.LoadInt32(&b.tripped) == 1
}

func (b *breaker) trip() {
	atomic.StoreInt32(&b.tripped, 1)
}

func (b *breaker) unTrip() {
	atomic.StoreInt32(&b.tripped, 0)
}
