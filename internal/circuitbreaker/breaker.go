package circuitbreaker

import (
	"context"
	"sync/atomic"
	"time"
)

type counts struct {
	consecutiveSuccesses uint32
	consecutiveFailures  uint32
}

func (c *counts) onSuccess() {
	// 1. Increment consecutiveSuccesses by using atomic.
	// 2. Clear consecutiveFailures by using atomic.
}
func (c *counts) onFailure() {
	// 1. Increment consecutiveFailures by using atomic.
	// 2. Clear consecutiveFailures by using atomic.
}
func (c *counts) clear() {
	// 1. Clear consecutiveSuccesses and consecutiveFailures by using atomic
}

type breaker struct {
	st         atomic.Value // type: state
	count      atomic.Value // type: *count
	lastFilure int32
	tripped    int32

	halfOpenTimeout int64
	expire          time.Time
}

func (b *breaker) do(ctx context.Context, fn func(ctx context.Context) (val interface{}, err error)) (val interface{}, err error) {
	// This function is a main flow of breaker.
	// When The current state is "Closed" or "Half-Open", this function executes the given user defined function variable.
	// If no error occurs, above function is treated as successful.

	// 1. Check current state. If current state is open, this function returns an error.

	// 2. Execute the given user defined function variable fn.

	// 3. Get an error of above function calling.

	// 4. Call success or fail function based on above error.
	//    - if error is nil or error is ErrCircuitBreakerIgnorable, calls success function.
	//    - if error is not nil, call fail function.

	// e.g. the following is an example flow code.
	/**
		if !b.ready() {
		    return nil, errors.ErrCircuitBreakerOpenState
		}

		val, err = fn(ctx)
		if err != nil {
		    if !errors.Is(err, ErrCircuitBreakerIgnorable) {
		        b.fail() <------- mark with fail
		        return nil, err.Unwrap()
		    }
		}

		b.success() <-------- mark with success
		return val, nil
	**/
	return nil, nil
}

func (b *breaker) ready() (ok bool) {
	// This function determines breaker is ready or not to call user defined function.

	// 1. If the circuit breaker state is "Closed" or "Half-Open", this function returns true.

	// e.g. the following is an example flow code.
	/**
	    st := b.currentState()
	    return st == stateClosed || st == stateHalfOpen
	**/
	return true
}

func (b *breaker) success() {
	// This function focus on the "Half-Open" state.
	// When current state is "Half-Open" state, this function records the number of successful attempts.
	// And the circuit breaker changes to the "Closed" state after a specified number of consecutive operation invocations have been successful.

	// 1. Increment the consecutiveSuccesses and clear consecutiveFailures by using count.onSuccess function
	// 2. Call b.unTrip function

	// e.g. the following is an example flow code.
	/**
		if st := b.currentState(); st == stateHalfOpen {
	        cnt := b.count.Load().(*counts).onSuccess()
		    if cnt >= b.halfOpenSuccessThreshold {
			    b.unTrip()
		    }
		}
	**/
}

func (b *breaker) fail() {
	// This function focus on the "Closed" and "Half-Open" state.
	// When current state is "Closed" state, this function records the number of failuer attempts.

	// 1. Increment the consecutiveFailures and clear consecutiveSuccesses by using count.onFailuer function
	// 2. Call b.trip function

	// e.g. the following is an example flow code.
	/**
		if st := b.currentState(); st == stateClosed {
			cnt := b.count.Load().(*counts).onFailure()
			if cnt >= b.halfClosedFailureThreshold {
				b.trip()
			}
		}
	**/
}

func (b *breaker) currentState() (st state) {
	// 1. Returns current state
	// if atomic.LoadInt32(&b.tripped) == 1 {
	// 	now := time.Now().UnixNano()
	// 	if now-int64(b.lastFilure) > b.halfOpenTimeout {
	// 		if b.st.CompareAndSwap(stateOpen, stateHalfOpen) {
	// 			b.reset()
	// 		}
	// 		return stateHalfOpen
	// 	}
	// 	return stateOpen // or return stateOpen
	// }
	return stateClosed
}

func (b *breaker) reset() {}

func (b *breaker) trip() (ok bool) {
	// st := b.currentState()
	// if st == stateClosed {
	// 	if b.st.CompareAndSwap(stateClosed, stateOpen) {
	// 		atomic.StoreInt32(&b.tripped, 1)
	// 	}
	// 	return true
	// }
	return false
}

func (b *breaker) unTrip() (ok bool) {
	return ok
}
